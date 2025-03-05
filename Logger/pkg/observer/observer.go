package observer

import "Logger/pkg/loggerfactory"

//Listen and update the log level
type Observer interface{
	UpdateLogger(logLevelMap *map[string]string, slogHandlerConfig loggerfactory.SlogHandlerConfig)
}