package main

import (
	"bufio"
	"fmt"
	convSta "m/conversorStack"
	"os"
	"regexp"
	"strings"
)

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
	stack := convSta.NewListStack()
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

func replacerListMarkdownToLatex(line string, tab int) string {
	// Regex para substituições
	replacements := []struct {
		re   *regexp.Regexp
		repl string
	}{
		{regexp.MustCompile(`^\s*\d+\.\s(.+)`), strings.Repeat("\t", tab) + `\item $1`},
		{regexp.MustCompile(`^\s*- (.+)`), strings.Repeat("\t", tab) + `\item $1`},
	}

	for _, repl := range replacements {
		line = repl.re.ReplaceAllString(line, repl.repl)
	}

	return line
}

func convertLists(line string, s *convSta.ListStack) string {
	if s.IsListType(line) {
		// Está em uma lista
		space := s.GetSpace(line)
		if len(space) > len(s.Space) || s.Ts.IsEmpty() {
			// Inicio de lista
			s.SetListType(line)
			line = replacerListMarkdownToLatex(line, s.Ts.Size())
			newLine := strings.Repeat("\t", s.Ts.Size()-1) + "\\begin{" + s.Ts.Peek().(string) + "}\n"
			line = newLine + line
		} else if len(space) < len(s.Space) {
			// Final de sub-lista
			line = replacerListMarkdownToLatex(line, s.Ts.Size())
			line += "\n" + strings.Repeat("\t", s.Ts.Size()-1) + "\\end{" + s.Ts.Pop().(string) + "}"
		} else {
			// Apenas mais um item da lista
			line = replacerListMarkdownToLatex(line, s.Ts.Size())
		}
		s.Space = space
	} else if len(line) > 2 && !s.Ts.IsEmpty() {
		// Finaliza lista por detectar linha válida fora da lista
		line = replacerListMarkdownToLatex(line, s.Ts.Size())
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
