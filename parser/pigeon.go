package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var g = &grammar{
	rules: []*rule{
		{
			name: "File",
			pos:  position{line: 4, col: 1, offset: 54},
			expr: &actionExpr{
				pos: position{line: 4, col: 9, offset: 62},
				run: (*parser).callonFile1,
				expr: &seqExpr{
					pos: position{line: 4, col: 9, offset: 62},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 4, col: 9, offset: 62},
							label: "lines",
							expr: &zeroOrMoreExpr{
								pos: position{line: 4, col: 15, offset: 68},
								expr: &ruleRefExpr{
									pos:  position{line: 4, col: 15, offset: 68},
									name: "Line",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 4, col: 21, offset: 74},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Line",
			pos:  position{line: 20, col: 1, offset: 407},
			expr: &actionExpr{
				pos: position{line: 20, col: 9, offset: 415},
				run: (*parser).callonLine1,
				expr: &seqExpr{
					pos: position{line: 20, col: 9, offset: 415},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 20, col: 9, offset: 415},
							label: "ws",
							expr: &ruleRefExpr{
								pos:  position{line: 20, col: 12, offset: 418},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 20, col: 14, offset: 420},
							label: "item",
							expr: &zeroOrOneExpr{
								pos: position{line: 20, col: 19, offset: 425},
								expr: &choiceExpr{
									pos: position{line: 20, col: 20, offset: 426},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 20, col: 20, offset: 426},
											name: "Comment",
										},
										&ruleRefExpr{
											pos:  position{line: 20, col: 30, offset: 436},
											name: "Section",
										},
										&ruleRefExpr{
											pos:  position{line: 20, col: 40, offset: 446},
											name: "KeyValuePair",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 20, col: 55, offset: 461},
							label: "le",
							expr: &ruleRefExpr{
								pos:  position{line: 20, col: 58, offset: 464},
								name: "LineEnd",
							},
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 29, col: 1, offset: 682},
			expr: &actionExpr{
				pos: position{line: 29, col: 12, offset: 693},
				run: (*parser).callonComment1,
				expr: &seqExpr{
					pos: position{line: 29, col: 12, offset: 693},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 29, col: 12, offset: 693},
							label: "cs",
							expr: &choiceExpr{
								pos: position{line: 29, col: 16, offset: 697},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 29, col: 16, offset: 697},
										val:        ";",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 29, col: 22, offset: 703},
										val:        "#",
										ignoreCase: false,
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 29, col: 27, offset: 708},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 29, col: 35, offset: 716},
								name: "CommentVal",
							},
						},
					},
				},
			},
		},
		{
			name: "Section",
			pos:  position{line: 37, col: 1, offset: 929},
			expr: &actionExpr{
				pos: position{line: 37, col: 12, offset: 940},
				run: (*parser).callonSection1,
				expr: &seqExpr{
					pos: position{line: 37, col: 12, offset: 940},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 37, col: 12, offset: 940},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 37, col: 16, offset: 944},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 37, col: 21, offset: 949},
								name: "SectionName",
							},
						},
						&litMatcher{
							pos:        position{line: 37, col: 33, offset: 961},
							val:        "]",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 37, col: 37, offset: 965},
							label: "ws",
							expr: &ruleRefExpr{
								pos:  position{line: 37, col: 40, offset: 968},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 37, col: 42, offset: 970},
							label: "comment",
							expr: &zeroOrOneExpr{
								pos: position{line: 37, col: 50, offset: 978},
								expr: &ruleRefExpr{
									pos:  position{line: 37, col: 50, offset: 978},
									name: "Comment",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "KeyValuePair",
			pos:  position{line: 46, col: 1, offset: 1206},
			expr: &actionExpr{
				pos: position{line: 46, col: 17, offset: 1222},
				run: (*parser).callonKeyValuePair1,
				expr: &seqExpr{
					pos: position{line: 46, col: 17, offset: 1222},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 46, col: 17, offset: 1222},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 46, col: 21, offset: 1226},
								name: "Key",
							},
						},
						&litMatcher{
							pos:        position{line: 46, col: 25, offset: 1230},
							val:        "=",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 46, col: 29, offset: 1234},
							label: "ws",
							expr: &ruleRefExpr{
								pos:  position{line: 46, col: 32, offset: 1237},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 46, col: 34, offset: 1239},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 46, col: 38, offset: 1243},
								name: "Value",
							},
						},
						&labeledExpr{
							pos:   position{line: 46, col: 44, offset: 1249},
							label: "comment",
							expr: &zeroOrOneExpr{
								pos: position{line: 46, col: 52, offset: 1257},
								expr: &ruleRefExpr{
									pos:  position{line: 46, col: 52, offset: 1257},
									name: "Comment",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CommentVal",
			pos:  position{line: 55, col: 1, offset: 1518},
			expr: &actionExpr{
				pos: position{line: 55, col: 15, offset: 1532},
				run: (*parser).callonCommentVal1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 55, col: 15, offset: 1532},
					expr: &seqExpr{
						pos: position{line: 55, col: 16, offset: 1533},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 55, col: 16, offset: 1533},
								expr: &ruleRefExpr{
									pos:  position{line: 55, col: 17, offset: 1534},
									name: "LineEnd",
								},
							},
							&anyMatcher{
								line: 55, col: 25, offset: 1542,
							},
						},
					},
				},
			},
		},
		{
			name: "SectionName",
			pos:  position{line: 63, col: 1, offset: 1709},
			expr: &actionExpr{
				pos: position{line: 63, col: 16, offset: 1724},
				run: (*parser).callonSectionName1,
				expr: &oneOrMoreExpr{
					pos: position{line: 63, col: 16, offset: 1724},
					expr: &charClassMatcher{
						pos:        position{line: 63, col: 16, offset: 1724},
						val:        "[^#;\\r\\n[\\]]",
						chars:      []rune{'#', ';', '\r', '\n', '[', ']'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 71, col: 1, offset: 1902},
			expr: &actionExpr{
				pos: position{line: 71, col: 8, offset: 1909},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 71, col: 8, offset: 1909},
					expr: &charClassMatcher{
						pos:        position{line: 71, col: 8, offset: 1909},
						val:        "[^#;=\\r\\n[\\]]",
						chars:      []rune{'#', ';', '=', '\r', '\n', '[', ']'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 79, col: 1, offset: 2080},
			expr: &choiceExpr{
				pos: position{line: 79, col: 10, offset: 2089},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 79, col: 10, offset: 2089},
						name: "QuotedValue",
					},
					&actionExpr{
						pos: position{line: 79, col: 24, offset: 2103},
						run: (*parser).callonValue3,
						expr: &ruleRefExpr{
							pos:  position{line: 79, col: 24, offset: 2103},
							name: "SimpleValue",
						},
					},
				},
			},
		},
		{
			name: "QuotedValue",
			pos:  position{line: 87, col: 1, offset: 2273},
			expr: &actionExpr{
				pos: position{line: 87, col: 16, offset: 2288},
				run: (*parser).callonQuotedValue1,
				expr: &seqExpr{
					pos: position{line: 87, col: 16, offset: 2288},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 87, col: 16, offset: 2288},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 87, col: 20, offset: 2292},
							expr: &ruleRefExpr{
								pos:  position{line: 87, col: 20, offset: 2292},
								name: "Char",
							},
						},
						&litMatcher{
							pos:        position{line: 87, col: 26, offset: 2298},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Char",
			pos:  position{line: 95, col: 1, offset: 2466},
			expr: &choiceExpr{
				pos: position{line: 95, col: 9, offset: 2474},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 95, col: 9, offset: 2474},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 95, col: 9, offset: 2474},
								expr: &choiceExpr{
									pos: position{line: 95, col: 11, offset: 2476},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 95, col: 11, offset: 2476},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 95, col: 17, offset: 2482},
											val:        "\\",
											ignoreCase: false,
										},
									},
								},
							},
							&anyMatcher{
								line: 95, col: 23, offset: 2488,
							},
						},
					},
					&actionExpr{
						pos: position{line: 95, col: 27, offset: 2492},
						run: (*parser).callonChar8,
						expr: &seqExpr{
							pos: position{line: 95, col: 27, offset: 2492},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 95, col: 27, offset: 2492},
									val:        "\\",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 95, col: 33, offset: 2498},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 95, col: 33, offset: 2498},
											val:        "[\\\\/bfnrt\"]",
											chars:      []rune{'\\', '/', 'b', 'f', 'n', 'r', 't', '"'},
											ignoreCase: false,
											inverted:   false,
										},
										&seqExpr{
											pos: position{line: 95, col: 47, offset: 2512},
											exprs: []interface{}{
												&litMatcher{
													pos:        position{line: 95, col: 47, offset: 2512},
													val:        "u",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 95, col: 51, offset: 2516},
													name: "HexDigit",
												},
												&ruleRefExpr{
													pos:  position{line: 95, col: 60, offset: 2525},
													name: "HexDigit",
												},
												&ruleRefExpr{
													pos:  position{line: 95, col: 69, offset: 2534},
													name: "HexDigit",
												},
												&ruleRefExpr{
													pos:  position{line: 95, col: 78, offset: 2543},
													name: "HexDigit",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 103, col: 1, offset: 2725},
			expr: &actionExpr{
				pos: position{line: 103, col: 13, offset: 2737},
				run: (*parser).callonHexDigit1,
				expr: &charClassMatcher{
					pos:        position{line: 103, col: 13, offset: 2737},
					val:        "[0-9a-f]i",
					ranges:     []rune{'0', '9', 'a', 'f'},
					ignoreCase: true,
					inverted:   false,
				},
			},
		},
		{
			name: "SimpleValue",
			pos:  position{line: 111, col: 1, offset: 2908},
			expr: &actionExpr{
				pos: position{line: 111, col: 16, offset: 2923},
				run: (*parser).callonSimpleValue1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 111, col: 16, offset: 2923},
					expr: &charClassMatcher{
						pos:        position{line: 111, col: 16, offset: 2923},
						val:        "[^;#\\r\\n]",
						chars:      []rune{';', '#', '\r', '\n'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "LineEnd",
			pos:  position{line: 119, col: 1, offset: 3098},
			expr: &choiceExpr{
				pos: position{line: 119, col: 12, offset: 3109},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 119, col: 12, offset: 3109},
						val:        "\r\n",
						ignoreCase: false,
					},
					&actionExpr{
						pos: position{line: 119, col: 21, offset: 3118},
						run: (*parser).callonLineEnd3,
						expr: &litMatcher{
							pos:        position{line: 119, col: 21, offset: 3118},
							val:        "\n",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 127, col: 1, offset: 3259},
			expr: &actionExpr{
				pos: position{line: 127, col: 19, offset: 3277},
				run: (*parser).callon_1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 127, col: 19, offset: 3277},
					expr: &charClassMatcher{
						pos:        position{line: 127, col: 19, offset: 3277},
						val:        "[ \\t]",
						chars:      []rune{' ', '\t'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 135, col: 1, offset: 3413},
			expr: &notExpr{
				pos: position{line: 135, col: 8, offset: 3420},
				expr: &anyMatcher{
					line: 135, col: 9, offset: 3421,
				},
			},
		},
	},
}

func (c *current) onFile1(lines interface{}) (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf("\n\n\n>> File: %s // '%s'", c.pos, string(c.text))

	// convert iface to []*Line
	lsSlice := toIfaceSlice(lines)
	ls := make([]*Line, len(lsSlice))
	for i, l := range lsSlice {
		ls[i] = l.(*Line)
	}

	return NewFile(ls), nil
}

func (p *parser) callonFile1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFile1(stack["lines"])
}

