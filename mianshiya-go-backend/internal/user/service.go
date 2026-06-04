package user

import "errors"

func Register(req *RegisterRequest) (int64, error) {
	if req.UserAccount == "" || req.UserPassword == "" || req.CheckPassword == "" {
		return 0, errors.New("账号或密码不能为空")
	}
	if req.UserPassword != req.CheckPassword {
		return 0, errors.New("两次输入的密码不一致")
	}

	if len(req.UserAccount) < 4 {
		return 0, errors.New("账号长度不能少于4位")
	}
	if len(req.UserPassword) < 8 {
		return 0, errors.New("密码长度不能少于8位")
	}
	return 1, nil
}
