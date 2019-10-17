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

func (store *eventsStore) AddEvent(height uint64, event events.Event) {
	store.pending.Lock()
	defer store.pending.Unlock()
	if store.pending.height != height {
		store.pending.items = events.Events{}
	}
	store.pending.items = append(store.pending.items, event)
	store.pending.height = height
}

func (store *eventsStore) LoadEvents(height uint64) events.Events {
	panic("implement me")
}

func (store *eventsStore) FlushEvents() error {
	store.pending.Lock()
	defer store.pending.Unlock()
	var data []interface{}
	for _, item := range store.pending.items {
		data = append(data, store.convert(item))
	}

	bytes, err := cdc.MarshalBinaryBare(data)
	if err != nil {
		return err
	}

	store.Lock()
	defer store.Unlock()
	store.db.Set(getKeyForHeight(store.pending.height), bytes)
	return nil
}

func (store *eventsStore) convert(event events.Event) interface{} {
	var res interface{}
	switch event.(type) {
	case *events.RewardEvent:
		res = store.convertReward(event.(*events.RewardEvent))
	case *events.UnbondEvent:
		// todo: res = store.convertUnbound(event.(*events.UnbondEvent))
	case *events.SlashEvent:
		// todo: res = store.convertSlash(event.(*events.SlashEvent))
	case *events.CoinLiquidationEvent:

	}
	return res
}

func (store *eventsStore) convertReward(rewardEvent *events.RewardEvent) interface{} {
	return rewardConvert(rewardEvent, store.Key(rewardEvent))
}

func (store *eventsStore) Key(rewardEvent *events.RewardEvent) uint8 {
	for id, v := range store.pubKeys {
		if string(v[:]) == string(rewardEvent.ValidatorPubKey) {
			return id
		}
	}
	var key [32]byte
	copy(key[:], rewardEvent.ValidatorPubKey)
	ID := uint8(len(store.pubKeys))
	store.pubKeys[ID] = key
	return ID
}

func getKeyForHeight(height uint64) []byte {
	var h = make([]byte, 8)
	binary.BigEndian.PutUint64(h, height)
	return h
}
