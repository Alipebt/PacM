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
	cmd := exec.Command("sh", "-c", "expac -H M --timefmt='%F %T' '%l\\t%n\\t%m' | sort -rn")
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
		if len(parts) >= 3 {
			// 第一部分是日期时间
			dateTimeStr := parts[0]
			
			// 第二部分是包名
			name := parts[1]
			
			// 第三部分是大小
			size := parts[2]
			
			// 解析日期时间
			installTime, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
			if err != nil {
				fmt.Printf("解析日期时间出错: %v\n", err)
				continue
			}
			
			// 存储包信息
			pkg := PackageInfo{
				InstallDate: installTime,
				Name:        name,
				Size:        size,
				Explicit:    explicit[name],
			}
			
			packages = append(packages, pkg)
		}
	}

	// 按安装日期反向排序
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].InstallDate.After(packages[j].InstallDate)
	})

	return packages
}

// GetDependencies 获取包的依赖列表
func GetDependencies(packageName string) []string {
	// 检查是否在支持 pactree 的系统上
	_, err := exec.LookPath("pactree")
	if err != nil {
		// 使用模拟数据进行测试
		fmt.Println("警告: 未找到 pactree")
		return []string{"依赖1", "依赖2", "依赖3"}
	}
	
	// 执行命令获取依赖
	cmd := exec.Command("pactree", "-d", "1", packageName)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("执行 pactree -d 1 %s 出错: %v\n", packageName, err)
		return []string{}
	}
	
	var dependencies []string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	// 跳过第一行（包名本身）
	if scanner.Scan() {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// 移除前导的符号（如├─, └─等）
			if strings.Contains(line, " ") {
				parts := strings.Split(line, " ")
				if len(parts) > 1 {
					dependencies = append(dependencies, parts[1])
				}
			} else {
				dependencies = append(dependencies, line)
			}
		}
	}
	
	return dependencies
}

// GetReverseDependencies 获取包的反向依赖列表
func GetReverseDependencies(packageName string) []string {
	// 检查是否在支持 pactree 的系统上
	_, err := exec.LookPath("pactree")
	if err != nil {
		// 使用模拟数据进行测试
		fmt.Println("警告: 未找到 pactree")
		return []string{"反向依赖1", "反向依赖2", "反向依赖3"}
	}
	
	// 执行命令获取反向依赖
	cmd := exec.Command("pactree", "-r", "-d", "1", packageName)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("执行 pactree -r -d 1 %s 出错: %v\n", packageName, err)
		return []string{}
	}
	
	var reverseDependencies []string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	// 跳过第一行（包名本身）
	if scanner.Scan() {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// 移除前导的符号（如├─, └─等）
			if strings.Contains(line, " ") {
				parts := strings.Split(line, " ")
				if len(parts) > 1 {
					reverseDependencies = append(reverseDependencies, parts[1])
				}
			} else {
				reverseDependencies = append(reverseDependencies, line)
			}
		}
	}
	
	return reverseDependencies
}

// GetPackageVersion 获取包的版本信息
func GetPackageVersion(packageName string) string {
	// 检查是否在支持 pacman 的系统上
	_, err := exec.LookPath("pacman")
	if err != nil {
		// 使用模拟数据进行测试
		fmt.Println("警告: 未找到 pacman")
		return "1.0.0"
	}
	
	// 执行命令获取版本信息
	cmd := exec.Command("pacman", "-Q", packageName)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("执行 pacman -Q %s 出错: %v\n", packageName, err)
		return "未知"
	}
	
	// 解析输出获取版本信息
	output := strings.TrimSpace(string(out))
	parts := strings.Fields(output)
	if len(parts) >= 2 {
		return parts[1]
	}
	
	return "未知"
}

// GetPackages 获取所有包信息
func GetPackages() []PackageInfo {
	explicitPackages := GetExplicitPackages()
	allPackages := GetAllPackages(explicitPackages)
	
	// 不再在初始化时获取版本信息，而是在需要时获取
	// 版本信息将在查看详情时获取
	
	return allPackages
}
