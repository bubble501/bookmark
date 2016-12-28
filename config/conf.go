package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//Configuration is the representation of conf.json in memory.
type Configuration struct {
	Base    map[string]string
	Release map[string]string
	Debug   map[string]string
	Mode    string
}

// func (c *Configuration) GetIntValue(key string, defaultValue int) int {
// 	return 0
// }

func getStringValue(base, derived map[string]string, key ...string) string {
	if val, ok := derived[key[0]]; ok != false {
		return val
	}
	if val, ok := base[key[0]]; ok != false {
		return val
	}

	if len(key) == 2 {
		return key[1]
	}

	return ""

}

//GetStringValue will get the value with specified key in the conf file.
func (c *Configuration) GetStringValue(key ...string) string {
	switch strings.ToLower(c.Mode) {
	case "debug", "":
		return getStringValue(c.Base, c.Debug, key...)
	case "release":
		return getStringValue(c.Base, c.Release, key...)
	default:
		return getStringValue(c.Base, c.Debug, key...)
	}
}

//Singleton is a singleton object.
var Singleton = loadConfig()

func loadConfig() *Configuration {
	file, err := os.Open("conf.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	decoder := json.NewDecoder(file)
	var conf Configuration
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Println(err)
	}
	return &conf
}
