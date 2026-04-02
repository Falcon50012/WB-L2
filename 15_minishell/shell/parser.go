package shell

import "os"

type Command struct {
	Args     []string
	RedirIn  string
	RedirOut string
}

type Pipeline struct {
	Cmds []*Command
}

type ListItem struct {
	Op       string
	Pipeline *Pipeline
}

type CommandList struct {
	Items []*ListItem
}

type tokKind int

const (
	tokWord tokKind = iota
	tokPipe
	tokAnd
	tokOr
	tokRIn
	tokROut
)

type tok struct {
	kind tokKind
	val  string
}

func Parse(line string) *CommandList {
	tokens := lex(line)
	cl, _ := parseList(tokens)
	return cl
}

func lex(line string) []tok {
	var tokens []tok
	i := 0

	for i < len(line) {
		switch line[i] {
		case ' ', '\t':
			i++

		case '|':
			if i+1 < len(line) && line[i+1] == '|' {
				tokens = append(tokens, tok{kind: tokOr})
				i += 2
			} else {
				tokens = append(tokens, tok{kind: tokPipe})
				i++
			}

		case '&':
			if i+1 < len(line) && line[i+1] == '&' {
				tokens = append(tokens, tok{kind: tokAnd})
				i += 2
			} else {
				tokens = append(tokens, tok{kind: tokWord, val: "&"})
				i++
			}

		case '<':
			tokens = append(tokens, tok{kind: tokRIn})
			i++

		case '>':
			tokens = append(tokens, tok{kind: tokROut})
			i++

		case '"':
			i++
			var buf []byte
			for i < len(line) && line[i] != '"' {
				if line[i] == '\\' && i+1 < len(line) {
					i++
					buf = append(buf, line[i])
				} else {
					buf = append(buf, line[i])
				}
				i++
			}
			if i < len(line) && line[i] == '"' {
				i++
			}
			tokens = append(tokens, tok{kind: tokWord, val: expandEnvSegment(string(buf))})

		case '\'':
			i++
			var buf []byte
			for i < len(line) && line[i] != '\'' {
				buf = append(buf, line[i])
				i++
			}
			if i < len(line) && line[i] == '\'' {
				i++
			}
			tokens = append(tokens, tok{kind: tokWord, val: string(buf)})

		default:
			j := i
			for j < len(line) {
				c := line[j]
				if c == ' ' || c == '\t' ||
					c == '|' || c == '&' ||
					c == '<' || c == '>' ||
					c == '"' || c == '\'' {
					break
				}
				j++
			}

			word := expandEnvSegment(line[i:j])
			tokens = append(tokens, tok{kind: tokWord, val: word})
			i = j
		}
	}

	return tokens
}

func expandEnvSegment(s string) string {
	var b []byte
	i := 0

	for i < len(s) {
		if s[i] != '$' {
			b = append(b, s[i])
			i++
			continue
		}

		if i+1 >= len(s) {
			b = append(b, '$')
			i++
			continue
		}

		if s[i+1] == '{' {
			j := i + 2
			for j < len(s) && s[j] != '}' {
				j++
			}
			if j >= len(s) {
				b = append(b, '$')
				i++
				continue
			}
			name := s[i+2 : j]
			if name != "" {
				b = append(b, os.Getenv(name)...)
			}
			i = j + 1
			continue
		}

		if !isVarStart(s[i+1]) {
			b = append(b, '$')
			i++
			continue
		}

		j := i + 1
		for j < len(s) && isVarChar(s[j]) {
			j++
		}
		name := s[i+1 : j]
		b = append(b, os.Getenv(name)...)
		i = j
	}

	return string(b)
}

func isVarStart(b byte) bool {
	return b == '_' ||
		(b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z')
}

func isVarChar(b byte) bool {
	return isVarStart(b) || (b >= '0' && b <= '9')
}

func parseList(tokens []tok) (*CommandList, []tok) {
	cl := &CommandList{}

	pl, rest := parsePipeline(tokens)
	if pl == nil {
		return cl, rest
	}

	cl.Items = append(cl.Items, &ListItem{Op: "", Pipeline: pl})

	for len(rest) > 0 {
		op := ""
		switch rest[0].kind {
		case tokAnd:
			op = "&&"
		case tokOr:
			op = "||"
		default:
			return cl, rest
		}

		rest = rest[1:]
		pl, rest = parsePipeline(rest)
		if pl == nil {
			return cl, rest
		}
		cl.Items = append(cl.Items, &ListItem{Op: op, Pipeline: pl})
	}

	return cl, rest
}

func parsePipeline(tokens []tok) (*Pipeline, []tok) {
	cmd, rest := parseCommand(tokens)
	if cmd == nil {
		return nil, tokens
	}

	pl := &Pipeline{Cmds: []*Command{cmd}}

	for len(rest) > 0 && rest[0].kind == tokPipe {
		rest = rest[1:]
		cmd, rest = parseCommand(rest)
		if cmd == nil {
			break
		}
		pl.Cmds = append(pl.Cmds, cmd)
	}

	return pl, rest
}

func parseCommand(tokens []tok) (*Command, []tok) {
	cmd := &Command{}
	i := 0

	for i < len(tokens) {
		switch tokens[i].kind {
		case tokWord:
			cmd.Args = append(cmd.Args, tokens[i].val)
			i++

		case tokRIn:
			if i+1 < len(tokens) && tokens[i+1].kind == tokWord {
				cmd.RedirIn = tokens[i+1].val
				i += 2
			} else {
				i++
			}

		case tokROut:
			if i+1 < len(tokens) && tokens[i+1].kind == tokWord {
				cmd.RedirOut = tokens[i+1].val
				i += 2
			} else {
				i++
			}

		default:
			return cmdIfValid(cmd), tokens[i:]
		}
	}

	return cmdIfValid(cmd), nil
}

func cmdIfValid(cmd *Command) *Command {
	if len(cmd.Args) == 0 {
		return nil
	}
	return cmd
}
