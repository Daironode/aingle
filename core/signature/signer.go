
package signature

import (
	"github.com/Daironode/aingle-crypto/keypair"
	"github.com/Daironode/aingle-crypto/signature"
)

// Signer is the abstract interface of user's information(Keys) for signing data.
type Signer interface {
	//get signer's private key
	PrivKey() keypair.PrivateKey

	//get signer's public key
	PubKey() keypair.PublicKey

	Scheme() signature.SignatureScheme
}
