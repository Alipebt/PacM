package main

import (
	"fmt"
	"os"
	"pacm/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// 初始化模型
	m := model.Model{
		ShowMenu:         false,
		Choice:           1,   // 1表示仅显示显式安装的包
		AllPackages:      nil, // 延迟获取
		FilteredPackages: nil, // 延迟获取
	}

	// 运行程序
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}
