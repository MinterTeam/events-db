package compact

import "github.com/MinterTeam/minter-go-node/eventsdb/events"

type CoinLiquidationEvent struct {
	Coin string `json:"coin"`
}

type reward struct {
	Role    byte
	Address [20]byte
	Amount  []byte
	ID      uint8
}

func rewardConvert(rewardEvent *events.RewardEvent, ID uint8) interface{} {
	reward := new(reward)
	copy(reward.Address[:], rewardEvent.Address[:])
	reward.Role = byte(rewardEvent.Role)
	reward.Amount = rewardEvent.Amount
	reward.ID = ID
	return reward
}

type slash struct {
	Address [20]byte
	Amount  []byte
	Coin    [10]byte
	ID      uint8
}

type unbond struct {
	Address [20]byte
	Amount  []byte
	Coin    [10]byte
	ID      uint8
}
