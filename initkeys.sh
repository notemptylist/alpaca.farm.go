#!/bin/bash
outputfile=.env.json
read -p "Enter vault server address : " vault_base_url
read -p "Enter username for $vault_base_url : " user
read -sp "Enter password : " pass

resp="$(curl -s -H 'Content-Type: application/json' -d "{\"password\": \"$pass\"}" -X POST $vault_base_url/v1/auth/userpass/login/$user)"
client_token="$(echo $resp | awk -F 'client_token":"' '{print $2}' | awk -F '"' '{print $1}')"
data="$(curl -s -H "X-Vault-Token: $client_token" --request GET $vault_base_url/v1/alpaca/data/keys )"
echo $data  | sed -n 's/.*"data":\(.*\),"metadata".*/\1/p' > $outputfile 
if [ $? -eq 0 ]; then
    echo "Written keys to $outputfile"
fi
