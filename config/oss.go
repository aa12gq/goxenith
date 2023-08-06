package config

import "goxenith/pkg/config"

func init() {

	config.Add("oss", func() map[string]interface{} {
		return map[string]interface{}{

			"accessKeyId":     config.Env("ACCESSKEYID", ""),
			"accessKeySecret": config.Env("ACCESSKEYSECRET", ""),
			"endpoint":        config.Env("ENDPOINT", ""),
			"bucketName":      config.Env("BUCKETNAME", ""),
			"region":          config.Env("REGION", "华北2（北京）"),
			// 图片访问路径
			"access_path": config.Env("ACCESSPATH", ""),
			// 最大闲置连接数。
			"maxIdleConns": config.Env("OSS_MAX_IDLE_CONNS", "100"),
			// 每个主机的最大闲置连接数
			"maxIdleConnsPerHost": config.Env("OSS_MAX_IDLE_CONNS_PER_HOST", "100"),
			// 每个主机的最大连接数
			"maxConnsPerHost": config.Env("OSS_MAX_CONNS_PER_HOST", "0"),
		}
	})
}
