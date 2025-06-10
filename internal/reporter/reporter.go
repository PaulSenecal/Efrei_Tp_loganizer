package reporter

import (
	"encoding/json"
	"fmt"

	"loganalyzer/internal/analyzer"
	"os"
)

func ExportResults(results []analyzer.LogResult, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("échec de la création du fichier de sortie : %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}
