# CLI: step (1-6) arg1, comma-separated field names arg2, filenames arg3+.
# Step: 1=read, 2=+parse, 3=+select, 4=+build, 5=+newline, 6=+write.
import strutils, tables, os, streams

proc handle(fileName: string, step: int, includeFields: seq[string]): bool =
  var inputStream: Stream
  if fileName == "-":
    inputStream = newFileStream(stdin)
  else:
    try:
      inputStream = newFileStream(fileName)
    except OSError as e:
      stderr.writeLine(e.msg)
      return false

  defer:
    if fileName != "-":
      inputStream.close()

  var line: string
  while inputStream.readLine(line):
    if step <= 1:
      continue

    # Step 2: line to map
    var oldmap = initOrderedTable[string, string]()
    for word in line.split(','):
      let pair = word.split('=', 2)
      if pair.len >= 2:
        oldmap[pair[0]] = pair[1]
    if step <= 2:
      continue

    # Step 3: map-to-map transform
    var newmap = initOrderedTable[string, string]()
    for k in includeFields:
      if k in oldmap:
        newmap[k] = oldmap[k]
    if step <= 3:
      continue

    # Step 4-5: map to string + newline
    var parts: seq[string]
    for k, v in newmap:
      parts.add(k & "=" & v)
    let outLine = parts.join(",") & "\n"
    if step <= 5:
      continue

    # Step 6: write to stdout
    stdout.write(outLine)

  return true

proc main =
  if paramCount() < 2:
    stderr.writeLine("usage: ", getAppFilename(), " <step 1-6> <field1,field2,...> [file ...]")
    quit(1)
  let stepStr = paramStr(1)
  let step = parseInt(stepStr)
  if step < 1 or step > 6:
    stderr.writeLine("step must be 1-6, got ", stepStr)
    quit(1)
  let includeFields = paramStr(2).split(',')
  var filenames: seq[string]
  for i in 3..paramCount():
    filenames.add(paramStr(i))
  if filenames.len == 0:
    filenames = @["-"]

  var ok = true
  for arg in filenames:
    if not handle(arg, step, includeFields):
      ok = false
  quit(if ok: 0 else: 1)

main()
