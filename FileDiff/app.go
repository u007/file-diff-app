package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type DiffLine struct {
	LineNumber int    `json:"lineNumber"`
	Content    string `json:"content"`
	Type       string `json:"type"` // "same", "added", "removed"
}

type DiffResult struct {
	LeftLines  []DiffLine `json:"leftLines"`
	RightLines []DiffLine `json:"rightLines"`
}

// CompareFiles compares two files and returns the differences
func (a *App) CompareFiles(leftPath string, rightPath string) (*DiffResult, error) {
	// Read files
	leftContent, err := os.ReadFile(leftPath)
	if err != nil {
		return nil, err
	}

	rightContent, err := os.ReadFile(rightPath)
	if err != nil {
		return nil, err
	}

	// Initialize diff match patch
	dmp := diffmatchpatch.New()

	// Get diffs
	diffs := dmp.DiffMain(string(leftContent), string(rightContent), false)

	// Process diffs into lines
	leftLines := []DiffLine{}
	rightLines := []DiffLine{}
	leftLineNum := 1
	rightLineNum := 1

	for _, diff := range diffs {
		lines := strings.Split(diff.Text, "\n")
		
		for i, line := range lines {
			if i == len(lines)-1 && line == "" {
				continue
			}

			switch diff.Type {
			case diffmatchpatch.DiffEqual:
				leftLines = append(leftLines, DiffLine{
					LineNumber: leftLineNum,
					Content:    line,
					Type:      "same",
				})
				rightLines = append(rightLines, DiffLine{
					LineNumber: rightLineNum,
					Content:    line,
					Type:      "same",
				})
				leftLineNum++
				rightLineNum++
			case diffmatchpatch.DiffDelete:
				leftLines = append(leftLines, DiffLine{
					LineNumber: leftLineNum,
					Content:    line,
					Type:      "removed",
				})
				leftLineNum++
			case diffmatchpatch.DiffInsert:
				rightLines = append(rightLines, DiffLine{
					LineNumber: rightLineNum,
					Content:    line,
					Type:      "added",
				})
				rightLineNum++
			}
		}
	}

	return &DiffResult{
		LeftLines:  leftLines,
		RightLines: rightLines,
	}, nil
}

// ApplyChanges applies changes from one file to another
func (a *App) ApplyChanges(targetPath string, sourcePath string, direction string) error {
	sourceContent, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	return os.WriteFile(targetPath, sourceContent, 0644)
}

// MergeFiles merges changes from both files
func (a *App) MergeFiles(leftPath string, rightPath string, outputPath string) error {
	// For now, we'll implement a simple merge that takes the right file's content
	// In a real implementation, this should handle conflicts and merging properly
	rightContent, err := os.ReadFile(rightPath)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, rightContent, 0644)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
