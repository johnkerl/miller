// CLI: step (1-6) arg1, comma-separated field names arg2, filenames arg3+.
// Step: 1=read, 2=+parse, 3=+select, 4=+build, 5=+newline, 6=+write.
const std = @import("std");

const pipeline_cap = 64;

const LineJob = struct { index: usize, line: []const u8 };
const OutJob = struct { index: usize, out: []const u8 };

fn BoundedQueue(comptime T: type) type {
    return struct {
        list: std.ArrayList(T),
        allocator: std.mem.Allocator,
        cap: usize,
        mutex: std.Thread.Mutex = .{},
        condition: std.Thread.Condition = .{},
        closed: bool = false,

        const Self = @This();
        fn init(allocator: std.mem.Allocator, capacity: usize) !Self {
            return .{
                .list = try std.ArrayList(T).initCapacity(allocator, capacity),
                .allocator = allocator,
                .cap = capacity,
            };
        }
        fn deinit(self: *Self) void {
            self.list.deinit(self.allocator);
        }
        fn push(self: *Self, item: T) bool {
            self.mutex.lock();
            defer self.mutex.unlock();
            while (self.list.items.len >= self.cap and !self.closed) {
                self.condition.wait(&self.mutex);
            }
            if (self.closed) return false;
            self.list.append(self.allocator, item) catch return false;
            self.condition.signal();
            return true;
        }
        fn pop(self: *Self) ?T {
            self.mutex.lock();
            defer self.mutex.unlock();
            while (self.list.items.len == 0 and !self.closed) {
                self.condition.wait(&self.mutex);
            }
            if (self.list.items.len == 0) return null;
            const item = self.list.orderedRemove(0);
            self.condition.signal();
            return item;
        }
        fn close(self: *Self) void {
            self.mutex.lock();
            defer self.mutex.unlock();
            self.closed = true;
            self.condition.broadcast();
        }
    };
}

fn readerRun(allocator: std.mem.Allocator, file: std.fs.File, read_queue: *BoundedQueue(LineJob)) void {
    var buf: [256 * 1024]u8 = undefined;
    var line_buf = std.ArrayList(u8).initCapacity(allocator, 4096) catch return;
    defer line_buf.deinit(allocator);
    var index: usize = 0;
    while (true) {
        const n = file.read(&buf) catch break;
        if (n == 0) break;
        var i: usize = 0;
        while (i < n) {
            if (buf[i] == '\n') {
                i += 1;
                const line = allocator.dupe(u8, line_buf.items) catch break;
                _ = read_queue.push(.{ .index = index, .line = line });
                index += 1;
                line_buf.clearRetainingCapacity();
            } else {
                line_buf.append(allocator, buf[i]) catch break;
                i += 1;
            }
        }
    }
    read_queue.close();
}

fn processorRun(
    allocator: std.mem.Allocator,
    step: u8,
    include_fields: []const []const u8,
    read_queue: *BoundedQueue(LineJob),
    write_queue: *BoundedQueue(OutJob),
) void {
    var mymap = std.StringHashMap([]const u8).init(allocator);
    defer {
        var it = mymap.iterator();
        while (it.next()) |e| {
            allocator.free(e.key_ptr.*);
            allocator.free(e.value_ptr.*);
        }
        mymap.deinit();
    }
    var newmap = std.StringHashMap([]const u8).init(allocator);
    defer {
        var it = newmap.iterator();
        while (it.next()) |e| {
            allocator.free(e.value_ptr.*);
        }
        newmap.deinit();
    }
    var out = std.ArrayList(u8).initCapacity(allocator, 256) catch return;
    defer out.deinit(allocator);

    while (read_queue.pop()) |job| {
        defer allocator.free(job.line);
        if (step <= 1) continue;
        {
            var it = mymap.iterator();
            while (it.next()) |e| {
                allocator.free(e.key_ptr.*);
                allocator.free(e.value_ptr.*);
            }
            mymap.clearRetainingCapacity();
        }
        var iter = std.mem.splitScalar(u8, job.line, ',');
        while (iter.next()) |field| {
            var kv_iter = std.mem.splitScalar(u8, field, '=');
            const k = kv_iter.next() orelse continue;
            const v = kv_iter.rest();
            const k_dup = allocator.dupe(u8, k) catch continue;
            const v_dup = allocator.dupe(u8, v) catch {
                allocator.free(k_dup);
                continue;
            };
            if (mymap.getPtr(k_dup)) |existing| {
                allocator.free(existing.*);
            }
            mymap.put(k_dup, v_dup) catch {
                allocator.free(k_dup);
                allocator.free(v_dup);
            };
        }
        if (step <= 2) continue;
        {
            var it = newmap.iterator();
            while (it.next()) |e| {
                allocator.free(e.value_ptr.*);
            }
            newmap.clearRetainingCapacity();
        }
        for (include_fields) |inc_k| {
            if (mymap.get(inc_k)) |val| {
                const val_dup = allocator.dupe(u8, val) catch continue;
                newmap.put(inc_k, val_dup) catch allocator.free(val_dup);
            }
        }
        if (step <= 3) continue;
        out.clearRetainingCapacity();
        var first = true;
        for (include_fields) |k| {
            if (newmap.get(k)) |v| {
                if (!first) out.appendSlice(allocator, ",") catch {};
                out.appendSlice(allocator, k) catch {};
                out.appendSlice(allocator, "=") catch {};
                out.appendSlice(allocator, v) catch {};
                first = false;
            }
        }
        out.append(allocator, '\n') catch {};
        if (step <= 5) continue;
        const out_slice = allocator.dupe(u8, out.items) catch continue;
        _ = write_queue.push(.{ .index = job.index, .out = out_slice });
    }
    write_queue.close();
}

