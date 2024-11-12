package client

import (
	"github.com/rancher/norman/types"
)

const (
	FaaType                 = "faa"
	FaaFieldAnnotations     = "annotations"
	FaaFieldAsf             = "asf1"
	FaaFieldCreated         = "created"
	FaaFieldCreatorID       = "creatorId"
	FaaFieldLabels          = "labels"
	FaaFieldName            = "name"
	FaaFieldOwnerReferences = "ownerReferences"
	FaaFieldRemoved         = "removed"
	FaaFieldSpec1           = "spec1"
	FaaFieldStatus1         = "status1"
	FaaFieldUUID            = "uuid"
)

type Faa struct {
	types.Resource
	Annotations     map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Asf             string            `json:"asf1,omitempty" yaml:"asf1,omitempty"`
	Created         string            `json:"created,omitempty" yaml:"created,omitempty"`
	CreatorID       string            `json:"creatorId,omitempty" yaml:"creatorId,omitempty"`
	Labels          map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Name            string            `json:"name,omitempty" yaml:"name,omitempty"`
	OwnerReferences []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed         string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	Spec1           *FaaSpec          `json:"spec1,omitempty" yaml:"spec1,omitempty"`
	Status1         *FaaStatus        `json:"status1,omitempty" yaml:"status1,omitempty"`
	UUID            string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type FaaCollection struct {
	types.Collection
	Data   []Faa `json:"data,omitempty"`
	client *FaaClient
}

type FaaClient struct {
	apiClient *Client
}

type FaaOperations interface {
	List(opts *types.ListOpts) (*FaaCollection, error)
	ListAll(opts *types.ListOpts) (*FaaCollection, error)
	Create(opts *Faa) (*Faa, error)
	Update(existing *Faa, updates interface{}) (*Faa, error)
	Replace(existing *Faa) (*Faa, error)
	ByID(id string) (*Faa, error)
	Delete(container *Faa) error

	ActionEcho1(resource *Faa, input *EchoActionInput) error
}

func newFaaClient(apiClient *Client) *FaaClient {
	return &FaaClient{
		apiClient: apiClient,
	}
}

func (c *FaaClient) Create(container *Faa) (*Faa, error) {
	resp := &Faa{}
	err := c.apiClient.Ops.DoCreate(FaaType, container, resp)
	return resp, err
}

func (c *FaaClient) Update(existing *Faa, updates interface{}) (*Faa, error) {
	resp := &Faa{}
	err := c.apiClient.Ops.DoUpdate(FaaType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *FaaClient) Replace(obj *Faa) (*Faa, error) {
	resp := &Faa{}
	err := c.apiClient.Ops.DoReplace(FaaType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *FaaClient) List(opts *types.ListOpts) (*FaaCollection, error) {
	resp := &FaaCollection{}
	err := c.apiClient.Ops.DoList(FaaType, opts, resp)
	resp.client = c
	return resp, err
}

func (c *FaaClient) ListAll(opts *types.ListOpts) (*FaaCollection, error) {
	resp := &FaaCollection{}
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

func (cc *FaaCollection) Next() (*FaaCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &FaaCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *FaaClient) ByID(id string) (*Faa, error) {
	resp := &Faa{}
	err := c.apiClient.Ops.DoByID(FaaType, id, resp)
	return resp, err
}

func (c *FaaClient) Delete(container *Faa) error {
	return c.apiClient.Ops.DoResourceDelete(FaaType, &container.Resource)
}

func (c *FaaClient) ActionEcho1(resource *Faa, input *EchoActionInput) error {
	err := c.apiClient.Ops.DoAction(FaaType, "echo1", &resource.Resource, input, nil)
	return err
}
