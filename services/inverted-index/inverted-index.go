package inverted_index

import (
	"errors"
	invertedIndex "prettyy-server-online/data/inverted-index"
	"prettyy-server-online/utils/tool"
)

type Client struct {
	manager *invertedIndex.Manager
}

var defaultClient *Client

func NewClient() (*Client, error) {
	manager, err := invertedIndex.NewManager()
	if err != nil {
		return nil, err
	}
	return &Client{manager: manager}, nil
}
func Get(typ, attrValue string) ([]*invertedIndex.InvertedIndex, error) {
	return defaultClient.Get(typ, attrValue, "")
}

func (c *Client) IsExist(typ, attrValue, index string) bool {
	ss, _ := c.Get(typ, attrValue, index)
	return len(ss) > 0
}

func IsExist(typ, attrValue, index string) bool {
	return defaultClient.IsExist(typ, attrValue, index)
}

func (c *Client) Get(typ, attrValue, index string) ([]*invertedIndex.InvertedIndex, error) {
	if typ == "" || attrValue == "" {
		return nil, tool.ErrParams
	}
	// key 由 typ和attr_value拼接而成
	i, err := c.manager.Get(typ, attrValue, index)
	if err != nil {
		return nil, errors.New("get inverted index from mysql failed: " + err.Error())
	}
	return i, nil
}

func Add(i *invertedIndex.InvertedIndex) (err error) {
	return defaultClient.Add(i)
}

func (c *Client) Add(i *invertedIndex.InvertedIndex) (err error) {
	if i == nil {
		return tool.ErrParams
	}
	if err = c.manager.Add(i); err != nil {
		return errors.New("add inverted index to mysql failed: " + err.Error())
	}
	return
}

func init() {
	var err error
	defaultClient, err = NewClient()
	if err != nil {
		panic(err)
	}
}
