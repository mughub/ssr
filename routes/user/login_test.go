package user

import (
	"github.com/golang/mock/gomock"
	"github.com/graphql-go/graphql"
	"github.com/mughub/mughub/db"
	"github.com/mughub/mughub/db/dbtest"
	"github.com/mughub/ssr/template"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	wd, _ := os.Getwd()
	uiDir := filepath.Dir(filepath.Dir(wd))

	viper.Set("gohub.ui.templates", filepath.Join(uiDir, "template"))
	template.ParseTmpls(viper.Sub("gohub.ui"))
	os.Exit(m.Run())
}

func TestGetLoginTmpl(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://localhost:8080", nil)

	GetLoginTmpl(w, req)

	resp := w.Result()
	n, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		t.Errorf("unexpected error from io.Copy: %s", err)
	}

	if n == 0 {
		t.Fail()
	}
}

const emailLoginStr = "{ login(email: $email, password: $password) { id } }"
const usernameLoginStr = "{ login(username: $username, password: $password) { id } }"

func TestPostLoginTmpl(t *testing.T) {
	// Create mock database for handling login queries
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := dbtest.NewMockInterface(ctrl)
	db.RegisterDB(mockDB)

	// Set up test cases
	testCases := []struct {
		Name string
		Req  string
		Res  graphql.Result // The data field in result does not reflect actual return type.
	}{
		{
			Name: "email",
			Req:  emailLoginStr,
			Res: graphql.Result{
				Data: "1",
			},
		},
	}

	// Run test cases
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(subT *testing.T) {
			mockDB.EXPECT().Do(gomock.Any(), testCase.Req, gomock.Any()).Return(testCase.Res)
		})
	}
}
