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

// purgeCmd represents the purge command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Remove the Rapt CRD and all associated resources from your Kubernetes cluster.",
	Long: `Purge all Rapt-related resources from your Kubernetes cluster, including the Rapt CustomResourceDefinition (CRD).

This command is useful for cleanup purposes when you no longer wish to use Rapt in a given cluster. 
It will remove the CRD and all associated custom resources managed by Rapt.

⚠️ Warning: This operation is destructive and cannot be undone. Ensure you have backups if needed before proceeding.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return rapt.Purge(namespace)
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)
}