func (c *current) onLine1(ws, item, le interface{}) (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Line: %s // '%s'", c.pos, string(c.text))
	it, _ := item.(Item)
	return NewLine(c.pos, ws.(string), it, le.(string)), nil
}

func (p *parser) callonLine1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLine1(stack["ws"], stack["item"], stack["le"])
}

func (c *current) onComment1(cs, comment interface{}) (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Comment: %s // '%s'\n", c.pos, string(c.text))
	return NewComment(c.pos, string(cs.([]byte)), comment.(string)), nil
}

func (p *parser) callonComment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onComment1(stack["cs"], stack["comment"])
}

func (c *current) onSection1(name, ws, comment interface{}) (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Section: %s // '%s'\n", c.pos, name)
	com, _ := comment.(*Comment)
	return NewSection(c.pos, name.(string), ws.(string), com), nil
}

func (p *parser) callonSection1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSection1(stack["name"], stack["ws"], stack["comment"])
}

func (c *current) onKeyValuePair1(key, ws, val, comment interface{}) (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> KeyValuePair: %s // '%s': '%s'\n", c.pos, key, val)
	com, _ := comment.(*Comment)
	return NewKeyValuePair(c.pos, key.(string), ws.(string), val.(string), com), nil
}

func (p *parser) callonKeyValuePair1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onKeyValuePair1(stack["key"], stack["ws"], stack["val"], stack["comment"])
}

func (c *current) onCommentVal1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> CommentVal: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonCommentVal1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCommentVal1()
}

func (c *current) onSectionName1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> SectionName: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonSectionName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSectionName1()
}

func (c *current) onKey1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Key: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonKey1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onKey1()
}

func (c *current) onValue3() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Value: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonValue3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onValue3()
}

func (c *current) onQuotedValue1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> QuotedValue: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonQuotedValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuotedValue1()
}

func (c *current) onChar8() (interface{}, error) {
	// " // ignore
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Char: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonChar8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onChar8()
}

func (c *current) onHexDigit1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> HexDigit: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonHexDigit1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHexDigit1()
}

func (c *current) onSimpleValue1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> SimpleValue: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonSimpleValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSimpleValue1()
}

func (c *current) onLineEnd3() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> LineEnd: %s\n", c.pos)
	return string(c.text), nil
}

func (p *parser) callonLineEnd3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLineEnd3()
}

func (c *current) on_1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> _ %s\n", c.pos)
	return string(c.text), nil
}

func (p *parser) callon_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.on_1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}