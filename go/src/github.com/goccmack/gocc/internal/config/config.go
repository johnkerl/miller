//Copyright 2013 Vastech SA (PTY) LTD
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package config

import (
	"errors"
	"flag"
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Config interface {
	Help() bool
	Verbose() bool
	Zip() bool
	AllowUnreachable() bool
	AutoResolveLRConf() bool
	SourceFile() string
	OutDir() string

	NoLexer() bool
	DebugLexer() bool
	DebugParser() bool

	ErrorsDir() string
	ParserDir() string
	ScannerDir() string
	TokenDir() string

	ProjectName() string
	Package() string

	PrintParams()
}

type ConfigRecord struct {
	workingDir string

	allowUnreachable  *bool
	autoResolveLRConf *bool
	debugLexer        *bool
	debugParser       *bool
	help              *bool
	noLexer           *bool
	outDir            string
	pkg               string
	srcFile           string
	verbose           *bool
	zip               *bool
}

func New() (Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cfg := &ConfigRecord{
		workingDir: wd,
	}

	if err := cfg.getFlags(); err != nil {
		cfg.PrintParams()
		return nil, err
	}

	return cfg, nil
}

func (this *ConfigRecord) Help() bool {
	return *this.help
}

func (this *ConfigRecord) Verbose() bool {
	return *this.verbose
}

func (this *ConfigRecord) Zip() bool {
	return *this.zip
}

func (this *ConfigRecord) AllowUnreachable() bool {
	return *this.allowUnreachable
}

func (this *ConfigRecord) AutoResolveLRConf() bool {
	return *this.autoResolveLRConf
}

func (this *ConfigRecord) NoLexer() bool {
	return *this.noLexer
}

func (this *ConfigRecord) DebugLexer() bool {
	return *this.debugLexer
}

func (this *ConfigRecord) DebugParser() bool {
	return *this.debugParser
}

func (this *ConfigRecord) SourceFile() string {
	return this.srcFile
}

func (this *ConfigRecord) OutDir() string {
	return this.outDir
}

func (this *ConfigRecord) ErrorsDir() string {
	return path.Join(this.outDir, "errors")
}

func (this *ConfigRecord) ParserDir() string {
	return path.Join(this.outDir, "parser")
}

func (this *ConfigRecord) ScannerDir() string {
	return path.Join(this.outDir, "scanner")
}

func (this *ConfigRecord) TokenDir() string {
	return path.Join(this.outDir, "token")
}

func (this *ConfigRecord) Package() string {
	return this.pkg
}

func (this *ConfigRecord) ProjectName() string {
	_, file := path.Split(this.pkg)
	return file
}

func (this *ConfigRecord) PrintParams() {
	fmt.Printf("-a             = %v\n", *this.autoResolveLRConf)
	fmt.Printf("-debug_lexer   = %v\n", *this.debugLexer)
	fmt.Printf("-debug_parser  = %v\n", *this.debugParser)
	fmt.Printf("-h             = %v\n", *this.help)
	fmt.Printf("-no_lexer      = %v\n", *this.noLexer)
	fmt.Printf("-o             = %v\n", this.outDir)
	fmt.Printf("-p             = %v\n", this.pkg)
	fmt.Printf("-u             = %v\n", *this.allowUnreachable)
	fmt.Printf("-v             = %v\n", *this.verbose)
	fmt.Printf("-zip           = %v\n", *this.zip)
}

/*** Utility routines ***/

func (this *ConfigRecord) getFlags() error {
	this.autoResolveLRConf = flag.Bool("a", false, "automatically resolve LR(1) conflicts")
	this.debugLexer = flag.Bool("debug_lexer", false, "enable debug logging in lexer")
	this.debugParser = flag.Bool("debug_parser", false, "enable debug logging in parser")
	this.help = flag.Bool("h", false, "help")
	this.noLexer = flag.Bool("no_lexer", false, "do not generate a lexer")
	flag.StringVar(&this.outDir, "o", this.workingDir, "output dir.")
	flag.StringVar(&this.pkg, "p", defaultPackage(this.outDir), "package")
	this.allowUnreachable = flag.Bool("u", false, "allow unreachable productions")
	this.verbose = flag.Bool("v", false, "verbose")
	this.zip = flag.Bool("zip", false, "zip the actiontable and gototable (experimental)")
	flag.Parse()

	if *this.noLexer && *this.debugLexer {
		return errors.New("no_lexer and debug_lexer cannot both be set")
	}

	this.outDir = getOutDir(this.outDir, this.workingDir)
	if this.outDir != this.workingDir {
		this.pkg = defaultPackage(this.outDir)
	}

	if len(flag.Args()) != 1 && !*this.help {
		return errors.New("Too few arguments")
	}

	this.srcFile = flag.Arg(0)

	return nil
}

func getOutDir(outDirSpec, wd string) string {
	if strings.HasPrefix(outDirSpec, wd) {
		return outDirSpec
	}
	return path.Join(wd, outDirSpec)
}

func defaultPackage(wd string) string {
	for _, srcDir := range build.Default.SrcDirs() {
		if strings.HasPrefix(wd, srcDir) {
			pkg := strings.Replace(wd, srcDir, "", -1)
			return sanitizePackage(pkg)
		}
	}
	return sanitizePackage(wd)
}

func sanitizePackage(pkg string) string {
	pkg = filepath.ToSlash(pkg)
	if strings.HasPrefix(pkg, "/") {
		pkg = pkg[1:]
	}
	return pkg
}
