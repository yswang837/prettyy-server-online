package inverted_index

import (
	"errors"
	xzfSnowflake "prettyy-server-online/custom-pkg/xzf-snowflake"
	invertedIndex "prettyy-server-online/data/inverted-index"
	"prettyy-server-online/utils/tool"
	"strconv"
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
func Get(attrValue, number string) (*invertedIndex.InvertedIndex, error) {
	return defaultClient.Get(attrValue, number)
}

func (c *Client) Get(attrValue, number string) (*invertedIndex.InvertedIndex, error) {
	if attrValue == "" || number == "" {
		return nil, tool.ErrParams
	}
	// key 由 attr_value和number拼接而成
	iMap, err := c.cacheManager.HGetAll(attrValue + number)
	if err != nil {
		return nil, errors.New("get inverted index from redis failed: " + err.Error())
	}
	if len(iMap) != 0 {
		return mapToInvertedIndex(iMap), nil
	}
	i, err := c.manager.Get(attrValue, number)
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
	i.CreateTime = time.Now()
	if _, err = c.cacheManager.HMSet(i.AttrValue+i.Number, invertedIndexToMap(i)); err != nil {
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
	m["number"] = i.Number
	m["uid"] = i.Uid
	m["create_time"] = i.CreateTime.Format(tool.DefaultDateTimeLayout)
	return m
}

func mapToInvertedIndex(m map[string]string) *invertedIndex.InvertedIndex {
	if len(m) == 0 {
		return nil
	}
	uid, _ := strconv.Atoi(m["uid"])
	a := &invertedIndex.InvertedIndex{}
	a.AttrValue = m["attr_value"]
	a.Number = m["number"]
	a.Uid = int64(uid)
	a.CreateTime = tool.StringToTime(m["create_time"])
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
