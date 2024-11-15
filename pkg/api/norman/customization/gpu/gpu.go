package gpu

import (
	"encoding/json"
	"fmt"

	"github.com/rancher/norman/types"
	gaccess "github.com/rancher/rancher/pkg/api/norman/customization/globalnamespaceaccess"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	mgmtv3 "github.com/rancher/rancher/pkg/generated/norman/management.cattle.io/v3"
	"k8s.io/apimachinery/pkg/labels"
)

type GPUWrapper struct {
	Users      mgmtv3.UserInterface
	GrbLister  mgmtv3.GlobalRoleBindingLister
	GrLister   mgmtv3.GlobalRoleLister
	NodeLister mgmtv3.NodeLister
}

// 给 rancher api 添加事件，如果这里没有添加 countGPU1 按钮是不能用的
// countGPU1 要与 schema 中的 MustImportAndCustomize 中 types.Action -> input ->  key 对应
func (w *GPUWrapper) Formatter(request *types.APIContext, resource *types.RawResource) {
	resource.AddAction(request, "countGPU1")
	// 在数据 data 中的 links 里面添加一个新的 url；
	resource.Links["gpuStats"] = "gpuStats"
}

// 用户使用 action 按钮功能后触发这个函数。本例中的 countGPU1 按钮
func (w *GPUWrapper) ActionHandler(actionName string, action *types.Action, request *types.APIContext) error {
	callerID := request.Request.Header.Get(gaccess.ImpersonateUserHeader)
	ma := gaccess.MemberAccess{
		Users:     w.Users,
		GrbLister: w.GrbLister,
		GrLister:  w.GrLister,
	}
	canAccess, err := ma.IsAdmin(callerID)
	if err != nil {
		return err
	}
	if !canAccess {
		return fmt.Errorf("GPUWrapper Access denied")
	}

	switch actionName {
	case "countGPU1":
		// 获取输入参数
		var input v3.GpuCountActionInput
		if err := json.NewDecoder(request.Request.Body).Decode(&input); err != nil {
			return err
		}

		// 检查 NodeLister 是否为 nil
		if w.NodeLister == nil {
			return fmt.Errorf("NodeLister is not initialized")
		}

		// 获取所有节点
		nodes, err := w.NodeLister.List("", labels.Everything())
		if err != nil {
			return err
		}

		// 根据输入参数筛选节点
		clusterGPUInfoMap := make(map[string]*v3.ClusterGPUInfo)
		var totalGPUCount int
		for _, node := range nodes {
			if (input.NodeHostName == "" || node.Spec.RequestedHostname == input.NodeHostName) &&
				(input.NodeName == "" || node.Name == input.NodeName) &&
				(input.NodeId == "" || node.Name == input.NodeId) &&
				(input.ClusterName == "" || node.ObjClusterName() == input.ClusterName) {
				clusterName := node.ObjClusterName()
				if _, ok := clusterGPUInfoMap[clusterName]; !ok {
					clusterGPUInfoMap[clusterName] = &v3.ClusterGPUInfo{
						ClusterName:   clusterName,
						TotalGPUCount: 0,
						NodeGPUInfo:   []v3.NodeGPUInfo{},
					}
				}

				nodeGPUInfo := v3.NodeGPUInfo{
					NodeId:       fmt.Sprintf("%s:%s", clusterName, node.Name),
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
		}

		// 构建响应
		var clusterGPUInfoList []v3.ClusterGPUInfo
		for _, clusterGPUInfo := range clusterGPUInfoMap {
			clusterGPUInfoList = append(clusterGPUInfoList, *clusterGPUInfo)
		}

		response := v3.GPU{
			Status: v3.GPUStatus{
				TotalGPUCount:  totalGPUCount,
				ClusterGPUInfo: clusterGPUInfoList,
			},
		}

		// 序列化响应并发送
		bytes, err := json.Marshal(response)
		if err != nil {
			return err
		}
		request.Response.Write(bytes)
		return nil
	default:
		return fmt.Errorf("Unknown action: %s", actionName)
	}
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
