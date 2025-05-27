package reporter

import (
	"encoding/json"
	"loganalyzer/internal/analyzer"
	"os"
)

func ExportResults(results []analyzer.LogResult, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}
