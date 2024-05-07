param (
    [Parameter(Mandatory=$true)]
    [int]$Select = 3 # Parameter to specify the number of times to execute the command
)

switch ($Select){
    0 {
        # 小小的批次
        .\run.ps1 -Args "--TotalDataSize 16000 --BatchSize 4000"
    }
    1{
        # 小小的批次 + 降低了Inject Speed
        .\run.ps1 -Args "--TotalDataSize 16000 --BatchSize 4000 --InjectSpeed 1000"
    }
    2{
        #  小批次，四个节点
        .\run.ps1 -Args "--TotalDataSize 16000 --BatchSize 4000" -N 4
    }
}