// CLI: step (1-6) arg1, comma-separated field names arg2, filenames arg3+.
// Step: 1=read, 2=+parse, 3=+select, 4=+build, 5=+newline, 6=+write.
#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <unordered_map>
#include <vector>
#include <queue>
#include <thread>
#include <mutex>
#include <condition_variable>
#include <atomic>
#include <cstdlib>

static const size_t PIPELINE_CAP = 64;

template<typename T>
struct BoundedQueue {
	std::queue<T> q;
	std::mutex m;
	std::condition_variable not_full;
	std::condition_variable not_empty;
	size_t cap;
	bool closed = false;

	explicit BoundedQueue(size_t capacity) : cap(capacity) {}

	bool push(T item) {
		std::unique_lock<std::mutex> lock(m);
		not_full.wait(lock, [this] { return q.size() < cap || closed; });
		if (closed) return false;
		q.push(std::move(item));
		not_empty.notify_one();
		return true;
	}

	bool pop(T& out) {
		std::unique_lock<std::mutex> lock(m);
		not_empty.wait(lock, [this] { return !q.empty() || closed; });
		if (q.empty()) return false;
		out = std::move(q.front());
		q.pop();
		not_full.notify_one();
		return true;
	}

	void close() {
		std::lock_guard<std::mutex> lock(m);
		closed = true;
		not_full.notify_all();
		not_empty.notify_all();
	}
};

static void usage(const char* prog) {
	std::cerr << "usage: " << prog << " <step 1-6> <field1,field2,...> [file ...]\n";
	std::exit(1);
}

static std::vector<std::string> split(const std::string& s, char delim) {
	std::vector<std::string> out;
	std::istringstream iss(s);
	std::string part;
	while (std::getline(iss, part, delim)) {
		out.push_back(part);
	}
	return out;
}

static bool splitKeyValue(const std::string& field, std::string& key, std::string& value) {
	size_t pos = field.find('=');
	if (pos == std::string::npos) return false;
	key = field.substr(0, pos);
	value = field.substr(pos + 1);
	return true;
}

static bool handle(const std::string& fileName, int step,
	const std::vector<std::string>& includeFields) {
	std::istream* in = &std::cin;
	std::ifstream fileStream;
	static char readBuf[256 * 1024];
	if (fileName != "-") {
		fileStream.open(fileName);
		if (!fileStream) {
			std::cerr << "open: " << fileName << "\n";
			return false;
		}
		fileStream.rdbuf()->pubsetbuf(readBuf, sizeof readBuf);
		in = &fileStream;
	}

	using Job = std::pair<size_t, std::string>;
	BoundedQueue<Job> readQueue(PIPELINE_CAP);
	BoundedQueue<Job> writeQueue(PIPELINE_CAP);
	std::mutex err_mutex;
	std::string err_msg;
	std::atomic<bool> has_error{false};

	std::thread reader([in, &readQueue]() {
		std::string line;
		line.reserve(4096);
		size_t index = 0;
		while (std::getline(*in, line)) {
			if (!readQueue.push({index, line})) break;
			index++;
		}
		readQueue.close();
	});

	std::thread processor([step, &includeFields, &readQueue, &writeQueue]() {
		std::unordered_map<std::string, std::string> mymap;
		std::unordered_map<std::string, std::string> newmap;
		newmap.reserve(includeFields.size());
		std::ostringstream buf;
		Job job;
		while (readQueue.pop(job)) {
			if (step <= 1) continue;
			mymap.clear();
			for (const auto& field : split(job.second, ',')) {
				std::string k, v;
				if (splitKeyValue(field, k, v))
					mymap[k] = v;
			}
			if (step <= 2) continue;
			newmap.clear();
			for (const auto& includeField : includeFields) {
				auto it = mymap.find(includeField);
				if (it != mymap.end())
					newmap[it->first] = it->second;
			}
			if (step <= 3) continue;
			buf.str("");
			buf.clear();
			bool first = true;
			for (const auto& k : includeFields) {
				auto it = newmap.find(k);
				if (it != newmap.end()) {
					if (!first) buf << ',';
					buf << it->first << '=' << it->second;
					first = false;
				}
			}
			buf << '\n';
			if (step <= 5) continue;
			writeQueue.push({job.first, buf.str()});
		}
		writeQueue.close();
	});

	std::thread writer([&writeQueue, &err_mutex, &err_msg, &has_error]() {
		Job job;
		while (writeQueue.pop(job)) {
			std::cout << job.second;
			if (!std::cout) {
				std::lock_guard<std::mutex> lock(err_mutex);
				err_msg = "write error";
				has_error.store(true);
				return;
			}
		}
	});

	reader.join();
	processor.join();
	writer.join();

	if (has_error.load()) {
		std::lock_guard<std::mutex> lock(err_mutex);
		if (!err_msg.empty())
			std::cerr << err_msg << "\n";
		return false;
	}
	return true;
}

int main(int argc, char* argv[]) {
	if (argc < 3) usage(argv[0]);

	int step = std::atoi(argv[1]);
	if (step < 1 || step > 6) {
		std::cerr << "step must be 1-6, got " << argv[1] << "\n";
		std::exit(1);
	}

	std::vector<std::string> includeFields = split(argv[2], ',');
	std::vector<std::string> filenames;
	for (int i = 3; i < argc; i++)
		filenames.push_back(argv[i]);
	if (filenames.empty())
		filenames.push_back("-");

	bool ok = true;
	for (const auto& arg : filenames) {
		if (!handle(arg, step, includeFields)) ok = false;
	}
	return ok ? 0 : 1;
}
