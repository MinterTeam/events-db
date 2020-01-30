package compact_db

import (
	"github.com/MinterTeam/minter-go-node/core/types"
	"github.com/tendermint/go-amino"
	"math/big"
)

func RegisterAminoEvents(codec *amino.Codec) {
	codec.RegisterInterface((*Event)(nil), nil)
	codec.RegisterConcrete(RewardEvent{},
		"minter/RewardEvent", nil)
	codec.RegisterConcrete(SlashEvent{},
		"minter/SlashEvent", nil)
	codec.RegisterConcrete(UnbondEvent{},
		"minter/UnbondEvent", nil)
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
	Role            Role          `json:"role"`
	Address         types.Address `json:"address"`
	Amount          string        `json:"amount"`
	ValidatorPubKey types.Pubkey  `json:"validator_pub_key"`
}

func rewardConvert(event *RewardEvent, pubKeyID uint16, addressID uint32) interface{} {
	result := new(reward)
	result.AddressID = addressID
	result.Role = byte(event.Role)
	bi, _ := big.NewInt(0).SetString(event.Amount, 10)
	result.Amount = bi.Bytes()
	result.PubKeyID = pubKeyID
	return result
}

func compileReward(item *reward, pubKey string, address [20]byte) interface{} {
	event := new(RewardEvent)
	copy(event.ValidatorPubKey[:], pubKey)
	copy(event.Address[:], address[:])
	event.Role = Role(item.Role)
	event.Amount = big.NewInt(0).SetBytes(item.Amount).String()
	return event
}

type slash struct {
	AddressID uint32
	Amount    []byte
	Coin      [10]byte
	PubKeyID  uint16
}

type SlashEvent struct {
	Address         types.Address    `json:"address"`
	Amount          string           `json:"amount"`
	Coin            types.CoinSymbol `json:"coin"`
	ValidatorPubKey types.Pubkey     `json:"validator_pub_key"`
}

func convertSlash(event *SlashEvent, pubKeyID uint16, addressID uint32) interface{} {
	result := new(slash)
	result.AddressID = addressID
	copy(result.Coin[:], event.Coin[:])
	bi, _ := big.NewInt(0).SetString(event.Amount, 10)
	result.Amount = bi.Bytes()
	result.PubKeyID = pubKeyID
	return result
}

func compileSlash(item *slash, pubKey string, address [20]byte) interface{} {
	event := new(SlashEvent)
	copy(event.ValidatorPubKey[:], pubKey)
	copy(event.Address[:], address[:])
	copy(event.Coin[:], item.Coin[:])
	event.Amount = big.NewInt(0).SetBytes(item.Amount).String()
	return event
}

type unbond struct {
	AddressID uint32
	Amount    []byte
	Coin      [10]byte
	PubKeyID  uint16
}

type UnbondEvent struct {
	Address         types.Address    `json:"address"`
	Amount          string           `json:"amount"`
	Coin            types.CoinSymbol `json:"coin"`
	ValidatorPubKey types.Pubkey     `json:"validator_pub_key"`
}

func convertUnbound(event *UnbondEvent, pubKeyID uint16, addressID uint32) interface{} {
	result := new(unbond)
	result.AddressID = addressID
	copy(result.Coin[:], event.Coin[:])
	bi, _ := big.NewInt(0).SetString(event.Amount, 10)
	result.Amount = bi.Bytes()
	result.PubKeyID = pubKeyID
	return result
}

func compileUnbond(item *unbond, pubKey string, address [20]byte) interface{} {
	event := new(UnbondEvent)
	copy(event.ValidatorPubKey[:], pubKey)
	copy(event.Address[:], address[:])
	copy(event.Coin[:], item.Coin[:])
	event.Amount = big.NewInt(0).SetBytes(item.Amount).String()
	return event
}
