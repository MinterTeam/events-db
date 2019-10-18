package compact

import "github.com/MinterTeam/minter-go-node/eventsdb/events"

type reward struct {
	Role      byte
	AddressID uint32
	Amount    []byte
	PubKeyID  uint16
}

func rewardConvert(rewardEvent *events.RewardEvent, pubKeyID uint16, addressID uint32) interface{} {
	reward := new(reward)
	reward.AddressID = addressID
	reward.Role = byte(rewardEvent.Role)
	reward.Amount = rewardEvent.Amount
	reward.PubKeyID = pubKeyID
	return reward
}

type slash struct {
	AddressID uint32
	Amount    []byte
	Coin      [10]byte
	PubKeyID  uint16
}

func convertSlash(rewardEvent *events.SlashEvent, pubKeyID uint16, addressID uint32) interface{} {
	reward := new(slash)
	reward.AddressID = addressID
	copy(reward.Coin[:], rewardEvent.Coin[:])
	reward.Amount = rewardEvent.Amount
	reward.PubKeyID = pubKeyID
	return reward
}

type unbond struct {
	AddressID uint32
	Amount    []byte
	Coin      [10]byte
	PubKeyID  uint16
}

func convertUnbound(rewardEvent *events.UnbondEvent, pubKeyID uint16, addressID uint32) interface{} {
	reward := new(unbond)
	reward.AddressID = addressID
	copy(reward.Coin[:], rewardEvent.Coin[:])
	reward.Amount = rewardEvent.Amount
	reward.PubKeyID = pubKeyID
	return reward
}
