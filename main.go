package main

import (
	"bufio"
	"fmt"
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

	listOpen := false
	mbOpen := false
	listType := ""
	typeCheck := regexp.MustCompile(`^\d+\.\s.+`)
	for scanner.Scan() {
		line := scanner.Text()
		if !listOpen {
			if typeCheck.MatchString(line) {
				listType = "enumerate"
			} else {
				listType = "itemize"
			}
		}
		// Realiza substituições
		line = replacerMarkdownToLatex(line)

		// Realiza montagem de listas de itens
		line, listOpen, listType = convertLists(line, listOpen, listType)

		// Montagem de blocos matemáticos
		line, mbOpen = convertMathBlocks(line, mbOpen)

		if listOpen && len(line) < 3 {
			continue
		}
		fmt.Fprintln(writer, line)
	}

	if listOpen {
		fmt.Fprintln(writer, "\\end{"+listType+"}\n")
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
		{regexp.MustCompile(`^\d+\.\s(.+)`), `\item $1`},
		{regexp.MustCompile(`^- (.+)`), `\item $1`},
		{regexp.MustCompile(`^---$`), `\hrulefill`},
	}

	for _, repl := range replacements {
		line = repl.re.ReplaceAllString(line, repl.repl)
	}

	return line
}

func convertLists(line string, listOpen bool, listType string) (string, bool, string) {
	// Construindo listas
	if strings.HasPrefix(line, `\item`) {
		line = "\t" + line
		if !listOpen {
			// Abre uma lista
			listOpen = true
			line = "\\begin{" + listType + "}\n" + line
		}
	} else if listOpen {
		if len(line) > 2 { // Evita linhas vazias
			// Fecha uma lista
			listOpen = false
			line = "\\end{" + listType + "}\n\n" + line
			listType = ""
		}
	}
	return line, listOpen, listType
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
