package yandex

import (
	"time"

	e "github.com/nikitaSstepanov/tools/error"
)

const (
	yandexTable  = "yandex"
	redisExpires = 3 * time.Hour
)

var (
	notFoundErr = e.New("This yandex integration wasn`t found.", e.NotFound)
)
