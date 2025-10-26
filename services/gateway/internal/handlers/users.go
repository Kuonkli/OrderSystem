package handlers

import (
	"OrderSystem/pkg/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type UsersProxy struct {
	baseURL string
	client  *http.Client
	log     *logger.Logger
}

func NewUsersProxy(baseURL string, log *logger.Logger) *UsersProxy {
	return &UsersProxy{
		baseURL: baseURL,
		client:  &http.Client{},
		log:     log,
	}
}

func (p *UsersProxy) ProxyTo(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := p.baseURL + path
		if q := c.Request.URL.RawQuery; q != "" {
			url += "?" + q
		}

		req, err := http.NewRequest(c.Request.Method, url, c.Request.Body)
		if err != nil {
			p.log.Error("proxy request error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "proxy error"})
			return
		}

		req.Header = c.Request.Header.Clone()
		req.Header.Set("X-User-ID", c.GetString("user_id"))
		req.Header.Set("X-User-Role", c.GetString("role"))

		resp, err := p.client.Do(req)
		if err != nil {
			p.log.Error("proxy request failed: ", err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "service unavailable"})
			return
		}
		defer resp.Body.Close()

		access := resp.Header.Get("Access-Token")
		refresh := resp.Header.Get("Refresh-Token")

		if access != "" {
			c.SetCookie("access_token", access, 900, "/", "", false, true)
			p.log.Info("Set access_token from header")
			resp.Header.Del("Access-Token") // Убираем, чтобы не отдавать клиенту
		}

		if refresh != "" {
			c.SetCookie("refresh_token", refresh, 86400, "/", "", false, true)
			p.log.Info("Set refresh_token from header")
			resp.Header.Del("Refresh-Token")
		}

		for k, v := range resp.Header {
			for _, vv := range v {
				c.Writer.Header().Add(k, vv)
			}
		}

		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}
