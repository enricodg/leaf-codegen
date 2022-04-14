package command

import (
	leafLogger "github.com/paulusrobin/leaf-utilities/logger/logger"
	"leaf-codegen/logger"
	"sync"
)

var (
	instance command
	once     sync.Once
)

type (
	command struct {
		log leafLogger.Logger
	}
)

func GetCommand() command {
	once.Do(func() {
		instance = command{
			log: logger.GetLogger(),
		}
	})
	return instance
}
