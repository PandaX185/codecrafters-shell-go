package commands

import (
	"strconv"
	"strings"
	"unicode"
)

func normalizeDQuotes(runes []rune) []rune {
	result := make([]rune, 0, len(runes))
	i := 0
	for i < len(runes) {
		if runes[i] == '\\' && i+1 < len(runes) {
			next := runes[i+1]
			switch next {
			case '\\':
				result = append(result, '\\')
				i += 2
				continue
			case '"':
				result = append(result, '"')
				i += 2
				continue
			case '\n':
				i += 2
				continue
			case '\r':
				if i+2 < len(runes) && runes[i+2] == '\n' {
					i += 3
					continue
				}
			case '$', '`':
				result = append(result, next)
				i += 2
				continue
			}
		}
		result = append(result, runes[i])
		i++
	}
	return result
}

func indexNonEscaped(runes []rune, quote rune) int {
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == quote {
			return i
		}
		if r == '\\' && i+1 < len(runes) {
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

		case '>':
			if esc {
				token = append(token, r)
				esc = false
				break
			}

			if len(token) != 0 {
				argc = append(argc, string(token))
				token = token[:0]
			}

			if i+1 < len(lineRunes) && lineRunes[i+1] == '>' {
				argc = append(argc, ">>")
				i++
			} else {
				argc = append(argc, ">")
			}

		case '1':
			if i < len(lineRunes)-2 && lineRunes[i+1] == '>' && lineRunes[i+2] == '>' && !esc {
				if len(token) != 0 {
					argc = append(argc, string(token))
					token = token[:0]
				}
				argc = append(argc, "1>>")
				i++
			} else if i < len(lineRunes)-1 && lineRunes[i+1] == '>' && !esc {
				if len(token) != 0 {
					argc = append(argc, string(token))
					token = token[:0]
				}
				argc = append(argc, "1>")
				i++
			} else {
				if esc {
					esc = false
				}
				token = append(token, r)
			}

		case '2':
			if i < len(lineRunes)-2 && lineRunes[i+1] == '>' && lineRunes[i+2] == '>' && !esc {
				if len(token) != 0 {
					argc = append(argc, string(token))
					token = token[:0]
				}
				argc = append(argc, "2>>")
				i++
			} else if i < len(lineRunes)-1 && lineRunes[i+1] == '>' && !esc {
				if len(token) != 0 {
					argc = append(argc, string(token))
					token = token[:0]
				}
				argc = append(argc, "2>")
				i++
			} else {
				if esc {
					esc = false
				}
				token = append(token, r)
			}

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
