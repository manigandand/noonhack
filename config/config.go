package config

import (
	"log"
	"os"
)

// Port ...
var (
	Port         string
	Home         string
	QueueFileDir string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	Port = "9090"
	Home = home
	QueueFileDir = home + "/.noonhack"

	if _, err := os.Stat(QueueFileDir); os.IsNotExist(err) {
		if err := os.Mkdir(QueueFileDir, os.ModePerm); err != nil {
			log.Fatal("Can't create config dir ", QueueFileDir, err)
		}
	}
}
