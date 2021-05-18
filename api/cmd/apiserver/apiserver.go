package apiserver

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

}
