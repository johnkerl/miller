package parser

import "github.com/goccmack/gocc/internal/ast"

const (
	NUM_STATES = 120
	NUM_TOKENS = 22
)

type (
	ActionTabU [NUM_STATES]ActionRowU
	ActionRowU [NUM_TOKENS]Action
)

type (
	CanRecover [NUM_STATES]bool
)

func getActionTableUncompressed() (act *ActionTabU) {
	act = new(ActionTabU)
	/* S0 */
	act[0][2] = Shift(6)   // tokId
	act[0][5] = Shift(7)   // regDefId
	act[0][6] = Shift(8)   // ignoredTokId
	act[0][17] = Shift(12) // prodId
	act[0][18] = Shift(13) // g_sdt_lit
	/* S1 */
	act[1][0] = Accept(0) // $
	/* S2 */
	act[2][0] = Reduce(2)  // $
	act[2][17] = Shift(12) // prodId
	act[2][18] = Shift(13) // g_sdt_lit
	/* S3 */
	act[3][0] = Reduce(3) // $
	/* S4 */
	act[4][0] = Reduce(4)  // $
	act[4][2] = Shift(6)   // tokId
	act[4][5] = Shift(7)   // regDefId
	act[4][6] = Shift(8)   // ignoredTokId
	act[4][17] = Reduce(4) // prodId
	act[4][18] = Reduce(4) // g_sdt_lit
	/* S5 */
	act[5][0] = Reduce(5)  // $
	act[5][2] = Reduce(5)  // tokId
	act[5][5] = Reduce(5)  // regDefId
	act[5][6] = Reduce(5)  // ignoredTokId
	act[5][17] = Reduce(5) // prodId
	act[5][18] = Reduce(5) // g_sdt_lit
	/* S6 */
	act[6][3] = Shift(16) // :
	/* S7 */
	act[7][3] = Shift(17) // :
	/* S8 */
	act[8][3] = Shift(18) // :
	/* S9 */
	act[9][17] = Shift(12) // prodId
	/* S10 */
	act[10][0] = Reduce(22) // $
	act[10][17] = Shift(12) // prodId
	/* S11 */
	act[11][0] = Reduce(23)  // $
	act[11][17] = Reduce(23) // prodId
	/* S12 */
	act[12][3] = Shift(21) // :
	/* S13 */
	act[13][17] = Reduce(39) // prodId
	/* S14 */
	act[14][0] = Reduce(1) // $
	/* S15 */
	act[15][0] = Reduce(6)  // $
	act[15][2] = Reduce(6)  // tokId
	act[15][5] = Reduce(6)  // regDefId
	act[15][6] = Reduce(6)  // ignoredTokId
	act[15][17] = Reduce(6) // prodId
	act[15][18] = Reduce(6) // g_sdt_lit
	/* S16 */
	act[16][5] = Shift(23)  // regDefId
	act[16][8] = Shift(26)  // .
	act[16][9] = Shift(27)  // char_lit
	act[16][11] = Shift(28) // [
	act[16][13] = Shift(29) // {
	act[16][15] = Shift(30) // (
	/* S17 */
	act[17][5] = Shift(23)  // regDefId
	act[17][8] = Shift(26)  // .
	act[17][9] = Shift(27)  // char_lit
	act[17][11] = Shift(28) // [
	act[17][13] = Shift(29) // {
	act[17][15] = Shift(30) // (
	/* S18 */
	act[18][5] = Shift(23)  // regDefId
	act[18][8] = Shift(26)  // .
	act[18][9] = Shift(27)  // char_lit
	act[18][11] = Shift(28) // [
	act[18][13] = Shift(29) // {
	act[18][15] = Shift(30) // (
	/* S19 */
	act[19][0] = Reduce(21) // $
	act[19][17] = Shift(12) // prodId
	/* S20 */
	act[20][0] = Reduce(24)  // $
	act[20][17] = Reduce(24) // prodId
	/* S21 */
	act[21][2] = Shift(33)  // tokId
	act[21][17] = Shift(34) // prodId
	act[21][19] = Shift(38) // error
	act[21][20] = Shift(39) // empty
	act[21][21] = Shift(41) // string_lit
	/* S22 */
	act[22][4] = Shift(42) // ;
	act[22][7] = Shift(43) // |
	/* S23 */
	act[23][4] = Reduce(17)  // ;
	act[23][5] = Reduce(17)  // regDefId
	act[23][7] = Reduce(17)  // |
	act[23][8] = Reduce(17)  // .
	act[23][9] = Reduce(17)  // char_lit
	act[23][11] = Reduce(17) // [
	act[23][13] = Reduce(17) // {
	act[23][15] = Reduce(17) // (
	/* S24 */
	act[24][4] = Reduce(10) // ;
	act[24][5] = Shift(23)  // regDefId
	act[24][7] = Reduce(10) // |
	act[24][8] = Shift(26)  // .
	act[24][9] = Shift(27)  // char_lit
	act[24][11] = Shift(28) // [
	act[24][13] = Shift(29) // {
	act[24][15] = Shift(30) // (
	/* S25 */
	act[25][4] = Reduce(12)  // ;
	act[25][5] = Reduce(12)  // regDefId
	act[25][7] = Reduce(12)  // |
	act[25][8] = Reduce(12)  // .
	act[25][9] = Reduce(12)  // char_lit
	act[25][11] = Reduce(12) // [
	act[25][13] = Reduce(12) // {
	act[25][15] = Reduce(12) // (
	/* S26 */
	act[26][4] = Reduce(14)  // ;
	act[26][5] = Reduce(14)  // regDefId
	act[26][7] = Reduce(14)  // |
	act[26][8] = Reduce(14)  // .
	act[26][9] = Reduce(14)  // char_lit
	act[26][11] = Reduce(14) // [
	act[26][13] = Reduce(14) // {
	act[26][15] = Reduce(14) // (
	/* S27 */
	act[27][4] = Reduce(15)  // ;
	act[27][5] = Reduce(15)  // regDefId
	act[27][7] = Reduce(15)  // |
	act[27][8] = Reduce(15)  // .
	act[27][9] = Reduce(15)  // char_lit
	act[27][10] = Shift(45)  // -
	act[27][11] = Reduce(15) // [
	act[27][13] = Reduce(15) // {
	act[27][15] = Reduce(15) // (
	/* S28 */
	act[28][5] = Shift(47)  // regDefId
	act[28][8] = Shift(50)  // .
	act[28][9] = Shift(51)  // char_lit
	act[28][11] = Shift(52) // [
	act[28][13] = Shift(53) // {
	act[28][15] = Shift(54) // (
	/* S29 */
	act[29][5] = Shift(56)  // regDefId
	act[29][8] = Shift(59)  // .
	act[29][9] = Shift(60)  // char_lit
	act[29][11] = Shift(61) // [
	act[29][13] = Shift(62) // {
	act[29][15] = Shift(63) // (
	/* S30 */
	act[30][5] = Shift(65)  // regDefId
	act[30][8] = Shift(68)  // .
	act[30][9] = Shift(69)  // char_lit
	act[30][11] = Shift(70) // [
	act[30][13] = Shift(71) // {
	act[30][15] = Shift(72) // (
	/* S31 */
	act[31][4] = Shift(73) // ;
	act[31][7] = Shift(43) // |
	/* S32 */
	act[32][4] = Shift(74) // ;
	act[32][7] = Shift(43) // |
	/* S33 */
	act[33][2] = Reduce(37)  // tokId
	act[33][4] = Reduce(37)  // ;
	act[33][7] = Reduce(37)  // |
	act[33][17] = Reduce(37) // prodId
	act[33][18] = Reduce(37) // g_sdt_lit
	act[33][21] = Reduce(37) // string_lit
	/* S34 */
	act[34][2] = Reduce(36)  // tokId
	act[34][4] = Reduce(36)  // ;
	act[34][7] = Reduce(36)  // |
	act[34][17] = Reduce(36) // prodId
	act[34][18] = Reduce(36) // g_sdt_lit
	act[34][21] = Reduce(36) // string_lit
	/* S35 */
	act[35][4] = Shift(75) // ;
	act[35][7] = Shift(76) // |
	/* S36 */
	act[36][4] = Reduce(26) // ;
	act[36][7] = Reduce(26) // |
	/* S37 */
	act[37][2] = Shift(33)  // tokId
	act[37][4] = Reduce(28) // ;
	act[37][7] = Reduce(28) // |
	act[37][17] = Shift(34) // prodId
	act[37][18] = Shift(77) // g_sdt_lit
	act[37][21] = Shift(41) // string_lit
	/* S38 */
	act[38][2] = Shift(33)  // tokId
	act[38][4] = Reduce(30) // ;
	act[38][7] = Reduce(30) // |
	act[38][17] = Shift(34) // prodId
	act[38][21] = Shift(41) // string_lit
	/* S39 */
	act[39][4] = Reduce(33) // ;
	act[39][7] = Reduce(33) // |
	/* S40 */
	act[40][2] = Reduce(34)  // tokId
	act[40][4] = Reduce(34)  // ;
	act[40][7] = Reduce(34)  // |
	act[40][17] = Reduce(34) // prodId
	act[40][18] = Reduce(34) // g_sdt_lit
	act[40][21] = Reduce(34) // string_lit
	/* S41 */
	act[41][2] = Reduce(38)  // tokId
	act[41][4] = Reduce(38)  // ;
	act[41][7] = Reduce(38)  // |
	act[41][17] = Reduce(38) // prodId
	act[41][18] = Reduce(38) // g_sdt_lit
	act[41][21] = Reduce(38) // string_lit
	/* S42 */
	act[42][0] = Reduce(7)  // $
	act[42][2] = Reduce(7)  // tokId
	act[42][5] = Reduce(7)  // regDefId
	act[42][6] = Reduce(7)  // ignoredTokId
	act[42][17] = Reduce(7) // prodId
	act[42][18] = Reduce(7) // g_sdt_lit
	/* S43 */
	act[43][5] = Shift(23)  // regDefId
	act[43][8] = Shift(26)  // .
	act[43][9] = Shift(27)  // char_lit
	act[43][11] = Shift(28) // [
	act[43][13] = Shift(29) // {
	act[43][15] = Shift(30) // (
	/* S44 */
	act[44][4] = Reduce(13)  // ;
	act[44][5] = Reduce(13)  // regDefId
	act[44][7] = Reduce(13)  // |
	act[44][8] = Reduce(13)  // .
	act[44][9] = Reduce(13)  // char_lit
	act[44][11] = Reduce(13) // [
	act[44][13] = Reduce(13) // {
	act[44][15] = Reduce(13) // (
	/* S45 */
	act[45][9] = Shift(81) // char_lit
	/* S46 */
	act[46][7] = Shift(82)  // |
	act[46][12] = Shift(83) // ]
	/* S47 */
	act[47][5] = Reduce(17)  // regDefId
	act[47][7] = Reduce(17)  // |
	act[47][8] = Reduce(17)  // .
	act[47][9] = Reduce(17)  // char_lit
	act[47][11] = Reduce(17) // [
	act[47][12] = Reduce(17) // ]
	act[47][13] = Reduce(17) // {
	act[47][15] = Reduce(17) // (
	/* S48 */
	act[48][5] = Shift(47)   // regDefId
	act[48][7] = Reduce(10)  // |
	act[48][8] = Shift(50)   // .
	act[48][9] = Shift(51)   // char_lit
	act[48][11] = Shift(52)  // [
	act[48][12] = Reduce(10) // ]
	act[48][13] = Shift(53)  // {
	act[48][15] = Shift(54)  // (
	/* S49 */
	act[49][5] = Reduce(12)  // regDefId
	act[49][7] = Reduce(12)  // |
	act[49][8] = Reduce(12)  // .
	act[49][9] = Reduce(12)  // char_lit
	act[49][11] = Reduce(12) // [
	act[49][12] = Reduce(12) // ]
	act[49][13] = Reduce(12) // {
	act[49][15] = Reduce(12) // (
	/* S50 */
	act[50][5] = Reduce(14)  // regDefId
	act[50][7] = Reduce(14)  // |
	act[50][8] = Reduce(14)  // .
	act[50][9] = Reduce(14)  // char_lit
	act[50][11] = Reduce(14) // [
	act[50][12] = Reduce(14) // ]
	act[50][13] = Reduce(14) // {
	act[50][15] = Reduce(14) // (
	/* S51 */
	act[51][5] = Reduce(15)  // regDefId
	act[51][7] = Reduce(15)  // |
	act[51][8] = Reduce(15)  // .
	act[51][9] = Reduce(15)  // char_lit
	act[51][10] = Shift(85)  // -
	act[51][11] = Reduce(15) // [
	act[51][12] = Reduce(15) // ]
	act[51][13] = Reduce(15) // {
	act[51][15] = Reduce(15) // (
	/* S52 */
	act[52][5] = Shift(47)  // regDefId
	act[52][8] = Shift(50)  // .
	act[52][9] = Shift(51)  // char_lit
	act[52][11] = Shift(52) // [
	act[52][13] = Shift(53) // {
	act[52][15] = Shift(54) // (
	/* S53 */
	act[53][5] = Shift(56)  // regDefId
	act[53][8] = Shift(59)  // .
	act[53][9] = Shift(60)  // char_lit
	act[53][11] = Shift(61) // [
	act[53][13] = Shift(62) // {
	act[53][15] = Shift(63) // (
	/* S54 */
	act[54][5] = Shift(65)  // regDefId
	act[54][8] = Shift(68)  // .
	act[54][9] = Shift(69)  // char_lit
	act[54][11] = Shift(70) // [
	act[54][13] = Shift(71) // {
	act[54][15] = Shift(72) // (
	/* S55 */
	act[55][7] = Shift(89)  // |
	act[55][14] = Shift(90) // }
	/* S56 */
	act[56][5] = Reduce(17)  // regDefId
	act[56][7] = Reduce(17)  // |
	act[56][8] = Reduce(17)  // .
	act[56][9] = Reduce(17)  // char_lit
	act[56][11] = Reduce(17) // [
	act[56][13] = Reduce(17) // {
	act[56][14] = Reduce(17) // }
	act[56][15] = Reduce(17) // (
	/* S57 */
	act[57][5] = Shift(56)   // regDefId
	act[57][7] = Reduce(10)  // |
	act[57][8] = Shift(59)   // .
	act[57][9] = Shift(60)   // char_lit
	act[57][11] = Shift(61)  // [
	act[57][13] = Shift(62)  // {
	act[57][14] = Reduce(10) // }
	act[57][15] = Shift(63)  // (
	/* S58 */
	act[58][5] = Reduce(12)  // regDefId
	act[58][7] = Reduce(12)  // |
	act[58][8] = Reduce(12)  // .
	act[58][9] = Reduce(12)  // char_lit
	act[58][11] = Reduce(12) // [
	act[58][13] = Reduce(12) // {
	act[58][14] = Reduce(12) // }
	act[58][15] = Reduce(12) // (
	/* S59 */
	act[59][5] = Reduce(14)  // regDefId
	act[59][7] = Reduce(14)  // |
	act[59][8] = Reduce(14)  // .
	act[59][9] = Reduce(14)  // char_lit
	act[59][11] = Reduce(14) // [
	act[59][13] = Reduce(14) // {
	act[59][14] = Reduce(14) // }
	act[59][15] = Reduce(14) // (
	/* S60 */
	act[60][5] = Reduce(15)  // regDefId
	act[60][7] = Reduce(15)  // |
	act[60][8] = Reduce(15)  // .
	act[60][9] = Reduce(15)  // char_lit
	act[60][10] = Shift(92)  // -
	act[60][11] = Reduce(15) // [
	act[60][13] = Reduce(15) // {
	act[60][14] = Reduce(15) // }
	act[60][15] = Reduce(15) // (
	/* S61 */
	act[61][5] = Shift(47)  // regDefId
	act[61][8] = Shift(50)  // .
	act[61][9] = Shift(51)  // char_lit
	act[61][11] = Shift(52) // [
	act[61][13] = Shift(53) // {
	act[61][15] = Shift(54) // (
	/* S62 */
	act[62][5] = Shift(56)  // regDefId
	act[62][8] = Shift(59)  // .
	act[62][9] = Shift(60)  // char_lit
	act[62][11] = Shift(61) // [
	act[62][13] = Shift(62) // {
	act[62][15] = Shift(63) // (
	/* S63 */
	act[63][5] = Shift(65)  // regDefId
	act[63][8] = Shift(68)  // .
	act[63][9] = Shift(69)  // char_lit
	act[63][11] = Shift(70) // [
	act[63][13] = Shift(71) // {
	act[63][15] = Shift(72) // (
	/* S64 */
	act[64][7] = Shift(96)  // |
	act[64][16] = Shift(97) // )
	/* S65 */
	act[65][5] = Reduce(17)  // regDefId
	act[65][7] = Reduce(17)  // |
	act[65][8] = Reduce(17)  // .
	act[65][9] = Reduce(17)  // char_lit
	act[65][11] = Reduce(17) // [
	act[65][13] = Reduce(17) // {
	act[65][15] = Reduce(17) // (
	act[65][16] = Reduce(17) // )
	/* S66 */
	act[66][5] = Shift(65)   // regDefId
	act[66][7] = Reduce(10)  // |
	act[66][8] = Shift(68)   // .
	act[66][9] = Shift(69)   // char_lit
	act[66][11] = Shift(70)  // [
	act[66][13] = Shift(71)  // {
	act[66][15] = Shift(72)  // (
	act[66][16] = Reduce(10) // )
	/* S67 */
	act[67][5] = Reduce(12)  // regDefId
	act[67][7] = Reduce(12)  // |
	act[67][8] = Reduce(12)  // .
	act[67][9] = Reduce(12)  // char_lit
	act[67][11] = Reduce(12) // [
	act[67][13] = Reduce(12) // {
	act[67][15] = Reduce(12) // (
	act[67][16] = Reduce(12) // )
	/* S68 */
	act[68][5] = Reduce(14)  // regDefId
	act[68][7] = Reduce(14)  // |
	act[68][8] = Reduce(14)  // .
	act[68][9] = Reduce(14)  // char_lit
	act[68][11] = Reduce(14) // [
	act[68][13] = Reduce(14) // {
	act[68][15] = Reduce(14) // (
	act[68][16] = Reduce(14) // )
	/* S69 */
	act[69][5] = Reduce(15)  // regDefId
	act[69][7] = Reduce(15)  // |
	act[69][8] = Reduce(15)  // .
	act[69][9] = Reduce(15)  // char_lit
	act[69][10] = Shift(99)  // -
	act[69][11] = Reduce(15) // [
	act[69][13] = Reduce(15) // {
	act[69][15] = Reduce(15) // (
	act[69][16] = Reduce(15) // )
	/* S70 */
	act[70][5] = Shift(47)  // regDefId
	act[70][8] = Shift(50)  // .
	act[70][9] = Shift(51)  // char_lit
	act[70][11] = Shift(52) // [
	act[70][13] = Shift(53) // {
	act[70][15] = Shift(54) // (
	/* S71 */
	act[71][5] = Shift(56)  // regDefId
	act[71][8] = Shift(59)  // .
	act[71][9] = Shift(60)  // char_lit
	act[71][11] = Shift(61) // [
	act[71][13] = Shift(62) // {
	act[71][15] = Shift(63) // (
	/* S72 */
	act[72][5] = Shift(65)  // regDefId
	act[72][8] = Shift(68)  // .
	act[72][9] = Shift(69)  // char_lit
	act[72][11] = Shift(70) // [
	act[72][13] = Shift(71) // {
	act[72][15] = Shift(72) // (
	/* S73 */
	act[73][0] = Reduce(8)  // $
	act[73][2] = Reduce(8)  // tokId
	act[73][5] = Reduce(8)  // regDefId
	act[73][6] = Reduce(8)  // ignoredTokId
	act[73][17] = Reduce(8) // prodId
	act[73][18] = Reduce(8) // g_sdt_lit
	/* S74 */
	act[74][0] = Reduce(9)  // $
	act[74][2] = Reduce(9)  // tokId
	act[74][5] = Reduce(9)  // regDefId
	act[74][6] = Reduce(9)  // ignoredTokId
	act[74][17] = Reduce(9) // prodId
	act[74][18] = Reduce(9) // g_sdt_lit
	/* S75 */
	act[75][0] = Reduce(25)  // $
	act[75][17] = Reduce(25) // prodId
	/* S76 */
	act[76][2] = Shift(33)  // tokId
	act[76][17] = Shift(34) // prodId
	act[76][19] = Shift(38) // error
	act[76][20] = Shift(39) // empty
	act[76][21] = Shift(41) // string_lit
	/* S77 */
	act[77][4] = Reduce(29) // ;
	act[77][7] = Reduce(29) // |
	/* S78 */
	act[78][2] = Reduce(35)  // tokId
	act[78][4] = Reduce(35)  // ;
	act[78][7] = Reduce(35)  // |
	act[78][17] = Reduce(35) // prodId
	act[78][18] = Reduce(35) // g_sdt_lit
	act[78][21] = Reduce(35) // string_lit
	/* S79 */
	act[79][2] = Shift(33)   // tokId
	act[79][4] = Reduce(31)  // ;
	act[79][7] = Reduce(31)  // |
	act[79][17] = Shift(34)  // prodId
	act[79][18] = Shift(104) // g_sdt_lit
	act[79][21] = Shift(41)  // string_lit
	/* S80 */
	act[80][4] = Reduce(11) // ;
	act[80][5] = Shift(23)  // regDefId
	act[80][7] = Reduce(11) // |
	act[80][8] = Shift(26)  // .
	act[80][9] = Shift(27)  // char_lit
	act[80][11] = Shift(28) // [
	act[80][13] = Shift(29) // {
	act[80][15] = Shift(30) // (
	/* S81 */
	act[81][4] = Reduce(16)  // ;
	act[81][5] = Reduce(16)  // regDefId
	act[81][7] = Reduce(16)  // |
	act[81][8] = Reduce(16)  // .
	act[81][9] = Reduce(16)  // char_lit
	act[81][11] = Reduce(16) // [
	act[81][13] = Reduce(16) // {
	act[81][15] = Reduce(16) // (
	/* S82 */
	act[82][5] = Shift(47)  // regDefId
	act[82][8] = Shift(50)  // .
	act[82][9] = Shift(51)  // char_lit
	act[82][11] = Shift(52) // [
	act[82][13] = Shift(53) // {
	act[82][15] = Shift(54) // (
	/* S83 */
	act[83][4] = Reduce(18)  // ;
	act[83][5] = Reduce(18)  // regDefId
	act[83][7] = Reduce(18)  // |
	act[83][8] = Reduce(18)  // .
	act[83][9] = Reduce(18)  // char_lit
	act[83][11] = Reduce(18) // [
	act[83][13] = Reduce(18) // {
	act[83][15] = Reduce(18) // (
	/* S84 */
	act[84][5] = Reduce(13)  // regDefId
	act[84][7] = Reduce(13)  // |
	act[84][8] = Reduce(13)  // .
	act[84][9] = Reduce(13)  // char_lit
	act[84][11] = Reduce(13) // [
	act[84][12] = Reduce(13) // ]
	act[84][13] = Reduce(13) // {
	act[84][15] = Reduce(13) // (
	/* S85 */
	act[85][9] = Shift(106) // char_lit
	/* S86 */
	act[86][7] = Shift(82)   // |
	act[86][12] = Shift(107) // ]
	/* S87 */
	act[87][7] = Shift(89)   // |
	act[87][14] = Shift(108) // }
	/* S88 */
	act[88][7] = Shift(96)   // |
	act[88][16] = Shift(109) // )
	/* S89 */
	act[89][5] = Shift(56)  // regDefId
	act[89][8] = Shift(59)  // .
	act[89][9] = Shift(60)  // char_lit
	act[89][11] = Shift(61) // [
	act[89][13] = Shift(62) // {
	act[89][15] = Shift(63) // (
	/* S90 */
	act[90][4] = Reduce(19)  // ;
	act[90][5] = Reduce(19)  // regDefId
	act[90][7] = Reduce(19)  // |
	act[90][8] = Reduce(19)  // .
	act[90][9] = Reduce(19)  // char_lit
	act[90][11] = Reduce(19) // [
	act[90][13] = Reduce(19) // {
	act[90][15] = Reduce(19) // (
	/* S91 */
	act[91][5] = Reduce(13)  // regDefId
	act[91][7] = Reduce(13)  // |
	act[91][8] = Reduce(13)  // .
	act[91][9] = Reduce(13)  // char_lit
	act[91][11] = Reduce(13) // [
	act[91][13] = Reduce(13) // {
	act[91][14] = Reduce(13) // }
	act[91][15] = Reduce(13) // (
	/* S92 */
	act[92][9] = Shift(111) // char_lit
	/* S93 */
	act[93][7] = Shift(82)   // |
	act[93][12] = Shift(112) // ]
	/* S94 */
	act[94][7] = Shift(89)   // |
	act[94][14] = Shift(113) // }
	/* S95 */
	act[95][7] = Shift(96)   // |
	act[95][16] = Shift(114) // )
	/* S96 */
	act[96][5] = Shift(65)  // regDefId
	act[96][8] = Shift(68)  // .
	act[96][9] = Shift(69)  // char_lit
	act[96][11] = Shift(70) // [
	act[96][13] = Shift(71) // {
	act[96][15] = Shift(72) // (
	/* S97 */
	act[97][4] = Reduce(20)  // ;
	act[97][5] = Reduce(20)  // regDefId
	act[97][7] = Reduce(20)  // |
	act[97][8] = Reduce(20)  // .
	act[97][9] = Reduce(20)  // char_lit
	act[97][11] = Reduce(20) // [
	act[97][13] = Reduce(20) // {
	act[97][15] = Reduce(20) // (
	/* S98 */
	act[98][5] = Reduce(13)  // regDefId
	act[98][7] = Reduce(13)  // |
	act[98][8] = Reduce(13)  // .
	act[98][9] = Reduce(13)  // char_lit
	act[98][11] = Reduce(13) // [
	act[98][13] = Reduce(13) // {
	act[98][15] = Reduce(13) // (
	act[98][16] = Reduce(13) // )
	/* S99 */
	act[99][9] = Shift(116) // char_lit
	/* S100 */
	act[100][7] = Shift(82)   // |
	act[100][12] = Shift(117) // ]
	/* S101 */
	act[101][7] = Shift(89)   // |
	act[101][14] = Shift(118) // }
	/* S102 */
	act[102][7] = Shift(96)   // |
	act[102][16] = Shift(119) // )
	/* S103 */
	act[103][4] = Reduce(27) // ;
	act[103][7] = Reduce(27) // |
	/* S104 */
	act[104][4] = Reduce(32) // ;
	act[104][7] = Reduce(32) // |
	/* S105 */
	act[105][5] = Shift(47)   // regDefId
	act[105][7] = Reduce(11)  // |
	act[105][8] = Shift(50)   // .
	act[105][9] = Shift(51)   // char_lit
	act[105][11] = Shift(52)  // [
	act[105][12] = Reduce(11) // ]
	act[105][13] = Shift(53)  // {
	act[105][15] = Shift(54)  // (
	/* S106 */
	act[106][5] = Reduce(16)  // regDefId
	act[106][7] = Reduce(16)  // |
	act[106][8] = Reduce(16)  // .
	act[106][9] = Reduce(16)  // char_lit
	act[106][11] = Reduce(16) // [
	act[106][12] = Reduce(16) // ]
	act[106][13] = Reduce(16) // {
	act[106][15] = Reduce(16) // (
	/* S107 */
	act[107][5] = Reduce(18)  // regDefId
	act[107][7] = Reduce(18)  // |
	act[107][8] = Reduce(18)  // .
	act[107][9] = Reduce(18)  // char_lit
	act[107][11] = Reduce(18) // [
	act[107][12] = Reduce(18) // ]
	act[107][13] = Reduce(18) // {
	act[107][15] = Reduce(18) // (
	/* S108 */
	act[108][5] = Reduce(19)  // regDefId
	act[108][7] = Reduce(19)  // |
	act[108][8] = Reduce(19)  // .
	act[108][9] = Reduce(19)  // char_lit
	act[108][11] = Reduce(19) // [
	act[108][12] = Reduce(19) // ]
	act[108][13] = Reduce(19) // {
	act[108][15] = Reduce(19) // (
	/* S109 */
	act[109][5] = Reduce(20)  // regDefId
	act[109][7] = Reduce(20)  // |
	act[109][8] = Reduce(20)  // .
	act[109][9] = Reduce(20)  // char_lit
	act[109][11] = Reduce(20) // [
	act[109][12] = Reduce(20) // ]
	act[109][13] = Reduce(20) // {
	act[109][15] = Reduce(20) // (
	/* S110 */
	act[110][5] = Shift(56)   // regDefId
	act[110][7] = Reduce(11)  // |
	act[110][8] = Shift(59)   // .
	act[110][9] = Shift(60)   // char_lit
	act[110][11] = Shift(61)  // [
	act[110][13] = Shift(62)  // {
	act[110][14] = Reduce(11) // }
	act[110][15] = Shift(63)  // (
	/* S111 */
	act[111][5] = Reduce(16)  // regDefId
	act[111][7] = Reduce(16)  // |
	act[111][8] = Reduce(16)  // .
	act[111][9] = Reduce(16)  // char_lit
	act[111][11] = Reduce(16) // [
	act[111][13] = Reduce(16) // {
	act[111][14] = Reduce(16) // }
	act[111][15] = Reduce(16) // (
	/* S112 */
	act[112][5] = Reduce(18)  // regDefId
	act[112][7] = Reduce(18)  // |
	act[112][8] = Reduce(18)  // .
	act[112][9] = Reduce(18)  // char_lit
	act[112][11] = Reduce(18) // [
	act[112][13] = Reduce(18) // {
	act[112][14] = Reduce(18) // }
	act[112][15] = Reduce(18) // (
	/* S113 */
	act[113][5] = Reduce(19)  // regDefId
	act[113][7] = Reduce(19)  // |
	act[113][8] = Reduce(19)  // .
	act[113][9] = Reduce(19)  // char_lit
	act[113][11] = Reduce(19) // [
	act[113][13] = Reduce(19) // {
	act[113][14] = Reduce(19) // }
	act[113][15] = Reduce(19) // (
	/* S114 */
	act[114][5] = Reduce(20)  // regDefId
	act[114][7] = Reduce(20)  // |
	act[114][8] = Reduce(20)  // .
	act[114][9] = Reduce(20)  // char_lit
	act[114][11] = Reduce(20) // [
	act[114][13] = Reduce(20) // {
	act[114][14] = Reduce(20) // }
	act[114][15] = Reduce(20) // (
	/* S115 */
	act[115][5] = Shift(65)   // regDefId
	act[115][7] = Reduce(11)  // |
	act[115][8] = Shift(68)   // .
	act[115][9] = Shift(69)   // char_lit
	act[115][11] = Shift(70)  // [
	act[115][13] = Shift(71)  // {
	act[115][15] = Shift(72)  // (
	act[115][16] = Reduce(11) // )
	/* S116 */
	act[116][5] = Reduce(16)  // regDefId
	act[116][7] = Reduce(16)  // |
	act[116][8] = Reduce(16)  // .
	act[116][9] = Reduce(16)  // char_lit
	act[116][11] = Reduce(16) // [
	act[116][13] = Reduce(16) // {
	act[116][15] = Reduce(16) // (
	act[116][16] = Reduce(16) // )
	/* S117 */
	act[117][5] = Reduce(18)  // regDefId
	act[117][7] = Reduce(18)  // |
	act[117][8] = Reduce(18)  // .
	act[117][9] = Reduce(18)  // char_lit
	act[117][11] = Reduce(18) // [
	act[117][13] = Reduce(18) // {
	act[117][15] = Reduce(18) // (
	act[117][16] = Reduce(18) // )
	/* S118 */
	act[118][5] = Reduce(19)  // regDefId
	act[118][7] = Reduce(19)  // |
	act[118][8] = Reduce(19)  // .
	act[118][9] = Reduce(19)  // char_lit
	act[118][11] = Reduce(19) // [
	act[118][13] = Reduce(19) // {
	act[118][15] = Reduce(19) // (
	act[118][16] = Reduce(19) // )
	/* S119 */
	act[119][5] = Reduce(20)  // regDefId
	act[119][7] = Reduce(20)  // |
	act[119][8] = Reduce(20)  // .
	act[119][9] = Reduce(20)  // char_lit
	act[119][11] = Reduce(20) // [
	act[119][13] = Reduce(20) // {
	act[119][15] = Reduce(20) // (
	act[119][16] = Reduce(20) // )
	return
}

