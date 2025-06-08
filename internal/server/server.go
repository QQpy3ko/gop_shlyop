package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gop_shlyop/internal/config"
	"gop_shlyop/internal/db"
	"gop_shlyop/internal/llm"
	"gop_shlyop/internal/metrics"
	"gop_shlyop/internal/repository"
	"gop_shlyop/internal/types"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

type Server struct {
	cfg             *config.Config
	log             zerolog.Logger
	dbpool          *pgxpool.Pool
	reviewsRepo     *repository.ReviewsRepository
	sentimentClient *llm.Client
	metrics         *metrics.Metrics
}

func New(ctx context.Context, cfg *config.Config, log zerolog.Logger) (*Server, error) {
	log.Info().Msg("Connecting to PostgreSQL...")
	dbpool, err := db.NewPostgresDB(ctx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	log.Info().Msg("Successfully connected to PostgreSQL")

	reviewsRepo := repository.NewReviewsRepository(dbpool)
	sentimentClient := llm.NewClient(cfg.Ollama)
	metrics := metrics.NewMetrics(prometheus.DefaultRegisterer)

	return &Server{
		cfg:             cfg,
		log:             log,
		dbpool:          dbpool,
		reviewsRepo:     reviewsRepo,
		sentimentClient: sentimentClient,
		metrics:         metrics,
	}, nil
}

// run server
func (s *Server) Run() error {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger) // Chi's own logger is fine for now
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(s.metricsMiddleware)

	// Metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	// API routes
	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/reviews", s.addReview)
		r.Get("/reviews", s.getReviews)
		r.Get("/items/{itemId}/rating", s.getItemRating)
	})

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.cfg.HTTPServer.Host, s.cfg.HTTPServer.Port),
		Handler: router,
	}

	s.log.Info().
		Str("host", s.cfg.HTTPServer.Host).
		Str("port", s.cfg.HTTPServer.Port).
		Msg("server is starting")

	// run web-server in separate goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // wait for signal

	s.log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		s.log.Error().Err(err).Msg("Server forced to shutdown")
	}

	s.dbpool.Close()

	s.log.Info().Msg("Server stopped gracefully")
	return nil
}

// --- Handlers ---

func (s *Server) addReview(w http.ResponseWriter, r *http.Request) {
	var req types.AddReviewRequest
	if err := s.decodeJSON(w, r, &req); err != nil {
		s.errorResponse(w, r, http.StatusBadRequest, err)
		return
	}

	// TODO: Add validation for the request fields (e.g., UserID, ItemID, Text not empty)

	sentiment, err := s.sentimentClient.AnalyzeSentiment(r.Context(), req.Text)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to analyze sentiment")
		// Save the review with a "neutral" sentiment or fail?
		// For now, fail the request.
		s.errorResponse(w, r, http.StatusInternalServerError, errors.New("failed to analyze review sentiment"))
		return
	}

	review, err := s.reviewsRepo.CreateReview(r.Context(), req, sentiment)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to create review")
		s.errorResponse(w, r, http.StatusInternalServerError, errors.New("failed to save review"))
		return
	}

	// Increment the counter for created reviews
	s.metrics.ReviewsCreatedTotal.With(prometheus.Labels{"sentiment": sentiment}).Inc()

	s.writeJSON(w, http.StatusCreated, review)
}

func (s *Server) getReviews(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	itemIDStr := r.URL.Query().Get("item_id")

	if userIDStr != "" && itemIDStr != "" {
		s.errorResponse(w, r, http.StatusBadRequest, errors.New("only one filter (user_id or item_id) can be applied at a time"))
		return
	}

	var userID, itemID int64
	var err error

	if userIDStr != "" {
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil || userID <= 0 {
			s.errorResponse(w, r, http.StatusBadRequest, errors.New("invalid user_id parameter"))
			return
		}
	}

	if itemIDStr != "" {
		itemID, err = strconv.ParseInt(itemIDStr, 10, 64)
		if err != nil || itemID <= 0 {
			s.errorResponse(w, r, http.StatusBadRequest, errors.New("invalid item_id parameter"))
			return
		}
	}

	reviews, err := s.reviewsRepo.GetReviews(r.Context(), userID, itemID)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to get reviews")
		s.errorResponse(w, r, http.StatusInternalServerError, errors.New("failed to retrieve reviews"))
		return
	}

	s.writeJSON(w, http.StatusOK, reviews)
}

func (s *Server) getItemRating(w http.ResponseWriter, r *http.Request) {
	itemIDStr := chi.URLParam(r, "itemId")
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	if err != nil || itemID <= 0 {
		s.errorResponse(w, r, http.StatusBadRequest, errors.New("invalid item_id parameter"))
		return
	}

	rating, err := s.reviewsRepo.GetItemRating(r.Context(), itemID)
	if err != nil {
		s.log.Error().Err(err).Int64("item_id", itemID).Msg("failed to get item rating")
		s.errorResponse(w, r, http.StatusInternalServerError, errors.New("failed to retrieve item rating"))
		return
	}

	s.writeJSON(w, http.StatusOK, rating)
}

// --- Helpers ---

// middleware that records HTTP request metrics.
func (s *Server) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		duration := time.Since(start).Seconds()

		// Record duration
		s.metrics.HttpRequestDuration.With(prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Observe(duration)

		// Record total requests
		s.metrics.HttpRequestsTotal.With(prometheus.Labels{
			"method": r.Method,
			"path":   r.URL.Path,
			"code":   strconv.Itoa(ww.Status()),
		}).Inc()
	})
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.log.Error().Err(err).Msg("failed to write JSON response")
	}
}

func (s *Server) errorResponse(w http.ResponseWriter, r *http.Request, status int, message error) {
	errPayload := map[string]string{"error": message.Error()}
	s.writeJSON(w, status, errPayload)
}

func (s *Server) decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		// ... error handling ...
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON object")
	}

	return nil
}
