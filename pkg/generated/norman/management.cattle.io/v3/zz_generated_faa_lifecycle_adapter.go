package v3

import (
	"github.com/rancher/norman/lifecycle"
	"github.com/rancher/norman/resource"
	"github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/runtime"
)

type FaaLifecycle interface {
	Create(obj *v3.Faa) (runtime.Object, error)
	Remove(obj *v3.Faa) (runtime.Object, error)
	Updated(obj *v3.Faa) (runtime.Object, error)
}

type faaLifecycleAdapter struct {
	lifecycle FaaLifecycle
}

func (w *faaLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *faaLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *faaLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v3.Faa))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *faaLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v3.Faa))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *faaLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v3.Faa))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewFaaLifecycleAdapter(name string, clusterScoped bool, client FaaInterface, l FaaLifecycle) FaaHandlerFunc {
	if clusterScoped {
		resource.PutClusterScoped(FaaGroupVersionResource)
	}
	adapter := &faaLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *v3.Faa) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}