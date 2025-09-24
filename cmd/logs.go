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
	"codeberg.org/lig/rapt/internal/app/rapt"
	"github.com/spf13/cobra"
)

var (
	logsFollow bool
	logsTail   int
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs <tool-name> [job-name]",
	Short: "View logs from tool executions",
	Long: `View logs from tool executions.

This command can list previous job runs for a tool or display logs for a specific job.
If no job name is provided, it lists all previous runs for the tool.
If a job name is provided, it displays the logs for that specific job.

Examples:
  rapt logs echo-tool                    # List previous runs for echo-tool
  rapt logs echo-tool echo-tool-20250115-143022  # Show logs for specific job
  rapt logs db-migrate --follow         # Follow logs for the latest job
  rapt logs my-tool --tail 100          # Show last 100 lines of logs`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		toolName := args[0]
		var jobName string
		if len(args) > 1 {
			jobName = args[1]
		}
		
		return rapt.ShowLogs(namespace, toolName, jobName, logsFollow, logsTail)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow logs in real-time (only for specific job)")
	logsCmd.Flags().IntVarP(&logsTail, "tail", "t", 0, "Number of lines to show from the end of logs (0 = all)")
}