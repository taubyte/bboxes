package build

import (
	"context"
	"fmt"
	"os/exec"
)

// TestOne runs a smoke test for the built image: docker run --rm <image> true.
// The image must already be built (e.g. by BuildOne). Fails if the container exits non-zero.
func TestOne(ctx context.Context, spec Spec) error {
	image := spec.VersionedImage()
	cmd := exec.CommandContext(ctx, "docker", "run", "--rm", image, "true")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("smoke test failed for %s: %w\n%s", image, err, out)
	}
	return nil
}
