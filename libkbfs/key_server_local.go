package libkbfs

import (
	"fmt"

	"github.com/keybase/client/protocol/go"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

// KeyServerLocal just stores key server halves in a local leveldb
// instance.
//
// Per-block key server halves are keyed by block ID, and per-TLF key
// server halves are keyed by (dir id, key generation, device encryption
// public key) tuples.
type KeyServerLocal struct {
	codec Codec
	db    *leveldb.DB
}

var _ KeyOps = (*KeyServerLocal)(nil)

func newKeyServerLocalWithStorage(codec Codec, storage storage.Storage) (*KeyServerLocal, error) {
	db, err := leveldb.Open(
		storage,
		&opt.Options{
			Compression: opt.NoCompression,
		})
	if err != nil {
		return nil, err
	}
	kops := &KeyServerLocal{codec, db}
	return kops, nil
}

// NewKeyServerLocal returns a KeyServerLocal with a leveldb instance at the
// given file.
func NewKeyServerLocal(codec Codec, dbfile string) (*KeyServerLocal, error) {
	storage, err := storage.OpenFile(dbfile)
	if err != nil {
		return nil, err
	}
	return newKeyServerLocalWithStorage(codec, storage)
}

// NewKeyServerMemory returns a KeyServerLocal with an in-memory leveldb
// instance.
func NewKeyServerMemory(codec Codec) (*KeyServerLocal, error) {
	return newKeyServerLocalWithStorage(codec, storage.NewMemStorage())
}

// GetBlockCryptKeyServerHalf implements the KeyOps interface for
// KeyServerLocal.
func (ks *KeyServerLocal) GetBlockCryptKeyServerHalf(id BlockID) (serverHalf BlockCryptKeyServerHalf, err error) {
	data, err := ks.db.Get(id[:], nil)
	if err != nil {
		return
	}
	if len(data) != len(serverHalf.ServerHalf) {
		err = fmt.Errorf("Expected length %d, got %d", len(serverHalf.ServerHalf), len(data))
		return
	}
	copy(serverHalf.ServerHalf[:], data)
	return
}

// PutBlockCryptKeyServerHalf implements the KeyOps interface for
// KeyServerLocal.
func (ks *KeyServerLocal) PutBlockCryptKeyServerHalf(id BlockID, serverHalf BlockCryptKeyServerHalf) error {
	return ks.db.Put(id[:], serverHalf.ServerHalf[:], nil)
}

// DeleteBlockCryptKeyServerHalf implements the KeyOps interface for
// KeyServerLocal.
func (ks *KeyServerLocal) DeleteBlockCryptKeyServerHalf(id BlockID) error {
	return ks.db.Delete(id[:], nil)
}

type serverHalfID struct {
	ID             DirID
	KeyGen         KeyGen
	CryptPublicKey CryptPublicKey
}

// GetTLFCryptKeyServerHalf implements the KeyOps interface for
// KeyServerLocal.
func (ks *KeyServerLocal) GetTLFCryptKeyServerHalf(
	id DirID, keyGen KeyGen, cryptPublicKey CryptPublicKey) (serverHalf TLFCryptKeyServerHalf, err error) {
	idData, err := ks.codec.Encode(serverHalfID{id, keyGen, cryptPublicKey})
	if err != nil {
		return
	}

	data, err := ks.db.Get(idData, nil)
	if err != nil {
		return
	}
	if len(data) != len(serverHalf.ServerHalf) {
		err = fmt.Errorf("Expected length %d, got %d", len(serverHalf.ServerHalf), len(data))
		return
	}
	copy(serverHalf.ServerHalf[:], data)
	return
}

// PutTLFCryptKeyServerHalf implements the KeyOps interface for
// KeyServerLocal.
func (ks *KeyServerLocal) PutTLFCryptKeyServerHalf(
	id DirID, keyGen KeyGen, cryptPublicKey CryptPublicKey, serverHalf TLFCryptKeyServerHalf) error {
	idData, err := ks.codec.Encode(serverHalfID{id, keyGen, cryptPublicKey})
	if err != nil {
		return err
	}

	return ks.db.Put(idData, serverHalf.ServerHalf[:], nil)
}

// GetMacPublicKey implements the KeyOps interface for KeyServerLocal.
func (ks *KeyServerLocal) GetMacPublicKey(uid keybase1.UID) (MacPublicKey, error) {
	return MacPublicKey{}, nil
}
