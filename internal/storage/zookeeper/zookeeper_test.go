package zookeeper

import (
	"errors"
	"testing"
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
)

var mockError = errors.New("mock error")

type zkMock struct {
	CreateFunc    func(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
	CreateTTLFunc func(path string, data []byte, flags int32, acl []zk.ACL, ttl time.Duration) (string, error)
	ExistsFunc    func(path string) (bool, *zk.Stat, error)
	ChildrenFunc  func(path string) ([]string, *zk.Stat, error)
}

func (z *zkMock) Create(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	return z.CreateFunc(path, data, flags, acl)
}

func (z *zkMock) CreateTTL(path string, data []byte, flags int32, acl []zk.ACL, ttl time.Duration) (string, error) {
	return z.CreateTTLFunc(path, data, flags, acl, ttl)
}

func (z *zkMock) Exists(path string) (bool, *zk.Stat, error) {
	return z.ExistsFunc(path)
}

func (z *zkMock) Children(path string) ([]string, *zk.Stat, error) {
	return z.ChildrenFunc(path)
}

func TestCreateParent(t *testing.T) {
	fixtures := map[string]struct {
		path       string
		createFunc func(path string, data []byte, flags int32, acl []zk.ACL) (string, error)
		expected   error
	}{
		"ok": {
			"/path/to/node",
			func(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
				return "", nil
			},
			nil,
		},
		"invalid path": {
			"invalid/path",
			nil,
			zk.ErrInvalidPath,
		},
		"node exists": {
			"/path",
			func(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
				return "", zk.ErrNodeExists
			},
			nil,
		},
		"internal error": {
			"/path",
			func(path string, data []byte, flags int32, acl []zk.ACL) (string, error) {
				return "", mockError
			},
			mockError,
		},
	}
	for name, f := range fixtures {
		t.Run(name, func(t *testing.T) {
			storage := storage{&zkMock{CreateFunc: f.createFunc}}
			actual := storage.CreateParent(f.path)
			assert.Equal(t, f.expected, actual)
		})
	}
}

func TestExists(t *testing.T) {
	fixtures := map[string]struct {
		path           string
		function       func(path string) (bool, *zk.Stat, error)
		expectedExists bool
		expectedError  error
	}{
		"path exists": {
			"/path/to/node",
			func(path string) (bool, *zk.Stat, error) {
				return true, nil, nil
			},
			true,
			nil,
		},
		"path not exists": {
			"/path/to/node",
			func(path string) (bool, *zk.Stat, error) {
				return false, nil, nil
			},
			false,
			nil,
		},
		"internal error": {
			"/path",
			func(path string) (bool, *zk.Stat, error) {
				return false, nil, mockError
			},
			false,
			mockError,
		},
	}
	for name, f := range fixtures {
		t.Run(name, func(t *testing.T) {
			storage := storage{&zkMock{ExistsFunc: f.function}}
			actualExists, actualError := storage.Exists(f.path)
			assert.Equal(t, f.expectedExists, actualExists)
			assert.Equal(t, f.expectedError, actualError)
		})
	}
}

func TestCreateSeq(t *testing.T) {
	fixtures := map[string]struct {
		path     string
		function func(path string, data []byte, flags int32, acl []zk.ACL, ttl time.Duration) (string, error)
		expected error
	}{
		"node created": {
			"/path/to/node",
			func(path string, data []byte, flags int32, acl []zk.ACL, ttl time.Duration) (string, error) {
				return "", nil
			},
			nil,
		},
		"path not exists": {
			"/path/to/node",
			func(path string, data []byte, flags int32, acl []zk.ACL, ttl time.Duration) (string, error) {
				return "", zk.ErrNoNode
			},
			zk.ErrNoNode,
		},
	}
	for name, f := range fixtures {
		t.Run(name, func(t *testing.T) {
			storage := storage{&zkMock{CreateTTLFunc: f.function}}
			actual := storage.CreateSeq(f.path, time.Second)
			assert.Equal(t, f.expected, actual)
		})
	}
}

func TestCreate(t *testing.T) {
	fixtures := map[string]struct {
		path     string
		function func(path string, data []byte, flags int32, acl []zk.ACL, ttl time.Duration) (string, error)
		expected error
	}{
		"node created": {
			"/path/to/node",
			func(path string, data []byte, flags int32, acl []zk.ACL, ttl time.Duration) (string, error) {
				return "", nil
			},
			nil,
		},
		"path not exists": {
			"/path/to/node",
			func(path string, data []byte, flags int32, acl []zk.ACL, ttl time.Duration) (string, error) {
				return "", zk.ErrNoNode
			},
			zk.ErrNoNode,
		},
	}
	for name, f := range fixtures {
		t.Run(name, func(t *testing.T) {
			storage := storage{&zkMock{CreateTTLFunc: f.function}}
			actual := storage.Create(f.path, time.Second)
			assert.Equal(t, f.expected, actual)
		})
	}
}

func TestCount(t *testing.T) {
	fixtures := map[string]struct {
		path          string
		function      func(path string) ([]string, *zk.Stat, error)
		expectedCount int
		expectedError error
	}{
		"got count of children": {
			"/path/to/node",
			func(path string) ([]string, *zk.Stat, error) {
				return []string{"a", "b", "c"}, nil, nil
			},
			3,
			nil,
		},
		"path not exists": {
			"/path/to/node",
			func(path string) ([]string, *zk.Stat, error) {
				return nil, nil, zk.ErrNoNode
			},
			0,
			zk.ErrNoNode,
		},
	}
	for name, f := range fixtures {
		t.Run(name, func(t *testing.T) {
			storage := storage{&zkMock{ChildrenFunc: f.function}}
			actualCount, actualError := storage.Count(f.path)
			assert.Equal(t, f.expectedCount, actualCount)
			assert.Equal(t, f.expectedError, actualError)
		})
	}
}
