package inverted_index

import (
	"errors"
	xzfSnowflake "prettyy-server-online/custom-pkg/xzf-snowflake"
	invertedIndex "prettyy-server-online/data/inverted-index"
	"prettyy-server-online/utils/tool"
	"time"
)

type Client struct {
	manager      *invertedIndex.Manager
	cacheManager *invertedIndex.ManagerRedis
}

var defaultClient *Client

func NewClient() (*Client, error) {
	manager, err := invertedIndex.NewManager()
	if err != nil {
		return nil, err
	}
	cache, err := invertedIndex.NewManagerRedis()
	if err != nil {
		return nil, err
	}
	return &Client{manager: manager, cacheManager: cache}, nil
}
func Get(typ, attrValue string) (*invertedIndex.InvertedIndex, error) {
	return defaultClient.Get(typ, attrValue)
}

func (c *Client) Get(typ, attrValue string) (*invertedIndex.InvertedIndex, error) {
	if typ == "" || attrValue == "" {
		return nil, tool.ErrParams
	}
	// key 由 typ和attr_value拼接而成
	iMap, err := c.cacheManager.HGetAll(typ + attrValue)
	if err != nil {
		return nil, errors.New("get inverted index from redis failed: " + err.Error())
	}
	if len(iMap) != 0 {
		return mapToInvertedIndex(iMap), nil
	}
	i, err := c.manager.Get(typ, attrValue)
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
	now := time.Now()
	i.CreateTime = now
	i.UpdateTime = now
	if _, err = c.cacheManager.HMSet(i.Typ+i.AttrValue, invertedIndexToMap(i)); err != nil {
		return errors.New("set inverted index to redis failed: " + err.Error())
	}
	return
}

func invertedIndexToMap(i *invertedIndex.InvertedIndex) map[string]interface{} {
	if i == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["attr_value"] = i.AttrValue
	m["typ"] = i.Typ
	m["index"] = i.Index
	m["create_time"] = i.CreateTime.Format(tool.DefaultDateTimeLayout)
	m["update_time"] = i.CreateTime.Format(tool.DefaultDateTimeLayout)
	return m
}

func mapToInvertedIndex(m map[string]string) *invertedIndex.InvertedIndex {
	if len(m) == 0 {
		return nil
	}
	a := &invertedIndex.InvertedIndex{}
	a.AttrValue = m["attr_value"]
	a.Typ = m["typ"]
	a.Index = m["index"]
	a.CreateTime = tool.StringToTime(m["create_time"])
	a.UpdateTime = tool.StringToTime(m["update_time"])
	return a
}

func init() {
	var err error
	defaultClient, err = NewClient()
	if err != nil {
		panic(err)
	}
	if err = xzfSnowflake.Init("2024-03-09", "1"); err != nil {
		panic(err)
	}
}
