package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
log files to analyze, processes them, and outputs a report.`,
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

		var results []LogResult

		fmt.Println("\n--- Log Configurations to Analyze ---")
		for _, cfg := range configs {
			fmt.Printf("ID: %s, Path: %s, Type: %s\n", cfg.ID, cfg.Path, cfg.Type)

			results = append(results, LogResult{
				LogID:        cfg.ID,
				FilePath:     cfg.Path,
				Status:       "SIMULATED_OK",
				Message:      "Simulated analysis successful.",
				ErrorDetails: "",
			})
		}

		if outputPath != "" {
			fmt.Printf("\nSimulating export of results to %s\n", outputPath)
			err := exportResults(outputPath, results)
			if err != nil {
				fmt.Printf("Error exporting results: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Simulated export complete.")
		}
	},
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