func getCanRecoverTableUncompressed() (cr *CanRecover) {
	cr = new(CanRecover)
	return
}

// NT is the set of non-terminal symbols of the target grammar
const (
	NUM_NT = 15
)

type (
	GotoTabU [NUM_STATES]GotoRowU
	GotoRowU [NUM_NT]State
)

const (
	NT_Grammar          = 0
	NT_LexicalPart      = 1
	NT_LexProductions   = 2
	NT_LexProduction    = 3
	NT_LexPattern       = 4
	NT_LexAlt           = 5
	NT_LexTerm          = 6
	NT_SyntaxPart       = 7
	NT_SyntaxProdList   = 8
	NT_SyntaxProduction = 9
	NT_Alternatives     = 10
	NT_SyntaxBody       = 11
	NT_Symbols          = 12
	NT_Symbol           = 13
	NT_FileHeader       = 14
)

func getGotoTableUncompressed() (gto *GotoTabU) {
	gto = new(GotoTabU)
	gto[0][NT_Grammar] = 1
	gto[0][NT_LexicalPart] = 2
	gto[0][NT_LexProductions] = 4
	gto[0][NT_LexProduction] = 5
	gto[0][NT_SyntaxPart] = 3
	gto[0][NT_SyntaxProdList] = 10
	gto[0][NT_SyntaxProduction] = 11
	gto[0][NT_FileHeader] = 9
	gto[2][NT_SyntaxPart] = 14
	gto[2][NT_SyntaxProdList] = 10
	gto[2][NT_SyntaxProduction] = 11
	gto[2][NT_FileHeader] = 9
	gto[4][NT_LexProduction] = 15
	gto[9][NT_SyntaxProdList] = 19
	gto[9][NT_SyntaxProduction] = 11
	gto[10][NT_SyntaxProduction] = 20
	gto[16][NT_LexPattern] = 22
	gto[16][NT_LexAlt] = 24
	gto[16][NT_LexTerm] = 25
	gto[17][NT_LexPattern] = 31
	gto[17][NT_LexAlt] = 24
	gto[17][NT_LexTerm] = 25
	gto[18][NT_LexPattern] = 32
	gto[18][NT_LexAlt] = 24
	gto[18][NT_LexTerm] = 25
	gto[19][NT_SyntaxProduction] = 20
	gto[21][NT_Alternatives] = 35
	gto[21][NT_SyntaxBody] = 36
	gto[21][NT_Symbols] = 37
	gto[21][NT_Symbol] = 40
	gto[24][NT_LexTerm] = 44
	gto[28][NT_LexPattern] = 46
	gto[28][NT_LexAlt] = 48
	gto[28][NT_LexTerm] = 49
	gto[29][NT_LexPattern] = 55
	gto[29][NT_LexAlt] = 57
	gto[29][NT_LexTerm] = 58
	gto[30][NT_LexPattern] = 64
	gto[30][NT_LexAlt] = 66
	gto[30][NT_LexTerm] = 67
	gto[37][NT_Symbol] = 78
	gto[38][NT_Symbols] = 79
	gto[38][NT_Symbol] = 40
	gto[43][NT_LexAlt] = 80
	gto[43][NT_LexTerm] = 25
	gto[48][NT_LexTerm] = 84
	gto[52][NT_LexPattern] = 86
	gto[52][NT_LexAlt] = 48
	gto[52][NT_LexTerm] = 49
	gto[53][NT_LexPattern] = 87
	gto[53][NT_LexAlt] = 57
	gto[53][NT_LexTerm] = 58
	gto[54][NT_LexPattern] = 88
	gto[54][NT_LexAlt] = 66
	gto[54][NT_LexTerm] = 67
	gto[57][NT_LexTerm] = 91
	gto[61][NT_LexPattern] = 93
	gto[61][NT_LexAlt] = 48
	gto[61][NT_LexTerm] = 49
	gto[62][NT_LexPattern] = 94
	gto[62][NT_LexAlt] = 57
	gto[62][NT_LexTerm] = 58
	gto[63][NT_LexPattern] = 95
	gto[63][NT_LexAlt] = 66
	gto[63][NT_LexTerm] = 67
	gto[66][NT_LexTerm] = 98
	gto[70][NT_LexPattern] = 100
	gto[70][NT_LexAlt] = 48
	gto[70][NT_LexTerm] = 49
	gto[71][NT_LexPattern] = 101
	gto[71][NT_LexAlt] = 57
	gto[71][NT_LexTerm] = 58
	gto[72][NT_LexPattern] = 102
	gto[72][NT_LexAlt] = 66
	gto[72][NT_LexTerm] = 67
	gto[76][NT_SyntaxBody] = 103
	gto[76][NT_Symbols] = 37
	gto[76][NT_Symbol] = 40
	gto[79][NT_Symbol] = 78
	gto[80][NT_LexTerm] = 44
	gto[82][NT_LexAlt] = 105
	gto[82][NT_LexTerm] = 49
	gto[89][NT_LexAlt] = 110
	gto[89][NT_LexTerm] = 58
	gto[96][NT_LexAlt] = 115
	gto[96][NT_LexTerm] = 67
	gto[105][NT_LexTerm] = 84
	gto[110][NT_LexTerm] = 91
	gto[115][NT_LexTerm] = 98
	return
}

