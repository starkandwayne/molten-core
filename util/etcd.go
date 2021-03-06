package util

import (
	"fmt"
	"time"

	"go.etcd.io/etcd/client"
	"go.etcd.io/etcd/clientv3"
)

func NewEtcdV2KeysAPI() (client.KeysAPI, error) {
	c, err := newEtcdV2Client()
	if err != nil {
		return nil, err
	}
	return client.NewKeysAPI(c), nil
}

func NewEtcdV2MembersAPI() (client.MembersAPI, error) {
	c, err := newEtcdV2Client()
	if err != nil {
		return nil, err
	}
	return client.NewMembersAPI(c), nil
}

func newEtcdV2Client() (client.Client, error) {
	cfg := client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd v2 client: %s", err)

	}
	return c, nil
}

func NewEtcdV3Client() (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd v3 client: %s", err)
	}

	return cli, nil
}

func NewAPIClient(c client.Client) (client.KeysAPI, error) {

	return client.NewKeysAPI(c), nil
}
