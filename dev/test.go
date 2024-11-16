package main

import (
	"github.com/go-redis/redis/v8"
	ominitor "github.com/stormi-li/ominitor-v1"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	c := ominitor.NewClient(&redis.Options{Addr: redisAddr, Password: password})
	c.Develop("118.25.196.166:9998")
}
