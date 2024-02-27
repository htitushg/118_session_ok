package controllers

import (
	"118_session_ok/internal/middlewares"
)

var HomeBundle = middlewares.Join(Home, middlewares.Log)
var LoginBundle = middlewares.Join(Login, middlewares.Log)
var SigninBundle = middlewares.Join(Signin, middlewares.Log)
var IndexBundle = middlewares.Join(Index, middlewares.Log, middlewares.Guard)
var LogoutBundle = middlewares.Join(Logout, middlewares.Log, middlewares.Guard, middlewares.Foo)
var RegisterBundle = middlewares.Join(Register, middlewares.Log, middlewares.Foo)
var AfficheUserInfoBundle = middlewares.Join(AfficheUserInfo, middlewares.Log, middlewares.Guard, middlewares.Foo)
var IndexHandlerNoMethBundle = middlewares.Join(IndexHandlerNoMeth, middlewares.Log, middlewares.Foo)
var IndexHandlerOtherBundle = middlewares.Join(IndexHandlerOther, middlewares.Log, middlewares.Foo)
