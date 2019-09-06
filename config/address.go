package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
)

// urlPath is the restful path to get nodes from broker.
const urlPath = "/broker/node/state"

// errNoAvailableNode is the error return then no sending node is available
var errNoAvailableNode = errors.New("no available node")

// AddressManager defines the way to get broker node.
type AddressManager interface {
	// BrokerAddressList return a list of available broker nodes.
	BrokerAddressList() []*models.Node
	// RandomNext randomly pick a node from available broker node list.
	// Return errNoAvailableNode when no node is available.
	RandomNext() (*models.Node, error)
}

// addressManager implements AddressManager
type addressManager struct {
	clientConfig *ClientConfig
	// url to fetch broker node
	url        string
	httpClient http.Client
	// latest updated nodes
	nodes      []*models.Node
	rand       *rand.Rand
	lock4nodes sync.RWMutex
	logger     *logger.Logger
}

// BrokerAddressList return a list of available broker nodes.
func (m *addressManager) BrokerAddressList() []*models.Node {
	m.lock4nodes.RLock()
	defer m.lock4nodes.RUnlock()
	list := make([]*models.Node, 0, len(m.nodes))
	list = append(list, m.nodes...)
	return list
}

// RandomNext randomly pick a node from available broker node list.
// Return errNoAvailableNode when no node is available.
func (m *addressManager) RandomNext() (*models.Node, error) {
	m.lock4nodes.RLock()
	defer m.lock4nodes.RUnlock()
	nodesNum := len(m.nodes)
	if nodesNum == 0 {
		return nil, errNoAvailableNode
	}
	index := m.rand.Intn(nodesNum)
	return m.nodes[index], nil
}

// syncBrokerAddressList synchronizes the available broker node list.
func (m *addressManager) syncBrokerAddressList() {
	resp, err := m.httpClient.Get(m.url)
	if err != nil {
		m.logger.Error("syncBrokerAddressList http get", logger.Error(err))
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			m.logger.Error("syncBrokerAddressList close response body", logger.Error(err))
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m.logger.Error("syncBrokerAddressList read resp body", logger.Error(err))
		return
	}

	var activeNodes []models.ActiveNode
	if err := json.Unmarshal(body, &activeNodes); err != nil {
		m.logger.Error("syncBrokerAddressList Unmarshal resp body", logger.Error(err))
		return
	}

	if len(activeNodes) == 0 {
		m.logger.Error("syncBrokerAddressList address list is empty")
		return
	}

	list := make([]*models.Node, len(activeNodes))
	for index, node := range activeNodes {
		list[index] = &node.Node
	}

	m.lock4nodes.Lock()
	defer m.lock4nodes.Unlock()
	m.nodes = list
}

func (m *addressManager) initSyncTask() {
	ticker := time.NewTicker(time.Duration(m.clientConfig.SyncAddressInterval) * time.Second)

	go func() {
		for range ticker.C {
			m.syncBrokerAddressList()
		}
	}()
}

// NewAddressManager creates a AddressManager with clientConfig.
func NewAddressManager(clientConfig *ClientConfig) AddressManager {
	m := &addressManager{
		clientConfig: clientConfig,
		url:          clientConfig.BrokerURL + urlPath,
		rand:         rand.New(rand.NewSource(time.Now().Unix())),
		httpClient: http.Client{
			Timeout: time.Duration(clientConfig.SyncAddressTimeout) * time.Second,
		},
		logger: logger.GetLogger("config", "AddressManager"),
	}

	m.syncBrokerAddressList()
	m.initSyncTask()
	return m
}

func BuildAddress(node *models.Node) string {
	return fmt.Sprintf("%s:%d", node.IP, node.TCPPort)
}
