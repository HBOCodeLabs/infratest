package k8s

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetKubeconfigPathE_NoPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	require.Nil(t, err)
	expectedPath := filepath.Join(homeDir, ".kube", "config")

	actualPath, err := getKubeconfigPathE("")

	require.Nil(t, err)
	require.Equal(t, expectedPath, actualPath)
}

func TestGetKubeconfigPathE_Path(t *testing.T) {
	expectedPath := filepath.Join("/tmp", ".kube", "config")

	actualPath, err := getKubeconfigPathE(expectedPath)

	require.Nil(t, err)
	require.Equal(t, expectedPath, actualPath)
}
