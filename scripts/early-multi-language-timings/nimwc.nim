import strutils

var word_count = 0
for line in stdin.lines:
  for word in line.split(","):
      word_count += 1
echo (word_count)
