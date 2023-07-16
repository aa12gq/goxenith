package config

import "goxenith/pkg/config"

func init() {

	config.Add("database", func() map[string]interface{} {
		return map[string]interface{}{

			"connection": config.Env("DB_CONNECTION", "mysql"),

			"mysql": map[string]interface{}{

				"host":     config.Env("DB_HOST", "127.0.0.1"),
				"port":     config.Env("DB_PORT", "3306"),
				"database": config.Env("DB_DATABASE", "goxenith"),
				"username": config.Env("DB_USERNAME", ""),
				"password": config.Env("DB_PASSWORD", ""),
				"charset":  "utf8mb4",
			},
		}
	})
}
