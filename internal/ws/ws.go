package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WSEventHandler func(event string, payload json.RawMessage) error

var (
	handlerMap   map[string]WSEventHandler = map[string]WSEventHandler{}
	clients      map[*websocket.Conn]bool  = map[*websocket.Conn]bool{}
	clientsMutex sync.Mutex
)

var upgrader = websocket.Upgrader{}

func addConnection(conn *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	clients[conn] = true
}

func removeConnection(conn *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	delete(clients, conn)
}

func handleClient(conn *websocket.Conn) error {
	var resp Response

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("client closed")
			return err
		}

		err = json.Unmarshal(data, &resp)
		if err != nil {
			conn.WriteJSON(map[string]string{
				"error": "json unmarshal error",
			})
			continue
		}

		handler, ok := handlerMap[resp.Action]
		if !ok {
			conn.WriteJSON(map[string]string{
				"error": "no handler",
			})
			continue
		}

		err = handler(resp.Action, resp.Payload)
		if err != nil {
			conn.WriteJSON(map[string]string{
				"error": err.Error(),
			})
			continue
		}
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	addConnection(conn)
	defer removeConnection(conn)

	handleClient(conn)
}

func SetHandler(event string, handler WSEventHandler) {
	handlerMap[event] = handler
}

func BroadcastEvent(event string, payload interface{}) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(map[string]interface{}{
			"event":   event,
			"payload": payload,
		})

		if err != nil {
			log.Println("broadcast error:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func Start(wsAddr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleWebSocket)

	srv := &http.Server{
		Addr:    wsAddr,
		Handler: mux,
	}

	return srv.ListenAndServe()
}

func Close() {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		client.Close()
		delete(clients, client)
	}
}
