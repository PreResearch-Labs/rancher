package faa

import (
	"encoding/json"
	"fmt"

	"github.com/rancher/norman/types"
	gaccess "github.com/rancher/rancher/pkg/api/norman/customization/globalnamespaceaccess"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	mgmtv3 "github.com/rancher/rancher/pkg/generated/norman/management.cattle.io/v3"
)

type FaaWrapper struct {
	Users     mgmtv3.UserInterface
	GrbLister mgmtv3.GlobalRoleBindingLister
	GrLister  mgmtv3.GlobalRoleLister
}

// 给 rancher api 添加事件，如果这里没有添加 echo1 按钮是不能用的
// echo1 要与 schema 中的 MustImportAndCustomize 中 types.Action -> input ->  key 对应
func (w *FaaWrapper) Formatter(request *types.APIContext, resource *types.RawResource) {
	resource.AddAction(request, "echo1")
	// 在数据 data 中的 links 里面添加一个新的 url； 如图 4-1
	resource.Links["aaa"] = "bbb"
}

// 用户使用 action 按钮功能后触发这个函数。本例中的 echo1 按钮
func (w *FaaWrapper) ActionHandler(actionName string, action *types.Action, request *types.APIContext) error {
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
		return fmt.Errorf("aaaaa")
	}

	switch actionName {
	case "echo1":
		bytes, _ := json.Marshal(v3.FaaStatus{Msg1: "hello11"})
		request.Response.Write(bytes)
		return nil
	default:
		return fmt.Errorf("bbbbbb")
	}
}
