package metadata_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/metadata"

	webapi2 "github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"

	"github.com/stretchr/testify/require"
)

func TestMetadataAPI_EntityDefinitions(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := newMetadataAPI(ctx, t).EntityDefinitions().Do(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Definitions)
}

func TestMetadataAPI_CustomEntityDefinitions(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := newMetadataAPI(ctx, t).CustomEntityDefinitions().Do(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Definitions)
}

func TestMetadataAPI_EntityDefinition(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	all, err := newMetadataAPI(ctx, t).EntityDefinitions().Do(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, all.Definitions)

	one, err := newMetadataAPI(ctx, t).EntityDefinition(all.Definitions[0].LogicalName).Do(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, one.LogicalName)
}

func TestMetadataAPI_EntityAttributes(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	all, err := newMetadataAPI(ctx, t).EntityDefinitions().Do(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, all.Definitions)

	one, err := newMetadataAPI(ctx, t).EntityDefinition(all.Definitions[0].LogicalName).Do(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, one.LogicalName)

	result, err := newMetadataAPI(ctx, t).EntityAttributes(one.LogicalName).Do(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, result.Attributes)
}

func newMetadataAPI(ctx context.Context, t *testing.T) *metadata.API {
	t.Helper()

	// Compose client config
	config := webapi2.ConfigFromENV()
	if config.IsEmpty() {
		t.Skip("missing ENVs: 'DEMO_MODEL_{TENANT_ID,CLIENT_ID,CLIENT_SECRET,API_HOST}'")
	}

	// Create OAuth2 client
	client, err := webapi2.NewClient(ctx, config, http.DefaultClient)
	require.NoError(t, err)

	return metadata.NewAPI(client)
}
