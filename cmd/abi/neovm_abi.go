
package abi

import "strings"

const (
	NEOVM_PARAM_TYPE_BOOL       = "boolean"
	NEOVM_PARAM_TYPE_STRING     = "string"
	NEOVM_PARAM_TYPE_INTEGER    = "integer"
	NEOVM_PARAM_TYPE_ARRAY      = "array"
	NEOVM_PARAM_TYPE_BYTE_ARRAY = "bytearray"
	NEOVM_PARAM_TYPE_VOID       = "void"
	NEOVM_PARAM_TYPE_ANY        = "any"
)

type NeovmContractAbi struct {
	Address    string                      `json:"hash"`
	EntryPoint string                      `json:"entrypoint"`
	Functions  []*NeovmContractFunctionAbi `json:"functions"`
	Events     []*NeovmContractEventAbi    `json:"events"`
}

func (this *NeovmContractAbi) GetFunc(method string) *NeovmContractFunctionAbi {
	method = strings.ToLower(method)
	for _, funcAbi := range this.Functions {
		if strings.ToLower(funcAbi.Name) == method {
			return funcAbi
		}
	}
	return nil
}

func (this *NeovmContractAbi) GetEvent(evt string) *NeovmContractEventAbi {
	evt = strings.ToLower(evt)
	for _, evtAbi := range this.Events {
		if strings.ToLower(evtAbi.Name) == evt {
			return evtAbi
		}
	}
	return nil
}

type NeovmContractFunctionAbi struct {
	Name       string                    `json:"name"`
	Parameters []*NeovmContractParamsAbi `json:"parameters"`
	ReturnType string                    `json:"returntype"`
}

type NeovmContractParamsAbi struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type NeovmContractEventAbi struct {
	Name       string                    `json:"name"`
	Parameters []*NeovmContractParamsAbi `json:"parameters"`
	ReturnType string                    `json:"returntype"`
}
