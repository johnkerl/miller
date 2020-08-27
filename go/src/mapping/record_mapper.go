package mapping

import (
	"containers"
)

type RecordMapper interface {
	Map(inrec *containers.Lrec, outrecs chan<- *containers.Lrec)
}
