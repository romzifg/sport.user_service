package controllers

import (
	"net/http"
	errWrap "user-service/common/error"
	"user-service/common/response"
	"user-service/domain/dto"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	service services.IServiceRegistry
}

type IUserController interface {
	Login(*gin.Context)
	Register(*gin.Context)
	Update(*gin.Context)
	UpdatePassword(*gin.Context)
	GetUserLogin(*gin.Context)
	GetUserByUUID(*gin.Context)
}

func NewUserController(service services.IServiceRegistry) IUserController {
	return &UserController{
		service: service,
	}
}

func (u *UserController) Login(ctx *gin.Context) {
	request := &dto.LoginRequest{}
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     ctx,
		})
		return
	}

	user, err := u.service.GetUser().Login(ctx, request)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHttpResponse{
		Code:  http.StatusOK,
		Data:  user.User,
		Token: &user.Token,
		Gin:   ctx,
	})
}

func (u *UserController) Register(ctx *gin.Context) {
	request := &dto.RegiterRequest{}
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     ctx,
		})
		return
	}

	user, err := u.service.GetUser().Register(ctx, request)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHttpResponse{
		Code: http.StatusOK,
		Data: user.User,
		Gin:  ctx,
	})
}

func (u *UserController) Update(ctx *gin.Context) {
	request := &dto.UpdateRequest{}
	uuid := ctx.Param("uuid")
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     ctx,
		})
		return
	}

	user, err := u.service.GetUser().Update(ctx, request, uuid)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHttpResponse{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

func (u *UserController) UpdatePassword(ctx *gin.Context) {
	request := &dto.UpdatePasswordRequest{}
	uuid := ctx.Param("uuid")
	err := ctx.ShouldBindJSON(request)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.HttpResponse(response.ParamHttpResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: &errMessage,
			Data:    errResponse,
			Error:   err,
			Gin:     ctx,
		})
		return
	}

	user, err := u.service.GetUser().UpdatePassword(ctx, request, uuid)
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHttpResponse{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

func (u *UserController) GetUserLogin(ctx *gin.Context) {
	user, err := u.service.GetUser().GetUserLogin(ctx.Request.Context())
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHttpResponse{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}

func (u *UserController) GetUserByUUID(ctx *gin.Context) {
	user, err := u.service.GetUser().GetUserByUUID(ctx.Request.Context(), ctx.Param("uuid"))
	if err != nil {
		response.HttpResponse(response.ParamHttpResponse{
			Code:  http.StatusBadRequest,
			Error: err,
			Gin:   ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHttpResponse{
		Code: http.StatusOK,
		Data: user,
		Gin:  ctx,
	})
}
