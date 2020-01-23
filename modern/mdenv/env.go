package mdenv

import (
	"fmt"
	"log"
	"os"

	"github.com/SentadMedia/elf/fw"
	"github.com/joho/godotenv"
)

var _ fw.Environment = (*GoDotEnv)(nil)

type GoDotEnv struct {
}

func (g GoDotEnv) GetEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	fmt.Printf("Trying to gen environement varyable for %s", key)
	if val == "" {
		return defaultValue
	}
	return val
}

func (g GoDotEnv) AutoLoadDotEnvFile() {
	_, err := os.Stat(".env")
	if os.IsNotExist(err) {
		return
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func NewGoDotEnv() GoDotEnv {
	return GoDotEnv{}
}
