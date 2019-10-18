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
	store.Lock()
	if len(store.pubKeys) == 0 {
		store.loadKeys()
	}
	store.Unlock()

}

func (store *eventsStore) FlushEvents() error {
	store.pending.Lock()
	defer store.pending.Unlock()
	var data []interface{}
	for _, item := range store.pending.items {
		data = append(data, store.convert(item))
	}

	bytes, err := amino.MarshalBinaryBare(data)
	if err != nil {
		return err
	}

	store.Lock()
	defer store.Unlock()
	store.db.Set(uint2byte(store.pending.height), bytes)
	return nil
}

func (store *eventsStore) convert(event events.Event) interface{} {
	var res interface{}
	switch event.(type) {
	case *events.RewardEvent:
		res = store.convertReward(event.(*events.RewardEvent))
	case *events.UnbondEvent:
		res = store.convertUnbound(event.(*events.UnbondEvent))
	case *events.SlashEvent:
		res = store.convertSlash(event.(*events.SlashEvent))
	//case *events.CoinLiquidationEvent:
	default:
		res = event
	}
	return res
}

func (store *eventsStore) convertReward(rewardEvent *events.RewardEvent) interface{} {
	return rewardConvert(rewardEvent, store.Key(rewardEvent.ValidatorPubKey[:]))
}

func (store *eventsStore) convertUnbound(unbondEvent *events.UnbondEvent) interface{} {
	return convertUnbound(unbondEvent, store.Key(unbondEvent.ValidatorPubKey[:]))
}

func (store *eventsStore) convertSlash(slashEvent *events.SlashEvent) interface{} {
	return convertSlash(slashEvent, store.Key(slashEvent.ValidatorPubKey[:]))
}

const validatorIdsKeyPrefix = "validator"
const validatorsCountKey = "validators"

func (store *eventsStore) Key(validatorPubKey []byte) uint8 {
	{
		strKey := string(validatorPubKey)
		for id, v := range store.pubKeys {
			if string(v[:]) == strKey {
				return id
			}
		}
	}
	var key [32]byte
	copy(key[:], validatorPubKey)
	id := uint8(len(store.pubKeys))
	store.pubKeys[id] = key

	store.db.Set(append([]byte(validatorIdsKeyPrefix), uint2byte(uint64(id))...), key[:])
	store.db.Set([]byte(validatorsCountKey), uint2byte(uint64(len(store.pubKeys))))
	return id
}

func (store *eventsStore) loadKeys() {
	count := store.db.Get([]byte(validatorsCountKey))
	for i := uint64(0); i < binary.BigEndian.Uint64(count); i++ {
		validatorPubKey := store.db.Get(append([]byte(validatorIdsKeyPrefix), uint2byte(i)...))
		var key [32]byte
		copy(key[:], validatorPubKey)
		store.pubKeys[uint8(i)] = key
	}
}

func uint2byte(height uint64) []byte {
	var h = make([]byte, 8)
	binary.BigEndian.PutUint64(h, height)
	return h
}
