package kvstorage

import (
	"context"
	"sync"
	"testing"
)

const (
	KEYNAME  string = "key1"
	KEYVALUE string = "key1 value"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		Name string
		Want *Storage
	}{
		{
			Name: "Normal run",
			Want: &Storage{
				Storage:     sync.Map{},
				Initialized: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var got interface{} = NewStorage()
			switch got.(type) {
			case KVStorage:
			default:
				t.Errorf("NewStorage() = %v, want %v", got, test.Want)
			}
		})
	}
}

func TestKVStoragePut(t *testing.T) {
	storage := NewStorage()
	ctx := context.Background()

	type args struct {
		Key   string
		Value string
	}

	tests := []struct {
		Name    string
		Fields  KVStorage
		Args    args
		WantErr bool
	}{
		{
			Name:    "Storage is not initialized",
			Fields:  &Storage{Initialized: false},
			Args:    args{Key: KEYNAME, Value: KEYVALUE},
			WantErr: true,
		},
		{
			Name:    "Empty key",
			Fields:  storage,
			Args:    args{Key: "", Value: KEYVALUE},
			WantErr: true,
		},
		{
			Name:    "Good key, empty value",
			Fields:  storage,
			Args:    args{Key: KEYNAME, Value: ""},
			WantErr: false,
		},
		{
			Name:    "Good key, good value",
			Fields:  storage,
			Args:    args{Key: "key2", Value: "value for key2"},
			WantErr: false,
		},
		{
			Name:    "Key already in storage (update key-value pair)",
			Fields:  storage,
			Args:    args{Key: "key2", Value: "new value for key2"},
			WantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			kv := &test.Fields
			if err := (*kv).Put(ctx, test.Args.Key, test.Args.Value); (err != nil) != test.WantErr {
				t.Errorf("KVStorage.Set() error = %v, wantErr %v", err, test.WantErr)
			}
		})
	}
}

func TestKVStorageGet(t *testing.T) {
	goodStorage := NewStorage()
	badStorage := NewStorage()
	badStorage.Initialized = false
	ctx := context.Background()
	keyForEmptyValue := "empty-key"

	type args struct {
		Key string
	}

	tests := []struct {
		Name    string
		Fields  KVStorage
		Args    args
		Want    string
		WantErr bool
	}{
		{
			Name:    "Storage is not initialized",
			Fields:  badStorage,
			Args:    args{KEYNAME},
			WantErr: true,
		},
		{
			Name:    "Empty key",
			Fields:  goodStorage,
			Args:    args{""},
			WantErr: true,
		},
		{
			Name:    "Key is in the storage, value is not empty",
			Fields:  goodStorage,
			Args:    args{KEYNAME},
			Want:    KEYVALUE,
			WantErr: false,
		},
		{
			Name:    "Key is in the storage, value is empty",
			Fields:  goodStorage,
			Args:    args{keyForEmptyValue},
			Want:    "",
			WantErr: false,
		},
		{
			Name:    "Key is not in the storage",
			Fields:  goodStorage,
			Args:    args{KEYVALUE + "xxx"},
			WantErr: false,
		},
	}

	goodStorage.Put(ctx, KEYNAME, KEYVALUE)
	goodStorage.Put(ctx, keyForEmptyValue, "")

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			kv := &test.Fields
			got, err := (*kv).Get(ctx, test.Args.Key)
			if (err != nil) != test.WantErr {
				t.Errorf("KVStorage.Get() error = %v, wantErr %v", err, test.WantErr)
				return
			}
			if got != nil && got != test.Want {
				t.Errorf("KVStorage.Get() = %v, want %v", got, test.Want)
			}
		})
	}
}

func TestKVStorageDelete(t *testing.T) {
	goodStorage := NewStorage()
	emptyStorage := NewStorage()
	badStorage := NewStorage()
	badStorage.Initialized = false
	ctx := context.Background()
	tests := []struct {
		Name    string
		Fields  KVStorage
		Key     string
		WantErr bool
	}{
		{
			Name:    "Storage is not initialized",
			Fields:  badStorage,
			Key:     KEYNAME,
			WantErr: true,
		},
		{
			Name:    "Empty key, normal storage",
			Fields:  goodStorage,
			Key:     "",
			WantErr: true,
		},
		{
			Name:    "Good key, empty storage",
			Fields:  emptyStorage,
			Key:     KEYNAME,
			WantErr: false,
		},
		{
			Name:    "Key is not found in storage (storage is not empty)",
			Fields:  goodStorage,
			Key:     KEYNAME + "xxx",
			WantErr: false,
		},
		{
			Name:    "Key is found in storage",
			Fields:  goodStorage,
			Key:     KEYNAME,
			WantErr: false,
		},
	}

	goodStorage.Put(ctx, KEYNAME, KEYVALUE)

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			kv := test.Fields
			err := kv.Delete(ctx, test.Key)
			got, _ := kv.Get(ctx, test.Key)
			if (err != nil) != test.WantErr {
				t.Errorf("KVStorage.Delete() error = %v, wantErr %v", err, test.WantErr)
				return
			}
			if err != nil && got != nil {
				t.Errorf("KVStorage.Get() = %v", got)
				return
			}
		})
	}
}
