# Cut benchmarks: same CLI for all (step 1-6, comma-separated fields, filenames).
# Usage: ./cutgo 6 a,b,c [file ...]   # step 6 = full pipeline

binaries: cutgo cutcpp cutnim cutzig cutrs

cutgo: cutgo.go
	go build -o cutgo cutgo.go

cutcpp: cutcpp.cpp
	$(CXX) $(CXXFLAGS) -O2 -o cutcpp cutcpp.cpp

cutnim: cutnim.nim
	nim c -d:release -o:cutnim cutnim.nim

cutzig: cutzig.zig
	zig build-exe -O ReleaseFast -femit-bin=cutzig cutzig.zig

cutrs: cutrs.rs
	rustc -C opt-level=3 -o cutrs cutrs.rs

clean:
	rm -f cutgo cutcpp cutnim cutzig cutrs cutd.o

.PHONY: all clean
