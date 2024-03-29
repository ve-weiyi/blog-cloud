package response

import (
	"time"
)

type MenuDetailsDTO struct {
	ID        int               `json:"id"`        // 主键
	ParentID  int               `json:"parent_id"` // 父id
	Title     string            `json:"title"`     // 菜单标题
	Type      int               `json:"type"`      // 菜单类型（0代表菜单、1代表iframe、2代表外链、3代表按钮）
	Path      string            `json:"path"`      // 路由地址
	Name      string            `json:"name"`      // 路由名字
	Component interface{}       `json:"component"` // Layout组件
	Redirect  string            `json:"redirect"`  // 路由重定向
	Meta      Meta              `json:"meta"`      // meta配置
	Children  []*MenuDetailsDTO `json:"children"`
	CreatedAt time.Time         `json:"created_at"` // 创建时间
	UpdatedAt time.Time         `json:"updated_at"` // 更新时间
}

// Meta 对应TypeScript中的meta对象
type Meta struct {
	Title    string `json:"title"`    // 菜单名称
	Icon     string `json:"icon"`     // 菜单图标
	ShowLink bool   `json:"showLink"` // 是否在菜单中显示
	Rank     int    `json:"rank"`     // 菜单升序排序

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
