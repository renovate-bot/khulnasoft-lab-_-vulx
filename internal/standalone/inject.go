//go:build wireinject
// +build wireinject

package standalone

import (
	"context"
	"time"

	"github.com/aquasecurity/fanal/cache"
	"github.com/google/wire"
	"github.com/khulnasoft-lab/vul/pkg/scanner"
	"github.com/khulnasoft-lab/vul/pkg/vulnerability"
)

func initializeDockerScanner(ctx context.Context, imageName string, layerCache cache.ImageCache, localImageCache cache.LocalImageCache,
	timeout time.Duration) (scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneDockerSet)
	return scanner.Scanner{}, nil, nil
}

func initializeArchiveScanner(ctx context.Context, filePath string, layerCache cache.ImageCache, localImageCache cache.LocalImageCache,
	timeout time.Duration) (scanner.Scanner, func(), error) {
	wire.Build(scanner.StandaloneArchiveSet)
	return scanner.Scanner{}, nil, nil
}

func initializeVulnerabilityClient() vulnerability.Client {
	wire.Build(vulnerability.SuperSet)
	return vulnerability.Client{}
}
