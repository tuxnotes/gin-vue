package dto

import "oceanlearn.teach/ginessential/model"

type UserDto struct { // 只返回给前端用户名和手机号，其他都不用返回
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:      user.Name,
		Telephone: user.Telephone,
	}
}
