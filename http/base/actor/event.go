
package actor

import (
	"github.com/Daironode/aingle-event/actor"
	" github.com/Daironode/aingle/events"
	" github.com/Daironode/aingle/events/message"
)

type EventActor struct {
	blockPersistCompleted func(v interface{})
	smartCodeEvt          func(v interface{})
}

//receive from subscribed actor
func (t *EventActor) Receive(c actor.Context) {
	switch msg := c.Message().(type) {
	case *message.SaveBlockCompleteMsg:
		t.blockPersistCompleted(*msg.Block)
	case *message.SmartCodeEventMsg:
		t.smartCodeEvt(*msg.Event)
	default:
	}
}

//Subscribe save block complete and smartcontract Event
func SubscribeEvent(topic string, handler func(v interface{})) {
	var props = actor.FromProducer(func() actor.Actor {
		if topic == message.TOPIC_SAVE_BLOCK_COMPLETE {
			return &EventActor{blockPersistCompleted: handler}
		} else if topic == message.TOPIC_SMART_CODE_EVENT {
			return &EventActor{smartCodeEvt: handler}
		} else {
			return &EventActor{}
		}
	})
	var pid = actor.Spawn(props)
	var sub = events.NewActorSubscriber(pid)
	sub.Subscribe(topic)
}
