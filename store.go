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
	cdc *amino.Codec
	sync.RWMutex
	db       db.DB
	pending  pendingEvents
	pubIDKey map[uint16][32]byte
	pubKeyID map[[32]byte]uint16
}

type pendingEvents struct {
	sync.Mutex
	height uint32
	items  events.Events
}

func NewEventsStore() eventsdb.IEventsDB {
	codec := amino.NewCodec()
	codec.RegisterInterface((*events.Event)(nil), nil)
	codec.RegisterConcrete(reward{},
		"reward", nil)
	codec.RegisterConcrete(slash{},
		"slash", nil)
	codec.RegisterConcrete(unbond{},
		"unbond", nil)
	codec.RegisterConcrete(events.CoinLiquidationEvent{},
		"minter/CoinLiquidationEvent", nil)

	return &eventsStore{
		db:  db.NewDB("events", db.GoLevelDBBackend, filepath.Join(".", "data")),
		cdc: codec,
	}
}

func (store *eventsStore) setPubKey(id uint16, key [32]byte) {
	store.pubIDKey[id] = key
	store.pubKeyID[key] = id
}

func (store *eventsStore) AddEvent(height uint64, event events.Event) {
	store.pending.Lock()
	defer store.pending.Unlock()
	if store.pending.height != uint32(height) {
		store.pending.items = events.Events{}
	}
	store.pending.items = append(store.pending.items, event)
	store.pending.height = uint32(height)
}

func (store *eventsStore) LoadEvents(height uint64) events.Events {
	store.Lock()
	if len(store.pubIDKey) == 0 {
		store.loadKeys()
	}
	store.Unlock()

	//todo: work in progress
	return nil
}

func (store *eventsStore) FlushEvents() error {
	store.pending.Lock()
	defer store.pending.Unlock()
	var data []interface{}
	for _, item := range store.pending.items {
		data = append(data, store.convert(item))
	}

	bytes, err := store.cdc.MarshalBinaryBare(data)
	if err != nil {
		return err
	}

	store.Lock()
	defer store.Unlock()
	store.db.Set(uint32ToBytes(store.pending.height), bytes)
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

func (store *eventsStore) Key(validatorPubKey []byte) uint16 {
	var key [32]byte
	copy(key[:], validatorPubKey)

	if id, ok := store.pubKeyID[key]; ok {
		return id
	}

	id := uint16(len(store.pubIDKey))
	store.setPubKey(id, key)

	store.db.Set(append([]byte(validatorIdsKeyPrefix), uint16ToBytes(id)...), key[:])
	store.db.Set([]byte(validatorsCountKey), uint32ToBytes(uint32(len(store.pubIDKey))))
	return id
}

func (store *eventsStore) loadKeys() {
	for id := uint16(0); id < binary.BigEndian.Uint16(store.db.Get([]byte(validatorsCountKey))); id++ {
		validatorPubKey := store.db.Get(append([]byte(validatorIdsKeyPrefix), uint16ToBytes(id)...))
		var key [32]byte
		copy(key[:], validatorPubKey)
		store.setPubKey(id, key)
	}
}

func uint32ToBytes(height uint32) []byte {
	var h = make([]byte, 4)
	binary.BigEndian.PutUint32(h, height)
	return h
}

func uint16ToBytes(height uint16) []byte {
	var h = make([]byte, 2)
	binary.BigEndian.PutUint16(h, height)
	return h
}
