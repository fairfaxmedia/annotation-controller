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
package v1

import (
	v1 "github.com/fairfaxmedia/annotation-controller/src/pkg/apis/example.com/v1"
	scheme "github.com/fairfaxmedia/annotation-controller/src/pkg/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AnnotationsGetter has a method to return a AnnotationInterface.
// A group's client should implement this interface.
type AnnotationsGetter interface {
	Annotations(namespace string) AnnotationInterface
}

// AnnotationInterface has methods to work with Annotation resources.
type AnnotationInterface interface {
	Create(*v1.Annotation) (*v1.Annotation, error)
	Update(*v1.Annotation) (*v1.Annotation, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.Annotation, error)
	List(opts meta_v1.ListOptions) (*v1.AnnotationList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Annotation, err error)
	AnnotationExpansion
}

// annotations implements AnnotationInterface
type annotations struct {
	client rest.Interface
	ns     string
}

// newAnnotations returns a Annotations
func newAnnotations(c *ExampleV1Client, namespace string) *annotations {
	return &annotations{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the annotation, and returns the corresponding annotation object, and an error if there is any.
func (c *annotations) Get(name string, options meta_v1.GetOptions) (result *v1.Annotation, err error) {
	result = &v1.Annotation{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("annotations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Annotations that match those selectors.
func (c *annotations) List(opts meta_v1.ListOptions) (result *v1.AnnotationList, err error) {
	result = &v1.AnnotationList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("annotations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested annotations.
func (c *annotations) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("annotations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a annotation and creates it.  Returns the server's representation of the annotation, and an error, if there is any.
func (c *annotations) Create(annotation *v1.Annotation) (result *v1.Annotation, err error) {
	result = &v1.Annotation{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("annotations").
		Body(annotation).
		Do().
		Into(result)
	return
}

// Update takes the representation of a annotation and updates it. Returns the server's representation of the annotation, and an error, if there is any.
func (c *annotations) Update(annotation *v1.Annotation) (result *v1.Annotation, err error) {
	result = &v1.Annotation{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("annotations").
		Name(annotation.Name).
		Body(annotation).
		Do().
		Into(result)
	return
}

// Delete takes name of the annotation and deletes it. Returns an error if one occurs.
func (c *annotations) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("annotations").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *annotations) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("annotations").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched annotation.
func (c *annotations) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Annotation, err error) {
	result = &v1.Annotation{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("annotations").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
