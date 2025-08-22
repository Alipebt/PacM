package ui

import (
	"fmt"
	"pacmanager/internal/packages"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// CreateTable 创建表格
func CreateTable(packages []packages.PackageInfo) table.Model {
	// 定义列
	columns := []table.Column{
		{Title: "序号", Width: 6},
		{Title: "安装日期", Width: 20},
		{Title: "包名", Width: 25},
		{Title: "类型", Width: 10},
	}

	// 填充行数据
	var rows []table.Row
	for i, pkg := range packages {
		installType := "依赖"
		if pkg.Explicit {
			installType = "显式"
		}
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			pkg.InstallDate.Format("2006-01-02 15:04:05"),
			pkg.Name,
			installType,
		})
	}

	// 创建表格
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15), // 固定高度以便更好地显示
	)

	// 设置样式
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
