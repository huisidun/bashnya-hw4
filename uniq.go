package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Options struct {
	Count       bool
	Duplicates  bool
	UniqueOnly  bool
	IgnoreCase  bool
	SkipFields  int
	SkipChars   int
}

func normalize(line string, opts Options) string {
	s := line
	if opts.SkipFields > 0 {
		fields := strings.Fields(s)
		if opts.SkipFields < len(fields) {
			s = strings.Join(fields[opts.SkipFields:], " ")
		} else {
			s = ""
		}
	}
	if opts.SkipChars > 0 {
		if opts.SkipChars < len(s) {
			s = s[opts.SkipChars:]
		} else {
			s = ""
		}
	}
	if opts.IgnoreCase {
		s = strings.ToLower(s)
	}
	return s
}

func shouldOutput(count int, opts Options) bool {
	if opts.Duplicates {
		return count >= 2
	}
	if opts.UniqueOnly {
		return count == 1
	}
	return true
}

func processLines(lines []string, opts Options, output *os.File) {
	if len(lines) == 0 {
		return
	}

	prevNorm := normalize(lines[0], opts)
	prevOrig := lines[0]
	countGroup := 1

	for i := 1; i < len(lines); i++ {
		currNorm := normalize(lines[i], opts)
		if currNorm == prevNorm {
			countGroup++
		} else {
			if shouldOutput(countGroup, opts) {
				if opts.Count {
					fmt.Fprintf(output, "%d %s\n", countGroup, prevOrig)
				} else {
					fmt.Fprintln(output, prevOrig)
				}
			}
			prevNorm = currNorm
			prevOrig = lines[i]
			countGroup = 1
		}
	}

	if shouldOutput(countGroup, opts) {
		if opts.Count {
			fmt.Fprintf(output, "%d %s\n", countGroup, prevOrig)
		} else {
			fmt.Fprintln(output, prevOrig)
		}
	}
}

func parseFlags() (Options, []string) {
	var (
		count      = flag.Bool("c", false, "подсчитать количество")
		dups       = flag.Bool("d", false, "только повторяющиеся")
		uniques    = flag.Bool("u", false, "только уникальные")
		ignoreCase = flag.Bool("i", false, "игнорировать регистр")
		skipFields = flag.Int("f", 0, "пропустить первые N полей")
		skipChars  = flag.Int("s", 0, "пропустить первые N символов")
	)

	flag.Parse()

	if (*count && (*dups || *uniques)) || (*dups && *uniques) {
		fmt.Fprintln(os.Stderr, "ошибка: нельзя использовать флаги -c, -d, -u вместе")
		os.Exit(1)
	}

	return Options{
		Count:      *count,
		Duplicates: *dups,
		UniqueOnly: *uniques,
		IgnoreCase: *ignoreCase,
		SkipFields: *skipFields,
		SkipChars:  *skipChars,
	}, flag.Args()
}

func openFiles(args []string) (*os.File, *os.File) {
	input := os.Stdin
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка открытия файла: %v\n", err)
			os.Exit(1)
		}
		input = f
	}

	output := os.Stdout
	if len(args) > 1 {
		f, err := os.Create(args[1])
		if err != nil {
			if input != os.Stdin {
				input.Close()
			}
			fmt.Fprintf(os.Stderr, "ошибка создания файла: %v\n", err)
			os.Exit(1)
		}
		output = f
	}

	return input, output
}

func readLines(input *os.File) []string {
	scanner := bufio.NewScanner(input)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка чтения: %v\n", err)
		os.Exit(1)
	}
	return lines
}

func main() {
	opts, args := parseFlags()
	input, output := openFiles(args)
	defer input.Close()
	defer output.Close()

	lines := readLines(input)
	processLines(lines, opts, output)
}