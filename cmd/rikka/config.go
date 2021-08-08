package main

import (
	"github.com/ictsc/ictsc-rikka/pkg/repository/mariadb"
	"github.com/ictsc/ictsc-rikka/pkg/seed"
)

type Config struct {
	Listen  ListenConfig          `yaml:"listen"`
	CORS    CORSConfig            `yaml:"cors"`
	MariaDB mariadb.MariaDBConfig `yaml:"mariadb"`
	Redis   RedisConfig           `yaml:"redis"`
	Minio   MinioConfig           `yaml:"minio"`
	Seed    seed.SeedConfig       `yaml:"seed"`
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
	BucketName      string `yaml:"bucketName`
	UseSSL          bool   `yaml:"useSSL"`
}
