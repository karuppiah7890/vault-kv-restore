package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

var usage = `
usage: vault-kv-restore [-quiet|--quiet] [-file|-file <vault-kv-backup-json-file-path>] <kv-mount-path>

Note that the flags MUST come before the arguments

arguments of vault-kv-restore:

  <kv-mount-path> string
    vault kv v2 secrets engine mount path where you want to
    restore the secrets to

flags of vault-kv-restore:

  -file / --file string
      vault kv backup json file path (default "vault_kv_backup.json")

  -quiet / --quiet
      quiet progress (default false).
      By default vault-kv-restore CLI will show a lot of details
      about the restore process and detailed progress during the
      restore process

  -h / -help / --help
      show help

examples:

# show help
vault-kv-restore -h

# restores all vault kv v2 secrets from the kv backup JSON file.
# also, any existing vault kv v2 secrets with the same secret
# path in the kv backup JSON file will be overwritten.
vault-kv-restore -file <path-to-vault-kv-backup-json-file> <kv-mount-path>

# OR you can use --file too instead of -file

vault-kv-restore --file <path-to-vault-kv-backup-json-file> <kv-mount-path>

# quietly restore all vault kv v2 secrets.
# this will just show dots (.) for progress
vault-kv-restore -quiet -file <path-to-vault-kv-backup-json-file> <kv-mount-path>

# OR you can use --quiet too instead of -quiet

vault-kv-restore --quiet --file <path-to-vault-kv-backup-json-file> <kv-mount-path>
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "%s", usage)
	}
	quietProgress := flag.Bool("quiet", false, "quiet progress")
	vaultKvBackupJsonFilePath := flag.String("file", "vault_kv_backup.json", "vault kv backup json file path")
	flag.Parse()

  if !(flag.NArg() == 1) {
		fmt.Fprintf(os.Stderr, "invalid number of arguments: %d. expected 1 argument.\n\n", flag.NArg())
		flag.Usage()
	}

	// TODO: Take a backup of the kv secrets in case it is already present,
	// regardless of if they (any or all) have the same secret path as the
	// secrets in the kv v2 secrets backup json file. So that we have a
	// backup just in case, especially before overwriting.

	// Question: Should we do a complete backup before doing kv secret restore one by one?
	// Or should we do a backup of each kv secret one by one? When we restore them one by one
	// that is. Basically - take a backup one by one, or take a complete backup of
	// all kv secrets?

	config := api.DefaultConfig()
	client, err := api.NewClient(config)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating vault client: %s\n", err)
		os.Exit(1)
	}

	vaultKvBackupJsonFileContent, err := readFile(*vaultKvBackupJsonFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading vault kv backup json file: %s\n", err)
		os.Exit(1)
	}

	vaultKvBackup, err := convertJSONToVaultKvBackup(vaultKvBackupJsonFileContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing vault kv backup json file content: %s\n", err)
		os.Exit(1)
	}

  kvMountPath := flag.Args()[0]

	err = restoreVaultKvSecrets(client, kvMountPath, vaultKvBackup, *quietProgress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error restoring vault policies: %s\n", err)
		os.Exit(1)
	}
}
