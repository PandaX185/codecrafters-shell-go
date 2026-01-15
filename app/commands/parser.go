package commands

import (
	"strings"
	"unicode"
)

func normalizeDQuotes(runes []rune) []rune {
	s := string(runes)
	s = unescapeString(s)
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

func tokenize(line string) (argc []string) {
	line = strings.TrimLeftFunc(line, unicode.IsSpace)
	if len(line) == 0 {
		return nil
	}

	token := make([]rune, 0, 256)
	lineRunes := []rune(line)
	escMode := false

	i := 0
	for i < len(lineRunes) {
		r := lineRunes[i]

		switch r {
		case ' ', '\t', '\n':
			if escMode {
				token = append(token, r)
				escMode = false
				break
			}

			if len(token) == 0 {
				break
			}

			argc = append(argc, string(token))
			token = token[:0]

		case '\'', '"':
			if escMode {
				token = append(token, r)
				escMode = false
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
			}

			token = append(token, toAppend...)
			i += nextQuoteInd

		case '\\':
			if escMode {
				token = append(token, r)
				escMode = false
				break
			}

			escMode = true

		default:
			if escMode {
				escMode = false
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
	return tokenize(line)
}
