package delta

import (
	"encoding/json"
	"strings"

	"github.com/r3labs/diff/v3"
)

func j(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// Operation is a JSON Patch operation
type Operation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value,omitempty"`
}

func (op Operation) String() string {
	return j(op)
}

func changesToOperations(changes []diff.Change) []Operation {
	// TODO:
	//  optimise changes to use moves and copies
	ops := make([]Operation, len(changes))
	for i, change := range changes {
		ops[i] = changeToOperation(change)
	}
	return ops
}

func changeToOperation(change diff.Change) Operation {
	op := Operation{}
	op.Path = "/" + strings.Join(change.Path, "/")
	switch change.Type {
	case diff.CREATE:
		op.Op = "add"
		op.Value = change.To
	case diff.UPDATE:
		op.Op = "replace"
		op.Value = change.To
	case diff.DELETE:
		op.Op = "remove"
	}
	return op
}
