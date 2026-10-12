// Package six standalone platform binaries with documentation and SHA-256 sums.
package main

import (
	"archive/zip"
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	version := flag.String("version", "dev", "Version embedded in binaries and archive filenames")
	out := flag.String("out", "dist", "Output directory")
	flag.Parse()
	if err := build(*version, *out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func build(version, out string) error {
	if !regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`).MatchString(version) {
		return fmt.Errorf("invalid version")
	}
	if err := os.MkdirAll(out, 0755); err != nil {
		return err
	}
	scratch, err := os.MkdirTemp("", "dlut-pan-build-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(scratch)
	var sums strings.Builder
	for _, platform := range []struct{ system, arch string }{{"linux", "amd64"}, {"linux", "arm64"}, {"darwin", "amd64"}, {"darwin", "arm64"}, {"windows", "amd64"}, {"windows", "arm64"}} {
		name := "dlut-pan"
		if platform.system == "windows" {
			name += ".exe"
		}
		binary := filepath.Join(scratch, name)
		command := exec.Command("go", "build", "-trimpath", "-ldflags", "-s -w -X main.version="+version, "-o", binary, "./cmd/dlut-pan")
		for _, value := range os.Environ() {
			key := strings.ToUpper(strings.SplitN(value, "=", 2)[0])
			if key != "GOOS" && key != "GOARCH" && key != "CGO_ENABLED" {
				command.Env = append(command.Env, value)
			}
		}
		command.Env = append(command.Env, "GOOS="+platform.system, "GOARCH="+platform.arch, "CGO_ENABLED=0")
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err = command.Run(); err != nil {
			return err
		}
		archive := fmt.Sprintf("dlut-pan_%s_%s_%s.zip", version, platform.system, platform.arch)
		path := filepath.Join(out, archive)
		if err = pack(path, binary, name); err != nil {
			return err
		}
		bytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Fprintf(&sums, "%x  %s\n", sha256.Sum256(bytes), archive)
		fmt.Println("Built", archive)
	}
	return os.WriteFile(filepath.Join(out, "checksums.txt"), []byte(sums.String()), 0644)
}

func pack(destination, binary, name string) error {
	f, err := os.Create(destination)
	if err != nil {
		return err
	}
	writer := zip.NewWriter(f)
	for _, item := range []struct {
		source, name string
		mode         os.FileMode
	}{{binary, name, 0755}, {"README.md", "README.md", 0644}, {"LICENSE", "LICENSE", 0644}, {"docs/authentication.md", "docs/authentication.md", 0644}} {
		bytes, readErr := os.ReadFile(item.source)
		if readErr != nil {
			writer.Close()
			f.Close()
			return readErr
		}
		header := &zip.FileHeader{Name: item.name, Method: zip.Deflate, Modified: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
		header.SetMode(item.mode)
		entry, writeErr := writer.CreateHeader(header)
		if writeErr == nil {
			_, writeErr = entry.Write(bytes)
		}
		if writeErr != nil {
			writer.Close()
			f.Close()
			return writeErr
		}
	}
	if err = writer.Close(); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
