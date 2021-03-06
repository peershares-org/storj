// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package miniogw

import (
	"context"
	"os"

	"github.com/minio/cli"
	minio "github.com/minio/minio/cmd"
	"github.com/vivint/infectious"
	"go.uber.org/zap"

	"czarcoin.org/czarcoin/pkg/eestream"
	"czarcoin.org/czarcoin/pkg/metainfo/kvmetainfo"
	"czarcoin.org/czarcoin/pkg/miniogw/logging"
	"czarcoin.org/czarcoin/pkg/overlay"
	"czarcoin.org/czarcoin/pkg/pointerdb/pdbclient"
	"czarcoin.org/czarcoin/pkg/provider"
	"czarcoin.org/czarcoin/pkg/storage/buckets"
	ecclient "czarcoin.org/czarcoin/pkg/storage/ec"
	"czarcoin.org/czarcoin/pkg/storage/segments"
	"czarcoin.org/czarcoin/pkg/storage/streams"
	"czarcoin.org/czarcoin/pkg/czarcoin"
)

// RSConfig is a configuration struct that keeps details about default
// redundancy strategy information
type RSConfig struct {
	MaxBufferMem     int `help:"maximum buffer memory (in bytes) to be allocated for read buffers" default:"0x400000"`
	ErasureShareSize int `help:"the size of each new erasure sure in bytes" default:"1024"`
	MinThreshold     int `help:"the minimum pieces required to recover a segment. k." default:"29"`
	RepairThreshold  int `help:"the minimum safe pieces before a repair is triggered. m." default:"35"`
	SuccessThreshold int `help:"the desired total pieces for a segment. o." default:"80"`
	MaxThreshold     int `help:"the largest amount of pieces to encode to. n." default:"95"`
}

// EncryptionConfig is a configuration struct that keeps details about
// encrypting segments
type EncryptionConfig struct {
	Key       string `help:"root key for encrypting the data"`
	BlockSize int    `help:"size (in bytes) of encrypted blocks" default:"1024"`
	DataType  int    `help:"Type of encryption to use for content and metadata (1=AES-GCM, 2=SecretBox)" default:"1"`
	PathType  int    `help:"Type of encryption to use for paths (0=Unencrypted, 1=AES-GCM, 2=SecretBox)" default:"1"`
}

// MinioConfig is a configuration struct that keeps details about starting
// Minio
type MinioConfig struct {
	AccessKey string `help:"Minio Access Key to use" default:"insecure-dev-access-key"`
	SecretKey string `help:"Minio Secret Key to use" default:"insecure-dev-secret-key"`
	Dir       string `help:"Minio generic server config path" default:"$CONFDIR/minio"`
}

// ClientConfig is a configuration struct for the miniogw that controls how
// the miniogw figures out how to talk to the rest of the network.
type ClientConfig struct {
	// TODO(jt): these should probably be the same
	OverlayAddr   string `help:"Address to contact overlay server through"`
	PointerDBAddr string `help:"Address to contact pointerdb server through"`

	APIKey        string `help:"API Key (TODO: this needs to change to macaroons somehow)"`
	MaxInlineSize int    `help:"max inline segment size in bytes" default:"4096"`
	SegmentSize   int64  `help:"the size of a segment in bytes" default:"64000000"`
}

// Config is a general miniogw configuration struct. This should be everything
// one needs to start a minio gateway.
type Config struct {
	Identity provider.IdentityConfig
	Minio    MinioConfig
	Client   ClientConfig
	RS       RSConfig
	Enc      EncryptionConfig
}

// Run starts a Minio Gateway given proper config
func (c Config) Run(ctx context.Context) (err error) {
	defer mon.Task()(&ctx)(&err)

	identity, err := c.Identity.Load()
	if err != nil {
		return err
	}

	err = minio.RegisterGatewayCommand(cli.Command{
		Name:  "czarcoin",
		Usage: "Czarcoin",
		Action: func(cliCtx *cli.Context) error {
			return c.action(ctx, cliCtx, identity)
		},
		HideHelpCommand: true,
	})
	if err != nil {
		return err
	}

	// TODO(jt): Surely there is a better way. This is so upsetting
	err = os.Setenv("MINIO_ACCESS_KEY", c.Minio.AccessKey)
	if err != nil {
		return err
	}
	err = os.Setenv("MINIO_SECRET_KEY", c.Minio.SecretKey)
	if err != nil {
		return err
	}

	minio.Main([]string{"czarcoin", "gateway", "czarcoin",
		"--address", c.Identity.Server.Address, "--config-dir", c.Minio.Dir, "--quiet"})
	return Error.New("unexpected minio exit")
}

