package index

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/config"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Index struct {
	Client *elastic.Client
}

var (
	ErrRecordNotFound     = errors.New("Record not found")
	ErrDatabaseConnection = errors.New("Could not connect to the search cluster")
)

func New() (*Index, error) {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(strings.Join(config.Get().ElasticSearch.Hosts, ",")),
		elastic.SetSniff(config.Get().ElasticSearch.Sniff),
		elastic.SetHealthcheck(config.Get().ElasticSearch.HealthCheck),
	}

	if config.Get().Debug {
		opts = append(opts, elastic.SetTraceLog(log.New(os.Stdout, "", 0)))
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		log.Print("Error: ", err)
		return nil, ErrDatabaseConnection
	}

	return &Index{client}, nil
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

func (i *Index) Flush(indexes ...string) *elastic.IndicesFlushService {
	return i.Client.Flush(indexes...)
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

		if !createIndex.Acknowledged {
			log.Printf("Created index: %s", index)
		}
	}

	return nil
}
