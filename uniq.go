package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
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

func normalizeLine(line string, opts Options) string {
	s := line

	// Пропуск полей
	if opts.SkipFields > 0 {
		fields := strings.Fields(s)
		if opts.SkipFields < len(fields) {
			s = strings.Join(fields[opts.SkipFields:], " ")
		} else {
			s = ""
		}
	}

	// Пропуск символов (только если строка не пустая)
	if opts.SkipChars > 0 && len(s) > 0 {
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

func buildResult(order []string, counts map[string]int, originals map[string]string, opts Options) []string {
	var result []string
	for _, key := range order {
		cnt := counts[key]
		orig := originals[key]

		if opts.Duplicates && cnt < 2 {
			continue
		}
		if opts.UniqueOnly && cnt != 1 {
			continue
		}

		if opts.Count {
			result = append(result, fmt.Sprintf("%d %s", cnt, orig))
		} else {
			result = append(result, orig)
		}
	}
	return result
}

func Process(lines []string, opts Options) ([]string, error) {
	if (opts.Count && (opts.Duplicates || opts.UniqueOnly)) ||
		(opts.Duplicates && opts.UniqueOnly) {
		return nil, fmt.Errorf("нельзя использовать флаги -c, -d, -u одновременно")
	}

	counts := make(map[string]int)
	order := []string{}
	originals := make(map[string]string)

	for _, line := range lines {
		key := normalizeLine(line, opts)
		if counts[key] == 0 {
			order = append(order, key)
			originals[key] = line
		}
		counts[key]++
	}

	return buildResult(order, counts, originals, opts), nil
}

func Run(r io.Reader, w io.Writer, opts Options) error {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	result, err := Process(lines, opts)
	if err != nil {
		return err
	}

	for _, line := range result {
		fmt.Fprintln(w, line)
	}
	return nil
}

func parseFlags() (Options, error) {
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
		return Options{}, fmt.Errorf("нельзя использовать флаги -c, -d, -u вместе")
	}

	return Options{
		Count:      *count,
		Duplicates: *dups,
		UniqueOnly: *uniques,
		IgnoreCase: *ignoreCase,
		SkipFields: *skipFields,
		SkipChars:  *skipChars,
	}, nil
}

func openInput(args []string) (*os.File, error) {
	if len(args) == 0 {
		return os.Stdin, nil
	}
	return os.Open(args[0])
}

func createOutput(args []string) (*os.File, error) {
	if len(args) < 2 {
		return os.Stdout, nil
	}
	return os.Create(args[1])
}

func main() {
	opts, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
		os.Exit(1)
	}

	args := flag.Args()
	input, err := openInput(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка открытия входного файла: %v\n", err)
		os.Exit(1)
	}
	defer input.Close()

	output, err := createOutput(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка создания выходного файла: %v\n", err)
		os.Exit(1)
	}
	defer output.Close()

	if err := Run(input, output, opts); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
		os.Exit(1)
	}
}