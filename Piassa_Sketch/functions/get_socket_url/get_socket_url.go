package lib

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/taubyte/go-sdk/event"
	http "github.com/taubyte/go-sdk/http/event"
	pubsub "github.com/taubyte/go-sdk/pubsub/node"
)

func fail(h http.Event, err error, code int) uint32 {
	h.Write([]byte(err.Error()))
	h.Return(code)
	return 1
}

//export getsocketurl
func getsocketurl(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		return 1
	}

	// 1. Get room from query (?room=piassa1)
	room, err := h.Query().Get("room")
	if err != nil {
		return fail(h, err, 500)
	}

	// 2. Hash the room to create a unique internal channel name
	hash := md5.New()
	hash.Write([]byte(room))
	roomHash := hex.EncodeToString(hash.Sum(nil))

	// 3. Open the channel (Must match the regex in messaging.yml)
	channel, err := pubsub.Channel("sketch-battle-" + roomHash)
	if err != nil {
		return fail(h, err, 500)
	}

	// 4. Get the websocket relay URL
	url, err := channel.WebSocket().Url()
	if err != nil {
		return fail(h, err, 500)
	}

	// 5. Write the path to the response
	h.Write([]byte(url.Path))

	return 0
}