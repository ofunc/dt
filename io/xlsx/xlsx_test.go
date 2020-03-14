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
}

func TestSave(t *testing.T) {
	xl, err := xlsx.OpenFile("工作簿1.xlsx")
	if err != nil {
		t.Fatal(err)
	}
	sheet := xl.Sheet("")
	data := sheet.Data()
	data.Rows[0].Cells[0].Type = "inlineStr"
	data.Rows[0].Cells[0].Value = "<测试成功！@#￥%……>"
	if err := sheet.Update(); err != nil {
		t.Fatal(err)
	}
	if err := xl.SaveFile("out.xlsx"); err != nil {
		t.Fatal(err)
	}
}
