package jsonstream

import "io"

type Reader struct {
	next      io.Reader
	depth     int
	prevC     byte
	inString  bool
	replacing bool
}

func New(r io.Reader) *Reader {
	return &Reader{next: r}
}

func (r *Reader) Read(bs []byte) (int, error) {
	n, err := r.next.Read(bs)
	for i, c := range bs[:n] {
		switch {
		case !r.inString && c == '\n':
			bs[i] = ' '
		case !r.inString && c == '"':
			r.inString = true
		case r.inString && c == '"' && r.prevC != '\\':
			r.inString = false
		case !r.inString && c == '{':
			r.depth++
		case !r.inString && c == '[':
			// first should be removed
			if r.depth == 0 {
				bs[i] = ' '
				continue
			}
			r.depth++
		case !r.inString && c == '}' || c == ']':
			if r.depth == 0 {
				bs[i] = ' '
				continue
			}
			r.depth--
			if r.depth == 0 {
				// We are now at the top level. The next comma should become a newline.
				r.replacing = true
			}
		case !r.inString && r.replacing && c == ',':
			bs[i] = '\n'
			r.replacing = false
		}
		r.prevC = c
	}
	return n, err
}
