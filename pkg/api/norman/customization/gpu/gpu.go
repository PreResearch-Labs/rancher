package gpu

import (
	"encoding/json"
	"fmt"

	"github.com/rancher/norman/types"
	gaccess "github.com/rancher/rancher/pkg/api/norman/customization/globalnamespaceaccess"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	mgmtv3 "github.com/rancher/rancher/pkg/generated/norman/management.cattle.io/v3"
)

type GPUWrapper struct {
	Users     mgmtv3.UserInterface
	GrbLister mgmtv3.GlobalRoleBindingLister
	GrLister  mgmtv3.GlobalRoleLister
}

// 给 rancher api 添加事件，如果这里没有添加 countGPU 按钮是不能用的
// countGPU 要与 schema 中的 MustImportAndCustomize 中 types.Action -> input ->  key 对应
func (w *GPUWrapper) Formatter(request *types.APIContext, resource *types.RawResource) {
	resource.AddAction(request, "countGPU1")
	// 在数据 data 中的 links 里面添加一个新的 url；
	resource.Links["gpuStats"] = "gpuStats"
}

// 用户使用 action 按钮功能后触发这个函数。本例中的 countGPU 按钮
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

		// 模拟节点和 GPU 数据
		nodeGPUInfo := []v3.NodeGPUInfo{
			{NodeName: "node1", TotalGPU: 8, UsedGPU: 3, UnusedGPU: 5},
			{NodeName: "node2", TotalGPU: 6, UsedGPU: 2, UnusedGPU: 4},
			{NodeName: "node3", TotalGPU: 4, UsedGPU: 1, UnusedGPU: 3},
		}

		// 根据输入参数筛选节点
		var filteredNodeGPUInfo []v3.NodeGPUInfo
		for _, node := range nodeGPUInfo {
			if input.NodeName == "" || node.NodeName == input.NodeName {
				filteredNodeGPUInfo = append(filteredNodeGPUInfo, node)
			}
		}

		// 计算总 GPU 数量
		var totalGPUCount int
		for _, node := range filteredNodeGPUInfo {
			totalGPUCount += node.TotalGPU
		}

		// 构建响应
		response := v3.GPU{
			TotalGPUCount: totalGPUCount,
			NodeGPUInfo:   filteredNodeGPUInfo,
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
