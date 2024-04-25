if( -not ($PSVersionTable.PSVersion.Major -ge 7 )){
    Write-Host "发现错误：" -ForegroundColor White -BackgroundColor Red
    Write-Host ""
    Write-Host "Powershell版本要求大于7.0。请先升级你的版本。程序现在退出。" -ForegroundColor White -BackgroundColor Red
    Write-Host ""
    exit
    Write-Output "After Exit"
}
$prefix = Get-Date -Format "MM-dd-yyyy-HH-mm-ss"

$Env:BCEM_OUTPUT_PREFIX=$prefix

Write-Host "现在启动若干个节点。输出前缀是时间戳。: $prefix" -ForegroundColor White -BackgroundColor Green
Write-Host ""


start cmd '/k .\blockchainEmulator.exe -n 1'

start cmd '/k .\blockchainEmulator.exe -n 2'

start cmd '/k .\blockchainEmulator.exe -c -f'

start cmd '/k .\blockchainEmulator.exe -n 0'

