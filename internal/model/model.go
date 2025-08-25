package model

import (
	"pacmanager/internal/packages"
	"pacmanager/internal/ui"

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
	ShowDetails      bool     // 是否显示详情信息
	Dependencies     []string // 依赖列表
	ReverseDependencies []string // 反向依赖列表
	NeedFetchDeps    bool     // 是否需要获取依赖信息
	Initialized      bool     // 是否已初始化包信息
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

// InitializePackages 初始化包信息
func (m *Model) InitializePackages() {
	if !m.Initialized {
		// 获取所有包信息
		m.AllPackages = packages.GetPackages()
		
		// 过滤包（仅显示显式安装的包）
		m.FilteredPackages = m.FilterPackages(1)
		
		// 创建表格
		m.Table = ui.CreateTable(m.FilteredPackages)
		
		// 默认选中第一个包
		if len(m.FilteredPackages) > 0 {
			m.SelectedPackage = &m.FilteredPackages[0]
		}
		
		m.Initialized = true
	}
}

// GetCurrentPackageVersion 获取当前选中包的版本信息
func (m *Model) GetCurrentPackageVersion() string {
	if m.SelectedPackage != nil {
		// 如果版本信息为空，获取版本信息
		if m.SelectedPackage.Version == "" {
			m.SelectedPackage.Version = packages.GetPackageVersion(m.SelectedPackage.Name)
		}
		return m.SelectedPackage.Version
	}
	return "未知"
}
