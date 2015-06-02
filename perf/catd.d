// Reads $(D stdin) and writes it to $(D stdout).
import std.stdio;

void main()
{
	string line;
	while ((line = stdin.readln()) !is null)
		write(line);
}
