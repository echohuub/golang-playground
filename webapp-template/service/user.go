package service

import (
	"webapp.demo/dao/mysql"
	"webapp.demo/model"
)

func SignUp(p model.ParamSignUp) {
	mysql.InsertUser()
}
