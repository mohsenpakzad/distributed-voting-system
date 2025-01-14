package counter

import (
	"encoding/json"
	"io"

	"github.com/hashicorp/raft"
)

type FSM struct {
	state *VoteState
}

func NewFSM() *FSM {
	return &FSM{
		state: &VoteState{
			Elections: make(map[string]map[string]int),
		},
	}
}

func (f *FSM) Apply(log *raft.Log) interface{} {
	var cmd Command
	if err := json.Unmarshal(log.Data, &cmd); err != nil {
		return err
	}

	switch cmd.Op {
	case "vote":
		f.state.mu.Lock()
		if f.state.Elections[cmd.ElectionID] == nil {
			f.state.Elections[cmd.ElectionID] = make(map[string]int)
		}
		f.state.Elections[cmd.ElectionID][cmd.CandidateID]++
		f.state.mu.Unlock()
	}
	return nil
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	f.state.mu.RLock()
	defer f.state.mu.RUnlock()

	// Create a deep copy of the state
	elections := make(map[string]map[string]int)
	for electionID, candidates := range f.state.Elections {
		elections[electionID] = make(map[string]int)
		for candidateID, count := range candidates {
			elections[electionID][candidateID] = count
		}
	}

	return &Snapshot{elections: elections}, nil
}

func (f *FSM) Restore(rc io.ReadCloser) error {
	defer rc.Close()

	elections := make(map[string]map[string]int)
	if err := json.NewDecoder(rc).Decode(&elections); err != nil {
		return err
	}

	f.state.mu.Lock()
	f.state.Elections = elections
	f.state.mu.Unlock()

	return nil
}
