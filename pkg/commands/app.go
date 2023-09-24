package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"github.com/khulnasoft-lab/vul-db/pkg/db"
	dbTypes "github.com/khulnasoft-lab/vul-db/pkg/types"
	"github.com/khulnasoft-lab/vul/pkg/commands/artifact"
	"github.com/khulnasoft-lab/vul/pkg/commands/client"
	"github.com/khulnasoft-lab/vul/pkg/commands/plugin"
	"github.com/khulnasoft-lab/vul/pkg/commands/server"
	tdb "github.com/khulnasoft-lab/vul/pkg/db"
	"github.com/khulnasoft-lab/vul/pkg/result"
	"github.com/khulnasoft-lab/vul/pkg/types"
	"github.com/khulnasoft-lab/vul/pkg/utils"
)

// VersionInfo holds the vul DB version Info
type VersionInfo struct {
	Version         string       `json:",omitempty"`
	VulnerabilityDB *db.Metadata `json:",omitempty"`
}

var (
	templateFlag = cli.StringFlag{
		Name:    "template",
		Aliases: []string{"t"},
		Value:   "",
		Usage:   "output template",
		EnvVars: []string{"VUL_TEMPLATE"},
	}

	formatFlag = cli.StringFlag{
		Name:    "format",
		Aliases: []string{"f"},
		Value:   "table",
		Usage:   "format (table, json, template)",
		EnvVars: []string{"VUL_FORMAT"},
	}

	inputFlag = cli.StringFlag{
		Name:    "input",
		Aliases: []string{"i"},
		Value:   "",
		Usage:   "input file path instead of image name",
		EnvVars: []string{"VUL_INPUT"},
	}

	severityFlag = cli.StringFlag{
		Name:    "severity",
		Aliases: []string{"s"},
		Value:   strings.Join(dbTypes.SeverityNames, ","),
		Usage:   "severities of vulnerabilities to be displayed (comma separated)",
		EnvVars: []string{"VUL_SEVERITY"},
	}

	outputFlag = cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "output file name",
		EnvVars: []string{"VUL_OUTPUT"},
	}

	exitCodeFlag = cli.IntFlag{
		Name:    "exit-code",
		Usage:   "Exit code when vulnerabilities were found",
		Value:   0,
		EnvVars: []string{"VUL_EXIT_CODE"},
	}

	skipDBUpdateFlag = cli.BoolFlag{
		Name:    "skip-db-update",
		Aliases: []string{"skip-update"},
		Usage:   "skip updating vulnerability database",
		EnvVars: []string{"VUL_SKIP_UPDATE", "VUL_SKIP_DB_UPDATE"},
	}

	skipPolicyUpdateFlag = cli.BoolFlag{
		Name:    "skip-policy-update",
		Usage:   "skip updating built-in policies",
		EnvVars: []string{"VUL_SKIP_POLICY_UPDATE"},
	}

	downloadDBOnlyFlag = cli.BoolFlag{
		Name:    "download-db-only",
		Usage:   "download/update vulnerability database but don't run a scan",
		EnvVars: []string{"VUL_DOWNLOAD_DB_ONLY"},
	}

	resetFlag = cli.BoolFlag{
		Name:    "reset",
		Usage:   "remove all caches and database",
		EnvVars: []string{"VUL_RESET"},
	}

	clearCacheFlag = cli.BoolFlag{
		Name:    "clear-cache",
		Aliases: []string{"c"},
		Usage:   "clear image caches without scanning",
		EnvVars: []string{"VUL_CLEAR_CACHE"},
	}

	quietFlag = cli.BoolFlag{
		Name:    "quiet",
		Aliases: []string{"q"},
		Usage:   "suppress progress bar and log output",
		EnvVars: []string{"VUL_QUIET"},
	}

	noProgressFlag = cli.BoolFlag{
		Name:    "no-progress",
		Usage:   "suppress progress bar",
		EnvVars: []string{"VUL_NO_PROGRESS"},
	}

	ignoreUnfixedFlag = cli.BoolFlag{
		Name:    "ignore-unfixed",
		Usage:   "display only fixed vulnerabilities",
		EnvVars: []string{"VUL_IGNORE_UNFIXED"},
	}

	debugFlag = cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"d"},
		Usage:   "debug mode",
		EnvVars: []string{"VUL_DEBUG"},
	}

	removedPkgsFlag = cli.BoolFlag{
		Name:    "removed-pkgs",
		Usage:   "detect vulnerabilities of removed packages (only for Alpine)",
		EnvVars: []string{"VUL_REMOVED_PKGS"},
	}

	vulnTypeFlag = cli.StringFlag{
		Name:    "vuln-type",
		Value:   strings.Join([]string{types.VulnTypeOS, types.VulnTypeLibrary}, ","),
		Usage:   "comma-separated list of vulnerability types (os,library)",
		EnvVars: []string{"VUL_VULN_TYPE"},
	}

	securityChecksFlag = cli.StringFlag{
		Name:    "security-checks",
		Value:   types.SecurityCheckVulnerability,
		Usage:   "comma-separated list of what security issues to detect (vuln,config)",
		EnvVars: []string{"VUL_SECURITY_CHECKS"},
	}

	cacheDirFlag = cli.StringFlag{
		Name:    "cache-dir",
		Value:   utils.DefaultCacheDir(),
		Usage:   "cache directory",
		EnvVars: []string{"VUL_CACHE_DIR"},
	}

	cacheBackendFlag = cli.StringFlag{
		Name:    "cache-backend",
		Value:   "fs",
		Usage:   "cache backend (e.g. redis://localhost:6379)",
		EnvVars: []string{"VUL_CACHE_BACKEND"},
	}

	ignoreFileFlag = cli.StringFlag{
		Name:    "ignorefile",
		Value:   result.DefaultIgnoreFile,
		Usage:   "specify .vulignore file",
		EnvVars: []string{"VUL_IGNOREFILE"},
	}

	timeoutFlag = cli.DurationFlag{
		Name:    "timeout",
		Value:   time.Second * 300,
		Usage:   "timeout",
		EnvVars: []string{"VUL_TIMEOUT"},
	}

	lightFlag = cli.BoolFlag{
		Name:    "light",
		Usage:   "light mode: it's faster, but vulnerability descriptions and references are not displayed",
		EnvVars: []string{"VUL_LIGHT"},
	}

	token = cli.StringFlag{
		Name:    "token",
		Usage:   "for authentication",
		EnvVars: []string{"VUL_TOKEN"},
	}

	tokenHeader = cli.StringFlag{
		Name:    "token-header",
		Value:   "Vul-Token",
		Usage:   "specify a header name for token",
		EnvVars: []string{"VUL_TOKEN_HEADER"},
	}

	ignorePolicy = cli.StringFlag{
		Name:    "ignore-policy",
		Usage:   "specify the Rego file to evaluate each vulnerability",
		EnvVars: []string{"VUL_IGNORE_POLICY"},
	}

	listAllPackages = cli.BoolFlag{
		Name:    "list-all-pkgs",
		Usage:   "enabling the option will output all packages regardless of vulnerability",
		EnvVars: []string{"VUL_LIST_ALL_PKGS"},
	}

	skipFiles = cli.StringSliceFlag{
		Name:    "skip-files",
		Usage:   "specify the file paths to skip traversal",
		EnvVars: []string{"VUL_SKIP_FILES"},
	}

	skipDirs = cli.StringSliceFlag{
		Name:    "skip-dirs",
		Usage:   "specify the directories where the traversal is skipped",
		EnvVars: []string{"VUL_SKIP_DIRS"},
	}

	// For misconfigurations
	configPolicy = cli.StringSliceFlag{
		Name:    "config-policy",
		Usage:   "specify paths to the Rego policy files directory, applying config files",
		EnvVars: []string{"VUL_CONFIG_POLICY"},
	}

	configPolicyAlias = cli.StringSliceFlag{
		Name:    "policy",
		Aliases: []string{"config-policy"},
		Usage:   "specify paths to the Rego policy files directory, applying config files",
		EnvVars: []string{"VUL_POLICY"},
	}

	configData = cli.StringSliceFlag{
		Name:    "config-data",
		Usage:   "specify paths from which data for the Rego policies will be recursively loaded",
		EnvVars: []string{"VUL_CONFIG_DATA"},
	}

	configDataAlias = cli.StringSliceFlag{
		Name:    "data",
		Aliases: []string{"config-data"},
		Usage:   "specify paths from which data for the Rego policies will be recursively loaded",
		EnvVars: []string{"VUL_DATA"},
	}

	filePatterns = cli.StringSliceFlag{
		Name:    "file-patterns",
		Usage:   "specify file patterns",
		EnvVars: []string{"VUL_FILE_PATTERNS"},
	}

	policyNamespaces = cli.StringSliceFlag{
		Name:    "policy-namespaces",
		Aliases: []string{"namespaces"},
		Usage:   "Rego namespaces",
		Value:   cli.NewStringSlice("users"),
		EnvVars: []string{"VUL_POLICY_NAMESPACES"},
	}

	includeNonFailures = cli.BoolFlag{
		Name:    "include-non-failures",
		Usage:   "include successes and exceptions",
		Value:   false,
		EnvVars: []string{"VUL_INCLUDE_NON_FAILURES"},
	}

	traceFlag = cli.BoolFlag{
		Name:    "trace",
		Usage:   "enable more verbose trace output for custom queries",
		Value:   false,
		EnvVars: []string{"VUL_TRACE"},
	}

	// Global flags
	globalFlags = []cli.Flag{
		&quietFlag,
		&debugFlag,
		&cacheDirFlag,
	}

	imageFlags = []cli.Flag{
		&templateFlag,
		&formatFlag,
		&inputFlag,
		&severityFlag,
		&outputFlag,
		&exitCodeFlag,
		&skipDBUpdateFlag,
		&downloadDBOnlyFlag,
		&resetFlag,
		&clearCacheFlag,
		&noProgressFlag,
		&ignoreUnfixedFlag,
		&removedPkgsFlag,
		&vulnTypeFlag,
		&securityChecksFlag,
		&ignoreFileFlag,
		&timeoutFlag,
		&lightFlag,
		&ignorePolicy,
		&listAllPackages,
		&cacheBackendFlag,
		stringSliceFlag(skipFiles),
		stringSliceFlag(skipDirs),
	}
)

