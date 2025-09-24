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
	deleteForce bool
	deleteAll   bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <tool-name>",
	Short: "Delete a tool definition from the cluster",
	Long: `Delete a tool definition from the Kubernetes cluster.

This command removes the specified tool definition, making it no longer available for execution.
You can delete multiple tools at once or use --all to delete all tools in the namespace.

Examples:
  rapt delete echo-tool
  rapt delete tool1 tool2 tool3
  rapt delete --all
  rapt delete --all --force`,
	Args: func(cmd *cobra.Command, args []string) error {
		if deleteAll && len(args) > 0 {
			return fmt.Errorf("cannot specify tool names when using --all")
		}
		if !deleteAll && len(args) == 0 {
			return fmt.Errorf("must specify tool name(s) or use --all")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if deleteAll {
			return rapt.DeleteAllTools(namespace, deleteForce)
		}
		return rapt.DeleteTools(namespace, args, deleteForce)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Skip confirmation prompt")
	deleteCmd.Flags().BoolVarP(&deleteAll, "all", "a", false, "Delete all tools in the namespace")
}