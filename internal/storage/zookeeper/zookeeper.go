package zookeeper

import (
	"strings"
	"time"

	"github.com/go-zookeeper/zk"
)

const sessionTimeout = 10 * time.Second

var acl = zk.WorldACL(zk.PermAll)
var payload = []byte("")

type zookeeper interface {
	Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
	CreateTTL(path string, data []byte, flags int32, acl []zk.ACL, ttl time.Duration) (string, error)
	Exists(path string) (bool, *zk.Stat, error)
	Children(path string) ([]string, *zk.Stat, error)
}

type storage struct {
	zookeeper
}

func New(hosts []string) *storage {
	connection, _, err := zk.Connect(hosts, sessionTimeout)
	if err != nil {
		panic(err)
	}
	return &storage{
		connection,
	}
}

func (s *storage) CreateParent(path string) error {
	if path[0] != '/' {
		return zk.ErrInvalidPath
	}
	pathString := ""
	pathNodes := strings.Split(path, "/")
	for i := 1; i < len(pathNodes); i++ {
		pathString += "/" + pathNodes[i]
		_, err := s.zookeeper.Create(pathString, payload, 0, acl)
		if err != nil && err != zk.ErrNodeExists && err != zk.ErrNoAuth {
			return err
		}
	}
	return nil
}

func (s *storage) Exists(path string) (bool, error) {
	exists, _, err := s.zookeeper.Exists(path)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *storage) CreateSeq(path string, ttl time.Duration) error {
	_, err := s.CreateTTL(path, payload, zk.FlagTTL|zk.FlagSequence, acl, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) Create(path string, ttl time.Duration) error {
	_, err := s.CreateTTL(path, payload, zk.FlagTTL|zk.FlagEphemeral, acl, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) Count(path string) (int, error) {
	children, _, err := s.Children(path)
	if err != nil {
		return 0, err
	}
	return len(children), nil
}
