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
	statusOutput string
	statusAll    bool
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Rapt installation status and cluster health",
	Long: `Show the current status of Rapt installation and cluster health.

This command checks the status of the Rapt CRD installation, cluster connectivity,
and provides information about the current Rapt setup.

Examples:
  rapt status
  rapt status --all-namespaces
  rapt status --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rapt.ShowStatus(namespace, statusOutput, statusAll)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().StringVarP(&statusOutput, "output", "o", "table", "Output format: table, json, yaml")
	statusCmd.Flags().BoolVarP(&statusAll, "all-namespaces", "A", false, "Show status for all namespaces")
}