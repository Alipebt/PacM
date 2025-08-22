package model

import (
	"fmt"
	"pacmanager/internal/ui"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//TODO：持久化存储备注或借助软件包备注程序
//TODO：备注编辑界面移动光标


func (m Model) Init() tea.Cmd {
	// 在初始化时获取所有包信息
	return tea.Batch(
		tea.SetWindowTitle("包管理器"),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.ShowMenu {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "down":
				if m.Choice < 2 {
					m.Choice++
				}
			case "up":
				if m.Choice > 0 {
					m.Choice--
				}
			case "enter":
				if m.Choice == 2 { // 退出选项
					return m, tea.Quit
				}
				m.ShowMenu = false
				// 根据选择过滤包
				m.FilteredPackages = m.FilterPackages(m.Choice)
				m.Table = ui.CreateTable(m.FilteredPackages)
				// 默认选中第一个包
				if len(m.FilteredPackages) > 0 {
					m.SelectedPackage = &m.FilteredPackages[0]
				}
			case "ctrl+c", "esc":
				return m, tea.Quit
			case "q":
				return m, tea.Quit
			}
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
				m.ShowMenu = true
				m.Table = table.Model{} // 清空表格
				m.SelectedPackage = nil
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
			default:
				// 当表格选中项改变时，自动更新选中的包
				m.SelectedPackage = m.GetCurrentSelectedPackage()
			}
		}
		return m, cmd
	}
	return m, nil
}

func (m Model) View() string {
	if m.ShowMenu {
		// 显示菜单
		menu := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render("选择要显示的包类型:")
		choice1 := "  显示所有包"
		choice2 := "  仅显示显式安装的包"
		choice3 := "  退出"
		
		if m.Choice == 0 {
			choice1 = "> " + lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("显示所有包")
		} else {
			choice1 = "  " + choice1
		}
		
		if m.Choice == 1 {
			choice2 = "> " + lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("仅显示显式安装的包")
		} else {
			choice2 = "  " + choice2
		}
		
		if m.Choice == 2 {
			choice3 = "> " + lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("退出")
		} else {
			choice3 = "  " + choice3
		}
		
		// 添加菜单描述
		description := "退出:q，选择:enter"
		
		return lipgloss.JoinVertical(lipgloss.Left, menu, choice1, choice2, choice3, "", description)
	} else if m.EditingNotes {
		// 显示备注编辑界面
		title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render("编辑备注")
		pkgInfo := fmt.Sprintf("包名: %s", m.SelectedPackage.Name)
		notesLabel := "备注:"
		notesInput := m.NewNotes + "█" // 添加光标
		
		instructions := "Enter: 保存, Esc: 取消"
		
		return lipgloss.JoinVertical(lipgloss.Left, title, "", pkgInfo, notesLabel, notesInput, "", instructions)
	} else {
		// 显示表格和选中包的详细信息
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
		
		// 显示选中包的详细信息
		if m.SelectedPackage != nil {
			detailsText := fmt.Sprintf("选中的包详情:\n包名: %s\n版本: %s\n安装日期: %s\n安装类型: %s",
				m.SelectedPackage.Name,
				m.SelectedPackage.Version,
				m.SelectedPackage.InstallDate.Format("2006-01-02 15:04:05"),
				func() string {
					if m.SelectedPackage.Explicit {
						return "显式安装"
					}
					return "依赖安装"
				}(),
			)
			
			// 添加备注信息
			if m.SelectedPackage.Notes != "" {
				detailsText += fmt.Sprintf("\n备注: %s", m.SelectedPackage.Notes)
			} else {
				detailsText += "\n备注: (无)"
			}
			
			details := lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")).
				Padding(1).
				Render(detailsText)
			view = lipgloss.JoinHorizontal(lipgloss.Top, view, "  ", details)
		}
		
		return view + "\n\n导航: ↑/↓ ，翻页: ←/→ ，编辑备注: e ，返回: q ,退出: esc"
	}
}
