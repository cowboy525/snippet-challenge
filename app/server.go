package app

import (
	"fmt"
	"net"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/ernie-mlg/ErniePJT-main-api-go/config"
	"github.com/ernie-mlg/ErniePJT-main-api-go/mlog"
	"github.com/ernie-mlg/ErniePJT-main-api-go/model"
	"github.com/ernie-mlg/ErniePJT-main-api-go/services/filestore"
	"github.com/ernie-mlg/ErniePJT-main-api-go/store"
	"github.com/ernie-mlg/ErniePJT-main-api-go/store/sqlstore"
	"github.com/ernie-mlg/ErniePJT-main-api-go/utils"
)

var MaxNotificationsPerChannelDefault int64 = 1000000

type Server struct {
	Store store.Store

	RootRouter *mux.Router
	Router     *mux.Router

	Server     *http.Server
	ListenAddr *net.TCPAddr
	Log        *mlog.Logger

	configStore config.Store
}

func NewServer(options ...Option) (*Server, error) {
	rootRouter := mux.NewRouter()

	s := &Server{
		RootRouter: rootRouter,
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, "failed to apply option")
		}
	}

	if s.configStore == nil {
		configStore, err := config.NewFileStore("config.json", true)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load config")
		}

		s.configStore = configStore
	}

	if s.Log == nil {
		s.Log = mlog.NewLogger(utils.MloggerConfigFromLoggerConfig(&s.Config().LogSettings, utils.GetLogFileLocation))
	}

	// Redirect default golang logger to this logger
	mlog.RedirectStdLog(s.Log)

	// Use this app logger as the global logger (eventually remove all instances of global logging)
	mlog.InitGlobalLogger(s.Log)

	if s.Store == nil {
		s.Store = sqlstore.NewSQLSupplier(s.Config().SQLSettings, false)
	}

	subpath := "/"
	s.Router = s.RootRouter.PathPrefix(subpath).Subrouter()

	// If configured with a subpath, redirect 404s at the root back into the subpath.
	if subpath != "/" {
		s.RootRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = path.Join(subpath, r.URL.Path)
			http.Redirect(w, r, r.URL.String(), http.StatusFound)
		})
	}
	s.Router.NotFoundHandler = http.HandlerFunc(s.FakeApp().Handle404)

	s.ReloadConfig()

	return s, nil
}

func (s *Server) Shutdown() error {
	mlog.Info("Stopping Server...")

	s.configStore.Close()

	if s.Store != nil {
		s.Store.Close()
	}

	mlog.Info("Server stopped")
	return nil
}

func (s *Server) Start() error {
	mlog.Info("Starting Server...")

	headersOk := handlers.AllowedHeaders([]string{"x-api-version", "authorization", "content-type", "client-id", "client-secretkey"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETe"})

	var handler http.Handler = handlers.CORS(headersOk, originsOk, methodsOk)(s.RootRouter)

	// Creating a logger for logging errors from http.Server at error level
	errStdLog, err := s.Log.StdLogAt(mlog.LevelError, mlog.String("source", "httpserver"))
	if err != nil {
		return err
	}

	s.Server = &http.Server{
		Handler:  handler,
		ErrorLog: errStdLog,
	}

	addr := *s.Config().ServiceSettings.ListenAddress
	if addr == "" {
		if *s.Config().ServiceSettings.ConnectionSecurity == model.CONN_SECURITY_TLS {
			addr = ":https"
		} else {
			addr = ":http"
		}
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		errors.Wrapf(err, "api.server.start_server.starting.critical")
		return err
	}
	s.ListenAddr = listener.Addr().(*net.TCPAddr)

	logListeningPort := fmt.Sprintf("Server is listening on %v", listener.Addr().String())
	mlog.Info(logListeningPort, mlog.String("address", listener.Addr().String()))

	go func() {
		if err := s.Server.Serve(listener); err != nil && err != http.ErrServerClosed {
			mlog.Critical("Error starting server", mlog.Err(err))
			time.Sleep(time.Second)
		}
	}()

	return nil
}

// A temporary bridge to deal with cases where the code is so tighly coupled that
// this is easier as a temporary solution
func (s *Server) FakeApp() *App {
	a := New(
		ServerConnector(s),
	)
	return a
}

// Global app options that should be applied to apps created by this server
func (s *Server) AppOptions() []AppOption {
	return []AppOption{
		ServerConnector(s),
	}
}

func (s *Server) FileBackend() (filestore.FileBackend, *model.AppError) {
	return filestore.NewFileBackend(s.Config())
}
