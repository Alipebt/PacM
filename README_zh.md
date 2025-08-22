# PacM

[English](README.md)

PacM 是一个基于expac的终端包管理器界面，专为 Arch Linux 系统设计。它提供了一个直观的 TUI（文本用户界面）来浏览和管理通过 pacman 安装的软件包。

## 功能特性

- **可视化包管理**：以表格形式展示所有已安装的软件包
- **包分类**：区分显式安装的包和依赖包
- **详细信息**：显示包的安装日期、版本等详细信息
- **备注功能**：可以为包添加自定义备注
- **交互式界面**：使用键盘导航和操作
- **过滤功能**：可以只显示显式安装的包

## 安装要求

- Arch Linux 或其他支持 pacman 的系统
- [expac](https://github.com/falconindy/expac) 工具（用于获取包的详细信息）
- Go 1.24 或更高版本（用于编译）

## 安装步骤

1. 安装 expac：
   ```bash
   sudo pacman -S expac
   ```

2. 克隆仓库：
   ```bash
   git clone https://github.com/Alipebt/pacm.git
   cd pacm
   ```

3. 安装依赖：
   ```bash
   go mod tidy
   ```

4. 构建项目：
   ```bash
   go build -o pacm
   ```


## 使用方法

运行程序：
```bash
./pacm
```



## 许可证

本项目采用 MIT 许可证。详情请见 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进此项目。

