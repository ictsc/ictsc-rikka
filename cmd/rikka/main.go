package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ListenTLSConfig struct {
	CertFilePath string `yaml:"certFilePath"`
	KeyFilePath  string `yaml:"keyFilePath"`
}
type ListenConfig struct {
	Address string           `yaml:"address"`
	Port    int              `yaml:"port"`
	TLS     *ListenTLSConfig `yaml:"tls"`
}
type Config struct {
	Listen ListenConfig `yaml:"listen"`
}

var (
	configPath string
	config     Config
)

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "config path")
	flag.Parse()

	f, err := os.Open(configPath)
	if err != nil {
		log.Fatalf(errors.Wrapf(err, "Failed to open config file `%s`.", configPath).Error())
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		f.Close()
		log.Fatalf(errors.Wrapf(err, "Failed to decode config.").Error())
	}
}
func main() {
	r := gin.Default()

	addr := fmt.Sprintf("%s:%d", config.Listen.Address, config.Listen.Port)
	if config.Listen.TLS == nil {
		// not serve tls
		if err := r.Run(addr); err != nil {
			log.Fatal(err.Error())
		}
	} else {
		// serve tls
		if err := r.RunTLS(addr, config.Listen.TLS.CertFilePath, config.Listen.TLS.KeyFilePath); err != nil {
			log.Fatal(err.Error())
		}
	}

}
