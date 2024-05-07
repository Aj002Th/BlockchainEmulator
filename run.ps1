param (
    [Parameter(Mandatory=$false)]
    [int]$N = 3, # Parameter to specify the number of times to execute the command

    [Parameter(Mandatory=$false)]
    [string] $Args, # Parameter to specify the arguments to pass to the command

    [Parameter(Mandatory=$false)]
    [switch] $Help # Parameter to specify the arguments to pass to the command
)

if ($Help) {
    Write-Host "启动脚本帮助"  -ForegroundColor White -BackgroundColor Green
    Write-Host ""
    Write-Host "The Blockchain Emulator (BCEM) is a tool for testing blockchain applications."
    Write-Host ""
    Write-Host "Usage: run.ps1 [-N <number>] [-Args <arguments>] [-Help]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -N <number>      Number of times to execute the command (default: 3)"
    Write-Host " -Args <arguments> Arguments to pass to the command"  
    Write-Host ""
    Write-Host "节点程序参数帮助"  -ForegroundColor White -BackgroundColor Green
    Write-Host ""
    .\blockchainEmulator.exe --help
    exit 0
}

# Check if $N is a valid positive integer
if ($N -le 2) {
    Write-Host "Error: 参数N错误。PBFT节点数量至少大于等于3"
    exit 1
}

$prefix = Get-Date -Format "MM-dd-yyyy-HH-mm-ss"

$Env:BCEM_OUTPUT_PREFIX=$prefix

Write-Host "现在准备启动若干个节点。输出前缀是时间戳 `"$prefix`"" -ForegroundColor White -BackgroundColor Green
Write-Host ""

# 检测文件是否存在
if (Test-Path ".\blockchainEmulator.exe") {
    Write-Host "构建产物blockchainEmulator.exe存在。"-ForegroundColor White -BackgroundColor Green
    Write-Host ""
} else {
    Write-Host "构建产物blockchainEmulator.exe不存在。请你先用go build构建。"  -ForegroundColor White -BackgroundColor Red
    Write-Host ""
    exit
}

# 获取构建产物的最后修改日期
$lastModified = (Get-Item ".\blockchainEmulator.exe").LastWriteTime

# 计算距离现在的时间间隔
$timeDifference = New-TimeSpan -Start $lastModified -End (Get-Date)


Write-Host "脚本正在获取解释的参数"
Write-Host ""
powershell .\blockchainEmulator.exe "--interpret -N $N $Args"

# 打印输出时间间隔
Write-Host "TIPS：这个二进制副本生成的时候，距离现在已经过去 $($timeDifference.Hours) 小时 $($timeDifference.Minutes) 分钟 $($timeDifference.Seconds) 秒。"
Write-Host ""

# 检测用户输入是否以N开头，如果是则退出
$userInput = Read-Host "按回车键开始执行，或输入以N开头的字符退出。"

if ($userInput -like "N*") {
    Write-Host "输入了N开头的，取消执行，现在退出。"
    exit
}

# Loop to execute the command N times
for ($i = 1; $i -lt $N; $i++) {
    # Execute the command (replace "print n" with your desired command)
    Write-Host "启动PBFT节点$i"
    start cmd ("/k .\blockchainEmulator.exe -n $i" + " "+ $Args)
    # Invoke the command here
}

Write-Host "启动Supervisor"
start cmd ('/k .\blockchainEmulator.exe -c -f' + " "+ $Args)

Write-Host "启动PBFT主节点"
# Write-Host ('/k .\blockchainEmulator.exe -n 0' + " "+$Args)
# Write-Host $Args
start cmd ('/k .\blockchainEmulator.exe -n 0' + " "+ $Args)

Write-Host "已启动若干节点"
