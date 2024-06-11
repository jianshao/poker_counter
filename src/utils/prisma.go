package utils

import (
	"log"

	"github.com/jianshao/poker_counter/prisma/db"
)

var (
	gPrisma *db.PrismaClient = nil
)

func GetPrismaClient() *db.PrismaClient {
	if gPrisma == nil {
		gPrisma = db.NewClient()
		if err := gPrisma.Prisma.Connect(); err != nil {
			log.Println(err.Error())
			gPrisma = nil
		}
	}
	return gPrisma
}

func closePrisma() {
	if gPrisma != nil {
		gPrisma.Prisma.Disconnect()
		gPrisma = nil
	}
}
