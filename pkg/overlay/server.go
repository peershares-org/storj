// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package overlay

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/spacemonkeygo/monkit.v2"

	"czarcoin.org/czarcoin/pkg/dht"
	"czarcoin.org/czarcoin/pkg/pb"
	"czarcoin.org/czarcoin/pkg/czarcoin"
	"czarcoin.org/czarcoin/storage"
)

// ServerError creates class of errors for stack traces
var ServerError = errs.Class("Server Error")

// Server implements our overlay RPC service
type Server struct {
	logger  *zap.Logger
	dht     dht.DHT
	cache   *Cache
	metrics *monkit.Registry
}

// NewServer creates a new Overlay Server
func NewServer(log *zap.Logger, cache *Cache, dht dht.DHT) *Server {
	return &Server{
		dht:     dht,
		cache:   cache,
		logger:  log,
		metrics: monkit.Default,
	}
}

// Lookup finds the address of a node in our overlay network
func (o *Server) Lookup(ctx context.Context, req *pb.LookupRequest) (*pb.LookupResponse, error) {
	na, err := o.cache.Get(ctx, req.NodeId)

	if err != nil {
		o.logger.Error("Error looking up node", zap.Error(err), zap.String("nodeID", req.NodeId.String()))
		return nil, err
	}

	return &pb.LookupResponse{
		Node: na,
	}, nil
}

// BulkLookup finds the addresses of nodes in our overlay network
func (o *Server) BulkLookup(ctx context.Context, reqs *pb.LookupRequests) (*pb.LookupResponses, error) {
	ns, err := o.cache.GetAll(ctx, lookupRequestsToNodeIDs(reqs))
	if err != nil {
		return nil, ServerError.New("could not get nodes requested %s\n", err)
	}
	return nodesToLookupResponses(ns), nil
}

// FindStorageNodes searches the overlay network for nodes that meet the provided requirements
func (o *Server) FindStorageNodes(ctx context.Context, req *pb.FindStorageNodesRequest) (resp *pb.FindStorageNodesResponse, err error) {
	opts := req.GetOpts()
	maxNodes := req.GetMaxNodes()
	if maxNodes <= 0 {
		maxNodes = opts.GetAmount()
	}

	excluded := opts.ExcludedNodes
	restrictions := opts.GetRestrictions()
	reputation := opts.GetMinStats()

	var startID czarcoin.NodeID
	result := []*pb.Node{}
	for {
		var nodes []*pb.Node
		nodes, startID, err = o.populate(ctx, req.Start, maxNodes, restrictions, reputation, excluded)
		if err != nil {
			return nil, Error.Wrap(err)
		}

		resultNodes := []*pb.Node{}
		usedAddrs := make(map[string]bool)
		for _, n := range nodes {
			addr := n.Address.GetAddress()
			excluded = append(excluded, n.Id) // exclude all nodes on next iteration
			if !usedAddrs[addr] {
				resultNodes = append(resultNodes, n)
				usedAddrs[addr] = true
			}
		}
		if len(resultNodes) <= 0 {
			break
		}

		result = append(result, resultNodes...)

		if len(result) >= int(maxNodes) || startID == (czarcoin.NodeID{}) {
			break
		}

	}

	if len(result) < int(maxNodes) {
		return nil, status.Errorf(codes.ResourceExhausted, fmt.Sprintf("requested %d nodes, only %d nodes matched the criteria requested", maxNodes, len(result)))
	}

	if len(result) > int(maxNodes) {
		result = result[:maxNodes]
	}

	return &pb.FindStorageNodesResponse{
		Nodes: result,
	}, nil
}

func (o *Server) getNodes(ctx context.Context, keys storage.Keys) ([]*pb.Node, error) {
	values, err := o.cache.DB.GetAll(keys)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	nodes := []*pb.Node{}
	for _, v := range values {
		n := &pb.Node{}
		if err := proto.Unmarshal(v, n); err != nil {
			return nil, Error.Wrap(err)
		}

		nodes = append(nodes, n)
	}

	return nodes, nil

}

func (o *Server) populate(ctx context.Context, startID czarcoin.NodeID, maxNodes int64,
	minRestrictions *pb.NodeRestrictions, minReputation *pb.NodeStats,
	excluded czarcoin.NodeIDList) ([]*pb.Node, czarcoin.NodeID, error) {

	limit := int(maxNodes * 2)
	keys, err := o.cache.DB.List(startID.Bytes(), limit)
	if err != nil {
		o.logger.Error("Error listing nodes", zap.Error(err))
		return nil, czarcoin.NodeID{}, Error.Wrap(err)
	}

	if len(keys) <= 0 {
		o.logger.Info("No Keys returned from List operation")
		return []*pb.Node{}, startID, nil
	}

	// TODO: should this be `var result []*pb.Node` ?
	result := []*pb.Node{}
	nodes, err := o.getNodes(ctx, keys)
	if err != nil {
		o.logger.Error("Error getting nodes", zap.Error(err))
		return nil, czarcoin.NodeID{}, Error.Wrap(err)
	}

	for _, v := range nodes {
		if v.Type != pb.NodeType_STORAGE {
			continue
		}

		nodeRestrictions := v.GetRestrictions()
		nodeReputation := v.GetReputation()

		if nodeRestrictions.GetFreeBandwidth() < minRestrictions.GetFreeBandwidth() ||
			nodeRestrictions.GetFreeDisk() < minRestrictions.GetFreeDisk() ||
			nodeReputation.GetUptimeRatio() < minReputation.GetUptimeRatio() ||
			nodeReputation.GetUptimeCount() < minReputation.GetUptimeCount() ||
			nodeReputation.GetAuditSuccessRatio() < minReputation.GetAuditSuccessRatio() ||
			nodeReputation.GetAuditCount() < minReputation.GetAuditCount() ||
			contains(excluded, v.Id) {
			continue
		}
		result = append(result, v)
	}

	var nextStart czarcoin.NodeID
	if len(keys) < limit {
		nextStart = czarcoin.NodeID{}
	} else {
		nextStart, err = czarcoin.NodeIDFromBytes(keys[len(keys)-1])
	}
	if err != nil {
		return nil, czarcoin.NodeID{}, Error.Wrap(err)
	}

	return result, nextStart, nil
}

// contains checks if item exists in list
func contains(nodeIDs czarcoin.NodeIDList, searchID czarcoin.NodeID) bool {
	for _, id := range nodeIDs {
		if bytes.Equal(id.Bytes(), searchID.Bytes()) {
			return true
		}
	}
	return false
}

// lookupRequestsToNodeIDs returns the nodeIDs from the LookupRequests
func lookupRequestsToNodeIDs(reqs *pb.LookupRequests) (ids czarcoin.NodeIDList) {
	for _, v := range reqs.LookupRequest {
		ids = append(ids, v.NodeId)
	}
	return ids
}

// nodesToLookupResponses returns LookupResponses from the nodes
func nodesToLookupResponses(nodes []*pb.Node) *pb.LookupResponses {
	var rs []*pb.LookupResponse
	for _, v := range nodes {
		r := &pb.LookupResponse{Node: v}
		rs = append(rs, r)
	}
	return &pb.LookupResponses{LookupResponse: rs}
}
