package main

import (
	"bytes"
	"fmt"
	"gopagelink/configs"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Load configuration
	config, err := configs.LoadSiteConfig("config.yml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Generate HTML
	err = generateHTML(config)
	if err != nil {
		log.Fatalf("Error generating HTML: %v", err)
	}

	err = copyAssets(config.Theme)
	if err != nil {
		log.Fatalf("Error copying and minifying assets: %v", err)
	}

	fmt.Println("Site generated successfully!")
}

func generateHTML(config *configs.SiteConfig) error {
	themeFile := fmt.Sprintf("themes/%s/index.html",
		config.Theme)

	// Load HTML template
	tmpl, err := template.ParseFiles(themeFile)
	if err != nil {
		return err
	}

	// Open output file
	outputFile, err := os.Create("index.html")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Define data to pass to the template
	data := struct {
		Config *configs.SiteConfig
	}{
		Config: config,
	}

	// Execute template with data
	return tmpl.Execute(outputFile, data)
}

func copyAssets(theme string) error {
	// Create assets directories if they don't exist
	if err := os.MkdirAll("assets/css", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create assets/css directory: %w", err)
	}
	if err := os.MkdirAll("assets/js", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create assets/js directory: %w", err)
	}
	if err := os.MkdirAll("assets/icons", os.ModePerm); err != nil {
		return fmt.Errorf("failed to create assets/icons directory: %w", err)
	}

	// Get theme assets files
	cssFiles, err := filepath.Glob(fmt.Sprintf("themes/%s/assets/css/*.css", theme))
	if err != nil {
		return fmt.Errorf("failed to list CSS files: %w", err)
	}
	if err := copyFiles(cssFiles, "assets/css"); err != nil {
		return fmt.Errorf("failed to copy css files: %w", err)
	}

	jsFiles, err := filepath.Glob(fmt.Sprintf("themes/%s/assets/js/*.js", theme))
	if err != nil {
		return fmt.Errorf("failed to list JS files: %w", err)
	}
	if err := copyFiles(jsFiles, "assets/js"); err != nil {
		return fmt.Errorf("failed to copy js files: %w", err)
	}

	// Copy favicons
	iconsFiles, err := filepath.Glob(fmt.Sprintf("themes/%s/assets/icons/*", theme))
	if err != nil {
		return fmt.Errorf("failed to list icon files: %w", err)
	}
	if err := copyFiles(iconsFiles, "assets/icons"); err != nil {
		return fmt.Errorf("failed to copy icon files: %w", err)
	}

	return nil
}

func copyFiles(files []string, outputDir string) error {
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		outputPath := filepath.Join(outputDir, filepath.Base(file))
		outFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", outputPath, err)
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, bytes.NewReader(data)); err != nil {
			return fmt.Errorf("failed to copy data to %s: %w", outputPath, err)
		}
	}

	return nil
}
