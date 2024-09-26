package column

import (
	"prettyy-server-online/data/column"
	"prettyy-server-online/utils/tool"
)

type Client struct {
	manager *column.Manager
}

var defaultClient *Client

func NewClient() (*Client, error) {
	manager, err := column.NewManager()
	if err != nil {
		return nil, err
	}
	return &Client{manager: manager}, nil
}

func (c *Client) Add(cidSlice, titleSlice []string, uid int64) (err error) {
	if len(cidSlice) != len(titleSlice) {
		return tool.ErrParams
	}
	for index, cid := range cidSlice {
		col := &column.Column{
			Cid:          cid,
			Title:        titleSlice[index],
			CoverImg:     "https://s21.ax1x.com/2024/07/05/pkRgyT0.jpg",
			Summary:      "",
			FrontDisplay: "1",
			IsFreeColumn: "1",
			Uid:          uid,
		}
		err = c.manager.Add(col)
	}
	return
}

func Add(cidSlice, titleSlice []string, uid int64) (err error) {
	return defaultClient.Add(cidSlice, titleSlice, uid)
}

func init() {
	var err error
	defaultClient, err = NewClient()
	if err != nil {
		panic(err)
	}
}
