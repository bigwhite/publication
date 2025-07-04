package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/your_org/shortlink/internal/api/http/handler"
	"github.com/your_org/shortlink/internal/shortener"
)

// Server 是我们的HTTP服务器结构
type Server struct {
	httpServer *http.Server
	service    *shortener.Service // Handler会用到
	logger     *log.Logger
}

// New 创建一个新的 Server 实例
func New(port string, svc *shortener.Service) *Server {
	logger := log.New(os.Stdout, "[HTTP Server] ", log.LstdFlags|log.Lshortfile)
	linkAPIHandler := handler.NewLinkAPI(svc, logger) // 创建 handler

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/links", linkAPIHandler.CreateLink)
	mux.HandleFunc("GET /", linkAPIHandler.RedirectLink) // RedirectLink 内部会解析短码
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	})

	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + port,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		service: svc,
		logger:  logger,
	}
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	s.logger.Printf("Server listening on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown 优雅地关闭HTTP服务器
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}
