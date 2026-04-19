package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envoy-diff/internal/diff"
	"github.com/yourorg/envoy-diff/internal/formatter"
	"github.com/yourorg/envoy-diff/internal/loader"
)

func main() {
	var (
		stagingFile    = flag.String("staging", "", "path to staging .env file (required)")
		productionFile = flag.String("production", "", "path to production .env file (required)")
		jsonOutput     = flag.Bool("json", false, "output diff as JSON")
		summaryOnly    = flag.Bool("summary", false, "print summary line only")
	)
	flag.Parse()

	if *stagingFile == "" || *productionFile == "" {
		fmt.Fprintln(os.Stderr, "error: --staging and --production flags are required")
		flag.Usage()
		os.Exit(1)
	}

	staging, err := loader.LoadEnvFile(*stagingFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading staging file: %v\n", err)
		os.Exit(1)
	}

	production, err := loader.LoadEnvFile(*productionFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading production file: %v\n", err)
		os.Exit(1)
	}

	result := diff.Compare(staging, production)

	if *summaryOnly {
		fmt.Println(formatter.Summary(result))
		return
	}

	if *jsonOutput {
		if err := formatter.JSONWriter(os.Stdout, result); err != nil {
			fmt.Fprintf(os.Stderr, "error writing JSON: %v\n", err)
			os.Exit(1)
		}
		return
	}

	formatter.TextWriter(os.Stdout, result)
}
