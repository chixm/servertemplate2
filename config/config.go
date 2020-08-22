package config

// load Configuration File in root path.
// config.json for default

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config Parameter Holder
var config *Configuration

// GetConfig get instance of Configuration, run after InitializeConfig() have finished.
func GetConfig() *Configuration {
	return config
}

// InitializeConfig read configuration file and load and instantiate Configuration
func InitializeConfig() {
	log.Println(`[Configuration]`)
	loadConfiguration()
}

// loading configuration file
// configuration file should be placed in root with name 'config.json' or 'config.production.json'
func loadConfiguration() {
	// first argument of binary is environment parameter.
	if len(os.Args) > 1 {
		environment := os.Args[1]
		viper.SetConfigName(`config.` + environment)
	} else {
		viper.SetConfigName(`config`)
	}

	viper.AddConfigPath(`../../`)
	viper.AddConfigPath(`.`)

	c := Configuration{}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	printConfiguration(&c)

	config = &c
}

func printConfiguration(c *Configuration) {
	for _, v := range c.Database {
		log.Println(`Loaded Database Configuration of ::` + v.Id + "[" + v.Host + ":" + strconv.Itoa(v.Port) + "]")
	}
	for _, r := range c.Redis {
		log.Println(`Loaded Redis Configuration of ::` + r.Id + "[" + r.Host + "]")
	}

}

// Configuration ALL Configuration File Contents Structure
type Configuration struct {
	Port     int            // server port
	Database []*DbConfig    //database configuration
	Redis    []*RedisConfig //redis configuration
	Email    *EmailConfig   // mail configuration
}

// DbConfig basic database configuration
type DbConfig struct {
	Id       string
	Name     string
	Host     string
	Port     int
	Username string
	Password string
	MaxIdle  int
	MaxOpen  int
}

// RedisConfig config for Redis
type RedisConfig struct {
	Id        string // redis connection identifier
	Host      string //redis server (ip or domain)
	Port      int    // redis port
	MaxIdle   int    // connection Idle max count
	MaxActive int    // connections Active limit
}

// EmailConfig config for email
type EmailConfig struct {
	Smtp         string // smtp server
	SmtpSvr      string // smtp server to access
	User         string // userid for authorization
	Password     string // password for authorization
	TestSendAddr string // sends email to this address for test when server launched.
}
