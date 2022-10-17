### Secrets management using Hashicorp Vault

To avoid hardcoding the API keys or leaving around dot files, the two scripts, `initkeys.ps1` and `initkeys.sh` are used to manage local copies of secret files that are required to talk the ALPACA API.

```
PS C:\src\alpaca.farm.go>.\initkeys.ps1
Enter vault server address : http://example.com:8200
Enter username for http://example.com:8200 : foo
Enter password : ******
Written keys to .env.json
```

```
./initkeys.sh
Enter vault server address : http://example.com:8200
Enter username for http://example.com:8200 : foo
Enter password : Written keys to .env.json
```

The scripts use a hardcoded path to abstract the location of the keys within the vault's `secret` module. 

It is also possible to call the vault API from the Go code.
