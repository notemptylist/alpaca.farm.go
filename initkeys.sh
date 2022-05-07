#!/bin/bash
outputfile=.env.json
read -p "Enter vault server address : " vault_base_url
read -p "Enter username for $vault_base_url : " user
read -sp "Enter password : " pass

resp="$(curl -s -H 'Content-Type: application/json' -d "{\"password\": \"$pass\"}" -X POST $vault_base_url/v1/auth/userpass/login/$user)"
client_token="$(echo $resp | grep -Po '"client_token":"\K([^,"])+')"
data="$(curl -s -H "X-Vault-Token: $client_token" --request GET $vault_base_url/v1/alpaca/data/keys )"
echo $data | grep -Po '[^meta]data":{"data":\K{(.+?})' > $outputfile
if [ $? -eq 0 ]; then
    echo "Written keys to $outputfile"
fi
