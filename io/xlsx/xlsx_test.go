package xlsx_test

import (
	"fmt"
	"testing"

	"github.com/ofunc/dt/io/xlsx"
)

func TestXLSX(t *testing.T) {
	f, err := xlsx.NewReader("/").Drop(1).Head(3).Tail(1).ReadFile("新建文档2020-03-14 17_56_39.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(f)

	err = xlsx.Saver{
		Template: "tmpl.xlsx",
		File:     "out.xlsx",
		Sheet:    "测试",
	}.Save(f)
	if err != nil {
		t.Fatal(err)
	}
}
