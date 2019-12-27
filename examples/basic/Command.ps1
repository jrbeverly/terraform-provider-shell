param(
    [String]$Action
)
$ProgressPreference='SilentlyContinue'
$ConfirmPreference ='None'

$Item = @{
    id = "1"
}
$Result = $Item | ConvertTo-JSON
switch($Action) {
   "create" { Write-Output $Result | Tee-Object here.txt; break}
   "read" { Get-Content -Path here.txt; break}
   "delete" { Write-Output $Result; break}
   "update" { Get-Content -Path here.txt; break}
   "list" { List-O365Member; break}
}
