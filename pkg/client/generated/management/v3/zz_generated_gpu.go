package client

import (
	"github.com/rancher/norman/types"
)

const (
	GPUType                      = "gpu"
	GPUFieldAnnotations          = "annotations"
	GPUFieldCreated              = "created"
	GPUFieldCreatorID            = "creatorId"
	GPUFieldLabels               = "labels"
	GPUFieldName                 = "name"
	GPUFieldNodeGPUs             = "nodeGPUs"
	GPUFieldOwnerReferences      = "ownerReferences"
	GPUFieldRemoved              = "removed"
	GPUFieldState                = "state"
	GPUFieldStatus               = "status"
	GPUFieldTransitioning        = "transitioning"
	GPUFieldTransitioningMessage = "transitioningMessage"
	GPUFieldUUID                 = "uuid"
)

type GPU struct {
	types.Resource
	Annotations          map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created              string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID            string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Labels               map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	NodeGPUs             map[string]int64  `json:"nodeGPUs,omitempty" yaml:"nodeGPUs,omitempty"`
	OwnerReferences      []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed              string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	State                string            `json:"state,omitempty" yaml:"state,omitempty"`
	Status               *GPUStatus        `json:"status,omitempty" yaml:"status,omitempty"`
	Transitioning        string            `json:"transitioning,omitempty" yaml:"transitioning,omitempty"`
	TransitioningMessage string            `json:"transitioningMessage,omitempty" yaml:"transitioningMessage,omitempty"`
	UUID                 string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type GPUCollection struct {
	types.Collection
	Data   []GPU `json:"data,omitempty"`
	client *GPUClient
}

type GPUClient struct {
	apiClient *Client
}

type GPUOperations interface {
	List(opts *types.ListOpts) (*GPUCollection, error)
	ListAll(opts *types.ListOpts) (*GPUCollection, error)
	Create(opts *GPU) (*GPU, error)
	Update(existing *GPU, updates interface{}) (*GPU, error)
	Replace(existing *GPU) (*GPU, error)
	ByID(id string) (*GPU, error)
	Delete(container *GPU) error

	ActionUpdateGPU(resource *GPU, input *GpuActionInput) error
}

func newGPUClient(apiClient *Client) *GPUClient {
	return &GPUClient{
		apiClient: apiClient,
	}
}

func (c *GPUClient) Create(container *GPU) (*GPU, error) {
	resp := &GPU{}
	err := c.apiClient.Ops.DoCreate(GPUType, container, resp)
	return resp, err
}

func (c *GPUClient) Update(existing *GPU, updates interface{}) (*GPU, error) {
	resp := &GPU{}
	err := c.apiClient.Ops.DoUpdate(GPUType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *GPUClient) Replace(obj *GPU) (*GPU, error) {
	resp := &GPU{}
	err := c.apiClient.Ops.DoReplace(GPUType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *GPUClient) List(opts *types.ListOpts) (*GPUCollection, error) {
	resp := &GPUCollection{}
	err := c.apiClient.Ops.DoList(GPUType, opts, resp)
	resp.client = c
	return resp, err
}

func (c *GPUClient) ListAll(opts *types.ListOpts) (*GPUCollection, error) {
	resp := &GPUCollection{}
	resp, err := c.List(opts)
	if err != nil {
		return resp, err
	}
	data := resp.Data
	for next, err := resp.Next(); next != nil && err == nil; next, err = next.Next() {
		data = append(data, next.Data...)
		resp = next
		resp.Data = data
	}
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (cc *GPUCollection) Next() (*GPUCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &GPUCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *GPUClient) ByID(id string) (*GPU, error) {
	resp := &GPU{}
	err := c.apiClient.Ops.DoByID(GPUType, id, resp)
	return resp, err
}

func (c *GPUClient) Delete(container *GPU) error {
	return c.apiClient.Ops.DoResourceDelete(GPUType, &container.Resource)
}

func (c *GPUClient) ActionUpdateGPU(resource *GPU, input *GpuActionInput) error {
	err := c.apiClient.Ops.DoAction(GPUType, "updateGPU", &resource.Resource, input, nil)
	return err
}
