package build

import (
	"fmt"
	"sync"
)

// Spec describes a single container image to build. Lang and Sub identify the
// container type; Version is the tag to apply (e.g. v0.1.0).
type Spec struct {
	Lang    string // go, rs, as
	Sub     string // func, lib (only for go)
	Version string
}

// ImageNameFormat is used to format org/repo:tag.
const ImageNameFormat = "%s/%s:%s"

// VersionedImage returns the full image name with version tag (e.g. taubyte/go-wasi:v0.1.0).
func (s Spec) VersionedImage() string {
	org, repo := orgRepo(s.Lang, s.Sub)
	return fmt.Sprintf(ImageNameFormat, org, repo, s.Version)
}

// LatestImage returns the full image name with latest tag.
func (s Spec) LatestImage() string {
	org, repo := orgRepo(s.Lang, s.Sub)
	return fmt.Sprintf(ImageNameFormat, org, repo, "latest")
}

// orgRepo returns organization and repo for (lang, sub). Sub is only used for go.
func orgRepo(lang, sub string) (org, repo string) {
	org = "taubyte"
	switch lang {
	case "go":
		if sub == "lib" {
			return org, "go-wasi-lib"
		}
		return org, "go-wasi"
	case "rs":
		return org, "rust-wasi"
	case "as":
		return org, "assembly-script-wasi"
	default:
		return "", ""
	}
}

// LangDir returns the containers subdir for this spec (e.g. go, go-lib, rs, as).
func (s Spec) LangDir() string {
	switch s.Lang {
	case "go":
		if s.Sub == "lib" {
			return "go-lib"
		}
		return "go"
	case "rs":
		return "rs"
	case "as":
		return "as"
	default:
		return ""
	}
}

type infoMessage struct {
	format string
	args   []interface{}
}

type infoMessages struct {
	msgs []infoMessage
	lock sync.Mutex
}

type errMsg struct {
	err  error
	lock sync.Mutex
}
