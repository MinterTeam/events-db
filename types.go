package compact

import "github.com/MinterTeam/minter-go-node/eventsdb/events"

type reward struct {
	Role     byte
	Address  [20]byte //todo: change to AddressID uint32
	Amount   []byte
	PubKeyID uint16
}

func rewardConvert(rewardEvent *events.RewardEvent, ID uint16) interface{} {
	reward := new(reward)
	copy(reward.Address[:], rewardEvent.Address[:])
	reward.Role = byte(rewardEvent.Role)
	reward.Amount = rewardEvent.Amount
	reward.PubKeyID = ID
	return reward
}

type slash struct {
	Address  [20]byte
	Amount   []byte
	Coin     [10]byte
	PubKeyID uint16
}

func convertSlash(rewardEvent *events.SlashEvent, ID uint16) interface{} {
	reward := new(slash)
	copy(reward.Address[:], rewardEvent.Address[:])
	copy(reward.Coin[:], rewardEvent.Coin[:])
	reward.Amount = rewardEvent.Amount
	reward.PubKeyID = ID
	return reward
}

type unbond struct {
	Address  [20]byte
	Amount   []byte
	Coin     [10]byte
	PubKeyID uint16
}

func convertUnbound(rewardEvent *events.UnbondEvent, ID uint16) interface{} {
	reward := new(unbond)
	copy(reward.Address[:], rewardEvent.Address[:])
	copy(reward.Coin[:], rewardEvent.Coin[:])
	reward.Amount = rewardEvent.Amount
	reward.PubKeyID = ID
	return reward
}
