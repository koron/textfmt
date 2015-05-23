package fmtinfo

import "golang.org/x/text/transform"

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

// Transformer creates new transform.Transformer to covert to the EOL.
func (c EOL) Transformer() transform.Transformer {
	// TODO:
	return nil
}
