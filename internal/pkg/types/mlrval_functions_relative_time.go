package types

import (
	"fmt"
	"math"
	"strings"
)

func BIF_dhms2sec(input1 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	var d, h, m, s int

	if strings.HasPrefix(input1.printrep, "-") {

		n, err := fmt.Sscanf(input1.printrep, "-%dd%dh%dm%ds", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return MlrvalFromInt(-(s + m*60 + h*60*60 + d*60*60*24))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%dh%dm%ds", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalFromInt(-(s + m*60 + h*60*60))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%dm%ds", &m, &s)
		if n == 2 && err == nil {
			return MlrvalFromInt(-(s + m*60))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%ds", &s)
		if n == 1 && err == nil {
			return MlrvalFromInt(-(s))
		}

	} else {

		n, err := fmt.Sscanf(input1.printrep, "%dd%dh%dm%ds", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return MlrvalFromInt(s + m*60 + h*60*60 + d*60*60*24)
		}
		n, err = fmt.Sscanf(input1.printrep, "%dh%dm%ds", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalFromInt(s + m*60 + h*60*60)
		}
		n, err = fmt.Sscanf(input1.printrep, "%dm%ds", &m, &s)
		if n == 2 && err == nil {
			return MlrvalFromInt(s + m*60)
		}
		n, err = fmt.Sscanf(input1.printrep, "%ds", &s)
		if n == 1 && err == nil {
			return MlrvalFromInt(s)
		}

	}
	return MLRVAL_ERROR
}

func BIF_dhms2fsec(input1 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}

	var d, h, m int
	var s float64

	if strings.HasPrefix(input1.printrep, "-") {

		n, err := fmt.Sscanf(input1.printrep, "-%dd%dh%dm%fs", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return MlrvalFromFloat64(-(s + float64(m*60+h*60*60+d*60*60*24)))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%dh%dm%fs", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalFromFloat64(-(s + float64(m*60+h*60*60)))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%dm%fs", &m, &s)
		if n == 2 && err == nil {
			return MlrvalFromFloat64(-(s + float64(m*60)))
		}
		n, err = fmt.Sscanf(input1.printrep, "-%fs", &s)
		if n == 1 && err == nil {
			return MlrvalFromFloat64(-(s))
		}

	} else {

		n, err := fmt.Sscanf(input1.printrep, "%dd%dh%dm%fs", &d, &h, &m, &s)
		if n == 4 && err == nil {
			return MlrvalFromFloat64(s + float64(m*60+h*60*60+d*60*60*24))
		}
		n, err = fmt.Sscanf(input1.printrep, "%dh%dm%fs", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalFromFloat64(s + float64(m*60+h*60*60))
		}
		n, err = fmt.Sscanf(input1.printrep, "%dm%fs", &m, &s)
		if n == 2 && err == nil {
			return MlrvalFromFloat64(s + float64(m*60))
		}
		n, err = fmt.Sscanf(input1.printrep, "%fs", &s)
		if n == 1 && err == nil {
			return MlrvalFromFloat64(s)
		}

	}
	return MLRVAL_ERROR
}

func BIF_hms2sec(input1 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}
	if input1.printrep == "" {
		return MLRVAL_ERROR
	}
	var h, m, s int

	if strings.HasPrefix(input1.printrep, "-") {
		n, err := fmt.Sscanf(input1.printrep, "-%d:%d:%d", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalFromInt(-(s + m*60 + h*60*60))
		}
	} else {
		n, err := fmt.Sscanf(input1.printrep, "%d:%d:%d", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalFromInt(s + m*60 + h*60*60)
		}
	}

	return MLRVAL_ERROR
}

func BIF_hms2fsec(input1 *Mlrval) *Mlrval {
	if input1.mvtype != MT_STRING {
		return MLRVAL_ERROR
	}

	var h, m int
	var s float64

	if strings.HasPrefix(input1.printrep, "-") {
		n, err := fmt.Sscanf(input1.printrep, "-%d:%d:%f", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalFromFloat64(-(s + float64(m*60+h*60*60)))
		}
	} else {
		n, err := fmt.Sscanf(input1.printrep, "%d:%d:%f", &h, &m, &s)
		if n == 3 && err == nil {
			return MlrvalFromFloat64(s + float64(m*60+h*60*60))
		}
	}

	return MLRVAL_ERROR
}

