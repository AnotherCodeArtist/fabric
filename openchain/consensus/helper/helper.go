/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package helper

import (
	"fmt"

	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"golang.org/x/net/context"

	"github.com/openblockchain/obc-peer/openchain/chaincode"
	"github.com/openblockchain/obc-peer/openchain/consensus"
	"github.com/openblockchain/obc-peer/openchain/peer"
	pb "github.com/openblockchain/obc-peer/protos"
)

// =============================================================================
// Init
// =============================================================================

var logger *logging.Logger // package-level logger

func init() {
	logger = logging.MustGetLogger("consensus/helper")
}

// =============================================================================
// Structure definitions go here
// =============================================================================

// Helper contains the reference to coordinator for broadcasts/unicasts.
type Helper struct {
	coordinator peer.MessageHandlerCoordinator
}

// =============================================================================
// Constructors go here
// =============================================================================

// NewHelper constructs the consensus helper object.
func NewHelper(mhc peer.MessageHandlerCoordinator) consensus.CPI {
	return &Helper{coordinator: mhc}
}

// =============================================================================
// Stack-facing implementation goes here
// =============================================================================

// GetReplicaHash returns the crypto IDs of the current replica and the whole network
// func (h *Helper) GetReplicaHash() (self []byte, network [][]byte, err error) {
// ATTN: Until the crypto package is integrated, this functions returns
// the <IP:port>s of the current replica and the whole network instead
func (h *Helper) GetReplicaHash() (self string, network []string, err error) {
	// v, _ := h.coordinator.GetValidator()
	// self = v.GetID()
	peer, err := peer.GetPeerEndpoint()
	if err != nil {
		return "", nil, err
	}
	self = peer.Address

	config := viper.New()
	config.SetConfigName("openchain")
	config.AddConfigPath("./")
	err = config.ReadInConfig()
	if err != nil {
		err = fmt.Errorf("Fatal error reading root config: %s", err)
		return self, nil, err
	}

	// encodedHashes := config.GetStringSlice("peer.validator.replicas.hashes")
	/* network = make([][]byte, len(encodedHashes))
	for i, v := range encodedHashes {
		network[i], _ = base64.StdEncoding.DecodeString(v)
	} */
	network = config.GetStringSlice("peer.validator.replicas.ips")

	return self, network, nil
}

// GetReplicaID returns the uint handle corresponding to a replica address
// func (h *Helper) GetReplicaID(hash []byte) (id uint64, err error) {
func (h *Helper) GetReplicaID(addr string) (id uint64, err error) {
	_, network, err := h.GetReplicaHash()
	if err != nil {
		return uint64(0), err
	}
	for i, v := range network {
		// if bytes.Equal(v, hash) {
		if v == addr {
			return uint64(i), nil
		}
	}

	//err = fmt.Errorf("Couldn't find crypto ID in list of VP IDs given in config")
	err = fmt.Errorf("Couldn't find IP:port in list of VP addresses given in config")
	return uint64(0), err
}

// Broadcast sends a message to all validating peers.
func (h *Helper) Broadcast(msg *pb.OpenchainMessage) error {
	_ = h.coordinator.Broadcast(msg) // TODO process the errors
	return nil
}

// Unicast sends a message to a specified receiver.
func (h *Helper) Unicast(msgPayload []byte, receiver string) error {
	// TODO Call a function in the comms layer; wait for Jeff's implementation.
	return nil
}

// ExecTXs executes all the transactions listed in the txs array one-by-one.
// If all the executions are successful, it returns the candidate global state hash, and nil error array.
func (h *Helper) ExecTXs(txs []*pb.Transaction) ([]byte, []error) {
	return chaincode.ExecuteTransactions(context.Background(), chaincode.DefaultChain, txs)
}
