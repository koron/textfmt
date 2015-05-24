package fmtinfo

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Encoding shows charset encoding.
type Encoding int

const (
	// Bin means binary file.
	Bin Encoding = iota
	// UTF8 means UTF-8 encoding.
	UTF8
	// EUCJP means EUC-JP encoding.
	EUCJP
	// JIS means ISO-2022-JP encoding.
	JIS
	// SHIFTJIS means Shift_JIS, CP932 encoding.
	SHIFTJIS
)

// String returns Encoding's string representation.
func (c Encoding) String() string {
	switch c {
	case Bin:
		return "binary"
	case UTF8:
		return "UTF-8"
	case EUCJP:
		return "EUC-JP"
	case JIS:
		return "ISO-2022-JP"
	case SHIFTJIS:
		return "Shift_JIS"
	default:
		return ""
	}
}

func (c Encoding) encoding() encoding.Encoding {
	switch c {
	case EUCJP:
		return japanese.EUCJP
	case JIS:
		return japanese.ISO2022JP
	case SHIFTJIS:
		return japanese.ShiftJIS
	default:
		return nil
	}
}

func (c Encoding) newDecoder() transform.Transformer {
	e := c.encoding()
	if e == nil {
		return nil
	}
	return e.NewDecoder()
}

func (c Encoding) newEncoder() transform.Transformer {
	e := c.encoding()
	if e == nil {
		return nil
	}
	return e.NewEncoder()
}
