package main

import (
	"database/sql"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"

	"github.com/kaatinga/env_loader"

	launcher "github.com/kaatinga/QuickHTTPServerLauncher"
)

// Модель данных для параметров окружения
type EnvironmentSettings struct {
	Port string `env:"PORT" validate:"required"`
}

const incorrectInputPhrase = "Incorrect input phrase.\nEnter an URL like http://<IP address>/world"

type server struct {
	launcher.Config
}

func main() {

	var err error

	// New web service
	var s = server{launcher.NewConfig()}

	var myEnvs EnvironmentSettings
	err = env_loader.LoadUsingReflect(&myEnvs)
	if err != nil {
		s.Config.Logger.SubMsg.Err(err).Msg("Environment variables have not been read")
		os.Exit(1)
	}
	s.Config.Logger.SubMsg.Info().Msg("Environment variables are set properly!")

	s.Config.SetDomain("kaatinga.ru")
	s.Config.SetEmail("info@kaatinga.ru")
	s.Config.SetLaunchMode("dev")
	s.Config.SetPort(myEnvs.Port)

	err = s.Config.Launch(s.SetUpHandlers)
	if err != nil {
		s.Config.Logger.SubMsg.Err(err).Msg("The s stopped")
	}
}

func (s server) SetUpHandlers(r *httprouter.Router, _ *sql.DB) {
	r.GET("/:phrase", s.HelloServer)
	r.GET("/", s.HelloServer)
}

func (s server) HelloServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	s.Config.Logger.Title.Info().Str("IP", r.RemoteAddr).Str("Method", r.Method).Str("URL", r.URL.String()).Msg("== A new request is received:")

	bytes := getPhraseBytes(ps)
	_, err := w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getPhraseBytes(ps httprouter.Params) []byte {
	if ps.ByName("phrase") == "" {
		return []byte(incorrectInputPhrase)
	}
	return []byte("Hello, " + ps.ByName("phrase"))
}
