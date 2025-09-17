package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
	"github.com/signintech/gopdf"
)

func checkErr(err error) {
	if err != nil {
		log.Print(err.Error())
	}
}

func writeFile(buf []byte) error {
	tmpPath := "./assets/imgs"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "line-chart-1-basic.png")
	return os.WriteFile(file, buf, 0600)
}

func Chart() {
	values := [][]float64{
		{120, 132, 101, charts.GetNullValue(), 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}

	opt := charts.NewLineChartOptionWithData(values)
	opt.Title.Text = "Line"
	opt.Title.FontStyle.FontSize = 16
	opt.XAxis.Labels = []string{
		"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
	}
	opt.Legend.SeriesNames = []string{
		"Email", "Union Ads", "Video Ads", "Direct", "Search Engine",
	}
	opt.Legend.Padding = charts.Box{
		Left: 100,
	}
	opt.Symbol = charts.SymbolCircle
	opt.LineStrokeWidth = 1.2

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        600,
		Height:       400,
	})
	if err := p.LineChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}

func main() {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	// 新增中文字體
	err := pdf.AddTTFFont("NotoSansTC", "assets/fonts/NotoSansTC-Black.ttf")
	checkErr(err)
	err = pdf.SetFont("NotoSansTC", "", 14)
	checkErr(err)

	pdf.AddHeader(func() {
		pdf.SetFontSize(14)
		pdf.SetX(160)
		pdf.SetY(20)
		pdf.Cell(nil, "Go Report Generator: 報表產生器")
		pdf.Image("assets/imgs/logo.png", 100, 10, &gopdf.Rect{W: 50, H: 50})

		pdf.SetXY(160, 50)
		pdf.Text("歡迎來到我的website(https://www.edwin.io)")
		pdf.AddExternalLink("https://www.edwin.io", 160, 50, 160, 10)
	})

	pdf.AddFooter(func() {
		pdf.SetFontSize(14)
		pdf.SetX(20)
		pdf.SetY(800)
		pdf.Cell(nil, "Footer")
	})

	pdf.AddPage()
	pdf.SetLineType("solid")
	pdf.Line(20, 70, 570, 70)
	pdf.SetY(90)
	pdf.Text("Table 範例:")

	// 設定表格Y起始位置
	tableStartY := 110.0
	// 設定表格左邊距
	marginLeft := 10.0
	table := pdf.NewTableLayout(marginLeft, tableStartY, 25, 0)
	// 表格欄位設定
	table.AddColumn("ITEM", 50, "center")
	table.AddColumn("DESCRIPTION", 200, "center")
	table.AddColumn("QTY.", 40, "center")
	table.AddColumn("PRICE", 60, "center")
	table.AddColumn("TOTAL", 60, "center")

	// 表格標題列
	table.AddRow([]string{"001", "Product A", "2", "10.00", "20.00"})
	table.AddRow([]string{"002", "Product B", "1", "15.00", "15.00"})
	table.AddRow([]string{"003", "Product C", "3", "5.00", "15.00"})

	// Set the style for table cells
	table.SetTableStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:    false,
			Left:   false,
			Bottom: true,
			Right:  false,
			Width:  1.0,
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      "NotoSansTC",
		FontSize:  12,
	})

	// Set the style for table header
	table.SetHeaderStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      false,
			Left:     false,
			Bottom:   true,
			Right:    false,
			Width:    2.0,
			RGBColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		},
		FillColor: gopdf.RGBColor{R: 165, G: 165, B: 165},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      "NotoSansTC",
		FontSize:  14,
	})

	table.SetCellStyle(gopdf.CellStyle{
		BorderStyle: gopdf.BorderStyle{
			Top:      true,
			Left:     false,
			Right:    false,
			Bottom:   true,
			Width:    0.5,
			RGBColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		},
		FillColor: gopdf.RGBColor{R: 255, G: 255, B: 255},
		TextColor: gopdf.RGBColor{R: 0, G: 0, B: 0},
		Font:      "NotoSansTC",
		FontSize:  10,
	})

	// 表格繪製
	table.DrawTable()

	pdf.AddPage()
	pdf.SetY(400)
	Chart()

	// 插入圖片
	err = pdf.Image("assets/imgs/line-chart-1-basic.png", 50, 100, &gopdf.Rect{W: 500, H: 300})
	checkErr(err)
	pdf.WritePdf("assets/pdf/header-footer.pdf")

}
