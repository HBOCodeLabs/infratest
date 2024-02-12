package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/vault/api"
	"github.com/hbocodelabs/infratest/pkg/vault"
	"github.com/stretchr/testify/require"
)

func getVaultDockerName(version string, uniqueID string) (name string) {
	return fmt.Sprintf("vault-%s-%s", version, uniqueID)
}

func runVaultDocker(t *testing.T, version string, uniqueID string) (rootToken string, containerName string) {
	containerName = getVaultDockerName(version, uniqueID)
	containerImage := fmt.Sprintf("vault:%s", version)
	rootToken = random.UniqueId()
	rootTokenEnvString := fmt.Sprintf("VAULT_DEV_ROOT_TOKEN_ID=%s", rootToken)
	runOpts := &docker.RunOptions{
		Detach: true,
		Name:   containerName,
		EnvironmentVariables: []string{
			rootTokenEnvString,
		},
		OtherOptions: []string{
			"--cap-add=IPC_LOCK",
			"-P",
		},
		Remove: true,
	}
	docker.Run(t, containerImage, runOpts)
	return
}

func getVaultDockerPort(t *testing.T, containerName string) (port uint16) {
	inspect := docker.Inspect(t, containerName)
	port = inspect.GetExposedHostPort(8200)
	require.NotZero(t, port, "Vault Docker container is not exposing the required port.")
	return
}

func stopVaultDocker(t *testing.T, name string) {
	containers := []string{
		name,
	}
	stopOpts := &docker.StopOptions{}
	docker.Stop(t, containers, stopOpts)
}

// This is added because tests in this suite do not currently support parallel execution
// due to the dependency on Docker containers. To be fixed in a future release!
//
//nolint:paralleltest
func TestVault(t *testing.T) {
	vaultVersion := os.Getenv("VAULT_VERSION")
	uniqueID := random.UniqueId()
	require.NotEmpty(t, vaultVersion, "Vault version must be specified.")

	rootToken, containerName := runVaultDocker(t, vaultVersion, uniqueID)
	defer stopVaultDocker(t, containerName)
	port := getVaultDockerPort(t, containerName)

	vaultAddress := fmt.Sprintf("http://localhost:%d", port)
	clientConfig := &api.Config{
		Address:    vaultAddress,
		MaxRetries: 100,
	}
	client, err := api.NewClient(clientConfig)
	require.Nil(t, err, "Vault NewClient method returned an unexpected error.")
	client.SetToken(rootToken)

	logicalClient := client.Logical()
	expectedPath := "secret/data/hello"
	expectedSecretData := map[string]interface{}{
		"data": map[string]interface{}{
			"username": "myname",
			"password": "password",
		},
	}
	_, err = logicalClient.Write(expectedPath, expectedSecretData)
	require.Nil(t, err, "Vault Write method returned an unexpected error.")

	ctx := context.TODO()

	vault.AssertSecretExists(ctx, t, vault.WithLogicalClient(logicalClient), vault.WithPath(expectedPath), vault.WithKey("username"), vault.WithValue("myname"))
	vault.AssertSecretExists(ctx, t, vault.WithLogicalClient(logicalClient), vault.WithPath(expectedPath), vault.WithKey("password"), vault.WithValue("password"))

}
