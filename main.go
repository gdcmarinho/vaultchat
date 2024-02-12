package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	UUID      string
	PublicKey string
}

type Session struct {
	UUID         string
	Participants []string
}

type Message struct {
	UUID          string
	SessionUUID   string
	SenderUUID    string
	RecipientUUID string
	Content       string
	SentAt        time.Time
	Signature     string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var sessions map[string]*Session = make(map[string]*Session)
var sessionMutex sync.Mutex

var userConnections map[string]*websocket.Conn = make(map[string]*websocket.Conn)
var userConnectionsMutex sync.Mutex

func main() {
	http.HandleFunc("/session", newSessionHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func newSessionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for new session")
	sessionUUID, err := generateUUID()
	if err != nil {
		log.Println("Error while trying to generate UUID:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Start a new session
	session := &Session{
		UUID:         sessionUUID,
		Participants: []string{},
	}

	// Add session to map
	sessionMutex.Lock()
	sessions[sessionUUID] = session
	sessionMutex.Unlock()

	// Upgrade to websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error while trying to upgrade:", err)
		return
	}

	// Create goroutine to manage session
	go handleSession(ws, session)
}

func handleSession(ws *websocket.Conn, session *Session) {
	// Read websocket message
	for {
		message, err := readMessage(ws)
		if err != nil {
			log.Println("Erro ao ler mensagem:", err)
			return
		}

		// Forward message to all participants
		for _, participant := range session.Participants {
			if participant != message.SenderUUID {
				sendMessage(participant, message)
			}
		}
	}

	// Remove session when websocket connection is closed
	sessionMutex.Lock()
	delete(sessions, session.UUID)
	sessionMutex.Unlock()
}

func readMessage(ws *websocket.Conn) (*Message, error) {
	// Decode JSON
	message := &Message{}
	err := ws.ReadJSON(message)
	if err != nil {
		return nil, err
	}

	log.Printf("Received message: UUID=%s, SessionUUID=%s, SenderUUID=%s, RecipientUUID=%s, Content=%s, SentAt=%s, Signature=%s\n",
		message.UUID, message.SessionUUID, message.SenderUUID, message.RecipientUUID, message.Content, message.SentAt.Format(time.RFC3339), message.Signature)

	// TODO: Verify message signature

	return message, nil
}

func sendMessage(participant string, message *Message) {
	userConnectionsMutex.Lock()
	ws, ok := userConnections[participant]
	userConnectionsMutex.Unlock()

	if !ok {
		log.Println("WebSocket connection not found for participant:", participant)
		return
	}

	// A mensagem pode precisar ser adaptada ou serializada antes do envio
	err := ws.WriteJSON(message)
	if err != nil {
		log.Println("Error sending message to participant:", participant, err)
	}
}

func generateUUID() (string, error) {
	log.Println("Generating UUID")

	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	uuid[0] = (uuid[0] | 0x40) & 0x7F

	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
