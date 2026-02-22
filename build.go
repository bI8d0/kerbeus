package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// Create build directory if it doesn't exist
	buildDir := filepath.Join(currentDir, "build")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		fmt.Printf("Error creating build directory: %v\n", err)
		return
	}

	// Executable name
	executableName := "kerbeus"

	// File name to compile
	filename := "main.go"

	// Compile for Linux
	fmt.Println("Compiling for Linux...")
	if err := buildForOS("linux", filepath.Join(buildDir, executableName), filename, currentDir); err != nil {
		fmt.Printf("Error compiling for Linux: %v\n", err)
	}

	fmt.Println("Compilation completed. Binaries are located in the 'build' directory.")
}

func buildForOS(goos, outputName, filename string, currentDir string) error {
	args := []string{"build", "-o", outputName}

	// Add compilation tag to exclude Windows-specific code on Linux
	if goos == "linux" {
		args = append(args, "-tags", "!windows")
	}

	args = append(args, filename)

	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(),
		"GOOS="+goos,
		"GOARCH=amd64",
	)
	cmd.Dir = currentDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error compiling for %s: %v", goos, err)
	}

	fmt.Printf("Compilation for %s completed.\n", goos)

	// Verify that the executable was created
	if _, err := os.Stat(outputName); os.IsNotExist(err) {
		return fmt.Errorf("the executable for %s was not generated correctly", goos)
	}

	return nil
}
