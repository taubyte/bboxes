package build

import (
	"context"
	"os/exec"
	"sync"

	"github.com/pterm/pterm"
	ci "github.com/taubyte/go-simple-container"
)

// Build will build the given images, and push if requested
func Build(push bool, latest bool, toBuild []Image) error {
	ci.ForceRebuild = true

	c, err := ci.New()
	if err != nil {
		return err
	}

	ctx := context.Background()
	var wg sync.WaitGroup
	var infoMsgs infoMessages
	var errMsgs errMsg

	wg.Add(len(toBuild))
	for _, build := range toBuild {
		go func(_infoMsgs *infoMessages, _image Image, _errMsgs *errMsg) {
			defer wg.Done()
			tarball, err := _image.Tarball()
			if err != nil {
				_errMsgs.append("opening tarball failed with: %s", err)
				return
			}

			versionedImage, err := _image.VersionedImage()
			if err != nil {
				_errMsgs.append("getting versioned image failed with: %s", err)
				return
			}

			if _, err = c.Image(ctx, versionedImage, ci.Build(tarball)); err != nil {
				_errMsgs.append("initializing image `%s` failed with: %s", versionedImage, err)
				return
			}

			builtImages := []string{versionedImage}
			_infoMsgs.appendMsg("`%s` image successfully built", pterm.White(versionedImage))

			if latest {
				latestImage := _image.LatestImage()
				if err := exec.Command("docker", []string{"tag", versionedImage, latestImage}...).Run(); err != nil {
					_errMsgs.append("tagging `%s` as `%s` failed with: %s", versionedImage, latestImage, err)
					return
				}

				builtImages = append(builtImages, latestImage)
				_infoMsgs.appendMsg("`%s` image tagged", pterm.White(latestImage))
			}

			if push {
				if err = login(); err != nil {
					_errMsgs.append("docker log in failed with: %s", err)
					return
				}

				for _, builtImage := range builtImages {
					if err := exec.Command("docker", []string{"push", builtImage}...).Run(); err != nil {
						_errMsgs.append("pushing image `%s` to docker hub failed with: %s", builtImage, err)
						return
					}

					_infoMsgs.appendMsg("`%s` pushed to docker hub", pterm.White(versionedImage))
				}
			}

		}(&infoMsgs, build, &errMsgs)
	}

	wg.Wait()

	for _, msg := range infoMsgs.msgs {
		msg.print()
	}

	if err = errMsgs.err; err != nil {
		return err
	}

	return nil
}
