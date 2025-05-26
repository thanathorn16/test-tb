package main

import (
	"be/config"
	"be/services/login"
	"be/services/profile"
	"be/services/register"
	"be/stores"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	appConf, err := config.Init()
	if err != nil {
		panic(fmt.Sprintf("error loading config: %v", err))
	}

	cfg := zap.NewProductionConfig()
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	e := initEcho()

	userDB, err := ConnectDB(appConf.DataBase)
	if err != nil {
		panic(fmt.Sprintf("error connecting to database: %v", err))
	}
	userRepo := stores.NewUserDB(userDB)

	parsedAccessTokenPrivatekey, err := parseRSAprivateKey([]byte(appConf.AccessTokenPrivatekey))
	if err != nil {
		panic(err)
	}

	loginService := login.NewLoginService(userRepo, parsedAccessTokenPrivatekey)
	loginHanlder := login.NewHandler(loginService, logger)
	e.POST("/login", loginHanlder.Login)

	registerService := register.NewRegisterService(userRepo)
	registerHandler := register.NewHandler(registerService, logger)
	e.POST("/register", registerHandler.Register)

	profileService := profile.NewProfileService(userRepo, parsedAccessTokenPrivatekey)
	profileHandler := profile.NewHandler(profileService, logger)
	e.GET("/profile", profileHandler.Register)

	defer gracefulShutdown(e, appConf.Port)

}

func gracefulShutdown(e *echo.Echo, port string) {
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func initEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e
}

func ConnectDB(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	fmt.Printf("connecting to %s:[secret]@tcp(%s:%s)/%s?parseTime=True\n", cfg.Username, cfg.Host, cfg.Port, cfg.DBName)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot connect db [%s:%s/%s] error: %w", cfg.Host, cfg.Port, cfg.DBName, err)
	}

	fmt.Printf("connected to db: %s, host: %s\n", cfg.DBName, cfg.Host)

	return db, nil
}

func parseRSAprivateKey(key []byte) (*rsa.PrivateKey, error) {
	keyBlock, _ := pem.Decode(key)
	if keyBlock == nil {
		return &rsa.PrivateKey{}, fmt.Errorf("failed to decode private key %s", key)
	}

	var parsedRSAKey *rsa.PrivateKey

	parsedRSAKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		parsedKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
		if err != nil {
			return &rsa.PrivateKey{}, fmt.Errorf("failed to parse public key due to %w", err)
		}
		parsedRSAKey = parsedKey.(*rsa.PrivateKey)
	}

	return parsedRSAKey, nil
}
