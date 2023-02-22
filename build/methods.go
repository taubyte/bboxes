package build

import (
	"fmt"
	"io"
	"os"

	"github.com/pterm/pterm"
	"github.com/taubyte/go-specs/builders/wasm/images"
)

func (i *infoMessages) appendMsg(format string, args ...interface{}) {
	i.msgs = append(i.msgs, infoMessage{format: format, args: args})
}

func (i infoMessage) print() {
	pterm.Info.Printfln(i.format, i.args...)
}

func (e *errMsg) append(format string, args ...any) {
	err := fmt.Errorf(format, args...)
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.err != nil {
		e.err = fmt.Errorf("%s && %s", e.err, err)
		return
	}

	e.err = err
}

// Tarball opens the DefaultImage tarball
func (d defaultImage) Tarball() (io.Reader, error) {
	return os.Open(d.TarBallPath(d.wd))
}

// VersionedImage returns the DefaultImage formatted image name with a version tag
func (d defaultImage) VersionedImage() (string, error) {
	version, err := version(d.Language())
	return d.Image(version), err
}

// LatestImage returns the DefaultImage formatted image name with latest tag
func (d defaultImage) LatestImage() string {
	return d.Image("latest")
}

// Tarball opens the given TarPath
func (c CustomImage) Tarball() (io.Reader, error) {
	return os.Open(c.TarPath)
}

// VersionedImage formats the given Organization, Repo, and Version
func (c CustomImage) VersionedImage() (string, error) {
	return fmt.Sprintf(images.ImageNameFormat, c.Organization, c.Repo, c.Version), nil
}

// LatestImage formats the given Organization, and Repo, with a `latest` tag
func (c CustomImage) LatestImage() string {
	return fmt.Sprintf(images.ImageNameFormat, c.Organization, c.Repo, "latest")
}
