package config

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"
	"time"

	"github.com/starkandwayne/molten-core/certs"
	"github.com/starkandwayne/molten-core/flannel"
	"github.com/starkandwayne/molten-core/util"

	"go.etcd.io/etcd/client"
)

const (
	etcdMolenCorePath  string = "/moltencore/nodes"
	dockerCertValidFor        = time.Hour * 24 * 365
	dockerTLSPort             = 2376
	singletonZoneIndex        = uint16(0)
)

type Docker struct {
	Endpoint string
	CA       certs.Cert
	Server   certs.Cert
	Client   certs.Cert
}

type NodeConfig struct {
	Subnet    flannel.Subnet
	ZoneIndex uint16
	Docker    Docker
	PrivateIP net.IP
	PublicIP  net.IP
}

func (nc NodeConfig) IsSingletonZone() bool {
	return nc.ZoneIndex == singletonZoneIndex
}

func (nc NodeConfig) Zone() string {
	return fmt.Sprintf("z%d", nc.ZoneIndex)
}

func (nc NodeConfig) CPIName() string {
	return fmt.Sprintf("docker-%s", nc.Zone())
}

func LoadNodeConfigs() (*[]NodeConfig, error) {
	kapi, err := util.NewEtcdV2KeysAPI()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	resp, err := kapi.Get(ctx, etcdMolenCorePath, &client.GetOptions{Recursive: true})
	if err != nil {
		return nil, fmt.Errorf("failed to load node configs from etcd: %s", err)
	}

	var confs []NodeConfig
	for _, node := range resp.Node.Nodes {
		var c NodeConfig
		err = json.Unmarshal([]byte(node.Value), &c)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal node config: %s", err)
		}
		confs = append(confs, c)
	}
	return &confs, nil
}

func LoadNodeConfig() (*NodeConfig, error) {
	privateIP, err := util.LookupIpV4Address(false)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup private node ip: %s", err)
	}

	kapi, err := util.NewEtcdV2KeysAPI()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	resp, err := kapi.Get(ctx, nodePath(privateIP), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load node config from etcd: %s", err)
	}

	var c NodeConfig
	err = json.Unmarshal([]byte(resp.Node.Value), &c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal node config: %s", err)
	}
	return &c, nil
}

func GenereateNodeConfig(index uint16) (*NodeConfig, error) {
	privateIP, err := util.LookupIpV4Address(false)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup private node ip: %s", err)
	}

	publicIP, err := util.LookupIpV4Address(true)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup public node ip: %s", err)
	}

	subnet, err := flannel.GetSubnetByIndex(index)
	if err != nil {
		return nil, fmt.Errorf("failed to generate docker certs: %s", err)
	}

	docker, err := newDocker(subnet, privateIP)
	if err != nil {
		return nil, fmt.Errorf("failed to generate docker certs: %s", err)
	}

	conf := NodeConfig{Subnet: subnet, Docker: docker, ZoneIndex: index,
		PrivateIP: privateIP, PublicIP: publicIP}

	err = conf.save()
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func (nc NodeConfig) save() error {
	kapi, err := util.NewEtcdV2KeysAPI()
	if err != nil {
		return err
	}

	ctx := context.Background()

	rawConf, err := json.Marshal(nc)
	if err != nil {
		return fmt.Errorf("failed to marshal node config: %s", err)
	}
	_, err = kapi.Set(ctx, nodePath(nc.PrivateIP), string(rawConf), nil)
	if err != nil {
		return fmt.Errorf("failed to update node config in etcd: %s", err)
	}
	return nil
}

func nodePath(privateIP net.IP) string {
	return filepath.Join(etcdMolenCorePath, privateIP.String())
}

func newDocker(s flannel.Subnet, hostIP net.IP) (Docker, error) {
	caCert, err := certs.Genereate(certs.GenArg{
		ValidFor: dockerCertValidFor,
	})
	if err != nil {
		return Docker{}, fmt.Errorf("failed to generate docker ca cert: %s", err)
	}

	serverCert, err := certs.Genereate(certs.GenArg{
		CA:          caCert,
		ValidFor:    dockerCertValidFor,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{hostIP, net.ParseIP("127.0.0.1")},
	})
	if err != nil {
		return Docker{}, fmt.Errorf("failed to generate docker server cert: %s", err)
	}

	clientCert, err := certs.Genereate(certs.GenArg{
		CA:          caCert,
		ValidFor:    dockerCertValidFor,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		return Docker{}, fmt.Errorf("failed to generate docker client cert: %s", err)
	}

	return Docker{
		Endpoint: fmt.Sprintf("%s:%d", hostIP, dockerTLSPort),
		CA:       caCert,
		Server:   serverCert,
		Client:   clientCert,
	}, nil
}
