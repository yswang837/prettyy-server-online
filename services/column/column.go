package column

import (
	"prettyy-server-online/data/column"
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

func (c *Client) Add(needInsertToColumn map[string]string, uid int64) (err error) {
	for cid, title := range needInsertToColumn {
		col := &column.Column{
			Cid:          cid,
			Title:        title,
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

func Add(needInsertToColumn map[string]string, uid int64) (err error) {
	return defaultClient.Add(needInsertToColumn, uid)
}

func init() {
	var err error
	defaultClient, err = NewClient()
	if err != nil {
		panic(err)
	}
}
