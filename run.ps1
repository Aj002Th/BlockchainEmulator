$prefix = Get-Date -Format "MM-dd-yyyy HH-mm"

$Env:BCEM_OUTPUT_PREFIX=$prefix

start cmd '/k .\blockchainEmulator.exe -n 1'

start cmd '/k .\blockchainEmulator.exe -n 2'

start cmd '/k .\blockchainEmulator.exe -c -f'

start cmd '/k .\blockchainEmulator.exe -n 0'

