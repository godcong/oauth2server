package oauth2server

import (
	"github.com/godcong/oauth2server/model"

	"strconv"
)

type UserAuthorizeInfo struct {
	User model.User
}

type ResponseMessage struct {
	Code    string
	Message string
	Data    interface{}
}

var (
	//登录
	A_LOGIN_SUCCESS                = ResponseMessage{strconv.Itoa(0), "登陆成功！", nil}
	A_NAME_OR_PASSWORD_CANNOT_NULL = ResponseMessage{strconv.Itoa(1), "用户名或密码不能为空!", nil}
	A_MOBILE_OR_CODE_CANNOT_NULL   = ResponseMessage{strconv.Itoa(2), "手机或验证码不能为空!", nil}
	A_NAME_OR_PASSWORD_WRONG       = ResponseMessage{strconv.Itoa(3), "用户名或密码错误!", nil}
	A_COMPLETE_USERINFO_SUCCESS    = ResponseMessage{strconv.Itoa(4), "成功完善用户信息！！！", nil}

	//注册
	A_REGISTER_SUCCESS = ResponseMessage{strconv.Itoa(0), "注册成功！", nil}

	A_MUST_NOT_BE_NULL       = ResponseMessage{strconv.Itoa(1), "必填字段不能为空！", nil}
	A_NAME_FORMAT_WRONG      = ResponseMessage{strconv.Itoa(2), "用户名格式错误！", nil}
	A_NAME_IS_EXISTS         = ResponseMessage{strconv.Itoa(3), "用户名已经存在！", nil}
	A_NAME_LENGTH_WRONG      = ResponseMessage{strconv.Itoa(4), "用户名不能少于六位！", nil}
	A_MOBILE_LENGTH_WRONG    = ResponseMessage{strconv.Itoa(5), "手机号填写错误！", nil}
	A_MOBILE_IS_EXISTS       = ResponseMessage{strconv.Itoa(3), "手机号已经被注册！", nil}
	A_MESSAGE_WRONG          = ResponseMessage{strconv.Itoa(6), "手机号验证码错误！", nil}
	A_PASSWORD_LENGTH_WRONG  = ResponseMessage{strconv.Itoa(7), "密码不能少于六位！", nil}
	A_PASSWORD_COMPARE_WRONG = ResponseMessage{strconv.Itoa(8), "两次密码输入不一致！", nil}

	//
	A_MESSAGE_SEND_SUCCESS  = ResponseMessage{strconv.Itoa(0), "信息发送成功！", nil}
	A_MESSAGE_SEND_FAILED   = ResponseMessage{strconv.Itoa(1), "信息发送失败！", nil}
	A_MESSAGE_CHECK_SUCCESS = ResponseMessage{strconv.Itoa(0), "验证码正确！", nil}
	A_MESSAGE_CHECK_FAILED  = ResponseMessage{strconv.Itoa(1), "验证码错误！", nil}

	A_FORGET_CHECK_SUCCESS        = ResponseMessage{strconv.Itoa(0), "验证通过！", nil}
	A_FORGET_CHANGE_SUCCESS       = ResponseMessage{strconv.Itoa(0), "密码变更成功!", nil}
	A_FORGET_NEED_IS_NULL         = ResponseMessage{strconv.Itoa(1), "缺少必填字段!！", nil}
	A_FORGET_ACCOUNT_WRONG        = ResponseMessage{strconv.Itoa(2), "账户填写错误!！", nil}
	A_FORGET_PASSWORD_CANNOT_NULL = ResponseMessage{strconv.Itoa(3), "密码不能为空!", nil}
	A_FORGET_PASSWORD_FILL_OUT    = ResponseMessage{strconv.Itoa(4), "输入密码不一致!", nil}
	A_FORGET_SYSTEM_ERROR         = ResponseMessage{strconv.Itoa(5), "系统错误,请重试！", nil}

	//修改密码
	A_CHANGE_SUCCESS            = ResponseMessage{strconv.Itoa(0), "密码修改成功!", nil}
	A_CHANGE_NEED_IS_NULL       = ResponseMessage{strconv.Itoa(1), "缺少必填字段!！", nil}
	A_CHANGE_OPASSWORD_WRONG    = ResponseMessage{strconv.Itoa(2), "原密码错误!", nil}
	A_CHANGE_NPASSWORD_FILL_OUT = ResponseMessage{strconv.Itoa(3), "新输入密码不一致!", nil}
)
