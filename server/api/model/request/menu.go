package request

// CustomizeRouteMeta 对应TypeScript中的CustomizeRouteMeta接口
//type CustomizeRouteMeta struct {
//	Title        string      `json:"title"`                  // 菜单名称
//	Icon         string      `json:"icon,omitempty"`         // 菜单图标
//	ExtraIcon    interface{} `json:"extraIcon,omitempty"`    // 菜单名称右侧的额外图标
//	ShowLink     bool        `json:"showLink,omitempty"`     // 是否在菜单中显示
//	ShowParent   bool        `json:"showParent,omitempty"`   // 是否显示父级菜单
//	Roles        []string    `json:"roles,omitempty"`        // 页面级别权限设置
//	Auths        []string    `json:"auths,omitempty"`        // 按钮级别权限设置
//	KeepAlive    bool        `json:"keepAlive,omitempty"`    // 路由组件缓存
//	FrameSrc     string      `json:"frameSrc,omitempty"`     // 内嵌的iframe链接
//	FrameLoading bool        `json:"frameLoading,omitempty"` // iframe页是否开启首次加载动画
//	Transition   Transition  `json:"transition,omitempty"`   // 页面加载动画
//	HiddenTag    bool        `json:"hiddenTag,omitempty"`    // 是否不添加信息到标签页
//	DynamicLevel int         `json:"dynamicLevel,omitempty"` // 动态路由可打开的最大数量
//	ActivePath   string      `json:"activePath,omitempty"`   // 将某个菜单激活
//}

// Meta 对应TypeScript中的meta对象
type Meta struct {
	Title    string `json:"title"`              // 菜单名称
	Icon     string `json:"icon,omitempty"`     // 菜单图标
	ShowLink bool   `json:"showLink,omitempty"` // 是否在菜单中显示
	Rank     int    `json:"rank,omitempty"`     // 菜单升序排序

	// 子菜单才有的属性
	ExtraIcon    interface{} `json:"extraIcon,omitempty"`    // 菜单名称右侧的额外图标
	ShowParent   bool        `json:"showParent,omitempty"`   // 是否显示父级菜单
	Roles        []string    `json:"roles,omitempty"`        // 页面级别权限设置
	Auths        []string    `json:"auths,omitempty"`        // 按钮级别权限设置
	KeepAlive    bool        `json:"keepAlive,omitempty"`    // 路由组件缓存
	FrameSrc     string      `json:"frameSrc,omitempty"`     // 内嵌的iframe链接
	FrameLoading bool        `json:"frameLoading,omitempty"` // iframe页是否开启首次加载动画
	Transition   Transition  `json:"transition,omitempty"`   // 页面加载动画
	HiddenTag    bool        `json:"hiddenTag,omitempty"`    // 是否不添加信息到标签页
	DynamicLevel int         `json:"dynamicLevel,omitempty"` // 动态路由可打开的最大数量
	ActivePath   string      `json:"activePath,omitempty"`   // 将某个菜单激活
}

// Transition 对应TypeScript中的transition对象
type Transition struct {
	Name            string `json:"name,omitempty"`            // 当前路由动画效果
	EnterTransition string `json:"enterTransition,omitempty"` // 进场动画
	LeaveTransition string `json:"leaveTransition,omitempty"` // 离场动画
}

// RouteChildrenConfigsTable 对应TypeScript中的RouteChildrenConfigsTable接口
//type RouteChildrenConfigsTable struct {
//	Path      string                      `json:"path"`                // 子路由地址
//	Name      string                      `json:"name,omitempty"`      // 路由名字
//	Redirect  string                      `json:"redirect,omitempty"`  // 路由重定向
//	Component interface{}                 `json:"component,omitempty"` // 按需加载组件
//	Meta      CustomizeRouteMeta          `json:"meta,omitempty"`      // meta配置
//	Children  []RouteChildrenConfigsTable `json:"children,omitempty"`  // 子路由配置项
//}

// RouteConfigsTable 对应TypeScript中的RouteConfigsTable接口
type RouteConfigsTable struct {
	Type      int                 `json:"type"`                // 菜单类型（0代表菜单、1代表iframe、2代表外链、3代表按钮）
	Path      string              `json:"path"`                // 路由地址
	Name      string              `json:"name,omitempty"`      // 路由名字
	Component interface{}         `json:"component,omitempty"` // Layout组件
	Redirect  string              `json:"redirect,omitempty"`  // 路由重定向
	Meta      Meta                `json:"meta,omitempty"`      // meta配置
	Children  []RouteConfigsTable `json:"children,omitempty"`  // 子路由配置项
}

type SyncMenuRequest struct {
	Menus []RouteConfigsTable `json:"menus"`
}
