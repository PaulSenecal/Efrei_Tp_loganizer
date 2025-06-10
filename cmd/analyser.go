package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type FileNotFoundError struct {
	Path string
}

func (e FileNotFoundError) Error() string {
	return "file not found: " + e.Path
}

type ParsingError struct {
	Details string
}

func (e ParsingError) Error() string {
	return "parsing error: " + e.Details
}

type LogConfig struct {
	ID   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type LogResult struct {
	LogID        string `json:"log_id"`
	FilePath     string `json:"file_path"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	ErrorDetails string `json:"error_details"`
	ProcessTime  string `json:"process_time"`
}

var (
	configPath string
	outputPath string
)

func init() {
	analyzeCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to the JSON configuration file")
	analyzeCmd.MarkFlagRequired("config")

	analyzeCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to export the analysis report to JSON")

	rootCmd.AddCommand(analyzeCmd)
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze log files based on a configuration file.",
	Long: `The analyze command reads a JSON configuration file specifying
log files to analyze, processes them concurrently, and outputs a report.`,
	Run: func(cmd *cobra.Command, args []string) {
		if configPath == "" {
			fmt.Println("Error: --config flag is required.")
			os.Exit(1)
		}

		configs, err := readConfigs(configPath)
		if err != nil {
			fmt.Printf("Error reading configuration file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully loaded %d log configurations from %s.\n", len(configs), configPath)
		fmt.Println("Starting concurrent analysis...")

		startTime := time.Now()
		results := AnalyzeLogs(configs)
		totalTime := time.Since(startTime)

		fmt.Printf("\nAnalysis completed in %v\n", totalTime)
		fmt.Printf("Processed %d log files\n", len(results))

		fmt.Println("\n--- Analysis Results ---")
		successCount := 0
		errorCount := 0

		for _, result := range results {
			fmt.Printf("ID: %s, Status: %s, Time: %s\n", result.LogID, result.Status, result.ProcessTime)
			if result.Status == "SUCCESS" {
				successCount++
			} else {
				errorCount++
				fmt.Printf("  Error: %s\n", result.ErrorDetails)
			}
		}

		fmt.Printf("\nSummary: %d successful, %d errors\n", successCount, errorCount)

		if outputPath != "" {
			fmt.Printf("\nExporting results to %s\n", outputPath)
			err := exportResults(outputPath, results)
			if err != nil {
				fmt.Printf("Error exporting results: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Export complete.")
		}
	},
}

func AnalyzeLogs(configs []LogConfig) []LogResult {
	var wg sync.WaitGroup
	results := make(chan LogResult, len(configs))

	// Lancement des goroutines
	for _, config := range configs {
		wg.Add(1)
		go func(cfg LogConfig) {
			defer wg.Done()
			result := analyzeLogFile(cfg)
			results <- result
		}(config)
	}

	// Fermeture du channel
	go func() {
		wg.Wait()
		close(results)
	}()

	var logResults []LogResult
	for result := range results {
		logResults = append(logResults, result)
	}

	return logResults
}

func analyzeLogFile(config LogConfig) LogResult {
	startTime := time.Now()

	// Simulation du temps de traitement (50-200ms)
	processingTime := time.Duration(50+rand.Intn(150)) * time.Millisecond
	time.Sleep(processingTime)

	result := LogResult{
		LogID:       config.ID,
		FilePath:    config.Path,
		ProcessTime: processingTime.String(),
	}

	if _, err := os.Stat(config.Path); err != nil {
		if os.IsNotExist(err) {
			fileErr := FileNotFoundError{Path: config.Path}
			result.Status = "ERROR"
			result.Message = "File analysis failed"
			result.ErrorDetails = fileErr.Error()
			return result
		}
	}

	data, err := os.ReadFile(config.Path)
	if err != nil {
		result.Status = "ERROR"
		result.Message = "File read failed"
		result.ErrorDetails = err.Error()
		return result
	}

	if len(data) == 0 {
		parseErr := ParsingError{Details: "empty log file"}
		result.Status = "ERROR"
		result.Message = "File analysis failed"
		result.ErrorDetails = parseErr.Error()
		return result
	}

	if rand.Intn(100) < 10 {
		parseErr := ParsingError{Details: "random parsing failure simulation"}
		result.Status = "ERROR"
		result.Message = "Random error occurred"
		result.ErrorDetails = parseErr.Error()
		return result
	}

	// Simulation d'analise
	linesCount := len(strings.Split(string(data), "\n"))
	var analysisDetails string

	switch config.Type {
	case "nginx-access":
		analysisDetails = fmt.Sprintf("Nginx access log analyzed: %d entries processed", linesCount)
	case "mysql-error":
		analysisDetails = fmt.Sprintf("MySQL error log analyzed: %d error entries found", linesCount)
	case "custom-app":
		analysisDetails = fmt.Sprintf("Custom application log analyzed: %d log entries processed", linesCount)
	default:
		analysisDetails = fmt.Sprintf("Generic log analyzed: %d lines processed", linesCount)
	}

	result.Status = "SUCCESS"
	result.Message = analysisDetails
	result.ErrorDetails = ""

	return result
}

func readConfigs(path string) ([]LogConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var configs []LogConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return nil, fmt.Errorf("could not unmarshal config JSON: %w", err)
	}
	return configs, nil
}

func exportResults(path string, results []LogResult) error {
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal results to JSON: %w", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write results to file: %w", err)
	}
	return nil
}

// Fonction utilitaire pour vérifier les types d'erreurs avec errors.Is()
func handleError(err error) {
	var fileNotFoundErr FileNotFoundError
	var parseErr ParsingError

	if errors.As(err, &fileNotFoundErr) {
		fmt.Printf("File not found error: %s\n", fileNotFoundErr.Path)
	} else if errors.As(err, &parseErr) {
		fmt.Printf("Parsing error: %s\n", parseErr.Details)
	} else {
		fmt.Printf("Unknown error: %v\n", err)
	}
}
