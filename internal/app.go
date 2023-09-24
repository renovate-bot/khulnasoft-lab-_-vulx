package internal

import (
	"strings"
	"time"

	"github.com/urfave/cli"

	"github.com/khulnasoft-lab/vul-db/pkg/types"
	"github.com/khulnasoft-lab/vul/internal/client"
	"github.com/khulnasoft-lab/vul/internal/server"
	"github.com/khulnasoft-lab/vul/internal/standalone"
	"github.com/khulnasoft-lab/vul/pkg/utils"
	"github.com/khulnasoft-lab/vul/pkg/vulnerability"
)

var (
	templateFlag = cli.StringFlag{
		Name:   "template, t",
		Value:  "",
		Usage:  "output template",
		EnvVar: "VUL_TEMPLATE",
	}

	formatFlag = cli.StringFlag{
		Name:   "format, f",
		Value:  "table",
		Usage:  "format (table, json, template)",
		EnvVar: "VUL_FORMAT",
	}

	inputFlag = cli.StringFlag{
		Name:   "input, i",
		Value:  "",
		Usage:  "input file path instead of image name",
		EnvVar: "VUL_INPUT",
	}

	severityFlag = cli.StringFlag{
		Name:   "severity, s",
		Value:  strings.Join(types.SeverityNames, ","),
		Usage:  "severities of vulnerabilities to be displayed (comma separated)",
		EnvVar: "VUL_SEVERITY",
	}

	outputFlag = cli.StringFlag{
		Name:   "output, o",
		Usage:  "output file name",
		EnvVar: "VUL_OUTPUT",
	}

	exitCodeFlag = cli.IntFlag{
		Name:   "exit-code",
		Usage:  "Exit code when vulnerabilities were found",
		Value:  0,
		EnvVar: "VUL_EXIT_CODE",
	}

	skipUpdateFlag = cli.BoolFlag{
		Name:   "skip-update",
		Usage:  "skip db update",
		EnvVar: "VUL_SKIP_UPDATE",
	}

	downloadDBOnlyFlag = cli.BoolFlag{
		Name:   "download-db-only",
		Usage:  "download/update vulnerability database but don't run a scan",
		EnvVar: "VUL_DOWNLOAD_DB_ONLY",
	}

	resetFlag = cli.BoolFlag{
		Name:   "reset",
		Usage:  "remove all caches and database",
		EnvVar: "VUL_RESET",
	}

	clearCacheFlag = cli.BoolFlag{
		Name:   "clear-cache, c",
		Usage:  "clear image caches without scanning",
		EnvVar: "VUL_CLEAR_CACHE",
	}

	quietFlag = cli.BoolFlag{
		Name:   "quiet, q",
		Usage:  "suppress progress bar and log output",
		EnvVar: "VUL_QUIET",
	}

	noProgressFlag = cli.BoolFlag{
		Name:   "no-progress",
		Usage:  "suppress progress bar",
		EnvVar: "VUL_NO_PROGRESS",
	}

	ignoreUnfixedFlag = cli.BoolFlag{
		Name:   "ignore-unfixed",
		Usage:  "display only fixed vulnerabilities",
		EnvVar: "VUL_IGNORE_UNFIXED",
	}

	debugFlag = cli.BoolFlag{
		Name:   "debug, d",
		Usage:  "debug mode",
		EnvVar: "VUL_DEBUG",
	}

	removedPkgsFlag = cli.BoolFlag{
		Name:   "removed-pkgs",
		Usage:  "detect vulnerabilities of removed packages (only for Alpine)",
		EnvVar: "VUL_REMOVED_PKGS",
	}

	vulnTypeFlag = cli.StringFlag{
		Name:   "vuln-type",
		Value:  "os,library",
		Usage:  "comma-separated list of vulnerability types (os,library)",
		EnvVar: "VUL_VULN_TYPE",
	}

	cacheDirFlag = cli.StringFlag{
		Name:   "cache-dir",
		Value:  utils.DefaultCacheDir(),
		Usage:  "cache directory",
		EnvVar: "VUL_CACHE_DIR",
	}

	ignoreFileFlag = cli.StringFlag{
		Name:   "ignorefile",
		Value:  vulnerability.DefaultIgnoreFile,
		Usage:  "specify .vulignore file",
		EnvVar: "VUL_IGNOREFILE",
	}

	timeoutFlag = cli.DurationFlag{
		Name:   "timeout",
		Value:  time.Second * 120,
		Usage:  "docker timeout",
		EnvVar: "VUL_TIMEOUT",
	}

	lightFlag = cli.BoolFlag{
		Name:   "light",
		Usage:  "light mode: it's faster, but vulnerability descriptions and references are not displayed",
		EnvVar: "VUL_LIGHT",
	}

	token = cli.StringFlag{
		Name:   "token",
		Usage:  "for authentication",
		EnvVar: "VUL_TOKEN",
	}

	tokenHeader = cli.StringFlag{
		Name:   "token-header",
		Value:  "Vul-Token",
		Usage:  "specify a header name for token",
		EnvVar: "VUL_TOKEN_HEADER",
	}
)

