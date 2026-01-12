package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type UserController struct {
	svcCtx *svctx.ServiceContext
}

func NewUserController(svcCtx *svctx.ServiceContext) *UserController {
	return &UserController{
		svcCtx: svcCtx,
	}
}

// @Tags		User
// @Summary		"删除用户绑定第三方平台账号"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.DeleteUserBindThirdPartyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/user/delete_user_bind_third_party [POST]
func (s *UserController) DeleteUserBindThirdParty(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.DeleteUserBindThirdPartyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUserLogic(s.svcCtx).DeleteUserBindThirdParty(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		User
// @Summary		"获取用户信息"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.UserInfoResp}	"返回信息"
// @Router		/blog-api/v1/user/get_user_info [GET]
func (s *UserController) GetUserInfo(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.EmptyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUserLogic(s.svcCtx).GetUserInfo(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		User
// @Summary		"获取用户点赞列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.UserLikeResp}	"返回信息"
// @Router		/blog-api/v1/user/get_user_like [GET]
func (s *UserController) GetUserLike(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.EmptyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUserLogic(s.svcCtx).GetUserLike(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		User
// @Summary		"修改用户头像"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateUserAvatarReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/user/update_user_avatar [POST]
func (s *UserController) UpdateUserAvatar(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateUserAvatarReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUserLogic(s.svcCtx).UpdateUserAvatar(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		User
// @Summary		"修改用户绑定邮箱"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateUserBindEmailReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/user/update_user_bind_email [POST]
func (s *UserController) UpdateUserBindEmail(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateUserBindEmailReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUserLogic(s.svcCtx).UpdateUserBindEmail(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		User
// @Summary		"修改用户绑定手机号"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateUserBindPhoneReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/user/update_user_bind_phone [POST]
func (s *UserController) UpdateUserBindPhone(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateUserBindPhoneReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUserLogic(s.svcCtx).UpdateUserBindPhone(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		User
// @Summary		"修改用户绑定第三方平台账号"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateUserBindThirdPartyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/user/update_user_bind_third_party [POST]
func (s *UserController) UpdateUserBindThirdParty(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateUserBindThirdPartyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUserLogic(s.svcCtx).UpdateUserBindThirdParty(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		User
// @Summary		"修改用户信息"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateUserInfoReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/user/update_user_info [POST]
func (s *UserController) UpdateUserInfo(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateUserInfoReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUserLogic(s.svcCtx).UpdateUserInfo(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		User
// @Summary		"修改用户密码"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateUserPasswordReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/user/update_user_password [POST]
func (s *UserController) UpdateUserPassword(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateUserPasswordReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewUserLogic(s.svcCtx).UpdateUserPassword(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
