package config

import (
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Instance   string
	Proxy      HttpProxy
	viper      *viper.Viper
	configPath string
}

type HttpProxy struct {
	Url  string
	Port string
}

func Init() Config {
	var conf = Config{}
	conf.configPath = "config.json"

	var vp = viper.New()
	vp.SetConfigFile(conf.configPath)
	vp.AddConfigPath(".")

	log.Println("Reading config.")
	var err = vp.ReadInConfig()
	if err != nil {
		log.Println(err)
	}
	log.Println("Reading config done.")
	conf.viper = vp

	conf.Instance = *(flag.String("config", vp.GetString("piped.instance"), "path to the config file"))
	conf.Proxy.Url = *(flag.String("proxy", vp.GetString("proxy.url"), "proxy url"))
	conf.Proxy.Port = *(flag.String("port", vp.GetString("proxy.port"), "proxy port"))
	flag.Parse()

	return conf
}

func (c *Config) Set(key string, value interface{}) {
	c.viper.Set(key, value)
	if err := c.viper.SafeWriteConfigAs(c.configPath); err != nil {
		if os.IsNotExist(err) {
			err = c.viper.WriteConfigAs(c.configPath)
		}
	}
	c.viper.WriteConfig()
}
