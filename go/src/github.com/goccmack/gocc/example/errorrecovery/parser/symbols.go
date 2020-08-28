package parser

type symbol struct {
	Name       string
	IsTerminal bool
}

var symbols = [numSymbols]symbol{

	{
		Name:       "INVALID",
		IsTerminal: true,
	},

	{
		Name:       "$",
		IsTerminal: true,
	},

	{
		Name:       "S'",
		IsTerminal: false,
	},

	{
		Name:       "StmtList",
		IsTerminal: false,
	},

	{
		Name:       "Stmt",
		IsTerminal: false,
	},

	{
		Name:       "id",
		IsTerminal: true,
	},

	{
		Name:       "error",
		IsTerminal: true,
	},
}
