package sns

import (
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randCode(l int) string {
	var code string
	for i := 0; i < l; i++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}
