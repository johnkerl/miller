package lib

func Plural(n int) string {
	if n == 1 {
		return ""
	} else {
		return "s"
	}
}
