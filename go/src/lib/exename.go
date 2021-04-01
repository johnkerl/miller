package lib

func MlrExeName() string {
	// TODO:
	// This is ideal, so if someone has a 'mlr.debug' or somesuch, the messages will reflect that:

	// return path.Base(os.Args[0])

	// ... however it makes automated regression-testing hard, cross-platform. For example,
	// 'mlr' vs 'C:\something\something\mlr.exe'.
	return "mlr"
}
