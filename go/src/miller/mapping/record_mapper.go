package mapping

import (
	"miller/containers"
)

type IRecordMapper interface {
	Name() string
	Map(inrec *containers.Lrec, outrecs chan<- *containers.Lrec)
}
