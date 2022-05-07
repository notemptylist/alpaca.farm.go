<#
Requests token from vault by authenticating via user name and password
Uses the token to query the v1/alpaca/data/keys path
Parses the json, selects data, converts back to json, stores to file using 'set-content'
#>

$outputfile = ".env.json"
$keypath = "/v1/alpaca/data/keys"

$vault_base_url = Read-Host -Prompt "Enter vault server address "
$login_url = $vault_base_url + "/v1/auth/userpass/login/"
$user = Read-Host -Prompt "Enter username for $vault_base_url "
$login_url = $login_url + $user
$pass = Read-Host -Prompt "Enter password " -AsSecureString
$pass = [System.Net.NetworkCredential]::new("", $pass).Password

$postParams = @{"password"=$pass} | ConvertTo-Json
$resp = Invoke-WebRequest -Uri "$login_url" -Method "POST" -ContentType "application/json" -Body $postParams
$token = $(ConvertFrom-Json -InputObject $resp.Content).auth.client_token

$vault_kv_url = $vault_base_url + $keypath
$headers = @{"X-Vault-Token" = $token }
$resp = Invoke-WebRequest -Uri "$vault_kv_url" -Method "GET" -ContentType "application/json" -Headers $headers
$($resp.Content | ConvertFrom-Json).data.data | ConvertTo-Json | Set-Content $outputfile

Write-Host "Written keys to $outputfile"