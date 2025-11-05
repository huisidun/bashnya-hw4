
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

func Process(lines []string, opts Options) ([]string, error) {
	if (opts.Count && (opts.Duplicates || opts.UniqueOnly)) ||
		(opts.Duplicates && opts.UniqueOnly) {
		return nil, fmt.Errorf("нельзя использовать флаги -c, -d, -u одновременно")
	}

	normalize := func(s string) string {
		fields := strings.Fields(s)
		if opts.SkipFields > 0 {
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

	counts := make(map[string]int)
	order := []string{}
	original := make(map[string]string)

	for _, line := range lines {
		key := normalize(line)
		if counts[key] == 0 {
			order = append(order, key)
			original[key] = line
		}
		counts[key]++
	}

	var result []string
	for _, key := range order {
		cnt := counts[key]
		orig := original[key]

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

	return result, nil
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

func main() {
	var (
		count      = flag.Bool("c", false, "подсчитать количество")
		dups       = flag.Bool("d", false, "только повторяющиеся")
		uniques    = flag.Bool("u", false, "только уникальные")
		ignoreCase = flag.Bool("i", false, "игнорировать регистр")
		skipFields = flag.Int("f", 0, "пропустить полей")
		skipChars  = flag.Int("s", 0, "пропустить символов")
	)

	flag.Parse()

	if (*count && (*dups || *uniques)) || (*dups && *uniques) {
		fmt.Fprintln(os.Stderr, "ошибка: нельзя использовать -c, -d, -u вместе")
		os.Exit(1)
	}

	args := flag.Args()

	var input *os.File = os.Stdin
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка открытия: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		input = f
	}

	var output *os.File = os.Stdout
	if len(args) > 1 {
		f, err := os.Create(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ошибка создания: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		output = f
	}

	opts := Options{
		Count:      *count,
		Duplicates: *dups,
		UniqueOnly: *uniques,
		IgnoreCase: *ignoreCase,
		SkipFields: *skipFields,
		SkipChars:  *skipChars,
	}

	if err := Run(input, output, opts); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
		os.Exit(1)
	}
}