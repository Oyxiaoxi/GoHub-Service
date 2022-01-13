// Package factories 存放工厂方法
package factories

import (
	"GoHub-Service/app/models/user"
	"GoHub-Service/pkg/helpers"

	"github.com/bxcodec/faker/v3"
)

func MakeUsers(times int) []user.User {

	var objs []user.User

	// 设置唯一值
	faker.SetGenerateUniqueValues(true)

	for i := 0; i < times; i++ {
		model := user.User{
			Name:     faker.Username(),
			Email:    faker.Email(),
			Phone:    helpers.RandomNumber(11),
			Password: "$2a$14$qU0sAfCMz3e73ZJOoLLU8eAppCkS/8P/LvU7FzUbwRa.RvEOga8bm", // secret
		}
		objs = append(objs, model)
	}

	return objs
}
