package index

import (
	"context"
	"encoding/json"
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
	Id     string
	Doc    interface{}
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
		log.Printf("Creating elastic search index: %s", index)

		if err = i.createIndex(index, b); err != nil {
			return err
		}
	}

	return nil
}

func (i *Index) AddRequest(index string, id string, doc interface{}) {
	if request := i.GetRequest(index, id); request != nil {
		request.Doc = doc
	} else {
		i.requests = append(i.requests, Request{
			Index: index,
			Id:    id,
			Doc:   doc,
		})
	}
}

func (i *Index) GetRequest(index string, id string) *Request {
	for idx, r := range i.requests {
		if r.Index == index && r.Id == id {
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
		bulk.Add(elastic.NewBulkIndexRequest().Index(r.Index).Id(r.Id).Doc(r.Doc))
	}
	actions := bulk.NumberOfActions()
	if actions != 0 {
		_, err := bulk.Do(context.Background())
		if err != nil {
			logrus.WithError(err).Fatal("Failed to persist request at height ", height)
		}

		logrus.WithFields(logrus.Fields{"actions": actions}).Info("Indexed height ", height)
	}

	i.requests = make([]Request, 0)
}

func (i *Index) GetById(index string, id string, record interface{}) error {
	result, err := i.Client.
		Get().
		Index(index).
		Id(id).
		Do(context.Background())
	if err != nil {
		return ErrRecordNotFound
	}

	return json.Unmarshal(result.Source, record)
}

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
