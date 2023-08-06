package oss

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

var Bucket *oss.Bucket

type OssConfig struct {
	Endpoint            string
	AccessKeyId         string
	AccessKeySecret     string
	Region              string
	BucketName          string
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
}

func NewOssService(cfg *OssConfig) error {
	clientOptions := []oss.ClientOption{
		oss.MaxConns(cfg.MaxIdleConns, cfg.MaxIdleConnsPerHost, cfg.MaxConnsPerHost),
	}
	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret, clientOptions...)
	if err != nil {
		return err
	}

	Bucket, err = client.Bucket(cfg.BucketName)
	if err != nil {
		return err
	}

	return nil
}
