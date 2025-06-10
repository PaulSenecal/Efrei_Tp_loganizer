package cmd

import (
	"fmt"
	"os"
	"time"

	"Efrei_Tp_loganizer/internal/analyzer"
	"Efrei_Tp_loganizer/internal/config"
	"Efrei_Tp_loganizer/internal/reporter"

	"github.com/spf13/cobra"
)

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

		configs, err := config.ReadConfigs(configPath)
		if err != nil {
			fmt.Printf("Error reading configuration file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully loaded %d log configurations from %s.\n", len(configs), configPath)
		fmt.Println("Starting concurrent analysis...")

		startTime := time.Now()
		results := analyzer.AnalyzeLogs(configs)
		totalTime := time.Since(startTime)

		fmt.Printf("\nAnalysis completed in %v\n", totalTime)
		fmt.Printf("Processed %d log files\n", len(results))

		fmt.Println("\n--- Analysis Results ---")
		successCount := 0
		errorCount := 0

		for _, result := range results {
			fmt.Printf("ID: %s, Status: %s, Time: %s", result.LogID, result.Status, result.ProcessTime)
			if result.Status == "FAILED" {
				errorCount++
				fmt.Printf(", Error: %s\n", result.ErrorDetails)
			} else {
				successCount++
				fmt.Printf(", Message: %s\n", result.Message)
			}
		}

		fmt.Printf("\nSummary: %d successful, %d failed\n", successCount, errorCount)

		if outputPath != "" {
			fmt.Printf("\nExporting results to %s...\n", outputPath)
			err := reporter.ExportReport(outputPath, results)
			if err != nil {
				fmt.Printf("Error exporting results: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Export complete.")
		} else {
			fmt.Println("\nOutput path not provided. Results will not be exported to a file.")
		}
	},
}
