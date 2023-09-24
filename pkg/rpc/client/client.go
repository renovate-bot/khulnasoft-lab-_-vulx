package client

import (
	"context"
	"net/http"

	ftypes "github.com/khulnasoft-lab/fanal/types"

	"github.com/khulnasoft-lab/vul/pkg/types"

	"github.com/google/wire"
	"golang.org/x/xerrors"

	"github.com/khulnasoft-lab/vul/pkg/report"
	r "github.com/khulnasoft-lab/vul/pkg/rpc"
	rpc "github.com/khulnasoft-lab/vul/rpc/scanner"
)

// SuperSet binds the dependencies for RPC client
var SuperSet = wire.NewSet(
	NewProtobufClient,
	NewScanner,
)

// RemoteURL for RPC remote host
type RemoteURL string

// NewProtobufClient is the factory method to return RPC scanner
func NewProtobufClient(remoteURL RemoteURL) rpc.Scanner {
	return rpc.NewScannerProtobufClient(string(remoteURL), &http.Client{})
}

// CustomHeaders for holding HTTP headers
type CustomHeaders http.Header

// Scanner implements the RPC scanner
type Scanner struct {
	customHeaders CustomHeaders
	client        rpc.Scanner
}

// NewScanner is the factory method to return RPC Scanner
func NewScanner(customHeaders CustomHeaders, s rpc.Scanner) Scanner {
	return Scanner{customHeaders: customHeaders, client: s}
}

// Scan scans the image
func (s Scanner) Scan(target, artifactKey string, blobKeys []string, options types.ScanOptions) (report.Results, *ftypes.OS, error) {
	ctx := WithCustomHeaders(context.Background(), http.Header(s.customHeaders))

	var res *rpc.ScanResponse
	err := r.Retry(func() error {
		var err error
		res, err = s.client.Scan(ctx, &rpc.ScanRequest{
			Target:     target,
			ArtifactId: artifactKey,
			BlobIds:    blobKeys,
			Options: &rpc.ScanOptions{
				VulnType:        options.VulnType,
				SecurityChecks:  options.SecurityChecks,
				ListAllPackages: options.ListAllPackages,
			},
		})
		return err
	})
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to detect vulnerabilities via RPC: %w", err)
	}

	return r.ConvertFromRPCResults(res.Results), r.ConvertFromRPCOS(res.Os), nil
}
