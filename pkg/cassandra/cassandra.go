// Copyright (c) WarnerMedia Direct, LLC. All rights reserved. Licensed under the MIT license.
// See the LICENSE file for license information.
package cassandra

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/HBOCodeLabs/hurley-kit/secrets"
)

func GetVaultSecrets() {
	ttl := secrets.CacheTTL(time.Duration(3600) * time.Second)
	vaultAddress := secrets.VaultAddress("https://vault.api.hbo.com")
	vaultTimeout := secrets.VaultTimeout(time.Duration(3600) * time.Second)
	vaultMaxRetries := secrets.VaultMaxRetries(3)
	appRole := secrets.AppRole("beta-demo")

	vaultStore, err := secrets.NewVaultStore(ttl, vaultAddress, vaultTimeout, vaultMaxRetries, appRole)
	if err != nil {
		t.Errorf("Failed to create secret store. Error %s", err)
		t.Fail()
	}
	//	assertion = false
	byts, err := vaultStore.Get("dre/service/all_rds/generic_read")
	if err != nil {
		return "false", err
	}

	fmt.Println("secret is", byts)

}

func fetchSecret(store *secrets.VaultStore, path string) (string, error) {
	secretBytes, err := store.Get(path)
	if err != nil {
		return "", err
	}

	var secretMap map[string]interface{}
	err = json.Unmarshal(secretBytes, &secretMap)
	if err != nil {
		return "", err
	}

	return secretMap["value"].(string), nil
}
