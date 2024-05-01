package misc

import (
	"log"
	"strconv"
)

type Address = string

// the default method
func Addr2Shard(addr Address) int {
	last16Addr := addr[len(addr)-8:]
	num, err := strconv.ParseUint(last16Addr, 16, 64)
	if err != nil {
		log.Panic(err)
	}
	return int(num) % 1
}
