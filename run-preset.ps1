param (
    [Parameter(Mandatory=$true)]
    [int]$Select # 选择用哪一组预设值来启动模拟器。
)

switch ($Select){
    0 {
        # 小小的批次
        .\run.ps1 -Args "--TotalDataSize 160000 --BatchSize 4000 --InjectSpeed 1000" -N 4
    }
    1{
        # 小小的批次 + 降低了Inject Speed
        .\run.ps1 -Args "--TotalDataSize 160000 --BatchSize 4000 --InjectSpeed 200" -N 4
    }
    2{
        #  小批次，四个节点
        .\run.ps1 -Args "--TotalDataSize 160000 --BatchSize 2000 --InjectSpeed 1000" -N 4
    }
}