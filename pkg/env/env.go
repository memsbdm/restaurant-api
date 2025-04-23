package env

import (
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

// GetString retrieves the value associated with the specified key from the .env file.
// If the key is not set, it panics with an error message.
func GetString(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("environment variable %s not set", key)
	}
	return val
}

// GetOptionalString retrieves the value associated with the specified key from the .env file.
// If the key is not set, it returns an empty string.
func GetOptionalString(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// GetBytes retrieves the value associated with the specified key from the .env file,
// converting it to an array of bytes. If the key is not set, it panics with an error message.
func GetBytes(key string) []byte {
	return []byte(GetString(key))
}

// GetOptionalBytes retrieves the value associated with the specified key from the .env file,
// converting it to an array of bytes. If the key is not set, it panics with an error message.
func GetOptionalBytes(key string, defaultVal []byte) []byte {
	val := GetOptionalString(key, "")
	if val == "" {
		return defaultVal
	}
	return []byte(GetString(key))
}

// GetInt retrieves the value associated with the specified key from the .env file,
// converting it to an integer. If the key is not set or the value cannot be parsed to an integer, it panics.
func GetInt(key string) int {
	val := GetString(key)
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("environment variable %s is not an int", key)
	}
	return i
}

// GetOptionalInt retrieves the value associated with the specified key from the .env file,
// converting it to an integer. If the key is not set or the value cannot be parsed to an integer, it returns 0.
func GetOptionalInt(key string, defaultVal int) int {
	val := GetOptionalString(key, "")
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

// GetDuration retrieves the value associated with the specified key from the .env file,
// converting it to a time.Duration. If the key is not set or the value cannot be parsed to a duration, it panics.
func GetDuration(key string) time.Duration {
	val := GetString(key)
	d, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("environment variable %s is not a duration", key)
	}
	return d
}

// GetOptionalDuration retrieves the value associated with the specified key from the .env file,
// converting it to a time.Duration. If the key is not set or the value cannot be parsed to a duration, it returns 0.
func GetOptionalDuration(key string, defaultVal time.Duration) time.Duration {
	val := GetOptionalString(key, "")
	if val == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
}

// GetFloat64 retrieves the value associated with the specified key from the .env file,
// converting it to a float64. If the key is not set or the value cannot be parsed to a float64, it panics.
func GetFloat64(key string) float64 {
	val := GetString(key)
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Fatalf("environment variable %s is not a float64", key)
	}
	return f
}

// GetOptionalFloat64 retrieves the value associated with the specified key from the .env file,
// converting it to a float64. If the key is not set or the value cannot be parsed to a float64, it returns 0.
func GetOptionalFloat64(key string, defaultVal float64) float64 {
	val := GetOptionalString(key, "")
	if val == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultVal
	}
	return f
}
