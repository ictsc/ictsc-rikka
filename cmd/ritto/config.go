package main

import "fmt"

type Config struct {
	MariaDB MariaDBConfig `yaml:"mariadb"`
	Rikka   Rikka         `yaml:"rikka"`
	Growi   GrowiConfig   `yaml:"growi"`
}

type MariaDBConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	MaxConn  int    `yaml:"maxConn"`
	Database string `yaml:"database"`
}

type Rikka struct {
	AuthorId string `yaml:"authorId"`
}

type GrowiConfig struct {
	Url         string `yaml:"url"`
	Token       string `yaml:"token"`
	ProblemPath string `yaml:"problemPath"`
	NoticePath  string `yaml:"noticePath"`
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
