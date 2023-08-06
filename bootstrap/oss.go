package bootstrap

import (
	"goxenith/pkg/config"
	"goxenith/pkg/logger"
	"goxenith/pkg/oss"
)

func SetupOss() {
	cfg := &oss.OssConfig{
		Endpoint:            config.Get("oss.endpoint"),
		AccessKeyId:         config.Get("oss.accessKeyId"),
		AccessKeySecret:     config.Get("oss.accessKeySecret"),
		BucketName:          config.Get("oss.bucketName"),
		Region:              config.Get("oss.region"),
		MaxIdleConns:        config.GetInt("oss.maxIdleConns"),
		MaxIdleConnsPerHost: config.GetInt("maxIdleConnsPerHost"),
		MaxConnsPerHost:     config.GetInt("maxConnsPerHost"),
	}

	err := oss.NewOssService(cfg)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	return
}
