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
	GPUGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "GPU",
	}
	GPUResource = metav1.APIResource{
		Name:         "gpus",
		SingularName: "gpu",
		Namespaced:   false,
		Kind:         GPUGroupVersionKind.Kind,
	}

	GPUGroupVersionResource = schema.GroupVersionResource{
		Group:    GroupName,
		Version:  Version,
		Resource: "gpus",
	}
)

func init() {
	resource.Put(GPUGroupVersionResource)
}

// Deprecated: use v3.GPU instead
type GPU = v3.GPU

func NewGPU(namespace, name string, obj v3.GPU) *v3.GPU {
	obj.APIVersion, obj.Kind = GPUGroupVersionKind.ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

type GPUHandlerFunc func(key string, obj *v3.GPU) (runtime.Object, error)

type GPUChangeHandlerFunc func(obj *v3.GPU) (runtime.Object, error)

type GPULister interface {
	List(namespace string, selector labels.Selector) (ret []*v3.GPU, err error)
	Get(namespace, name string) (*v3.GPU, error)
}

type GPUController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() GPULister
	AddHandler(ctx context.Context, name string, handler GPUHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GPUHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler GPUHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, handler GPUHandlerFunc)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, after time.Duration)
}

type GPUInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*v3.GPU) (*v3.GPU, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v3.GPU, error)
	Get(name string, opts metav1.GetOptions) (*v3.GPU, error)
	Update(*v3.GPU) (*v3.GPU, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*v3.GPUList, error)
	ListNamespaced(namespace string, opts metav1.ListOptions) (*v3.GPUList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() GPUController
	AddHandler(ctx context.Context, name string, sync GPUHandlerFunc)
	AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GPUHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle GPULifecycle)
	AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle GPULifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GPUHandlerFunc)
	AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync GPUHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GPULifecycle)
	AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle GPULifecycle)
}

type gpuLister struct {
	ns         string
	controller *gpuController
}

func (l *gpuLister) List(namespace string, selector labels.Selector) (ret []*v3.GPU, err error) {
	if namespace == "" {
		namespace = l.ns
	}
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*v3.GPU))
	})
	return
}

func (l *gpuLister) Get(namespace, name string) (*v3.GPU, error) {
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
			Group:    GPUGroupVersionKind.Group,
			Resource: GPUGroupVersionResource.Resource,
		}, key)
	}
	return obj.(*v3.GPU), nil
}

type gpuController struct {
	ns string
	controller.GenericController
}

func (c *gpuController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *gpuController) Lister() GPULister {
	return &gpuLister{
		ns:         c.ns,
		controller: c,
	}
}

func (c *gpuController) AddHandler(ctx context.Context, name string, handler GPUHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.GPU); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *gpuController) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, handler GPUHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.GPU); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *gpuController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler GPUHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.GPU); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *gpuController) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, cluster string, handler GPUHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if !enabled() {
			return nil, nil
		} else if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*v3.GPU); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type gpuFactory struct {
}

func (c gpuFactory) Object() runtime.Object {
	return &v3.GPU{}
}

func (c gpuFactory) List() runtime.Object {
	return &v3.GPUList{}
}

func (s *gpuClient) Controller() GPUController {
	genericController := controller.NewGenericController(s.ns, GPUGroupVersionKind.Kind+"Controller",
		s.client.controllerFactory.ForResourceKind(GPUGroupVersionResource, GPUGroupVersionKind.Kind, false))

	return &gpuController{
		ns:                s.ns,
		GenericController: genericController,
	}
}

type gpuClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   GPUController
}

func (s *gpuClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *gpuClient) Create(o *v3.GPU) (*v3.GPU, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*v3.GPU), err
}

func (s *gpuClient) Get(name string, opts metav1.GetOptions) (*v3.GPU, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*v3.GPU), err
}

func (s *gpuClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*v3.GPU, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*v3.GPU), err
}

func (s *gpuClient) Update(o *v3.GPU) (*v3.GPU, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*v3.GPU), err
}

func (s *gpuClient) UpdateStatus(o *v3.GPU) (*v3.GPU, error) {
	obj, err := s.objectClient.UpdateStatus(o.Name, o)
	return obj.(*v3.GPU), err
}

func (s *gpuClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *gpuClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *gpuClient) List(opts metav1.ListOptions) (*v3.GPUList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*v3.GPUList), err
}

func (s *gpuClient) ListNamespaced(namespace string, opts metav1.ListOptions) (*v3.GPUList, error) {
	obj, err := s.objectClient.ListNamespaced(namespace, opts)
	return obj.(*v3.GPUList), err
}

func (s *gpuClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *gpuClient) Patch(o *v3.GPU, patchType types.PatchType, data []byte, subresources ...string) (*v3.GPU, error) {
	obj, err := s.objectClient.Patch(o.Name, o, patchType, data, subresources...)
	return obj.(*v3.GPU), err
}

func (s *gpuClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *gpuClient) AddHandler(ctx context.Context, name string, sync GPUHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *gpuClient) AddFeatureHandler(ctx context.Context, enabled func() bool, name string, sync GPUHandlerFunc) {
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *gpuClient) AddLifecycle(ctx context.Context, name string, lifecycle GPULifecycle) {
	sync := NewGPULifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *gpuClient) AddFeatureLifecycle(ctx context.Context, enabled func() bool, name string, lifecycle GPULifecycle) {
	sync := NewGPULifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddFeatureHandler(ctx, enabled, name, sync)
}

func (s *gpuClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GPUHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *gpuClient) AddClusterScopedFeatureHandler(ctx context.Context, enabled func() bool, name, clusterName string, sync GPUHandlerFunc) {
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}

func (s *gpuClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GPULifecycle) {
	sync := NewGPULifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *gpuClient) AddClusterScopedFeatureLifecycle(ctx context.Context, enabled func() bool, name, clusterName string, lifecycle GPULifecycle) {
	sync := NewGPULifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedFeatureHandler(ctx, enabled, name, clusterName, sync)
}
