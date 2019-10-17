package compact

import (
	"encoding/binary"
	"github.com/MinterTeam/go-amino"
	"github.com/MinterTeam/minter-go-node/eventsdb"
	"github.com/MinterTeam/minter-go-node/eventsdb/events"
	"github.com/tendermint/tm-db"
	"path/filepath"
	"sync"
)

var cdc = amino.NewCodec()

type eventsStore struct {
	sync.RWMutex
	db      db.DB
	pending pendingEvents
	pubKeys map[uint8][32]byte
}

type pendingEvents struct {
	sync.Mutex
	height uint64
	items  events.Events
}

func NewEventsStore() eventsdb.IEventsDB {
	return &eventsStore{db: db.NewDB("events", db.GoLevelDBBackend, filepath.Join(".", "data"))}
}

func (e *eventsStore) AddEvent(height uint64, event events.Event) {
	e.pending.Lock()
	defer e.pending.Unlock()
	if e.pending.height != height {
		e.pending.items = events.Events{}
	}
	e.pending.items = append(e.pending.items, event)
	e.pending.height = height
}

func (e *eventsStore) LoadEvents(height uint64) events.Events {
	panic("implement me")
}

func (e *eventsStore) FlushEvents() error {
	e.pending.Lock()
	defer e.pending.Unlock()
	for _, item := range e.pending.items {
		switch item.(type) {
		case events.RewardEvent:

		case events.UnbondEvent:

		case events.SlashEvent:

		case events.CoinLiquidationEvent:

		}
	}
	var data []byte
	// todo: data =
	e.Lock()
	defer e.Unlock()
	e.db.Set(getKeyForHeight(e.pending.height), data)
	return nil
}

func getKeyForHeight(height uint64) []byte {
	var h = make([]byte, 8)
	binary.BigEndian.PutUint64(h, height)
	return h
}
