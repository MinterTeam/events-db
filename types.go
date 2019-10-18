package compact

import "github.com/MinterTeam/minter-go-node/eventsdb/events"

type reward struct {
	Role      byte
	AddressID uint32
	Amount    []byte
	PubKeyID  uint16
}

func rewardConvert(event *events.RewardEvent, pubKeyID uint16, addressID uint32) interface{} {
	result := new(reward)
	result.AddressID = addressID
	result.Role = byte(event.Role)
	result.Amount = event.Amount
	result.PubKeyID = pubKeyID
	return result
}

func compileReward(item *reward, pubKey string, address [20]byte) interface{} {
	event := new(events.RewardEvent)
	event.ValidatorPubKey = []byte(pubKey)
	copy(event.Address[:], address[:])
	event.Role = events.Role(item.Role)
	event.Amount = item.Amount
	return event
}

type slash struct {
	AddressID uint32
	Amount    []byte
	Coin      [10]byte
	PubKeyID  uint16
}

func convertSlash(event *events.SlashEvent, pubKeyID uint16, addressID uint32) interface{} {
	result := new(slash)
	result.AddressID = addressID
	copy(result.Coin[:], event.Coin[:])
	result.Amount = event.Amount
	result.PubKeyID = pubKeyID
	return result
}

func compileSlash(item *slash, pubKey string, address [20]byte) interface{} {
	event := new(events.SlashEvent)
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

func convertUnbound(event *events.UnbondEvent, pubKeyID uint16, addressID uint32) interface{} {
	result := new(unbond)
	result.AddressID = addressID
	copy(result.Coin[:], event.Coin[:])
	result.Amount = event.Amount
	result.PubKeyID = pubKeyID
	return result
}

func compileUnbond(item *unbond, pubKey string, address [20]byte) interface{} {
	event := new(events.UnbondEvent)
	event.ValidatorPubKey = []byte(pubKey)
	copy(event.Address[:], address[:])
	copy(event.Coin[:], item.Coin[:])
	event.Amount = item.Amount
	return event
}
