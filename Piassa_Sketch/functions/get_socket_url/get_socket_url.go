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

//export GetSketchSocket
func getsocketurl(e event.Event) uint32 {
    h, err := e.HTTP()
    if err != nil {
        return 1
    }

    // Get room name from query (?room=piassa-room-1)
    room, err := h.Query().Get("room")
    if err != nil {
        return fail(h, err, 500)
    }

    // Hash the room name to ensure valid channel naming
    hash := md5.New()
    hash.Write([]byte(room))
    roomHash := hex.EncodeToString(hash.Sum(nil))

    // Open the channel sketch-battle-<hash>
    channel, err := pubsub.Channel("sketch-battle-" + roomHash)
    if err != nil {
        return fail(h, err, 500)
    }

    // Generate the relay URL
    url, err := channel.WebSocket().Url()
    if err != nil {
        return fail(h, err, 500)
    }

    // Return the path so the Next.js app can connect
    h.Write([]byte(url.Path))
    return 0
}