# CLI: step (1-6) arg1, comma-separated field names arg2, filenames arg3+.
# Step: 1=read, 2=+parse, 3=+select, 4=+build, 5=+newline, 6=+write.
import strutils, tables, os, streams, std/typedthreads

const pipelineCap = 64

type
  LineJob = tuple[index: int, line: string]
  OutJob = tuple[index: int, data: string]
  PipelineContext = ref object
    readChan: Channel[LineJob]
    writeChan: Channel[OutJob]
    step: int
    includeFields: seq[string]
    inputStream: Stream

proc readerProc(ctx: PipelineContext) {.thread.} =
  var index = 0
  var line: string
  while ctx.inputStream.readLine(line):
    ctx.readChan.send((index, line))
    index += 1
  ctx.readChan.send((-1, ""))

proc processorProc(ctx: PipelineContext) {.thread.} =
  var oldmap = initOrderedTable[string, string]()
  var newmap = initOrderedTable[string, string]()
  while true:
    let job = ctx.readChan.recv()
    if job.index < 0:
      break
    if ctx.step <= 1:
      continue
    oldmap.clear()
    for word in job.line.split(','):
      let pair = word.split('=', 2)
      if pair.len >= 2:
        oldmap[pair[0]] = pair[1]
    if ctx.step <= 2:
      continue
    newmap.clear()
    for k in ctx.includeFields:
      if k in oldmap:
        newmap[k] = oldmap[k]
    if ctx.step <= 3:
      continue
    var outLine: string
    for k in ctx.includeFields:
      if k in newmap:
        if outLine.len > 0:
          outLine.add(',')
        outLine.add(k)
        outLine.add('=')
        outLine.add(newmap[k])
    outLine.add('\n')
    if ctx.step <= 5:
      continue
    ctx.writeChan.send((job.index, outLine))
  ctx.writeChan.send((-1, ""))

proc writerProc(ctx: PipelineContext) {.thread.} =
  while true:
    let job = ctx.writeChan.recv()
    if job.index < 0:
      break
    stdout.write(job.data)

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

  var ctx: PipelineContext
  new ctx
  ctx.step = step
  ctx.includeFields = includeFields
  ctx.inputStream = inputStream
  ctx.readChan.open(pipelineCap)
  ctx.writeChan.open(pipelineCap)

  var readerThread: Thread[PipelineContext]
  var processorThread: Thread[PipelineContext]
  var writerThread: Thread[PipelineContext]
  createThread(readerThread, readerProc, ctx)
  createThread(processorThread, processorProc, ctx)
  createThread(writerThread, writerProc, ctx)
  joinThreads(readerThread, processorThread, writerThread)

  ctx.readChan.close()
  ctx.writeChan.close()

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
