package main

import (
	"bufio"
	"fmt"
	stack "m/Stack"
	"os"
	"regexp"
	"strings"
)

type listStack struct {
	Space string
	Ts    stack.Stack
}

func NewListStack() listStack {
	return listStack{
		Space: "",
		Ts:    stack.Stack{},
	}
}

func (ls *listStack) IsListType(line string) bool {
	sub, err := regexp.MatchString(`(^\s*\d+\.\s.+)|(^\s*-\s.+)`, line)
	if err != nil {
		fmt.Printf("Erro: %v\n", err)
	}
	return sub
}

func (ls *listStack) SetListType(line string) {
	if ls.IsListType(line) {
		typeCheck := regexp.MustCompile(`^\s*\d+\.\s.+`)
		if typeCheck.MatchString(line) {
			ls.Ts.Push("enumerate")
		} else {
			ls.Ts.Push("itemize")
		}
	}
}

func (ls *listStack) GetSpace(line string) string {
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

func main() {
	// Abre o arquivo Markdown para leitura
	inputFile, err := os.Open("input.md")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer inputFile.Close()

	// Cria um novo arquivo LaTeX para saída
	outputFile, err := os.Create("output.tex")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	mbOpen := false
	stack := NewListStack()
	for scanner.Scan() {
		line := scanner.Text()

		// Realiza substituições
		line = replacerMarkdownToLatex(line)

		// Montagem de blocos matemáticos
		line, mbOpen = convertMathBlocks(line, mbOpen)

		// Realiza montagem de listas de itens
		line = convertLists(line, &stack)

		if stack.Ts.Size() > 0 && len(line) < 2 {
			continue
		}
		fmt.Fprintln(writer, line)
	}

	for stack.Ts.Size() > 0 {
		stg := strings.Repeat("\t", stack.Ts.Size()-1) + "\\end{" + stack.Ts.Pop().(string) + "}\n"
		fmt.Fprintln(writer, stg)
	}

	if mbOpen {
		fmt.Fprintln(writer, "\\]")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	writer.Flush()
}

func replacerMarkdownToLatex(line string) string {
	// Regex para substituições
	replacements := []struct {
		re   *regexp.Regexp
		repl string
	}{
		{regexp.MustCompile(`^# (.+)`), `\section*{$1}`},
		{regexp.MustCompile(`^## (.+)`), `\subsection*{$1}`},
		{regexp.MustCompile(`^### (.+)`), `\subsubsection*{$1}`},
		{regexp.MustCompile(`^#### (.+)`), `\paragraph*{$1}`},
		{regexp.MustCompile(`\*\*(.+?)\*\*`), `\textbf{$1}`},
		{regexp.MustCompile(`\*(.+?)\*`), `\textit{$1}`},
		{regexp.MustCompile(`^---$`), `\hrulefill`},
	}

	for _, repl := range replacements {
		line = repl.re.ReplaceAllString(line, repl.repl)
	}

	return line
}

func replacerListMarkdownToLatex(line string, s *listStack) string {
	// Regex para substituições
	replacements := []struct {
		re   *regexp.Regexp
		repl string
	}{
		{regexp.MustCompile(`^\s*\d+\.\s(.+)`), strings.Repeat("\t", s.Ts.Size()) + `\item $1`},
		{regexp.MustCompile(`^\s*- (.+)`), strings.Repeat("\t", s.Ts.Size()) + `\item $1`},
	}

	for _, repl := range replacements {
		line = repl.re.ReplaceAllString(line, repl.repl)
	}

	return line
}

func convertLists(line string, s *listStack) string {
	if s.IsListType(line) {
		// Está em uma lista
		space := s.GetSpace(line)
		if len(space) > len(s.Space) || s.Ts.IsEmpty() {
			// Inicio de lista
			s.SetListType(line)
			line = replacerListMarkdownToLatex(line, s)
			newLine := strings.Repeat("\t", s.Ts.Size()-1) + "\\begin{" + s.Ts.Peek().(string) + "}\n"
			line = newLine + line
		} else if len(space) < len(s.Space) {
			// Final de sub-lista
			line = replacerListMarkdownToLatex(line, s)
			line += "\n" + strings.Repeat("\t", s.Ts.Size()-1) + "\\end{" + s.Ts.Pop().(string) + "}"
		} else {
			// Apenas mais um item da lista
			line = replacerListMarkdownToLatex(line, s)
		}
		s.Space = space
	} else if len(line) > 2 && !s.Ts.IsEmpty() {
		// Finaliza lista por detectar linha válida fora da lista
		line = replacerListMarkdownToLatex(line, s)
		line = strings.Repeat("\t", s.Ts.Size()-1) + "\\end{" + s.Ts.Pop().(string) + "}\n\n" + line
	}
	return line
}

func convertMathBlocks(line string, mathBlockOpen bool) (string, bool) {
	// Realiza replaces considerando abertura e fechamento de blocos
	if strings.Contains(line, "$$") {
		mathBlockOpen = !mathBlockOpen
	}
	if mathBlockOpen {
		line = strings.Replace(line, "$$", "\\[", 1)
	} else {
		line = strings.Replace(line, "$$", "\\]", 1)
	}
	return line, mathBlockOpen
}
