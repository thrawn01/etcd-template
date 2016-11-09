package etcdTemplate_test

import (
	"context"
	"testing"
	"time"

	"bytes"
	"encoding/base32"
	"fmt"
	"os"
	"path"

	"strings"

	etcd "github.com/coreos/etcd/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
	"github.com/thrawn01/etcd-template"
)

func okToTestEtcd() {
	if os.Getenv("ETCD_ENDPOINTS") == "" {
		Skip("ETCD_ENDPOINTS not set, skipped....")
	}
}

func newRootPath() string {
	var buf bytes.Buffer
	encoder := base32.NewEncoder(base32.StdEncoding, &buf)
	encoder.Write(uuid.NewRandom())
	encoder.Close()
	buf.Truncate(26)
	return path.Join("/etcd-template-tests", buf.String())
}

func etcdClientFactory() etcd.Client {
	if os.Getenv("ETCD_ENDPOINTS") == "" {
		return nil
	}

	client, err := etcd.New(etcd.Config{
		Endpoints: strings.Split(os.Getenv("ETCD_ENDPOINTS"), ","),
	})
	if err != nil {
		Fail(fmt.Sprintf("etcdApiFactory() - %s", err.Error()))
	}
	return client
}

func etcdPut(client etcd.Client, root, key, value string) {
	// Context Timeout for 2 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	api := etcd.NewKeysAPI(client)
	// Set the value in the etcd store
	_, err := api.Set(ctx, path.Join(root, key), value, nil)
	if err != nil {
		Fail(fmt.Sprintf("etcdPut() - %s", err.Error()))
	}
}

func TestWatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WatcherTests")
}

var _ = Describe("Watcher", func() {
	var client etcd.Client
	var watcher *etcdTemplate.Watcher
	var etcdRoot string

	BeforeEach(func() {
		etcdRoot = newRootPath()
		client = etcdClientFactory()
		watcher = etcdTemplate.NewWatcher(client)
	})

	AfterEach(func() {
		if watcher != nil {
			watcher.Close()
		}
	})

	Describe("Watch()", func() {
		It("Should fetch new value from when etcd key changes", func() {
			okToTestEtcd()

			etcdPut(client, etcdRoot, "/environment/config", "some-value")
			watchChan := watcher.Watch(path.Join(etcdRoot, "/environment/config"))
			time.Sleep(time.Second)
			etcdPut(client, etcdRoot, "/environment/config", "some-new-value")
			newPair := <-watchChan
			Expect(newPair.Key).To(Equal(path.Join(etcdRoot, "/environment/config")))
			Expect(newPair.Value).To(Equal("some-new-value"))
		})
	})
})
