package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

type VaultKvBackup struct {
	Secrets map[string]map[string]interface{} `json:"secrets"`
}

func convertJSONToVaultKvBackup(JSONData []byte) (*VaultKvBackup, error) {
	vaultKvBackup, err := fromJSON(JSONData)
	if err != nil {
		return nil, err
	}
	return vaultKvBackup, nil
}

func restoreVaultKvSecrets(client *api.Client, kvMountPath string, vaultKvBackup *VaultKvBackup, quietProgress bool) error {
	kvV2Client := client.KVv2(kvMountPath)

	for secretsPath, secrets := range vaultKvBackup.Secrets {
		if quietProgress {
			fmt.Fprintf(os.Stdout, ".")
		} else {
			fmt.Fprintf(os.Stdout, "\nrestoring secrets to `%s` secrets path in vault\n", secretsPath)
		}

		restoredKvSecret, err := kvV2Client.Put(context.TODO(), secretsPath, secrets)

		if err != nil {
			return fmt.Errorf("error occurred while putting/writing the secrets at path `%s` in vault: %v", secretsPath, err)
		}

		if restoredKvSecret == nil {
			return fmt.Errorf("no secret at path `%s` in vault after write operation", secretsPath)
		}
	}

	return nil
}
