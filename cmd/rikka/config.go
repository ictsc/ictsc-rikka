package main

import (
	"fmt"

	"github.com/ictsc/ictsc-rikka/pkg/seed"
)

type Config struct {
	Contest ContestConfig   `yaml:"contest"`
	Listen  ListenConfig    `yaml:"listen"`
	CORS    CORSConfig      `yaml:"cors"`
	MariaDB MariaDBConfig   `yaml:"mariadb"`
	Redis   RedisConfig     `yaml:"redis"`
	Minio   MinioConfig     `yaml:"minio"`
	Seed    seed.SeedConfig `yaml:"seed"`
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

type CORSConfig struct {
	Origins []string `yaml:"origins"`
}

type RedisConfig struct {
	IdleConnectionSize int    `yaml:"idleConnectionSize"`
	Address            string `yaml:"address"`
	Port               int    `yaml:"port"`
	Password           string `yaml:"password"`
	KeyPair            string `yaml:"keyPair"`
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
