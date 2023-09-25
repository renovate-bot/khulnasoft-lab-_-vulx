package client

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	"github.com/aquasecurity/fanal/analyzer"
	"github.com/aquasecurity/fanal/analyzer/config"
	"github.com/aquasecurity/fanal/artifact"
	"github.com/khulnasoft-lab/vul/pkg/cache"
	"github.com/khulnasoft-lab/vul/pkg/commands/operation"
	"github.com/khulnasoft-lab/vul/pkg/log"
	pkgReport "github.com/khulnasoft-lab/vul/pkg/report"
	"github.com/khulnasoft-lab/vul/pkg/rpc/client"
	"github.com/khulnasoft-lab/vul/pkg/scanner"
	"github.com/khulnasoft-lab/vul/pkg/types"
	"github.com/khulnasoft-lab/vul/pkg/utils"
)

const defaultPolicyNamespace = "appshield"

// Run runs the scan
func Run(cliCtx *cli.Context) error {
	opt, err := NewOption(cliCtx)
	if err != nil {
		return xerrors.Errorf("option error: %w", err)
	}

	ctx, cancel := context.WithTimeout(cliCtx.Context, opt.Timeout)
	defer cancel()

	// Disable the lock file scanning
	opt.DisabledAnalyzers = analyzer.TypeLockfiles

	err = runWithTimeout(ctx, opt)
	if xerrors.Is(err, context.DeadlineExceeded) {
		log.Logger.Warn("Increase --timeout value")
	}
	return err
}

func runWithTimeout(ctx context.Context, opt Option) error {
	if err := initialize(&opt); err != nil {
		return xerrors.Errorf("initialize error: %w", err)
	}

	if opt.ClearCache {
		log.Logger.Warn("A client doesn't have image cache")
		return nil
	}

	s, cleanup, err := initializeScanner(ctx, opt)
	if err != nil {
		return xerrors.Errorf("scanner initialize error: %w", err)
	}
	defer cleanup()

	scanOptions := types.ScanOptions{
		VulnType:            opt.VulnType,
		SecurityChecks:      opt.SecurityChecks,
		ScanRemovedPackages: opt.ScanRemovedPkgs,
		ListAllPackages:     opt.ListAllPkgs,
	}
	log.Logger.Debugf("Vulnerability type:  %s", scanOptions.VulnType)

	report, err := s.ScanArtifact(ctx, scanOptions)
	if err != nil {
		return xerrors.Errorf("error in image scan: %w", err)
	}

	resultClient := initializeResultClient()
	results := report.Results
	for i := range results {
		vulns, misconfSummary, misconfs, err := resultClient.Filter(ctx, results[i].Vulnerabilities, results[i].Misconfigurations,
			opt.Severities, opt.IgnoreUnfixed, opt.IncludeNonFailures, opt.IgnoreFile, opt.IgnorePolicy)
		if err != nil {
			return xerrors.Errorf("filter error: %w", err)
		}
		results[i].Vulnerabilities = vulns
		results[i].Misconfigurations = misconfs
		results[i].MisconfSummary = misconfSummary
	}

	if err = pkgReport.Write(report, pkgReport.Option{
		Format:             opt.Format,
		Output:             opt.Output,
		Severities:         opt.Severities,
		OutputTemplate:     opt.Template,
		Light:              false,
		IncludeNonFailures: opt.IncludeNonFailures,
		Trace:              opt.Trace,
	}); err != nil {
		return xerrors.Errorf("unable to write results: %w", err)
	}

	exit(opt, results)

	return nil
}

func initialize(opt *Option) error {
	// Initialize logger
	if err := log.InitLogger(opt.Debug, opt.Quiet); err != nil {
		return xerrors.Errorf("failed to initialize a logger: %w", err)
	}

	// Initialize options
	if err := opt.Init(); err != nil {
		return xerrors.Errorf("failed to initialize options: %w", err)
	}

	// configure cache dir
	utils.SetCacheDir(opt.CacheDir)
	log.Logger.Debugf("cache dir:  %s", utils.CacheDir())

	return nil
}

func disabledAnalyzers(opt Option) []analyzer.Type {
	// Specified analyzers to be disabled depending on scanning modes
	// e.g. The 'image' subcommand should disable the lock file scanning.
	analyzers := opt.DisabledAnalyzers

	// It doesn't analyze apk commands by default.
	if !opt.ScanRemovedPkgs {
		analyzers = append(analyzers, analyzer.TypeApkCommand)
	}

	// Don't analyze programming language packages when not running in 'library' mode
	if !utils.StringInSlice(types.VulnTypeLibrary, opt.VulnType) {
		analyzers = append(analyzers, analyzer.TypeLanguages...)
	}

	return analyzers
}

func initializeScanner(ctx context.Context, opt Option) (scanner.Scanner, func(), error) {
	remoteCache := cache.NewRemoteCache(cache.RemoteURL(opt.RemoteAddr), opt.CustomHeaders)

	// ScannerOptions is filled only when config scanning is enabled.
	var configScannerOptions config.ScannerOption
	if utils.StringInSlice(types.SecurityCheckConfig, opt.SecurityChecks) {
		builtinPolicyPaths, err := operation.InitBuiltinPolicies(ctx, false)
		if err != nil {
			return scanner.Scanner{}, nil, xerrors.Errorf("failed to initialize default policies: %w", err)
		}

		configScannerOptions = config.ScannerOption{
			Trace:        opt.Trace,
			Namespaces:   append(opt.PolicyNamespaces, defaultPolicyNamespace),
			PolicyPaths:  append(opt.PolicyPaths, builtinPolicyPaths...),
			DataPaths:    opt.DataPaths,
			FilePatterns: opt.FilePatterns,
		}
	}

	artifactOpt := artifact.Option{
		DisabledAnalyzers: disabledAnalyzers(opt),
		SkipFiles:         opt.SkipFiles,
		SkipDirs:          opt.SkipDirs,
	}

	if opt.Input != "" {
		// Scan tar file
		s, err := initializeArchiveScanner(ctx, opt.Input, remoteCache, client.CustomHeaders(opt.CustomHeaders),
			client.RemoteURL(opt.RemoteAddr), opt.Timeout, artifactOpt, configScannerOptions)
		if err != nil {
			return scanner.Scanner{}, nil, xerrors.Errorf("unable to initialize the archive scanner: %w", err)
		}
		return s, func() {}, nil
	}

	// Scan an image in Docker Engine or Docker Registry
	s, cleanup, err := initializeDockerScanner(ctx, opt.Target, remoteCache, client.CustomHeaders(opt.CustomHeaders),
		client.RemoteURL(opt.RemoteAddr), opt.Timeout, artifactOpt, configScannerOptions)
	if err != nil {
		return scanner.Scanner{}, nil, xerrors.Errorf("unable to initialize the docker scanner: %w", err)
	}

	return s, cleanup, nil
}

func exit(c Option, results pkgReport.Results) {
	if c.ExitCode != 0 {
		for _, result := range results {
			if len(result.Vulnerabilities) > 0 {
				os.Exit(c.ExitCode)
			}
		}
	}
}
