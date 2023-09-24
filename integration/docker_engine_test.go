//go:build integration
// +build integration

package integration

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/khulnasoft-lab/vul/pkg/commands"
)

func TestRun_WithDockerEngine(t *testing.T) {
	testCases := []struct {
		name                string
		withImageSubcommand bool
		imageTag            string
		invalidImage        bool
		ignoreUnfixed       bool
		severity            []string
		ignoreIDs           []string
		testfile            string
		wantOutputFile      string
		wantError           string
	}{
		// All of these cases should pass for either
		// $ vul <args>
		// $ vul image <args>
		{
			name:           "happy path, valid image path, alpine:3.10",
			imageTag:       "alpine:3.10",
			wantOutputFile: "testdata/alpine-310.json.golden",
			testfile:       "testdata/fixtures/images/alpine-310.tar.gz",
		},
		{
			name:                "happy path, valid image path, with image subcommand, alpine:3.10",
			withImageSubcommand: true,
			imageTag:            "alpine:3.10",
			wantOutputFile:      "testdata/alpine-310.json.golden",
			testfile:            "testdata/fixtures/images/alpine-310.tar.gz",
		},
		{
			name:           "happy path, valid image path, alpine:3.10, ignore unfixed",
			ignoreUnfixed:  true,
			imageTag:       "alpine:3.10",
			wantOutputFile: "testdata/alpine-310-ignore-unfixed.json.golden",
			testfile:       "testdata/fixtures/images/alpine-310.tar.gz",
		},
		{
			name:           "happy path, valid image path, alpine:3.10, ignore unfixed, with medium and high severity",
			ignoreUnfixed:  true,
			severity:       []string{"MEDIUM", "HIGH"},
			imageTag:       "alpine:3.10",
			wantOutputFile: "testdata/alpine-310-medium-high.json.golden",
			testfile:       "testdata/fixtures/images/alpine-310.tar.gz",
		},
		{
			name:           "happy path, valid image path, alpine:3.10, with .vulignore",
			imageTag:       "alpine:3.10",
			ignoreIDs:      []string{"CVE-2019-1549", "CVE-2019-1563"},
			wantOutputFile: "testdata/alpine-310-ignore-cveids.json.golden",
			testfile:       "testdata/fixtures/images/alpine-310.tar.gz",
		},
		{
			name:           "happy path, valid image path, alpine:3.9",
			imageTag:       "alpine:3.9",
			wantOutputFile: "testdata/alpine-39.json.golden",
			testfile:       "testdata/fixtures/images/alpine-39.tar.gz",
		},
		{
			name:           "happy path, valid image path, amazonlinux:1",
			imageTag:       "amazonlinux:1",
			wantOutputFile: "testdata/amazon-1.json.golden",
			testfile:       "testdata/fixtures/images/amazon-1.tar.gz",
		},
		{
			name:           "happy path, valid image path, amazonlinux:2",
			imageTag:       "amazonlinux:2",
			wantOutputFile: "testdata/amazon-2.json.golden",
			testfile:       "testdata/fixtures/images/amazon-2.tar.gz",
		},
		{
			name:           "happy path, valid image path, centos:6",
			imageTag:       "centos:6",
			wantOutputFile: "testdata/centos-6.json.golden",
			testfile:       "testdata/fixtures/images/centos-6.tar.gz",
		},
		{
			name:           "happy path, valid image path, centos:7",
			imageTag:       "centos:7",
			wantOutputFile: "testdata/centos-7.json.golden",
			testfile:       "testdata/fixtures/images/centos-7.tar.gz",
		},
		{
			name:           "happy path, valid image path, centos:7, with --ignore-unfixed option",
			imageTag:       "centos:7",
			ignoreUnfixed:  true,
			wantOutputFile: "testdata/centos-7-ignore-unfixed.json.golden",
			testfile:       "testdata/fixtures/images/centos-7.tar.gz",
		},
		{
			name:           "happy path, valid image path, centos:7, with --ignore-unfixed option, with low and high severity",
			imageTag:       "centos:7",
			ignoreUnfixed:  true,
			severity:       []string{"LOW", "HIGH"},
			wantOutputFile: "testdata/centos-7-low-high.json.golden",
			testfile:       "testdata/fixtures/images/centos-7.tar.gz",
		},
		{
			name:           "happy path, valid image path, debian:buster",
			imageTag:       "debian:buster",
			wantOutputFile: "testdata/debian-buster.json.golden",
			testfile:       "testdata/fixtures/images/debian-buster.tar.gz",
		},
		{
			name:           "happy path, valid image path, debian:buster, with --ignore-unfixed option",
			ignoreUnfixed:  true,
			imageTag:       "debian:buster",
			wantOutputFile: "testdata/debian-buster-ignore-unfixed.json.golden",
			testfile:       "testdata/fixtures/images/debian-buster.tar.gz",
		},
		{
			name:           "happy path, valid image path, debian:stretch",
			imageTag:       "debian:stretch",
			wantOutputFile: "testdata/debian-stretch.json.golden",
			testfile:       "testdata/fixtures/images/debian-stretch.tar.gz",
		},
		{
			name:           "happy path, valid image path, distroless:base",
			imageTag:       "gcr.io/distroless/base:latest",
			wantOutputFile: "testdata/distroless-base.json.golden",
			testfile:       "testdata/fixtures/images/distroless-base.tar.gz",
		},
		{
			name:           "happy path, valid image path, distroless:base",
			imageTag:       "gcr.io/distroless/base:latest",
			wantOutputFile: "testdata/distroless-base.json.golden",
			testfile:       "testdata/fixtures/images/distroless-base.tar.gz",
		},
		{
			name:           "happy path, valid image path, distroless:base, with --ignore-unfixed option",
			imageTag:       "gcr.io/distroless/base:latest",
			ignoreUnfixed:  true,
			wantOutputFile: "testdata/distroless-base-ignore-unfixed.json.golden",
			testfile:       "testdata/fixtures/images/distroless-base.tar.gz",
		},
		{
			name:           "happy path, valid image path, distroless:python2.7",
			imageTag:       "gcr.io/distroless/python2.7:latest",
			wantOutputFile: "testdata/distroless-python27.json.golden",
			testfile:       "testdata/fixtures/images/distroless-python27.tar.gz",
		},
		{
			name:           "happy path, valid image path, oraclelinux:6-slim",
			imageTag:       "oraclelinux:6-slim",
			wantOutputFile: "testdata/oraclelinux-6-slim.json.golden",
			testfile:       "testdata/fixtures/images/oraclelinux-6-slim.tar.gz",
		},
		{
			name:           "happy path, valid image path, oraclelinux:7-slim",
			imageTag:       "oraclelinux:7-slim",
			wantOutputFile: "testdata/oraclelinux-7-slim.json.golden",
			testfile:       "testdata/fixtures/images/oraclelinux-7-slim.tar.gz",
		},
		{
			name:           "happy path, valid image path, oraclelinux:8-slim",
			imageTag:       "oraclelinux:8-slim",
			wantOutputFile: "testdata/oraclelinux-8-slim.json.golden",
			testfile:       "testdata/fixtures/images/oraclelinux-8-slim.tar.gz",
		},
		{
			name:           "happy path, valid image path, ubuntu:16.04",
			imageTag:       "ubuntu:16.04",
			wantOutputFile: "testdata/ubuntu-1604.json.golden",
			testfile:       "testdata/fixtures/images/ubuntu-1604.tar.gz",
		},
		{
			name:           "happy path, valid image path, ubuntu:18.04",
			imageTag:       "ubuntu:18.04",
			wantOutputFile: "testdata/ubuntu-1804.json.golden",
			testfile:       "testdata/fixtures/images/ubuntu-1804.tar.gz",
		},
		{
			name:           "happy path, valid image path, ubuntu:18.04, with --ignore-unfixed option",
			imageTag:       "ubuntu:18.04",
			ignoreUnfixed:  true,
			wantOutputFile: "testdata/ubuntu-1804-ignore-unfixed.json.golden",
			testfile:       "testdata/fixtures/images/ubuntu-1804.tar.gz",
		},
		{
			name:           "happy path, valid image path, registry.redhat.io/ubi7",
			imageTag:       "registry.redhat.io/ubi7",
			wantOutputFile: "testdata/ubi-7.json.golden",
			testfile:       "testdata/fixtures/images/ubi-7.tar.gz",
		},
		{
			name:           "happy path, valid image path, opensuse leap 15.1",
			imageTag:       "opensuse/leap:latest",
			wantOutputFile: "testdata/opensuse-leap-151.json.golden",
			testfile:       "testdata/fixtures/images/opensuse-leap-151.tar.gz",
		},
		{
			name:           "happy path, valid image path, opensuse leap 42.3",
			imageTag:       "opensuse/leap:42.3",
			wantOutputFile: "testdata/opensuse-leap-423.json.golden",
			testfile:       "testdata/fixtures/images/opensuse-leap-423.tar.gz",
		},
		{
			name:           "happy path, valid image path, photon 1.0",
			imageTag:       "photon:1.0-20190823",
			wantOutputFile: "testdata/photon-10.json.golden",
			testfile:       "testdata/fixtures/images/photon-10.tar.gz",
		},
		{
			name:           "happy path, valid image path, photon 2.0",
			imageTag:       "photon:2.0-20190726",
			wantOutputFile: "testdata/photon-20.json.golden",
			testfile:       "testdata/fixtures/images/photon-20.tar.gz",
		},
		{
			name:           "happy path, valid image path, photon 3.0",
			imageTag:       "photon:3.0-20190823",
			wantOutputFile: "testdata/photon-30.json.golden",
			testfile:       "testdata/fixtures/images/photon-30.tar.gz",
		},
		{
			name:           "buxybox with Cargo.lock integration",
			imageTag:       "busy-cargo:latest",
			wantOutputFile: "testdata/busybox-with-lockfile.json.golden",
			testfile:       "testdata/fixtures/images/busybox-with-lockfile.tar.gz",
		},
		{
			name:         "sad path, invalid image",
			invalidImage: true,
			testfile:     "badimage:latest",
			wantError:    "unable to inspect the image (index.docker.io/library/badimage:latest)",
		},
	}

	// Set up testing DB
	cacheDir := gunzipDB(t)

	ctx := context.Background()
	defer ctx.Done()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.invalidImage {
				testfile, err := os.Open(tc.testfile)
				require.NoError(t, err, tc.name)

				// ensure image doesnt already exists
				_, _ = cli.ImageRemove(ctx, tc.testfile, types.ImageRemoveOptions{
					Force:         true,
					PruneChildren: true,
				})

				// load image into docker engine
				res, err := cli.ImageLoad(ctx, testfile, true)
				require.NoError(t, err, tc.name)
				io.Copy(io.Discard, res.Body)

				// tag our image to something unique
				err = cli.ImageTag(ctx, tc.imageTag, tc.testfile)
				require.NoError(t, err, tc.name)
			}

			tmpDir := t.TempDir()
			output := filepath.Join(tmpDir, "result.json")

			// run vul
			app := commands.NewApp("dev")
			vulArgs := []string{"vul"}
			vulArgs = append(vulArgs, "--cache-dir", cacheDir)
			if tc.withImageSubcommand {
				vulArgs = append(vulArgs, "image")
			}

			vulArgs = append(vulArgs, []string{"--skip-update", "--format=json", "--output", output}...)

			if tc.ignoreUnfixed {
				vulArgs = append(vulArgs, "--ignore-unfixed")
			}
			if len(tc.severity) != 0 {
				vulArgs = append(vulArgs,
					[]string{"--severity", strings.Join(tc.severity, ",")}...,
				)
			}
			if len(tc.ignoreIDs) != 0 {
				vulIgnore := ".vulignore"
				err := os.WriteFile(vulIgnore, []byte(strings.Join(tc.ignoreIDs, "\n")), 0444)
				assert.NoError(t, err, "failed to write .vulignore")
				defer os.Remove(vulIgnore)
			}
			vulArgs = append(vulArgs, tc.testfile)

			err = app.Run(vulArgs)
			switch {
			case tc.wantError != "":
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.wantError, tc.name)
				return
			default:
				assert.NoError(t, err, tc.name)
			}

			// check for vulnerability output info
			got := readReport(t, output)
			want := readReport(t, tc.wantOutputFile)
			assert.Equal(t, want, got)

			// cleanup
			_, err = cli.ImageRemove(ctx, tc.testfile, types.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			})
			_, err = cli.ImageRemove(ctx, tc.imageTag, types.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			})
			assert.NoError(t, err, tc.name)
		})
	}
}
