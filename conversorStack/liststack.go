package conversorstack

import (
	"fmt"
	"m/stack"
	"regexp"
)

type ListStack struct {
	Space string
	Ts    stack.Stack
}

func NewListStack() ListStack {
	return ListStack{
		Space: "",
		Ts:    stack.Stack{},
	}
}

func (ls *ListStack) IsListType(line string) bool {
	sub, err := regexp.MatchString(`(^\s*\d+\.\s.+)|(^\s*-\s.+)`, line)
	if err != nil {
		fmt.Printf("Erro: %v\n", err)
	}
	return sub
}

func (ls *ListStack) SetListType(line string) {
	if ls.IsListType(line) {
		typeCheck := regexp.MustCompile(`^\s*\d+\.\s.+`)
		if typeCheck.MatchString(line) {
			ls.Ts.Push("enumerate")
		} else {
			ls.Ts.Push("itemize")
		}
	}
}

func (ls *ListStack) GetSpace(line string) string {
	space := ""
	var spaceRegex *regexp.Regexp
	if ls.Ts.Peek() == "enumerate" {
		spaceRegex = regexp.MustCompile(`^(\s+)\d+\.\s.+`)
	} else {
		spaceRegex = regexp.MustCompile(`^(\s+)-\s.+`)
	}
	spaceMatch := spaceRegex.FindStringSubmatch(line)
	if len(spaceMatch) > 1 {
		space = spaceMatch[1]
	}
	return space
}
