package counter

import (
	"encoding/json"

	"github.com/hashicorp/raft"
)

type Snapshot struct {
	elections map[string]map[string]int
}

func (s *Snapshot) Persist(sink raft.SnapshotSink) error {
	err := json.NewEncoder(sink).Encode(s.elections)
	if err != nil {
		sink.Cancel()
		return err
	}
	return sink.Close()
}

func (s *Snapshot) Release() {}
