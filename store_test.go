package compact_db

import (
	"encoding/hex"
	"github.com/MinterTeam/go-amino"
	"github.com/MinterTeam/minter-go-node/eventsdb/events"
	db "github.com/tendermint/tm-db"
	"math/big"
	"path/filepath"
	"sync"
	"testing"
)

func TestIEventsDB(t *testing.T) {
	codec := amino.NewCodec()
	codec.RegisterInterface((*interface{})(nil), nil)
	codec.RegisterConcrete(reward{}, "reward", nil)
	codec.RegisterConcrete(slash{}, "slash", nil)
	codec.RegisterConcrete(unbond{}, "unbond", nil)
	codec.RegisterConcrete(events.CoinLiquidationEvent{}, "minter/CoinLiquidationEvent", nil)

	eventsDB := IEventsDB(&eventsStore{
		cdc:       codec,
		RWMutex:   sync.RWMutex{},
		db:        db.NewDB("events", db.GoLevelDBBackend, filepath.Join(".", "data_test")),
		pending:   pendingEvents{},
		idPubKey:  make(map[uint16]string),
		pubKeyID:  make(map[string]uint16),
		idAddress: make(map[uint32][20]byte),
		addressID: make(map[[20]byte]uint32),
	})
	{
		amount, _ := big.NewInt(0).SetString("111497225000000000000", 10)
		event := events.RewardEvent{
			Role:            events.RoleDevelopers,
			Address:         [20]byte{},
			Amount:          amount.Bytes(),
			ValidatorPubKey: []byte{},
		}
		bytesAddress, err := hex.DecodeString("Mx04bea23efb744dc93b4fda4c20bf4a21c6e195f1"[2:])
		if err != nil {
			t.Fatal(err)
		}
		copy(event.Address[:], bytesAddress)
		bytesPubKey, err := hex.DecodeString("Mp9e13f2f5468dd782b316444fbd66595e13dba7d7bd3efa1becd50b42045f58c6"[2:])
		if err != nil {
			t.Fatal(err)
		}
		event.ValidatorPubKey = bytesPubKey
		eventsDB.AddEvent(12, event)
	}
	{
		amount, _ := big.NewInt(0).SetString("891977800000000000000", 10)
		event := events.RewardEvent{
			Role:            events.RoleValidator,
			Address:         [20]byte{},
			Amount:          amount.Bytes(),
			ValidatorPubKey: []byte{},
		}
		bytesAddress, err := hex.DecodeString("Mx18467bbb64a8edf890201d526c35957d82be3d95"[2:])
		if err != nil {
			t.Fatal(err)
		}
		copy(event.Address[:], bytesAddress)
		bytesPubKey, err := hex.DecodeString("Mp738da41ba6a7b7d69b7294afa158b89c5a1b410cbf0c2443c85c5fe24ad1dd1c"[2:])
		if err != nil {
			t.Fatal(err)
		}
		event.ValidatorPubKey = bytesPubKey
		eventsDB.AddEvent(12, event)
	}
	err := eventsDB.CommitEvents()
	if err != nil {
		t.Fatal(err)
	}

	loadEvents := eventsDB.LoadEvents(12)
	for _, v := range loadEvents {
		t.Logf("%+v", v)
		t.Logf("%+v", big.NewInt(0).SetBytes(v.(*events.RewardEvent).Amount).String())
		t.Logf("%+v", v.(*events.RewardEvent).Address.String())
	}
}
