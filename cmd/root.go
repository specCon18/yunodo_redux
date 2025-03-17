/*
Copyright © 2025 Steven Carpenter <steven.carpenter@skdevstudios.com>

*/
package cmd

import (
    "fmt"
    "os"
    "io/fs"
    "path/filepath"
    "github.com/spf13/cobra"
    "strings"
    "bufio"
    "regexp"
)


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "v2",
    Short: "Tooling to make dealing with TODO statements in comments easier to manage.",
    Long: `YUNODO redux is a better written version of a tool I wrote for internal use.
This application is a tool that takes comments that start with TODO: and outputs them in a given format.

format for comment is as follows:
    <COMMENT_SYMBOL>TODO: P<PRIORITY_VALUE> <COMMENT>

where COMMENT_SYMBOL is your languages comment symbol,
PRIORITY_VALUE is a priority from 0-9 where 9 is least critical and 0 is needs patched NOW!!!,
and PATH_TO_PROJECT is the path to your codebase's root.

Usage:
    yunodo -p <PATH_TO_PROJECT>

`,
    // Uncomment the following line if your bare application
    // has an action associated with it:
    Run: func(cmd *cobra.Command, args []string) {
        path, _ := cmd.Flags().GetString("path")
	    if path != "" {
	    //TODO: recursivly read each file in the path and store content as string in slice
	   	readFileTree(path) 
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
func readFileTree(path string) {
    fsys := os.DirFS(path)

    // Walk through all files in the directory
    fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        // Check if file has .go extension
        if filepath.Ext(p) == ".go" {
            // Read the file and check for comments
            file, err := os.Open(filepath.Join(path, p))
            if err != nil {
                return err
            }
            defer file.Close()

            // Process each line in the file
            scanner := bufio.NewScanner(file)
            for scanner.Scan() {
                line := scanner.Text()

		//TODO: P:0 add support for TODO's not at start of line
                // Check for C-style comments (// and /* */)
                if strings.HasPrefix(strings.TrimSpace(line), "//TODO:") {
                    // Regular expression to match "P:" followed by a digit
                    re := regexp.MustCompile(`P:(\d)`)

                    // Trim spaces and remove the "//TODO:" prefix
                    trimmedLine := strings.TrimSpace(line)
                    trimmedLine = strings.TrimPrefix(trimmedLine, "//TODO: ")

                    // Try to find the priority
                    matches := re.FindStringSubmatch(trimmedLine)
    
                    // If a priority is found, remove it and print the comment and priority
                    if len(matches) > 1 {
                        // Remove the "P:<digit>" from the beginning of the line
                        trimmedLine = strings.Replace(trimmedLine, matches[0], "", 1) // Replace only the first occurrence
                        fmt.Printf("comment: %s | priority: %s \n", strings.TrimSpace(trimmedLine), matches[1])
                    } else {
                        fmt.Printf("comment: %s | priority: 9 \n", trimmedLine)
                    }
                }

            }

            // Check if there was any error reading the file
            if err := scanner.Err(); err != nil {
                return err
            }
        }
        return nil
    })
}
