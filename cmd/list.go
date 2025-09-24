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
	listOutput string
	listAll    bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available tools in the cluster",
	Long: `List all available tools registered in the Kubernetes cluster.

This command shows all tool definitions that have been created with 'rapt add'.
You can specify output format and whether to show tools from all namespaces.

Examples:
  rapt list
  rapt list --output table
  rapt list --all-namespaces
  rapt list --output json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rapt.ListTools(namespace, listOutput, listAll)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&listOutput, "output", "o", "table", "Output format: table, json, yaml")
	listCmd.Flags().BoolVarP(&listAll, "all-namespaces", "A", false, "List tools from all namespaces")
}