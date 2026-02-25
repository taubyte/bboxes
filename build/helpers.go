package build

import (
	"os"
	"os/exec"
)

// Login runs docker login using DOCKER_USER and DOCKER_TOKEN (or DOCKER_PASSWORD) env vars.
func Login() error {
	user := os.Getenv("DOCKER_USER")
	if user == "" {
		user = os.Getenv("DOCKER_USERNAME")
	}
	token := os.Getenv("DOCKER_TOKEN")
	if token == "" {
		token = os.Getenv("DOCKER_PASSWORD")
	}
	if user == "" || token == "" {
		return exec.Command("docker", "login").Run() // use default docker config
	}
	return exec.Command("docker", "login", "-u", user, "-p", token).Run()
}
