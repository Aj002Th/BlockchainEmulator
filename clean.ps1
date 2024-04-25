# 用户输入 Y/n 决定是否运行脚本
$confirm = Read-Host "是否运行清理脚本？(输入 Y 确认，输入其他任意字符取消)"

if ($confirm -eq "Y") {
    # 创建一个列表dir_to_remove，内容是record, log, result这三个文件夹。
    $dir_to_remove = @("record", "log", "result")

    # 遍历列表
    foreach ($dir in $dir_to_remove) {
        $dir_path = Join-Path -Path "." -ChildPath $dir

        # 检测路径是否存在
        If (Test-Path -Path $dir_path) {
            # 如果存在，则删除文件夹
            Remove-Item -Path $dir_path -Force -Recurse
            Write-Host "文件夹 $dir_path 被删除。"
        } Else {
            Write-Host "文件夹 $dir_path 不存在，不用删除。"
        }
    }
} else {
    Write-Host "取消运行清理脚本。"
}
