package handlers

import (
	"OrderSystem/pkg/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type OrdersProxy struct {
	baseURL string
	client  *http.Client
	log     *logger.Logger
}

func NewOrdersProxy(baseURL string, log *logger.Logger) *OrdersProxy {
	return &OrdersProxy{
		baseURL: baseURL,
		client:  &http.Client{},
		log:     log,
	}
}

// ProxyTo — универсальный прокси для всех роутов /orders/*
func (p *OrdersProxy) ProxyTo(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := p.baseURL + path
		if q := c.Request.URL.RawQuery; q != "" {
			url += "?" + q
		}

		req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
		if err != nil {
			p.log.Error("orders proxy: failed to create request", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "proxy error"})
			return
		}

		req.Header = c.Request.Header.Clone()

		if userID := c.GetString("user_id"); userID != "" {
			req.Header.Set("X-User-ID", userID)
		}
		if role := c.GetString("role"); role != "" {
			req.Header.Set("X-User-Role", role)
		}

		resp, err := p.client.Do(req)
		if err != nil {
			p.log.Error("orders proxy: service unavailable", "error", err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "orders service unavailable"})
			return
		}
		defer resp.Body.Close()

		for k, values := range resp.Header {
			for _, v := range values {
				c.Writer.Header().Add(k, v)
			}
		}

		c.Status(resp.StatusCode)
		if _, err := io.Copy(c.Writer, resp.Body); err != nil {
			p.log.Warn("orders proxy: failed to copy response body", "error", err)
		}
	}
}