// NewApp is the factory method to return Vul CLI
func NewApp(version string) *cli.App {
	cli.VersionPrinter = func(c *cli.Context) {
		showVersion(c.String("cache-dir"), c.String("format"), c.App.Version, c.App.Writer)
	}

	app := cli.NewApp()
	app.Name = "vul"
	app.Version = version
	app.ArgsUsage = "target"
	app.Usage = "A simple and comprehensive vulnerability scanner for containers"
	app.EnableBashCompletion = true

	flags := append(globalFlags, setHidden(imageFlags, true)...)

	app.Flags = flags
	app.Commands = []*cli.Command{
		NewImageCommand(),
		NewFilesystemCommand(),
		NewRootfsCommand(),
		NewRepositoryCommand(),
		NewClientCommand(),
		NewServerCommand(),
		NewConfigCommand(),
		NewPluginCommand(),
	}
	app.Commands = append(app.Commands, plugin.LoadCommands()...)

	runAsPlugin := os.Getenv("VUL_RUN_AS_PLUGIN")
	if runAsPlugin == "" {
		app.Action = artifact.ImageRun
	} else {
		app.Action = func(ctx *cli.Context) error {
			return plugin.RunWithArgs(ctx.Context, runAsPlugin, ctx.Args().Slice())
		}
	}
	return app
}

