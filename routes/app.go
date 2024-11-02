package routes

import (
	"github.com/ForkBench/Inseki-Core/tools"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
	"inseki-desk/components"
	"inseki-desk/core"
	"inseki-desk/pages"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

/*
Create a new chi router, configure it and return it.
*/
func NewChiRouter(configJson string) *chi.Mux {

	err, config := tools.ReadEmbedConfigFile(configJson)
	if err != nil {
		panic(err)
	}

	// Check if the folder exists
	err = tools.CheckIfConfigFolderExists(config)
	if err != nil {
		panic(err)
	}

	err, insekiIgnore := tools.ReadInsekiIgnore(config)
	if err != nil {
		panic(err)
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	analyzer := core.Analyze{Config: config, InsekiIgnore: insekiIgnore, Home: homedir}

	mainFolders := analyzer.GetMainFolders()

	r := chi.NewRouter()

	// Useful middleware, see : https://pkg.go.dev/github.com/go-chi/httplog/v2@v2.1.1#NewLogger
	logger := httplog.NewLogger("app-logger", httplog.Options{
		// All log
		LogLevel: slog.LevelInfo,
		Concise:  true,
	})

	// Use the logger and recoverer middleware.
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)

	// ULTRA IMPORTANT : This middleware is used to prevent caching of the pages.
	// Sometimes, HX requests may be cached by the browser, which may cause unexpected behavior.
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			next.ServeHTTP(w, r)
		})
	})

	// Serve static files.
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		HXRender(w, r, pages.HomePage(), mainFolders)

		// 200 OK status
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/query", func(w http.ResponseWriter, r *http.Request) {
		// Get the path from the query
		path := r.URL.Query().Get("path")

		// Unb64 the path
		file, err := core.FileFromUrl(path)
		if err != nil {
			HXRender(w, r, pages.PathNotFound(), mainFolders)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Check if the file is a directory
		if file.IconPath != "/folder.png" {
			HXRender(w, r, pages.Reader(file), mainFolders)
			w.WriteHeader(http.StatusOK)
			return
		}

		// List all subfiles
		files := analyzer.ListAllSubFiles(file.Path)

		HXRender(w, r, pages.QueryPage(file, files), mainFolders)

		// 200 OK status
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/load-content", func(w http.ResponseWriter, r *http.Request) {

		// TODO: Make it customizable
		err, responses := analyzer.Process("~/Documents/")
		if err != nil {
			HXRender(w, r, components.ErrorMsg(err), mainFolders)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		projects := make([]core.File, 0)

		for _, response := range responses {
			projects = append(projects, core.File{
				FileName: filepath.Base(response.Filepath),
				Path:     response.Filepath,
			})
			println(response.Filepath)
		}

		HXRender(w, r, components.Projects(projects), mainFolders)

		// 200 OK status
		w.WriteHeader(http.StatusOK)
	})

	// Listen to port 3000.
	go func() {
		err := http.ListenAndServe(":9245", r)
		if err != nil {
			panic(err)
		}
	}()

	return r
}
