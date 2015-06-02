// Reads $(D stdin) and writes it to $(D stdout).
// http://dlang.org/hash-map.html
import std.stdio;
import std.string;
import std.array;

void main() {
	string[] includeFields = ["a", "x"];
	string line;
	while ((line = stdin.readln()) !is null) {
		// Input string to hashmap.
		string[string] oldmap;
		string[] fields = split(line, ',');
		foreach (field; fields) {
			string[] kvps = split(field, '='); // really want splitN with max #parts = 2
			oldmap[kvps[0]] = kvps[1];
		}

		// Hashmap-to-hashmap transform.
		// Note: unordered hashmap here.
		string[string] newmap;
		foreach (includeField; includeFields) {
			if (includeField in oldmap) {
				newmap[includeField] = oldmap[includeField];
			}
		}

		// Hashmap to output strings.
		int i = 0;
		foreach (key; newmap.keys) {
			if (i > 0)
				write(',');
			write(key);
			write('=');
			write(newmap[key]);
			i++;
		}
		write('\n');
	}
}
