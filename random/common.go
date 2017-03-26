package random

import (
	"math/rand"
	"time"
)

var RandomSource = rand.New(rand.NewSource(time.Now().UnixNano()))
