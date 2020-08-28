package parser

import "github.com/goccmack/gocc/internal/ast"

var ProductionsTable = ProdTab{
	// [0]
	ProdTabEntry{
		"S! : Grammar ;",
		"S!",
		1,
		func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	},
	// [1]
	ProdTabEntry{
		"Grammar : LexicalPart SyntaxPart << ast.NewGrammar(X[0], X[1]) >> ;",
		"Grammar",
		2,
		func(X []Attrib) (Attrib, error) {
			return ast.NewGrammar(X[0], X[1])
		},
	},
	// [2]
	ProdTabEntry{
		"Grammar : LexicalPart << ast.NewGrammar(X[0], nil) >> ;",
		"Grammar",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewGrammar(X[0], nil)
		},
	},
	// [3]
	ProdTabEntry{
		"Grammar : SyntaxPart << ast.NewGrammar(nil, X[0]) >> ;",
		"Grammar",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewGrammar(nil, X[0])
		},
	},
	// [4]
	ProdTabEntry{
		"LexicalPart : LexProductions << ast.NewLexPart(nil, nil, X[0]) >> ;",
		"LexicalPart",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexPart(nil, nil, X[0])
		},
	},
	// [5]
	ProdTabEntry{
		"LexProductions : LexProduction << ast.NewLexProductions(X[0]) >> ;",
		"LexProductions",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexProductions(X[0])
		},
	},
	// [6]
	ProdTabEntry{
		"LexProductions : LexProductions LexProduction << ast.AppendLexProduction(X[0], X[1]) >> ;",
		"LexProductions",
		2,
		func(X []Attrib) (Attrib, error) {
			return ast.AppendLexProduction(X[0], X[1])
		},
	},
	// [7]
	ProdTabEntry{
		"LexProduction : tokId : LexPattern ; << ast.NewLexTokDef(X[0], X[2]) >> ;",
		"LexProduction",
		4,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexTokDef(X[0], X[2])
		},
	},
	// [8]
	ProdTabEntry{
		"LexProduction : regDefId : LexPattern ; << ast.NewLexRegDef(X[0], X[2]) >> ;",
		"LexProduction",
		4,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexRegDef(X[0], X[2])
		},
	},
	// [9]
	ProdTabEntry{
		"LexProduction : ignoredTokId : LexPattern ; << ast.NewLexIgnoredTokDef(X[0], X[2]) >> ;",
		"LexProduction",
		4,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexIgnoredTokDef(X[0], X[2])
		},
	},
	// [10]
	ProdTabEntry{
		"LexPattern : LexAlt << ast.NewLexPattern(X[0]) >> ;",
		"LexPattern",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexPattern(X[0])
		},
	},
	// [11]
	ProdTabEntry{
		"LexPattern : LexPattern | LexAlt << ast.AppendLexAlt(X[0], X[2]) >> ;",
		"LexPattern",
		3,
		func(X []Attrib) (Attrib, error) {
			return ast.AppendLexAlt(X[0], X[2])
		},
	},
	// [12]
	ProdTabEntry{
		"LexAlt : LexTerm << ast.NewLexAlt(X[0]) >> ;",
		"LexAlt",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexAlt(X[0])
		},
	},
	// [13]
	ProdTabEntry{
		"LexAlt : LexAlt LexTerm << ast.AppendLexTerm(X[0], X[1]) >> ;",
		"LexAlt",
		2,
		func(X []Attrib) (Attrib, error) {
			return ast.AppendLexTerm(X[0], X[1])
		},
	},
	// [14]
	ProdTabEntry{
		"LexTerm : . << ast.LexDOT, nil >> ;",
		"LexTerm",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.LexDOT, nil
		},
	},
	// [15]
	ProdTabEntry{
		"LexTerm : char_lit << ast.NewLexCharLit(X[0]) >> ;",
		"LexTerm",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexCharLit(X[0])
		},
	},
	// [16]
	ProdTabEntry{
		"LexTerm : char_lit - char_lit << ast.NewLexCharRange(X[0], X[2]) >> ;",
		"LexTerm",
		3,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexCharRange(X[0], X[2])
		},
	},
	// [17]
	ProdTabEntry{
		"LexTerm : regDefId << ast.NewLexRegDefId(X[0]) >> ;",
		"LexTerm",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexRegDefId(X[0])
		},
	},
	// [18]
	ProdTabEntry{
		"LexTerm : [ LexPattern ] << ast.NewLexOptPattern(X[1]) >> ;",
		"LexTerm",
		3,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexOptPattern(X[1])
		},
	},
	// [19]
	ProdTabEntry{
		"LexTerm : { LexPattern } << ast.NewLexRepPattern(X[1]) >> ;",
		"LexTerm",
		3,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexRepPattern(X[1])
		},
	},
	// [20]
	ProdTabEntry{
		"LexTerm : ( LexPattern ) << ast.NewLexGroupPattern(X[1]) >> ;",
		"LexTerm",
		3,
		func(X []Attrib) (Attrib, error) {
			return ast.NewLexGroupPattern(X[1])
		},
	},
	// [21]
	ProdTabEntry{
		"SyntaxPart : FileHeader SyntaxProdList << ast.NewSyntaxPart(X[0], X[1]) >> ;",
		"SyntaxPart",
		2,
		func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxPart(X[0], X[1])
		},
	},
	// [22]
	ProdTabEntry{
		"SyntaxPart : SyntaxProdList << ast.NewSyntaxPart(nil, X[0]) >> ;",
		"SyntaxPart",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxPart(nil, X[0])
		},
	},
	// [23]
	ProdTabEntry{
		"SyntaxProdList : SyntaxProduction << ast.NewSyntaxProdList(X[0]) >> ;",
		"SyntaxProdList",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxProdList(X[0])
		},
	},
	// [24]
	ProdTabEntry{
		"SyntaxProdList : SyntaxProdList SyntaxProduction << ast.AddSyntaxProds(X[0], X[1]) >> ;",
		"SyntaxProdList",
		2,
		func(X []Attrib) (Attrib, error) {
			return ast.AddSyntaxProds(X[0], X[1])
		},
	},
	// [25]
	ProdTabEntry{
		"SyntaxProduction : prodId : Alternatives ; << ast.NewSyntaxProd(X[0], X[2]) >> ;",
		"SyntaxProduction",
		4,
		func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxProd(X[0], X[2])
		},
	},
	// [26]
	ProdTabEntry{
		"Alternatives : SyntaxBody << ast.NewSyntaxAlts(X[0]) >> ;",
		"Alternatives",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxAlts(X[0])
		},
	},
	// [27]
	ProdTabEntry{
		"Alternatives : Alternatives | SyntaxBody << ast.AddSyntaxAlt(X[0], X[2]) >> ;",
		"Alternatives",
		3,
		func(X []Attrib) (Attrib, error) {
			return ast.AddSyntaxAlt(X[0], X[2])
		},
	},
	// [28]
	ProdTabEntry{
		"SyntaxBody : Symbols << ast.NewSyntaxBody(X[0], nil) >> ;",
		"SyntaxBody",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxBody(X[0], nil)
		},
	},
	// [29]
	ProdTabEntry{
		"SyntaxBody : Symbols g_sdt_lit << ast.NewSyntaxBody(X[0], X[1]) >> ;",
		"SyntaxBody",
		2,
		func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxBody(X[0], X[1])
		},
	},
	// [30]
	ProdTabEntry{
		"SyntaxBody : error << ast.NewErrorBody(nil, nil) >> ;",
		"SyntaxBody",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewErrorBody(nil, nil)
		},
	},
	// [31]
	ProdTabEntry{
		"SyntaxBody : error Symbols << ast.NewErrorBody(X[1], nil) >> ;",
		"SyntaxBody",
		2,
		func(X []Attrib) (Attrib, error) {
			return ast.NewErrorBody(X[1], nil)
		},
	},
	// [32]
	ProdTabEntry{
		"SyntaxBody : error Symbols g_sdt_lit << ast.NewErrorBody(X[1], X[2]) >> ;",
		"SyntaxBody",
		3,
		func(X []Attrib) (Attrib, error) {
			return ast.NewErrorBody(X[1], X[2])
		},
	},
	// [33]
	ProdTabEntry{
		"SyntaxBody : empty << ast.NewEmptyBody() >> ;",
		"SyntaxBody",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewEmptyBody()
		},
	},
	// [34]
	ProdTabEntry{
		"Symbols : Symbol << ast.NewSyntaxSymbols(X[0]) >> ;",
		"Symbols",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxSymbols(X[0])
		},
	},
	// [35]
	ProdTabEntry{
		"Symbols : Symbols Symbol << ast.AddSyntaxSymbol(X[0], X[1]) >> ;",
		"Symbols",
		2,
		func(X []Attrib) (Attrib, error) {
			return ast.AddSyntaxSymbol(X[0], X[1])
		},
	},
	// [36]
	ProdTabEntry{
		"Symbol : prodId << ast.NewSyntaxProdId(X[0]) >> ;",
		"Symbol",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxProdId(X[0])
		},
	},
	// [37]
	ProdTabEntry{
		"Symbol : tokId << ast.NewTokId(X[0]) >> ;",
		"Symbol",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewTokId(X[0])
		},
	},
	// [38]
	ProdTabEntry{
		"Symbol : string_lit << ast.NewStringLit(X[0]) >> ;",
		"Symbol",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewStringLit(X[0])
		},
	},
	// [39]
	ProdTabEntry{
		"FileHeader : g_sdt_lit << ast.NewFileHeader(X[0]) >> ;",
		"FileHeader",
		1,
		func(X []Attrib) (Attrib, error) {
			return ast.NewFileHeader(X[0])
		},
	},
}