func NewApp(version string) *cli.App {
	cli.AppHelpTemplate = `NAME:
  {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}
USAGE:
  {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}
VERSION:
  {{.Version}}{{end}}{{end}}{{if .Description}}
DESCRIPTION:
  {{.Description}}{{end}}{{if len .Authors}}
AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
  {{range $index, $author := .Authors}}{{if $index}}
  {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}
OPTIONS:
  {{range $index, $option := .VisibleFlags}}{{if $index}}
  {{end}}{{$option}}{{end}}{{end}}
`
	app := cli.NewApp()
	app.Name = "vul"
	app.Version = version
	app.ArgsUsage = "image_name"

	app.Usage = "A simple and comprehensive vulnerability scanner for containers"

	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		templateFlag,
		formatFlag,
		inputFlag,
		severityFlag,
		outputFlag,
		exitCodeFlag,
		skipUpdateFlag,
		downloadDBOnlyFlag,
		resetFlag,
		clearCacheFlag,
		quietFlag,
		noProgressFlag,
		ignoreUnfixedFlag,
		debugFlag,
		removedPkgsFlag,
		vulnTypeFlag,
		cacheDirFlag,
		ignoreFileFlag,
		timeoutFlag,
		lightFlag,

		// deprecated options
		cli.StringFlag{
			Name:   "only-update",
			Usage:  "deprecated",
			EnvVar: "VUL_ONLY_UPDATE",
		},
		cli.BoolFlag{
			Name:   "refresh",
			Usage:  "deprecated",
			EnvVar: "VUL_REFRESH",
		},
		cli.BoolFlag{
			Name:   "auto-refresh",
			Usage:  "deprecated",
			EnvVar: "VUL_AUTO_REFRESH",
		},
	}

	app.Commands = []cli.Command{
		NewClientCommand(),
		NewServerCommand(),
	}

	app.Action = standalone.Run
	return app
}

func NewClientCommand() cli.Command {
	return cli.Command{
		Name:    "client",
		Aliases: []string{"c"},
		Usage:   "client mode",
		Action:  client.Run,
		Flags: []cli.Flag{
			templateFlag,
			formatFlag,
			inputFlag,
			severityFlag,
			outputFlag,
			exitCodeFlag,
			clearCacheFlag,
			quietFlag,
			ignoreUnfixedFlag,
			debugFlag,
			removedPkgsFlag,
			vulnTypeFlag,
			ignoreFileFlag,
			cacheDirFlag,
			timeoutFlag,

			// original flags
			token,
			tokenHeader,
			cli.StringFlag{
				Name:   "remote",
				Value:  "http://localhost:4954",
				Usage:  "server address",
				EnvVar: "VUL_REMOTE",
			},
			cli.StringSliceFlag{
				Name:   "custom-headers",
				Usage:  "custom headers",
				EnvVar: "VUL_CUSTOM_HEADERS",
			},
		},
	}
}

func NewServerCommand() cli.Command {
	return cli.Command{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "server mode",
		Action:  server.Run,
		Flags: []cli.Flag{
			skipUpdateFlag,
			downloadDBOnlyFlag,
			resetFlag,
			quietFlag,
			debugFlag,
			cacheDirFlag,

			// original flags
			token,
			tokenHeader,
			cli.StringFlag{
				Name:   "listen",
				Value:  "localhost:4954",
				Usage:  "listen address",
				EnvVar: "VUL_LISTEN",
			},
		},
	}
}
