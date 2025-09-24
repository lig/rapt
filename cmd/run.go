/*
Copyright © 2025 Serge Matveenko <lig@countzero.co>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	
	"codeberg.org/lig/rapt/internal/app/rapt"
	"github.com/spf13/cobra"
)

var (
	runArgs    []string
	runEnv     []string
	runWait    bool
	runFollow  bool
	runTimeout int
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <tool-name>",
	Short: "Execute a tool by creating a Kubernetes Job",
	Long: `Execute a tool by creating a Kubernetes Job from the tool definition.

This command creates a Kubernetes Job that runs the specified tool with the given arguments and environment variables.

Examples:
  rapt run echo-tool --arg message="Hello World"
  rapt run db-migrate --arg database=production --arg script=migration.sql
  rapt run file-processor --env DEBUG=true --wait
  rapt run my-tool --arg input=file.txt --arg output=result.txt --follow`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolName := args[0]
		
		// Parse arguments into key-value pairs
		argMap := make(map[string]string)
		for _, arg := range runArgs {
			parts := splitKeyValue(arg, "=")
			if len(parts) != 2 {
				return fmt.Errorf("invalid argument format: %s (expected key=value)", arg)
			}
			argMap[parts[0]] = parts[1]
		}
		
		// Parse environment variables
		envMap := make(map[string]string)
		for _, env := range runEnv {
			parts := splitKeyValue(env, "=")
			if len(parts) != 2 {
				return fmt.Errorf("invalid environment variable format: %s (expected key=value)", env)
			}
			envMap[parts[0]] = parts[1]
		}
		
		return rapt.RunTool(namespace, toolName, argMap, envMap, runWait, runFollow, runTimeout)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringArrayVarP(&runArgs, "arg", "a", nil, "Tool argument in the form key=value. Can be specified multiple times.")
	runCmd.Flags().StringArrayVarP(&runEnv, "env", "e", nil, "Environment variable in the form key=value. Can be specified multiple times.")
	runCmd.Flags().BoolVarP(&runWait, "wait", "w", false, "Wait for the job to complete before exiting")
	runCmd.Flags().BoolVarP(&runFollow, "follow", "f", false, "Follow job logs in real-time (implies --wait)")
	runCmd.Flags().IntVarP(&runTimeout, "timeout", "t", 300, "Timeout in seconds when waiting for job completion (0 = no timeout)")
}

// splitKeyValue splits a string by the first occurrence of the separator
func splitKeyValue(s, sep string) []string {
	for i := 0; i < len(s); i++ {
		if s[i:i+len(sep)] == sep {
			return []string{s[:i], s[i+len(sep):]}
		}
	}
	return []string{s}
}