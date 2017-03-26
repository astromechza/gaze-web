package random

import (
	"time"

	"github.com/oklog/ulid"
)

func NewUlidNow() ulid.ULID {
	return NewUlidThen(time.Now())
}

func NewUlidThen(t time.Time) ulid.ULID {
	return ulid.MustNew(ulid.Timestamp(t), RandomSource)
}
