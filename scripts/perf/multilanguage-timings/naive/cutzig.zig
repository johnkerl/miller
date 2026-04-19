// CLI: step (1-6) arg1, comma-separated field names arg2, filenames arg3+.
// Step: 1=read, 2=+parse, 3=+select, 4=+build, 5=+newline, 6=+write.
const std = @import("std");

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

    var buf: [256 * 1024]u8 = undefined;
    var line_buf = try std.ArrayList(u8).initCapacity(allocator, 4096);
    defer line_buf.deinit(allocator);

    while (true) {
        const n = file.read(&buf) catch |err| {
            std.debug.print("read: {}\n", .{err});
            return false;
        };
        if (n == 0) break;
        var i: usize = 0;
        while (i < n) {
            if (buf[i] == '\n') {
                i += 1;
                const line = line_buf.items;
                line_buf.clearRetainingCapacity();

                if (step <= 1) continue;

                // Step 2: line to map
                var mymap = std.StringHashMap([]const u8).init(allocator);
                defer {
                    var it = mymap.iterator();
                    while (it.next()) |e| {
                        allocator.free(e.key_ptr.*);
                        allocator.free(e.value_ptr.*);
                    }
                    mymap.deinit();
                }
                var iter = std.mem.splitScalar(u8, line, ',');
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

                // Step 3: map-to-map transform (output in include_fields order)
                var newmap = std.StringHashMap([]const u8).init(allocator);
                defer {
                    var it = newmap.iterator();
                    while (it.next()) |e| {
                        allocator.free(e.value_ptr.*);
                    }
                    newmap.deinit();
                }
                for (include_fields) |inc_k| {
                    if (mymap.get(inc_k)) |val| {
                        const val_dup = allocator.dupe(u8, val) catch continue;
                        newmap.put(inc_k, val_dup) catch allocator.free(val_dup);
                    }
                }
                if (step <= 3) continue;

                // Step 4-5: map to string + newline
                var out = try std.ArrayList(u8).initCapacity(allocator, 256);
                defer out.deinit(allocator);
                var first = true;
                for (include_fields) |k| {
                    if (newmap.get(k)) |v| {
                        if (!first) try out.appendSlice(allocator, ",");
                        try out.appendSlice(allocator, k);
                        try out.appendSlice(allocator, "=");
                        try out.appendSlice(allocator, v);
                        first = false;
                    }
                }
                try out.append(allocator, '\n');
                if (step <= 5) continue;

                // Step 6: write to stdout
                const stdout_file = std.fs.File.stdout();
                stdout_file.writeAll(out.items) catch |err| {
                    std.debug.print("write: {}\n", .{err});
                    return false;
                };
            } else {
                line_buf.append(allocator, buf[i]) catch return false;
                i += 1;
            }
        }
    }

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
