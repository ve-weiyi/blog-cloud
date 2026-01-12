package types

type ApiBackVO struct {
	Id        int64        `json:"id,optional"`     // 主键id
	ParentId  int64        `json:"parent_id"`       // 分组id
	Name      string       `json:"name"`            // api名称
	Path      string       `json:"path"`            // api路径
	Method    string       `json:"method"`          // api请求方法
	Traceable int64        `json:"traceable"`       // 是否追溯操作记录 0需要，1是
	Status    int64        `json:"status,optional"` // 状态 0正常 1禁用
	CreatedAt int64        `json:"created_at"`      // 创建时间
	UpdatedAt int64        `json:"updated_at"`      // 更新时间
	Children  []*ApiBackVO `json:"children"`
}

type NewApiReq struct {
	Id        int64  `json:"id,optional"`     // 主键id
	ParentId  int64  `json:"parent_id"`       // 分组id
	Name      string `json:"name"`            // api名称
	Path      string `json:"path"`            // api路径
	Method    string `json:"method"`          // api请求方法
	Traceable int64  `json:"traceable"`       // 是否追溯操作记录 0需要，1是
	Status    int64  `json:"status,optional"` // 状态 0正常 1禁用
}

type QueryApiReq struct {
	PageQuery
	Name   string `json:"name,optional"`   // api名称
	Path   string `json:"path,optional"`   // api路径
	Method string `json:"method,optional"` // api请求方法
}

type SyncApiReq struct {
}
