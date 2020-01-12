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

package get

import (
	"fmt"
	"github.com/n3wscott/kubectl-duck/pkg/client"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"k8s.io/klog"
)

var (
	gray   = color.New(color.FgHiBlack)
	red    = color.New(color.FgRed)
	yellow = color.New(color.FgYellow)
	green  = color.New(color.FgGreen)
)

type Get struct {

	// Client holds the context objects to talk to the api server.
	// +required
	Client *client.Client

	// Namespace the default namespace for namespaced resources.
	// +required
	Namespace string

	Duck string

	Verbose bool
}

func (g *Get) Do() error {
	duck, err := client.ToKnownDuck(g.Duck)
	if err != nil {
		return err
	}

	crds := g.Client.Ducks(duck)
	for i, crd := range crds {
		if i > 0 {
			_, _ = fmt.Fprintln(color.Output, "")
		}

		gvrk := client.CRDToGVRK(crd)

		err := g.DoGVK(gvrk)
		if err != nil {
			return err
		}

	}

	return nil
}

func (g *Get) DoGVK(gvrk client.GroupVersionResourceKind) error {
	tbl := uitable.New()
	tbl.Separator = "  "
	tbl.AddRow("NAMESPACE", "NAME", "READY", "REASON", "AGE")

	krs := g.Client.KResources(gvrk.GVR())

	for _, obj := range krs {
		klog.V(6).Infof("found %v.%v/%v", obj.Kind, obj.Name)
		rc := obj.Status.GetCondition("Ready")

		ready := "-"
		reason := ""
		if rc != nil {
			ready = string(rc.Status)
			reason = rc.Reason
		}

		var readyColor *color.Color
		switch ready {
		case "True":
			readyColor = green
		case "False":
			readyColor = red
		case "Unknown":
			readyColor = yellow
		default:
			readyColor = gray
		}
		if ready == "" {
			ready = "-"
		}

		tbl.AddRow(obj.GetNamespace(), fmt.Sprintf("%s/%s",
			obj.Kind,
			color.New(color.Bold).Sprint(obj.GetName())),
			readyColor.Sprint(ready),
			readyColor.Sprint(reason),
			client.Age(&obj))
	}

	if len(krs) > 0 {
		_, _ = fmt.Fprintln(color.Output, gray.Sprintf("%s", gvrk))
		_, _ = fmt.Fprintln(color.Output, tbl)
	}
	return nil
}
