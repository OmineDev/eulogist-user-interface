package main

import (
	"encoding/hex"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/df-mc/goleveldb/leveldb"
)

var keyHexString = "8afdfad6dec562be40a205e366ddc435"

func main() {
	var dbKeys [][]byte

	keyBytes, err := hex.DecodeString(keyHexString)
	if err != nil {
		panic(err)
	}
	keyString := string(keyBytes)

	db, err := leveldb.OpenFile("world/db", nil)
	if err != nil {
		panic(err)
	}

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		dbKeys = append(dbKeys, append([]byte{}, iter.Key()...))
	}
	iter.Release()

	for _, dbKey := range dbKeys {
		value, err := db.Get(dbKey, nil)
		if err != nil {
			panic(err)
		}

		value = crypto.
			FromBytes(value).
			SetKey(keyString).
			Aes().
			ECB().
			PKCS7Padding().
			Decrypt().
			ToBytes()

		err = db.Put(dbKey, value, nil)
		if err != nil {
			panic(err)
		}
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}
}
