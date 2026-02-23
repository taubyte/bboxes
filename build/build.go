package build

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pterm/pterm"
)

// BuildOne builds a single image from spec. If push is true, logs in and pushes to the registry.
// Uses a directory build context (no tarballs). Runs docker build and optionally docker push.
func BuildOne(ctx context.Context, wd string, spec Spec, push bool) error {
	var infoMsgs infoMessages
	var errMsgs errMsg

	contextDir, cleanup, err := PrepareContext(wd, spec)
	if err != nil {
		return err
	}
	defer cleanup()

	versionedImage := spec.VersionedImage()
	// docker build -f <contextDir>/Dockerfile -t <versionedImage> <contextDir>
	cmd := exec.CommandContext(ctx, "docker", "build", "-f", filepath.Join(contextDir, "Dockerfile"), "-t", versionedImage, contextDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	infoMsgs.appendMsg("built %s", pterm.White(versionedImage))

	if push {
		if err := Login(); err != nil {
			errMsgs.append("docker login: %w", err)
		} else {
			cmd := exec.CommandContext(ctx, "docker", "push", versionedImage)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				errMsgs.append("docker push %s: %w", versionedImage, err)
			} else {
				infoMsgs.appendMsg("pushed %s", pterm.White(versionedImage))
			}
		}
	}

	for _, m := range infoMsgs.msgs {
		m.print()
	}
	if errMsgs.err != nil {
		return errMsgs.err
	}
	return nil
}
