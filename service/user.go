package service

import "chat/serializer"

type UserRegisterService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=30"`
	Password string `form:"password" json:"password" binding:"required,min=5,max=30"`
}

func (service *UserRegisterService) Register() serializer.Response {

}
