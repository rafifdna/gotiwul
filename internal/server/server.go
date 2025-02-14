package server

import (
    "crypto/tls"
    "fmt"
    "net/http"
    
    "gotiwul/config"
    "gotiwul/internal/handlers"
    "gotiwul/internal/middleware"
)

type Server struct {
    cfg *config.Config
}

func New(cfg *config.Config) *Server {
    return &Server{cfg: cfg}
}

func (s *Server) Start() error {
    proxyHandler := handlers.NewProxyHandler(s.cfg.Backends)
    handler := middleware.Logging(proxyHandler)

    server := &http.Server{
        Addr:    fmt.Sprintf(":%d", s.cfg.Server.Port),
        Handler: handler,
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS12,
            CurvePreferences: []tls.CurveID{
                tls.CurveP256,
                tls.X25519,
            },
            PreferServerCipherSuites: true,
            CipherSuites: []uint16{
                tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
                tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
                tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
                tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            },
        },
    }

    return server.ListenAndServeTLS(s.cfg.Server.CertFile, s.cfg.Server.KeyFile)
}