fn writerRun(allocator: std.mem.Allocator, write_queue: *BoundedQueue(OutJob)) void {
    const stdout_file = std.fs.File.stdout();
    while (write_queue.pop()) |job| {
        defer allocator.free(job.out);
        stdout_file.writeAll(job.out) catch return;
    }
}

fn handle(
    allocator: std.mem.Allocator,
    filename: []const u8,
    step: u8,
    include_fields: []const []const u8,
) !bool {
    var file = if (std.mem.eql(u8, filename, "-"))
        std.fs.File.stdin()
    else
        std.fs.cwd().openFile(filename, .{}) catch |err| {
            std.debug.print("open {s}: {}\n", .{ filename, err });
            return false;
        };
    defer if (!std.mem.eql(u8, filename, "-")) file.close();

    var read_queue = try BoundedQueue(LineJob).init(allocator, pipeline_cap);
    defer read_queue.deinit();
    var write_queue = try BoundedQueue(OutJob).init(allocator, pipeline_cap);
    defer write_queue.deinit();

    const h_reader = try std.Thread.spawn(.{}, readerRun, .{ allocator, file, &read_queue });
    const h_processor = try std.Thread.spawn(.{}, processorRun, .{
        allocator,
        step,
        include_fields,
        &read_queue,
        &write_queue,
    });
    const h_writer = try std.Thread.spawn(.{}, writerRun, .{ allocator, &write_queue });

    h_reader.join();
    h_processor.join();
    h_writer.join();

    return true;
}

pub fn main() !void {
    // GPA is slow (debug allocator); page_allocator is ~10–20x faster for allocation-heavy workloads.
    const allocator = std.heap.page_allocator;

    var args = try std.process.argsAlloc(allocator);
    defer std.process.argsFree(allocator, args);

    if (args.len < 3) {
        std.debug.print("usage: {s} <step 1-6> <field1,field2,...> [file ...]\n", .{args[0]});
        std.process.exit(1);
    }

    const step_n = std.fmt.parseInt(u8, args[1], 10) catch {
        std.debug.print("step must be 1-6, got {s}\n", .{args[1]});
        std.process.exit(1);
    };
    if (step_n < 1 or step_n > 6) {
        std.debug.print("step must be 1-6, got {s}\n", .{args[1]});
        std.process.exit(1);
    }

    var include_fields = try std.ArrayList([]const u8).initCapacity(allocator, 8);
    defer include_fields.deinit();
    var field_iter = std.mem.splitScalar(u8, args[2], ',');
    while (field_iter.next()) |f| {
        if (f.len > 0) include_fields.append(allocator, f) catch {};
    }

    var ok = true;
    if (args.len > 3) {
        for (args[3..]) |arg| {
            const result = handle(allocator, arg, step_n, include_fields.items) catch false;
            if (!result) ok = false;
        }
    } else {
        const result = handle(allocator, "-", step_n, include_fields.items) catch false;
        if (!result) ok = false;
    }
    std.process.exit(if (ok) 0 else 1);
}
