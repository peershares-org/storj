// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.
package transport

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"czarcoin.org/czarcoin/internal/identity"
	"czarcoin.org/czarcoin/internal/testczarcoin"
	"czarcoin.org/czarcoin/pkg/pb"
)

var ctx = context.Background()

func TestDialNode(t *testing.T) {
	ca, err := testidentity.NewTestCA(ctx)
	assert.NoError(t, err)
	identity, err := ca.NewIdentity()
	assert.NoError(t, err)

	oc := Transport{
		identity: identity,
	}

	// node.Address.Address == "" condition test
	node := pb.Node{
		Id: testczarcoin.NodeIDFromString("DUMMYID1"),
		Address: &pb.NodeAddress{
			Transport: pb.NodeTransport_TCP_TLS_GRPC,
			Address:   "",
		},
	}
	conn, err := oc.DialNode(ctx, &node)
	assert.Error(t, err)
	assert.Nil(t, conn)

	// node.Address == nil condition test
	node = pb.Node{
		Id:      testczarcoin.NodeIDFromString("DUMMYID2"),
		Address: nil,
	}
	conn, err = oc.DialNode(ctx, &node)
	assert.Error(t, err)
	assert.Nil(t, conn)

	// node is valid argument condition test
	node = pb.Node{
		Id: testczarcoin.NodeIDFromString("DUMMYID3"),
		Address: &pb.NodeAddress{
			Transport: pb.NodeTransport_TCP_TLS_GRPC,
			Address:   "127.0.0.0:9000",
		},
	}
	conn, err = oc.DialNode(ctx, &node)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
}
