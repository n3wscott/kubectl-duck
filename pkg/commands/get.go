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
	"github.com/n3wscott/kubectl-duck/pkg/client"
	"github.com/spf13/cobra"

	"github.com/n3wscott/kubectl-duck/pkg/get"
)

func addGet(topLevel *cobra.Command, ns string, c *client.Client) {

	duck := ""

	cmd := &cobra.Command{
		Use:       "get",
		ValidArgs: []string{},
		Short:     "Get the resource instances related to a ducktype.",
		Example: `
  To get resource instances that are of the given ducktype:
  $ kubectl duck get [ducktype]
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
			g := &get.Get{
				Duck:      duck,
				Client:    c,
				Namespace: ns,
			}

			// Run it.
			if err := g.Do(); err != nil {
				return fmt.Errorf("failed to run get command: %w", err)
			}
			return nil
		},
	}

	// options.AddVerboseArg(cmd, vo)
	// options.AddFilenameArg(cmd, fo, true)

	topLevel.AddCommand(cmd)
}