var ActionTable ActionTab = ActionTab{
	// state 0
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			18: Shift(13), // g_sdt_lit
			2:  Shift(6),  // tokId
			5:  Shift(7),  // regDefId
			6:  Shift(8),  // ignoredTokId
			17: Shift(12), // prodId
		},
	},

	// state 1
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			0: Accept(0), // $
		},
	},

	// state 2
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			0:  Reduce(2), // $
			18: Shift(13), // g_sdt_lit
			17: Shift(12), // prodId
		},
	},

	// state 3
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			0: Reduce(3), // $
		},
	},

	// state 4
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			18: Reduce(4), // g_sdt_lit
			17: Reduce(4), // prodId
			0:  Reduce(4), // $
			2:  Shift(6),  // tokId
			5:  Shift(7),  // regDefId
			6:  Shift(8),  // ignoredTokId
		},
	},

	// state 5
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			18: Reduce(5), // g_sdt_lit
			17: Reduce(5), // prodId
			0:  Reduce(5), // $
			2:  Reduce(5), // tokId
			5:  Reduce(5), // regDefId
			6:  Reduce(5), // ignoredTokId
		},
	},

	// state 6
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			3: Shift(16), // :
		},
	},

	// state 7
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			3: Shift(17), // :
		},
	},

	// state 8
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			3: Shift(18), // :
		},
	},

	// state 9
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			17: Shift(12), // prodId
		},
	},

	// state 10
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			0:  Reduce(22), // $
			17: Shift(12),  // prodId
		},
	},

	// state 11
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			0:  Reduce(23), // $
			17: Reduce(23), // prodId
		},
	},

	// state 12
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			3: Shift(21), // :
		},
	},

	// state 13
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			17: Reduce(39), // prodId
		},
	},

	// state 14
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			0: Reduce(1), // $
		},
	},

	// state 15
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			18: Reduce(6), // g_sdt_lit
			17: Reduce(6), // prodId
			0:  Reduce(6), // $
			2:  Reduce(6), // tokId
			5:  Reduce(6), // regDefId
			6:  Reduce(6), // ignoredTokId
		},
	},

	// state 16
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(26), // .
			9:  Shift(27), // char_lit
			5:  Shift(23), // regDefId
			11: Shift(28), // [
			13: Shift(29), // {
			15: Shift(30), // (
		},
	},

	// state 17
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(26), // .
			9:  Shift(27), // char_lit
			5:  Shift(23), // regDefId
			11: Shift(28), // [
			13: Shift(29), // {
			15: Shift(30), // (
		},
	},

	// state 18
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(26), // .
			9:  Shift(27), // char_lit
			5:  Shift(23), // regDefId
			11: Shift(28), // [
			13: Shift(29), // {
			15: Shift(30), // (
		},
	},

	// state 19
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			0:  Reduce(21), // $
			17: Shift(12),  // prodId
		},
	},

	// state 20
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			0:  Reduce(24), // $
			17: Reduce(24), // prodId
		},
	},

	// state 21
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			19: Shift(38), // error
			20: Shift(39), // empty
			17: Shift(34), // prodId
			2:  Shift(33), // tokId
			21: Shift(41), // string_lit
		},
	},

	// state 22
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4: Shift(42), // ;
			7: Shift(43), // |
		},
	},

	// state 23
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(17), // ;
			8:  Reduce(17), // .
			9:  Reduce(17), // char_lit
			5:  Reduce(17), // regDefId
			11: Reduce(17), // [
			13: Reduce(17), // {
			15: Reduce(17), // (
			7:  Reduce(17), // |
		},
	},

	// state 24
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(10), // ;
			7:  Reduce(10), // |
			8:  Shift(26),  // .
			9:  Shift(27),  // char_lit
			5:  Shift(23),  // regDefId
			11: Shift(28),  // [
			13: Shift(29),  // {
			15: Shift(30),  // (
		},
	},

	// state 25
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(12), // ;
			8:  Reduce(12), // .
			9:  Reduce(12), // char_lit
			5:  Reduce(12), // regDefId
			11: Reduce(12), // [
			13: Reduce(12), // {
			15: Reduce(12), // (
			7:  Reduce(12), // |
		},
	},

	// state 26
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(14), // ;
			8:  Reduce(14), // .
			9:  Reduce(14), // char_lit
			5:  Reduce(14), // regDefId
			11: Reduce(14), // [
			13: Reduce(14), // {
			15: Reduce(14), // (
			7:  Reduce(14), // |
		},
	},

	// state 27
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			10: Shift(45),  // -
			9:  Reduce(15), // char_lit
			5:  Reduce(15), // regDefId
			13: Reduce(15), // {
			15: Reduce(15), // (
			4:  Reduce(15), // ;
			8:  Reduce(15), // .
			11: Reduce(15), // [
			7:  Reduce(15), // |
		},
	},

	// state 28
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(50), // .
			9:  Shift(51), // char_lit
			5:  Shift(47), // regDefId
			11: Shift(52), // [
			13: Shift(53), // {
			15: Shift(54), // (
		},
	},

	// state 29
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(59), // .
			9:  Shift(60), // char_lit
			5:  Shift(56), // regDefId
			11: Shift(61), // [
			13: Shift(62), // {
			15: Shift(63), // (
		},
	},

	// state 30
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(68), // .
			9:  Shift(69), // char_lit
			5:  Shift(65), // regDefId
			11: Shift(70), // [
			13: Shift(71), // {
			15: Shift(72), // (
		},
	},

	// state 31
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4: Shift(73), // ;
			7: Shift(43), // |
		},
	},

	// state 32
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4: Shift(74), // ;
			7: Shift(43), // |
		},
	},

	// state 33
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(37), // ;
			18: Reduce(37), // g_sdt_lit
			17: Reduce(37), // prodId
			2:  Reduce(37), // tokId
			21: Reduce(37), // string_lit
			7:  Reduce(37), // |
		},
	},

	// state 34
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(36), // ;
			18: Reduce(36), // g_sdt_lit
			17: Reduce(36), // prodId
			2:  Reduce(36), // tokId
			21: Reduce(36), // string_lit
			7:  Reduce(36), // |
		},
	},

	// state 35
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4: Shift(75), // ;
			7: Shift(76), // |
		},
	},

	// state 36
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4: Reduce(26), // ;
			7: Reduce(26), // |
		},
	},

	// state 37
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(28), // ;
			18: Shift(77),  // g_sdt_lit
			7:  Reduce(28), // |
			17: Shift(34),  // prodId
			2:  Shift(33),  // tokId
			21: Shift(41),  // string_lit
		},
	},

	// state 38
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(30), // ;
			7:  Reduce(30), // |
			17: Shift(34),  // prodId
			2:  Shift(33),  // tokId
			21: Shift(41),  // string_lit
		},
	},

	// state 39
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4: Reduce(33), // ;
			7: Reduce(33), // |
		},
	},

	// state 40
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(34), // ;
			18: Reduce(34), // g_sdt_lit
			17: Reduce(34), // prodId
			2:  Reduce(34), // tokId
			21: Reduce(34), // string_lit
			7:  Reduce(34), // |
		},
	},

	// state 41
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(38), // ;
			18: Reduce(38), // g_sdt_lit
			17: Reduce(38), // prodId
			2:  Reduce(38), // tokId
			21: Reduce(38), // string_lit
			7:  Reduce(38), // |
		},
	},

	// state 42
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			18: Reduce(7), // g_sdt_lit
			17: Reduce(7), // prodId
			0:  Reduce(7), // $
			2:  Reduce(7), // tokId
			5:  Reduce(7), // regDefId
			6:  Reduce(7), // ignoredTokId
		},
	},

	// state 43
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(26), // .
			9:  Shift(27), // char_lit
			5:  Shift(23), // regDefId
			11: Shift(28), // [
			13: Shift(29), // {
			15: Shift(30), // (
		},
	},

	// state 44
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(13), // ;
			8:  Reduce(13), // .
			9:  Reduce(13), // char_lit
			5:  Reduce(13), // regDefId
			11: Reduce(13), // [
			13: Reduce(13), // {
			15: Reduce(13), // (
			7:  Reduce(13), // |
		},
	},

	// state 45
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			9: Shift(81), // char_lit
		},
	},

	// state 46
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Shift(83), // ]
			7:  Shift(82), // |
		},
	},

	// state 47
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(17), // ]
			8:  Reduce(17), // .
			9:  Reduce(17), // char_lit
			5:  Reduce(17), // regDefId
			11: Reduce(17), // [
			13: Reduce(17), // {
			15: Reduce(17), // (
			7:  Reduce(17), // |
		},
	},

	// state 48
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(10), // ]
			7:  Reduce(10), // |
			8:  Shift(50),  // .
			9:  Shift(51),  // char_lit
			5:  Shift(47),  // regDefId
			11: Shift(52),  // [
			13: Shift(53),  // {
			15: Shift(54),  // (
		},
	},

	// state 49
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(12), // ]
			8:  Reduce(12), // .
			9:  Reduce(12), // char_lit
			5:  Reduce(12), // regDefId
			11: Reduce(12), // [
			13: Reduce(12), // {
			15: Reduce(12), // (
			7:  Reduce(12), // |
		},
	},

	// state 50
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(14), // ]
			8:  Reduce(14), // .
			9:  Reduce(14), // char_lit
			5:  Reduce(14), // regDefId
			11: Reduce(14), // [
			13: Reduce(14), // {
			15: Reduce(14), // (
			7:  Reduce(14), // |
		},
	},

	// state 51
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(15), // ]
			10: Shift(85),  // -
			15: Reduce(15), // (
			8:  Reduce(15), // .
			9:  Reduce(15), // char_lit
			5:  Reduce(15), // regDefId
			11: Reduce(15), // [
			13: Reduce(15), // {
			7:  Reduce(15), // |
		},
	},

	// state 52
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(50), // .
			9:  Shift(51), // char_lit
			5:  Shift(47), // regDefId
			11: Shift(52), // [
			13: Shift(53), // {
			15: Shift(54), // (
		},
	},

	// state 53
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(59), // .
			9:  Shift(60), // char_lit
			5:  Shift(56), // regDefId
			11: Shift(61), // [
			13: Shift(62), // {
			15: Shift(63), // (
		},
	},

	// state 54
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(68), // .
			9:  Shift(69), // char_lit
			5:  Shift(65), // regDefId
			11: Shift(70), // [
			13: Shift(71), // {
			15: Shift(72), // (
		},
	},

	// state 55
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Shift(90), // }
			7:  Shift(89), // |
		},
	},

	// state 56
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(17), // }
			8:  Reduce(17), // .
			9:  Reduce(17), // char_lit
			5:  Reduce(17), // regDefId
			11: Reduce(17), // [
			13: Reduce(17), // {
			15: Reduce(17), // (
			7:  Reduce(17), // |
		},
	},

	// state 57
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(10), // }
			7:  Reduce(10), // |
			8:  Shift(59),  // .
			9:  Shift(60),  // char_lit
			5:  Shift(56),  // regDefId
			11: Shift(61),  // [
			13: Shift(62),  // {
			15: Shift(63),  // (
		},
	},

	// state 58
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(12), // }
			8:  Reduce(12), // .
			9:  Reduce(12), // char_lit
			5:  Reduce(12), // regDefId
			11: Reduce(12), // [
			13: Reduce(12), // {
			15: Reduce(12), // (
			7:  Reduce(12), // |
		},
	},

	// state 59
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(14), // }
			8:  Reduce(14), // .
			9:  Reduce(14), // char_lit
			5:  Reduce(14), // regDefId
			11: Reduce(14), // [
			13: Reduce(14), // {
			15: Reduce(14), // (
			7:  Reduce(14), // |
		},
	},

	// state 60
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(15), // }
			9:  Reduce(15), // char_lit
			11: Reduce(15), // [
			15: Reduce(15), // (
			10: Shift(92),  // -
			8:  Reduce(15), // .
			5:  Reduce(15), // regDefId
			13: Reduce(15), // {
			7:  Reduce(15), // |
		},
	},

	// state 61
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(50), // .
			9:  Shift(51), // char_lit
			5:  Shift(47), // regDefId
			11: Shift(52), // [
			13: Shift(53), // {
			15: Shift(54), // (
		},
	},

	// state 62
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(59), // .
			9:  Shift(60), // char_lit
			5:  Shift(56), // regDefId
			11: Shift(61), // [
			13: Shift(62), // {
			15: Shift(63), // (
		},
	},

	// state 63
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(68), // .
			9:  Shift(69), // char_lit
			5:  Shift(65), // regDefId
			11: Shift(70), // [
			13: Shift(71), // {
			15: Shift(72), // (
		},
	},

	// state 64
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Shift(97), // )
			7:  Shift(96), // |
		},
	},

	// state 65
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Reduce(17), // )
			8:  Reduce(17), // .
			9:  Reduce(17), // char_lit
			5:  Reduce(17), // regDefId
			11: Reduce(17), // [
			13: Reduce(17), // {
			15: Reduce(17), // (
			7:  Reduce(17), // |
		},
	},

	// state 66
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Reduce(10), // )
			7:  Reduce(10), // |
			8:  Shift(68),  // .
			9:  Shift(69),  // char_lit
			5:  Shift(65),  // regDefId
			11: Shift(70),  // [
			13: Shift(71),  // {
			15: Shift(72),  // (
		},
	},

	// state 67
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Reduce(12), // )
			8:  Reduce(12), // .
			9:  Reduce(12), // char_lit
			5:  Reduce(12), // regDefId
			11: Reduce(12), // [
			13: Reduce(12), // {
			15: Reduce(12), // (
			7:  Reduce(12), // |
		},
	},

	// state 68
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Reduce(14), // )
			8:  Reduce(14), // .
			9:  Reduce(14), // char_lit
			5:  Reduce(14), // regDefId
			11: Reduce(14), // [
			13: Reduce(14), // {
			15: Reduce(14), // (
			7:  Reduce(14), // |
		},
	},

	// state 69
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Reduce(15), // .
			9:  Reduce(15), // char_lit
			5:  Reduce(15), // regDefId
			11: Reduce(15), // [
			13: Reduce(15), // {
			7:  Reduce(15), // |
			16: Reduce(15), // )
			10: Shift(99),  // -
			15: Reduce(15), // (
		},
	},

	// state 70
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(50), // .
			9:  Shift(51), // char_lit
			5:  Shift(47), // regDefId
			11: Shift(52), // [
			13: Shift(53), // {
			15: Shift(54), // (
		},
	},

	// state 71
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(59), // .
			9:  Shift(60), // char_lit
			5:  Shift(56), // regDefId
			11: Shift(61), // [
			13: Shift(62), // {
			15: Shift(63), // (
		},
	},

	// state 72
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(68), // .
			9:  Shift(69), // char_lit
			5:  Shift(65), // regDefId
			11: Shift(70), // [
			13: Shift(71), // {
			15: Shift(72), // (
		},
	},

	// state 73
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			18: Reduce(8), // g_sdt_lit
			17: Reduce(8), // prodId
			0:  Reduce(8), // $
			2:  Reduce(8), // tokId
			5:  Reduce(8), // regDefId
			6:  Reduce(8), // ignoredTokId
		},
	},

	// state 74
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			18: Reduce(9), // g_sdt_lit
			17: Reduce(9), // prodId
			0:  Reduce(9), // $
			2:  Reduce(9), // tokId
			5:  Reduce(9), // regDefId
			6:  Reduce(9), // ignoredTokId
		},
	},

	// state 75
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			0:  Reduce(25), // $
			17: Reduce(25), // prodId
		},
	},

	// state 76
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			19: Shift(38), // error
			20: Shift(39), // empty
			17: Shift(34), // prodId
			2:  Shift(33), // tokId
			21: Shift(41), // string_lit
		},
	},

	// state 77
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4: Reduce(29), // ;
			7: Reduce(29), // |
		},
	},

	// state 78
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(35), // ;
			18: Reduce(35), // g_sdt_lit
			17: Reduce(35), // prodId
			2:  Reduce(35), // tokId
			21: Reduce(35), // string_lit
			7:  Reduce(35), // |
		},
	},

	// state 79
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(31), // ;
			18: Shift(104), // g_sdt_lit
			7:  Reduce(31), // |
			17: Shift(34),  // prodId
			2:  Shift(33),  // tokId
			21: Shift(41),  // string_lit
		},
	},

	// state 80
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(11), // ;
			7:  Reduce(11), // |
			8:  Shift(26),  // .
			9:  Shift(27),  // char_lit
			5:  Shift(23),  // regDefId
			11: Shift(28),  // [
			13: Shift(29),  // {
			15: Shift(30),  // (
		},
	},

	// state 81
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(16), // ;
			8:  Reduce(16), // .
			9:  Reduce(16), // char_lit
			5:  Reduce(16), // regDefId
			11: Reduce(16), // [
			13: Reduce(16), // {
			15: Reduce(16), // (
			7:  Reduce(16), // |
		},
	},

	// state 82
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(50), // .
			9:  Shift(51), // char_lit
			5:  Shift(47), // regDefId
			11: Shift(52), // [
			13: Shift(53), // {
			15: Shift(54), // (
		},
	},

	// state 83
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(18), // ;
			8:  Reduce(18), // .
			9:  Reduce(18), // char_lit
			5:  Reduce(18), // regDefId
			11: Reduce(18), // [
			13: Reduce(18), // {
			15: Reduce(18), // (
			7:  Reduce(18), // |
		},
	},

	// state 84
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(13), // ]
			8:  Reduce(13), // .
			9:  Reduce(13), // char_lit
			5:  Reduce(13), // regDefId
			11: Reduce(13), // [
			13: Reduce(13), // {
			15: Reduce(13), // (
			7:  Reduce(13), // |
		},
	},

	// state 85
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			9: Shift(106), // char_lit
		},
	},

	// state 86
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Shift(107), // ]
			7:  Shift(82),  // |
		},
	},

	// state 87
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Shift(108), // }
			7:  Shift(89),  // |
		},
	},

	// state 88
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Shift(109), // )
			7:  Shift(96),  // |
		},
	},

	// state 89
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(59), // .
			9:  Shift(60), // char_lit
			5:  Shift(56), // regDefId
			11: Shift(61), // [
			13: Shift(62), // {
			15: Shift(63), // (
		},
	},

	// state 90
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(19), // ;
			8:  Reduce(19), // .
			9:  Reduce(19), // char_lit
			5:  Reduce(19), // regDefId
			11: Reduce(19), // [
			13: Reduce(19), // {
			15: Reduce(19), // (
			7:  Reduce(19), // |
		},
	},

	// state 91
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(13), // }
			8:  Reduce(13), // .
			9:  Reduce(13), // char_lit
			5:  Reduce(13), // regDefId
			11: Reduce(13), // [
			13: Reduce(13), // {
			15: Reduce(13), // (
			7:  Reduce(13), // |
		},
	},

	// state 92
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			9: Shift(111), // char_lit
		},
	},

	// state 93
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Shift(112), // ]
			7:  Shift(82),  // |
		},
	},

	// state 94
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Shift(113), // }
			7:  Shift(89),  // |
		},
	},

	// state 95
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Shift(114), // )
			7:  Shift(96),  // |
		},
	},

	// state 96
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			8:  Shift(68), // .
			9:  Shift(69), // char_lit
			5:  Shift(65), // regDefId
			11: Shift(70), // [
			13: Shift(71), // {
			15: Shift(72), // (
		},
	},

	// state 97
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4:  Reduce(20), // ;
			8:  Reduce(20), // .
			9:  Reduce(20), // char_lit
			5:  Reduce(20), // regDefId
			11: Reduce(20), // [
			13: Reduce(20), // {
			15: Reduce(20), // (
			7:  Reduce(20), // |
		},
	},

	// state 98
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Reduce(13), // )
			8:  Reduce(13), // .
			9:  Reduce(13), // char_lit
			5:  Reduce(13), // regDefId
			11: Reduce(13), // [
			13: Reduce(13), // {
			15: Reduce(13), // (
			7:  Reduce(13), // |
		},
	},

	// state 99
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			9: Shift(116), // char_lit
		},
	},

	// state 100
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Shift(117), // ]
			7:  Shift(82),  // |
		},
	},

	// state 101
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Shift(118), // }
			7:  Shift(89),  // |
		},
	},

	// state 102
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Shift(119), // )
			7:  Shift(96),  // |
		},
	},

	// state 103
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4: Reduce(27), // ;
			7: Reduce(27), // |
		},
	},

	// state 104
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			4: Reduce(32), // ;
			7: Reduce(32), // |
		},
	},

	// state 105
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(11), // ]
			7:  Reduce(11), // |
			8:  Shift(50),  // .
			9:  Shift(51),  // char_lit
			5:  Shift(47),  // regDefId
			11: Shift(52),  // [
			13: Shift(53),  // {
			15: Shift(54),  // (
		},
	},

	// state 106
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(16), // ]
			8:  Reduce(16), // .
			9:  Reduce(16), // char_lit
			5:  Reduce(16), // regDefId
			11: Reduce(16), // [
			13: Reduce(16), // {
			15: Reduce(16), // (
			7:  Reduce(16), // |
		},
	},

	// state 107
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(18), // ]
			8:  Reduce(18), // .
			9:  Reduce(18), // char_lit
			5:  Reduce(18), // regDefId
			11: Reduce(18), // [
			13: Reduce(18), // {
			15: Reduce(18), // (
			7:  Reduce(18), // |
		},
	},

	// state 108
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(19), // ]
			8:  Reduce(19), // .
			9:  Reduce(19), // char_lit
			5:  Reduce(19), // regDefId
			11: Reduce(19), // [
			13: Reduce(19), // {
			15: Reduce(19), // (
			7:  Reduce(19), // |
		},
	},

	// state 109
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			12: Reduce(20), // ]
			8:  Reduce(20), // .
			9:  Reduce(20), // char_lit
			5:  Reduce(20), // regDefId
			11: Reduce(20), // [
			13: Reduce(20), // {
			15: Reduce(20), // (
			7:  Reduce(20), // |
		},
	},

	// state 110
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(11), // }
			7:  Reduce(11), // |
			8:  Shift(59),  // .
			9:  Shift(60),  // char_lit
			5:  Shift(56),  // regDefId
			11: Shift(61),  // [
			13: Shift(62),  // {
			15: Shift(63),  // (
		},
	},

	// state 111
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(16), // }
			8:  Reduce(16), // .
			9:  Reduce(16), // char_lit
			5:  Reduce(16), // regDefId
			11: Reduce(16), // [
			13: Reduce(16), // {
			15: Reduce(16), // (
			7:  Reduce(16), // |
		},
	},

	// state 112
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(18), // }
			8:  Reduce(18), // .
			9:  Reduce(18), // char_lit
			5:  Reduce(18), // regDefId
			11: Reduce(18), // [
			13: Reduce(18), // {
			15: Reduce(18), // (
			7:  Reduce(18), // |
		},
	},

	// state 113
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(19), // }
			8:  Reduce(19), // .
			9:  Reduce(19), // char_lit
			5:  Reduce(19), // regDefId
			11: Reduce(19), // [
			13: Reduce(19), // {
			15: Reduce(19), // (
			7:  Reduce(19), // |
		},
	},

	// state 114
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			14: Reduce(20), // }
			8:  Reduce(20), // .
			9:  Reduce(20), // char_lit
			5:  Reduce(20), // regDefId
			11: Reduce(20), // [
			13: Reduce(20), // {
			15: Reduce(20), // (
			7:  Reduce(20), // |
		},
	},

	// state 115
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Reduce(11), // )
			7:  Reduce(11), // |
			8:  Shift(68),  // .
			9:  Shift(69),  // char_lit
			5:  Shift(65),  // regDefId
			11: Shift(70),  // [
			13: Shift(71),  // {
			15: Shift(72),  // (
		},
	},

	// state 116
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Reduce(16), // )
			8:  Reduce(16), // .
			9:  Reduce(16), // char_lit
			5:  Reduce(16), // regDefId
			11: Reduce(16), // [
			13: Reduce(16), // {
			15: Reduce(16), // (
			7:  Reduce(16), // |
		},
	},

	// state 117
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Reduce(18), // )
			8:  Reduce(18), // .
			9:  Reduce(18), // char_lit
			5:  Reduce(18), // regDefId
			11: Reduce(18), // [
			13: Reduce(18), // {
			15: Reduce(18), // (
			7:  Reduce(18), // |
		},
	},

	// state 118
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Reduce(19), // )
			8:  Reduce(19), // .
			9:  Reduce(19), // char_lit
			5:  Reduce(19), // regDefId
			11: Reduce(19), // [
			13: Reduce(19), // {
			15: Reduce(19), // (
			7:  Reduce(19), // |
		},
	},

	// state 119
	&ActionRow{
		canRecover: false,
		Actions: Actions{
			16: Reduce(20), // )
			8:  Reduce(20), // .
			9:  Reduce(20), // char_lit
			5:  Reduce(20), // regDefId
			11: Reduce(20), // [
			13: Reduce(20), // {
			15: Reduce(20), // (
			7:  Reduce(20), // |
		},
	},
}

