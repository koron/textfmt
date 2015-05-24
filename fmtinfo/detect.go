package fmtinfo

import (
	"io"
	"os"
)

const (
	defaultBuffer = 32 * 1024
)

// Detect detects format of text file: encoding and end of line.
func Detect(path string) (*Info, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return detect(f)
}

func detect(r io.Reader) (*Info, error) {
	// Setup detectors.
	enc := newEncodingDetector()
	eol := newEolDetector()
	d := newMultiDetector(enc, eol)
	// Do detection.
	buf := make([]byte, defaultBuffer)
	for {
		nr, err := r.Read(buf)
		if nr > 0 {
			d.parse(buf[0:nr], err == io.EOF)
			if !d.isParsing() {
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	// TODO: get result of encoding detection.
	return &Info{Encode: enc.encode, EOL: eol.eol}, nil
}

type detector interface {
	isParsing() bool
	parse(data []byte, atEOF bool)
}

type multiDetector struct {
	d []detector
}

func newMultiDetector(d ...detector) detector {
	return &multiDetector{
		d: d,
	}
}

func (m *multiDetector) isParsing() bool {
	for _, d := range m.d {
		if d.isParsing() {
			return true
		}
	}
	return false
}

func (m *multiDetector) parse(data []byte, atEOF bool) {
	for _, d := range m.d {
		if d.isParsing() {
			d.parse(data, atEOF)
		}
	}
}

type eolDetector struct {
	parsing bool
	eol     EOL
	prevCR  bool
}

func newEolDetector() *eolDetector {
	return &eolDetector{
		parsing: true,
		eol:     Mix,
		prevCR:  false,
	}
}

func (d *eolDetector) isParsing() bool {
	return d.parsing
}

func (d *eolDetector) parse(data []byte, atEOF bool) {
	for len(data) > 0 && d.parsing {
		nSkip, eol := findEOL(data)
		switch eol {
		case '\n':
			if d.prevCR {
				if nSkip == 0 {
					d.emit(CRLF)
				} else {
					d.emit(CR)
					d.emit(LF)
				}
			} else {
				d.emit(LF)
			}
			nSkip++
			d.prevCR = false
		case '\r':
			if d.prevCR {
				d.emit(CR)
			}
			nSkip++
			d.prevCR = true
		default:
			d.prevCR = false
		}
		data = data[nSkip:]
	}
	if atEOF && d.prevCR {
		d.emit(CR)
	}
}

func (d *eolDetector) emit(eol EOL) {
	if !d.parsing {
		return
	}
	if d.eol == Mix {
		d.eol = eol
		return
	}
	if d.eol != eol {
		d.eol = Mix
		d.parsing = false
	}
}

type encodingDetector struct {
	parsing bool
	// TODO:
	encoding Encode
}

func newEncodingDetector() *encodingDetector {
	return &encodingDetector{
		parsing: true,
		// TODO:
		encoding: Bin,
	}
}

func (d *encodingDetector) isParsing() bool {
	// TODO:
	return d.parsing
}

func (d *encodingDetector) parse(data []byte, atEOF bool) {
	// TODO:
}
