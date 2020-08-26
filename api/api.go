package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/fasmide/go-mess/database"
	"github.com/gorilla/websocket"
	"github.com/r3labs/diff"
)

var upgrader = websocket.Upgrader{} // use default options

// API provides a http API
type API struct {
	http.ServeMux
	sync.RWMutex

	Db *database.Database

	clients []*websocket.Conn
	loop    sync.Once
}

// New initializes a new *API
func New(d *database.Database) *API {
	a := &API{
		Db:      d,
		clients: make([]*websocket.Conn, 0),
	}

	a.HandleFunc("/active", a.active)
	a.HandleFunc("/previous", a.previous)
	a.HandleFunc("/changes", a.changes)
	return a
}

func (a *API) changeLoop() {
	current, err := a.Db.ActiveOrders()
	if err != nil {
		err = fmt.Errorf("could not lookup active orders: %w", err)
		log.Printf(err.Error())
		return
	}

	for {
		time.Sleep(time.Second)
		next, err := a.Db.ActiveOrders()
		if err != nil {
			err = fmt.Errorf("could not lookup active orders: %w", err)
			log.Printf(err.Error())
			return
		}

		set, err := diff.Diff(current, next)
		if err != nil {
			log.Printf("unable to calculate change set: %s", err)
			return
		}

		go a.emit(set)
		current = next
	}
}

func (a *API) emit(s diff.Changelog) {
	payload, err := json.Marshal(s)
	if err != nil {
		log.Printf("unable to marshal change set: %s", err)
		return
	}

	a.RLock()
	defer a.RUnlock()

	for _, c := range a.clients {
		go func(c *websocket.Conn) {
			err := c.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				log.Printf("unable to write data to websocket peer: %s", err)
				go a.removeClient(c)
			}
		}(c)
	}
}

func (a *API) removeClient(c *websocket.Conn) {
	a.Lock()
	defer a.Unlock()

	// find index of client
	var i int
	var x *websocket.Conn
	var found bool
	for i, x = range a.clients {
		if x == c {
			found = true
			break
		}
	}

	// another routine may have already removed this
	if !found {
		return
	}

	_ = x.Close()

	// move the client and cut him off the slice
	a.clients[len(a.clients)-1], a.clients[i] = a.clients[i], a.clients[len(a.clients)-1]
	a.clients = a.clients[:len(a.clients)-1]

}

func (a *API) changes(w http.ResponseWriter, r *http.Request) {

	// startup the change loop
	a.loop.Do(a.changeLoop)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		err = fmt.Errorf("http websocket upgrade failed: %w", err)
		log.Print(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.Lock()
	a.clients = append(a.clients, c)
	a.Unlock()

	// throw away incoming messages
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
	}

	a.removeClient(c)
}

func (a *API) active(w http.ResponseWriter, r *http.Request) {
	data, err := a.Db.ActiveOrders()
	if err != nil {
		err = fmt.Errorf("could not lookup active orders: %w", err)
		log.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		log.Printf("could not json encode active orders: %s", err)
	}
}

func (a *API) previous(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	ids, exists := values["id"]
	if !exists {
		http.Error(w, "no id query parameters provided", http.StatusBadRequest)
		return
	}

	data, err := a.Db.PreviousOrders(ids...)
	if err != nil {
		errStr := fmt.Sprintf("could not lookup previous orders: %s", err)
		log.Printf(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(data)
	if err != nil {
		log.Printf("could not encode previous orders: %s", err)
	}
}
