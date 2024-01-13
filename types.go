package main

/*
{
  "isCurrentSequencer": 0,
  "writer": {
    "txHash": "0x6dabe8aa50b97f5eb58e2104c95995e4c5e2fbe7ea52051160f4005c8be2f7dc",
    "timestamp": 1704808290,
    "ret": ""
  },
  "l2": { "blockNumber": 461877, "timestamp": 1704809059 },
  "rpc": {
    "txHash": "0xfbfc31e944092c201f2dcac88264a808065afbf36521548f53ae9d13298f7a38",
    "timestamp": 1704809040
  },
  "mpc": {
    "isMpcProposer": 1,
    "signId": "07e3924b-f390-442b-b54c-a6d13b63be87",
    "signSuccess": 1,
    "timestamp": 1704808945
  }
}
*/

type WriterState struct {
	TxHash    string  `json:"txHash"`
	Timestamp float64 `json:"timestamp"`
	Ret       string  `json:"ret"`
}

type MPCState struct {
	IsMpcProposer int     `json:"isMpcProposer"`
	SignId        string  `json:"signId"`
	SignSuccess   int     `json:"signSuccess"`
	Timestamp     float64 `json:"timestamp"`
}

type L2State struct {
	BlockNumber float64 `json:"blockNumber"`
	Timestamp   float64 `json:"timestamp"`
}

type RPCState struct {
	TxHash    string  `json:"txHash"`
	Timestamp float64 `json:"timestamp"`
}

type NodeHealthResp struct {
	L2  L2State   `json:"l2"`
	RPC RPCState  `json:"rpc"`
	MPC *MPCState `json:"mpc"`
}

type LastestSpanResp struct {
	Height float64 `json:"height,string"`
}
