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
		line, listOpen, listType, mbOpen = convertMarkdownToLatex(line, listOpen, listType, mbOpen)
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

func convertMarkdownToLatex(line string, listOpen bool, listType string, mbOpen bool) (string, bool, string, bool) {
	// Regex para substituições
	replacements := []struct {
		re   *regexp.Regexp
		repl string
	}{
		{regexp.MustCompile(`^# (.+)`), `\section{$1}`},
		{regexp.MustCompile(`^## (.+)`), `\subsection{$1}`},
		{regexp.MustCompile(`^### (.+)`), `\subsubsection{$1}`},
		{regexp.MustCompile(`^#### (.+)`), `\paragraph{$1}`},
		{regexp.MustCompile(`\*\*(.+?)\*\*`), `\textbf{$1}`},
		{regexp.MustCompile(`\*(.+?)\*`), `\textit{$1}`},
		{regexp.MustCompile(`^\d+\.\s(.+)`), `\item $1`},
		{regexp.MustCompile(`^- (.+)`), `\item $1`},
		{regexp.MustCompile(`^---$`), `\hrulefill`},
	}

	for _, repl := range replacements {
		line = repl.re.ReplaceAllString(line, repl.repl)
	}

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

	// Montagem de blocos matemáticos
	line, mbOpen = handleMathBlocks(line, mbOpen)

	return line, listOpen, listType, mbOpen
}

func handleMathBlocks(line string, mathBlockOpen bool) (string, bool) {
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
