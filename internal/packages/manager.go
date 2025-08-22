package packages

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// GetExplicitPackages 获取显式安装的包列表
func GetExplicitPackages() map[string]bool {
	// 检查是否在支持 pacman 的系统上
	_, err := exec.LookPath("pacman")
	if err != nil {
		// 使用模拟数据进行测试
		fmt.Println("警告: 未找到 pacman")
		return nil
	}
	
	// 执行命令
	cmd := exec.Command("pacman", "-Qe")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("执行 pacman -Qe 出错: %v\n", err)
		os.Exit(1)
	}

	explicit := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 1 {
			explicit[fields[0]] = true
		}
	}

	return explicit
}

// GetAllPackages 获取所有包信息
func GetAllPackages(explicit map[string]bool) []PackageInfo {
	// 检查是否在支持 expac 的系统上
	_, err := exec.LookPath("expac")
	if err != nil {
		// 使用模拟数据进行测试
		fmt.Println("警告: 未找到 expac")
		return nil
	}
	
	// 尝试执行实际命令
	cmd := exec.Command("sh", "-c", "expac --timefmt='%F %T' '%l\\t%n %v' | sort -rn")
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("执行 expac 出错: %v\n", err)
		os.Exit(1)
	}

	var packages []PackageInfo
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			// 第一部分是日期时间
			dateTimeStr := parts[0]
			
			// 第二部分包含包名和版本号
			nameVersion := parts[1]
			nameVersionParts := strings.Fields(nameVersion)
			if len(nameVersionParts) >= 2 {
				name := nameVersionParts[0]
				version := strings.Join(nameVersionParts[1:], " ")
				
				// 解析日期时间
				installTime, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
				if err != nil {
					fmt.Printf("解析日期时间出错: %v\n", err)
					continue
				}
				
				pkg := PackageInfo{
					InstallDate: installTime,
					Name:        name,
					Version:     version,
					Explicit:    explicit[name],
				}
				
				packages = append(packages, pkg)
			}
		}
	}

	// 按安装日期反向排序
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].InstallDate.After(packages[j].InstallDate)
	})

	return packages
}

// GetPackages 获取所有包信息
func GetPackages() []PackageInfo {
	explicitPackages := GetExplicitPackages()
	allPackages := GetAllPackages(explicitPackages)
	return allPackages
}
