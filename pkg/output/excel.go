package output

import (
	"fmt"
	"github.com/TARI0510/gosca/pkg/config"
	"github.com/TARI0510/gosca/pkg/util/utils"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func ExcelVuln(vulIdList []string, vulnDbMap map[string]config.VulnDb, stdPkgLocationMap map[string][]string, f F) {
	xlsx := excelize.NewFile()

	SHEET1 := "Sheet1"

	cellBorder := []excelize.Border{
		{
			Type:  "left",
			Color: "#000000",
			Style: 1,
		}, {
			Type:  "top",
			Color: "#000000",
			Style: 1,
		}, {
			Type:  "bottom",
			Color: "#000000",
			Style: 1,
		}, {
			Type:  "right",
			Color: "#000000",
			Style: 1,
		},
	}
	centerCell := &excelize.Alignment{
		Vertical:   "center",
		Horizontal: "left",
		WrapText:   true,
	}

	globalStyle, _ := xlsx.NewStyle(&excelize.Style{
		Alignment: centerCell,
		Border:    cellBorder,
	})

	firstLineStyle, _ := xlsx.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Alignment: centerCell,
		Border:    cellBorder,
	})

	severityMap := map[string]string{
		"CRITICAL": "#000000",
		"HIGH":     "#d9534f",
		"MEDIUM":   "#ec971f",
		"LOW":      "#497398",
		"else":     "#DCDCDC",
	}

	severityPriority := map[string]int{
		"CRITICAL": 0,
		"HIGH":     1,
		"MEDIUM":   2,
		"LOW":      3,
	}

	titles := []string{"ID", "CVE", "Severity", "Module", "Package", "Affected Version", "Detail"}
	for i, title := range titles {
		// A B C
		xlsx.SetCellValue(SHEET1, fmt.Sprintf("%s1", utils.ExcelDiv(i+1)), title)
		xlsx.SetCellStyle(SHEET1, fmt.Sprintf("%s1", utils.ExcelDiv(i+1)), fmt.Sprintf("%s1", utils.ExcelDiv(i+1)), firstLineStyle)
	}

	for i, vulId := range vulIdList {
		vulnDb := vulnDbMap[vulId]

		// Get the most severe vulnerability color
		var color string
		var colorPriority int = 9999
		for _, severity := range vulnDb.Severities {
			if _, ok := severityMap[severity]; ok {
				if severityPriority[severity] <= colorPriority {
					color = severityMap[severity]
					colorPriority = severityPriority[severity]
				}
			} else {
				color = severityMap["else"]
			}
		}

		styleSev, _ := xlsx.NewStyle(&excelize.Style{
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{color},
				Pattern: 1,
			},
			Alignment: centerCell,
			Border:    cellBorder,
		})

		xlsx.SetCellValue(SHEET1, "A"+strconv.Itoa(i+2), vulId)
		xlsx.SetCellStyle(SHEET1, "A"+strconv.Itoa(i+2), "A"+strconv.Itoa(i+2), globalStyle)
		if len(vulnDb.Cves) > 0 {
			xlsx.SetCellValue(SHEET1, "B"+strconv.Itoa(i+2), strings.Join(vulnDb.Cves, "\n"))
			xlsx.SetCellStyle(SHEET1, "B"+strconv.Itoa(i+2), "B"+strconv.Itoa(i+2), globalStyle)
			if vulnDb.Severities != nil {
				xlsx.SetCellValue(SHEET1, "C"+strconv.Itoa(i+2), strings.Join(vulnDb.Severities, "\n"))
				xlsx.SetCellStyle(SHEET1, "C"+strconv.Itoa(i+2), "C"+strconv.Itoa(i+2), globalStyle)
				xlsx.SetCellStyle(SHEET1, "C"+strconv.Itoa(i+2), "C"+strconv.Itoa(i+2), styleSev)
			}
		}
		xlsx.SetCellValue(SHEET1, "D"+strconv.Itoa(i+2), vulnDb.Module)
		xlsx.SetCellStyle(SHEET1, "D"+strconv.Itoa(i+2), "D"+strconv.Itoa(i+2), globalStyle)
		xlsx.SetCellValue(SHEET1, "E"+strconv.Itoa(i+2), vulnDb.Package)
		xlsx.SetCellStyle(SHEET1, "E"+strconv.Itoa(i+2), "E"+strconv.Itoa(i+2), globalStyle)

		var versionsList []string
		for _, v := range vulnDb.Versions {
			version := ""
			if _, ok := v["introduced"]; ok {
				version += ">= " + v["introduced"] + ", "
			}
			if _, ok := v["fixed"]; ok {
				version += "< " + v["fixed"]
			}

			versionsList = append(versionsList, version)
		}
		versions := strings.Join(versionsList, "\n")
		xlsx.SetCellValue(SHEET1, "F"+strconv.Itoa(i+2), versions)
		xlsx.SetCellStyle(SHEET1, "F"+strconv.Itoa(i+2), "F"+strconv.Itoa(i+2), globalStyle)

		xlsx.SetCellValue(SHEET1, "G"+strconv.Itoa(i+2), "https://github.com/golang/vulndb/tree/master/reports/"+vulId+".yaml")
		xlsx.SetCellStyle(SHEET1, "G"+strconv.Itoa(i+2), "G"+strconv.Itoa(i+2), globalStyle)
	}

	// cell adaptive width and height
	cols, err := xlsx.GetCols(SHEET1)
	if err != nil {
		return
	}
	for idx, col := range cols {
		largestWidth := 0
		for _, rowCell := range col {
			cellWidth := utf8.RuneCountInString(rowCell) + 1 // + 2 for margin
			if cellWidth > largestWidth {
				largestWidth = cellWidth
			}
		}
		name, err := excelize.ColumnNumberToName(idx + 1)
		if err != nil {
			return
		}
		xlsx.SetColWidth(SHEET1, name, name, float64(largestWidth))
	}

	var outPath string
	if f.OutPath == "" {
		middle := filepath.Base(f.WorkDir)
		outPath = fmt.Sprintf("GoSCA_%s_%s.xlsx", middle, time.Now().Format("20060102150405"))
	} else {
		outPath = f.OutPath
	}

	if err := xlsx.SaveAs(outPath); err != nil {
		fmt.Fprintf(os.Stderr, "[-] Failed to save %s: %s\n", outPath, err)
		os.Exit(1)
	}
	fmt.Println("[+] Excel file created successfully.")
}
