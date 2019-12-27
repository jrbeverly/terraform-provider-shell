param(
    [String]$Action
)
$ProgressPreference='SilentlyContinue'
$ConfirmPreference ='None'

$Item = @{
    id = "1"
}
$Result = $Item | ConvertTo-JSON
Write-Output $Result | Tee-Object here.txt
# switch($Action) {
#    "create" { Write-Output $Result; break}
#    "read" { Get-O365Member; break}
#    "delete" { Remove-O365Member; break}
#    "update" { Set-O365Member; break}
#    "list" { List-O365Member; break}
# }
