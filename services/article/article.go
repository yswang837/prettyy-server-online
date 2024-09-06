package article

import (
	"errors"
	xzfSnowflake "prettyy-server-online/custom-pkg/xzf-snowflake"
	"prettyy-server-online/data/article"
	"prettyy-server-online/utils/tool"
	"strconv"
	"time"
)

type Client struct {
	manager      *article.Manager
	cacheManager *article.ManagerRedis
}

var defaultClient *Client

func NewClient() (*Client, error) {
	manager, err := article.NewManager()
	if err != nil {
		return nil, err
	}
	cache, err := article.NewManagerRedis()
	if err != nil {
		return nil, err
	}
	return &Client{manager: manager, cacheManager: cache}, nil
}

func (c *Client) Add(a *article.Article) (err error) {
	if a == nil {
		return tool.ErrParams
	}
	a.Aid = xzfSnowflake.GenID("AA")
	if err = c.manager.Add(a); err != nil {
		return errors.New("add article to mysql failed: " + err.Error())
	}
	a.CreateTime = time.Now()
	a.UpdateTime = time.Now()
	if _, err = c.cacheManager.HMSet(a.Aid, articleToMap(a)); err != nil {
		return errors.New("set article to redis failed: " + err.Error())
	}
	return
}

func (c *Client) Get(aid string) (*article.Article, error) {
	if aid == "" {
		return nil, tool.ErrParams
	}
	articleMap, err := c.cacheManager.HGetAll(aid)
	if err != nil {
		return nil, errors.New("get article from redis failed: " + err.Error())
	}
	if len(articleMap) != 0 {
		return mapToArticle(articleMap), nil
	}
	art, err := c.manager.Get(aid)
	if err != nil {
		return nil, errors.New("get article from mysql failed: " + err.Error())
	}
	art.Content = tool.Base64Decode(art.Content)
	return art, nil
}

func Get(aid string) (*article.Article, error) {
	return defaultClient.Get(aid)
}

func GetArticleList(uid int64, page, pageSize int, visibility, typ string) ([]*article.Article, error) {
	return defaultClient.GetArticleList(uid, page, pageSize, visibility, typ)
}

func (c *Client) GetArticleList(uid int64, page, pageSize int, visibility, typ string) ([]*article.Article, error) {
	articleList, err := c.manager.GetArticleList(uid, page, pageSize, visibility, typ)
	if err != nil {
		return nil, errors.New("get article list from mysql failed: " + err.Error())
	}
	for _, art := range articleList {
		art.Content = tool.Base64Decode(art.Content)
	}
	return articleList, nil
}

func articleToMap(a *article.Article) map[string]interface{} {
	if a == nil {
		return nil
	}
	m := make(map[string]interface{})
	m["aid"] = a.Aid
	m["title"] = a.Title
	m["content"] = a.Content
	m["cover_img"] = a.CoverImg
	m["summary"] = a.Summary
	m["uid"] = a.Uid
	m["create_time"] = a.CreateTime.Format(tool.DefaultDateTimeLayout)
	m["update_time"] = a.UpdateTime.Format(tool.DefaultDateTimeLayout)
	return m
}

func mapToArticle(m map[string]string) *article.Article {
	if len(m) == 0 {
		return nil
	}
	rn, _ := strconv.Atoi(m["read_num"])
	cn, _ := strconv.Atoi(m["comment_num"])
	con, _ := strconv.Atoi(m["collect_num"])
	uid, _ := strconv.Atoi(m["uid"])
	a := &article.Article{}
	a.Aid = m["aid"]
	a.Title = m["title"]
	a.Content = tool.Base64Decode(m["content"])
	a.CoverImg = m["cover_img"]
	a.Summary = m["summary"]
	a.ReadNum = rn
	a.CommentNum = cn
	a.CollectNum = con
	a.Uid = int64(uid)
	a.CreateTime = tool.StringToTime(m["create_time"])
	a.UpdateTime = tool.StringToTime(m["update_time"])
	return a
}

func Add(a *article.Article) (err error) {
	return defaultClient.Add(a)
}

func Delete(aid string, uid int64) (err error) {
	return defaultClient.Delete(aid, uid)
}

func (c *Client) Delete(aid string, uid int64) (err error) {
	if aid == "" {
		return tool.ErrParams
	}
	if err = c.manager.Delete(aid, uid); err != nil {
		return errors.New("del article from mysql failed: " + err.Error())
	}
	if _, err = c.cacheManager.Del(aid); err != nil {
		return errors.New("del article from redis failed: " + err.Error())
	}
	return nil
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
