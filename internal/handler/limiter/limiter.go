package limiter

import (
	"github.com/misgorod/antibot/internal/config"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"

	"github.com/misgorod/antibot/internal/common"

	"github.com/pkg/errors"
)

var mask = net.IPv4Mask(255, 255, 0, 0)

type storage interface {
	Exists(path string) (bool, error)
	CreateParent(path string) error
	CreateSeq(path string, ttl time.Duration) error
	Create(path string, ttl time.Duration) error
	Count(path string) (int, error)
}

type handler struct {
	storage
	log *logrus.Logger
	banTime time.Duration
	requestTime time.Duration
	requestLimit int
	zkHosts []string
	zkPrefix string
	zkRequestNode string
	zkBanNode string
}

func New(storage storage, logger *logrus.Logger, config config.Config) *handler {
	return &handler{
		storage: storage,
		log: logger,
		banTime: config.BanTime,
		requestTime: config.RequestTime,
		requestLimit: config.RequestLimit,
		zkHosts: config.ZkHosts,
		zkPrefix: config.ZkPrefix,
		zkRequestNode: config.ZkRequestNode,
		zkBanNode: config.ZkBanNode,
	}
}

func parseIP(ipHeader string) (net.IP, error) {
	ip := net.ParseIP(ipHeader)
	if ip == nil {
		return nil, errors.New("invalid ip")
	}
	ipMasked := ip.Mask(mask)
	return ipMasked, nil
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	ipHeader := r.Header.Get("X-Forwarded-For")

	ipMasked, err := parseIP(ipHeader)
	if err != nil {
		common.RespondOK(w)
		return
	}

	log := h.log.WithField("ip", ipHeader).Logger
	log = h.log.WithField("mask", ipMasked).Logger

	containerPath := h.zkPrefix + "/" + ipMasked.String()
	queuePath := containerPath + "/queue"
	nodePath := queuePath + "/" + h.zkRequestNode
	banNodePath := containerPath + "/" + h.zkBanNode

	log.Tracef("container path: %v", containerPath)
	log.Tracef("node path: %v", nodePath)
	log.Tracef("ban path: %v", banNodePath)

	err = h.storage.CreateParent(queuePath)
	if err != nil {
		common.RespondInternalServerError(log, w, errors.Wrap(err, "failed to create parent node"))
		return
	}

	// Check if subnet is banned
	exists, err := h.storage.Exists(banNodePath)
	if err != nil {
		common.RespondInternalServerError(log, w, errors.Wrap(err, "failed to check existence of ban node"))
		return
	}
	if exists {
		log.Debug("ban node exists")
		common.RespondBadRequest(w)
		return
	}

	// Create sequential node
	err = h.storage.CreateSeq(nodePath, h.requestTime)
	if err != nil {
		common.RespondInternalServerError(log, w, errors.Wrap(err, "failed to create ttl sequential node"))
		return
	}

	// Check nodes count
	count, err := h.storage.Count(queuePath)
	if err != nil {
		common.RespondInternalServerError(log, w, errors.Wrap(err, "failed to get children of container path"))
		return
	}
	if count > h.requestLimit {
		err = h.storage.Create(banNodePath, h.banTime)
		if err != nil {
			common.RespondInternalServerError(log, w, errors.Wrap(err, "failed to create ban node"))
			return
		}
		common.RespondBadRequest(w)
		return
	}

	common.RespondOK(w)
}
