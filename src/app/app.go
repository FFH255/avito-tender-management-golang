package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	bidrepository "tms/src/core/data/bid-repository"
	decisionrepository "tms/src/core/data/decision-repository"
	employeerepository "tms/src/core/data/employee-repository"
	organizationresponsiblerepository "tms/src/core/data/organization-responsible-repository"
	tenderrepository "tms/src/core/data/tender-repository"
	bidusecases "tms/src/core/services/use-cases/bid"
	usecases "tms/src/core/services/use-cases/tender"
	"tms/src/pkg/logger/sl"
	"tms/src/pkg/pg"
	httpserver "tms/src/transport/http-server"
	"tms/src/transport/http-server/handlers"
	bidhandlers "tms/src/transport/http-server/handlers/bid"
	tenderhandlers "tms/src/transport/http-server/handlers/tender"
)

func Run() {

	// Infra
	log := sl.SetupLogger()
	cfg := mustLoadConfig(*log)

	log.With("config", cfg).Info("application started")

	psqlClient, err := pg.New(
		log,
		pg.Config{
			URL:         cfg.Postgres.Conn,
			AutoMigrate: cfg.Postgres.AutoMigrate,
			Migrations:  cfg.Postgres.Migration,
		})
	defer psqlClient.Close()

	if err != nil {
		panic(err)
	}

	// Repositories
	tenderRepository := tenderrepository.New(*psqlClient)
	employeeRepository := employeerepository.New(*psqlClient)
	orgResponsibleRepository := organizationresponsiblerepository.New(*psqlClient)
	bidRepository := bidrepository.New(*psqlClient)
	decisionRepository := decisionrepository.New(*psqlClient)

	// UseCases
	getAllTendersUseCase := usecases.NewGetAllTendersUseCase(tenderRepository)
	createTenderUseCase := usecases.NewCreateTenderUseCase(
		orgResponsibleRepository,
		tenderRepository,
		employeeRepository,
	)
	getUserTendersUseCase := usecases.NewGetUserTendersUseCase(
		employeeRepository,
		orgResponsibleRepository,
		tenderRepository,
	)
	getTenderStatusUseCase := usecases.NewGetTenderStatusUseCase(
		employeeRepository,
		orgResponsibleRepository,
		tenderRepository,
	)
	changeTenderStatusUseCase := usecases.NewChangeTenderStatusUseCase(
		employeeRepository,
		orgResponsibleRepository,
		tenderRepository,
	)
	editTenderUseCase := usecases.NewEditTenderUseCase(
		employeeRepository,
		orgResponsibleRepository,
		tenderRepository,
	)
	rollbackTenderUseCase := usecases.NewRollBackTenderUseCase(
		employeeRepository,
		orgResponsibleRepository,
		tenderRepository,
	)
	createBidUseCase := bidusecases.NewCreateBidUseCase(
		employeeRepository,
		tenderRepository,
		bidRepository,
	)
	getUserBidsUseCase := bidusecases.NewGetUserBidsUseCase(
		employeeRepository,
		bidRepository,
	)
	getBidsOfTenderUseCase := bidusecases.NewGetBidsOfTenderUseCase(
		employeeRepository,
		tenderRepository,
		bidRepository,
		orgResponsibleRepository,
	)
	getBidStatusUseCase := bidusecases.NewGetBidStatusUseCase(
		employeeRepository,
		bidRepository,
	)
	changeBidStatusUseCase := bidusecases.NewChangeBidStatusUseCase(
		employeeRepository,
		bidRepository,
	)
	editBidUseCase := bidusecases.NewEditBidUseCase(
		employeeRepository,
		bidRepository,
	)
	submitDecisionUseCase := bidusecases.NewSubmitDecisionUseCase(
		employeeRepository,
		orgResponsibleRepository,
		bidRepository,
		tenderRepository,
		decisionRepository,
	)
	rollbackBidUseCase := bidusecases.NewRollbackBidUseCase(
		employeeRepository,
		bidRepository,
	)

	// Handlers
	pingHandler := handlers.NewPingHandler()
	getAllTendersHandler := tenderhandlers.NewGetAllTendersHandler(*log, getAllTendersUseCase)
	createTenderHandler := tenderhandlers.NewCreateTenderHandler(*log, createTenderUseCase)
	getMyTendersHandler := tenderhandlers.NewGetMyTendersHandlers(*log, getUserTendersUseCase)
	getTenderStatusHandler := tenderhandlers.NewGetTenderStatus(*log, getTenderStatusUseCase)
	changeTenderStatusHandler := tenderhandlers.NewChangeTenderStatusHandler(*log, changeTenderStatusUseCase)
	editTenderUseHandler := tenderhandlers.NewEditTenderHandler(*log, editTenderUseCase)
	rollbackTenderHandler := tenderhandlers.NewRollbackTenderHandler(*log, rollbackTenderUseCase)
	createBidHandler := bidhandlers.NewCreateBidHandler(*log, createBidUseCase)
	getUserBidsHandler := bidhandlers.NewGetUserBidsHandler(*log, getUserBidsUseCase)
	getBidsOfTenderHandler := bidhandlers.NewGetBidsOfTender(*log, getBidsOfTenderUseCase)
	getBidStatusHandler := bidhandlers.NewGetBidStatusHandler(*log, getBidStatusUseCase)
	changeBidStatusHandler := bidhandlers.NewChangeBidStatusHandler(*log, changeBidStatusUseCase)
	editBidHandler := bidhandlers.NewEditBidHandler(*log, editBidUseCase)
	submitDecisionHandler := bidhandlers.NewSubmitDecisionHandler(*log, submitDecisionUseCase)
	rollbackBidHandler := bidhandlers.NewRollBackHandler(*log, rollbackBidUseCase)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	h := httpserver.Handlers{
		Ping:               pingHandler,
		GetAllTenders:      getAllTendersHandler,
		CreateTenders:      createTenderHandler,
		GetMyTenders:       getMyTendersHandler,
		GetTenderStatus:    getTenderStatusHandler,
		ChangeTenderStatus: changeTenderStatusHandler,
		EditTender:         editTenderUseHandler,
		RollbackTender:     rollbackTenderHandler,
		CreateBid:          createBidHandler,
		GetUserBid:         getUserBidsHandler,
		GetBidsOfTender:    getBidsOfTenderHandler,
		GetBidStatus:       getBidStatusHandler,
		ChangeBidStatus:    changeBidStatusHandler,
		EditBid:            editBidHandler,
		SubmitDecision:     submitDecisionHandler,
		RollbackBid:        rollbackBidHandler,
	}

	srv := httpserver.New(
		h,
		*log,
		httpserver.Config{
			Address:      cfg.HTTPServer.Address,
			ReadTimeout:  cfg.HTTPServer.ReadTimeout,
			WriteTimeout: cfg.HTTPServer.WriteTimeout,
			IdleTimeout:  cfg.HTTPServer.IdleTimeout,
		},
	)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error(fmt.Sprintf("failed to maintain server because: %v", err))
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))
		return
	}

	log.Info("application stopped")
}
