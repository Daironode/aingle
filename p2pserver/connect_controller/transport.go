
package connect_controller

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
	"strconv"
	"time"

	" github.com/Daironode/aingle/common/config"
	" github.com/Daironode/aingle/p2pserver/common"
)

type Dialer interface {
	Dial(nodeAddr string) (net.Conn, error)
}

func NewDialer(config *config.P2PNodeConfig) (Dialer, error) {
	if config.IsTLS {
		return newTlsDialer(config)
	}

	return &noTlsDialer{}, nil
}

type tlsDialer struct {
	config *tls.Config
}

func newTlsDialer(config *config.P2PNodeConfig) (*tlsDialer, error) {
	clientCertPool := x509.NewCertPool()
	cacert, err := ioutil.ReadFile(config.CAPath)
	if err != nil {
		return nil, err
	}
	cert, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
	if err != nil {
		return nil, err
	}

	ret := clientCertPool.AppendCertsFromPEM(cacert)
	if !ret {
		return nil, errors.New("[p2p]failed to parse root certificate")
	}

	conf := &tls.Config{
		RootCAs:      clientCertPool,
		Certificates: []tls.Certificate{cert},
	}

	return &tlsDialer{config: conf}, nil
}

func (self *tlsDialer) Dial(nodeAddr string) (net.Conn, error) {
	var dialer net.Dialer
	dialer.Timeout = time.Second * common.DIAL_TIMEOUT
	return tls.DialWithDialer(&dialer, "tcp", nodeAddr, self.config)
}

type noTlsDialer struct{}

func (self *noTlsDialer) Dial(nodeAddr string) (net.Conn, error) {
	return net.DialTimeout("tcp", nodeAddr, time.Second*common.DIAL_TIMEOUT)
}

func NewListener(port uint16, config *config.P2PNodeConfig) (net.Listener, error) {
	if config == nil || !config.IsTLS {
		return net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	}

	// load cert
	cert, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
	if err != nil {
		return nil, err
	}
	// load root ca
	caData, err := ioutil.ReadFile(config.CAPath)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	ret := pool.AppendCertsFromPEM(caData)
	if !ret {
		return nil, errors.New("[p2p]failed to parse root certificate")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	}

	return tls.Listen("tcp", ":"+strconv.Itoa(int(port)), tlsConfig)
}
