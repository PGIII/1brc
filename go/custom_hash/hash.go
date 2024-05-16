package main

import (
	"bytes"
	"fmt"
)

type bucket[T any] struct {
	key []byte
	val *T
}

type hashMap[T any] []bucket[T]

func (hm *hashMap[T]) find(key []byte) (*T, uint64) {
	hash := hash(key)
	pos := hash % uint64(len(*hm))
	if bytes.Compare((*hm)[pos].key, key) == 0 {
		return (*hm)[pos].val, hash
	} else {
		for i := pos + 1; i != pos; i++ {
			if bytes.Compare((*hm)[pos].key, key) == 0 {
				return (*hm)[pos].val, hash
			}
		}
	}
	return nil, hash
}

func (hm *hashMap[T]) insertByHash(hash uint64, key []byte, val *T) {
	pos := hash % uint64(len(*hm))
	fmt.Printf("Len: %d", len(*hm))
	if (*hm)[pos].key == nil {
		(*hm)[pos].key = key
		(*hm)[pos].val = val
	} else {
		// exploiting integer overflow here
		for i := pos + 1; i != pos; i++ {
			if (*hm)[pos].key == nil {
				(*hm)[pos].key = key
				(*hm)[pos].val = val
				break
			}
		}
		panic("Couldn't find slot in hashmap")
	}
}

func (hm *hashMap[T]) insert(key []byte, val *T) {
	hash := hash(key)
	hm.insertByHash(hash, key, val)
}

// returns 1 if collison occured, 0 otherwise
func (hm *hashMap[T]) insertChecked(key []byte, val *T) uint {
	hash := hash(key)
	if (*hm)[hash].key == nil {
		(*hm)[hash].key = key
		(*hm)[hash].val = val
	} else {
		// exploiting integer overflow here
		for i := hash + 1; i != hash; i++ {
			if (*hm)[hash].key == nil {
				(*hm)[hash].key = key
				(*hm)[hash].val = val
				return 1
			}
		}
		panic("Couldn't find slot in hashmap")
	}
	return 0
}

func hash(data []byte) uint64 {
	const offsetBasis uint64 = 0xcbf29ce484222325
	const prime uint64 = 0x100000001b3
	hash := offsetBasis

	for _, b := range data {
		hash = hash ^ uint64(b)
		hash = hash * prime
	}
	return hash
}
