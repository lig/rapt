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
	"errors"
	"strings"

	"codeberg.org/lig/rapt/internal/app/rapt"
	"github.com/spf13/cobra"
)

// Add command flags
var (
	addImage   string
	addCommand string
	addEnv     []string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <tool-name>",
	Short: "Add a new tool definition to your Kubernetes cluster.",
	Long: `Add a tool (containerized job/command) to the Rapt system in your Kubernetes cluster.

This command registers a new tool by specifying its container image, the command to run inside the image, and (optionally) a set of environment variables.

Examples:
  rapt add lstool -i alpine --command "ls -la"
  rapt add echo --image busybox -e FOO=bar -e BAZ=qux --command "echo $FOO $BAZ"
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, e := range addEnv {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" {
				return errors.New("each --env/-e argument must be in NAME=VALUE format")
			}
		}
		toolName := args[0]
		return rapt.Add(namespace, toolName, addImage, addCommand, addEnv)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addImage, "image", "i", "", "(Required) Container image to run.")
	addCmd.MarkFlagRequired("image")
	addCmd.Flags().StringVarP(&addCommand, "command", "c", "", "Command to run (overrides ENTRYPOINT). Specify as a single string.")
	addCmd.Flags().StringArrayVarP(&addEnv, "env", "e", nil, "Environment variable in the form NAME=VALUE. Can be specified multiple times.")
}
