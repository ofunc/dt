package xlsx_test

import (
	"fmt"
	"testing"

	"github.com/ofunc/dt/io/xlsx"
)

func TestXLSX(t *testing.T) {
	f, err := xlsx.NewReader("/").Drop(1).Head(2).Tail(1).ReadFile("工作簿1.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(f)

	fmt.Println(xlsx.Save(f, "save.xlsx", "测试"))
}
