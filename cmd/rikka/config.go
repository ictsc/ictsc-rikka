package main

import (
	"fmt"

	"github.com/ictsc/ictsc-rikka/pkg/seed"
)

type Config struct {
	Rikka    RikkaConfig     `yaml:"rikka"`
	Contest  ContestConfig   `yaml:"contest"`
	Listen   ListenConfig    `yaml:"listen"`
	CORS     CORSConfig      `yaml:"cors"`
	Store    StoreConfig     `yaml:"store"`
	Notify   NotifyConfig    `yaml:"notify"`
	MariaDB  MariaDBConfig   `yaml:"mariadb"`
	Redis    RedisConfig     `yaml:"redis"`
	Minio    MinioConfig     `yaml:"minio"`
	Seed     seed.SeedConfig `yaml:"seed"`
	Sentry   SentryConfig    `yaml:"sentry"`
	Recreate RecreateConfig  `yaml:"recreate"`
}

type RikkaConfig struct {
	// 点数やランキングがユーザーには表示されなくなる予選用モード
	PreRoundMode bool `yaml:"preRoundMode"`
}

type ContestConfig struct {
	AnswerLimit int `yaml:"answerLimit"`
}

type ListenTLSConfig struct {
	CertFilePath string `yaml:"certFilePath"`
	KeyFilePath  string `yaml:"keyFilePath"`
}

type ListenConfig struct {
	Address string           `yaml:"address"`
	Port    int              `yaml:"port"`
	TLS     *ListenTLSConfig `yaml:"tls"`
}

type NotifyConfig struct {
	Answer string `yaml:"answer"`
}

type CORSConfig struct {
	Origins []string `yaml:"origins"`
}

type StoreConfig struct {
	Secure             bool   `yaml:"secure"`
	SameSiteStrictMode string `yaml:"sameSiteStrictMode"`
	Domain             string `yaml:"domain"`
}

type RedisConfig struct {
	IdleConnectionSize int    `yaml:"idleConnectionSize"`
	Address            string `yaml:"address"`
	Port               int    `yaml:"port"`
	Password           string `yaml:"password"`
	KeyPair            string `yaml:"keyPair"`
}

type RecreateConfig struct {
	URL string `yaml:"url"`
}

type MinioConfig struct {
	Endpoint        string `yaml:"endPoint"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	BucketName      string `yaml:"bucketName"`
	UseSSL          bool   `yaml:"useSSL"`
}

type MariaDBConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	MaxConn  int    `yaml:"maxConn"`
	Database string `yaml:"database"`
}

func (c *MariaDBConfig) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Address,
		c.Port,
		c.Database,
	)
}

type SentryConfig struct {
	Dsn              string  `yaml:"dsn"`
	Environment      string  `yaml:"environment"`
	TracesSampleRate float64 `yaml:"sampleRate"`
}
