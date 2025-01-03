package ares_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/ares"
	"github.com/stretchr/testify/assert"
)

func TestClient_ByIC(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := ares.NewClient(http.DefaultClient)

	// Not found
	_, err := client.ByIC(ctx, "12345678")
	assert.NotEqual(t, ares.EntityNotFoundError{IC: "12345678"}, err)

	// Found
	result, err := client.ByIC(ctx, "00023337")
	if assert.NoError(t, err) {
		assert.Equal(t, "Národní divadlo", result.Name)
		assert.Equal(t, "00023337", result.IC)
		assert.Equal(t, "CZ00023337", result.VatID)
		assert.Equal(t, "Ostrovní 225/1\nNové Město\n11000 Praha 1", result.Address.String())
	}
}
