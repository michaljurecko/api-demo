package player_test

import (
	"connectrpc.com/connect"
	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	"github.com/michaljurecko/ch-demo/api/gen/go/demo/v1/apiconnect"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/service/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"testing"
)

func TestAggregation(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aggregation")
}

var _ = Describe("Aggregation Operations", Ordered, func() {
	// Run server on a random port and create generated client
	var svcClient apiconnect.ApiServiceClient
	BeforeAll(func(ctx SpecContext) {
		baseURL := test.StartTestServer(ctx)
		svcClient = apiconnect.NewApiServiceClient(http.DefaultClient, baseURL)
	})

	var playersInitialCount int
	It("can list players and characters", func(ctx SpecContext) {
		response, err := svcClient.ListPlayersAndCharacters(ctx, connect.NewRequest(&emptypb.Empty{}))
		Expect(err).NotTo(HaveOccurred())
		playersInitialCount = len(response.Msg.Players)
	})

	Context("a new player and character is created", func() {
		// Create a player
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

		// Create a character
		var character *api.Character
		BeforeAll(func(ctx SpecContext) {
			// Load warrior class
			var warriorClass *api.Class
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

			// Load human race
			var humanRace *api.Race
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

			// Create a character
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
		})
		AfterAll(func(ctx SpecContext) {
			_, err := svcClient.DeleteCharacter(ctx, connect.NewRequest(
				&api.DeleteCharacterRequest{
					Id: character.GetId(),
				},
			))
			Expect(err).NotTo(HaveOccurred())
		})

		It("response should contains new player and character, cache is invalidated", func(ctx SpecContext) {
			response, err := svcClient.ListPlayersAndCharacters(ctx, connect.NewRequest(&emptypb.Empty{}))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(response.Msg.Players)).To(Equal(playersInitialCount + 1))

			By("get player response")
			var subResponse *api.CharactersPerPlayer
			for _, r := range response.Msg.Players {
				if r.Player.Id == player.Id {
					subResponse = r
					break
				}
			}
			Expect(subResponse).NotTo(BeNil())
			Expect(len(subResponse.Characters)).To(Equal(1))
			Expect(subResponse.Characters[0].Id).To(Equal(character.Id))
		})
	})
})
