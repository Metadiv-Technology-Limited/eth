package eth

type TransactOpts struct {
	Value    uint
	GasLimit uint
	GasPrice uint
}

type DeployContractResult struct {
	Address string
	TxHash  string
	Nonce   uint64
}

type WriteContractResult struct {
	TxHash string
	Nonce  uint64
}

type Abi struct {
	Name            string      `json:"name,omitempty"`
	Type            string      `json:"type,omitempty"`
	StateMutability string      `json:"stateMutability,omitempty"`
	Inputs          []AbiInput  `json:"inputs,omitempty"`
	Outputs         []AbiOutput `json:"outputs,omitempty"`
}

type AbiInput struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Indexed bool   `json:"indexed,omitempty"`
}

type AbiOutput struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}
