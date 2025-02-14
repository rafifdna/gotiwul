package health

import (
    "crypto/tls"
    "net/http"
    "sync"
    "time"
)

type Checker struct {
    backends    []string
    healthCheck map[string]bool
    mutex       sync.Mutex
    client      *http.Client
}

func NewChecker(backends []string) *Checker {
    c := &Checker{
        backends:    backends,
        healthCheck: make(map[string]bool),
        client: &http.Client{
            Timeout: 5 * time.Second,
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{
                    InsecureSkipVerify: true,
                },
            },
        },
    }
    
    for _, backend := range backends {
        c.healthCheck[backend] = true
    }
    
    go c.start()
    return c
}

func (c *Checker) IsHealthy(backend string) bool {
    c.mutex.Lock()
    defer c.mutex.Unlock()
    return c.healthCheck[backend]
}

func (c *Checker) start() {
    ticker := time.NewTicker(10 * time.Second)
    for range ticker.C {
        c.checkAll()
    }
}

func (c *Checker) checkAll() {
    for _, backend := range c.backends {
        go c.checkBackend(backend)
    }
}

func (c *Checker) checkBackend(backend string) {
    resp, err := c.client.Get(backend + "/health")
    
    c.mutex.Lock()
    defer c.mutex.Unlock()
    
    if err != nil || resp.StatusCode != http.StatusOK {
        c.healthCheck[backend] = false
    } else {
        c.healthCheck[backend] = true
    }
}