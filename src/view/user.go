package view

import (
	"context"

	"github.com/jianshao/poker_counter/prisma/db"
	"github.com/jianshao/poker_counter/src/utils"
)

func GetUserInfoByOpenid(openid string) (*db.UserModel, error) {
	client := utils.GetPrismaClient()
	return client.User.FindUnique(db.User.Openid.Equals(openid)).Exec(context.Background())
}

func GetUserById(userId int) (*db.UserModel, error) {
	client := utils.GetPrismaClient()
	return client.User.FindUnique(db.User.ID.Equals(userId)).Exec(context.Background())
}

func CreateOneUser(name, openId string) (*db.UserModel, error) {
	client := utils.GetPrismaClient()
	return client.User.CreateOne(
		db.User.Openid.Set(openId),
		db.User.Name.Set(name),
	).Exec(context.Background())
}
