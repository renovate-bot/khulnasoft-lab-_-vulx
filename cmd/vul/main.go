package main

import (
	l "log"
	"os"

	"github.com/khulnasoft-lab/vul/internal"

	"github.com/khulnasoft-lab/vul/pkg/log"
)

var (
	version = "dev"
)

func main() {
	app := internal.NewApp(version)
	err := app.Run(os.Args)
	if err != nil {
		if log.Logger != nil {
			log.Fatal(err)
		}
		l.Fatal(err)
	}
}
