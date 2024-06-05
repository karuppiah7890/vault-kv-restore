package main

import "encoding/json"

func fromJSON(JSONData []byte) (*VaultKvBackup, error) {
	vaultKvBackup := &VaultKvBackup{}
	err := json.Unmarshal(JSONData, vaultKvBackup)
	if err != nil {
		return nil, err
	}
	return vaultKvBackup, nil
}
