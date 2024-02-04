package config

import (
	"log"
	"os"
)

var TMP_DIR = "tmp"

func init() {
	err := os.MkdirAll(TMP_DIR, 0777)
	if err != nil {
		log.Fatalln(err)
	}
}
