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
	Port       string `env:"PORT" validate:"required"`
}

const incorrectInputPhrase = "Incorrect input phrase.\nEnter an URL like http://<IP address>/world"

func main() {

	var err error

	// New web service
	config := launcher.NewConfig()

	var myEnvs EnvironmentSettings
	err = env_loader.LoadUsingReflect(&myEnvs)
	if err != nil {
		config.Logger.SubMsg.Err(err).Msg("Environment variables have not been read")
		os.Exit(1)
	}
	config.Logger.SubMsg.Info().Msg("Ok!")

	config.SetDomain("kaatinga.ru")
	config.SetEmail("info@kaatinga.ru")
	config.SetLaunchMode("dev")
	config.SetPort(myEnvs.Port)

	err = config.Launch(SetUpHandlers)
	if err != nil {
		config.Logger.SubMsg.Err(err).Msg("The server stopped")
	}
}

func SetUpHandlers(r *httprouter.Router, db *sql.DB) {
	r.GET("/:phrase", HelloServer)
	r.GET("/", HelloServer)
}

func HelloServer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
