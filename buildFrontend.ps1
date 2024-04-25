# 先同步一下子模块
git submodule update

# 切换工作目录到bcEmSupMonitor
Set-Location "bcEmSupMonitor"

# 执行 npm install
npm install

# 检查是否有错误，如果有则停止并报错
if ($LastExitCode -ne 0) {
    Write-Error "npm install 失败"
    exit $LastExitCode
}

# 删除原来的前端部分

$dir_path =  "../web/out"

If (Test-Path -Path $dir_path) {
    # 如果存在，则删除文件夹
    Remove-Item -Path $dir_path -Force -Recurse
    Write-Host "文件夹 $dir_path 被删除。"
} Else {
    Write-Host "文件夹 $dir_path 不存在，不用删除。"
}

# 执行 next build 等等
next build

# 检查是否有错误，如果有则停止并报错
if ($LastExitCode -ne 0) {
    Write-Error "next build 失败"
    exit $LastExitCode
}

# 将输出目录out拷贝到原工作目录的web/out
Copy-Item -Path "out" -Destination "../web/out" -Recurse
