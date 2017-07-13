package main

import (
	"encoding/json"
	"fmt"

	"github.com/apex/go-apex"
)

type request struct {
	// TODO
}

func main() {
	apex.HandleFunc(func(msg json.RawMessage, lambdaContext *apex.Context) (interface{}, error) {
		return nil, fmt.Errorf("not implemented")
	})
}
