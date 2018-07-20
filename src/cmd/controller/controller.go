/*
Copyright 2017 The Kubernetes Authors.

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

package main

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	clientset "github.com/fairfaxmedia/annotation-controller/src/pkg/client/clientset/versioned"
	annotationscheme "github.com/fairfaxmedia/annotation-controller/src/pkg/client/clientset/versioned/scheme"
	informers "github.com/fairfaxmedia/annotation-controller/src/pkg/client/informers/externalversions"
	listers "github.com/fairfaxmedia/annotation-controller/src/pkg/client/listers/example.com/v1"
)

const controllerAgentName = "annotation-controller"

const (
	SuccessSynced         = "Synced"
	ErrResourceExists     = "ErrResourceExists"
	MessageResourceExists = "Resource %q already exists and is not managed by Annotation"
	MessageResourceSynced = "Annotation synced successfully"
)

type T struct {
	client kubernetes.Interface
}

type Controller struct {
	kubeclientset       kubernetes.Interface
	annotationclientset clientset.Interface
	annotationsLister   listers.AnnotationLister
	annotationsSynced   cache.InformerSynced
	workqueue           workqueue.RateLimitingInterface
	recorder            record.EventRecorder
}

func NewController(
	kubeclientset kubernetes.Interface,
	annotationclientset clientset.Interface,
	kubeInformerFactory kubeinformers.SharedInformerFactory,
	annotationInformerFactory informers.SharedInformerFactory) *Controller {

	annotationInformer := annotationInformerFactory.Example().V1().Annotations()

	annotationscheme.AddToScheme(scheme.Scheme)
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:       kubeclientset,
		annotationclientset: annotationclientset,
		annotationsLister:   annotationInformer.Lister(),
		annotationsSynced:   annotationInformer.Informer().HasSynced,
		workqueue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Annotations"),
		recorder:            recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when Annotation resources change
	annotationInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueAnnotation,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueAnnotation(new)
		},
	})

	return controller
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	glog.Info("Starting Annotation controller")

	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.annotationsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		if key, ok = obj.(string); !ok {
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	annotation, err := c.annotationsLister.Annotations(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("annotation '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	var t T
	t.client = c.kubeclientset

	for _, target := range annotation.Spec.Targets {

		params := []reflect.Value{reflect.ValueOf(annotation.GetNamespace()), reflect.ValueOf(target.Data)}
		reflect.ValueOf(&t).MethodByName(strings.Title(target.Kind)).Call(params)

	}

	c.recorder.Event(annotation, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) enqueueAnnotation(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		runtime.HandleError(err)
		return
	}
	c.workqueue.AddRateLimited(key)
}

func (t *T) Namespace(ns string, data map[string]string) {
	obj, err := t.client.CoreV1().Namespaces().Get(ns, meta_v1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	}
	if !annotationsExist(obj.ObjectMeta) {
		obj.ObjectMeta.Annotations = make(map[string]string)
	}
	obj.ObjectMeta.Annotations = mergeAnnotations(obj.ObjectMeta.Annotations, data)
	obj, err = t.client.CoreV1().Namespaces().Update(obj)
	if err != nil {
		fmt.Println(err)
	}
}

func mergeAnnotations(old map[string]string, new map[string]string) map[string]string {
	for k, v := range new {
		_, present := old[k]
		if !present {
			old[k] = v
		}
	}
	return old
}

func annotationsExist(meta meta_v1.ObjectMeta) bool {
	if meta.Annotations == nil {
		return false
	}
	return true
}
