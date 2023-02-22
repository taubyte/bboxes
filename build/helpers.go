package build

import (
	"fmt"
	"os/exec"

	wasmSpec "github.com/taubyte/go-specs/builders/wasm"
	imageSpec "github.com/taubyte/go-specs/builders/wasm/images"
)

func version(lang wasmSpec.SupportedLanguage) (string, error) {
	switch lang {
	case wasmSpec.Go:
		return goImageVersion, nil
	case wasmSpec.AssemblyScript:
		return assemblyScriptVersion, nil
	case wasmSpec.Rust:
		return rustImageVersion, nil
	default:
		return "", fmt.Errorf("`%s` is not a supported language", lang)
	}
}

func login() error {
	user, err := imageSpec.UserEnvVar.Get()
	if err != nil {
		return err
	}

	token, err := imageSpec.TokenEnvVar.Get()
	if err != nil {
		return err
	}

	return exec.Command("docker", []string{"login", "-u", user, "-p", token}...).Run()
}
