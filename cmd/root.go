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
	var count int
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

				// Check for C-style comments (// and /* */)
				if strings.Contains(strings.TrimSpace(line), "//") {
					count++
					fmt.Printf("comment: %s | current comment count:%d\n",strings.TrimSpace(line),count)
				}
				if strings.Contains(line, "/*") && strings.Contains(line, "*/") {
					count++
				}
			}

			// Check if there was any error reading the file
			if err := scanner.Err(); err != nil {
				return err
			}
		}
		if filepath.Ext(p) == ".lua" {
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

				// Check for Lua-style comments (--)
				if strings.Contains(strings.TrimSpace(line), "--") {
					count++
					fmt.Printf("comment: %s | current comment count:%d\n",strings.TrimSpace(line),count)
				}
			}
			// Check if there was any error reading the file
			if err := scanner.Err(); err != nil {
				return err
			}
		}
		if filepath.Ext(p) == ".py" {
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

				// Check for Python-style comments (#)
				if strings.Contains(strings.TrimSpace(line), "#") {
					count++ 
					fmt.Printf("comment: %s | current comment count:%d\n",strings.TrimSpace(line),count)
				}
			}
			// Check if there was any error reading the file
			if err := scanner.Err(); err != nil {
				return err
			}
		}
		return nil
	})

	// Print the number of comments found
	fmt.Println("Comment count:", count)
}

