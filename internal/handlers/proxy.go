package handlers

import (
    "crypto/tls"
    "net/http"
    "net/http/httputil"
    "net/url"
    "sync"
    
    "gotiwul/pkg/health"
)

type ProxyHandler struct {
    checker     *health.Checker
    mutex       sync.Mutex
    current     int
    backends    []string
}

func NewProxyHandler(backends []string) *ProxyHandler {
    h := &ProxyHandler{
        backends: backends,
        current:  0,
    }
    h.checker = health.NewChecker(backends)
    return h
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    backend := h.getNextBackend()
    if backend == "" {
        http.Error(w, "No available backends", http.StatusServiceUnavailable)
        return
    }

    backendURL, err := url.Parse(backend)
    if err != nil {
        http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
        return
    }

    proxy := &httputil.ReverseProxy{
        Director: func(req *http.Request) {
            req.URL.Scheme = backendURL.Scheme
            req.URL.Host = backendURL.Host
            req.Host = backendURL.Host
            
            req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
            req.Header.Set("X-Forwarded-Proto", "https")
            req.Header.Set("X-Forwarded-For", req.RemoteAddr)
        },
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true,
            },
        },
    }

    proxy.ServeHTTP(w, r)
}

func (h *ProxyHandler) getNextBackend() string {
    h.mutex.Lock()
    defer h.mutex.Unlock()

    attempts := 0
    for attempts < len(h.backends) {
        h.current = (h.current + 1) % len(h.backends)
        backend := h.backends[h.current]
        if h.checker.IsHealthy(backend) {
            return backend
        }
        attempts++
    }
    return ""
}