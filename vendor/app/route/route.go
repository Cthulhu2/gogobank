package route

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"app/handler"

	"github.com/gorilla/mux"
)

// Load returns the routes and middleware
func Load() http.Handler {
	return routes()
}

// LoadHTTPS returns the HTTP routes and middleware
func LoadHTTPS() http.Handler {
	return routes()
}

// LoadHTTP returns the HTTPS routes and middleware
func LoadHTTP() http.Handler {
	return routes()

	// Uncomment this and comment out the line above to always redirect to HTTPS
	//return http.HandlerFunc(redirectToHTTPS)
}

// Optional method to make it easy to redirect from HTTP to HTTPS
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://"+req.Host, http.StatusMovedPermanently)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(
			time.Now().Format("2006-01-02 03:04:05 PM"),
			r.RemoteAddr, r.Method, r.URL)
		// Call the next handler, which can be another middleware in the chain,
		// or the final handler.
		next.ServeHTTP(w, r)
	})
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	r.NotFoundHandler = loggingMiddleware(
		http.HandlerFunc(handler.Error404))
	r.MethodNotAllowedHandler = loggingMiddleware(
		http.HandlerFunc(handler.Error405))

	// Home page
	r.HandleFunc("/", handler.IndexGET)
	r.HandleFunc("/index", handler.IndexGET)

	// Account
	sa := r.PathPrefix("/v1/account").Subrouter()
	sa.HandleFunc("/{id}", handler.AccountGET).
		Methods("GET")
	sa.HandleFunc("/", handler.AccountPOST).
		Methods("POST")
	sa.HandleFunc("/trans/", handler.TransferPOST).
		Methods("POST")

	// Enable Pprof
	sd := r.PathPrefix("/debug/pprof/").Subrouter()
	sd.HandleFunc("/cmdline", pprof.Cmdline)
	sd.HandleFunc("/profile", pprof.Profile)
	sd.HandleFunc("/symbol", pprof.Symbol)
	sd.HandleFunc("/", pprof.Index)

	// SwaggerUI
	fs := http.FileServer(http.Dir("./api/swaggerui/"))
	r.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", fs))

	return r
}
