package main

import (
	"bytes"
)

type bucket[T any] struct {
	key []byte
	val T
}

type hashMap[T any] []bucket[T]

func (hm *hashMap[T]) contains(key []byte) bool {
	length := len(*hm)
	hash := hash(key)
	pos := int(hash % uint64(length))
	if bytes.Compare((*hm)[pos].key, key) == 0 {
		return true
	} else {
		for i := pos + 1; i < length; i++ {
			if bytes.Compare((*hm)[i].key, key) == 0 {
				return true
			}
		}
		for i := 0; i < pos; i++ {
			if bytes.Compare((*hm)[i].key, key) == 0 {
				return true
			}
		}
	}
	return false
}

func (hm *hashMap[T]) find(key []byte) (T, uint64) {
	length := len(*hm)
	hash := hash(key)
	pos := int(hash % uint64(length))
	if bytes.Compare((*hm)[pos].key, key) == 0 {
		return (*hm)[pos].val, hash
	} else {
		for i := pos + 1; i < len(*hm); i++ {
			if (*hm)[i].key == nil {
				return (*hm)[i].val, hash
			}
		}
		for i := 0; i < pos; i++ {
			if (*hm)[i].key == nil {
				return (*hm)[i].val, hash
			}
		}
	}
	return *new(T), hash
}

func (hm *hashMap[T]) insertByHash(hash uint64, key []byte, val T) {
	length := len(*hm)
	pos := int(hash % uint64(length))
	(*hm)[pos].key = key
	(*hm)[pos].val = val
}

func (hm *hashMap[T]) insert(key []byte, val T) {
	hash := hash(key)
	hm.insertByHash(hash, key, val)
}

// returns 1 if collison occured, 0 otherwise
func (hm *hashMap[T]) insertChecked(key []byte, val T) int {
	length := len(*hm)
	hash := hash(key)
	pos := int(hash % uint64(length))
	if (*hm)[pos].key == nil {
		(*hm)[pos].key = key
		(*hm)[pos].val = val
	} else {
		for i := pos + 1; i < len(*hm); i++ {
			if (*hm)[i].key == nil {
				(*hm)[i].key = key
				(*hm)[i].val = val
				return 1
			}
		}
		for i := 0; i < pos; i++ {
			if (*hm)[i].key == nil {
				(*hm)[i].key = key
				(*hm)[i].val = val
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
