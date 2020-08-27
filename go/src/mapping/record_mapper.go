package mapping

import (
	"containers"
)

type RecordMapper interface {
	Name() string
	Map(inrec *containers.Lrec, outrecs chan<- *containers.Lrec)
}
