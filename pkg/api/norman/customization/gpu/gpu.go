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

func (w *GPUWrapper) Formatter(request *types.APIContext, resource *types.RawResource) {
	resource.AddAction(request, "gpustats")
	resource.Links["gpuStats"] = "gpuStats"
}

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
		return fmt.Errorf("access denied")
	}

	switch actionName {
	case "gpustats":
		// Mock data for GPU stats
		gpuStats := v3.GPUStatus{
			TotalGPUs: 4,
		}
		bytes, _ := json.Marshal(gpuStats)
		request.Response.Write(bytes)
		return nil
	default:
		return fmt.Errorf("unknown action")
	}
}
