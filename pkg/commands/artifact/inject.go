//go:build wireinject
// +build wireinject

package artifact

import (
	"context"
	"time"

	"github.com/google/wire"

	"github.com/khulnasoft-lab/fanal/analyzer/config"
	"github.com/khulnasoft-lab/fanal/artifact"
	"github.com/khulnasoft-lab/fanal/cache"
	"github.com/khulnasoft-lab/vul/pkg/result"
	"github.com/khulnasoft-lab/vul/pkg/scanner"
)

func initializeDockerScanner(ctx context.Context, imageName string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.LocalArtifactCache, timeout time.Duration, artifactOption artifact.Option,
	configScannerOption config.ScannerOption) (scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneDockerSet)
	return scanner.Scanner{}, nil, nil
}

func initializeArchiveScanner(ctx context.Context, filePath string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.LocalArtifactCache, timeout time.Duration, artifactOption artifact.Option,
	configScannerOption config.ScannerOption) (scanner.Scanner, error) {
	wire.Build(scanner.StandaloneArchiveSet)
	return scanner.Scanner{}, nil
}

func initializeFilesystemScanner(ctx context.Context, dir string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.LocalArtifactCache, artifactOption artifact.Option,
	configScannerOption config.ScannerOption) (scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneFilesystemSet)
	return scanner.Scanner{}, nil, nil
}

func initializeRepositoryScanner(ctx context.Context, url string, artifactCache cache.ArtifactCache,
	localArtifactCache cache.LocalArtifactCache, artifactOption artifact.Option,
	configScannerOption config.ScannerOption) (scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneRepositorySet)
	return scanner.Scanner{}, nil, nil
}

func initializeResultClient() result.Client {
	wire.Build(result.SuperSet)
	return result.Client{}
}
