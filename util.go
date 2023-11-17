package eth

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func ParseArgs(args []string, abi Abi) ([]any, error) {
	parsedArgs := make([]any, 0)

	if len(args) != len(abi.Inputs) {
		return nil, errors.New("args length does not match abi inputs length")
	}

	for i, arg := range args {
		abiArg := abi.Inputs[i]
		switch abiArg.Type {
		case "address":
			address := common.HexToAddress(arg)
			parsedArgs = append(parsedArgs, address)
		case "string":
			parsedArgs = append(parsedArgs, arg)
		case "bool":
			parsedArgs = append(parsedArgs, arg == "true")
		case "[]byte":
			parsedArgs = append(parsedArgs, common.FromHex(arg))
		default:
			n, ok := big.NewInt(0).SetString(arg, 10)
			if !ok {
				return nil, errors.New("failed to parse arg as big int")
			}
			parsedArgs = append(parsedArgs, n)
		}
	}

	return parsedArgs, nil
}
