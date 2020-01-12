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

package options

import (
	"github.com/spf13/cobra"
)

// FilenameOptions
type FilenameOptions struct {
	Filename string
}

func AddFilenameArg(cmd *cobra.Command, fo *FilenameOptions, required bool) {
	cmd.Flags().StringVarP(&fo.Filename, "filename", "f", "",
		"The path to the file to use.")
	if required {
		if err := cmd.MarkFlagRequired("filename"); err != nil {
			panic(err)
		}
	}
}
