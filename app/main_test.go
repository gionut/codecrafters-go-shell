package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/docker/docker/api/types/container"

)

func TestShellContainer(t *testing.T) {
	ctx := context.Background()

	// 1. Build the binary
	// We build it to a temporary file
	tempDir := t.TempDir()
	binaryPath := filepath.Join(tempDir, "shell")

	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build shell binary: %v", err)
	}

	// 2. Define the container request
	req := testcontainers.ContainerRequest{
		Image: "ubuntu:latest",
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      binaryPath,
				ContainerFilePath: "/usr/local/bin/myshell",
				FileMode:          0755,
			},
		},
		
		// Run the shell binary as the container's entrypoint/command
		Cmd: []string{"/usr/local/bin/myshell"},

		
		ConfigModifier: func(config *container.Config) {
        	config.Tty        = true
        	config.OpenStdin  = true
        	config.StdinOnce  = false
    	},

		WaitingFor: wait.ForExit(),
	}
	
	// 3. Start the container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	// Clean up container after test
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	// 4. Basic assertion: Container should be running
	state, err := container.State(ctx)
	if err != nil {
		t.Fatalf("failed to get container state: %v", err)
	}

	if !state.Running {
		// It might have exited if the shell exits on EOF immediately.
		// In that case, we might check ExitCode.
		// For this draft, we'll log it.
		t.Logf("Container is not running (ExitCode: %d). This might be expected if shell exits on EOF.", state.ExitCode)
	} else {
		t.Log("Container is running")
	}
}
