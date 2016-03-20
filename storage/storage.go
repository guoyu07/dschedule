package storage

import (
	// "fmt"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/consul"
	// log "github.com/omidnikta/logrus"
)

func init() {
	// Register consul store to libkv
	consul.Register()

	// We can register as many backends that are supported by libkv
	// etcd.Register()
	// zookeeper.Register()
	// boltdb.Register()
}

type StorageConfig struct {
	Address     string
	BackendName string
	Config      *store.Config
}

type Storage struct {
	Client store.Store
}

func NewStorage(storageConfig *StorageConfig) (*Storage, error) {
	// Initialize a new store with consul
	kv, err := libkv.NewStore(
		store.Backend(storageConfig.BackendName), // or "consul","zookeeper","etcd"
		[]string{storageConfig.Address},
		// &store.Config{

		// ConnectionTimeout: 10 * time.Second,
		// },
		storageConfig.Config,
	)
	if err != nil {
		return nil, err
	}
	storage := &Storage{
		Client: kv,
	}
	return storage, nil
}

// func ( storage *Storage)
