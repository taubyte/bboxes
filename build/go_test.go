//go:build linux

package build

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

const (
	goWasiImage   = "taubyte/go-wasi:local-test"
	rustWasiImage = "taubyte/rust-wasi:local-test"
)

func findModuleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find module root (go.mod)")
		}
		dir = parent
	}
}

// ensureGoWasiImage builds taubyte/go-wasi:local-test from this project if the image is not present.
func ensureGoWasiImage(t *testing.T, root string) {
	t.Helper()
	if err := exec.CommandContext(t.Context(), "docker", "image", "inspect", goWasiImage).Run(); err == nil {
		return
	}
	t.Logf("image %s not found, building from project...", goWasiImage)
	spec := Spec{Lang: "go", Version: "local-test"}
	if err := BuildOne(t.Context(), root, spec, false); err != nil {
		t.Fatalf("build image %s: %v", goWasiImage, err)
	}
}

func TestGoBuild(t *testing.T) {
	root := findModuleRoot(t)
	ensureGoWasiImage(t, root)

	runGoBuild := func(t *testing.T, name, srcSubdir string) {
		t.Helper()
		srcDir := filepath.Join(root, "testdata", "go", srcSubdir)
		outDir := t.TempDir()

		srcAbs, err := filepath.Abs(srcDir)
		if err != nil {
			t.Fatalf("abs src: %v", err)
		}
		outAbs, err := filepath.Abs(outDir)
		if err != nil {
			t.Fatalf("abs out: %v", err)
		}

		if _, err := os.Stat(srcDir); err != nil {
			t.Fatalf("testdata src dir missing: %v", err)
		}

		cmd := exec.CommandContext(t.Context(), "docker", "run", "--rm",
			"-v", srcAbs+":/src",
			"-v", outAbs+":/out",
			"-e", "OUT=/out", "-e", "SRC=/src",
			"-w", "/src",
			goWasiImage,
			"/bin/sh", "/src/.taubyte/build.sh",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("docker run: %v\n%s", err, out)
		}

		artifact := filepath.Join(outDir, "artifact.wasm")
		if _, err := os.Stat(artifact); err != nil {
			t.Fatalf("artifact.wasm not produced: %v\n%s", err, out)
		}
	}

	t.Run("wasm-ping", func(t *testing.T) {
		runGoBuild(t, "wasm-ping", "wasm-ping")
	})

	t.Run("wasm-ping-old-style", func(t *testing.T) {
		runGoBuild(t, "wasm-ping-old-style", "wasm-ping-old-style")
	})
}

// ensureRustWasiImage builds taubyte/rust-wasi:local-test from this project if the image is not present.
func ensureRustWasiImage(t *testing.T, root string) {
	t.Helper()
	if err := exec.CommandContext(t.Context(), "docker", "image", "inspect", rustWasiImage).Run(); err == nil {
		return
	}
	t.Logf("image %s not found, building from project...", rustWasiImage)
	spec := Spec{Lang: "rs", Version: "local-test"}
	if err := BuildOne(t.Context(), root, spec, false); err != nil {
		t.Fatalf("build image %s: %v", rustWasiImage, err)
	}
}

func TestRustBuild(t *testing.T) {
	root := findModuleRoot(t)
	ensureRustWasiImage(t, root)

	runRustBuild := func(t *testing.T, name, srcSubdir string) {
		t.Helper()
		srcDir := filepath.Join(root, "testdata", "rust", srcSubdir)
		outDir := t.TempDir()

		srcAbs, err := filepath.Abs(srcDir)
		if err != nil {
			t.Fatalf("abs src: %v", err)
		}
		outAbs, err := filepath.Abs(outDir)
		if err != nil {
			t.Fatalf("abs out: %v", err)
		}

		if _, err := os.Stat(srcDir); err != nil {
			t.Fatalf("testdata src dir missing: %v", err)
		}

		cmd := exec.CommandContext(t.Context(), "docker", "run", "--rm",
			"-v", srcAbs+":/src",
			"-v", outAbs+":/out",
			"-e", "OUT=/out", "-e", "SRC=/src",
			"-w", "/src",
			rustWasiImage,
			"/bin/sh", "/src/.taubyte/build.sh",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("docker run: %v\n%s", err, out)
		}

		artifact := filepath.Join(outDir, "artifact.wasm")
		if _, err := os.Stat(artifact); err != nil {
			t.Fatalf("artifact.wasm not produced: %v\n%s", err, out)
		}
	}

	t.Run("wasm-ping", func(t *testing.T) {
		runRustBuild(t, "wasm-ping", "wasm-ping")
	})
}

// TestRewritePackage verifies that the rewriter inside the image changes
// non-main package declarations to "package main" in the mounted source.
// It runs the rewriter on a temp dir containing a single package lib file,
// then checks the file content.
func TestRewritePackage(t *testing.T) {
	root := findModuleRoot(t)
	ensureGoWasiImage(t, root)

	dir := t.TempDir()
	libGo := filepath.Join(dir, "lib.go")
	const content = "package lib\n\nfunc F() {}\n"
	if err := os.WriteFile(libGo, []byte(content), 0644); err != nil {
		t.Fatalf("write lib.go: %v", err)
	}
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("abs dir: %v", err)
	}

	// Run rewriter in container; then cat the file so we can assert on output.
	cmd := exec.CommandContext(t.Context(), "docker", "run", "--rm",
		"-v", dirAbs+":/src",
		"-e", "SRC=/src",
		goWasiImage,
		"/bin/sh", "-c", "/utils/rewrite-package /src && cat /src/lib.go",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("docker run (rewrite + cat): %v\n%s", err, out)
	}

	if !bytes.Contains(out, []byte("package main")) {
		t.Errorf("rewritten file should contain 'package main'; got:\n%s", out)
	}
	if bytes.Contains(out, []byte("package lib")) {
		t.Errorf("rewritten file should not contain 'package lib'; got:\n%s", out)
	}
}
