package packwrap

import (
	"github.com/jlgerber/logger"
)

func GetLogger() *logger.Logger { return logger.GetLogger() }

var log = GetLogger()
