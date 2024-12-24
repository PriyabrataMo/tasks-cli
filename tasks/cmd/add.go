// Package cmd /*
package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
	"todo/tasks/models"
)

const taskFile = "tasks.csv"

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [task name]",
	Short: "Add a new task to your to-do list",
	Long: `The "add" command allows you to add a new task to your to-do list. 
You can specify the task name as an argument.

Examples:
  tasks add "Buy groceries"
  tasks add "Prepare presentation for Monday"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("add called")
		taskName := args[0]
		task := taskModel.Task{
			Name:        taskName,
			CreatedAt:   time.Now().Format(time.RFC3339),
			IsCompleted: false,
		}
		if err := saveTaskToCSV(task); err != nil {
			fmt.Println("Error saving task:", err)
			return
		}

		fmt.Printf("Task \"%s\" added successfully\n", taskName)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func saveTaskToCSV(task taskModel.Task) error {
	file, err := os.OpenFile(taskFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the task to the CSV file
	record := []string{task.Name, task.CreatedAt, fmt.Sprintf("%t", task.IsCompleted)}
	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}
