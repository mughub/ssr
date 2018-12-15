package routes

import (
	"github.com/mughub/ssr/template"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

var (
	goOnce sync.Once
)

func TestMain(m *testing.M) {
	wd, _ := os.Getwd()
	uiDir := filepath.Dir(wd)

	viper.Set("gohub.ui.templates", filepath.Join(uiDir, "template"))
	template.ParseTmpls(viper.Sub("gohub.ui"))
	os.Exit(m.Run())
}

func TestGetHomeTmpl(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost:8080", nil)

	GetHome(w, req)

	resp := w.Result()
	n, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		t.Errorf("unexpected error from io.Copy: %s", err)
	}

	if n == 0 {
		t.Fail()
	}
}
