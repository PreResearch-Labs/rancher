package v3

import (
	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type GPULifecycle interface {
	Create(obj *v3.GPU) (runtime.Object, error)
	Remove(obj *v3.GPU) (runtime.Object, error)
	Updated(obj *v3.GPU) (runtime.Object, error)
}

type gpuLifecycleAdapter struct {
	lifecycle GPULifecycle
}

func (w *gpuLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *gpuLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *gpuLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v3.GPU))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *gpuLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v3.GPU))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *gpuLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v3.GPU))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGPULifecycleAdapter(name string, clusterScoped bool, client GPUInterface, l GPULifecycle) GPUHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(GPUGroupVersionResource)
	}
	adapter := &gpuLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.GPU) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
