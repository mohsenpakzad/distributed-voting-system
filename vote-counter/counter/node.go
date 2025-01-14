package counter

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"

	"github.com/mohsenpakzad/distributed-voting-system/shared/models"
)

type Node struct {
	raft        *raft.Raft
	fsm         *FSM
	nodeAddress string
}

func NewNode(nodeID, nodeAddress, dataDir string, bootstrap bool) (*Node, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	addr, err := net.ResolveTCPAddr("tcp", nodeAddress)
	if err != nil {
		return nil, err
	}

	transport, err := raft.NewTCPTransport(nodeAddress, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	snapshots, err := raft.NewFileSnapshotStore(filepath.Join(dataDir, "snapshots"), 2, os.Stderr)
	if err != nil {
		return nil, err
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft.db"))
	if err != nil {
		return nil, err
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "stable.db"))
	if err != nil {
		return nil, err
	}

	fsm := NewFSM()

	r, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, err
	}

	if bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		r.BootstrapCluster(configuration)
	}

	return &Node{
		raft:        r,
		fsm:         fsm,
		nodeAddress: nodeAddress,
	}, nil
}

func (n *Node) AddVote(vote *models.Vote) error {
	cmd := Command{
		Op:          "vote",
		ElectionID:  vote.ElectionID,
		CandidateID: vote.CandidateID,
	}

	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return n.raft.Apply(data, 500*time.Millisecond).Error()
}

func (n *Node) GetResults(electionID string) *models.ElectionResult {
	n.fsm.state.mu.RLock()
	defer n.fsm.state.mu.RUnlock()

	if counts, exists := n.fsm.state.Elections[electionID]; exists {
		candidates := make([]models.CandidateCount, 0, len(counts))
		for candidateID, count := range counts {
			candidates = append(candidates, models.CandidateCount{
				CandidateID: candidateID,
				Count:       count,
			})
		}
		return &models.ElectionResult{
			ElectionID: electionID,
			Candidates: candidates,
		}
	}
	return nil
}

func (n *Node) GetAllResults() []models.ElectionResult {
	n.fsm.state.mu.RLock()
	defer n.fsm.state.mu.RUnlock()

	results := make([]models.ElectionResult, 0, len(n.fsm.state.Elections))
	for electionID, counts := range n.fsm.state.Elections {
		candidates := make([]models.CandidateCount, 0, len(counts))
		for candidateID, count := range counts {
			candidates = append(candidates, models.CandidateCount{
				CandidateID: candidateID,
				Count:       count,
			})
		}
		results = append(results, models.ElectionResult{
			ElectionID: electionID,
			Candidates: candidates,
		})
	}
	return results
}

func (n *Node) JoinCluster(nodeID, nodeAddr string) error {
	configFuture := n.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return err
	}

	serverID := raft.ServerID(nodeID)
	serverAddr := raft.ServerAddress(nodeAddr)

	for _, srv := range configFuture.Configuration().Servers {
		if srv.ID == serverID || srv.Address == serverAddr {
			if srv.ID == serverID && srv.Address == serverAddr {
				return nil
			}
			removeFuture := n.raft.RemoveServer(serverID, 0, 0)
			if err := removeFuture.Error(); err != nil {
				return err
			}
		}
	}

	addFuture := n.raft.AddVoter(serverID, serverAddr, 0, 0)
	return addFuture.Error()
}
