package input

import (
	"containers"
	"strings"
)

// xxx to do: ifs and ips
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
		//lrec.Put(&key, &value)
		lrec.PutAtEnd(&key, &value)
	}
	return lrec
}
