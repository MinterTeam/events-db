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
	db        db.DB
	pending   pendingEvents
	idPubKey  map[uint16][32]byte
	pubKeyID  map[[32]byte]uint16
	idAddress map[uint32][20]byte
	addressID map[[20]byte]uint32
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

func (store *eventsStore) cachePubKey(id uint16, key [32]byte) {
	store.idPubKey[id] = key
	store.pubKeyID[key] = id
}

func (store *eventsStore) cacheAddress(id uint32, address [20]byte) {
	store.idAddress[id] = address
	store.addressID[address] = id
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
	store.loadCache()

	//todo: work in progress
	return nil
}

func (store *eventsStore) FlushEvents() error {
	store.loadCache()

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

func (store *eventsStore) loadCache() {
	store.Lock()
	if len(store.idPubKey) == 0 {
		store.loadPubKeys()
		store.loadAddresses()
	}
	store.Unlock()
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
	return rewardConvert(rewardEvent, store.savePubKey(rewardEvent.ValidatorPubKey[:]), store.saveAddress(rewardEvent.Address))
}

func (store *eventsStore) convertUnbound(unbondEvent *events.UnbondEvent) interface{} {
	return convertUnbound(unbondEvent, store.savePubKey(unbondEvent.ValidatorPubKey[:]), store.saveAddress(unbondEvent.Address))
}

func (store *eventsStore) convertSlash(slashEvent *events.SlashEvent) interface{} {
	return convertSlash(slashEvent, store.savePubKey(slashEvent.ValidatorPubKey[:]), store.saveAddress(slashEvent.Address))
}

const pubKeyPrefix = "pubKey"
const addressPrefix = "address"
const pubKeysCountKey = "pubKeys"
const addressesCountKey = "addresses"

func (store *eventsStore) saveAddress(address [20]byte) uint32 {

	if id, ok := store.addressID[address]; ok {
		return id
	}

	id := uint32(len(store.idPubKey))
	store.cacheAddress(id, address)

	store.db.Set(append([]byte(addressPrefix), uint32ToBytes(id)...), address[:])
	store.db.Set([]byte(addressesCountKey), uint32ToBytes(uint32(len(store.idPubKey))))
	return id
}

func (store *eventsStore) savePubKey(validatorPubKey []byte) uint16 {
	var key [32]byte
	copy(key[:], validatorPubKey)

	if id, ok := store.pubKeyID[key]; ok {
		return id
	}

	id := uint16(len(store.idPubKey))
	store.cachePubKey(id, key)

	store.db.Set(append([]byte(pubKeyPrefix), uint16ToBytes(id)...), key[:])
	store.db.Set([]byte(pubKeysCountKey), uint32ToBytes(uint32(len(store.idPubKey))))
	return id
}

func (store *eventsStore) loadPubKeys() {
	for id := uint16(0); id < binary.BigEndian.Uint16(store.db.Get([]byte(pubKeysCountKey))); id++ {
		pubKey := store.db.Get(append([]byte(pubKeyPrefix), uint16ToBytes(id)...))
		var key [32]byte
		copy(key[:], pubKey)
		store.cachePubKey(id, key)
	}
}

func (store *eventsStore) loadAddresses() {
	for id := uint32(0); id < binary.BigEndian.Uint32(store.db.Get([]byte(addressesCountKey))); id++ {
		address := store.db.Get(append([]byte(addressPrefix), uint32ToBytes(id)...))
		var key [20]byte
		copy(key[:], address)
		store.cacheAddress(id, key)
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
