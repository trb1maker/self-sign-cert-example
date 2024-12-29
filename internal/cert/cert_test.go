package cert

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrivateMethods(t *testing.T) {
	g := new(Generator)

	t.Run("createPrivrateKey", func(t *testing.T) {
		require.NoError(t, g.createPrivrateKey())
		require.NotNil(t, g.privateKey)
	})

	t.Run("createCertificate", func(t *testing.T) {
		require.NoError(t, g.createCertificate("pravo.tech"))
		require.NotNil(t, g.certificate)
		require.Equal(t, []string{"pravo.tech", "localhost"}, g.certificate.DNSNames)
	})

	t.Run("store", func(t *testing.T) {
		require.NoError(t, g.store("."))

		require.FileExists(t, "./pravo.tech.key")
		require.FileExists(t, "./pravo.tech.pem")

		require.NoError(t, os.Remove("./pravo.tech.key"))
		require.NoError(t, os.Remove("./pravo.tech.pem"))
	})
}

func TestCreateCertificate(t *testing.T) {
	g := new(Generator)
	require.NoError(t, g.CreateCertificate(".", "pravo.tech"))
	require.NoError(t, os.Remove("./pravo.tech.key"))
	require.NoError(t, os.Remove("./pravo.tech.pem"))
}
