package main

import (
	"encoding/hex"
	"fmt"

	"github.com/deatil/go-cryptobin/cryptobin/crypto"
	"github.com/df-mc/goleveldb/leveldb"
)

func main() {
	var dbKeys [][]byte

	keyHexString := ReadStringFromPanel("请输入存档解密密钥: ")
	mcworldPath := ReadStringFromPanel("请输入存档路径: ")

	keyBytes, err := hex.DecodeString(keyHexString)
	if err != nil {
		panic(err)
	}
	keyString := string(keyBytes)

	db, err := leveldb.OpenFile(fmt.Sprintf("%s/db", mcworldPath), nil)
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
