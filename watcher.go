package etcdTemplate

import (
	"context"

	"sync"

	etcd "github.com/coreos/etcd/client"
	"github.com/mailgun/log"
)

type Watcher struct {
	client     etcd.Client
	api        etcd.KeysAPI
	cancels    []context.CancelFunc
	changeChan chan Pair
}

func NewWatcher(client etcd.Client) *Watcher {
	return &Watcher{
		client:     client,
		api:        etcd.NewKeysAPI(client),
		changeChan: make(chan Pair, 1),
	}
}

type Pair struct {
	Key   string
	Value string
}

// Get retrieves a value from a K/V store for the provided key.
func (self *Watcher) Get(ctx context.Context, key string) (Pair, error) {
	resp, err := self.api.Get(ctx, key, nil)
	if err != nil {
		return Pair{}, err
	}
	return Pair{Key: string(resp.Node.Key), Value: resp.Node.Value}, nil
}

// Watch monitors store for changes to key path.
func (self *Watcher) Watch(path string) <-chan Pair {
	var isRunning sync.WaitGroup
	var once sync.Once

	watcher := self.api.Watcher(path, &etcd.WatcherOptions{AfterIndex: 0, Recursive: false})
	ctx, cancel := context.WithCancel(context.Background())
	self.cancels = append(self.cancels, cancel)

	isRunning.Add(1)
	go func() {
		for {
			once.Do(func() { isRunning.Done() })
			r, err := watcher.Next(ctx)
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Errorf("etcd watcher.Next() - %s", err)
				continue
			}
			if r.Action == "get" || r.Action == "delete" {
				continue
			}
			pair, err := self.Get(ctx, path)
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Errorf("etcd self.Get(%s) - %s", path, err)
				continue
			}
			self.changeChan <- pair
		}
	}()
	// Wait until the go-routine is running before returning
	isRunning.Wait()
	return self.changeChan
}

func (self *Watcher) Close() {
	if len(self.cancels) != 0 {
		for _, cancel := range self.cancels {
			cancel()
		}
	}
}
