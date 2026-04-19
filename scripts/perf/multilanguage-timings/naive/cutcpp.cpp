// CLI: step (1-6) arg1, comma-separated field names arg2, filenames arg3+.
// Step: 1=read, 2=+parse, 3=+select, 4=+build, 5=+newline, 6=+write.
#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <unordered_map>
#include <vector>
#include <cstdlib>

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
	if (fileName != "-") {
		fileStream.open(fileName);
		if (!fileStream) {
			std::cerr << "open: " << fileName << "\n";
			return false;
		}
		in = &fileStream;
	}

	std::string line;
	while (std::getline(*in, line)) {
		if (step <= 1) continue;

		// Step 2: line to map
		std::unordered_map<std::string, std::string> mymap;
		for (const auto& field : split(line, ',')) {
			std::string k, v;
			if (splitKeyValue(field, k, v))
				mymap[k] = v;
		}
		if (step <= 2) continue;

		// Step 3: map-to-map transform
		std::unordered_map<std::string, std::string> newmap;
		for (const auto& includeField : includeFields) {
			auto it = mymap.find(includeField);
			if (it != mymap.end())
				newmap[it->first] = it->second;
		}
		if (step <= 3) continue;

		// Step 4-5: map to string + newline (emit in includeFields order)
		std::ostringstream buf;
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

		// Step 6: write to stdout
		std::cout << buf.str();
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
