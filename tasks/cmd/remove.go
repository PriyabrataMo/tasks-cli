/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// removeCmd represents the remove command

var removeCmd = &cobra.Command{
	Use:   "remove [task name]",
	Short: "Remove a task from your to-do list",
	Long: `The "remove" command allows you to remove a specific task from your to-do list.
You need to specify the exact task name as an argument.

Examples:
  tasks remove "Buy groceries"
  tasks remove "Prepare presentation for Monday"`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("remove called")
		taskName := args[0]

		if err := removeTaskFromCSV(taskName); err != nil {
			fmt.Println("Task Not Found...", err)
			return
		}

		fmt.Printf("Task \"%s\" removed successfully\n", taskName)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func removeTaskFromCSV(taskName string) error {
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
			continue // Skip adding this line to the output
		}

		// Add the line to the remaining lines
		remainingLines = append(remainingLines, line)
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
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing to file: %w", err)
	}

	return nil
}
