package types

type NewRoleReq struct {
	Id          int64  `json:"id,optional"`        // 主键id
	ParentId    int64  `json:"parent_id,optional"` // 父角色id
	RoleKey     string `json:"role_key"`           // 角色名
	RoleLabel   string `json:"role_label"`         // 角色标签
	RoleComment string `json:"role_comment"`       // 角色备注
	IsDefault   int64  `json:"is_default"`         // 是否默认角色 0否 1是
	Status      int64  `json:"status"`             // 状态 0正常 1禁用
}

type QueryRoleReq struct {
	PageQuery
	RoleKey   string `json:"role_key,optional"`   // 角色名
	RoleLabel string `json:"role_label,optional"` // 角色标签
	Status    int64  `json:"status,optional"`     // 状态 0正常 1禁用
}

type RoleBackVO struct {
	Id          int64  `json:"id,optional"`  // 主键id
	ParentId    int64  `json:"parent_id"`    // 父角色id
	RoleKey     string `json:"role_key"`     // 角色名
	RoleLabel   string `json:"role_label"`   // 角色标签
	RoleComment string `json:"role_comment"` // 角色备注
	IsDefault   int64  `json:"is_default"`   // 是否默认角色 0否 1是
	Status      int64  `json:"status"`       // 状态 0正常 1禁用
	CreatedAt   int64  `json:"created_at"`   // 创建时间
	UpdatedAt   int64  `json:"updated_at"`   // 更新时间
}

type RoleResourcesResp struct {
	RoleId  int64   `json:"role_id"`
	ApiIds  []int64 `json:"api_ids"`
	MenuIds []int64 `json:"menu_ids"`
}

type UpdateRoleApisReq struct {
	RoleId int64   `json:"role_id"`
	ApiIds []int64 `json:"api_ids"`
}

type UpdateRoleMenusReq struct {
	RoleId  int64   `json:"role_id"`
	MenuIds []int64 `json:"menu_ids"`
}
