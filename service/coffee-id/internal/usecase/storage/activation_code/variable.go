package activation_code

import (
	"time"

	e "github.com/nikitaSstepanov/tools/error"
)

const (
	redisExpires = 5 * time.Minute
)

var (
	notFoundErr = e.New("This code wasn`t found.", e.NotFound)
)
