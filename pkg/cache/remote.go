package cache

import (
	"context"
	"net/http"

	"golang.org/x/xerrors"

	"github.com/khulnasoft-lab/fanal/cache"
	"github.com/khulnasoft-lab/fanal/types"
	"github.com/khulnasoft-lab/vul/pkg/rpc"
	"github.com/khulnasoft-lab/vul/pkg/rpc/client"
	rpcCache "github.com/khulnasoft-lab/vul/rpc/cache"
)

// RemoteCache implements remote cache
type RemoteCache struct {
	ctx    context.Context // for custom header
	client rpcCache.Cache
}

// RemoteURL to hold remote host
type RemoteURL string

// NewRemoteCache is the factory method for RemoteCache
func NewRemoteCache(url RemoteURL, customHeaders http.Header) cache.ArtifactCache {
	ctx := client.WithCustomHeaders(context.Background(), customHeaders)
	c := rpcCache.NewCacheProtobufClient(string(url), &http.Client{})
	return &RemoteCache{ctx: ctx, client: c}
}

// PutArtifact sends artifact to remote client
func (c RemoteCache) PutArtifact(imageID string, artifactInfo types.ArtifactInfo) error {
	_, err := c.client.PutArtifact(c.ctx, rpc.ConvertToRPCArtifactInfo(imageID, artifactInfo))
	if err != nil {
		return xerrors.Errorf("unable to store cache on the server: %w", err)
	}
	return nil
}

// PutBlob sends blobInfo to remote client
func (c RemoteCache) PutBlob(diffID string, blobInfo types.BlobInfo) error {
	_, err := c.client.PutBlob(c.ctx, rpc.ConvertToRPCBlobInfo(diffID, blobInfo))
	if err != nil {
		return xerrors.Errorf("unable to store cache on the server: %w", err)
	}
	return nil
}

// MissingBlobs fetches missing blobs from RemoteCache
func (c RemoteCache) MissingBlobs(imageID string, layerIDs []string) (bool, []string, error) {
	layers, err := c.client.MissingBlobs(c.ctx, rpc.ConvertToMissingBlobsRequest(imageID, layerIDs))
	if err != nil {
		return false, nil, xerrors.Errorf("unable to fetch missing layers: %w", err)
	}
	return layers.MissingArtifact, layers.MissingBlobIds, nil
}
