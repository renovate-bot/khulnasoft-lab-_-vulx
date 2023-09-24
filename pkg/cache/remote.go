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

type RemoteCache struct {
	ctx    context.Context // for custom header
	client rpcCache.Cache
}

type RemoteURL string

func NewRemoteCache(url RemoteURL, customHeaders http.Header) cache.ImageCache {
	ctx := client.WithCustomHeaders(context.Background(), customHeaders)
	c := rpcCache.NewCacheProtobufClient(string(url), &http.Client{})
	return &RemoteCache{ctx: ctx, client: c}
}

func (c RemoteCache) PutImage(imageID string, imageInfo types.ImageInfo) error {
	_, err := c.client.PutImage(c.ctx, rpc.ConvertToRpcImageInfo(imageID, imageInfo))
	if err != nil {
		return xerrors.Errorf("unable to store cache on the server: %w", err)
	}
	return nil
}

func (c RemoteCache) PutLayer(layerID, decompressedLayerID string, layerInfo types.LayerInfo) error {
	_, err := c.client.PutLayer(c.ctx, rpc.ConvertToRpcLayerInfo(layerID, decompressedLayerID, layerInfo))
	if err != nil {
		return xerrors.Errorf("unable to store cache on the server: %w", err)
	}
	return nil
}

func (c RemoteCache) MissingLayers(imageID string, layerIDs []string) (bool, []string, error) {
	layers, err := c.client.MissingLayers(c.ctx, rpc.ConvertToMissingLayersRequest(imageID, layerIDs))
	if err != nil {
		return false, nil, xerrors.Errorf("unable to fetch missing layers: %w", err)
	}
	return layers.MissingImage, layers.MissingLayerIds, nil
}