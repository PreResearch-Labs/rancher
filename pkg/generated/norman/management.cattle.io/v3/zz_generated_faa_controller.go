package v3

import (
	"context"
	"time"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	FaaGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Faa",
	}
	FaaResource = metav1.APIResource{
		Name:         "faas",
		SingularName: "faa",
		Namespaced:   false,
		Kind:         FaaGroupVersionKind.Kind,
	}

	FaaGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "faas",
	}
)

func init() {
	resource.Put(FaaGroupVersionResource)
}

// Deprecated: use v3.Faa instead
type Faa = v3.Faa

func NewFaa(namespace, name string, obj v3.Faa) *v3.Faa {
	obj.APIVersion, obj.Kind = FaaGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type FaaHandlerFunc func(key string, obj *v3.Faa) (runtime.Object, error)

type FaaChangeHandlerFunc func(obj *v3.Faa) (runtime.Object, error)

type FaaLister interface {
	List(namespace string, selector labels.Selector) (ret []*v3.Faa, err error)
	Get(namespace, name string) (*v3.Faa, error)
}

type FaaController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() FaaLister
	AddHandler(ctx context.Context, name string, handler FaaHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync FaaHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler FaaHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler FaaHandlerFunc)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, after time.Duration)
}

type FaaInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v3.Faa) (*v3.Faa, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v3.Faa, error)
	Get(name string, opts metav1.GetOptions) (*v3.Faa, error)
	Update(*v3.Faa) (*v3.Faa, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*v3.FaaList, error)
	ListNamespaced(namespace string, opts metav1.ListOptions) (*v3.FaaList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() FaaController
	AddHandler(ctx context.Context, name string, sync FaaHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync FaaHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle FaaLifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle FaaLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync FaaHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync FaaHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle FaaLifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle FaaLifecycle)
}

type faaLister struct {
	ns         string
	controller *faaController
}

func (l *faaLister) List(namespace string, selector labels.Selector) (ret []*v3.Faa, err error) {
	if namespace == "" {
		namespace = l.ns
	}
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v3.Faa))
	})
	return
}

func (l *faaLister) Get(namespace, name string) (*v3.Faa, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    FaaGroupVersionKind.Group,
			Resource: FaaGroupVersionResource.Resource,
		}, key)
	}
	return obj.(*v3.Faa), nil
}

type faaController struct {
	ns string
	controller.GenericController
}

func (c *faaController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *faaController) Lister() FaaLister {
	return &faaLister{
		ns:         c.ns,
		controller: c,
	}
}

func (c *faaController) AddHandler(ctx context.Context, name string, handler FaaHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.Faa); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *faaController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler FaaHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.Faa); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *faaController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler FaaHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.Faa); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *faaController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler FaaHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.Faa); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type faaFactory struct {
}

func (c faaFactory) Object() runtime.Object {
	return &v3.Faa{}
}

func (c faaFactory) List() runtime.Object {
	return &v3.FaaList{}
}

func (s *faaClient) Controller() FaaController {
	genericController := controller.NewGenericController(s.ns, FaaGroupVersionKind.Kind+"Controller",
		s.client.controllerFactory.ForResourceKind(FaaGroupVersionResource, FaaGroupVersionKind.Kind, false))

	return &faaController{
		ns:                s.ns,
		GenericController: genericController,
	}
}

type faaClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   FaaController
}

func (s *faaClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *faaClient) Create(o *v3.Faa) (*v3.Faa, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v3.Faa), err
}

func (s *faaClient) Get(name string, opts metav1.GetOptions) (*v3.Faa, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v3.Faa), err
}

func (s *faaClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v3.Faa, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v3.Faa), err
}

func (s *faaClient) Update(o *v3.Faa) (*v3.Faa, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v3.Faa), err
}

func (s *faaClient) UpdateStatus(o *v3.Faa) (*v3.Faa, error) {
	obj, err := s.objectClient.UpdateStatus(o.Name, o)
	return obj.(*v3.Faa), err
}

func (s *faaClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *faaClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *faaClient) List(opts metav1.ListOptions) (*v3.FaaList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*v3.FaaList), err
}

func (s *faaClient) ListNamespaced(namespace string, opts metav1.ListOptions) (*v3.FaaList, error) {
	obj, err := s.objectClient.ListNamespaced(namespace, opts)
	return obj.(*v3.FaaList), err
}

func (s *faaClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *faaClient) Patch(o *v3.Faa, patchType types.PatchType, data []byte, subresources ...string) (*v3.Faa, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v3.Faa), err
}

func (s *faaClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *faaClient) AddHandler(ctx context.Context, name string, sync FaaHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *faaClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync FaaHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *faaClient) AddLifecycle(ctx context.Context, name string, lifecycle FaaLifecycle) {
	sync := NewFaaLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *faaClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle FaaLifecycle) {
	sync := NewFaaLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *faaClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync FaaHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *faaClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync FaaHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *faaClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle FaaLifecycle) {
	sync := NewFaaLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *faaClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle FaaLifecycle) {
	sync := NewFaaLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}
