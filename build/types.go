package build

import (
	"io"
	"sync"

	"github.com/taubyte/go-specs/builders/wasm"
	"github.com/taubyte/go-specs/builders/wasm/images"
)

// CustomImage defines an image that is built outside of our production build images
type CustomImage struct {
	TarPath      string
	Organization string
	Repo         string
	Version      string
}

type defaultImage struct {
	images.LanguageConfig
	wd string
}

// New creates a defaultImage, a predefined Taubyte container image
func New(lang wasm.SupportedLanguage, wd string) (image defaultImage, err error) {
	config, err := images.Config(lang)
	if err != nil {
		return
	}

	return defaultImage{
		LanguageConfig: config,
		wd:             wd,
	}, nil
}

// Image defines the interface of images that can be built by bbox
type Image interface {
	Tarball() (io.Reader, error)
	VersionedImage() (string, error)
	LatestImage() string
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
