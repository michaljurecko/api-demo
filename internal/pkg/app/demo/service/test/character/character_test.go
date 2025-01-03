package character_test

import (
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	"github.com/michaljurecko/ch-demo/api/gen/go/demo/v1/apiconnect"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/service/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestCharacter(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Character")
}

var _ = Describe("Character Operations", Ordered, func() {
	// Run server on a random port and create generated client
	var svcClient apiconnect.ApiServiceClient
	BeforeAll(func(ctx SpecContext) {
		baseURL := test.StartTestServer(ctx)
		svcClient = apiconnect.NewApiServiceClient(http.DefaultClient, baseURL)
	})

	// Create a player - character parent entity
	var player *api.Player
	BeforeAll(func(ctx SpecContext) {
		response, err := svcClient.CreatePlayer(ctx, connect.NewRequest(
			&api.CreatePlayerRequest{
				FirstName: "John",
				LastName:  "Brown",
				Phone:     "+420123456",
				Email:     "john.brown@game.com",
				Ic:        "00023337",
			},
		))
		Expect(err).NotTo(HaveOccurred())
		player = response.Msg
	})
	AfterAll(func(ctx SpecContext) {
		_, err := svcClient.DeletePlayer(ctx, connect.NewRequest(
			&api.DeletePlayerRequest{
				Id: player.GetId(),
			},
		))
		Expect(err).NotTo(HaveOccurred())
	})

	// Clear character if test fails and the entity hasn't been deleted
	var character *api.Character
	AfterAll(func(ctx SpecContext) {
		if character != nil {
			_, err := svcClient.DeleteCharacter(ctx, connect.NewRequest(
				&api.DeleteCharacterRequest{
					Id: player.GetId(),
				},
			))
			Expect(err).NotTo(HaveOccurred())
		}
	})

	// Load warrior class
	var warriorClass *api.Class
	BeforeAll(func(ctx SpecContext) {
		classesResp, err := svcClient.ListClasses(ctx, connect.NewRequest(&emptypb.Empty{}))
		Expect(err).NotTo(HaveOccurred())
		classes := classesResp.Msg.GetClasses()
		Expect(classes).NotTo(BeEmpty())

		for _, class := range classes {
			if class.GetName() == "Warrior" {
				warriorClass = class
				break
			}
		}
		Expect(warriorClass).NotTo(BeNil())
	})

	// Load human race
	var humanRace *api.Race
	BeforeAll(func(ctx SpecContext) {
		racesResp, err := svcClient.ListRaces(ctx, connect.NewRequest(&emptypb.Empty{}))
		Expect(err).NotTo(HaveOccurred())
		races := racesResp.Msg.GetRaces()
		Expect(races).NotTo(BeEmpty())

		for _, race := range races {
			if race.GetName() == "Human" {
				humanRace = race
				break
			}
		}
		Expect(humanRace).NotTo(BeNil())
	})

	It("cannot create character with empty request", func(ctx SpecContext) {
		_, err := svcClient.CreateCharacter(ctx, connect.NewRequest(
			&api.CreateCharacterRequest{},
		))
		Expect(err).To(MatchError(
			ContainSubstring("name: value length must be at least 1 characters [string.min_len]"),
			ContainSubstring("strength: value must be greater than or equal to 1 and less than or equal to 20 [int32.gte_lte]"),
			ContainSubstring("dexterity: value must be greater than or equal to 1 and less than or equal to 20 [int32.gte_lte]"),
			ContainSubstring("intelligence: value must be greater than or equal to 1 and less than or equal to 20 [int32.gte_lte]"),
			ContainSubstring("charisma: value must be greater than or equal to 1 and less than or equal to 20 [int32.gte_lte]"),
			ContainSubstring("class_id: value length must be at least 1 characters [string.min_len]"),
			ContainSubstring("race_id: value length must be at least 1 characters [string.min_len]"),
			ContainSubstring("player_id: value length must be at least 1 characters [string.min_len]"),
		))
	})

	It("can create character", func(ctx SpecContext) {
		response, err := svcClient.CreateCharacter(ctx, connect.NewRequest(
			&api.CreateCharacterRequest{
				Name:         "Boromir",
				Strength:     1,
				Dexterity:    2,
				Intelligence: 3,
				Charisma:     4,
				ClassId:      warriorClass.GetId(),
				RaceId:       humanRace.GetId(),
				PlayerId:     player.GetId(),
			},
		))
		Expect(err).NotTo(HaveOccurred())

		character = response.Msg
		Expect(character).Should(PointTo(MatchFields(IgnoreExtras, Fields{
			"Name":         Equal("Boromir"),
			"Strength":     BeNumerically(">=", 1+warriorClass.GetStrengthBase()+humanRace.GetStrengthBase()),
			"Dexterity":    BeNumerically(">=", 2+warriorClass.GetDexterityBase()+humanRace.GetDexterityBase()),
			"Intelligence": BeNumerically(">=", 3+warriorClass.GetIntelligenceBase()+humanRace.GetIntelligenceBase()),
			"Charisma":     BeNumerically(">=", 4+warriorClass.GetCharismaBase()+humanRace.GetCharismaBase()),
			"ClassId":      Equal(warriorClass.GetId()),
			"RaceId":       Equal(humanRace.GetId()),
			"PlayerId":     Equal(player.GetId()),
		})))
	})

	It("can update character", func(ctx SpecContext) {
		response, err := svcClient.UpdateCharacter(ctx, connect.NewRequest(
			&api.UpdateCharacterRequest{
				Id:           character.Id,
				UpdateMask:   &fieldmaskpb.FieldMask{Paths: []string{"name", "intelligence"}},
				Name:         "Faramir",
				Intelligence: 20,
			},
		))
		Expect(err).NotTo(HaveOccurred())

		character = response.Msg
		Expect(character).Should(PointTo(MatchFields(IgnoreExtras, Fields{
			"Name":         Equal("Faramir"),
			"Intelligence": BeNumerically(">=", 20+warriorClass.GetIntelligenceBase()+humanRace.GetIntelligenceBase()),
		})))
	})

	It("get character by id", func(ctx SpecContext) {
		response, err := svcClient.GetCharacter(ctx, connect.NewRequest(
			&api.GetCharacterRequest{
				Id: character.Id,
			},
		))
		Expect(err).NotTo(HaveOccurred())
		Expect(response.Msg).Should(PointTo(MatchFields(IgnoreExtras, Fields{
			"Name":         Equal("Faramir"),
			"Strength":     BeNumerically(">=", 1+warriorClass.GetStrengthBase()+humanRace.GetStrengthBase()),
			"Dexterity":    BeNumerically(">=", 2+warriorClass.GetDexterityBase()+humanRace.GetDexterityBase()),
			"Intelligence": BeNumerically(">=", 20+warriorClass.GetIntelligenceBase()+humanRace.GetIntelligenceBase()),
			"Charisma":     BeNumerically(">=", 4+warriorClass.GetCharismaBase()+humanRace.GetCharismaBase()),
			"ClassId":      Equal(warriorClass.GetId()),
			"RaceId":       Equal(humanRace.GetId()),
			"PlayerId":     Equal(player.GetId()),
		})))
	})

	It("delete character", func(ctx SpecContext) {
		_, err := svcClient.DeleteCharacter(ctx, connect.NewRequest(
			&api.DeleteCharacterRequest{
				Id: character.Id,
			},
		))
		Expect(err).NotTo(HaveOccurred())
		character = nil
	})
})
