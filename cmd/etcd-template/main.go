package main

import (
	"fmt"
	"os"
	"time"

	"os/signal"

	"context"

	etcd "github.com/coreos/etcd/client"
	"github.com/thrawn01/args"
	"github.com/thrawn01/etcd-template"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	parser := args.NewParser(args.Name("etcd-template"),
		args.Desc("Read mailgun compatable etcd dictionaries from etcd and generate files from a template"))

	parser.AddOption("--watch").Alias("-w").IsTrue().
		Help("Watches the specified etcd key for changes and regenerates templates if the key value changes")
	parser.AddOption("--etcd-endpoints").Alias("-e").Default("http://localhost:2379").Env("ETCD_ENDPOINTS").
		Help("A Comma Separated list of etcd server endpoints")
	parser.AddArgument("etcd-path").Required().
		Help("The etcd path to the key where our config is stored")
	parser.AddArgument("template-dir").Required().
		Help("The directory where template files suffixed with .tpl are located")
	parser.AddArgument("output-dir").
		Help("Output directory for generated files (Defaults to template-dir if not provided)")

	options := parser.ParseArgsSimple(nil)

	client, err := etcd.New(etcd.Config{Endpoints: options.StringSlice("etcd-endpoints")})
	checkErr(err)

	etcdPath := options.String("etcd-path")
	watcher := etcdTemplate.NewWatcher(client)

	// Get the config from etcd
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	pair, err := watcher.Get(ctx, etcdPath)
	checkErr(err)

	// Build our files from the templates
	err = etcdTemplate.Generate(options, pair)
	checkErr(err)

	if !options.IsSeen("watch") {
		os.Exit(0)
	}

	watchChan := watcher.Watch(etcdPath)

	done := make(chan struct{})
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, os.Kill)
		sig := <-signalChan
		fmt.Printf("Captured %v. Exiting...", sig)
		watcher.Close()
		close(done)
	}()

	select {
	case pair := <-watchChan:
		fmt.Printf("%s Config Updated", pair.Key)
		if err := etcdTemplate.Generate(options, pair); err != nil {
			fmt.Printf("Error %s\n", err)
		}
	case <-done:
		os.Exit(1)
	}
}
