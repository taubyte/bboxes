package build

import (
	"fmt"
	"sync"
)

// Spec describes a single container image to build. Lang identifies the
// container type; Version is the tag to apply (e.g. v0.1.0).
type Spec struct {
	Lang    string // go, rs, as
	Version string
}

// ImageNameFormat is used to format org/repo:tag.
const ImageNameFormat = "%s/%s:%s"

// VersionedImage returns the full image name with version tag (e.g. taubyte/go-wasi:v0.1.0).
func (s Spec) VersionedImage() string {
	org, repo := orgRepo(s.Lang)
	return fmt.Sprintf(ImageNameFormat, org, repo, s.Version)
}

// LatestImage returns the full image name with latest tag.
func (s Spec) LatestImage() string {
	org, repo := orgRepo(s.Lang)
	return fmt.Sprintf(ImageNameFormat, org, repo, "latest")
}

// orgRepo returns organization and repo for the given language.
func orgRepo(lang string) (org, repo string) {
	org = "taubyte"
	switch lang {
	case "go":
		return org, "go-wasi"
	case "rs":
		return org, "rust-wasi"
	case "as":
		return org, "assembly-script-wasi"
	default:
		return "", ""
	}
}

// LangDir returns the containers subdir for this spec (e.g. go, rs, as). Go has a single container dir.
func (s Spec) LangDir() string {
	switch s.Lang {
	case "go":
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
