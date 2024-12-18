package test_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	model "github.com/michaljurecko/ch-demo/internal/pkg/app/demo/model/gen"
	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/webapi"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	// Create client and repo
	config := webapi.ConfigFromENV()
	if config.IsEmpty() {
		t.Skip("missing ENVs: 'DEMO_MODEL_{TENANT_ID,CLIENT_ID,CLIENT_SECRET,API_HOST}'")
	}

	config.DebugRequest = true
	config.DebugResponse = true
	client, err := webapi.NewClient(ctx, config, http.DefaultClient)
	require.NoError(t, err)
	repo := model.NewRepository(client)

	// Player ----------------------------------------------------------------------------------------------------------

	// Create a player
	player1 := &model.Player{
		FirstName: "John",
		LastName:  "Test",
		Phone:     "123456789",
		Email:     "john.test@company.com",
		IC:        "123456789",
		VATID:     "123456789",
		Address:   "Test street 123",
	}
	_, err = repo.Player().Create(player1).Do(ctx)
	require.NoError(t, err)

	// Delete player after test
	t.Cleanup(func() {
		_, err = repo.Player().Delete(player1.ID).Do(ctx)
		assert.NoError(t, err)
	})

	// Update player
	player1.TrackChanges()
	player1.LastName = "Test123"
	_, err = repo.Player().Update(player1).Do(ctx)
	require.NoError(t, err)

	// Get player by ID
	player1Clone, err := repo.Player().ByID(player1.ID).Do(ctx)
	assert.Equal(t, player1, player1Clone)

	// Character -------------------------------------------------------------------------------------------------------

	// Load class, there is no alternate key for the name field, and we shouldn't modify the schema.
	var class *model.Class
	classes, err := repo.Class().All().Do(ctx)
	require.NoError(t, err)
	for _, c := range classes.Items {
		if c.Name == "Warrior" {
			class = c
			break
		}
	}
	require.NotNil(t, class)

	// Load race, there is no alternate key for the name field, and we shouldn't modify the schema.
	var race *model.Race
	races, err := repo.Race().All().Do(ctx)
	require.NoError(t, err)
	for _, r := range races.Items {
		if r.Name == "Human" {
			race = r
			break
		}
	}
	require.NotNil(t, race)

	// Create atomically: 1. character, 2. dice roll,  use change set
	character1 := &model.Character{
		Name:         "Test Character",
		Level:        3,
		Strength:     10,
		Dexterity:    2,
		Intelligence: 5,
		Charisma:     4,
		Class:        webapi.LookupValue[model.Class](class.ID),
		Race:         webapi.LookupValue[model.Race](race.ID),
		Player:       webapi.LookupValue[model.Player](player1.ID),
	}
	diceRoll1 := &model.DiceRoll{
		DiceRoll:   "strength",
		RollResult: 4,
		RollDate:   time.Now(),
		ActionType: "test",
	}
	changes := repo.NewChangeSet()
	characterRef := changes.Add(repo.Character().Create(character1))
	diceRoll1.Character.SetContentID(characterRef)
	changes.Add(repo.DiceRoll().Create(diceRoll1))
	require.NoError(t, changes.Do(ctx))

	// Delete character and dice roll after test
	t.Cleanup(func() {
		_, err = repo.Character().Delete(character1.ID).Do(ctx)
		assert.NoError(t, err)
		_, err = repo.DiceRoll().Delete(diceRoll1.ID).Do(ctx)
		assert.NoError(t, err)
	})
}
