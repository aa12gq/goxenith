package bootstrap

import (
	"goxenith/pkg/config"
	"goxenith/pkg/logger"
)

func SetupLogger() {

	logger.InitLogger(
		config.GetString("log.filename"),
		config.GetInt("log.max_size"),
		config.GetInt("log.max_backup"),
		config.GetInt("log.max_age"),
		config.GetBool("log.compress"),
		config.GetString("log.type"),
		config.GetString("log.level"),
	)
}
