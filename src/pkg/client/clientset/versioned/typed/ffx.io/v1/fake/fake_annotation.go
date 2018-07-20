/*
Copyright 2018 pickledrick

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
package fake

import (
	ffx_io_v1 "github.com/fairfaxmedia/annotation-controller/src/pkg/apis/ffx.io/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAnnotations implements AnnotationInterface
type FakeAnnotations struct {
	Fake *FakeExampleV1
	ns   string
}

var annotationsResource = schema.GroupVersionResource{Group: "ffx.io", Version: "v1", Resource: "annotations"}

var annotationsKind = schema.GroupVersionKind{Group: "ffx.io", Version: "v1", Kind: "Annotation"}

// Get takes name of the annotation, and returns the corresponding annotation object, and an error if there is any.
func (c *FakeAnnotations) Get(name string, options v1.GetOptions) (result *ffx_io_v1.Annotation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(annotationsResource, c.ns, name), &ffx_io_v1.Annotation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*ffx_io_v1.Annotation), err
}

// List takes label and field selectors, and returns the list of Annotations that match those selectors.
func (c *FakeAnnotations) List(opts v1.ListOptions) (result *ffx_io_v1.AnnotationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(annotationsResource, annotationsKind, c.ns, opts), &ffx_io_v1.AnnotationList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &ffx_io_v1.AnnotationList{}
	for _, item := range obj.(*ffx_io_v1.AnnotationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested annotations.
func (c *FakeAnnotations) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(annotationsResource, c.ns, opts))

}

// Create takes the representation of a annotation and creates it.  Returns the server's representation of the annotation, and an error, if there is any.
func (c *FakeAnnotations) Create(annotation *ffx_io_v1.Annotation) (result *ffx_io_v1.Annotation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(annotationsResource, c.ns, annotation), &ffx_io_v1.Annotation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*ffx_io_v1.Annotation), err
}

// Update takes the representation of a annotation and updates it. Returns the server's representation of the annotation, and an error, if there is any.
func (c *FakeAnnotations) Update(annotation *ffx_io_v1.Annotation) (result *ffx_io_v1.Annotation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(annotationsResource, c.ns, annotation), &ffx_io_v1.Annotation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*ffx_io_v1.Annotation), err
}

// Delete takes name of the annotation and deletes it. Returns an error if one occurs.
func (c *FakeAnnotations) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(annotationsResource, c.ns, name), &ffx_io_v1.Annotation{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAnnotations) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(annotationsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &ffx_io_v1.AnnotationList{})
	return err
}

// Patch applies the patch and returns the patched annotation.
func (c *FakeAnnotations) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *ffx_io_v1.Annotation, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(annotationsResource, c.ns, name, data, subresources...), &ffx_io_v1.Annotation{})

	if obj == nil {
		return nil, err
	}
	return obj.(*ffx_io_v1.Annotation), err
}
