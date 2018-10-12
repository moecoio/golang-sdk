package prot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Client struct {
	client        http.Client
	masterNodeUrl string
	hash          string
	apiKey        string
	log           *logrus.Logger
}

func NewClient(url, apiKey, hash string) Client {
	return Client{
		client:        http.Client{},
		masterNodeUrl: url,
		hash:          hash,
		apiKey:        apiKey,
		log:           nil,
	}
}

func (c *Client) Init(log *logrus.Logger) error {
	c.log = log
	path := "/api/gate/auth"
	bodyReq, err := json.Marshal(InitGateReq{c.apiKey, Gate{Hash: c.hash}})
	if err != nil {
		return err
	}

	_, err = c.sendRequest("POST", path, bodyReq)
	return err
}

func (c *Client) sendRequest(method, path string, body []byte) ([]byte, error) {
	c.log.Debugf("gate sendRequest path: %s reqBody: %s", path, body);
	req, err := http.NewRequest(method, c.masterNodeUrl+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Gateway "+c.hash)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	c.log.Debugf("gate sendRequest path: %s response: %s", path, body);

	return body, nil
}

func (c *Client) SyncTransaction(transactions Transactions, lastSync int) (*SyncResponse, error) {
	path := "/api/gate/sync"
	if lastSync != 0 {
		path += "?last_sync=" + string(lastSync)
	}

	reqBody, err := json.Marshal(transactions)
	if err != nil {
		return nil, err
	}

	body, err := c.sendRequest("POST", path, reqBody)
	if err != nil {
		return nil, err
	}

	var res SyncResponse
	err = json.Unmarshal(body, &res)
	return &res, err
}

func (c *Client) GetDevices() (*DeviceResponse, error) {
	path := "/api/gate/v2/devices"

	body, err := c.sendRequest("GET", path, []byte{})
	if err != nil {
		return nil, err
	}

	var res DeviceResponse
	err = json.Unmarshal(body, &res)
	return &res, err
}
