package internal

import (
	subv2 "github.com/threefoldtech/substrate-client-dev"
)

// Identity is the user identity to be used in substrate
type Identity subv2.Identity

// NewIdentityFromSr25519Phrase generates a new Sr25519 identity from mnemonics
func NewIdentityFromSr25519Phrase(phrase string) (Identity, error) {
	id, err := subv2.NewIdentityFromSr25519Phrase(phrase)
	return Identity(id), err
}
