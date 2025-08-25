package model

import (
	"fmt"
	"pacmanager/internal/packages"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//TODO：持久化存储备注或借助软件包备注程序
//TODO：备注编辑界面移动光标


func (m Model) Init() tea.Cmd {
	// 不在初始化时获取所有包信息，而是在需要时获取
	return tea.Batch(
		tea.SetWindowTitle("包管理器"),
	)
}

// 定义自定义消息类型
type fetchDepsMsg struct {
	dependencies     []string
	reverseDependencies []string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// 如果还没有初始化包信息，先初始化
	if !m.Initialized {
		// 直接调用初始化方法
		modelPtr := &m
		modelPtr.InitializePackages()
		
		// 返回更新后的模型和nil命令
		return m, nil
	}
	
	if m.ShowDetails {
		// 处理详情信息显示状态
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q":
				// 返回表格视图
				m.ShowDetails = false
				m.Dependencies = nil
				m.ReverseDependencies = nil
				m.NeedFetchDeps = false
			case "ctrl+c":
				return m, tea.Quit
			}
		case fetchDepsMsg:
			// 接收依赖信息
			m.Dependencies = msg.dependencies
			m.ReverseDependencies = msg.reverseDependencies
			m.NeedFetchDeps = false
		}
	} else if m.EditingNotes {
		// 处理备注编辑状态
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// 保存备注
				if m.SelectedPackage != nil {
					m.SelectedPackage.Notes = m.NewNotes
					// 更新所有包列表中的备注
					for i := range m.AllPackages {
						if m.AllPackages[i].Name == m.SelectedPackage.Name {
							m.AllPackages[i].Notes = m.NewNotes
							break
						}
					}
				}
				m.EditingNotes = false
				m.NewNotes = ""
			case "esc":
				// 取消编辑
				m.EditingNotes = false
				m.NewNotes = ""
			case "backspace":
				if len(m.NewNotes) > 0 {
					m.NewNotes = m.NewNotes[:len(m.NewNotes)-1]
				}
			default:
				// 添加字符到备注
				if msg.Type == tea.KeyRunes {
					m.NewNotes += msg.String()
				}
				// 显式处理空格键
				if msg.String() == " " {
					m.NewNotes += " "
				}
			}
		}
	} else {
		var cmd tea.Cmd
		m.Table, cmd = m.Table.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				return m, tea.Quit
			case "ctrl+c":
				return m, tea.Quit
			case "q":
				return m, tea.Quit
			case "left":
				m.Table.MoveUp(13)
				m.SelectedPackage = &m.FilteredPackages[m.GetCurrentSelectedIndex()]
			case "right":
				m.Table.MoveDown(13)
				m.SelectedPackage = &m.FilteredPackages[m.GetCurrentSelectedIndex()]
			case "e":
				// 编辑备注
				if m.SelectedPackage != nil {
					m.EditingNotes = true
					m.NewNotes = m.SelectedPackage.Notes
				}
			case "enter":
				// 显示详情信息
				if m.SelectedPackage != nil {
					m.ShowDetails = true
					// 不再立即获取依赖和反向依赖信息，延后到查看详情时获取
					m.Dependencies = nil
					m.ReverseDependencies = nil
					m.NeedFetchDeps = true
					// 返回获取依赖的命令
					return m, fetchDependencies(m.SelectedPackage.Name)
				}
			default:
				// 当表格选中项改变时，自动更新选中的包
				m.SelectedPackage = m.GetCurrentSelectedPackage()
			}
		}
		return m, cmd
	}
	return m, nil
}

// 获取依赖信息的命令
func fetchDependencies(packageName string) tea.Cmd {
	return func() tea.Msg {
		dependencies := packages.GetDependencies(packageName)
		reverseDependencies := packages.GetReverseDependencies(packageName)
		return fetchDepsMsg{
			dependencies:     dependencies,
			reverseDependencies: reverseDependencies,
		}
	}
}

