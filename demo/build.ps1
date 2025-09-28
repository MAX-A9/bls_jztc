# 设置错误时停止执行
$ErrorActionPreference = "Stop"

# 设置输出编码为UTF-8
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

# 定义变量
$projectName = "demo"
$buildDir = "build"
$version = "1.0.0"
$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"

# 定义颜色函数
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    else {
        $input | Write-Output
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Write-Info($message) {
    Write-ColorOutput Green "[INFO] $message"
}

function Write-Warning($message) {
    Write-ColorOutput Yellow "[WARN] $message"
}

function Write-Error($message) {
    Write-ColorOutput Red "[ERROR] $message"
}

# 检查Go环境
function Check-GoEnvironment {
    try {
        $goVersion = go version
        Write-Info "检测到Go环境: $goVersion"
    }
    catch {
        Write-Error "未检测到Go环境，请先安装Go"
        exit 1
    }
}

# 创建构建目录
function Create-BuildDirectory {
    Write-Info "创建构建目录..."
    if (Test-Path $buildDir) {
        Remove-Item -Path $buildDir -Recurse -Force
    }
    New-Item -ItemType Directory -Path $buildDir | Out-Null
    New-Item -ItemType Directory -Path "$buildDir\manifest\config" -Force | Out-Null
    New-Item -ItemType Directory -Path "$buildDir\manifest\docker" -Force | Out-Null
    New-Item -ItemType Directory -Path "$buildDir\manifest\sql" -Force | Out-Null
    # 创建Docker构建所需的目录
    New-Item -ItemType Directory -Path "$buildDir\temp\linux_amd64" -Force | Out-Null
}

# 构建Go项目
function Build-GoProject {
    Write-Info "开始构建Go项目(仅Linux版本)..."
    try {
        # 设置编译环境变量
        $env:GOOS = "linux"
        $env:GOARCH = "amd64"
        $env:CGO_ENABLED = "0"
        
        # 构建Linux版本
        go build -o "$buildDir\${projectName}" -ldflags="-s -w" main.go
        
        # 同时复制一份到Docker构建目录
        Copy-Item -Path "$buildDir\${projectName}" -Destination "$buildDir\temp\linux_amd64\main" -Force
        
        # 恢复环境变量
        $env:GOOS = "windows"
        $env:GOARCH = "amd64"
        $env:CGO_ENABLED = "1"
        
        Write-Info "Go项目构建完成"
    }
    catch {
        Write-Error "Go项目构建失败: $_"
        exit 1
    }
}

# 复制必要文件
function Copy-RequiredFiles {
    Write-Info "复制必要文件..."
    try {
        # 复制配置文件
        Copy-Item -Path "manifest\config\config.yaml" -Destination "$buildDir\manifest\config\config.yaml" -Force
        
        # 复制SQL文件
        if (Test-Path "manifest\sql") {
            Copy-Item -Path "manifest\sql\*" -Destination "$buildDir\manifest\sql\" -Recurse -Force
        }
        
        # 复制Docker部署文件
        Copy-Item -Path "manifest\docker\*" -Destination "$buildDir\manifest\docker\" -Recurse -Force
        
        # 复制资源文件
        Copy-Item -Path "resource" -Destination "$buildDir\resource" -Recurse -Force
        
        # 复制其他必要文件
        Copy-Item -Path "go.mod" -Destination "$buildDir\go.mod" -Force
        Copy-Item -Path "go.sum" -Destination "$buildDir\go.sum" -Force
        Copy-Item -Path "README.MD" -Destination "$buildDir\README.MD" -Force
        
        Write-Info "文件复制完成"
    }
    catch {
        Write-Error "文件复制失败: $_"
        exit 1
    }
}

# 创建生产环境配置文件
function Create-ProductionConfig {
    Write-Info "创建生产环境配置文件..."
    $configContent = @"
# 生产环境配置
server:
  address: ":8000"
  openapiPath: "/api.json"
  swaggerPath: "/swagger"

logger:
  level: "info"
  stdout: true
  file: "logs/server.log"

database:
  default:
    link: "mysql:go-jztc:BeAdwSKMbeiDpkHw@tcp(mysql:3306)/go-jztc?charset=utf8mb4&parseTime=true&loc=Local"
    debug: false
    prefix: ""
    maxIdle: 10
    maxOpen: 100
    maxLifetime: 30
"@

    $configContent | Out-File -FilePath "$buildDir\manifest\config\config.prod.yaml" -Encoding UTF8
    Write-Info "生产环境配置文件创建完成"
}

# 创建部署说明文档
function Create-DeploymentGuide {
    Write-Info "创建部署说明文档..."
    $readmeContent = @"
# 部署说明

## CentOS部署
1. 解压文件
2. 修改 manifest/config/config.yaml 配置文件
3. 给程序添加执行权限: chmod +x demo
4. 运行程序: ./demo

## Docker部署
1. 进入 manifest/docker 目录
2. 给脚本添加执行权限: chmod +x *.sh
3. 执行部署脚本: ./deploy.sh
4. 查看服务状态: ./status.sh

## 配置说明
- 数据库配置在 manifest/config/config.yaml 中
- 生产环境配置在 manifest/config/config.prod.yaml 中

## 常见问题
请参考 manifest/docker/README.md
"@

    $readmeContent | Out-File -FilePath "$buildDir\DEPLOY.md" -Encoding UTF8
    Write-Info "部署说明文档创建完成"
}

# 主函数
function Main {
    Write-Info "开始构建项目 $projectName v$version ($timestamp)..."
    
    # 检查环境
    Check-GoEnvironment
    
    # 创建构建目录
    Create-BuildDirectory
    
    # 构建Go项目
    Build-GoProject
    
    # 复制必要文件
    Copy-RequiredFiles
    
    # 创建生产环境配置文件
    Create-ProductionConfig
    
    # 创建部署说明文档
    Create-DeploymentGuide
    
    Write-Info "构建完成！请在 $buildDir 目录中查看构建结果"
}

# 执行主函数
Main 