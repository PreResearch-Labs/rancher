package faa

import (
	"context"
	"reflect"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/rancher/pkg/clustermanager"
	mgmtv3 "github.com/rancher/rancher/pkg/generated/norman/management.cattle.io/v3"
	"github.com/rancher/rancher/pkg/types/config"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Controller struct {
	ctx            context.Context
	clusterManager *clustermanager.Manager
	clusters       mgmtv3.ClusterInterface
	clusterLister  mgmtv3.ClusterLister
	faas           mgmtv3.FaaInterface
	faaLister      mgmtv3.FaaLister
}

// controller 初始化
func Register(ctx context.Context, scaledContext *config.ScaledContext, manager *clustermanager.Manager) error {
	c := &Controller{
		ctx:            ctx,
		clusterManager: manager,
		clusters:       scaledContext.Management.Clusters(""),
		clusterLister:  scaledContext.Management.Clusters("").Controller().Lister(),
		faas:           scaledContext.Management.Faas(""),
		faaLister:      scaledContext.Management.Faas("").Controller().Lister(),
	}

	scaledContext.Management.Clusters("").AddHandler(ctx, "faa-controller", c.syncByCluster)
	return nil
}

func (c *Controller) syncByCluster(key string, cluster *v3.Cluster) (runtime.Object, error) {
	if cluster == nil || cluster.DeletionTimestamp != nil {
		return nil, nil
	}
	return nil, c.InitFaa()
}

func (c *Controller) InitFaa() error {
	payload := &v3.Faa{
		TypeMeta: v1.TypeMeta{
			Kind: "Faa",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "faa-global1",
			Namespace: "cattle-global-data",
		},
		Spec1: v3.FaaSpec{
			Bar1: false,
		},
		Status1: v3.FaaStatus{
			Msg1: "",
		},
	}
	_, err := c.faas.GetNamespaced("cattle-global-data", "faa-global1", v1.GetOptions{})
	if err != nil && !k8sErrors.IsNotFound(err) {
		return err
	}
	if k8sErrors.IsNotFound(err) {
		if _, err := c.faas.Create(payload); err != nil {
			return err
		}
	} else {
		exist, err := c.faas.GetNamespaced("cattle-global-data", "faa-global1", v1.GetOptions{})
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(payload.Status1.Msg1, exist.Status1.Msg1) {
			exist.Spec1.Bar1 = payload.Spec1.Bar1
			exist.Status1.Msg1 = payload.Status1.Msg1
			_, err = c.faas.Update(exist)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
