package fmtinfo

import "golang.org/x/text/transform"

func findEOL(bytes []byte) (int, byte) {
	for i, b := range bytes {
		switch b {
		case '\r', '\n':
			return i, b
		}
	}
	return len(bytes), 0
}

// EOL is type of end of line codes.
type EOL int

const (
	// Mix is mixed EOL.
	Mix EOL = iota
	// LF shows lines are ended with LF.
	LF
	// CRLF shows lines are ended with CRLF.
	CRLF
	// CR shows lines are ended with CR.
	CR
)

// String returns EOL's string representation.
func (c EOL) String() string {
	switch c {
	case Mix:
		return "mixed"
	case LF:
		return "LF"
	case CRLF:
		return "CR+LF"
	case CR:
		return "CR"
	default:
		return ""
	}
}

func (c EOL) bytes() []byte {
	switch c {
	case LF:
		return []byte{'\n'}
	case CRLF:
		return []byte{'\r', '\n'}
	case CR:
		return []byte{'\r'}
	default:
		return []byte{}
	}
}

// Transformer creates new transform.Transformer to covert to the EOL.
func (c EOL) Transformer() transform.Transformer {
	return &eolTransformer{replacement: c.bytes()}
}

type eolTransformer struct {
	replacement []byte
	prevCR      bool
}

func (t *eolTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	w := &buffer{bytes: dst}
	lenDst, lenSrc := len(dst), len(src)

	emitEOL := func() error {
		if w.Capacity() < len(t.replacement) {
			return transform.ErrShortDst
		}
		m, _ := w.Write(t.replacement)
		nDst += m
		return nil
	}

	for nDst < lenDst && nSrc < lenSrc {
		nSkip, eol := findEOL(src)
		if nSkip > 0 {
			m, _ := w.Write(src[0:nSkip])
			nDst += m
			nSrc += m
			src = src[m:]
			if err != nil {
				return nDst, nSrc, err
			}
		}
		switch eol {
		case '\n':
			if !(nSkip == 0 && t.prevCR) {
				if err := emitEOL(); err != nil {
					return nDst, nSrc, err
				}
			}
			nSrc++
			src = src[1:]
			t.prevCR = false
		case '\r':
			if err := emitEOL(); err != nil {
				return nDst, nSrc, err
			}
			nSrc++
			src = src[1:]
			t.prevCR = true
		default:
			t.prevCR = false
		}
	}
	return nDst, nSrc, nil
}

func (t *eolTransformer) Reset() {
	t.prevCR = false
}

type buffer struct {
	bytes []byte
	index int
}

func (b *buffer) Write(src []byte) (n int, err error) {
	l := len(src)
	if r := b.Capacity(); r < l {
		l = r
		err = transform.ErrShortDst
	}
	n = copy(b.bytes[b.index:], src[0:l])
	b.index += n
	return n, err
}

func (b *buffer) Capacity() int {
	return len(b.bytes) - b.index
}
