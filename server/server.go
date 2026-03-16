package server

import (
	"net/http"
	"os"
	httphandlers "test_project/httpHandlers"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *httphandlers.HTTPHandlers
}

func NewHTTPServer(httpHanlders *httphandlers.HTTPHandlers) *HTTPServer {
	return &HTTPServer{httpHandlers: httpHanlders}
}
func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()

	router.Path("/api/v1/wallet").Methods("POST").HandlerFunc(s.httpHandlers.UpdateWallet)
	router.Path("/api/v1/wallets/{WALLET_UUID}").Methods("GET").HandlerFunc(s.httpHandlers.GetBalance)
	port := os.Getenv("PORT")
	if port == "" {
		port = "9091"
	}
	addr := ":" + port
	return http.ListenAndServe(addr, LoggingMiddleware(router))
}
