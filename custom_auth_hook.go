package main

import (
	"bytes"
	"fmt"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
)

type CustomAuthHook struct {
	mqtt.HookBase
}

func (h *CustomAuthHook) ID() string {
	return "events-example"
}

func (h *CustomAuthHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnectAuthenticate,
		mqtt.OnACLCheck,
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribed,
	}, []byte{b})
}

// OnConnectAuthenticate returns true if the connecting client has rules which provide access in the auth ledger.
func (h *CustomAuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {

	AuthOk := func(cl *mqtt.Client, pk packets.Packet) (n int, ok bool) {
		username := string(cl.Properties.Username)
		password := string(pk.Connect.Password)
		fmt.Println("[OnConnectAuthenticate]username:", username)
		fmt.Println("[OnConnectAuthenticate]password:", password)

		// Verify JWT
		if username == "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwicm9sZXMiOlsiYWRtaW4iXX0.orOzAgQ6dneDCLN0cssDexmbnsqiLB_XseDsbWpMoe4" {
			return 0, true
		}

		return 0, false
	}

	if _, ok := AuthOk(cl, pk); ok {
		return true
	}

	h.Log.Info().
		Str("username", string(pk.Connect.Username)).
		Str("remote", cl.Net.Remote).
		Msg("client failed authentication check")

	return false
}

// OnACLCheck returns true if the connecting client has matching read or write access to subscribe or publish to a given topic.
func (h *CustomAuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {

	ACLOk := func(cl *mqtt.Client, topic string, write bool) (n int, ok bool) {
		username := string(cl.Properties.Username)
		fmt.Println("[OnACLCheck]username:", username)
		fmt.Println("[OnACLCheck]topic:", topic)
		fmt.Println("[OnACLCheck]write:", write)

		// TODO check can access topic by username JWT

		return 0, true
	}

	if _, ok := ACLOk(cl, topic, write); ok {
		return true
	}

	h.Log.Debug().
		Str("client", cl.ID).
		Str("username", string(cl.Properties.Username)).
		Str("topic", topic).
		Msg("client failed allowed ACL check")

	return false
}

func (h *CustomAuthHook) OnConnect(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Msgf("client connected")
}

func (h *CustomAuthHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Info().Str("client", cl.ID).Bool("expire", expire).Err(err).Msg("client disconnected")
}

func (h *CustomAuthHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msgf("subscribed qos=%v", reasonCodes)
}

func (h *CustomAuthHook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msg("unsubscribed")
}
