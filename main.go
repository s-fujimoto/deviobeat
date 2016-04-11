package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/s-fujimoto/deviobeat/beater"
)

func main() {
	err := beat.Run("deviobeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
