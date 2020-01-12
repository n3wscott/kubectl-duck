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
	"flag"
	"fmt"
	"github.com/n3wscott/kubectl-duck/pkg/client"
	"os"

	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var cf *genericclioptions.ConfigFlags

func AddDuckCommands(topLevel *cobra.Command) error {
	cf = genericclioptions.NewConfigFlags(true)
	cf.AddFlags(topLevel.Flags())
	if err := flag.Set("logtostderr", "true"); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to set logtostderr flag: %v\n", err)
		os.Exit(1)
	}

	restConfig, err := cf.ToRESTConfig()
	if err != nil {
		return err
	}
	c, err := client.GetClient(restConfig)
	if err != nil {
		return err
	}

	ns := *cf.Namespace
	if ns == "" {
		clientConfig := cf.ToRawKubeConfigLoader()
		defaultNamespace, _, err := clientConfig.Namespace()
		if err != nil {
			defaultNamespace = "default"
		}
		ns = defaultNamespace
	}

	addList(topLevel, ns, c)
	addGet(topLevel, ns, c)
	return nil
}
