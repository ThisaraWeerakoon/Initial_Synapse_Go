package config

import (
	"context"
	"encoding"
	"fmt"
	"log"
	"net/url"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


type Config struct {
	koanf *koanf.Koanf
}

type validator interface {
	Validate() error
}

func MustReadFile(filename string) *Config {
	cfg, err := ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return cfg
}

func ReadFile(filename string) (*Config, error) {
	k := koanf.New(".")
	f := file.Provider(filename)
	if err := k.Load(f, yaml.Parser()); err != nil {
		return nil, err
	}

	cfg := &Config{
		koanf: k,
	}

	return cfg, nil
}


func (c *Config) IsSet(key string) bool {
	return c.koanf.Exists(key)
}

func (c *Config) Watch(ctx context.Context, filename string, loggerConfig *zap.Config) { 
	f:= file.Provider(filename)
	f.Watch(func(event interface{}, err error) {
		if err != nil {
			log.Printf("watch error: %v", err)
			return
		}
		// Throw away the old config and load a fresh copy.
		log.Println("config changed. Reloading ...")
		new_k := koanf.New(".")
		if err := new_k.Load(f, yaml.Parser()); err != nil {
			log.Printf("error loading new config: %v",err)
			return
		}
		c.koanf = new_k
		newLevelAsString := c.koanf.Get("logger.level").(string)
		level,err:= zapcore.ParseLevel(newLevelAsString)
		if err != nil {
			log.Printf("error parsing new level: %v",err)
			return
		}
		if loggerConfig != nil { // Check if loggerConfig is set
			loggerConfig.Level.SetLevel(level)
		}
	})	
}

func (c *Config) Unmarshal(key string, out interface{}) error {

	// err := c.koanf.UnmarshalWithConf(key, out, koanf.UnmarshalConf{
	// 	DecoderConfig: &mapstructure.DecoderConfig{
	// 		DecodeHook: mapstructure.ComposeDecodeHookFunc(
	// 			mapstructure.StringToTimeDurationHookFunc(),
	// 			mapstructure.StringToSliceHookFunc(","),
	// 			StringToFunctionHookFunc(),
	// 			StringToUrlHookFunc(),
	// 		),
	// 	},
	// })

	err := c.koanf.Unmarshal(key, out)

	if err != nil {
		return fmt.Errorf("cannot unmarshal config for key %q: %v", key, err)
	}
	if v, ok := out.(validator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("invalid configuration for key %q: %v", key, err)
		}
	}
	return nil
}

func (c *Config) MustUnmarshal(key string, out interface{}) {
	err := c.Unmarshal(key, out)
	if err != nil {
		panic(err)
	}
}

func StringToFunctionHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || !(t.Kind() == reflect.Func || t.Kind() == reflect.Struct) {
			return data, nil
		}
		v := reflect.New(t)
		u, ok := v.Interface().(encoding.TextUnmarshaler)
		if ok {
			err := u.UnmarshalText([]byte(data.(string)))
			if err != nil {
				return nil, err
			}
			return v.Elem().Interface(), nil
		}
		return data, nil
	}
}

func StringToUrlHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		if t != reflect.TypeOf(url.URL{}) {
			return data, nil
		}

		return url.Parse(data.(string))
	}
}

