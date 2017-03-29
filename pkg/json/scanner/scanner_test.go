package scanner

import (
	"testing"

	"github.com/hashicorp/hcl/json/token"
)

func TestRealExample(t *testing.T) {
	simpleReal := `
  {
    "aa": 11,
    "bb": "測試"
  }
  `
	literals := []struct {
		tokenType token.Type
		literal   string
	}{
		{token.LBRACE, `{`},
		{token.STRING, `"aa"`},
		{token.COLON, `:`},
		{token.NUMBER, `11`},
		{token.COMMA, `,`},
		{token.STRING, `"bb"`},
		{token.COLON, `:`},
		{token.STRING, `"測試"`},
		{token.RBRACE, `}`},
		{token.EOF, ``},
	}

	s := New([]byte(simpleReal))
	// for i := 0; i < 20; i++ {
	// 	tok := s.Scan()
	// 	fmt.Println(tok)
	// }
	for _, l := range literals {
		tok := s.Scan()
		if l.literal != tok.Text {
			t.Errorf("actual: %s, expected: %s", tok.Text, l.literal)
		}
	}
}

// 	complexReal := `
// {
//     "variable": {
//         "foo": {
//             "default": "bar",
//             "description": "bar",
//             "depends_on": ["something"]
//         }
//     }
// }`
