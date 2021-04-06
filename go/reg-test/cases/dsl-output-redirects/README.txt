Do either 'put -q' (no record-stream output) or use --opprint (record-stream
output is all at end of stream) since '> stdout' redirection decouples
record-stream output from print output, resulting in non-deterministic
output, which makes regtests fail.
