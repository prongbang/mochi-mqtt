package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/auth"
	"github.com/mochi-co/mqtt/v2/listeners"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	server := mqtt.New(nil)

	authRules := &auth.Ledger{
		Auth: auth.AuthRules{ // Auth disallows all by default
			{Username: "innotech", Password: "password", Allow: true},
		},
		ACL: auth.ACLRules{
			{
				// user melon can read and write to their own topic
				Username: "innotech", Filters: auth.Filters{
					"#": auth.ReadWrite,
				},
			},
		},
	}

	// you may also find this useful...
	// d, _ := authRules.ToYAML()
	// d, _ := authRules.ToJSON()
	// fmt.Println(string(d))

	err := server.AddHook(new(auth.Hook), &auth.Options{
		Ledger: authRules,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = server.AddHook(new(CustomAuthHook), map[string]any{})
	if err != nil {
		log.Fatal(err)
	}

	tcp := listeners.NewTCP("t1", ":1883", nil)
	err = server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	ws := listeners.NewWebsocket("ws1", ":8083", nil)
	err = server.AddListener(ws)
	if err != nil {
		log.Fatal(err)
	}

	stats := listeners.NewHTTPStats("stats", ":8080", nil, server.Info)
	err = server.AddListener(stats)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-done
	server.Log.Warn().Msg("caught signal, stopping...")
	server.Close()
	server.Log.Info().Msg("main.go finished")
}
