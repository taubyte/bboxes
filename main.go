package main

import (
	"log"
	"os"
	"path"

	"github.com/taubyte/bboxes/build"
)

//go:generate bash containers/build.sh

var (
	wd string
)

func init() {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
}

type releaseType string

const (
	production    releaseType = "production"
	test_examples releaseType = "test_examples"

	buildsRelDir = "containers/_builds"
)

func (r releaseType) String() string {
	return string(r)
}

func (r releaseType) Version() string {
	switch r {
	case production:
		return "v0"
	case test_examples:
		return "test-examples"
	default:
		return "testing"
	}
}

func goImage(release releaseType, version string) build.Image {
	return build.CustomImage{
		TarPath:      path.Join(wd, buildsRelDir, release.String(), "go.tar"),
		Organization: "taubyte",
		Repo:         "go-wasi",
		Version:      version,
	}
}

func rustImage(release releaseType, version string) build.Image {
	return build.CustomImage{
		TarPath:      path.Join(wd, buildsRelDir, release.String(), "rs.tar"),
		Organization: "taubyte",
		Repo:         "rust-wasi",
		Version:      version,
	}
}

func assemblyImage(release releaseType, version string) build.Image {
	return build.CustomImage{
		TarPath:      path.Join(wd, buildsRelDir, release.String(), "rs.tar"),
		Organization: "taubyte",
		Repo:         "assembly-script-wasi",
		Version:      version,
	}
}

func main() {
	err := build.Build(true, true, []build.Image{
		goImage(production, "v0"),
	})
	if err != nil {
		panic(err)
	}
}

// func setDefaultToBuild(rust, golang, as bool) (toBuild []build.Image, err error) {
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		return nil, err
// 	}

// 	wd = path.Join(wd, "containers")
// 	if rust {
// 		_build, err := build.New(wasm.Rust, wd)
// 		if err != nil {
// 			return nil, err
// 		}
// 		toBuild = append(toBuild, _build)
// 	}

// 	if golang {
// 		_build, err := build.New(wasm.Go, wd)
// 		if err != nil {
// 			return nil, err
// 		}
// 		toBuild = append(toBuild, _build)
// 	}

// 	if as {
// 		_build, err := build.New(wasm.AssemblyScript, wd)
// 		if err != nil {
// 			return nil, err
// 		}
// 		toBuild = append(toBuild, _build)
// 	}

// 	return
// }
