package http_server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
	"tms/src/transport/http-server/middleware/logger"
)

type Config struct {
	Address      string
	ReadTimeout  uint
	WriteTimeout uint
	IdleTimeout  uint
}

type Handlers struct {
	Ping http.HandlerFunc
	// Tender handlers
	GetAllTenders      http.HandlerFunc
	CreateTenders      http.HandlerFunc
	GetMyTenders       http.HandlerFunc
	GetTenderStatus    http.HandlerFunc
	ChangeTenderStatus http.HandlerFunc
	EditTender         http.HandlerFunc
	RollbackTender     http.HandlerFunc
	// Bid handlers
	CreateBid       http.HandlerFunc
	GetUserBid      http.HandlerFunc
	GetBidsOfTender http.HandlerFunc
	GetBidStatus    http.HandlerFunc
	ChangeBidStatus http.HandlerFunc
	EditBid         http.HandlerFunc
	SubmitDecision  http.HandlerFunc
	RollbackBid     http.HandlerFunc
}

func New(handlers Handlers, log slog.Logger, cfg Config) *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(logger.New(&log))

	router.Route("/api", func(r chi.Router) {
		r.Get("/ping", handlers.Ping)
		// Tender endpoints
		r.Get("/tenders", handlers.GetAllTenders)
		r.Post("/tenders/new", handlers.CreateTenders)
		r.Get("/tenders/my", handlers.GetMyTenders)
		r.Get("/tenders/{tenderId}/status", handlers.GetTenderStatus)
		r.Put("/tenders/{tenderId}/status", handlers.ChangeTenderStatus)
		r.Patch("/tenders/{tenderId}/edit", handlers.EditTender)
		r.Put("/tenders/{tenderId}/rollback/{version}", handlers.RollbackTender)
		// Bid endpoints
		r.Post("/bids/new", handlers.CreateBid)
		r.Get("/bids/my", handlers.GetMyTenders)
		r.Get("/bids/{tenderId}/list", handlers.GetBidsOfTender)
		r.Get("/bids/{bidId}/status", handlers.GetBidStatus)
		r.Put("/bids/{bidId}/status", handlers.ChangeBidStatus)
		r.Patch("/bids/{bidId}/edit", handlers.EditBid)
		r.Put("/bids/{bidId}/submit_decision", handlers.SubmitDecision)
		r.Put("/bids/{bidId}/rollback/{version}", handlers.RollbackBid)
	})

	return &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}
}
