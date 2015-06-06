package packwrap

import (
	"github.com/jlgerber/logger"
)

// GetLogger - function to return the instance of the logger, created in this
// file. There should only be one logger in packwrap
func GetLogger() *logger.Logger { return logger.GetLogger() }

// Our logger
var log = GetLogger()
