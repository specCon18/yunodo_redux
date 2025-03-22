/*
Copyright © 2025 Steven Carpenter <steven.carpenter@skdevstudios.com>

*/
package cmd

import ( //TODO: P:5 THIS IS A TEST FOR INLINE
    "fmt"
    "os"
    "io/fs"
    "path/filepath"
    "github.com/spf13/cobra"
    "strings"
    "strconv"
    "bufio"
    "regexp"
)

type Comment struct {
    FilePath     string
    LineNumber   int
    Priority     int
    Task         string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "yunodo_redux",
    Short: "Tooling to make dealing with TODO statements in comments easier to manage.",
    Long: `YUNODO redux is a better written version of a tool I wrote for internal use.
This application is a tool that takes comments that start with TODO: and outputs them in a given format.

format for comment is as follows:
    <COMMENT_SYMBOL>TODO: P<PRIORITY_VALUE> <COMMENT>

where COMMENT_SYMBOL is your languages comment symbol,
PRIORITY_VALUE is a priority from 0-9 where 9 is least critical and 0 is needs patched NOW!!!,
and PATH_TO_PROJECT is the path to your codebase's root.

Usage:
    yunodo_redux -p <PATH_TO_PROJECT>

`,
    // Uncomment the following line if your bare application
    // has an action associated with it:
    Run: func(cmd *cobra.Command, args []string) {
        path, _ := cmd.Flags().GetString("path")
	    if path != "" {
		comments := readFileTree(path)
		table := mdTableFormatter(comments)
		fmt.Print(table)
	    } else {
	        fmt.Println("ERROR: No path provided, please provide path using the -p flag")
	        os.Exit(1)
	    }
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    err := rootCmd.Execute()
    if err != nil {
	os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().StringP("path", "p", "","assign the path to get comments from.")
}

// Function to extract priority from the TODO comment
func extractPriority(line string) (int, error) {
    re := regexp.MustCompile(`P:(\d)`)
    matches := re.FindStringSubmatch(line)

    if len(matches) > 1 {
        // Convert the priority string (matches[1]) to an integer
        priority, err := strconv.Atoi(matches[1])
        if err != nil {
            return 0, err
        }
        return priority, nil
    }
    return 9, nil // Default priority if none is found
}

// Function to extract the task from the TODO comment and remove the priority
func extractTask(line string) string {
    trimmedLine := strings.TrimSpace(line)
    trimmedLine = strings.TrimPrefix(trimmedLine, "//TODO: ")

    // Remove the priority part (P:<digit>) if it exists
    re := regexp.MustCompile(`P:\d+`)
    trimmedLine = re.ReplaceAllString(trimmedLine, "")

    return strings.TrimSpace(trimmedLine) // Return the task without priority
}

// Function to extract comment from a line that might have code before the comment
func extractCommentFromLine(path, p string, ln int, line string) (Comment, error) {
    // Use a regular expression to look for a TODO comment anywhere in the line
    re := regexp.MustCompile(`//TODO:(.*)`)
    matches := re.FindStringSubmatch(line)

    if len(matches) > 0 {
        // Extract priority and task using helper functions
        priority, err := extractPriority(matches[0])
        if err != nil {
            return Comment{}, err
        }

        task := extractTask(matches[0])

        // Return the Comment object
        return Comment{
            FilePath:  filepath.Join(path, p),
            LineNumber: ln,
            Priority: priority,
            Task:     task,
        }, nil
    }
    return Comment{}, nil // No comment found, return an empty Comment
}

func readFileTree(path string) []Comment {
    fsys := os.DirFS(path)
    var comments []Comment

    // Walk through all files in the directory
    fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        // Check if file has a valid extension
	ext := filepath.Ext(p)
	if ext == ".go" || ext == ".rs" || ext == ".java" || ext == ".c"{
            // Read the file and check for comments
            file, err := os.Open(filepath.Join(path, p))
            if err != nil {
                return err
            }
            defer file.Close()

            // Process each line in the file
            scanner := bufio.NewScanner(file)
            ln := 0
            for scanner.Scan() {
                line := scanner.Text()
                ln++
                // Use the new helper function to extract the comment
                comment, err := extractCommentFromLine(path, p, ln, line)
                if err != nil {
                    return err
                }
                if (comment != Comment{}) {
                    comments = append(comments, comment)
                }
            }

            // Check if there was any error reading the file
            if err := scanner.Err(); err != nil {
                return err
            }
        }
        return nil
    })
    return comments
}

func mdTableFormatter(comments []Comment) string {
	table := fmt.Sprintf("| File Path | Line Number | Priority | Task |\n| --------- | ----------- | -------- | ---- |\n")
	for _, comment := range comments {
    		// You can access each `comment` here
		row := fmt.Sprintf("| %s | %d | %d | %s |\n",comment.FilePath,comment.LineNumber,comment.Priority,comment.Task)
		table = table + row
	}
	return table
}
