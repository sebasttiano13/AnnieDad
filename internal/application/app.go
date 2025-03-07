package application

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano13/AnnieDad/internal/config"
	"github.com/sebasttiano13/AnnieDad/internal/server"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Run starts application
func Run(ctx context.Context, cfgPath string) error {
	cfg, err := config.NewDadConfig(cfgPath)
	if err != nil {
		return err
	}

	_, err = logger.NewLogger(cfg.LoggerCfg.File, cfg.LoggerCfg.Level)
	if err != nil {
		return err
	}

	var db *sqlx.DB
	db, err = sqlx.Connect("pgx", cfg.DBCfg.GetDSN())
	if err != nil {
		logger.Errorf("database openning failed: %e", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	wg := &sync.WaitGroup{}

	settings := &server.GRPSServerSettings{
		SecretKey:     cfg.GRPSServerCfg.Secret,
		TokenDuration: time.Duration(cfg.GRPSServerCfg.TokenDuration) * time.Second,
		CertFile:      cfg.Cert.Cert,
		CertKey:       cfg.Cert.Key,
	}

	grpcSrv := server.NewGRPSServer(settings, db)

	wg.Add(1)
	go grpcSrv.Start(cfg.GetGRPSAddress())
	go grpcSrv.HandleShutdown(ctx, wg)
	wg.Wait()

	return nil
}
