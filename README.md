# Markdown to LaTeX Converter

🌍 *[Português](README.md) ∙ [English](README_en.md)*

Este repositório contém um programa em Go que converte arquivos em Markdown para o formato LaTeX. O programa foi projetado para tratar elementos comuns de Markdown, como cabeçalhos, texto em itálico e negrito, listas enumeradas e não-enumeradas, bem como blocos de código matemático.

## Recursos

O conversor suporta os seguintes recursos de Markdown:

- Cabeçalhos de diferentes níveis (de `#` até `####`)
- Texto em itálico (`*texto*`) e negrito (`**texto**`)
- Listas enumeradas (`1. item`) e não-enumeradas (`- item`)
- Blocos de código matemático delimitados por `$$`

## Instalação

Para utilizar este conversor, você precisa ter o Go instalado em sua máquina. Se você ainda não tem o Go, pode instalar a partir do [site oficial](https://golang.org/dl/).

Clone este repositório utilizando o seguinte comando:

```bash
git clone https://github.com/JunioCesarFerreira/md2tex.git
cd md2tex
```

## Uso

Para converter um arquivo Markdown para LaTeX, siga os passos abaixo:

1. Coloque o arquivo Markdown que você deseja converter dentro do diretório do projeto e renomeie-o para `input.md`.

2. Execute o comando:

```bash
go run main.go
```

3. O arquivo convertido será gerado com o nome `output.tex` no mesmo diretório.

4. Uso do código gerado. Utilize o seguinte modelo para gerar seus documentos com LaTeX e o resultado deste conversor.

```tex
\documentclass{article}
\usepackage{titlesec}

\setcounter{secnumdepth}{4}

\titleformat{\paragraph}
{\normalfont\normalsize\bfseries}{\theparagraph}{1em}{}
\titlespacing*{\paragraph}
{0pt}{3.25ex plus 1ex minus .2ex}{1.5ex plus .2ex}

\begin{document}

% Insira o código gerado aqui!

\end{document}
```

## Futuras Melhorias

- Permitir uso listas dentro de listas.

## Contribuições

Contribuições para o projeto são bem-vindas.

## Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE).
