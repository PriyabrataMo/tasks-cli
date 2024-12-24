// Package cmd /*
package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// doneCmd represents the done command
var doneCmd = &cobra.Command{
	Use:   "done [task name]",
	Short: "Mark a task as completed",
	Long: `The "done" command allows you to mark a specific task as completed in your to-do list.
You need to specify the exact task name as an argument.

Examples:
  tasks done "Buy groceries"
  tasks done "Prepare presentation for Monday"`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("done called")
		taskName := args[0]
		if err := markAsDoneInCSV(taskName); err != nil {
			fmt.Println("Failed to mark as Done... ", err)
			return
		}

		fmt.Printf("Task \"%s\" is Marked as Done!\n", taskName)
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// doneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// doneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func markAsDoneInCSV(taskName string) error {
	const taskFile = "tasks.csv"
	file, err := os.Open(taskFile)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Use bufio.Scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	var remainingLines []string
	found := false

	for scanner.Scan() {
		line := scanner.Text()
		// Parse the CSV record
		record, err := csv.NewReader(strings.NewReader(line)).Read()
		if err != nil {
			return fmt.Errorf("error parsing line as CSV: %w", err)
		}

		// Check if the task name matches
		if record[0] == taskName {
			found = true
			// Mark as done by setting IsCompleted to "true"
			if len(record) >= 3 {
				record[2] = "true"
			} else {
				return fmt.Errorf("invalid record format: %v", record)
			}
		}

		// Rebuild the CSV line
		var builder strings.Builder
		writer := csv.NewWriter(&builder)
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing updated record: %w", err)
		}
		writer.Flush()

		// Append the updated line
		remainingLines = append(remainingLines, builder.String())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// If the task wasn't found, return an error
	if !found {
		return fmt.Errorf("task with name '%s' not found", taskName)
	}

	// Write the remaining lines back to the file
	outputFile, err := os.Create(taskFile) // Truncate the file
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	for _, line := range remainingLines {
		_, err := writer.WriteString(line)
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing to file: %w", err)
	}

	return nil
}
