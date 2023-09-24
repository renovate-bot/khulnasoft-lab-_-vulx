package debian_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	fake "k8s.io/utils/clock/testing"

	ftypes "github.com/khulnasoft-lab/fanal/types"
	"github.com/khulnasoft-lab/vul-db/pkg/db"
	dbTypes "github.com/khulnasoft-lab/vul-db/pkg/types"
	"github.com/khulnasoft-lab/vul-db/pkg/vulnsrc/vulnerability"
	"github.com/khulnasoft-lab/vul/pkg/dbtest"
	"github.com/khulnasoft-lab/vul/pkg/detector/ospkg/debian"
	"github.com/khulnasoft-lab/vul/pkg/types"
)

func TestScanner_Detect(t *testing.T) {
	type args struct {
		osVer string
		pkgs  []ftypes.Package
	}
	tests := []struct {
		name     string
		args     args
		fixtures []string
		want     []types.DetectedVulnerability
		wantErr  string
	}{
		{
			name:     "happy path",
			fixtures: []string{"testdata/fixtures/debian.yaml"},
			args: args{
				osVer: "9.1",
				pkgs: []ftypes.Package{
					{
						Name:       "htpasswd",
						Version:    "2.4.24",
						SrcName:    "apache2",
						SrcVersion: "2.4.24",
						Layer: ftypes.Layer{
							DiffID: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
						},
					},
				},
			},
			want: []types.DetectedVulnerability{
				{
					PkgName:          "htpasswd",
					VulnerabilityID:  "CVE-2020-11985",
					VendorIDs:        []string{"DSA-4884-1"},
					InstalledVersion: "2.4.24",
					FixedVersion:     "2.4.25-1",
					Layer: ftypes.Layer{
						DiffID: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
					},
				},
				{
					PkgName:          "htpasswd",
					VulnerabilityID:  "CVE-2021-31618",
					InstalledVersion: "2.4.24",
					SeveritySource:   vulnerability.Debian,
					Vulnerability: dbTypes.Vulnerability{
						Severity: dbTypes.SeverityMedium.String(),
					},
					Layer: ftypes.Layer{
						DiffID: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
					},
				},
			},
		},
		{
			name:     "invalid bucket",
			fixtures: []string{"testdata/fixtures/invalid.yaml"},
			args: args{
				osVer: "9.1",
				pkgs: []ftypes.Package{
					{
						Name:       "htpasswd",
						Version:    "2.4.24",
						SrcName:    "apache2",
						SrcVersion: "2.4.24",
						Layer: ftypes.Layer{
							DiffID: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
						},
					},
				},
			},
			wantErr: "failed to get Debian OVAL advisories",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = dbtest.InitDB(t, tt.fixtures)
			defer db.Close()

			s := debian.NewScanner()
			got, err := s.Detect(tt.args.osVer, tt.args.pkgs)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestScanner_IsSupportedVersion(t *testing.T) {
	type args struct {
		osFamily string
		osVer    string
	}
	tests := []struct {
		name string
		now  time.Time
		args args
		want bool
	}{
		{
			name: "debian 7",
			now:  time.Date(2018, 3, 31, 23, 59, 59, 0, time.UTC),
			args: args{
				osFamily: "debian",
				osVer:    "7",
			},
			want: true,
		},
		{
			name: "debian 8 EOL",
			now:  time.Date(2020, 7, 31, 23, 59, 59, 0, time.UTC),
			args: args{
				osFamily: "debian",
				osVer:    "8.2",
			},
			want: false,
		},
		{
			name: "unknown",
			now:  time.Date(2020, 7, 31, 23, 59, 59, 0, time.UTC),
			args: args{
				osFamily: "debian",
				osVer:    "unknown",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := debian.NewScanner(debian.WithClock(fake.NewFakeClock(tt.now)))
			got := s.IsSupportedVersion(tt.args.osFamily, tt.args.osVer)
			assert.Equal(t, tt.want, got)
		})
	}
}
