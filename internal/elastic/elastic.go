package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"strings"
)

type Elastic struct {
	Client *elastic.Client
}

type Config struct {
	Value string `json:"value"`
}

var (
	ErrRecordNotFound     = errors.New("Record not found")
	ErrConfigNotFound     = errors.New("Config not found")
	ErrDatabaseConnection = errors.New("Could not connect to the search cluster")
	ErrNoResultFound      = errors.New("No results found")
)

func New(hosts []string, sniff bool, healthCheck bool, debug bool) (*Elastic, error) {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(strings.Join(hosts, ",")),
		elastic.SetSniff(sniff),
		elastic.SetHealthcheck(healthCheck),
	}

	if debug {
		opts = append(opts, elastic.SetTraceLog(log.New(os.Stdout, "", 0)))
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		log.Print("Error: ", err)
		return nil, ErrDatabaseConnection
	}

	return &Elastic{Client: client}, nil
}

func (e *Elastic) Flush(indexes ...string) *elastic.IndicesFlushService {
	return e.Client.Flush(indexes...)
}

func (e *Elastic) GetConfig(index string, key string) ([]byte, error) {
	result, err := e.Client.
		Get().
		Index(index).
		Id(key).
		Do(context.Background())

	if err != nil {
		return nil, ErrConfigNotFound
	}

	var config Config
	if err := json.Unmarshal(result.Source, &config); err != nil {
		return nil, err
	}

	return []byte(config.Value), nil
}

func (e *Elastic) SetConfig(index string, key string, val []byte) error {
	_, err := e.Client.
		Index().
		Index(index).
		Id(key).
		BodyString(fmt.Sprintf("{\"value\": \"%s\"}", string(val))).
		Do(context.Background())

	return err
}

func (e *Elastic) GetById(index string, id string, record interface{}) error {
	result, err := e.Client.
		Get().
		Index(index).
		Id(id).
		Do(context.Background())

	if err != nil {
		return ErrRecordNotFound
	}

	return json.Unmarshal(result.Source, record)
}

func (e *Elastic) GetIndexHeight() (uint64, error) {
	return 0, nil
}