const NUM_PRODS = 40

type (
	ProdTabU      [NUM_PRODS]*ProdTabUEntry
	ProdTabUEntry struct {
		String     string
		Head       NT
		HeadIndex  int
		NumSymbols int
		ReduceFunc func([]Attrib) (Attrib, error)
	}
)

func getProductionsTableUncompressed() (pt *ProdTabU) {
	pt = new(ProdTabU)
	pt[0] = &ProdTabUEntry{
		String:     "S! : Grammar ;",
		Head:       "S!",
		HeadIndex:  -1,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return X[0], nil
		},
	}
	pt[1] = &ProdTabUEntry{
		String:     "Grammar : LexicalPart SyntaxPart << ast.NewGrammar(X[0], X[1]) >> ;",
		Head:       "Grammar",
		HeadIndex:  NT_Grammar,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewGrammar(X[0], X[1])
		},
	}
	pt[2] = &ProdTabUEntry{
		String:     "Grammar : LexicalPart << ast.NewGrammar(X[0], nil) >> ;",
		Head:       "Grammar",
		HeadIndex:  NT_Grammar,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewGrammar(X[0], nil)
		},
	}
	pt[3] = &ProdTabUEntry{
		String:     "Grammar : SyntaxPart << ast.NewGrammar(nil, X[0]) >> ;",
		Head:       "Grammar",
		HeadIndex:  NT_Grammar,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewGrammar(nil, X[0])
		},
	}
	pt[4] = &ProdTabUEntry{
		String:     "LexicalPart : LexProductions << ast.NewLexPart(nil, nil, X[0]) >> ;",
		Head:       "LexicalPart",
		HeadIndex:  NT_LexicalPart,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexPart(nil, nil, X[0])
		},
	}
	pt[5] = &ProdTabUEntry{
		String:     "LexProductions : LexProduction << ast.NewLexProductions(X[0]) >> ;",
		Head:       "LexProductions",
		HeadIndex:  NT_LexProductions,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexProductions(X[0])
		},
	}
	pt[6] = &ProdTabUEntry{
		String:     "LexProductions : LexProductions LexProduction << ast.AppendLexProduction(X[0], X[1]) >> ;",
		Head:       "LexProductions",
		HeadIndex:  NT_LexProductions,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.AppendLexProduction(X[0], X[1])
		},
	}
	pt[7] = &ProdTabUEntry{
		String:     "LexProduction : tokId : LexPattern ; << ast.NewLexTokDef(X[0], X[2]) >> ;",
		Head:       "LexProduction",
		HeadIndex:  NT_LexProduction,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexTokDef(X[0], X[2])
		},
	}
	pt[8] = &ProdTabUEntry{
		String:     "LexProduction : regDefId : LexPattern ; << ast.NewLexRegDef(X[0], X[2]) >> ;",
		Head:       "LexProduction",
		HeadIndex:  NT_LexProduction,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexRegDef(X[0], X[2])
		},
	}
	pt[9] = &ProdTabUEntry{
		String:     "LexProduction : ignoredTokId : LexPattern ; << ast.NewLexIgnoredTokDef(X[0], X[2]) >> ;",
		Head:       "LexProduction",
		HeadIndex:  NT_LexProduction,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexIgnoredTokDef(X[0], X[2])
		},
	}
	pt[10] = &ProdTabUEntry{
		String:     "LexPattern : LexAlt << ast.NewLexPattern(X[0]) >> ;",
		Head:       "LexPattern",
		HeadIndex:  NT_LexPattern,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexPattern(X[0])
		},
	}
	pt[11] = &ProdTabUEntry{
		String:     "LexPattern : LexPattern | LexAlt << ast.AppendLexAlt(X[0], X[2]) >> ;",
		Head:       "LexPattern",
		HeadIndex:  NT_LexPattern,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.AppendLexAlt(X[0], X[2])
		},
	}
	pt[12] = &ProdTabUEntry{
		String:     "LexAlt : LexTerm << ast.NewLexAlt(X[0]) >> ;",
		Head:       "LexAlt",
		HeadIndex:  NT_LexAlt,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexAlt(X[0])
		},
	}
	pt[13] = &ProdTabUEntry{
		String:     "LexAlt : LexAlt LexTerm << ast.AppendLexTerm(X[0], X[1]) >> ;",
		Head:       "LexAlt",
		HeadIndex:  NT_LexAlt,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.AppendLexTerm(X[0], X[1])
		},
	}
	pt[14] = &ProdTabUEntry{
		String:     "LexTerm : . << ast.LexDOT, nil >> ;",
		Head:       "LexTerm",
		HeadIndex:  NT_LexTerm,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.LexDOT, nil
		},
	}
	pt[15] = &ProdTabUEntry{
		String:     "LexTerm : char_lit << ast.NewLexCharLit(X[0]) >> ;",
		Head:       "LexTerm",
		HeadIndex:  NT_LexTerm,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexCharLit(X[0])
		},
	}
	pt[16] = &ProdTabUEntry{
		String:     "LexTerm : char_lit - char_lit << ast.NewLexCharRange(X[0], X[2]) >> ;",
		Head:       "LexTerm",
		HeadIndex:  NT_LexTerm,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexCharRange(X[0], X[2])
		},
	}
	pt[17] = &ProdTabUEntry{
		String:     "LexTerm : regDefId << ast.NewLexRegDefId(X[0]) >> ;",
		Head:       "LexTerm",
		HeadIndex:  NT_LexTerm,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexRegDefId(X[0])
		},
	}
	pt[18] = &ProdTabUEntry{
		String:     "LexTerm : [ LexPattern ] << ast.NewLexOptPattern(X[1]) >> ;",
		Head:       "LexTerm",
		HeadIndex:  NT_LexTerm,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexOptPattern(X[1])
		},
	}
	pt[19] = &ProdTabUEntry{
		String:     "LexTerm : { LexPattern } << ast.NewLexRepPattern(X[1]) >> ;",
		Head:       "LexTerm",
		HeadIndex:  NT_LexTerm,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexRepPattern(X[1])
		},
	}
	pt[20] = &ProdTabUEntry{
		String:     "LexTerm : ( LexPattern ) << ast.NewLexGroupPattern(X[1]) >> ;",
		Head:       "LexTerm",
		HeadIndex:  NT_LexTerm,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewLexGroupPattern(X[1])
		},
	}
	pt[21] = &ProdTabUEntry{
		String:     "SyntaxPart : FileHeader SyntaxProdList << ast.NewSyntaxPart(X[0], X[1]) >> ;",
		Head:       "SyntaxPart",
		HeadIndex:  NT_SyntaxPart,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxPart(X[0], X[1])
		},
	}
	pt[22] = &ProdTabUEntry{
		String:     "SyntaxPart : SyntaxProdList << ast.NewSyntaxPart(nil, X[0]) >> ;",
		Head:       "SyntaxPart",
		HeadIndex:  NT_SyntaxPart,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxPart(nil, X[0])
		},
	}
	pt[23] = &ProdTabUEntry{
		String:     "SyntaxProdList : SyntaxProduction << ast.NewSyntaxProdList(X[0]) >> ;",
		Head:       "SyntaxProdList",
		HeadIndex:  NT_SyntaxProdList,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxProdList(X[0])
		},
	}
	pt[24] = &ProdTabUEntry{
		String:     "SyntaxProdList : SyntaxProdList SyntaxProduction << ast.AddSyntaxProds(X[0], X[1]) >> ;",
		Head:       "SyntaxProdList",
		HeadIndex:  NT_SyntaxProdList,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.AddSyntaxProds(X[0], X[1])
		},
	}
	pt[25] = &ProdTabUEntry{
		String:     "SyntaxProduction : prodId : Alternatives ; << ast.NewSyntaxProd(X[0], X[2]) >> ;",
		Head:       "SyntaxProduction",
		HeadIndex:  NT_SyntaxProduction,
		NumSymbols: 4,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxProd(X[0], X[2])
		},
	}
	pt[26] = &ProdTabUEntry{
		String:     "Alternatives : SyntaxBody << ast.NewSyntaxAlts(X[0]) >> ;",
		Head:       "Alternatives",
		HeadIndex:  NT_Alternatives,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxAlts(X[0])
		},
	}
	pt[27] = &ProdTabUEntry{
		String:     "Alternatives : Alternatives | SyntaxBody << ast.AddSyntaxAlt(X[0], X[2]) >> ;",
		Head:       "Alternatives",
		HeadIndex:  NT_Alternatives,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.AddSyntaxAlt(X[0], X[2])
		},
	}
	pt[28] = &ProdTabUEntry{
		String:     "SyntaxBody : Symbols << ast.NewSyntaxBody(X[0], nil) >> ;",
		Head:       "SyntaxBody",
		HeadIndex:  NT_SyntaxBody,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxBody(X[0], nil)
		},
	}
	pt[29] = &ProdTabUEntry{
		String:     "SyntaxBody : Symbols g_sdt_lit << ast.NewSyntaxBody(X[0], X[1]) >> ;",
		Head:       "SyntaxBody",
		HeadIndex:  NT_SyntaxBody,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxBody(X[0], X[1])
		},
	}
	pt[30] = &ProdTabUEntry{
		String:     "SyntaxBody : error << ast.NewErrorBody(nil, nil) >> ;",
		Head:       "SyntaxBody",
		HeadIndex:  NT_SyntaxBody,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewErrorBody(nil, nil)
		},
	}
	pt[31] = &ProdTabUEntry{
		String:     "SyntaxBody : error Symbols << ast.NewErrorBody(X[1], nil) >> ;",
		Head:       "SyntaxBody",
		HeadIndex:  NT_SyntaxBody,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewErrorBody(X[1], nil)
		},
	}
	pt[32] = &ProdTabUEntry{
		String:     "SyntaxBody : error Symbols g_sdt_lit << ast.NewErrorBody(X[1], X[2]) >> ;",
		Head:       "SyntaxBody",
		HeadIndex:  NT_SyntaxBody,
		NumSymbols: 3,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewErrorBody(X[1], X[2])
		},
	}
	pt[33] = &ProdTabUEntry{
		String:     "SyntaxBody : empty << ast.NewEmptyBody() >> ;",
		Head:       "SyntaxBody",
		HeadIndex:  NT_SyntaxBody,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewEmptyBody()
		},
	}
	pt[34] = &ProdTabUEntry{
		String:     "Symbols : Symbol << ast.NewSyntaxSymbols(X[0]) >> ;",
		Head:       "Symbols",
		HeadIndex:  NT_Symbols,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxSymbols(X[0])
		},
	}
	pt[35] = &ProdTabUEntry{
		String:     "Symbols : Symbols Symbol << ast.AddSyntaxSymbol(X[0], X[1]) >> ;",
		Head:       "Symbols",
		HeadIndex:  NT_Symbols,
		NumSymbols: 2,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.AddSyntaxSymbol(X[0], X[1])
		},
	}
	pt[36] = &ProdTabUEntry{
		String:     "Symbol : prodId << ast.NewSyntaxProdId(X[0]) >> ;",
		Head:       "Symbol",
		HeadIndex:  NT_Symbol,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewSyntaxProdId(X[0])
		},
	}
	pt[37] = &ProdTabUEntry{
		String:     "Symbol : tokId << ast.NewTokId(X[0]) >> ;",
		Head:       "Symbol",
		HeadIndex:  NT_Symbol,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewTokId(X[0])
		},
	}
	pt[38] = &ProdTabUEntry{
		String:     "Symbol : string_lit << ast.NewStringLit(X[0]) >> ;",
		Head:       "Symbol",
		HeadIndex:  NT_Symbol,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewStringLit(X[0])
		},
	}
	pt[39] = &ProdTabUEntry{
		String:     "FileHeader : g_sdt_lit << ast.NewFileHeader(X[0]) >> ;",
		Head:       "FileHeader",
		HeadIndex:  NT_FileHeader,
		NumSymbols: 1,
		ReduceFunc: func(X []Attrib) (Attrib, error) {
			return ast.NewFileHeader(X[0])
		},
	}
	return
}