func setHidden(flags []cli.Flag, hidden bool) []cli.Flag {
	var newFlags []cli.Flag
	for _, flag := range flags {
		var f cli.Flag
		switch pf := flag.(type) {
		case *cli.StringFlag:
			stringFlag := *pf
			stringFlag.Hidden = hidden
			f = &stringFlag
		case *cli.StringSliceFlag:
			stringSliceFlag := *pf
			stringSliceFlag.Hidden = hidden
			f = &stringSliceFlag
		case *cli.BoolFlag:
			boolFlag := *pf
			boolFlag.Hidden = hidden
			f = &boolFlag
		case *cli.IntFlag:
			intFlag := *pf
			intFlag.Hidden = hidden
			f = &intFlag
		case *cli.DurationFlag:
			durationFlag := *pf
			durationFlag.Hidden = hidden
			f = &durationFlag
		}
		newFlags = append(newFlags, f)
	}
	return newFlags
}

func showVersion(cacheDir, outputFormat, version string, outputWriter io.Writer) {
	var dbMeta *db.Metadata

	metadata, _ := tdb.NewMetadata(afero.NewOsFs(), cacheDir).Get() // nolint: errcheck
	if !metadata.UpdatedAt.IsZero() && !metadata.NextUpdate.IsZero() && metadata.Version != 0 {
		dbMeta = &db.Metadata{
			Version:      metadata.Version,
			Type:         metadata.Type,
			NextUpdate:   metadata.NextUpdate.UTC(),
			UpdatedAt:    metadata.UpdatedAt.UTC(),
			DownloadedAt: metadata.DownloadedAt.UTC(),
		}
	}

	switch outputFormat {
	case "json":
		b, _ := json.Marshal(VersionInfo{ // nolint: errcheck
			Version:         version,
			VulnerabilityDB: dbMeta,
		})
		fmt.Fprintln(outputWriter, string(b))
	default:
		output := fmt.Sprintf("Version: %s\n", version)
		if dbMeta != nil {
			var dbType string
			switch dbMeta.Type {
			case db.TypeFull:
				dbType = "Full"
			case db.TypeLight:
				dbType = "Light"
			}
			output += fmt.Sprintf(`Vulnerability DB:
  Type: %s
  Version: %d
  UpdatedAt: %s
  NextUpdate: %s
  DownloadedAt: %s
`, dbType, dbMeta.Version, dbMeta.UpdatedAt.UTC(), dbMeta.NextUpdate.UTC(), dbMeta.DownloadedAt.UTC())
		}
		fmt.Fprintf(outputWriter, output)
	}
}

