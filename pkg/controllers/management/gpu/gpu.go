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
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

type Controller struct {
	ctx            context.Context
	clusterManager *clustermanager.Manager
	clusters       mgmtv3.ClusterInterface
	clusterLister  mgmtv3.ClusterLister
	gpus           mgmtv3.GPUInterface
	gpuLister      mgmtv3.GPULister
	nodeLister     mgmtv3.NodeLister
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
		nodeLister:     management.Management.Nodes("").Controller().Lister(),
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
	// 获取所有节点
	nodes, err := c.nodeLister.List("", labels.Everything())
	if err != nil {
		return err
	}

	// 统计每个集群的 GPU 信息
	clusterGPUInfoMap := make(map[string]*v3.ClusterGPUInfo)
	var totalGPUCount int
	for _, node := range nodes {
		clusterName := node.ObjClusterName()
		if _, ok := clusterGPUInfoMap[clusterName]; !ok {
			clusterGPUInfoMap[clusterName] = &v3.ClusterGPUInfo{
				ClusterName:   clusterName,
				TotalGPUCount: 0,
				NodeGPUInfo:   []v3.NodeGPUInfo{},
			}
		}

		nodeGPUInfo := v3.NodeGPUInfo{
			NodeId:       node.Name,
			NodeHostName: node.Spec.RequestedHostname,
			NodeName:     node.Name,
			TotalGPU:     getGPUCountFromNode(node),
			UsedGPU:      getUsedGPUCountFromNode(node),
			UnusedGPU:    getUnusedGPUCountFromNode(node),
		}
		clusterGPUInfoMap[clusterName].NodeGPUInfo = append(clusterGPUInfoMap[clusterName].NodeGPUInfo, nodeGPUInfo)
		clusterGPUInfoMap[clusterName].TotalGPUCount += nodeGPUInfo.TotalGPU
		totalGPUCount += nodeGPUInfo.TotalGPU
	}

	// 构建 ClusterGPUInfo 列表
	var clusterGPUInfoList []v3.ClusterGPUInfo
	for _, clusterGPUInfo := range clusterGPUInfoMap {
		clusterGPUInfoList = append(clusterGPUInfoList, *clusterGPUInfo)
	}

	payload := &v3.GPU{
		TypeMeta: v1.TypeMeta{
			Kind: "Gpu",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      "gpu-global1",
			Namespace: "cattle-global-data",
		},
		Spec: v3.GPUSpec{
			Enabled: false,
		},
		Status: v3.GPUStatus{
			Message:        "",
			TotalGPUCount:  totalGPUCount,
			ClusterGPUInfo: clusterGPUInfoList,
		},
	}

	_, err = c.gpus.GetNamespaced("cattle-global-data", "gpu-global1", v1.GetOptions{})
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
		if !reflect.DeepEqual(payload.Status, exist.Status) {
			exist.Spec.Enabled = payload.Spec.Enabled
			exist.Status = payload.Status
			_, err = c.gpus.Update(exist)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getGPUCountFromNode(node *v3.Node) int {
	// 从节点中获取 GPU 数量
	gpuQuantity, ok := node.Status.InternalNodeStatus.Capacity["nvidia.com/gpu"]
	if !ok {
		return 0
	}
	gpuCount := gpuQuantity.Value()
	return int(gpuCount)
}

func getUsedGPUCountFromNode(node *v3.Node) int {
	// 从节点中获取已使用的 GPU 数量
	gpuQuantity, ok := node.Status.InternalNodeStatus.Allocatable["nvidia.com/gpu"]
	if !ok {
		return 0
	}
	gpuCount := gpuQuantity.Value()
	return int(gpuCount)
}

func getUnusedGPUCountFromNode(node *v3.Node) int {
	// 从节点中获取未使用的 GPU 数量
	totalGPU := getGPUCountFromNode(node)
	usedGPU := getUsedGPUCountFromNode(node)
	return totalGPU - usedGPU
}
