param(
    [String]$Action
)
$ProgressPreference='SilentlyContinue'
$ConfirmPreference ='None'

$Item = @{
    id = "1";
    outputs = @{
        a = 1;
        b = 2;
        c = 3;
    }
}

$Result = $Item | ConvertTo-JSON
switch($Action) {
   "create" { Write-Output $Result | Tee-Object $env:TF_DATA_FILE; break}
   "read" { Get-Content -Path $env:TF_DATA_FILE; break}
   "delete" { Write-Output $Result; break}
   "update" { Get-Content -Path $env:TF_DATA_FILE; break}
   "list" { List-O365Member; break}
}
