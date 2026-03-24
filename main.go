package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"sql-migrator/internal/convert"
	"sql-migrator/internal/dialect"
)

func main() {
	inputPath := flag.String("input", "", "path to input .sql file (required)")
	outputPath := flag.String("output", "", "path to write converted SQL (default: stdout)")
	from := flag.String("from", "", "source dialect: sqlite, postgres, mariadb or mysql (same), or sqlserver (required)")
	to := flag.String("to", "", "target dialect: sqlite, postgres, mariadb or mysql (same), or sqlserver (required)")
	flag.Parse()

	if strings.TrimSpace(*from) == "" || strings.TrimSpace(*to) == "" || strings.TrimSpace(*inputPath) == "" {
		fmt.Fprintln(os.Stderr, "usage: sql-migrator -input file.sql -from <dialect> -to <dialect> [-output out.sql]")
		flag.PrintDefaults()
		os.Exit(2)
	}

	src, err := dialect.Parse(*from)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid -from: %v\n", err)
		os.Exit(2)
	}
	dst, err := dialect.Parse(*to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid -to: %v\n", err)
		os.Exit(2)
	}

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read input: %v\n", err)
		os.Exit(1)
	}

	out, err := convert.Convert(string(data), src, dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "convert: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(*outputPath) == "" {
		_, err = io.WriteString(os.Stdout, out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write stdout: %v\n", err)
			os.Exit(1)
		}
		return
	}

	err = os.WriteFile(*outputPath, []byte(out), 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write output: %v\n", err)
		os.Exit(1)
	}
}
