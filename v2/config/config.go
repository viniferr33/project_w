package config

import (
	"log"
	"os"
	"time"
)

var TMP_DIR = "tmp"
var GCS_TIMEOUT_S = time.Second * 50
var GCS_BUCKET = "project_w_audio"

func init() {
	err := os.MkdirAll(TMP_DIR, 0777)
	if err != nil {
		log.Fatalln(err)
	}
}