func (m Model) View() string {
	if m.ShowDetails {
		// 显示详情信息
		if m.SelectedPackage == nil {
			return "未选择包"
		}
		
		// 获取当前选中包的版本信息
		modelPtr := &m
		pkgVersion := modelPtr.GetCurrentPackageVersion()
		
		// 延后获取依赖和反向依赖信息
		if m.Dependencies == nil && m.ReverseDependencies == nil && !m.NeedFetchDeps {
			m.NeedFetchDeps = true
			// 不再在这里直接获取依赖信息，而是在Update函数中通过命令获取
		}
		
		// 构建详情信息
		details := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render("包详情信息")
		pkgName := fmt.Sprintf("包名: %s", m.SelectedPackage.Name)
		pkgVersionStr := fmt.Sprintf("版本: %s", pkgVersion)
		pkgSize := fmt.Sprintf("大小: %s", m.SelectedPackage.Size)
		pkgInstallDate := fmt.Sprintf("安装日期: %s", m.SelectedPackage.InstallDate.Format("2006-01-02 15:04:05"))
		pkgInstallType := fmt.Sprintf("安装类型: %s", func() string {
			if m.SelectedPackage.Explicit {
				return "显式安装"
			}
			return "依赖安装"
		}())
		
		// 添加备注信息
		notes := "备注: "
		if m.SelectedPackage.Notes != "" {
			notes += m.SelectedPackage.Notes
		} else {
			notes += "(无)"
		}
		
		// 构建依赖信息
		dependenciesTitle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render("\n依赖 (显式安装的软件):")
		var dependenciesList string
		if len(m.Dependencies) > 0 {
			for _, dep := range m.Dependencies {
				dependenciesList += fmt.Sprintf("  %s\n", dep)
			}
		} else {
			dependenciesList = "  (无)\n"
		}
		
		// 构建反向依赖信息
		reverseDependenciesTitle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render("\n反向依赖 (依赖安装):")
		var reverseDependenciesList string
		if len(m.ReverseDependencies) > 0 {
			for _, dep := range m.ReverseDependencies {
				reverseDependenciesList += fmt.Sprintf("  %s\n", dep)
			}
		} else {
			reverseDependenciesList = "  (无)\n"
		}
		
		instructions := "\n导航: Esc/q 返回, Ctrl+C 退出"
		
		return lipgloss.JoinVertical(lipgloss.Left, 
			details, 
			pkgName, 
			pkgVersionStr, 
			pkgSize,
			pkgInstallDate, 
			pkgInstallType, 
			notes, 
			dependenciesTitle, 
			dependenciesList, 
			reverseDependenciesTitle, 
			reverseDependenciesList, 
			instructions)
	} else if m.EditingNotes {
		// 显示备注编辑界面
		title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render("编辑备注")
		pkgInfo := fmt.Sprintf("包名: %s", m.SelectedPackage.Name)
		notesLabel := "备注:"
		notesInput := m.NewNotes + "█" // 添加光标
		
		instructions := "Enter: 保存, Esc: 取消"
		
		return lipgloss.JoinVertical(lipgloss.Left, title, "", pkgInfo, notesLabel, notesInput, "", instructions)
	} else {
		// 显示表格
		view := m.Table.View()
		
		// 计算软件包总数和当前选中软件包的百分比
		totalPackages := len(m.FilteredPackages)
		currentIndex := m.Table.Cursor() + 1
		percentage := 0
		if totalPackages > 0 {
			percentage = (currentIndex * 100) / totalPackages
		}
		
		// 添加统计信息
		stats := fmt.Sprintf("当前位置: %d/%d (%d%%)", currentIndex, totalPackages, percentage)
		view = lipgloss.JoinVertical(lipgloss.Left, view, stats)
		
		return view + "\n\n导航: ↑/↓ ，翻页: ←/→ ，编辑备注: e ，查看详情: Enter ，退出: q"
	}
}
