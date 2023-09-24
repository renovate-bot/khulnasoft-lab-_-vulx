//go:build wireinject
// +build wireinject

package operation

import (
	"github.com/google/wire"
	"github.com/khulnasoft-lab/vul/pkg/db"
)

func initializeDBClient(cacheDir string, quiet bool) db.Client {
	wire.Build(db.SuperSet)
	return db.Client{}
}
