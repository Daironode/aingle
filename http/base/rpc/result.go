
package rpc

import (
	Err " github.com/Daironode/aingle/http/base/error"
)

func ResponseSuccess(result interface{}) map[string]interface{} {
	return ResponsePack(Err.SUCCESS, result)
}
func ResponsePack(errcode int64, result interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"error":  errcode,
		"desc":   Err.ErrMap[errcode],
		"result": result,
	}
	return resp
}
