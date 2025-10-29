package application

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sebasttiano13/AnnieDad/internal/config"
	"github.com/sebasttiano13/AnnieDad/internal/repository"
	"github.com/sebasttiano13/AnnieDad/internal/server"
	"github.com/sebasttiano13/AnnieDad/pkg/clients"
	"github.com/sebasttiano13/AnnieDad/pkg/logger"
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

	db, err := sqlx.Connect("pgx", cfg.DBCfg.GetDSN())
	if err != nil {
		logger.Errorf("database openning failed: %e", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	wg := &sync.WaitGroup{}

	repo := repository.NewDBStorage(db)
	awsConfig := aws.Config{
		Region:       cfg.S3Cfg.Region,
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.S3Cfg.AccessKey, cfg.S3Cfg.SecretKey, ""),
		BaseEndpoint: aws.String(cfg.S3Cfg.Endpoint),
	}

	s3client := clients.NewS3Client(awsConfig, time.Duration(cfg.S3Cfg.ExpiresURLIn)*time.Minute)
	settings := &server.GRPSServerSettings{
		AccessSecretKey:      cfg.GRPSServerCfg.AccessSecret,
		RefreshSecretKey:     cfg.GRPSServerCfg.RefreshSecret,
		AccessTokenDuration:  time.Duration(cfg.GRPSServerCfg.AccessTokenDuration) * time.Minute,
		RefreshTokenDuration: time.Duration(cfg.GRPSServerCfg.RefreshTokenDuration) * time.Minute,
		CertFile:             cfg.Cert.Cert,
		CertKey:              cfg.Cert.Key,
	}

	grpcSrv := server.NewGRPSServer(settings, repo, repo, repo, s3client)

	wg.Add(1)
	go grpcSrv.Start(cfg.GetGRPSAddress())
	go grpcSrv.HandleShutdown(ctx, wg)
	wg.Wait()

	return nil
}
