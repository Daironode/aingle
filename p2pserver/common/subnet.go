
package common

type SubnetMemberInfo struct {
	PubKey     string `json:"pubKey"`
	ListenAddr string `json:"listenAddr"`
	Connected  bool   `json:"connected"`
	Height     uint64 `json:"height"`
	Version    string `json:"version"`
}
