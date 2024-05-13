# Markdown to LaTeX Converter

🌍 *[Português](README.md) ∙ [English](README_en.md)*

This repository contains a Go program that converts Markdown files into LaTeX format. The program is designed to handle common Markdown elements such as headers, italicized and bold text, ordered and unordered lists, as well as mathematical code blocks.

## Features

The converter supports the following Markdown features:

- Headers of different levels (from `#` to `####`)
- Italicized text (`*text*`) and bold text (`**text**`)
- Ordered lists (`1. item`) and unordered lists (`- item`)
- Mathematical code blocks delimited by `$$`

## Installation

To use this converter, you must have Go installed on your machine. If you do not have Go, you can install it from the [official site](https://golang.org/dl/).

Clone this repository using the following command:

```bash
git clone https://github.com/JunioCesarFerreira/md2tex.git
cd md2tex
```

## Usage

To convert a Markdown file to LaTeX, follow the steps below:

1. Place the Markdown file you want to convert into the project directory and rename it to `input.md`.

2. Run the command:

```bash
go run main.go
```

3. The converted file will be generated with the name `output.tex` in the same directory.

4. Use of the generated code. Use the following template to generate your documents with LaTeX and the results of this converter.

```tex
\documentclass{article}
\usepackage{titlesec}

\setcounter{secnumdepth}{4}

\titleformat{\paragraph}
{\normalfont\normalsize\bfseries}{\theparagraph}{1em}{}
\titlespacing*{\paragraph}
{0pt}{3.25ex plus 1ex minus .2ex}{1.5ex plus .2ex}

\begin{document}

% Insert the generated code here!

\end{document}
```

## Contributions

Contributions to the project are welcome.

## License

This project is licensed under the [MIT License](LICENSE).