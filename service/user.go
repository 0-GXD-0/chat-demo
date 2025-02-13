package service

import (
	"chat/model"
	"chat/serializer"
)

type UserRegisterService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=5,max=30"`
}

func (service *UserRegisterService) Register() serializer.Response {
	var user model.User
	code := 200
	count := 0
	model.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).Count(&count)
	if count != 0 {
		return serializer.Response{
			Status: 400,
			Msg:    "用户名已存在",
		}
	}
	user = model.User{
		UserName: service.UserName,
	}
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "密码加密失败",
			Error:  err.Error(),
		}
	}
	model.DB.Create(&user)
	return serializer.Response{
		Status: code,
		Msg:    "注册成功",
	}
}
