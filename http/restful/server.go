 
// Package restful privides a function to start restful server
package restful

import (
	" github.com/Daironode/aingle/http/restful/restful"
)

//start restful
func StartServer() {
	rt := restful.InitRestServer()
	go rt.Start()
}
