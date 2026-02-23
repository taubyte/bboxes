package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// PrepareContext copies containers/<langDir> then containers/common into a temp
// dir, copies the lang Dockerfile as "Dockerfile", and returns the temp dir path
// and a cleanup function. Merge order: lang first, then common (common overwrites).
func PrepareContext(wd string, spec Spec) (contextDir string, cleanup func(), err error) {
	langDir := spec.LangDir()
	if langDir == "" {
		return "", nil, fmt.Errorf("invalid lang/sub: %s/%s", spec.Lang, spec.Sub)
	}
	tmp, err := os.MkdirTemp("", "bboxes-build-")
	if err != nil {
		return "", nil, err
	}
	cleanup = func() { _ = os.RemoveAll(tmp) }

	containers := filepath.Join(wd, "containers")
	langPath := filepath.Join(containers, langDir)
	commonPath := filepath.Join(containers, "common")
	dockerfileSrc := filepath.Join(langPath, "Dockerfile")

	for _, src := range []string{langPath, commonPath} {
		if err := copyDir(src, tmp); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("copy %s: %w", src, err)
		}
	}
	// Copy Dockerfile to temp root as "Dockerfile"
	df, err := os.ReadFile(dockerfileSrc)
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("read Dockerfile: %w", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "Dockerfile"), df, 0644); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("write Dockerfile: %w", err)
	}
	return tmp, cleanup, nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(dest, info.Mode())
		}
		return copyFile(path, dest)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
