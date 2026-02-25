package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/taubyte/bboxes/build"
)

const usage = `Usage: bboxes <build|test|publish> <language> <version>

  Subcommands:
    build    Build the image only (tag with version).
    test     Build then run a smoke test.
    publish  Build and push to the registry.

  Arguments:
    language   go | rs | as
    version    e.g. v0.1.0

Examples:
  bboxes build go v0.1.0
  bboxes publish go v0.2.0
  bboxes test rs v0.1.0
`

func main() {
	if len(os.Args) < 4 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
	subcmd := strings.ToLower(os.Args[1])
	lang := strings.ToLower(os.Args[2])
	version := os.Args[3]

	if subcmd != "build" && subcmd != "test" && subcmd != "publish" {
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n%s", subcmd, usage)
		os.Exit(1)
	}
	if lang != "go" && lang != "rs" && lang != "as" {
		fmt.Fprintf(os.Stderr, "language must be go, rs, or as (got %q)\n", lang)
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	spec := build.Spec{Lang: lang, Version: version}
	if spec.LangDir() == "" {
		fmt.Fprintf(os.Stderr, "invalid lang: %s\n", lang)
		os.Exit(1)
	}

	ctx := context.Background()
	switch subcmd {
	case "build":
		if err := build.BuildOne(ctx, wd, spec, false); err != nil {
			log.Fatal(err)
		}
	case "publish":
		if err := build.BuildOne(ctx, wd, spec, true); err != nil {
			log.Fatal(err)
		}
	case "test":
		if err := build.BuildOne(ctx, wd, spec, false); err != nil {
			log.Fatal(err)
		}
		if err := build.TestOne(ctx, spec); err != nil {
			log.Fatal(err)
		}
	}
}
