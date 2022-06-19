package storage

import (
	"math/rand"
	"time"
)

func Random() int {
	randomRange := 4
	rand.Seed(time.Now().UnixNano())
	res := rand.Intn(randomRange)
	return res
}
