package main

import (
	"flag"

	"github.com/benmizrahi/gobig/internal/common"
	"github.com/benmizrahi/gobig/internal/master"
	log "github.com/sirupsen/logrus"
)

func main() {

	host := flag.String("host", "localhost", "# Host to listen")
	port := flag.Int("port", 9999, "#port to listen")
	isLocal := flag.Bool("isLocal", true, "# Run locally with processes K8S/Local")
	flag.Parse()

	log.Info("gobig Status: Starting")
	log.Info("******************")
	log.Info("Status: Ready!")
	log.Info("******************")

	data := [][]string{
		{"HELLO"},
		{"WORLD"},
		{"GOLANG"},
		{"IS"},
		{"THE"},
		{"BEST"},
	}
	// word count example
	master.
		NewMaster(*isLocal, *host, *port, 2).
		Parallelize(data, common.Options{}).
		Count()
}
