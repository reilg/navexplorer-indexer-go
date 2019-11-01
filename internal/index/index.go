package index

import (
	"context"
	"errors"
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Index struct {
	Client        *elastic.Client
	requests      []Request
	bulkIndexSize uint
}

type Request struct {
	Height uint64
	Index  string
	Name   string
	Doc    interface{}
	Type   string
	Id     string
}

var (
	ErrRecordNotFound = errors.New("Record not found")
)

func New() (*Index, error) {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(strings.Join(config.Get().ElasticSearch.Hosts, ",")),
		elastic.SetSniff(config.Get().ElasticSearch.Sniff),
		elastic.SetHealthcheck(config.Get().ElasticSearch.HealthCheck),
	}

	if config.Get().ElasticSearch.Username != "" {
		opts = append(opts, elastic.SetBasicAuth(
			config.Get().ElasticSearch.Username,
			config.Get().ElasticSearch.Password,
		))
	}

	if config.Get().ElasticSearch.Debug {
		opts = append(opts, elastic.SetTraceLog(log.New(os.Stdout, "", 0)))
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		log.Println("Error: ", err)
	}

	return &Index{
		Client:        client,
		requests:      make([]Request, 0),
		bulkIndexSize: config.Get().BulkIndexSize,
	}, err
}

func (i *Index) Init() error {
	log.Println("Install Mappings")
	files, err := ioutil.ReadDir(config.Get().ElasticSearch.MappingDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", config.Get().ElasticSearch.MappingDir, name))
		if err != nil {
			return err
		}

		index := fmt.Sprintf("%s.%s", config.Get().Network, name[0:len(name)-len(filepath.Ext(name))])
		if err = i.createIndex(index, b); err != nil {
			return err
		}
	}

	return nil
}

func (i *Index) AddIndexRequest(index string, name string, doc interface{}) {
	i.AddRequest(index, name, doc, "index", "")
}

func (i *Index) AddUpdateRequest(index string, name string, doc interface{}, id string) {
	i.AddRequest(index, name, doc, "update", id)
}

func (i *Index) AddRequest(index string, name string, doc interface{}, reqType string, id string) {
	if request := i.GetRequest(index, name, id); request != nil {
		request.Doc = doc
	} else {
		i.requests = append(i.requests, Request{
			Index: index,
			Name:  name,
			Doc:   doc,
			Type:  reqType,
			Id:    id,
		})
	}
}

func (i *Index) GetRequest(index string, name string, id string) *Request {
	for idx, r := range i.requests {
		if r.Index == index && r.Name == name && r.Id == id {
			return &i.requests[idx]
		}
	}

	return nil
}

func (i *Index) PersistRequests(height uint64) {
	if height%uint64(i.bulkIndexSize) != 0 || len(i.requests) == 0 {
		return
	}

	bulk := i.Client.Bulk()
	for _, r := range i.requests {
		if r.Type == "index" {
			bulk.Add(elastic.NewBulkIndexRequest().Index(r.Index).Doc(r.Doc))
		} else if r.Type == "update" {
			bulk.Add(elastic.NewBulkUpdateRequest().Index(r.Index).Id(r.Id).Doc(r.Doc))
		}
	}

	actions := bulk.NumberOfActions()
	if actions != 0 {
		_, err := bulk.Do(context.Background())
		if err != nil {
			logrus.WithError(err).Fatal("Failed to persist request at height ", height)
		}
	}

	logrus.WithFields(logrus.Fields{"actions": actions}).Info("Indexed height ", height)

	i.requests = make([]Request, 0)
}

func (i *Index) RewindTo(height uint64) *Index {
	rewind := func(index string) {
		_, err := i.Client.DeleteByQuery().
			Index(index).
			Query(elastic.NewRangeQuery("height").Gt(height)).
			Do(context.Background())
		if err != nil {
			logrus.WithError(err).Fatal("Could not rewind block index to ", height)
		}
	}

	rewind(BlockIndex.Get())
	rewind(BlockTransactionIndex.Get())
	rewind(AddressTransactionIndex.Get())
	rewind(SignalIndex.Get())
	rewind(ProposalIndex.Get())
	rewind(PaymentRequestIndex.Get())

	return i
}

//func (i *Index) GetById(index string, id string, record interface{}) error {
//	result, err := i.Client.
//		Get().
//		Index(index).
//		Id(id).
//		Do(context.Background())
//	if err != nil {
//		return ErrRecordNotFound
//	}
//
//	return json.Unmarshal(result.Source, record)
//}

func (i *Index) createIndex(index string, mapping []byte) error {
	ctx := context.Background()
	client := i.Client

	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}

	if exists && config.Get().Reindex {
		_, err = client.DeleteIndex(index).Do(ctx)
		if err != nil {
			return err
		}
	}

	if !exists {
		createIndex, err := client.CreateIndex(index).BodyString(string(mapping)).Do(ctx)
		if err != nil {
			return err
		}

		if createIndex.Acknowledged {
			logrus.Info("Created index: ", index)
		}
	}

	return nil
}
