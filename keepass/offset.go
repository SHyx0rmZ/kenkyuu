package keepass

import (
	"encoding/base64"
)

func FindOffset(g *Group, e *Entry) int {
	n, ok := findOffset(g, e)
	if !ok {
		return 0
	}
	return n
}

func findOffset(g *Group, s *Entry) (n int, ok bool) {
	for _, g := range g.Group {
		gn, gok := findOffset(g, s)
		if gok {
			return n + gn, true
		}
	}

	for _, e := range g.Entry {
		if e == s {
			return n, true
		}

		for _, f := range e.String {
			if !f.Value.Protected || len(f.Value.Value) == 0 {
				continue
			}

			bs, err := base64.StdEncoding.DecodeString(f.Value.Value)
			if err != nil {
				panic("corrupted database")
			}

			n += len(bs)
		}

		for _, h := range e.History {
			for _, f := range h.String {
				if !f.Value.Protected || len(f.Value.Value) == 0 {
					continue
				}

				bs, err := base64.StdEncoding.DecodeString(f.Value.Value)
				if err != nil {
					panic("corrupted database")
				}

				n += len(bs)
			}
		}
	}

	return n, false
}
