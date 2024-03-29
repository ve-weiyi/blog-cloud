package request

type UpdateUserRolesReq struct {
	UserId  int   `json:"user_id"`
	RoleIds []int `json:"role_ids"`
}

type UpdateRoleMenusReq struct {
	RoleId  int   `json:"role_id"`
	MenuIds []int `json:"menu_ids"`
}

type UpdateRoleApisReq struct {
	RoleId      int   `json:"role_id"`
	ResourceIds []int `json:"resource_ids"`
}