func BIF_sec2dhms(input1 *Mlrval) *Mlrval {
	isec, ok := input1.GetIntValue()
	if !ok {
		return MLRVAL_ERROR
	}

	var d, h, m, s int

	splitIntToDHMS(isec, &d, &h, &m, &s)
	if d != 0 {
		return MlrvalFromString(
			fmt.Sprintf("%dd%02dh%02dm%02ds", d, h, m, s),
		)
	} else if h != 0 {
		return MlrvalFromString(
			fmt.Sprintf("%dh%02dm%02ds", h, m, s),
		)
	} else if m != 0 {
		return MlrvalFromString(
			fmt.Sprintf("%dm%02ds", m, s),
		)
	} else {
		return MlrvalFromString(
			fmt.Sprintf("%ds", s),
		)
	}

	return MLRVAL_ERROR
}

func BIF_sec2hms(input1 *Mlrval) *Mlrval {
	isec, ok := input1.GetIntValue()
	if !ok {
		return MLRVAL_ERROR
	}
	sign := ""
	if isec < 0 {
		sign = "-"
		isec = -isec
	}

	var d, h, m, s int

	splitIntToDHMS(isec, &d, &h, &m, &s)
	h += d * 24

	return MlrvalFromString(
		fmt.Sprintf("%s%02d:%02d:%02d", sign, h, m, s),
	)

	return MLRVAL_ERROR
}

func BIF_fsec2dhms(input1 *Mlrval) *Mlrval {
	fsec, ok := input1.GetNumericToFloatValue()
	if !ok {
		return MLRVAL_ERROR
	}

	sign := 1
	if fsec < 0 {
		sign = -1
		fsec = -fsec
	}
	isec := int(math.Trunc(fsec))
	fractional := fsec - float64(isec)

	var d, h, m, s int

	splitIntToDHMS(isec, &d, &h, &m, &s)

	if d != 0 {
		d = sign * d
		return MlrvalFromString(
			fmt.Sprintf(
				"%dd%02dh%02dm%09.6fs",
				d, h, m, float64(s)+fractional),
		)
	} else if h != 0 {
		h = sign * h
		return MlrvalFromString(
			fmt.Sprintf(
				"%dh%02dm%09.6fs",
				h, m, float64(s)+fractional),
		)
	} else if m != 0 {
		m = sign * m
		return MlrvalFromString(
			fmt.Sprintf(
				"%dm%09.6fs",
				m, float64(s)+fractional),
		)
	} else {
		s = sign * s
		fractional = float64(sign) * fractional
		return MlrvalFromString(
			fmt.Sprintf(
				"%.6fs",
				float64(s)+fractional),
		)
	}
}

func BIF_fsec2hms(input1 *Mlrval) *Mlrval {
	fsec, ok := input1.GetNumericToFloatValue()
	if !ok {
		return MLRVAL_ERROR
	}

	sign := ""
	if fsec < 0 {
		sign = "-"
		fsec = -fsec
	}
	isec := int(math.Trunc(fsec))
	fractional := fsec - float64(isec)

	var d, h, m, s int

	splitIntToDHMS(isec, &d, &h, &m, &s)
	h += d * 24

	// "%02.6f" does not exist so we have to do our own zero-pad
	if s < 10 {
		return MlrvalFromString(
			fmt.Sprintf("%s%02d:%02d:0%.6f", sign, h, m, float64(s)+fractional),
		)
	} else {
		return MlrvalFromString(
			fmt.Sprintf("%s%02d:%02d:%.6f", sign, h, m, float64(s)+fractional),
		)
	}

	return MLRVAL_ERROR
}

// Helper function
func splitIntToDHMS(u int, pd, ph, pm, ps *int) {
	d := 0
	h := 0
	m := 0
	s := 0
	sign := 1
	if u < 0 {
		u = -u
		sign = -1
	}
	s = u % 60
	u = u / 60
	if u == 0 {
		s = s * sign
	} else {
		m = u % 60
		u = u / 60
		if u == 0 {
			m = m * sign
		} else {
			h = u % 24
			u = u / 24
			if u == 0 {
				h = h * sign
			} else {
				d = u * sign
			}
		}
	}
	*pd = d
	*ph = h
	*pm = m
	*ps = s
}
