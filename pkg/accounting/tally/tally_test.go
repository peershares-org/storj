// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package tally

import (
	"crypto/ecdsa"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"czarcoin.org/czarcoin/internal/identity"
	"czarcoin.org/czarcoin/internal/testcontext"
	"czarcoin.org/czarcoin/internal/testczarcoin"
	"czarcoin.org/czarcoin/pkg/accounting"
	"czarcoin.org/czarcoin/pkg/bwagreement"
	"czarcoin.org/czarcoin/pkg/bwagreement/test"
	"czarcoin.org/czarcoin/pkg/kademlia"
	"czarcoin.org/czarcoin/pkg/overlay"
	"czarcoin.org/czarcoin/pkg/overlay/mocks"
	"czarcoin.org/czarcoin/pkg/pb"
	"czarcoin.org/czarcoin/pkg/pointerdb"
	"czarcoin.org/czarcoin/pkg/czarcoin"
	"czarcoin.org/czarcoin/satellite/satellitedb"
	"czarcoin.org/czarcoin/storage/teststore"
)

func TestIdentifyActiveNodes(t *testing.T) {

}
func TestOnlineNodes(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	logger := zap.NewNop()
	pointerdb := pointerdb.NewServer(teststore.New(), &overlay.Cache{}, logger, pointerdb.Config{}, nil)

	const N = 50
	nodes := []*pb.Node{}
	nodeIDs := czarcoin.NodeIDList{}
	expectedOnline := []*pb.Node{}
	for i := 0; i < N; i++ {
		nodeID := testczarcoin.NodeIDFromString(strconv.Itoa(i))
		n := &pb.Node{Id: nodeID, Type: pb.NodeType_STORAGE, Address: &pb.NodeAddress{Address: ""}}
		nodes = append(nodes, n)
		if i%(rand.Intn(5)+2) == 0 {
			id := testczarcoin.NodeIDFromString("id" + nodeID.String())
			nodeIDs = append(nodeIDs, id)
		} else {
			nodeIDs = append(nodeIDs, nodeID)
			expectedOnline = append(expectedOnline, n)
		}
	}
	overlayServer := mocks.NewOverlay(nodes)
	kad := &kademlia.Kademlia{}
	limit := 0
	interval := time.Second

	accountingDb, err := accounting.NewDb("sqlite3://file::memory:?mode=memory&cache=shared")
	assert.NoError(t, err)
	defer ctx.Check(accountingDb.Close)

	masterDB, err := satellitedb.NewDB("sqlite3://file::memory:?mode=memory&cache=shared")
	assert.NoError(t, err)
	defer ctx.Check(masterDB.Close)

	tally := newTally(logger, accountingDb, masterDB.BandwidthAgreement(), pointerdb, overlayServer, kad, limit, interval)

	online, err := tally.onlineNodes(ctx, nodeIDs)
	assert.NoError(t, err)
	assert.Equal(t, expectedOnline, online)
}

func TestTallyAtRestStorage(t *testing.T) {

}

func TestNeedToContact(t *testing.T) {

}

func TestUpdateGranularTable(t *testing.T) {

}

func TestQueryNoAgreements(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	//get stuff we need
	pointerdb := pointerdb.NewServer(teststore.New(), &overlay.Cache{}, zap.NewNop(), pointerdb.Config{}, nil)
	overlayServer := mocks.NewOverlay([]*pb.Node{})
	kad := &kademlia.Kademlia{}
	accountingDb, err := accounting.NewDb("sqlite3://file::memory:?mode=memory&cache=shared")
	assert.NoError(t, err)
	defer ctx.Check(accountingDb.Close)

	masterDB, err := satellitedb.NewDB("sqlite3://file::memory:?mode=memory&cache=shared")
	assert.NoError(t, err)
	defer ctx.Check(masterDB.Close)

	tally := newTally(zap.NewNop(), accountingDb, masterDB.BandwidthAgreement(), pointerdb, overlayServer, kad, 0, time.Second)

	//check the db
	err = tally.Query(ctx)
	assert.NoError(t, err)
}

func TestQueryWithBw(t *testing.T) {
	ctx := testcontext.New(t)
	defer ctx.Cleanup()

	//get stuff we need
	pointerdb := pointerdb.NewServer(teststore.New(), &overlay.Cache{}, zap.NewNop(), pointerdb.Config{}, nil)
	overlayServer := mocks.NewOverlay([]*pb.Node{})
	kad := &kademlia.Kademlia{}
	accountingDb, err := accounting.NewDb("sqlite3://file::memory:?mode=memory&cache=shared")
	assert.NoError(t, err)
	defer ctx.Check(accountingDb.Close)

	masterDB, err := satellitedb.NewDB("sqlite3://file::memory:?mode=memory&cache=shared")
	assert.NoError(t, err)
	err = masterDB.CreateTables()
	assert.NoError(t, err)
	defer ctx.Check(masterDB.Close)

	bwDb := masterDB.BandwidthAgreement()
	tally := newTally(zap.NewNop(), accountingDb, bwDb, pointerdb, overlayServer, kad, 0, time.Second)

	//get a private key
	fiC, err := testidentity.NewTestIdentity()
	assert.NoError(t, err)
	k, ok := fiC.Key.(*ecdsa.PrivateKey)
	assert.True(t, ok)
	//generate an agreement with the key
	pba, err := test.GeneratePayerBandwidthAllocation(pb.PayerBandwidthAllocation_GET, k)
	assert.NoError(t, err)
	rba, err := test.GenerateRenterBandwidthAllocation(pba, k)
	assert.NoError(t, err)
	//save to db
	err = bwDb.CreateAgreement(ctx, bwagreement.Agreement{Signature: rba.GetSignature(), Agreement: rba.GetData()})
	assert.NoError(t, err)

	//check the db
	err = tally.Query(ctx)
	assert.NoError(t, err)
}
