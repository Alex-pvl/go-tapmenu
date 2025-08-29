package tapmenu

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alex-pvl/go-tapmenu/internal/app/config"
	"github.com/alex-pvl/go-tapmenu/internal/app/store"
	"github.com/alex-pvl/go-tapmenu/internal/app/tapmenu/kafka"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	configuration *config.Configuration
	logger        *logrus.Logger
	router        *mux.Router
	db            *store.Store
	producer      *kafka.Producer
}

type Response struct {
	MenuUrl     string `json:"menu_url"`
	TableNumber uint   `json:"table_number"`
}

func New(
	configuration *config.Configuration,
	db *store.Store,
	producer *kafka.Producer,
) *Server {
	return &Server{
		configuration: configuration,
		logger:        logrus.New(),
		router:        mux.NewRouter(),
		db:            db,
		producer:      producer,
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configureRouter()

	s.logger.Info("starting server on port ", s.configuration.BindAddress)

	return http.ListenAndServe(s.configuration.BindAddress, s.router)
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.configuration.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	return nil
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/tapmenu/{hash}", s.handleHash()).Methods(http.MethodGet)
	s.router.HandleFunc("/tapmenu/{hash}/call", s.handleCall()).Methods(http.MethodPost)
}

func (s *Server) handleHash() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := mux.Vars(r)["hash"]
		table, err := s.db.GetTable(hash)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		s.logger.Info(*table)

		renderJSON(w, *table)
	}
}

func (s *Server) handleCall() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := mux.Vars(r)["hash"]
		table, err := s.db.GetTable(hash)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		now := time.Now().UTC()

		if !table.LastCall.IsZero() && now.Sub(table.LastCall) < 5*time.Minute {
			msg := "Too many requests: please wait 5 minutes between calls"
			s.logger.Error(msg)
			http.Error(w, msg, http.StatusTooManyRequests)
			return
		}

		order := store.Order{
			Id:             uuid.New(),
			RestaurantName: table.RestaurantName,
			TableNumber:    table.Number,
			CreatedAt:      now,
			UpdatedAt:      now,
			Accepted:       false,
		}

		s.db.FindAndDeleteExistingCall(int8(order.TableNumber))
		err = s.db.CreateCall(order)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.logger.Info("order [", order.Id, "] created")

		err = s.producer.SendMessage(order)
		if err != nil {
			s.logger.Error(err)
			http.Error(w, "Failed to produce Kafka message", http.StatusInternalServerError)
			return
		}
		s.logger.Info("order [", order.Id, "] pushed to kafka")

		table.LastCall = now
		if err = s.db.UpdateTable(hash, table); err != nil {
			s.logger.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderJSON(w, order)
	}
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
