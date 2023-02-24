package main

import (
	"os"
	"path"

	"github.com/taubyte/bboxes/build"
	"github.com/taubyte/go-specs/builders/wasm"
)

//go:generate bash containers/build.sh

func main() {
	// builds, err := setDefaultToBuild(true, true, true)
	// if err != nil {
	// 	panic(err)
	// }

	err := build.Build(false, false, []build.Image{
		build.CustomImage{
			TarPath:      "/home/tafkhan/Documents/Work/Taubyte/Repos/bboxes/containers/_builds/go_test_examples.tar",
			Organization: "taubyte",
			Repo:         "go-wasi",
			Version:      "test-examples",
		},
	})
	if err != nil {
		panic(err)
	}
	// app := &cli.App{
	// 	Name:  "bbox",
	// 	Usage: "Build and Push Docker Images",
	// 	Flags: []cli.Flag{
	// 		&cli.BoolFlag{
	// 			Name:    "all",
	// 			Aliases: []string{"a"},
	// 			Action: func(ctx *cli.Context, b bool) error {
	// 				if b {
	// 					toBuild, err := setDefaultToBuild(true, true, true)
	// 					if err != nil {
	// 						return err
	// 					}

	// 					return build.Build(false, true, toBuild)
	// 				}

	// 				return nil
	// 			},
	// 		},
	// 	},
	// }

	// if err := app.Run(os.Args); err != nil {
	// 	log.Fatal(err)
	// }
}

func setDefaultToBuild(rust, golang, as bool) (toBuild []build.Image, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	wd = path.Join(wd, "containers")
	if rust {
		_build, err := build.New(wasm.Rust, wd)
		if err != nil {
			return nil, err
		}
		toBuild = append(toBuild, _build)
	}

	if golang {
		_build, err := build.New(wasm.Go, wd)
		if err != nil {
			return nil, err
		}
		toBuild = append(toBuild, _build)
	}

	if as {
		_build, err := build.New(wasm.AssemblyScript, wd)
		if err != nil {
			return nil, err
		}
		toBuild = append(toBuild, _build)
	}

	return
}
