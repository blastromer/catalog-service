package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"micro-services/catalog/internal/config"
	"micro-services/catalog/internal/database"

	"github.com/gin-gonic/gin"
)

// Dependencies are shared by all catalog microservices.
type Dependencies struct {
	Config *config.Config
	Store  *database.Store
}

// NewRouter builds a Gin engine with health checks and recovery middleware.
func NewRouter(deps Dependencies) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), requestLogger())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": deps.Config.ServiceName,
		})
	})

	r.GET("/readyz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := deps.Store.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "not_ready",
				"service": deps.Config.ServiceName,
				"error":   "database unavailable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ready",
			"service": deps.Config.ServiceName,
		})
	})

	return r
}

// Run starts the HTTP server and shuts down gracefully on SIGINT/SIGTERM.
func Run(deps Dependencies, handler http.Handler) error {
	srv := &http.Server{
		Addr:         deps.Config.HTTPAddr,
		Handler:      handler,
		ReadTimeout:  deps.Config.ReadTimeout,
		WriteTimeout: deps.Config.WriteTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("%s listening on %s", deps.Config.ServiceName, deps.Config.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-stop:
		log.Printf("received %s, shutting down", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), deps.Config.ShutdownGrace)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return deps.Store.Close(ctx)
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Printf("%s %s %d %s", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start))
	}
}
