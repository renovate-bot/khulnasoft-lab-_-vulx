package amazon

import (
	"strings"
	"time"

	"k8s.io/utils/clock"

	version "github.com/knqyf263/go-deb-version"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/khulnasoft-lab/vul-db/pkg/vulnsrc/amazon"
	ftypes "github.com/khulnasoft-lab/vul/pkg/fanal/types"
	"github.com/khulnasoft-lab/vul/pkg/log"
	"github.com/khulnasoft-lab/vul/pkg/scanner/utils"
	"github.com/khulnasoft-lab/vul/pkg/types"
)

var (
	eolDates = map[string]time.Time{
		// https://aws.amazon.com/jp/blogs/aws/update-on-amazon-linux-ami-end-of-life/
		"1": time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
		// https://aws.amazon.com/amazon-linux-2/faqs/?nc1=h_ls
		"2": time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC),
		// Amazon Linux 2022 was renamed to 2023. AL2022 is not currently supported.
		"2023": time.Date(2028, 3, 15, 23, 59, 59, 0, time.UTC),
	}
)

type options struct {
	clock clock.Clock
	l     *zap.SugaredLogger
}

type option func(*options)

func WithClock(clock clock.Clock) option {
	return func(opts *options) {
		opts.clock = clock
	}
}

// Scanner to scan amazon vulnerabilities
type Scanner struct {
	ac amazon.VulnSrc
	options
}

// NewScanner is the factory method to return Amazon scanner
func NewScanner(opts ...option) *Scanner {
	o := &options{
		l:     log.Logger,
		clock: clock.RealClock{},
	}

	for _, opt := range opts {
		opt(o)
	}
	return &Scanner{
		ac:      amazon.NewVulnSrc(),
		options: *o,
	}
}

// Detect scans the packages using amazon scanner
func (s *Scanner) Detect(osVer string, _ *ftypes.Repository, pkgs []ftypes.Package) ([]types.DetectedVulnerability, error) {
	log.Logger.Info("Detecting Amazon Linux vulnerabilities...")

	osVer = strings.Fields(osVer)[0]
	if osVer != "2" && osVer != "2022" && osVer != "2023" {
		osVer = "1"
	}
	log.Logger.Debugf("amazon: os version: %s", osVer)
	log.Logger.Debugf("amazon: the number of packages: %d", len(pkgs))

	var vulns []types.DetectedVulnerability
	for _, pkg := range pkgs {
		advisories, err := s.ac.Get(osVer, pkg.Name)
		if err != nil {
			return nil, xerrors.Errorf("failed to get amazon advisories: %w", err)
		}

		installed := utils.FormatVersion(pkg)
		if installed == "" {
			continue
		}

		installedVersion, err := version.NewVersion(installed)
		if err != nil {
			log.Logger.Debugf("failed to parse Amazon Linux installed package version: %s", err)
			continue
		}

		for _, adv := range advisories {
			fixedVersion, err := version.NewVersion(adv.FixedVersion)
			if err != nil {
				log.Logger.Debugf("failed to parse Amazon Linux package version: %s", err)
				continue
			}

			if installedVersion.LessThan(fixedVersion) {
				vuln := types.DetectedVulnerability{
					VulnerabilityID:  adv.VulnerabilityID,
					PkgID:            pkg.ID,
					PkgName:          pkg.Name,
					InstalledVersion: installed,
					FixedVersion:     adv.FixedVersion,
					PkgRef:           pkg.Ref,
					Layer:            pkg.Layer,
					Custom:           adv.Custom,
					DataSource:       adv.DataSource,
				}
				vulns = append(vulns, vuln)
			}
		}
	}
	return vulns, nil
}

// IsSupportedVersion checks if os can be scanned using amazon scanner
func (s *Scanner) IsSupportedVersion(osFamily, osVer string) bool {
	osVer = strings.Fields(osVer)[0]
	if osVer != "2" && osVer != "2022" && osVer != "2023" {
		osVer = "1"
	}
	eol, ok := eolDates[osVer]
	if !ok {
		log.Logger.Warnf("This OS version is not on the EOL list: %s %s", osFamily, osVer)
		return false
	}

	return s.clock.Now().Before(eol)
}
