package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type LogResult struct {
	LogID        string `json:"log_id"`
	FilePath     string `json:"file_path"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorDetails string `json:"error_details"`
}

// ExportReport exporte les résultats d’analyse dans un fichier JSON
func ExportReport(outputPath string, results []LogResult) error {
	dir := filepath.Dir(outputPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("could not create directory %q: %w", dir, err)
		}
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal results to JSON: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("could not write results to file %q: %w", outputPath, err)
	}

	return nil
}
