// Package cmd /*
package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/mergestat/timediff"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
	taskModel "todo/tasks/models"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks in your to-do list",
	Long: `The "list" command displays all the tasks in your to-do list, including their
status (completed or not completed).

Examples:
  tasks list`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("list called")
		tasks, err := loadTasksFromCSV()
		if err != nil {
			fmt.Println("Error Loading Tasks from CSV", err)
			return
		}
		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return
		}
		fmt.Println("Tasks:")
		for i, task := range tasks {
			status := "pending"
			if task.IsCompleted {
				status = "completed"
			}
			createdTime, err := time.Parse(time.RFC3339, task.CreatedAt)
			if err != nil {
				fmt.Printf("%d. %s (Created at: Unknown) [%s]\n", i+1, task.Name, status)
				continue
			}
			friendlyTime := formatFriendlyTime(createdTime)
			fmt.Printf("%d. %s (Created: %s) [%s]\n", i+1, task.Name, friendlyTime, status)

		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loadTasksFromCSV() ([]taskModel.Task, error) {
	const taskFile = "tasks.csv"

	// Open the CSV file
	file, err := os.Open(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []taskModel.Task{}, nil // Return an empty list if the file doesn't exist
		}
		return nil, err
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Parse records into tasks
	var tasks []taskModel.Task
	for _, record := range records {
		isCompleted, _ := strconv.ParseBool(record[2])
		tasks = append(tasks, taskModel.Task{
			Name:        record[0],
			CreatedAt:   record[1],
			IsCompleted: isCompleted,
		})
	}

	return tasks, nil
}
func formatFriendlyTime(t time.Time) string {
	//now := time.Now()

	diff := timediff.TimeDiff(t)

	// Return the formatted relative time
	return diff
}
