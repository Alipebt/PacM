package packages

import (
	"time"
)

// PackageInfo represents information about a package
type PackageInfo struct {
	InstallDate time.Time
	Name        string
	Version     string
	Explicit    bool
	Notes       string // 添加备注字段
}
