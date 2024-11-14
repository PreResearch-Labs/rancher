package v3

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type GPU struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GPUSpec   `json:"spec"`
	Status GPUStatus `json:"status"`
	// 这里的 json:"gpuCount" 定义 api 字段名；如图："定义返回的数据结构" data 中的 gpuCount
	GPUCount int `json:"gpuCount" norman:"default=0"`
}

type GPUSpec struct {
	// norman:"nullable" norman 规范
	Enabled bool `json:"enabled" norman:"nullable"`
}

type GPUStatus struct {
	Message string `json:"message"`
}

// action 入参结构体
type GpuCountActionInput struct {
	Count bool `json:"count"`
	Debug bool `json:"debug"`
}
