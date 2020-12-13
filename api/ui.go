package api

import (
	"html/template"
	"net/http"
	"os"

	"github.com/bodenr/opsyc/util"
)

type UIHandler struct{}

var AssetsDir string

func init() {
	assetsDir, exists := os.LookupEnv("OPSYC_ASSETS_DIR")
	if !exists {
		assetsDir = "assets"
	}

	AssetsDir = assetsDir
}

func (ui *UIHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var page string
	page, req.URL.Path = util.ShiftPath(req.URL.Path)
	if req.URL.Path != "/" {
		http.Error(res, "Not Found", http.StatusNotFound)
		return
	}
	if req.Method != "GET" {
		http.Error(res, "Only GET supported", http.StatusBadRequest)
		return
	}
	switch page {
	case "runtime.html":
		runtimePageHandler(res, req)
	default:
		http.Error(res, "Not Found", http.StatusNotFound)
	}
}

func runtimePageHandler(resp http.ResponseWriter, req *http.Request) {
	page := template.Must(template.ParseFiles(AssetsDir + "/html/runtime.html"))
	page.Execute(resp, NewRuntimeEnv())
}
