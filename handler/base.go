package handler

import (
	leafLogger "github.com/paulusrobin/leaf-utilities/logger/logger"
	"leaf-codegen/logger"
	"sync"
)

var (
	instance handler
	once     sync.Once
)

type (
	handler struct {
		log leafLogger.Logger
	}
)

func GetHandler() handler {
	once.Do(func() {
		instance = handler{
			log: logger.GetLogger(),
		}
	})
	return instance
}
