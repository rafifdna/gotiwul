package server

import (
    "context"
    "crypto/tls"
    "fmt"
    "log"
    "net/http"
    "path/filepath"
    
    "golang.org/x/crypto/acme/autocert"
    "gotiwul/config"
    "gotiwul/internal/handlers"
    "gotiwul/internal/middleware"
)

type Server struct {
    cfg          *config.Config
    httpServer   *http.Server
    httpsServer  *http.Server
}

func New(cfg *config.Config) *Server {
    return &Server{cfg: cfg}
}

func (s *Server) StartWithAutoTLS() error {
    certManager := autocert.Manager{
        Prompt:     autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist(s.cfg.Server.Domain),
        Cache:      autocert.DirCache(filepath.Join("certs", "acme")),
        Email:      s.cfg.Server.Email,
    }

    proxyHandler := handlers.NewProxyHandler(s.cfg.Backends)
    handler := middleware.Logging(proxyHandler)

    tlsConfig := &tls.Config{
        GetCertificate:           certManager.GetCertificate,
        MinVersion:               tls.VersionTLS12,
        CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
    }

    s.httpsServer = &http.Server{
        Addr:      fmt.Sprintf(":%d", s.cfg.Server.Port),
        Handler:   handler,
        TLSConfig: tlsConfig,
    }

    s.httpServer = &http.Server{
        Addr:    ":80",
        Handler: certManager.HTTPHandler(nil),
    }

    go func() {
        log.Printf("Starting HTTP server on port 80")
        if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
            log.Printf("HTTP server error: %v", err)
        }
    }()

    log.Printf("Starting HTTPS server on port %d", s.cfg.Server.Port)
    return s.httpsServer.ListenAndServeTLS("", "")
}

func (s *Server) Shutdown(ctx context.Context) error {
    if err := s.httpServer.Shutdown(ctx); err != nil {
        log.Printf("HTTP server shutdown error: %v", err)
    }
    return s.httpsServer.Shutdown(ctx)
}