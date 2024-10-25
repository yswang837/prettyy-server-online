package article

import (
	"errors"
	"prettyy-server-online/data/article"
	invertedIndex "prettyy-server-online/data/inverted-index"
	invertedIndex2 "prettyy-server-online/services/inverted-index"
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
	if a.Visibility == "" {
		a.Visibility = "1"
	}
	if a.Typ == "" {
		a.Typ = "1"
	}
	// 发文的时候，这些数都是0
	a.ShareNum = 0
	a.CommentNum = 0
	a.LikeNum = 0
	a.ReadNum = 0
	a.CollectNum = 0
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

func IncrReadNum(aid string) error {
	return defaultClient.IncrReadNum(aid)
}

func (c *Client) IncrReadNum(aid string) error {
	if aid == "" {
		return tool.ErrParams
	}
	a, err := c.Get(aid)
	if err != nil {
		return err
	}
	a.ReadNum++
	m := map[string]interface{}{"read_num": a.ReadNum}
	if _, err = c.cacheManager.HMSet(aid, m); err != nil {
		return errors.New("redis incr read num failed: " + err.Error())
	}
	if err = c.manager.IncrReadNum(aid); err != nil {
		return errors.New("mysql incr read num failed: " + err.Error())
	}
	return nil
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

func UpdateLikeNum(aid string, isAddClick bool) (int, error) {
	return defaultClient.UpdateLikeNum(aid, isAddClick)
}

func (c *Client) UpdateLikeNum(aid string, isAddClick bool) (int, error) {
	if aid == "" {
		return 0, tool.ErrParams
	}
	a, err := c.Get(aid)
	if err != nil {
		return 0, err
	}
	if isAddClick {
		a.LikeNum++
	} else {
		a.LikeNum--
	}
	m := map[string]interface{}{"like_num": a.LikeNum}
	if _, err = c.cacheManager.HMSet(aid, m); err != nil {
		return 0, errors.New("redis update like num failed: " + err.Error())
	}
	if isAddClick {
		if err = c.manager.IncrLikeNum(aid); err != nil {
			return 0, errors.New("mysql incr like num failed: " + err.Error())
		}
	} else {
		if err = c.manager.DeIncrLikeNum(aid); err != nil {
			return 0, errors.New("mysql deincr like num failed: " + err.Error())
		}
	}
	return a.LikeNum, nil
}

func UpdateCollectNum(aid string, isAddClick bool) (int, error) {
	return defaultClient.UpdateCollectNum(aid, isAddClick)
}

func (c *Client) UpdateCollectNum(aid string, isAddClick bool) (int, error) {
	if aid == "" {
		return 0, tool.ErrParams
	}
	a, err := c.Get(aid)
	if err != nil {
		return 0, err
	}
	if isAddClick {
		a.CollectNum++
	} else {
		a.CollectNum--
	}
	m := map[string]interface{}{"collect_num": a.CollectNum}
	if _, err = c.cacheManager.HMSet(aid, m); err != nil {
		return 0, errors.New("redis update collect num failed: " + err.Error())
	}
	if isAddClick {
		if err = c.manager.IncrCollectNum(aid); err != nil {
			return 0, errors.New("mysql incr collect num failed: " + err.Error())
		}
	} else {
		if err = c.manager.DeIncrCollectNum(aid); err != nil {
			return 0, errors.New("mysql deincr collect num failed: " + err.Error())
		}
	}
	return a.CollectNum, nil
}

func GetArticleList(uid int64, page, pageSize int, visibility, typ string) ([]*article.Article, int64, error) {
	return defaultClient.GetArticleList(uid, page, pageSize, visibility, typ)
}

func (c *Client) GetArticleList(uid int64, page, pageSize int, visibility, typ string) ([]*article.Article, int64, error) {
	// 如果uid合法，则通过uid查询aid列表，否则随机查询一个article表，前者用于管理我的文章，后者用于主页显示，后者后期可优化成推荐
	if uid >= 10000 {
		iList, err := invertedIndex2.Get(invertedIndex.TypUidAid, strconv.FormatInt(uid, 10))
		if err != nil {
			return nil, 0, err
		}
		var aids []string
		for _, i := range iList {
			if i.Idx != "" {
				aids = append(aids, i.Idx)
			}
		}
		if len(aids) <= 0 {
			return nil, 0, errors.New("inverted index not found")
		}
		return c.manager.GetContentManageArticleList(aids, visibility, typ)
	} else {
		// 主页仅展示”全部可见“的文章
		articleList, err := c.manager.GetHomeArticleList(page, pageSize, "1")
		if err != nil {
			return nil, 0, errors.New("get article list from mysql failed: " + err.Error())
		}
		for _, art := range articleList {
			art.Content = tool.Base64Decode(art.Content)
		}
		return articleList, 0, nil
	}
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
	m["tags"] = a.Tags
	m["visibility"] = a.Visibility
	m["typ"] = a.Typ
	m["share_num"] = a.ShareNum
	m["comment_num"] = a.CommentNum
	m["like_num"] = a.LikeNum
	m["read_num"] = a.ReadNum
	m["collect_num"] = a.CollectNum
	m["uid"] = a.Uid
	m["status"] = a.Status
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
	sn, _ := strconv.Atoi(m["share_num"])
	ln, _ := strconv.Atoi(m["like_num"])
	uid, _ := strconv.Atoi(m["uid"])
	return &article.Article{
		Aid:        m["aid"],
		Title:      m["title"],
		Content:    tool.Base64Decode(m["content"]),
		CoverImg:   m["cover_img"],
		Summary:    m["summary"],
		Tags:       m["tags"],
		Visibility: m["visibility"],
		Typ:        m["typ"],
		ShareNum:   sn,
		CommentNum: cn,
		LikeNum:    ln,
		ReadNum:    rn,
		CollectNum: con,
		Status:     m["status"],
		Uid:        int64(uid),
		CreateTime: tool.StringToTime(m["create_time"]),
		UpdateTime: tool.StringToTime(m["update_time"]),
	}
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
}