func (c Config) action(ctx context.Context, cliCtx *cli.Context, identity *provider.FullIdentity) (err error) {
	defer mon.Task()(&ctx)(&err)

	gw, err := c.NewGateway(ctx, identity)
	if err != nil {
		return err
	}

	minio.StartGateway(cliCtx, logging.Gateway(gw, zap.L()))
	return Error.New("unexpected minio exit")
}

// GetMetainfo returns an implementation of czarcoin.Metainfo
func (c Config) GetMetainfo(ctx context.Context, identity *provider.FullIdentity) (db czarcoin.Metainfo, ss streams.Store, err error) {
	defer mon.Task()(&ctx)(&err)

	oc, err := overlay.NewOverlayClient(identity, c.Client.OverlayAddr)
	if err != nil {
		return nil, nil, err
	}

	pdb, err := pdbclient.NewClient(identity, c.Client.PointerDBAddr, c.Client.APIKey)
	if err != nil {
		return nil, nil, err
	}

	ec := ecclient.NewClient(identity, c.RS.MaxBufferMem)
	fc, err := infectious.NewFEC(c.RS.MinThreshold, c.RS.MaxThreshold)
	if err != nil {
		return nil, nil, err
	}
	rs, err := eestream.NewRedundancyStrategy(eestream.NewRSScheme(fc, c.RS.ErasureShareSize), c.RS.RepairThreshold, c.RS.SuccessThreshold)
	if err != nil {
		return nil, nil, err
	}

	segments := segments.NewSegmentStore(oc, ec, pdb, rs, c.Client.MaxInlineSize)

	if c.RS.ErasureShareSize*c.RS.MinThreshold%c.Enc.BlockSize != 0 {
		err = Error.New("EncryptionBlockSize must be a multiple of ErasureShareSize * RS MinThreshold")
		return nil, nil, err
	}

	key := new(czarcoin.Key)
	copy(key[:], c.Enc.Key)

	streams, err := streams.NewStreamStore(segments, c.Client.SegmentSize, key, c.Enc.BlockSize, czarcoin.Cipher(c.Enc.DataType))
	if err != nil {
		return nil, nil, err
	}

	buckets := buckets.NewStore(streams)

	return kvmetainfo.New(buckets, streams, segments, pdb, key), streams, nil
}

// GetRedundancyScheme returns the configured redundancy scheme for new uploads
func (c Config) GetRedundancyScheme() czarcoin.RedundancyScheme {
	return czarcoin.RedundancyScheme{
		Algorithm:      czarcoin.ReedSolomon,
		RequiredShares: int16(c.RS.MinThreshold),
		RepairShares:   int16(c.RS.RepairThreshold),
		OptimalShares:  int16(c.RS.SuccessThreshold),
		TotalShares:    int16(c.RS.MaxThreshold),
	}
}

// GetEncryptionScheme returns the configured encryption scheme for new uploads
func (c Config) GetEncryptionScheme() czarcoin.EncryptionScheme {
	return czarcoin.EncryptionScheme{
		Cipher:    czarcoin.Cipher(c.Enc.DataType),
		BlockSize: int32(c.Enc.BlockSize),
	}
}

// NewGateway creates a new minio Gateway
func (c Config) NewGateway(ctx context.Context, identity *provider.FullIdentity) (gw minio.Gateway, err error) {
	defer mon.Task()(&ctx)(&err)

	metainfo, streams, err := c.GetMetainfo(ctx, identity)
	if err != nil {
		return nil, err
	}

	return NewCzarcoinGateway(metainfo, streams, czarcoin.Cipher(c.Enc.PathType), c.GetEncryptionScheme(), c.GetRedundancyScheme()), nil
}