// NewImageCommand is the factory method to add image command
func NewImageCommand() *cli.Command {
	return &cli.Command{
		Name:      "image",
		Aliases:   []string{"i"},
		ArgsUsage: "image_name",
		Usage:     "scan an image",
		Action:    artifact.ImageRun,
		Flags:     imageFlags,
	}
}

// NewFilesystemCommand is the factory method to add filesystem command
func NewFilesystemCommand() *cli.Command {
	return &cli.Command{
		Name:      "filesystem",
		Aliases:   []string{"fs"},
		ArgsUsage: "dir",
		Usage:     "scan local filesystem for language-specific dependencies and config files",
		Action:    artifact.FilesystemRun,
		Flags: []cli.Flag{
			&templateFlag,
			&formatFlag,
			&severityFlag,
			&outputFlag,
			&exitCodeFlag,
			&skipDBUpdateFlag,
			&skipPolicyUpdateFlag,
			&clearCacheFlag,
			&ignoreUnfixedFlag,
			&vulnTypeFlag,
			&securityChecksFlag,
			&ignoreFileFlag,
			&cacheBackendFlag,
			&timeoutFlag,
			&noProgressFlag,
			&ignorePolicy,
			&listAllPackages,
			stringSliceFlag(skipFiles),
			stringSliceFlag(skipDirs),
			stringSliceFlag(configPolicy),
			stringSliceFlag(configData),
			stringSliceFlag(policyNamespaces),
		},
	}
}

// NewRootfsCommand is the factory method to add filesystem command
func NewRootfsCommand() *cli.Command {
	return &cli.Command{
		Name:      "rootfs",
		ArgsUsage: "dir",
		Usage:     "scan rootfs",
		Action:    artifact.RootfsRun,
		Flags: []cli.Flag{
			&templateFlag,
			&formatFlag,
			&severityFlag,
			&outputFlag,
			&exitCodeFlag,
			&skipDBUpdateFlag,
			&skipPolicyUpdateFlag,
			&clearCacheFlag,
			&ignoreUnfixedFlag,
			&vulnTypeFlag,
			&securityChecksFlag,
			&ignoreFileFlag,
			&cacheBackendFlag,
			&timeoutFlag,
			&noProgressFlag,
			&ignorePolicy,
			&listAllPackages,
			stringSliceFlag(skipFiles),
			stringSliceFlag(skipDirs),
			stringSliceFlag(configPolicy),
			stringSliceFlag(configData),
			stringSliceFlag(policyNamespaces),
		},
	}
}

// NewRepositoryCommand is the factory method to add repository command
func NewRepositoryCommand() *cli.Command {
	return &cli.Command{
		Name:      "repository",
		Aliases:   []string{"repo"},
		ArgsUsage: "repo_url",
		Usage:     "scan remote repository",
		Action:    artifact.RepositoryRun,
		Flags: []cli.Flag{
			&templateFlag,
			&formatFlag,
			&inputFlag,
			&severityFlag,
			&outputFlag,
			&exitCodeFlag,
			&skipDBUpdateFlag,
			&skipPolicyUpdateFlag,
			&clearCacheFlag,
			&ignoreUnfixedFlag,
			&removedPkgsFlag,
			&vulnTypeFlag,
			&securityChecksFlag,
			&ignoreFileFlag,
			&cacheBackendFlag,
			&timeoutFlag,
			&noProgressFlag,
			&ignorePolicy,
			&listAllPackages,
			stringSliceFlag(skipFiles),
			stringSliceFlag(skipDirs),
		},
	}
}

