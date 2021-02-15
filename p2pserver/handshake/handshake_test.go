 package handshake

import (
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"

	" github.com/Daironode/aingle/p2pserver/common"
	" github.com/Daironode/aingle/p2pserver/message/types"
	" github.com/Daironode/aingle/p2pserver/peer"
	"github.com/stretchr/testify/assert"
)

func init() {
	common.Difficulty = 1
	HANDSHAKE_DURATION = 1 * time.Second
}

type Node struct {
	Id   *common.PeerKeyId
	Info *peer.PeerInfo
	Conn net.Conn
}

func NewNode(conn net.Conn) Node {
	node := Node{
		Id:   common.RandPeerKeyId(),
		Info: &peer.PeerInfo{},
		Conn: conn,
	}
	node.Info.Id = node.Id.Id
	node.Info.SoftVersion = "v1.9.0-beta"

	return node
}

func NewPair() (client Node, server Node) {
	c, s := net.Pipe()

	client = NewNode(c)
	server = NewNode(s)
	return
}

func TestHandshakeNormal(t *testing.T) {
	client, server := NewPair()
	versions := []string{"v1.8.0", "v1.7.0", "v1.9.0", "v1.9.0-beta", "v1.20"}

	for i := 0; i < 100; i++ {
		client.Info.SoftVersion = versions[rand.Intn(len(versions))]
		server.Info.SoftVersion = versions[rand.Intn(len(versions))]

		wg := sync.WaitGroup{}
		wg.Add(2)
		result := make([]struct {
			info [2]*peer.PeerInfo
			err  error
		}, 2)
		go func() {
			info, err := HandshakeClient(client.Info, client.Id, client.Conn)
			result[0].err = err
			result[0].info = [2]*peer.PeerInfo{info, server.Info}
			wg.Done()
		}()
		go func() {
			info, err := HandshakeServer(server.Info, server.Id, server.Conn)
			result[1].err = err
			result[1].info = [2]*peer.PeerInfo{info, client.Info}
			wg.Done()
		}()
		wg.Wait()

		for _, res := range result {
			assert.Nil(t, res.err)
			assert.Equal(t, res.info[0].Id.ToUint64(), res.info[1].Id.ToUint64())
		}
	}
}

func TestHandshakeTimeout(t *testing.T) {
	client, _ := NewPair()

	_, err := HandshakeClient(client.Info, client.Id, client.Conn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "deadline exceeded")
}

func TestHandshakeWrongMsg(t *testing.T) {
	client, server := NewPair()
	go func() {
		err := sendMsg(client.Conn, &types.Addr{})
		assert.Nil(t, err)
	}()

	_, err := HandshakeServer(server.Info, server.Id, server.Conn)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "expected version message")
}

func TestVersion(t *testing.T) {
	assert.True(t, supportDHT(common.MIN_VERSION_FOR_DHT))
	assert.True(t, supportDHT("1.9.1"))
	assert.True(t, supportDHT("v1.10.0"))
	assert.True(t, supportDHT("v1.10"))
	assert.True(t, supportDHT("v2.0"))
	assert.True(t, supportDHT("v1.9.1"))
	assert.True(t, supportDHT("1.9.1-beta"))
	assert.True(t, supportDHT("v1.9.1-beta"))
	assert.True(t, supportDHT("1.9.1-beta-9"))
	assert.True(t, supportDHT("1.9.1-beta-9-geeaeewwf"))

	assert.False(t, supportDHT("1.9.1-alpha"))
	assert.False(t, supportDHT("1.8.0-beta-9-geeaeewwf"))
	assert.False(t, supportDHT("1.8.0"))
}
