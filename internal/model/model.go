package model

import (
	"pacmanager/internal/packages"

	"github.com/charmbracelet/bubbles/table"
)

// Model represents the application state
type Model struct {
	Table            table.Model
	Cursor           int
	Choice           int
	ShowMenu         bool
	AllPackages      []packages.PackageInfo
	FilteredPackages []packages.PackageInfo
	SelectedPackage  *packages.PackageInfo
	EditingNotes     bool     // 是否正在编辑备注
	NewNotes         string   // 新备注内容
}

// GetCurrentSelectedIndex 获取当前选中的行
func (m Model) GetCurrentSelectedIndex() int {
    // 获取当前选中的行
    return m.Table.Cursor()
}

// FilterPackages 根据选择过滤包
func (m Model) FilterPackages(choice int) []packages.PackageInfo {
	if choice == 1 {
		// 仅返回显式安装的包
		var explicitOnly []packages.PackageInfo
		for _, pkg := range m.AllPackages {
			if pkg.Explicit {
				explicitOnly = append(explicitOnly, pkg)
			}
		}
		return explicitOnly
	}
	
	// 返回所有包
	return m.AllPackages
}

// GetCurrentSelectedPackage 获取当前选中的包
func (m Model) GetCurrentSelectedPackage() *packages.PackageInfo {
	if len(m.FilteredPackages) == 0 {
		return nil
	}
	
	// 获取当前表格选中的行索引
	cursor := m.Table.Cursor()
	if cursor >= 0 && cursor < len(m.FilteredPackages) {
		return &m.FilteredPackages[cursor]
	}
	
	return nil
}
