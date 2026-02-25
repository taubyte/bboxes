package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/taubyte/bboxes/build"
)

const usage = `Usage: bboxes <build|test|publish> <language> [sub] <version>

  Subcommands:
    build    Build the image only (tag with version).
    test     Build then run a smoke test.
    publish  Build and push to the registry.

  Arguments:
    language   go | rs | as
    sub        func | lib  (only for go; default func)
    version    e.g. v0.1.0

Examples:
  bboxes build go func v0.1.0
  bboxes publish go lib v0.2.0
  bboxes test rs v0.1.0
`

func main() {
	if len(os.Args) < 4 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
	subcmd := strings.ToLower(os.Args[1])
	lang := strings.ToLower(os.Args[2])
	arg3 := strings.ToLower(os.Args[3])

	var sub, version string
	if arg3 == "func" || arg3 == "lib" {
		// bboxes build go func v0.1.0
		if len(os.Args) < 5 {
			fmt.Fprint(os.Stderr, usage)
			os.Exit(1)
		}
		sub = arg3
		version = os.Args[4]
	} else {
		// bboxes build rs v0.1.0  (sub defaults to func)
		sub = "func"
		version = arg3
	}

	if subcmd != "build" && subcmd != "test" && subcmd != "publish" {
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n%s", subcmd, usage)
		os.Exit(1)
	}
	if lang != "go" && lang != "rs" && lang != "as" {
		fmt.Fprintf(os.Stderr, "language must be go, rs, or as (got %q)\n", lang)
		os.Exit(1)
	}
	if lang != "go" {
		sub = "func"
	} else if sub != "func" && sub != "lib" {
		fmt.Fprintf(os.Stderr, "sub must be func or lib for go (got %q)\n", sub)
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	spec := build.Spec{Lang: lang, Sub: sub, Version: version}
	if spec.LangDir() == "" {
		fmt.Fprintf(os.Stderr, "invalid lang/sub: %s/%s\n", lang, sub)
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
