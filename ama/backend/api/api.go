package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/mrocha98/go-studies/ama/backend/internal/store/pgstore"
)

type subscribersMap = map[string]map[*websocket.Conn]context.CancelFunc

type apiHandler struct {
	queries     *pgstore.Queries
	router      *chi.Mux
	upgrader    websocket.Upgrader
	subscribers subscribersMap
	mutex       *sync.Mutex
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func NewHandler(queries *pgstore.Queries) http.Handler {
	handler := apiHandler{
		queries: queries,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		subscribers: make(subscribersMap),
		mutex:       &sync.Mutex{},
	}

	r := chi.NewMux()

	r.Use(
		middleware.RequestID,
		middleware.Logger,
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
			AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           int((365 * 24 * time.Hour).Seconds()),
		}),
		middleware.Recoverer,
		middleware.Heartbeat("/api/health"),
	)

	r.Route("/ws", func(r chi.Router) {
		r.Get("/subscription/{room_id}", handler.handleSubscribeToRoom)
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", handler.handleCreateRoom)
			r.Get("/", handler.handleListRooms)

			r.Route("/{room_id}/messages", func(r chi.Router) {
				r.Post("/", handler.handleCreateRoomMessage)
				r.Get("/", handler.handleListRoomMessages)

				r.Route("/{message_id}", func(r chi.Router) {
					r.Get("/", handler.handleGetRoomMessage)
					r.Patch("/reaction", handler.handleReactToRoomMessage)
					r.Delete("/reaction", handler.handleRemoveReactionFromRoomMessage)
					r.Delete("/answer", handler.handleAnswerRoomMessage)
				})
			})
		})
	})

	handler.router = r
	return handler
}

const (
	MessageKindMessageCreated = "message_created"
)

type MessageMessageCreated struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type Message struct {
	Kind   string `json:"kind"`
	Value  any    `json:"value"`
	RoomId string `json:"-"`
}

func (h apiHandler) notifyClients(msg Message) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	subscribers, ok := h.subscribers[msg.RoomId]
	if !ok || len(subscribers) == 0 {
		return
	}

	for conn, cancel := range subscribers {
		if err := conn.WriteJSON(msg); err != nil {
			slog.Error("failed to send message to client", slog.Any("error", err))
			cancel()
		}
	}
}

func (h apiHandler) handleSubscribeToRoom(w http.ResponseWriter, r *http.Request) {
	rawRoomId := chi.URLParam(r, "room_id")
	roomId, err := uuid.Parse(rawRoomId)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	_, err = h.queries.GetRoom(r.Context(), roomId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return
		}

		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	wsConn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		const errorMessage = "failed to upgrade to websocket connection"
		slog.Warn(errorMessage, slog.Any("error", err))
		w.Header().Set("Upgrade", "websocket")
		http.Error(w, errorMessage, http.StatusUpgradeRequired)
		return
	}
	defer wsConn.Close()

	slog.Info("New client connected", slog.String("room_id", rawRoomId), slog.String("client_ip", r.RemoteAddr))

	ctx, cancel := context.WithCancel(r.Context())

	h.mutex.Lock()
	if _, ok := h.subscribers[rawRoomId]; !ok {
		h.subscribers[rawRoomId] = make(map[*websocket.Conn]context.CancelFunc)
	}
	h.subscribers[rawRoomId][wsConn] = cancel
	h.mutex.Unlock()

	<-ctx.Done()

	h.mutex.Lock()
	delete(h.subscribers[rawRoomId], wsConn)
	h.mutex.Unlock()
}

func (h apiHandler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	type _body = struct {
		Theme string `json:"theme"`
	}
	var body _body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	roomId, err := h.queries.InsertRoom(r.Context(), body.Theme)
	if err != nil {
		slog.Error("failed to insert room", slog.Any("error", err))
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	type response struct {
		Id string `json:"id"`
	}

	data, _ := json.Marshal(response{Id: roomId.String()})
	w.Header().Set("Content/Type", "application/json")
	w.Write(data)
}

func (h apiHandler) handleListRooms(w http.ResponseWriter, r *http.Request) {

}

func (h apiHandler) handleListRoomMessages(w http.ResponseWriter, r *http.Request) {

}

func (h apiHandler) handleCreateRoomMessage(w http.ResponseWriter, r *http.Request) {
	rawRoomId := chi.URLParam(r, "room_id")
	roomId, err := uuid.Parse(rawRoomId)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	_, err = h.queries.GetRoom(r.Context(), roomId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return
		}

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	type _body = struct {
		Message string `json:"message"`
	}
	var body _body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	messageId, err := h.queries.InsertMessage(r.Context(), pgstore.InsertMessageParams{
		RoomID:  roomId,
		Message: body.Message,
	})
	if err != nil {
		slog.Error("failed to insert message", slog.Any("error", err))
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	type response struct {
		Id string `json:"id"`
	}

	data, _ := json.Marshal(response{Id: messageId.String()})
	w.Header().Set("Content/Type", "application/json")
	w.Write(data)

	go h.notifyClients(Message{
		Kind:   MessageKindMessageCreated,
		RoomId: rawRoomId,
		Value: MessageMessageCreated{
			Id:      messageId.String(),
			Message: body.Message,
		},
	})
}

func (h apiHandler) handleGetRoomMessage(w http.ResponseWriter, r *http.Request) {

}

func (h apiHandler) handleReactToRoomMessage(w http.ResponseWriter, r *http.Request) {

}

func (h apiHandler) handleRemoveReactionFromRoomMessage(w http.ResponseWriter, r *http.Request) {

}

func (h apiHandler) handleAnswerRoomMessage(w http.ResponseWriter, r *http.Request) {

}
