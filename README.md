# PacM

[中文](README_zh.md)

PacM is a terminal-based package manager interface based on expac designed for Arch Linux systems . It provides an intuitive TUI (Text User Interface) to browse and manage software packages installed via pacman.

## Features

- **Visual Package Management**: Display all installed packages in a table format
- **Package Classification**: Distinguish between explicitly installed packages and dependencies
- **Detailed Information**: Show installation date, version, and other details of packages
- **Note Functionality**: Add custom notes to packages
- **Interactive Interface**: Navigate and operate using keyboard
- **Filter Function**: Display only explicitly installed packages

## Requirements

- Arch Linux or other systems supporting pacman
- [expac](https://github.com/falconindy/expac) tool (for retrieving package details)
- Go 1.24 or higher (for compilation)

## Installation

1. Install expac :
   ```bash
   sudo pacman -S expac
   ```

2. Clone the repository:
   ```bash
   git clone https://github.com/Alipebt/pacm.git
   cd pacm
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Build the project:
   ```bash
   go build -o pacm
   ```



## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Issues and Pull Requests are welcome to improve this project.

