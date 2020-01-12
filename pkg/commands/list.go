/*
Copyright 2020 Scott Nichols <author@n3wscott.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package commands

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/n3wscott/kubectl-duck/pkg/client"
	"github.com/n3wscott/kubectl-duck/pkg/list"
)

func addList(topLevel *cobra.Command, ns string, c *client.Client) {

	duck := ""

	cmd := &cobra.Command{
		Use:       "list",
		ValidArgs: []string{},
		Short:     "List CustomResourceDefinitions that implement ducktypes.",
		Example: `
  To list the known ducktypes:
  $ kubectl duck list

  To list which resources a given ducktype maps to:
  $ kubectl duck list <ducktype>
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) >= 1 {
				duck = args[0]
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Build up command.
			l := &list.List{
				Duck:      duck,
				Client:    c,
				Namespace: ns,
			}

			// Run it.
			if err := l.Do(); err != nil {
				return fmt.Errorf("failed to run list command: %w", err)
			}
			return nil
		},
	}

	// options.AddVerboseArg(cmd, vo)
	// options.AddFilenameArg(cmd, fo, true)

	topLevel.AddCommand(cmd)
}
