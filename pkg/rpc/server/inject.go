//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"
	"github.com/khulnasoft-lab/fanal/cache"
)

func initializeScanServer(localArtifactCache cache.LocalArtifactCache) *ScanServer {
	wire.Build(ScanSuperSet)
	return &ScanServer{}
}

func initializeDBWorker(cacheDir string, quiet bool) dbWorker {
	wire.Build(DBWorkerSuperSet)
	return dbWorker{}
}
