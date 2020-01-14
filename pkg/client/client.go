package client

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/duration"
	"k8s.io/client-go/dynamic"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // combined authprovider import
	"k8s.io/client-go/rest"
	"k8s.io/klog"

	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type Client struct {
	DynamicClient dynamic.Interface
}

var client *Client
var onceError error
var once sync.Once

func GetClient(cfg *rest.Config) (*Client, error) {
	once.Do(func() {
		cfg.QPS = 1000
		cfg.Burst = 1000
		dyn, err := dynamic.NewForConfig(cfg)
		if err != nil {
			onceError = fmt.Errorf("failed to construct dynamic client: %w", err)
			return
		}

		// Setup client.
		client = &Client{
			DynamicClient: dyn,
		}
	})
	if onceError != nil {
		return nil, onceError
	}
	return client, nil
}

// Ducks returns CRDs labeled as given.
// labelSelector should be in the form "duck.knative.dev/source=true"
func (c *Client) Ducks(labelSelector string) []v1beta1.CustomResourceDefinition {
	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.k8s.io",
		Version:  "v1beta1",
		Resource: "customresourcedefinitions",
	}
	like := v1beta1.CustomResourceDefinition{}

	list, err := c.DynamicClient.Resource(gvr).List(metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		klog.Errorf("failed to list customresourcedefinitions, %v", err)
		return nil
	}

	all := make([]v1beta1.CustomResourceDefinition, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			klog.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj
	}
	return all
}

func (c *Client) KResources(gvr schema.GroupVersionResource) []duckv1.KResource {
	like := duckv1.KResource{}

	list, err := c.DynamicClient.Resource(gvr).List(metav1.ListOptions{})
	if err != nil {
		klog.Errorf("failed to list %s, %v", gvr.String(), err)
		return nil
	}

	all := make([]duckv1.KResource, len(list.Items))

	for i, item := range list.Items {
		obj := like.DeepCopy()
		if err = runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, obj); err != nil {
			klog.Fatalf("Error DefaultUnstructuree.Dynamiconverter. %v", err)
		}
		obj.ResourceVersion = gvr.Version
		obj.APIVersion = gvr.GroupVersion().String()
		all[i] = *obj
	}
	return all
}

type GroupVersionResourceKind struct {
	Group    string
	Version  string
	Resource string
	Kind     string
}

func (g GroupVersionResourceKind) GVR() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    g.Group,
		Version:  g.Version,
		Resource: g.Resource,
	}
}

func (g GroupVersionResourceKind) GVK() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   g.Group,
		Version: g.Version,
		Kind:    g.Kind,
	}
}

func (g GroupVersionResourceKind) String() string {
	return fmt.Sprintf("%s %s.%s/%s", g.Kind, g.Resource, g.Group, g.Version)
}

func CRDsToGVRKs(crds []v1beta1.CustomResourceDefinition) []GroupVersionResourceKind {
	gvrks := make([]GroupVersionResourceKind, 0)
	for _, crd := range crds {
		for _, v := range crd.Spec.Versions {
			if !v.Served {
				continue
			}

			gvrk := GroupVersionResourceKind{
				Group:    crd.Spec.Group,
				Version:  v.Name,
				Resource: crd.Spec.Names.Plural,
				Kind:     crd.Spec.Names.Kind,
			}
			gvrks = append(gvrks, gvrk)
		}
	}
	return gvrks
}

func CRDToGVRK(crd v1beta1.CustomResourceDefinition) GroupVersionResourceKind {
	for _, v := range crd.Spec.Versions {
		if !v.Served {
			continue
		}

		return GroupVersionResourceKind{
			Group:    crd.Spec.Group,
			Version:  v.Name,
			Resource: crd.Spec.Names.Plural,
			Kind:     crd.Spec.Names.Kind,
		}
	}
	return GroupVersionResourceKind{}
}

type Ageable interface {
	GetCreationTimestamp() metav1.Time
}

func Age(obj Ageable) string {
	c := obj.GetCreationTimestamp()
	age := duration.HumanDuration(time.Since(c.Time))
	if c.IsZero() {
		age = "<unknown>"
	}
	return age
}

// Duck related methods

func DuckName(duck string) string {
	parts := strings.Split(duck, "=")
	return parts[0]
}

func DuckShortName(duck string) string {
	parts := strings.Split(DuckName(duck), "/")
	return parts[1]
}

// TODO: it would be neat if there was a file locally that could add more.
// Like kubectl duck learn <ducktype>
//
var knownDucks = []string{
	"duck.knative.dev/source=true",
	"messaging.knative.dev/subscribable=true",
	"duck.knative.dev/addressable=true",
	"duck.knative.dev/podspecable=true",
}

func KnownDucks() []string {
	return knownDucks
}

func ToKnownDuck(query string) (string, error) {
	duckMap := make(map[string]string, 0)
	for _, kd := range knownDucks {
		duckMap[kd] = kd
		duckMap[DuckName(kd)] = kd
		duckMap[DuckShortName(kd)] = kd
	}
	if duck, known := duckMap[query]; known {
		return duck, nil
	}
	return "", fmt.Errorf("unknown ducktype: %s", query)
}
