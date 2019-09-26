package indexer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

func (i *Indexer) InstallMappings() *Indexer {
	log.Println("Install Mappings")
	files, err := ioutil.ReadDir(i.MappingDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", i.MappingDir, name))
		if err != nil {
			log.Fatal(err)
		}

		index := fmt.Sprintf("%s.%s", i.Network, name[0:len(name)-len(filepath.Ext(name))])
		if err = i.createIndex(index, b, i.Reindex); err != nil {
			log.Fatal(err)
		}
	}

	return i
}

func (i *Indexer) createIndex(index string, mapping []byte, purge bool) error {
	ctx := context.Background()
	client := i.Elastic.Client

	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}

	if exists && purge {
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
