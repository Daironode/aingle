 package common

import (
	"github.com/Daironode/aingle-crypto/keypair"
	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/core/payload"
	" github.com/Daironode/aingle/core/types"
)

type PayloadInfo interface{}

type InvokeCodeInfo struct {
	Code string
}
type DeployCodeInfo struct {
	Code        string
	VmType      byte
	Name        string
	CodeVersion string
	Author      string
	Email       string
	Description string
}

type BookkeeperInfo struct {
	PubKey     string
	Action     string
	Issuer     string
	Controller string
}

//get tranasction payload data
func TransPayloadToHex(p types.Payload) PayloadInfo {
	switch object := p.(type) {
	case *payload.Bookkeeper:
		obj := new(BookkeeperInfo)
		pubKeyBytes := keypair.SerializePublicKey(object.PubKey)
		obj.PubKey = common.ToHexString(pubKeyBytes)
		if object.Action == payload.BookkeeperAction_ADD {
			obj.Action = "add"
		} else if object.Action == payload.BookkeeperAction_SUB {
			obj.Action = "sub"
		} else {
			obj.Action = "nil"
		}
		pubKeyBytes = keypair.SerializePublicKey(object.Issuer)
		obj.Issuer = common.ToHexString(pubKeyBytes)

		return obj
	case *payload.InvokeCode:
		obj := new(InvokeCodeInfo)
		obj.Code = common.ToHexString(object.Code)
		return obj
	case *payload.DeployCode:
		obj := new(DeployCodeInfo)
		obj.Code = common.ToHexString(object.GetRawCode())
		obj.VmType = byte(object.VmType())
		obj.Name = object.Name
		obj.CodeVersion = object.Version
		obj.Author = object.Author
		obj.Email = object.Email
		obj.Description = object.Description
		return obj
	}
	return nil
}
