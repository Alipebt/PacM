package main

import (
	"fmt"
	"os"
	"pacmanager/internal/model"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// 不在程序启动时获取所有包信息，而是在需要时获取
	// 初始化模型，延迟获取包信息
	m := model.Model{
		ShowMenu:         false,
		Choice:           1, // 1表示仅显示显式安装的包
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
