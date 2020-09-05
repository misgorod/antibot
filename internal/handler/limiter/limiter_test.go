package limiter

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var errMock = errors.New("mock error")

type storageMock struct {
	ExistsFunc       func(path string) (bool, error)
	CreateParentFunc func(path string) error
	CreateSeqFunc    func(path string, ttl time.Duration) error
	CreateFunc       func(path string, ttl time.Duration) error
	CountFunc        func(path string) (int, error)
}

func (s *storageMock) Exists(path string) (bool, error) {
	return s.ExistsFunc(path)
}

func (s *storageMock) CreateParent(path string) error {
	return s.CreateParentFunc(path)
}

func (s *storageMock) CreateSeq(path string, ttl time.Duration) error {
	return s.CreateSeqFunc(path, ttl)
}

func (s *storageMock) Create(path string, ttl time.Duration) error {
	return s.CreateFunc(path, ttl)
}

func (s *storageMock) Count(path string) (int, error) {
	return s.CountFunc(path)
}

func TestHandle(t *testing.T) {
	fixtures := map[string]struct {
		existsFunc       func(path string) (bool, error)
		createParentFunc func(path string) error
		createSeqFunc    func(path string, ttl time.Duration) error
		createFunc       func(path string, ttl time.Duration) error
		countFunc        func(path string) (int, error)
		ipHeader         string
		responseWriter   http.ResponseWriter
		request          *http.Request
		expectedCode     int
	}{
		"valid request": {
			ipHeader: "127.0.0.1",
			existsFunc: func(path string) (bool, error) {
				return false, nil
			},
			createParentFunc: func(path string) error {
				return nil
			},
			createSeqFunc: func(path string, ttl time.Duration) error {
				return nil
			},
			createFunc: func(path string, ttl time.Duration) error {
				return nil
			},
			countFunc: func(path string) (int, error) {
				return 0, nil
			},
			expectedCode: http.StatusOK,
		},
		"banned request": {
			ipHeader: "127.0.0.1",
			existsFunc: func(path string) (bool, error) {
				return true, nil
			},
			createParentFunc: func(path string) error {
				return nil
			},
			createSeqFunc: func(path string, ttl time.Duration) error {
				return nil
			},
			createFunc: func(path string, ttl time.Duration) error {
				return nil
			},
			countFunc: func(path string) (int, error) {
				return 0, nil
			},
			expectedCode: http.StatusForbidden,
		},
		"max request": {
			ipHeader: "127.0.0.1",
			existsFunc: func(path string) (bool, error) {
				return false, nil
			},
			createParentFunc: func(path string) error {
				return nil
			},
			createSeqFunc: func(path string, ttl time.Duration) error {
				return nil
			},
			createFunc: func(path string, ttl time.Duration) error {
				return nil
			},
			countFunc: func(path string) (int, error) {
				return 111, nil
			},
			expectedCode: http.StatusForbidden,
		},
		"invalid header": {
			ipHeader: "invalid",
			existsFunc: func(path string) (bool, error) {
				return false, nil
			},
			createParentFunc: func(path string) error {
				return nil
			},
			createSeqFunc: func(path string, ttl time.Duration) error {
				return nil
			},
			createFunc: func(path string, ttl time.Duration) error {
				return nil
			},
			countFunc: func(path string) (int, error) {
				return 0, nil
			},
			expectedCode: http.StatusOK,
		},
	}
	for name, f := range fixtures {
		t.Run(name, func(t *testing.T) {
			logger := logrus.New()
			logger.SetOutput(ioutil.Discard)
			limitHandler := &handler{
				&storageMock{
					ExistsFunc:       f.existsFunc,
					CreateParentFunc: f.createParentFunc,
					CreateSeqFunc:    f.createSeqFunc,
					CreateFunc:       f.createFunc,
					CountFunc:        f.countFunc,
				},
				logger,
				time.Minute,
				time.Minute,
				100,
				[]string{"test-zk01", "test-zk02"},
				"/prefix",
				"req",
				"ban",
			}
			handler := http.HandlerFunc(limitHandler.Handle)
			request, _ := http.NewRequest("GET", "/api", nil)
			request.Header.Add("X-Forwarded-For", f.ipHeader)
			recorder := httptest.NewRecorder()

			handler.ServeHTTP(recorder, request)
			assert.Equal(t, f.expectedCode, recorder.Code)
		})
	}
}
