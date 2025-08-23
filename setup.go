package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	name := "crawlx"
	fmt.Println("🔧 Installing", name, "on", runtime.GOOS)

	// Determine the source and destination paths based on the OS.
	var sourcePath, destPath string
	if runtime.GOOS == "windows" {
		sourcePath = "dist\\" + name + ".exe"
		destPath = "C:\\Windows\\System32\\" + name + ".exe"
	} else {
		sourcePath = "dist/" + name
		destPath = "/usr/local/bin/" + name
	}

	// Check if the pre-built binary exists in the 'dist' folder.
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		fmt.Printf("❌ Pre-built binary not found at %s. Please build the project first.\n", sourcePath)
		os.Exit(1)
	}

	// Copy the pre-built binary to the system path.
	fmt.Println("📂 Copying pre-built binary to", destPath)
	var copyCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		copyCmd = exec.Command("powershell", "Copy-Item", "-Path", sourcePath, "-Destination", destPath, "-Force")
	} else {
		// Use a shell to run 'sudo' for proper permissions on Unix-like systems.
		copyCmd = exec.Command("sh", "-c", "sudo cp "+sourcePath+" "+destPath)
	}

	copyCmd.Stdout = os.Stdout
	copyCmd.Stderr = os.Stderr

	if err := copyCmd.Run(); err != nil {
		fmt.Println("❌ Installation failed:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Installation complete. You can now run '" + name + "' from any folder!")
}