var GotoTable GotoTab = GotoTab{
	// state 0
	GotoRow{
		"Grammar":          State(1),
		"LexicalPart":      State(2),
		"SyntaxPart":       State(3),
		"LexProductions":   State(4),
		"LexProduction":    State(5),
		"FileHeader":       State(9),
		"SyntaxProdList":   State(10),
		"SyntaxProduction": State(11),
	},
	// state 1
	GotoRow{},
	// state 2
	GotoRow{
		"SyntaxPart":       State(14),
		"FileHeader":       State(9),
		"SyntaxProdList":   State(10),
		"SyntaxProduction": State(11),
	},
	// state 3
	GotoRow{},
	// state 4
	GotoRow{
		"LexProduction": State(15),
	},
	// state 5
	GotoRow{},
	// state 6
	GotoRow{},
	// state 7
	GotoRow{},
	// state 8
	GotoRow{},
	// state 9
	GotoRow{
		"SyntaxProdList":   State(19),
		"SyntaxProduction": State(11),
	},
	// state 10
	GotoRow{
		"SyntaxProduction": State(20),
	},
	// state 11
	GotoRow{},
	// state 12
	GotoRow{},
	// state 13
	GotoRow{},
	// state 14
	GotoRow{},
	// state 15
	GotoRow{},
	// state 16
	GotoRow{
		"LexPattern": State(22),
		"LexAlt":     State(24),
		"LexTerm":    State(25),
	},
	// state 17
	GotoRow{
		"LexPattern": State(31),
		"LexAlt":     State(24),
		"LexTerm":    State(25),
	},
	// state 18
	GotoRow{
		"LexPattern": State(32),
		"LexAlt":     State(24),
		"LexTerm":    State(25),
	},
	// state 19
	GotoRow{
		"SyntaxProduction": State(20),
	},
	// state 20
	GotoRow{},
	// state 21
	GotoRow{
		"Alternatives": State(35),
		"SyntaxBody":   State(36),
		"Symbols":      State(37),
		"Symbol":       State(40),
	},
	// state 22
	GotoRow{},
	// state 23
	GotoRow{},
	// state 24
	GotoRow{
		"LexTerm": State(44),
	},
	// state 25
	GotoRow{},
	// state 26
	GotoRow{},
	// state 27
	GotoRow{},
	// state 28
	GotoRow{
		"LexPattern": State(46),
		"LexAlt":     State(48),
		"LexTerm":    State(49),
	},
	// state 29
	GotoRow{
		"LexPattern": State(55),
		"LexAlt":     State(57),
		"LexTerm":    State(58),
	},
	// state 30
	GotoRow{
		"LexPattern": State(64),
		"LexAlt":     State(66),
		"LexTerm":    State(67),
	},
	// state 31
	GotoRow{},
	// state 32
	GotoRow{},
	// state 33
	GotoRow{},
	// state 34
	GotoRow{},
	// state 35
	GotoRow{},
	// state 36
	GotoRow{},
	// state 37
	GotoRow{
		"Symbol": State(78),
	},
	// state 38
	GotoRow{
		"Symbols": State(79),
		"Symbol":  State(40),
	},
	// state 39
	GotoRow{},
	// state 40
	GotoRow{},
	// state 41
	GotoRow{},
	// state 42
	GotoRow{},
	// state 43
	GotoRow{
		"LexAlt":  State(80),
		"LexTerm": State(25),
	},
	// state 44
	GotoRow{},
	// state 45
	GotoRow{},
	// state 46
	GotoRow{},
	// state 47
	GotoRow{},
	// state 48
	GotoRow{
		"LexTerm": State(84),
	},
	// state 49
	GotoRow{},
	// state 50
	GotoRow{},
	// state 51
	GotoRow{},
	// state 52
	GotoRow{
		"LexPattern": State(86),
		"LexAlt":     State(48),
		"LexTerm":    State(49),
	},
	// state 53
	GotoRow{
		"LexPattern": State(87),
		"LexAlt":     State(57),
		"LexTerm":    State(58),
	},
	// state 54
	GotoRow{
		"LexPattern": State(88),
		"LexAlt":     State(66),
		"LexTerm":    State(67),
	},
	// state 55
	GotoRow{},
	// state 56
	GotoRow{},
	// state 57
	GotoRow{
		"LexTerm": State(91),
	},
	// state 58
	GotoRow{},
	// state 59
	GotoRow{},
	// state 60
	GotoRow{},
	// state 61
	GotoRow{
		"LexPattern": State(93),
		"LexAlt":     State(48),
		"LexTerm":    State(49),
	},
	// state 62
	GotoRow{
		"LexPattern": State(94),
		"LexAlt":     State(57),
		"LexTerm":    State(58),
	},
	// state 63
	GotoRow{
		"LexPattern": State(95),
		"LexAlt":     State(66),
		"LexTerm":    State(67),
	},
	// state 64
	GotoRow{},
	// state 65
	GotoRow{},
	// state 66
	GotoRow{
		"LexTerm": State(98),
	},
	// state 67
	GotoRow{},
	// state 68
	GotoRow{},
	// state 69
	GotoRow{},
	// state 70
	GotoRow{
		"LexPattern": State(100),
		"LexAlt":     State(48),
		"LexTerm":    State(49),
	},
	// state 71
	GotoRow{
		"LexPattern": State(101),
		"LexAlt":     State(57),
		"LexTerm":    State(58),
	},
	// state 72
	GotoRow{
		"LexPattern": State(102),
		"LexAlt":     State(66),
		"LexTerm":    State(67),
	},
	// state 73
	GotoRow{},
	// state 74
	GotoRow{},
	// state 75
	GotoRow{},
	// state 76
	GotoRow{
		"SyntaxBody": State(103),
		"Symbols":    State(37),
		"Symbol":     State(40),
	},
	// state 77
	GotoRow{},
	// state 78
	GotoRow{},
	// state 79
	GotoRow{
		"Symbol": State(78),
	},
	// state 80
	GotoRow{
		"LexTerm": State(44),
	},
	// state 81
	GotoRow{},
	// state 82
	GotoRow{
		"LexAlt":  State(105),
		"LexTerm": State(49),
	},
	// state 83
	GotoRow{},
	// state 84
	GotoRow{},
	// state 85
	GotoRow{},
	// state 86
	GotoRow{},
	// state 87
	GotoRow{},
	// state 88
	GotoRow{},
	// state 89
	GotoRow{
		"LexAlt":  State(110),
		"LexTerm": State(58),
	},
	// state 90
	GotoRow{},
	// state 91
	GotoRow{},
	// state 92
	GotoRow{},
	// state 93
	GotoRow{},
	// state 94
	GotoRow{},
	// state 95
	GotoRow{},
	// state 96
	GotoRow{
		"LexAlt":  State(115),
		"LexTerm": State(67),
	},
	// state 97
	GotoRow{},
	// state 98
	GotoRow{},
	// state 99
	GotoRow{},
	// state 100
	GotoRow{},
	// state 101
	GotoRow{},
	// state 102
	GotoRow{},
	// state 103
	GotoRow{},
	// state 104
	GotoRow{},
	// state 105
	GotoRow{
		"LexTerm": State(84),
	},
	// state 106
	GotoRow{},
	// state 107
	GotoRow{},
	// state 108
	GotoRow{},
	// state 109
	GotoRow{},
	// state 110
	GotoRow{
		"LexTerm": State(91),
	},
	// state 111
	GotoRow{},
	// state 112
	GotoRow{},
	// state 113
	GotoRow{},
	// state 114
	GotoRow{},
	// state 115
	GotoRow{
		"LexTerm": State(98),
	},
	// state 116
	GotoRow{},
	// state 117
	GotoRow{},
	// state 118
	GotoRow{},
	// state 119
	GotoRow{},
}
