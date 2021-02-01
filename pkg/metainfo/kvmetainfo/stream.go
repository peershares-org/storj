// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package kvmetainfo

import (
	"context"
	"errors"

	"github.com/gogo/protobuf/proto"

	"czarcoin.org/czarcoin/pkg/encryption"
	"czarcoin.org/czarcoin/pkg/pb"
	"czarcoin.org/czarcoin/pkg/czarcoin"
)

var _ czarcoin.ReadOnlyStream = (*readonlyStream)(nil)

type readonlyStream struct {
	db *DB

	info          czarcoin.Object
	encryptedPath czarcoin.Path
	streamKey     *czarcoin.Key // lazySegmentReader derivedKey
}

func (stream *readonlyStream) Info() czarcoin.Object { return stream.info }

func (stream *readonlyStream) SegmentsAt(ctx context.Context, byteOffset int64, limit int64) (infos []czarcoin.Segment, more bool, err error) {
	defer mon.Task()(&ctx)(&err)

	if stream.info.FixedSegmentSize <= 0 {
		return nil, false, errors.New("not implemented")
	}

	index := byteOffset / stream.info.FixedSegmentSize
	return stream.Segments(ctx, index, limit)
}

func (stream *readonlyStream) segment(ctx context.Context, index int64) (segment czarcoin.Segment, err error) {
	defer mon.Task()(&ctx)(&err)

	segment = czarcoin.Segment{
		Index: index,
	}

	var segmentPath czarcoin.Path
	isLastSegment := segment.Index+1 == stream.info.SegmentCount
	if !isLastSegment {
		segmentPath = getSegmentPath(stream.encryptedPath, index)
		_, meta, err := stream.db.segments.Get(ctx, segmentPath)
		if err != nil {
			return segment, err
		}

		segmentMeta := pb.SegmentMeta{}
		err = proto.Unmarshal(meta.Data, &segmentMeta)
		if err != nil {
			return segment, err
		}

		segment.Size = stream.info.FixedSegmentSize
		copy(segment.EncryptedKeyNonce[:], segmentMeta.KeyNonce)
		segment.EncryptedKey = segmentMeta.EncryptedKey
	} else {
		segmentPath = czarcoin.JoinPaths("l", stream.encryptedPath)
		segment.Size = stream.info.LastSegment.Size
		segment.EncryptedKeyNonce = stream.info.LastSegment.EncryptedKeyNonce
		segment.EncryptedKey = stream.info.LastSegment.EncryptedKey
	}

	contentKey, err := encryption.DecryptKey(segment.EncryptedKey, stream.Info().EncryptionScheme.Cipher, stream.streamKey, &segment.EncryptedKeyNonce)
	if err != nil {
		return segment, err
	}

	nonce := new(czarcoin.Nonce)
	_, err = encryption.Increment(nonce, index+1)
	if err != nil {
		return segment, err
	}

	pointer, _, _, err := stream.db.pointers.Get(ctx, segmentPath)
	if err != nil {
		return segment, err
	}

	if pointer.GetType() == pb.Pointer_INLINE {
		segment.Inline, err = encryption.Decrypt(pointer.InlineSegment, stream.info.EncryptionScheme.Cipher, contentKey, nonce)
	} else {
		segment.PieceID = czarcoin.PieceID(pointer.Remote.PieceId)
		segment.Pieces = make([]czarcoin.Piece, 0, len(pointer.Remote.RemotePieces))
		for _, piece := range pointer.Remote.RemotePieces {
			var nodeID czarcoin.NodeID
			copy(nodeID[:], piece.NodeId.Bytes())
			segment.Pieces = append(segment.Pieces, czarcoin.Piece{Number: byte(piece.PieceNum), Location: nodeID})
		}
	}

	return segment, nil
}

func (stream *readonlyStream) Segments(ctx context.Context, index int64, limit int64) (infos []czarcoin.Segment, more bool, err error) {
	defer mon.Task()(&ctx)(&err)

	if index < 0 {
		return nil, false, errors.New("invalid argument")
	}
	if limit <= 0 {
		limit = defaultSegmentLimit
	}
	if index >= stream.info.SegmentCount {
		return nil, false, nil
	}

	infos = make([]czarcoin.Segment, 0, limit)
	for ; index < stream.info.SegmentCount && limit > 0; index++ {
		limit--
		segment, err := stream.segment(ctx, index)
		if err != nil {
			return nil, false, err
		}
		infos = append(infos, segment)
	}

	more = index < stream.info.SegmentCount
	return infos, more, nil
}

type mutableStream struct {
	db   *DB
	info czarcoin.Object
}

func (stream *mutableStream) Info() czarcoin.Object { return stream.info }

func (stream *mutableStream) AddSegments(ctx context.Context, segments ...czarcoin.Segment) error {
	return errors.New("not implemented")
}

func (stream *mutableStream) UpdateSegments(ctx context.Context, segments ...czarcoin.Segment) error {
	return errors.New("not implemented")
}
