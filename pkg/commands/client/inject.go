//go:build wireinject
// +build wireinject

package client

import (
	"context"
	"time"

	"github.com/google/wire"

	"github.com/aquasecurity/fanal/analyzer/config"
	"github.com/aquasecurity/fanal/artifact"
	"github.com/aquasecurity/fanal/cache"
	"github.com/khulnasoft-lab/vul/pkg/result"
	"github.com/khulnasoft-lab/vul/pkg/rpc/client"
	"github.com/khulnasoft-lab/vul/pkg/scanner"
)

func initializeDockerScanner(ctx context.Context, imageName string, artifactCache cache.ArtifactCache, customHeaders client.CustomHeaders,
	url client.RemoteURL, timeout time.Duration, artifactOption artifact.Option, configScannerOption config.ScannerOption) (
	scanner.Scanner, func(), error) {
	wire.Build(scanner.RemoteDockerSet)
	return scanner.Scanner{}, nil, nil
}

func initializeArchiveScanner(ctx context.Context, filePath string, artifactCache cache.ArtifactCache,
	customHeaders client.CustomHeaders, url client.RemoteURL, timeout time.Duration, artifactOption artifact.Option,
	configScannerOption config.ScannerOption) (scanner.Scanner, error) {
	wire.Build(scanner.RemoteArchiveSet)
	return scanner.Scanner{}, nil
}

func initializeResultClient() result.Client {
	wire.Build(result.SuperSet)
	return result.Client{}
}
