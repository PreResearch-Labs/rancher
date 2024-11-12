package v3

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Faa struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec1   FaaSpec   `json:"spec1"`
	Status1 FaaStatus `json:"status1"`
	// 这里的 json:"asf1" 定义 api 字段名；如图："定义返回的数据结构" data 中的 asf1
	Asf string `json:"asf1" norman:"default=bbbb"`
}

type FaaSpec struct {
	// norman:"nullable" norman 规范
	Bar1 bool `json:"bar1" norman:"nullable"`
}

type FaaStatus struct {
	Msg1 string `json:"msg1"`
}

// action 入参结构体
type EchoActionInput struct {
	Echo bool `json:"echo"`
	Aaaa bool `json:"aaaa"`
}
