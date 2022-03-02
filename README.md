# Golang Multiclient Websocket Example

Try it on your local:
---
```
go run cmd/main.go
```


Address:
---
Connect to:
```
localhost:8000/ws
```


Socket Payload:
---
```
{
    "type": "message type",
    "message": "message interface"
}
```


Available Request Type:
---
1. **SUBSCRIBE**
```
{
    "type": "SUBSCRIBE",
    "message": "asset_id"
}

send this request to websocket to subscribe an asset.
websocket will update client detail and return subscribed asset data to the client every server intended delay time(10s).
if client already subscribed before, new subscribed asset_id will replace current config and send that asset data instead.
```

2. **UNSUBSCRIBE**
```
{
    "type": "UNSUBSCRIBE"
}

send this request to websocket to Unsubscribe an asset.
websocket will update client detail.
```


Available Asset ID:
---
1. BBCA
2. BUKA
3. TKPD


Reference:
---
- https://tutorialedge.net/projects/chat-system-in-go-and-react/part-4-handling-multiple-clients/
