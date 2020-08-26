package input

import (
	"containers"
	"strings"
)

func LrecFromDKVPLine(
	line *string,
	ifs *string,
	ips *string,
) *containers.Lrec {
	lrec := containers.LrecAlloc()
	pairs := strings.Split(*line, *ifs)
	for _, pair := range(pairs) {
		kv := strings.SplitN(pair, *ips, 2)
		// xxx range-check
		key := kv[0]
		value := kv[1]
		// to do: avoid re-walk ...
		lrec.Put(&key, &value)
	}
	return lrec
}
