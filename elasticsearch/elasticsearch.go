package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/XC-Zero/zero_common/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/conflicts"
	"github.com/pkg/errors"
	"io"
	"log"
)

type EsIndex interface {
	// TableName ES 的索引名
	TableName() string
	// Mapping 返回 ES 的 mapping
	Mapping() Json
	// ToDoc 转为 ES 的一条记录
	ToDoc() Json
	// Search 搜索
	Search(searchContent string, offset, limit int) (total int64, data []Json, err error)
}

type Json map[string]any

const (
	IK_SMART    = "ik_smart"
	IK_MAX_WORD = "ik_max_word"
)

const (
	PRE_TAG  = "<hl>"
	POST_TAG = "</hl>"
)

type Client elasticsearch.TypedClient

// InitElasticsearch ...
func InitElasticsearch(config config.ElasticSearchConfig) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: config.Host,
		Username:  config.User,
		Password:  config.Password,
	}
	client, err := elasticsearch.NewTypedClient(cfg)

	if err != nil {
		panic(err)
	}

	return (*Client)(client), nil
}

func (c *Client) ExistIndex(ctx context.Context, doc EsIndex) (exist bool, err error) {
	return c.Indices.Exists(doc.TableName()).Do(ctx)
}

// CreateIndex 创建 Index
func (c *Client) CreateIndex(ctx context.Context, doc EsIndex) error {

	_, err := c.Indices.Create(doc.TableName()).Do(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	marshal, _ := json.Marshal(doc.Mapping())

	_, err = c.Indices.PutMapping(doc.TableName()).Raw(bytes.NewReader(marshal)).Do(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// DeleteIndex 删除 Index
func (c *Client) DeleteIndex(ctx context.Context, doc EsIndex) error {
	_, err := c.Indices.Delete(doc.TableName()).Do(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// CreateDoc 创建 Doc
func (c *Client) CreateDoc(ctx context.Context, doc EsIndex) error {
	_, err := c.Index(doc.TableName()).Request(doc.ToDoc()).Do(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpsertDoc 存在就更新不存在就创建
func (c *Client) UpsertDoc(ctx context.Context, uniqueID string, doc EsIndex) error {

	_, err := c.Index(doc.TableName()).Id(uniqueID).Document(doc.ToDoc()).Do(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil

}

// UpdateDoc term模式查询更新
func (c *Client) UpdateDoc(ctx context.Context, term map[string]string, doc EsIndex) error {
	document := doc.ToDoc()
	pattern := `ctx._source['%s'] = %v;`
	var source = ""
	for k, v := range document {
		if str, ok := v.(string); ok {
			source += fmt.Sprintf(pattern, k, `'`+str+`'`)
		}
	}
	var terms = make(map[string]types.TermQuery, len(term))
	for k, v := range term {
		terms[k] = types.TermQuery{Value: v}
	}
	_, err := c.UpdateByQuery(doc.TableName()).Query(&types.Query{
		Term: terms,
	}).Script(map[string]string{
		"source": source,
	}).Conflicts(conflicts.Proceed).Do(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	//fmt.Println(*do.Updated)
	return nil
}

// DeleteDoc term模式查询删除
func (c *Client) DeleteDoc(ctx context.Context, term map[string]string, doc EsIndex) error {

	var terms = make(map[string]types.TermQuery, len(term))
	for k, v := range term {
		terms[k] = types.TermQuery{Value: v}
	}
	_, err := c.DeleteByQuery(doc.TableName()).Query(&types.Query{
		Term: terms,
	}).Do(ctx)
	if err != nil {
		return err
	}
	//fmt.Println(do.Deleted)
	return nil
}

func (c *Client) DeleteDocByID(ctx context.Context, ID string, doc EsIndex) error {
	do, err := c.Delete(doc.TableName(), ID).Do(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	log.Printf(`[INFO]Delete from elasticsearch index %s  id is %s`, doc.TableName(), do.Id_)
	return nil

}

func (c *Client) SearchDoc(ctx context.Context, doc EsIndex, query io.Reader) (data [][]byte, hit int64, err error) {
	res, er := c.Search().Index(doc.TableName()).Raw(query).Do(ctx)
	if er != nil {
		err = er
		return
	}
	var result = make([][]byte, 0, len(res.Hits.Hits))
	hit = res.Hits.Total.Value
	for i := range res.Hits.Hits {
		hit := res.Hits.Hits[i]
		var b []byte
		err := hit.UnmarshalJSON(b)
		if err != nil {
			log.Println(err)
			continue
		}
		result = append(result, b)
	}
	data = result
	return
}

func Unmarshal(data *search.Response) (res []Json, err error) {
	for i := range data.Hits.Hits {

		hit := data.Hits.Hits[i]
		marshalJSON, err := hit.Source_.MarshalJSON()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		var temp Json
		err = json.Unmarshal(marshalJSON, &temp)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		for key, val := range hit.Highlight {
			if len(val) > 0 {
				// FIXME *.keyword 要考虑下怎么办
				temp[key] = val[0]
			}
		}
		temp["index"] = hit.Index_
		res = append(res, temp)
	}
	return
}
