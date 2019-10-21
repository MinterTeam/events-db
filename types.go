package compact_db

import (
	"github.com/MinterTeam/go-amino"
	"github.com/MinterTeam/minter-go-node/core/types"
)

func RegisterAminoEvents(codec *amino.Codec) {
	codec.RegisterInterface((*Event)(nil), nil)
	codec.RegisterConcrete(RewardEvent{},
		"RewardEvent", nil)
	codec.RegisterConcrete(SlashEvent{},
		"SlashEvent", nil)
	codec.RegisterConcrete(UnbondEvent{},
		"UnbondEvent", nil)
	codec.RegisterConcrete(CoinLiquidationEvent{},
		"CoinLiquidationEvent", nil)
}

type Event interface{}
type Events []Event

type Role byte

const (
	RoleValidator Role = iota
	RoleDelegator
	RoleDAO
	RoleDevelopers
)

func (r Role) String() string {
	switch r {
	case RoleValidator:
		return "Validator"
	case RoleDelegator:
		return "Delegator"
	case RoleDAO:
		return "DAO"
	case RoleDevelopers:
		return "Developers"
	}

	return "Undefined"
}

type reward struct {
	Role      byte
	AddressID uint32
	Amount    []byte
	PubKeyID  uint16
}

type RewardEvent struct {
	Role            Role
	Address         types.Address
	Amount          []byte
	ValidatorPubKey types.Pubkey
}

func rewardConvert(event *RewardEvent, pubKeyID uint16, addressID uint32) interface{} {
	result := new(reward)
	result.AddressID = addressID
	result.Role = byte(event.Role)
	result.Amount = event.Amount
	result.PubKeyID = pubKeyID
	return result
}

func compileReward(item *reward, pubKey string, address [20]byte) interface{} {
	event := new(RewardEvent)
	event.ValidatorPubKey = []byte(pubKey)
	copy(event.Address[:], address[:])
	event.Role = Role(item.Role)
	event.Amount = item.Amount
	return event
}

type slash struct {
	AddressID uint32
	Amount    []byte
	Coin      [10]byte
	PubKeyID  uint16
}

type SlashEvent struct {
	Address         types.Address
	Amount          []byte
	Coin            types.CoinSymbol
	ValidatorPubKey types.Pubkey
}

func convertSlash(event *SlashEvent, pubKeyID uint16, addressID uint32) interface{} {
	result := new(slash)
	result.AddressID = addressID
	copy(result.Coin[:], event.Coin[:])
	result.Amount = event.Amount
	result.PubKeyID = pubKeyID
	return result
}

func compileSlash(item *slash, pubKey string, address [20]byte) interface{} {
	event := new(SlashEvent)
	event.ValidatorPubKey = []byte(pubKey)
	copy(event.Address[:], address[:])
	copy(event.Coin[:], item.Coin[:])
	event.Amount = item.Amount
	return event
}

type unbond struct {
	AddressID uint32
	Amount    []byte
	Coin      [10]byte
	PubKeyID  uint16
}

type UnbondEvent struct {
	Address         types.Address
	Amount          []byte
	Coin            types.CoinSymbol
	ValidatorPubKey types.Pubkey
}

func convertUnbound(event *UnbondEvent, pubKeyID uint16, addressID uint32) interface{} {
	result := new(unbond)
	result.AddressID = addressID
	copy(result.Coin[:], event.Coin[:])
	result.Amount = event.Amount
	result.PubKeyID = pubKeyID
	return result
}

func compileUnbond(item *unbond, pubKey string, address [20]byte) interface{} {
	event := new(UnbondEvent)
	event.ValidatorPubKey = []byte(pubKey)
	copy(event.Address[:], address[:])
	copy(event.Coin[:], item.Coin[:])
	event.Amount = item.Amount
	return event
}

type CoinLiquidationEvent struct {
	Coin types.CoinSymbol
}
