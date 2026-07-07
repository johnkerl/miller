package bifs

import (
	"github.com/johnkerl/miller/v6/pkg/mlrval"
)

// sparklineTicks are the eighth-height Unicode block characters used to
// render array/map values as a compact one-line bar chart.
var sparklineTicks = []rune("▁▂▃▄▅▆▇█")

func BIF_sparkline(collection *mlrval.Mlrval) *mlrval.Mlrval {
	ok, valueIfNot := check_collection(collection, "sparkline")
	if !ok {
		return valueIfNot
	}

	var floatValues []float64
	if collection.IsArray() {
		array := collection.AcquireArrayValue()
		floatValues = make([]float64, 0, len(array))
		for _, element := range array {
			floatValue, isFloat := element.GetNumericToFloatValue()
			if !isFloat {
				return mlrval.FromNotNumericError("sparkline", element)
			}
			floatValues = append(floatValues, floatValue)
		}
	} else {
		m := collection.AcquireMapValue()
		floatValues = make([]float64, 0, m.FieldCount)
		for pe := m.Head; pe != nil; pe = pe.Next {
			floatValue, isFloat := pe.Value.GetNumericToFloatValue()
			if !isFloat {
				return mlrval.FromNotNumericError("sparkline", pe.Value)
			}
			floatValues = append(floatValues, floatValue)
		}
	}

	if len(floatValues) == 0 {
		return mlrval.VOID
	}

	lo := floatValues[0]
	hi := floatValues[0]
	for _, floatValue := range floatValues[1:] {
		if floatValue < lo {
			lo = floatValue
		}
		if floatValue > hi {
			hi = floatValue
		}
	}

	numTicks := len(sparklineTicks)
	runes := make([]rune, len(floatValues))
	for i, floatValue := range floatValues {
		if hi == lo {
			runes[i] = sparklineTicks[0]
			continue
		}
		tickIndex := int(float64(numTicks-1)*(floatValue-lo)/(hi-lo) + 0.5)
		if tickIndex < 0 {
			tickIndex = 0
		} else if tickIndex >= numTicks {
			tickIndex = numTicks - 1
		}
		runes[i] = sparklineTicks[tickIndex]
	}

	return mlrval.FromString(string(runes))
}
