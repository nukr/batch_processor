package scanner

import (
	"bytes"
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/nukr/batch_processor/pkg/json/token"
)

const eof = rune(0)

// Scanner ...
type Scanner struct {
	src []byte
	buf *bytes.Buffer

	srcPos      token.Pos
	prevPos     token.Pos
	tokPos      token.Pos
	Error       func(pos token.Pos, msg string)
	ErrorCount  int
	lastCharLen int
}

// New ...
func New(src []byte) *Scanner {
	b := bytes.NewBuffer(src)
	s := &Scanner{
		src: src,
		buf: b,
	}
	s.srcPos.Line = 1
	return s
}

// Scan scans the next token and returns the token
func (s *Scanner) Scan() token.Token {
	ch := s.next()
	for unicode.IsSpace(ch) {
		ch = s.next()
	}

	fmt.Println("all ch", string(ch))

	var tok token.Type

	switch {
	case unicode.IsLetter(ch):
		fmt.Println(string(ch))
	case unicode.IsDigit(ch):
		fmt.Println("digit!!", string(ch))
	default:
		switch ch {
		case '{':
			tok = token.LBRACE
		}
	}
	return token.Token{
		Type: tok,
	}
}

func (s *Scanner) next() rune {
	ch, size, err := s.buf.ReadRune()
	if err != nil {
		s.srcPos.Column++
		s.srcPos.Offset += size
		s.lastCharLen = size
		return eof
	}
	if ch == utf8.RuneError && size == 1 {
		s.srcPos.Column++
		s.srcPos.Offset += size
		s.lastCharLen = size
		s.err("illegal UTF-8 encoding")
		return ch
	}
	fmt.Printf("ch: %q, size: %d\n", ch, size)
	return ch
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (s *Scanner) err(msg string) {
	s.ErrorCount++
	pos := s.recentPosition()

	if s.Error != nil {
		s.Error(pos, msg)
		return
	}
	fmt.Fprintf(os.Stderr, "%s: %s\n", pos, msg)
}

func (s *Scanner) recentPosition() (pos token.Pos) {
}
