package middleware

import (
	"database/sql"
	"github.com/TechBuilder-360/Auth_Server/logger"
	"net/http"
	"strings"
)

// Middleware chains middleware
type Middleware struct {
	Next   http.Handler
	db     *sql.DB
}

// New starts a middleware build / chain
func New(defaultHandler http.Handler, db *sql.DB) *Middleware {

	return &Middleware{defaultHandler, db}
}

// Build returns a handler to all chained middleware
func (m *Middleware) Build() http.Handler {
	return m.Next
}

// UseClientValidation validates the call client
func (m *Middleware) UseClientValidation() *Middleware {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.URL.Path, "/api/") {
			m.Next.ServeHTTP(w, r)
			return
		}

		//response := utility.NewResponse()
		//client := model.ApiClient{}
		//validated := client.ValidateRequest(r, m.db)
		//
		//m.logger.Info("value for validate  %+v", validated)

		//if !validated {
		//	w.WriteHeader(http.StatusUnauthorized)
		//	json.NewEncoder(w).Encode(response.Error(utility.CLIENTERROR, utility.GetCodeMsg(utility.CLIENTERROR)))
		//	return
		//}

		m.Next.ServeHTTP(w, r)
	})

	logger.Info("UseClientValidation middleware registered successfully.")

	return &Middleware{nextHandler, m.db}
}

