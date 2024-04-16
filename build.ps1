Write-Host "现在在编译程序……"
go build -gcflags=all="-N -l"
if  ($LASTEXITCODE  -eq 0) {
    Write-Host "编译成功。"
}
else{
    Write-Host "编译失败，错误码: $LASTEXITCODE"
}
