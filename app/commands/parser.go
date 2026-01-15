package commands

import (
	"strconv"
	"strings"
	"unicode"
)

func normalizeDQuotes(runes []rune) []rune {
	s := string(runes)
	s = strings.ReplaceAll(s, `\\`, `\`)
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, "\\\n", "")
	s = strings.ReplaceAll(s, "\\\r\n", "")
	return []rune(s)
}

func indexNonEscaped(runes []rune, quote rune) int {
	for i, r := range runes {
		if r == quote {
			return i
		}
		if r == '\\' && i+1 < len(runes) && runes[i+1] == quote {
			i++
		}
	}
	return -1
}

func TokenizeCommand(line string) (argc []string) {
	line = strings.TrimLeftFunc(line, unicode.IsSpace)
	if len(line) == 0 {
		return nil
	}

	token := make([]rune, 0, 256)
	lineRunes := []rune(line)
	esc := false

	i := 0
	for i < len(lineRunes) {
		r := lineRunes[i]

		switch r {
		case ' ', '\t', '\n':
			if esc {
				token = append(token, r)
				esc = false
				break
			}

			if len(token) == 0 {
				break
			}

			argc = append(argc, string(token))
			token = token[:0]

		case '\'', '"':
			if esc {
				token = append(token, r)
				esc = false
				break
			}

			i++
			nextQuoteInd := 0
			if r == '\'' {
				for j, rr := range lineRunes[i:] {
					if rr == r {
						nextQuoteInd = j
						break
					}
				}
				if nextQuoteInd == 0 && lineRunes[i+nextQuoteInd] != r {
					nextQuoteInd = -1
				}
			} else {
				nextQuoteInd = indexNonEscaped(lineRunes[i:], r)
			}

			if nextQuoteInd == -1 {
				toAppend := lineRunes[i:]
				token = append(token, toAppend...)
				i = len(lineRunes)
				continue
			}

			toAppend := lineRunes[i : i+nextQuoteInd]
			if r == '"' {
				toAppend = normalizeDQuotes(toAppend)
			} else if r == '\'' {
				s := string(toAppend)
				s = UnescapeString(s)
				toAppend = []rune(s)
			}

			token = append(token, toAppend...)
			i += nextQuoteInd

		case '\\':
			if esc {
				token = append(token, r)
				esc = false
				break
			}

			esc = true

		default:
			if esc {
				esc = false
			}

			token = append(token, r)
		}

		i++
	}

	if len(token) != 0 {
		argc = append(argc, string(token))
	}
	return
}

func Parse(line string) []string {
	line = strings.TrimSuffix(line, "\n")
	return TokenizeCommand(line)
}

func UnescapeString(s string) string {
	unescaped, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
		return s
	}
	return unescaped
}
