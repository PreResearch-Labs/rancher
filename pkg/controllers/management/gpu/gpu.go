package gpu

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
	gpus           mgmtv3.GPUInterface
	gpuLister      mgmtv3.GPULister
}

// controller 初始化
func Register(ctx context.Context, management *config.ManagementContext, manager *clustermanager.Manager) error {
	c := &Controller{
		ctx:            ctx,
		clusterManager: manager,
		clusters:       management.Management.Clusters(""),
		clusterLister:  management.Management.Clusters("").Controller().Lister(),
		gpus:           management.Management.GPUs(""),
		gpuLister:      management.Management.GPUs("").Controller().Lister(),
	}

	management.Management.Clusters("").AddHandler(ctx, "gpu-controller", c.syncByCluster)
	return nil
}

func (c *Controller) syncByCluster(key string, cluster *v3.Cluster) (runtime.Object, error) {
	if cluster == nil || cluster.DeletionTimestamp != nil {
		return nil, nil
	}
	return nil, c.InitGPU()
}

func (c *Controller) InitGPU() error {
	payload := &v3.GPU{
		TypeMeta: v1.TypeMeta{
			Kind: "Gpu",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "gpu-global1",
			Namespace: "cattle-global-data",
		},
		Spec: v3.GPUSpec{
			Enabled: true,
		},
		Status: v3.GPUStatus{
			Message: "",
		},
	}
	_, err := c.gpus.GetNamespaced("cattle-global-data", "gpu-global1", v1.GetOptions{})
	if err != nil && !k8sErrors.IsNotFound(err) {
		return err
	}
	if k8sErrors.IsNotFound(err) {
		if _, err := c.gpus.Create(payload); err != nil {
			return err
		}
	} else {
		exist, err := c.gpus.GetNamespaced("cattle-global-data", "gpu-global1", v1.GetOptions{})
		if err != nil {
			return err
		}
		if !reflect.DeepEqual(payload.Status.Message, exist.Status.Message) {
			exist.Spec.Enabled = payload.Spec.Enabled
			exist.Status.Message = payload.Status.Message
			_, err = c.gpus.Update(exist)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
