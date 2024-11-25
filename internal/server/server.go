package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"semaright/internal/database"
	"semaright/internal/entities"
	"semaright/internal/pricing"
	"time"
)

type server struct {
	port             int
	connectionString string
	db               *database.DB
}

func New(port int, cs string) (*server, error) {
	db, err := database.New(cs)
	if err != nil {
		return nil, err
	}

	err = db.Connect()
	if err != nil {
		return nil, err
	}

	return &server{
		port:             port,
		connectionString: cs,
		db:               db,
	}, nil
}

func (s *server) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {

		rq := entities.Transaction{}
		err := json.NewDecoder(r.Body).Decode(&rq)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Calculate spend for request
		spend, err := pricing.CalculateSpend(rq.Model, float64(rq.InputTokens), float64(rq.OutputTokens))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Model pricing metadata not found"))
			return
		}

		rq.Spend = spend
		rq.Date = time.Now()

		// Save to usecase-transactions collection using a Mongo transaction & update total spend for the usecase
		err = s.db.SaveTransaction(r.Context(), &rq)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to save transaction " + err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	mux.HandleFunc("/allow/{usecaseId}", func(w http.ResponseWriter, r *http.Request) {
		// Parse usecaseId from path
		uc := r.PathValue("usecaseId")
		budget := pricing.GetBudget(uc)

		currentSpend, err := s.db.GetCurrentSpend(r.Context(), uc)
		if err != nil {
			if errors.Is(err, database.ErrNoEntryFound) {
				w.WriteHeader(http.StatusOK)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to get current spend " + err.Error()))
			return
		}

		if currentSpend/1000000 > budget {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux)
	if err != nil {
		return err
	}

	return nil
}
