package rbacx

import (
	"context"
	"fmt"
	"sync"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/permissionrpc"
)

type RbacHolder struct {
	sync.RWMutex

	pr       permissionrpc.PermissionRpc
	policies map[string]IApiPolicy
}

func NewRbacHolder(pr permissionrpc.PermissionRpc) *RbacHolder {
	return &RbacHolder{
		pr:       pr,
		policies: make(map[string]IApiPolicy),
	}
}

func (m *RbacHolder) LoadPolicy() error {
	m.Lock()
	defer m.Unlock()

	in := &permissionrpc.FindApiListReq{}
	list, err := m.pr.FindApiList(context.Background(), in)
	if err != nil {
		return err
	}

	m.policies = make(map[string]IApiPolicy)
	for _, v := range list.List {
		m.addPolicy(nil, v)
	}

	return nil
}

func (m *RbacHolder) addPolicy(root, api *permissionrpc.ApiDetails) {
	if api.Method != "" {
		ap := ApiPolicy{
			desc:      api.Name,
			disable:   api.IsDisable == 1,
			traceable: api.Traceable == 1,
			roles:     make([]string, 0),
		}

		if root != nil {
			ap.module = root.Name
		}

		url := fmt.Sprintf("%s %s", api.Path, api.Method)
		m.policies[url] = &ap
	}

	for _, v := range api.Children {
		m.addPolicy(api, v)
	}
}

func (m *RbacHolder) GetPolicy(path string, method string) IApiPolicy {
	m.RLock()
	defer m.RUnlock()

	url := fmt.Sprintf("%s %s", path, method)
	if v, ok := m.policies[url]; ok {
		return v
	}

	return nil
}
