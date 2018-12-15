package template

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestParseTmpls(t *testing.T) {
	wd, _ := os.Getwd()
	viper.Set("gohub.ui.templates", wd)

	ParseTmpls(viper.Sub("gohub.ui"))

	if store.Name() != "store" {
		t.Fail()
	}

	fmt.Println(store.DefinedTemplates())
}
