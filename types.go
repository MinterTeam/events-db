package compact

type CoinLiquidationEvent struct {
	Coin string `json:"coin"`
}

type reward struct {
	Role            byte
	Address         [20]byte
	Amount          []byte
	ValidatorPubKey uint8
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
