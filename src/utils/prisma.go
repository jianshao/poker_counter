package utils

import (
	"log"
	"private/backend/poker_counter/prisma/db"
)

var (
	gPrisma *db.PrismaClient = nil
)

func GetPrismaClient() *db.PrismaClient {
	if gPrisma != nil {
		return gPrisma
	}

	gPrisma = db.NewClient()
	if err := gPrisma.Prisma.Connect(); err != nil {
		log.Println(err.Error())
		gPrisma = nil
	}
	return gPrisma
}

func Init() {

}

func Close() {
	gPrisma.Prisma.Disconnect()
	gPrisma = nil
}
