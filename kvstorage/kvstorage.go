package kvstorage

import (
	"context"
	"errors"
	"sync"
)

type KVStorage interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Put(ctx context.Context, key string, val interface{}) error
	Delete(ctx context.Context, key string) error
}

type Storage struct {
	Storage     sync.Map
	Initialized bool
}

const ERROR_MESSAGE string = "Storage is not initialized or key is empty"

func NewStorage() *Storage {
	return &Storage{
		Storage:     sync.Map{},
		Initialized: true,
	}
}

func (storage *Storage) Get(ctx context.Context, key string) (interface{}, error) {
	if !storage.Initialized || len(key) == 0 {
		return nil, errors.New(ERROR_MESSAGE)
	}
	val, ok := storage.Storage.Load(key)
	if ok {
		return val, nil
	}
	return nil, nil
}

func (storage *Storage) Put(ctx context.Context, key string, val interface{}) error {
	if !storage.Initialized || len(key) == 0 {
		return errors.New(ERROR_MESSAGE)
	}
	storage.Storage.Store(key, val)
	return nil
}

func (storage *Storage) Delete(ctx context.Context, key string) error {
	if !storage.Initialized || len(key) == 0 {
		return errors.New(ERROR_MESSAGE)
	}
	storage.Storage.Delete(key)
	return nil
}
