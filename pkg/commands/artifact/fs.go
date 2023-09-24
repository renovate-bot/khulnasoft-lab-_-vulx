package artifact

import (
	"context"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	"github.com/khulnasoft-lab/fanal/analyzer"
	"github.com/khulnasoft-lab/fanal/analyzer/config"
	"github.com/khulnasoft-lab/fanal/artifact"
	"github.com/khulnasoft-lab/fanal/cache"
	"github.com/khulnasoft-lab/vul/pkg/scanner"
)

func filesystemScanner(ctx context.Context, dir string, ac cache.ArtifactCache, lac cache.LocalArtifactCache,
	_ time.Duration, artifactOpt artifact.Option, scannerOpt config.ScannerOption) (scanner.Scanner, func(), error) {
	s, cleanup, err := initializeFilesystemScanner(ctx, dir, ac, lac, artifactOpt, scannerOpt)
	if err != nil {
		return scanner.Scanner{}, func() {}, xerrors.Errorf("unable to initialize a filesystem scanner: %w", err)
	}
	return s, cleanup, nil
}

// FilesystemRun runs scan on filesystem for language-specific dependencies and config files
func FilesystemRun(ctx *cli.Context) error {
	opt, err := initOption(ctx)
	if err != nil {
		return xerrors.Errorf("option error: %w", err)
	}

	// Disable the individual package scanning
	opt.DisabledAnalyzers = analyzer.TypeIndividualPkgs

	return Run(ctx.Context, opt, filesystemScanner, initFSCache)
}

// RootfsRun runs scan on rootfs.
func RootfsRun(ctx *cli.Context) error {
	opt, err := initOption(ctx)
	if err != nil {
		return xerrors.Errorf("option error: %w", err)
	}

	// Disable the lock file scanning
	opt.DisabledAnalyzers = analyzer.TypeLockfiles

	return Run(ctx.Context, opt, filesystemScanner, initFSCache)
}
