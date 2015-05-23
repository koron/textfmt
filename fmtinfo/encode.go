package fmtinfo

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// Encode shows charset encoding.
type Encode int

const (
	// Bin means binary file.
	Bin Encode = iota
	// UTF8 means UTF-8 encoding.
	UTF8
	// EUCJP means EUC-JP encoding.
	EUCJP
	// JIS means ISO-2022-JP encoding.
	JIS
	// SHIFTJIS means Shift_JIS, CP932 encoding.
	SHIFTJIS
)

// String returns Encode's string representation.
func (c Encode) String() string {
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

func (c Encode) encoding() encoding.Encoding {
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

func (c Encode) newDecoder() transform.Transformer {
	e := c.encoding()
	if e == nil {
		return nil
	}
	return e.NewDecoder()
}

func (c Encode) newEncoder() transform.Transformer {
	e := c.encoding()
	if e == nil {
		return nil
	}
	return e.NewEncoder()
}
