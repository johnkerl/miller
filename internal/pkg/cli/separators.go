package cli

const COLON = ":"
const COMMA = ","
const CR = "\\r"
const CRCR = "\\r\\r"
const CRLF = "\\r\\n"
const CRLFCRLF = "\\r\\n\\r\\n"
const EQUALS = "="
const LF = "\\n"
const LFLF = "\\n\\n"
const NEWLINE = "\\n"
const PIPE = "|"
const SEMICOLON = ";"
const SLASH = "/"
const SPACE = " "
const SPACES = "( )+"
const TAB = "\\t"
const TABS = "(\\t)+"
const WHITESPACE = "([ \\t])+"

const ASCII_ESC = "\\x1b"
const ASCII_ETX = "\\x04"
const ASCII_FS = "\\x1c"
const ASCII_GS = "\\x1d"
const ASCII_NULL = "\\x01"
const ASCII_RS = "\\x1e"
const ASCII_SOH = "\\x02"
const ASCII_STX = "\\x03"
const ASCII_US = "\\x1f"

const ASV_FS = "\\x1f"
const ASV_RS = "\\x1e"
const USV_FS = "\\xe2\\x90\\x9f"
const USV_RS = "\\xe2\\x90\\x9e"

const ASV_FS_FOR_HELP = "\\x1f"
const ASV_RS_FOR_HELP = "\\x1e"
const USV_FS_FOR_HELP = "U+241F (UTF-8 \\xe2\\x90\\x9f)"
const USV_RS_FOR_HELP = "U+241E (UTF-8 \\xe2\\x90\\x9e)"

const DEFAULT_JSON_FLATTEN_SEPARATOR = "."

var SEPARATOR_NAMES_TO_VALUES = map[string]string{
	"ascii_esc":  ASCII_ESC,
	"ascii_etx":  ASCII_ETX,
	"ascii_fs":   ASCII_FS,
	"ascii_gs":   ASCII_GS,
	"ascii_null": ASCII_NULL,
	"ascii_rs":   ASCII_RS,
	"ascii_soh":  ASCII_SOH,
	"ascii_stx":  ASCII_STX,
	"ascii_us":   ASCII_US,
	"asv_fs":     ASV_FS,
	"asv_rs":     ASV_RS,
	"colon":      COLON,
	"comma":      COMMA,
	"cr":         CR,
	"crcr":       CRCR,
	"crlf":       CRLF,
	"crlfcrlf":   CRLFCRLF,
	"equals":     EQUALS,
	"lf":         LF,
	"lflf":       LFLF,
	"newline":    NEWLINE,
	"pipe":       PIPE,
	"semicolon":  SEMICOLON,
	"slash":      SLASH,
	"space":      SPACE,
	"spaces":     SPACES,
	"tab":        TAB,
	"tabs":       TABS,
	"usv_fs":     USV_FS,
	"usv_rs":     USV_RS,
	"whitespace": WHITESPACE,
}

// E.g. if IFS isn't specified, it's space for NIDX and comma for DKVP, etc.

var defaultFSes = map[string]string{
	"csv":      ",",
	"csvlite":  ",",
	"dkvp":     ",",
	"json":     "N/A", // not alterable; not parameterizable in JSON format
	"nidx":     " ",
	"markdown": " ",
	"pprint":   " ",
	"xtab":     "\n", // todo: windows-dependent ...
}

var defaultPSes = map[string]string{
	"csv":      "N/A",
	"csvlite":  "N/A",
	"dkvp":     "=",
	"json":     "N/A", // not alterable; not parameterizable in JSON format
	"markdown": "N/A",
	"nidx":     "N/A",
	"pprint":   "N/A",
	"xtab":     " ", // todo: windows-dependent ...
}

var defaultRSes = map[string]string{
	"csv":      "\n",
	"csvlite":  "\n",
	"dkvp":     "\n",
	"json":     "N/A", // not alterable; not parameterizable in JSON format
	"markdown": "\n",
	"nidx":     "\n",
	"pprint":   "\n",
	"xtab":     "\n\n", // todo: maybe jettison the idea of this being alterable
}

var defaultAllowRepeatIFSes = map[string]bool{
	"csv":      false,
	"csvlite":  false,
	"dkvp":     false,
	"json":     false,
	"markdown": false,
	"nidx":     false,
	"pprint":   true,
	"xtab":     false,
}

var defaultAllowRepeatIPSes = map[string]bool{
	"csv":      false,
	"csvlite":  false,
	"dkvp":     false,
	"json":     false,
	"markdown": false,
	"nidx":     false,
	"pprint":   false,
	"xtab":     true,
}
