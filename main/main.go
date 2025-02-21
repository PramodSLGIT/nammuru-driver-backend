package main

import (
	"nammuru-driver-backend/api"
	"nammuru-driver-backend/utils"
)

func main() {
	var redis = new(utils.RedisConfig)
	redis.InitRedis()
	api.DriverApi()
}
