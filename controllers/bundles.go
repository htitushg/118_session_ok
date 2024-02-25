package controllers

//import "118_session_ok/controllers"

var HomeBundle = Join(Home, Log())
var LoginBundle = Join(Login, Log())
var SigninBundle = Join(Signin, Log())
var IndexBundle = Join(Index, Log(), Guard())
var LogoutBundle = Join(Logout, Log(), Guard(), Foo())
var RegisterBundle = Join(Register, Log(), Foo())
var AfficheUserInfoBundle = Join(AfficheUserInfo, Log(), Guard(), Foo())
var IndexHandlerNoMethBundle = Join(IndexHandlerNoMeth, Log(), Guard(), Foo())
var IndexHandlerOtherBundle = Join(IndexHandlerOther, Log(), Guard(), Foo())
