package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <github-url>")
		os.Exit(1)
	}

	url := os.Args[1]
	if !strings.Contains(url, "github.com") {
		fmt.Println("Error: only GitHub URLs are supported")
		os.Exit(1)
	}

	repoName := strings.TrimPrefix(url, "https://github.com/")
	repoName = strings.TrimSuffix(repoName, "/")
	filename := strings.ReplaceAll(repoName, "/", "-") + ".txt"

	tmpDir, err := os.MkdirTemp("", "repo-*")
	if err != nil {
		fmt.Printf("Error creating temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		fmt.Printf("Error cloning repo: %v\n", err)
		os.Exit(1)
	}

	if err := writeContentFile(tmpDir, filename); err != nil {
		fmt.Printf("Error writing content: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Content saved to %s\n", filename)
}

func writeContentFile(dir, filename string) error {
	var content strings.Builder

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") || shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasPrefix(info.Name(), ".") || shouldSkipFile(info.Name()) {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		content.WriteString(fmt.Sprintf("--- File: %s ---\n", relPath))

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		content.Write(fileContent)
		content.WriteString("\n\n")
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %v", err)
	}

	return os.WriteFile(filename, []byte(content.String()), 0644)
}

func shouldSkipDir(name string) bool {
	skipDirs := []string{
		"node_modules", "vendor", "dist", "build",
	}
	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}
	return false
}

func shouldSkipFile(name string) bool {
	skipFiles := []string{
		"package-lock.json", "yarn.lock", "pnpm-lock.yaml",
		"go.sum", "Cargo.lock", "Gemfile.lock", "composer.lock",
	}
	for _, skip := range skipFiles {
		if name == skip {
			return true
		}
	}
	return false
}
