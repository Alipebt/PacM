package main

import (
	"fmt"
	"os"
	"pacmanager/internal/model"
	"pacmanager/internal/packages"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// 在程序启动时就获取所有包信息
	allPackages := packages.GetPackages()
	
	// 初始化模型
	m := model.Model{
		ShowMenu:    true,
		Choice:      0,
		AllPackages: allPackages,
	}

	// 运行程序
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}
}
