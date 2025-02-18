package tools

import (
	"strconv"
	"time"

	"golang.org/x/exp/rand"
)

func GenSix() string {
	rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	min := 100000
	max := 999999
	c := rand.Intn(max-min+1) + min
	return strconv.Itoa(c)
}
