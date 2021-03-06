package fmtinfo

import (
	"fmt"

	"golang.org/x/text/transform"
)

// Info is text format information.
type Info struct {
	Encoding Encoding
	EOL      EOL
}

// String returns Info's string representation.
func (n *Info) String() string {
	if n.Encoding == Bin {
		return "binary file"
	}
	return fmt.Sprintf("%s, %s", n.Encoding.String(), n.EOL.String())
}

// Transformer build and return a transformer of text format.
func (n *Info) Transformer(to *Info) transform.Transformer {
	if n == nil || to == nil {
		return nil
	}
	var d, m, e transform.Transformer
	if to.Encoding != Bin {
		if n.Encoding != Bin && n.Encoding != to.Encoding {
			d = n.Encoding.newDecoder()
			e = to.Encoding.newEncoder()
		}
	}
	if to.EOL != Mix {
		if n.EOL != Mix && n.EOL != to.EOL {
			m = to.EOL.Transformer()
		}
	}
	t := make([]transform.Transformer, 0, 3)
	if d != nil {
		t = append(t, d)
	}
	if m != nil {
		t = append(t, m)
	}
	if e != nil {
		t = append(t, e)
	}
	if len(t) == 0 {
		return nil
	}
	return transform.Chain(t...)
}
