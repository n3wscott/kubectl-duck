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

package list

import (
	"fmt"
	"github.com/n3wscott/kubectl-duck/pkg/client"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"k8s.io/klog"
)

var (
	gray  = color.New(color.FgHiBlack)
	red   = color.New(color.FgRed)
	green = color.New(color.FgGreen)
)

type List struct {

	// Client holds the context objects to talk to the api server.
	// +required
	Client *client.Client

	// Namespace the default namespace for namespaced resources.
	// +required
	Namespace string

	Duck string

	Verbose bool
}

func (l *List) Do() error {
	if l.Duck == "" {
		return l.DoDucks()
	}

	duck, err := client.ToKnownDuck(l.Duck)
	if err != nil {
		return err
	}

	return l.DoCRDs(duck)
}

func (l *List) DoDucks() error {
	tbl := uitable.New()
	tbl.Separator = "  "
	tbl.AddRow("NAME", "SHORT NAME", "SELECTOR")
	for _, duck := range client.KnownDucks() {
		tbl.AddRow(client.DuckName(duck), client.DuckShortName(duck), duck)
	}

	_, _ = fmt.Fprintln(color.Output, tbl)
	return nil
}

func (l *List) DoCRDs(duck string) error {
	crds := l.Client.Ducks(duck)

	tbl := uitable.New()
	tbl.Separator = "  "
	tbl.AddRow("NAME", "KIND", "DUCK", "CREATED AT")

	for _, crd := range crds {
		gvrk := client.CRDToGVRK(crd)
		klog.V(6).Infof("found :%v", gvrk)
		tbl.AddRow(crd.Name, gvrk.Kind, client.DuckShortName(duck), client.Age(&crd))
	}

	_, _ = fmt.Fprintln(color.Output, tbl)
	return nil
}
