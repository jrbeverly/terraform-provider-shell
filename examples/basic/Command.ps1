param(
    [String]$Action
)
$ErrorActionPreference = "Stop"
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
   "create" { Set-Content -Path $env:TF_DATA_FILE -Value $Result; break}
   "read" { Set-Content -Path $env:TF_DATA_FILE -Value $Result; break}
   "delete" { Set-Content -Path $env:TF_DATA_FILE -Value ""; break}
   "update" { Set-Content -Path $env:TF_DATA_FILE -Value $Result; break}
}