// NewClientCommand is the factory method to add client command
func NewClientCommand() *cli.Command {
	return &cli.Command{
		Name:      "client",
		Aliases:   []string{"c"},
		ArgsUsage: "image_name",
		Usage:     "client mode",
		Action:    client.Run,
		Flags: []cli.Flag{
			&templateFlag,
			&formatFlag,
			&inputFlag,
			&severityFlag,
			&outputFlag,
			&exitCodeFlag,
			&clearCacheFlag,
			&ignoreUnfixedFlag,
			&removedPkgsFlag,
			&vulnTypeFlag,
			&securityChecksFlag,
			&ignoreFileFlag,
			&timeoutFlag,
			&ignorePolicy,
			stringSliceFlag(skipFiles),
			stringSliceFlag(skipDirs),
			stringSliceFlag(configPolicy),
			&listAllPackages,

			// original flags
			&token,
			&tokenHeader,
			&cli.StringFlag{
				Name:    "remote",
				Value:   "http://localhost:4954",
				Usage:   "server address",
				EnvVars: []string{"VUL_REMOTE"},
			},
			&cli.StringSliceFlag{
				Name:    "custom-headers",
				Usage:   "custom headers",
				EnvVars: []string{"VUL_CUSTOM_HEADERS"},
			},
		},
	}
}

// NewServerCommand is the factory method to add server command
func NewServerCommand() *cli.Command {
	return &cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "server mode",
		Action:  server.Run,
		Flags: []cli.Flag{
			&skipDBUpdateFlag,
			&downloadDBOnlyFlag,
			&resetFlag,
			&cacheBackendFlag,

			// original flags
			&token,
			&tokenHeader,
			&cli.StringFlag{
				Name:    "listen",
				Value:   "localhost:4954",
				Usage:   "listen address",
				EnvVars: []string{"VUL_LISTEN"},
			},
		},
	}
}

// NewConfigCommand adds config command
func NewConfigCommand() *cli.Command {
	return &cli.Command{
		Name:      "config",
		Aliases:   []string{"conf"},
		ArgsUsage: "dir",
		Usage:     "scan config files",
		Action:    artifact.ConfigRun,
		Flags: []cli.Flag{
			&templateFlag,
			&formatFlag,
			&severityFlag,
			&outputFlag,
			&exitCodeFlag,
			&skipPolicyUpdateFlag,
			&resetFlag,
			&clearCacheFlag,
			&ignoreFileFlag,
			&timeoutFlag,
			stringSliceFlag(skipFiles),
			stringSliceFlag(skipDirs),
			stringSliceFlag(configPolicyAlias),
			stringSliceFlag(configDataAlias),
			stringSliceFlag(policyNamespaces),
			stringSliceFlag(filePatterns),
			&includeNonFailures,
			&traceFlag,
		},
	}
}

// NewPluginCommand is the factory method to add plugin command
func NewPluginCommand() *cli.Command {
	return &cli.Command{
		Name:      "plugin",
		Aliases:   []string{"p"},
		ArgsUsage: "plugin_uri",
		Usage:     "manage plugins",
		Subcommands: cli.Commands{
			{
				Name:      "install",
				Aliases:   []string{"i"},
				Usage:     "install a plugin",
				ArgsUsage: "URL | FILE_PATH",
				Action:    plugin.Install,
			},
			{
				Name:      "uninstall",
				Aliases:   []string{"u"},
				Usage:     "uninstall a plugin",
				ArgsUsage: "PLUGIN_NAME",
				Action:    plugin.Uninstall,
			},
			{
				Name:      "run",
				Aliases:   []string{"r"},
				Usage:     "run a plugin on the fly",
				ArgsUsage: "PLUGIN_NAME [PLUGIN_OPTIONS]",
				Action:    plugin.Run,
			},
		},
	}
}

// StringSliceFlag is defined globally. When the app runs multiple times,
// the previous value will be retained and it causes unexpected results.
// The flag value is copied through this function to prevent the issue.
func stringSliceFlag(f cli.StringSliceFlag) *cli.StringSliceFlag {
	return &f
}
