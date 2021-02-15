 package ontid

import (
	"encoding/hex"

	" github.com/Daironode/aingle/common"
	" github.com/Daironode/aingle/smartcontract/event"
	" github.com/Daironode/aingle/smartcontract/service/native"
)

func newEvent(srvc *native.NativeService, st interface{}) {
	e := event.NotifyEventInfo{}
	e.ContractAddress = srvc.ContextRef.CurrentContext().ContractAddress
	e.States = st
	srvc.Notifications = append(srvc.Notifications, &e)
}

func triggerRegisterEvent(srvc *native.NativeService, id []byte) {
	newEvent(srvc, []string{"Register", string(id)})
}

func triggerPublicEvent(srvc *native.NativeService, op string, id, pub []byte, keyID uint32) {
	st := []interface{}{"PublicKey", op, string(id), keyID, hex.EncodeToString(pub)}
	newEvent(srvc, st)
}

func triggerAttributeEvent(srvc *native.NativeService, op string, id []byte, path [][]byte) {
	var attr interface{}
	if op == "remove" {
		attr = hex.EncodeToString(path[0])
	} else {
		t := make([]string, len(path))
		for i, v := range path {
			t[i] = hex.EncodeToString(v)
		}
		attr = t
	}
	st := []interface{}{"Attribute", op, string(id), attr}
	newEvent(srvc, st)
}

func triggerRecoveryEvent(srvc *native.NativeService, op string, id []byte, addr common.Address) {
	st := []string{"Recovery", op, string(id), addr.ToHexString()}
	newEvent(srvc, st)
}

func triggerContextEvent(srvc *native.NativeService, op string, id []byte, contexts [][]byte) {
	t := make([]string, len(contexts))
	var c interface{}
	for i := 0; i < len(contexts); i++ {
		t[i] = hex.EncodeToString(contexts[i])
	}
	c = t
	st := []interface{}{"Context", op, string(id), c}
	newEvent(srvc, st)
}

func triggerServiceEvent(srvc *native.NativeService, op string, id []byte, serviceId []byte) {
	st := []string{"Service", op, string(id), common.ToHexString(serviceId)}
	newEvent(srvc, st)
}

func triggerAuthKeyEvent(srvc *native.NativeService, op string, id []byte, keyID uint32) {
	st := []interface{}{"AuthKey", op, string(id), keyID}
	newEvent(srvc, st)
}
