package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/koron/textfmt/convert"
	"github.com/koron/textfmt/fmtinfo"
)

type enc int

const (
	noenc enc = iota
	utf8
	eucjp
	jis
	shiftjis
)

func (c enc) fmtinfo() fmtinfo.Encode {
	switch c {
	case utf8:
		return fmtinfo.UTF8
	case eucjp:
		return fmtinfo.EUCJP
	case jis:
		return fmtinfo.JIS
	case shiftjis:
		return fmtinfo.SHIFTJIS
	default:
		return fmtinfo.Bin
	}
}

func str2enc(s string) (enc, error) {
	switch strings.ToUpper(s) {
	case "UTF8", "UTF-8":
		return utf8, nil
	case "EUCJP", "EUC-JP", "EUC_JP":
		return eucjp, nil
	case "EUC":
		// NOTE: "EUC" might not be "JP" in future.
		return eucjp, nil
	case "JIS", "ISO2022JP":
		return jis, nil
	case "CP932", "SJIS", "Shift_JIS", "WIN31J":
		return shiftjis, nil
	case "U":
		return utf8, nil
	case "E":
		return eucjp, nil
	case "J":
		return jis, nil
	case "S":
		return shiftjis, nil
	case "":
		return noenc, nil
	default:
		return noenc, fmt.Errorf("unknown -enc: %s", s)
	}
}

type eol int

const (
	noeol eol = iota
	lf
	crlf
	cr
)

func (c eol) fmtinfo() fmtinfo.EOL {
	switch c {
	case lf:
		return fmtinfo.LF
	case crlf:
		return fmtinfo.CRLF
	case cr:
		return fmtinfo.CR
	default:
		return fmtinfo.Mix
	}
}

func str2eol(s string) (eol, error) {
	switch strings.ToUpper(s) {
	case "LF", "UNIX", "OSX":
		return lf, nil
	case "CRLF", "WIN", "DOS":
		return crlf, nil
	case "CR", "MAC":
		return cr, nil
	case "U":
		return lf, nil
	case "W", "D":
		return crlf, nil
	case "M":
		return lf, nil
	case "":
		return noeol, nil
	default:
		return noeol, fmt.Errorf("unknown -eol: %s", s)
	}
}

const excludeDefaults = `\.git$|\.svn$|\.hg$|\.o$|\.obj$|\.exe$`

var (
	help    bool
	optenc  string
	opteol  string
	exclude string
)

var (
	toEnc          enc
	toEOL          eol
	excludePattern *regexp.Regexp
)

func usage(err error) {
	c := 0
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n\n", err)
		c = 1
	}
	fmt.Fprintf(os.Stderr, "USAGE: textfmt [OPTIONS] {files or directories}\n\nOPTIONS:\n")
	flag.PrintDefaults()
	os.Exit(c)
}

func main() {
	// Parse flags.
	flag.BoolVar(&help, "h", false, "show this help")
	flag.StringVar(&optenc, "enc", "", "encoding (UTF8, EUC, JIS, CP932)")
	flag.StringVar(&opteol, "eol", "", "end of line (LF, CRLF, CR)")
	flag.StringVar(&exclude, "exclude", excludeDefaults, "excludes file/dir pattern")
	flag.Parse()
	if help {
		usage(nil)
	}
	// Check options and args.
	var err error
	if toEnc, err = str2enc(optenc); err != nil {
		usage(err)
	}
	if toEOL, err = str2eol(opteol); err != nil {
		usage(err)
	}
	if exclude != "" {
		if excludePattern, err = regexp.Compile(exclude); err != nil {
			usage(err)
		}
	}
	targets := flag.Args()
	if len(targets) == 0 {
		usage(errors.New("require at least one file or directory"))
	}
	for _, v := range targets {
		err := procPath(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", v, err)
		}
	}
}

func peelError(err error) error {
	if pe, ok := err.(*os.PathError); ok {
		err = pe.Err
	}
	return err
}

func procPath(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		return peelError(err)
	}
	if fi.IsDir() {
		return filepath.Walk(path, procWalk)
	}
	return procFile(path)
}

func procWalk(path string, info os.FileInfo, err error) error {
	if err != nil || info == nil {
		return err
	}
	if excludePattern != nil && excludePattern.MatchString(path) {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}
	if info.IsDir() {
		return nil
	}
	return procFile(path)
}

func procFile(path string) error {
	srcInfo, err := fmtinfo.Detect(path)
	if err != nil {
		return err
	}
	t := srcInfo.Transformer(dstInfo())
	if t == nil {
		fmt.Printf("%s (%s)\n", path, srcInfo)
		return nil
	}
	return convert.Convert(path, t)
}

func dstInfo() *fmtinfo.Info {
	return &fmtinfo.Info{
		Encode: toEnc.fmtinfo(),
		EOL:    toEOL.fmtinfo(),
	}
}
