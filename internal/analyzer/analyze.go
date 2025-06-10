package analyzer

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"Efrei_Tp_loganizer/internal/config"
	"Efrei_Tp_loganizer/internal/reporter"
)

type FileNotFoundError struct {
	Path string
}

func (e FileNotFoundError) Error() string {
	return fmt.Sprintf("file not found or inaccessible: %s", e.Path)
}

type ParsingError struct {
	Details string
}

func (e ParsingError) Error() string {
	return fmt.Sprintf("parsing error: %s", e.Details)
}

func AnalyzeLogs(configs []config.LogConfig) []reporter.LogResult {
	var wg sync.WaitGroup
	resultsChan := make(chan reporter.LogResult, len(configs))

	for _, cfg := range configs {
		wg.Add(1)
		go func(cfgInstance config.LogConfig) {
			defer wg.Done()
			result := analyzeLogFile(cfgInstance)
			resultsChan <- result
		}(cfg)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var logResults []reporter.LogResult
	for result := range resultsChan {
		logResults = append(logResults, result)
	}

	return logResults
}

func analyzeLogFile(cfg config.LogConfig) reporter.LogResult {
	processingTime := time.Duration(50+rand.Intn(150)) * time.Millisecond
	time.Sleep(processingTime)

	result := reporter.LogResult{
		LogID:       cfg.ID,
		FilePath:    cfg.Path,
		ProcessTime: processingTime.String(),
	}

	if _, err := os.Stat(cfg.Path); os.IsNotExist(err) {
		fileErr := FileNotFoundError{Path: cfg.Path}
		result.Status = "FAILED"
		result.Message = "File access failed"
		result.ErrorDetails = fileErr.Error()
		return result
	} else if err != nil {
		result.Status = "FAILED"
		result.Message = "File access failed"
		result.ErrorDetails = fmt.Sprintf("could not access file: %v", err)
		return result
	}

	data, err := os.ReadFile(cfg.Path)
	if err != nil {
		result.Status = "FAILED"
		result.Message = "File read failed"
		result.ErrorDetails = err.Error()
		return result
	}

	if len(data) == 0 {
		parseErr := ParsingError{Details: "empty log file, no content to parse"}
		result.Status = "FAILED"
		result.Message = "Parsing failed"
		result.ErrorDetails = parseErr.Error()
		return result
	}

	if rand.Intn(100) < 10 {
		parseErr := ParsingError{Details: "simulated random parsing failure"}
		result.Status = "FAILED"
		result.Message = "Random parsing error occurred"
		result.ErrorDetails = parseErr.Error()
		return result
	}

	linesCount := len(strings.Split(string(data), "\n"))
	var analysisDetails string

	switch cfg.Type {
	case "nginx-access":
		analysisDetails = fmt.Sprintf("Nginx access log: %d entries processed.", linesCount)
	case "mysql-error":
		analysisDetails = fmt.Sprintf("MySQL error log: %d error entries found.", linesCount)
	case "custom-app":
		analysisDetails = fmt.Sprintf("Custom application log: %d log entries processed.", linesCount)
	default:
		analysisDetails = fmt.Sprintf("Generic log: %d lines processed.", linesCount)
	}

	result.Status = "SUCCESS"
	result.Message = analysisDetails
	result.ErrorDetails = ""

	return result
}

func HandleError(err error) {
	var fileNotFoundErr FileNotFoundError
	var parseErr ParsingError

	if errors.As(err, &fileNotFoundErr) {
		fmt.Printf("Custom File Not Found Error: %s\n", fileNotFoundErr.Path)
	} else if errors.As(err, &parseErr) {
		fmt.Printf("Custom Parsing Error: %s\n", parseErr.Details)
	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Standard OS 'Does Not Exist' Error: %v\n", err)
	} else {
		fmt.Printf("Other Error: %v\n", err)
	}
}
