package configuration

import (
	"os"
	"errors"
	"strconv"
)

type Config struct {
	Parallelism int
	Delay int // in milliseconds
}

// Retrieve configuration from environment
func GetConfig() (*Config, error) {

	// Get parallelism
	parallelismVar, ok := os.LookupEnv("parallelism");
	if !ok {
		return nil , errors.New("parallelism is not present")
	}
	parallelism, err := strconv.Atoi(parallelismVar)
	if err != nil {
		return nil , errors.New("Cannot convert to int " + parallelismVar)
	}

	// Get delay
	delay := 0
	delayVar, ok := os.LookupEnv("delay");
	if ok {
		d, err := strconv.Atoi(delayVar)
		if err == nil {
			delay = d
		}
	}

	return &Config{parallelism, delay}, nil
}