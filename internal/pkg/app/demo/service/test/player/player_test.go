package player_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	api "github.com/michaljurecko/ch-demo/api/gen/go/demo/v1"
	"github.com/michaljurecko/ch-demo/api/gen/go/demo/v1/apiconnect"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/service/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func TestPlayer(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Player")
}

var _ = Describe("Player Operations", Ordered, func() {
	// Run server on a random port and create generated client
	var svcClient apiconnect.ApiServiceClient
	BeforeAll(func(ctx SpecContext) {
		baseURL := test.StartTestServer(ctx)
		svcClient = apiconnect.NewApiServiceClient(http.DefaultClient, baseURL)
	})

	// Clear player if test fails and the entity hasn't been deleted
	var player *api.Player
	AfterAll(func(ctx SpecContext) {
		if player != nil {
			_, err := svcClient.DeletePlayer(ctx, connect.NewRequest(
				&api.DeletePlayerRequest{
					Id: player.GetId(),
				},
			))
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("cannot create player with empty request", func(ctx SpecContext) {
		_, err := svcClient.CreatePlayer(ctx, connect.NewRequest(
			&api.CreatePlayerRequest{},
		))
		Expect(err).To(MatchError(
			ContainSubstring("first_name: value length must be at least 1 characters [string.min_len]"),
			ContainSubstring("last_name: value length must be at least 1 characters [string.min_len]"),
			ContainSubstring("phone: phone number must be in international format, e.g. '+421905123456'. [phone]"),
			ContainSubstring("email: value is empty, which is not a valid email address [string.email_empty]"),
			ContainSubstring("ic: value is not valid IČO. [ic]"),
		))
	})

	It("cannot create player with invalid IC", func(ctx SpecContext) {
		_, err := svcClient.CreatePlayer(ctx, connect.NewRequest(
			&api.CreatePlayerRequest{
				FirstName: "John",
				LastName:  "Brown",
				Phone:     "+420123456",
				Email:     "john.brown@game.com",
				Ic:        "12345678",
			},
		))
		Expect(err).To(MatchError(ContainSubstring("ic: value is not valid IČO")))
	})

	It("can create player with valid IC", func(ctx SpecContext) {
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
		Expect(player).Should(PointTo(MatchFields(IgnoreExtras, Fields{
			"FirstName": Equal("John"),
			"LastName":  Equal("Brown"),
			"Phone":     Equal("+420123456"),
			"Email":     Equal("john.brown@game.com"),
			"Ic":        Equal("00023337"),
			"VatId":     Equal("CZ00023337"),
			"Address":   Equal("Ostrovní 225/1\nNové Město\n11000 Praha 1"),
		})))
	})

	It("cannot update player if not exists", func(ctx SpecContext) {
		_, err := svcClient.UpdatePlayer(ctx, connect.NewRequest(
			&api.UpdatePlayerRequest{
				Id:         "00000000-0000-0000-0000-000000000001",
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"phone"}},
				Phone:      "+420987654",
			},
		))
		Expect(err).To(MatchError(ContainSubstring("player '00000000-0000-0000-0000-000000000001' not found")))
	})

	It("can update player", func(ctx SpecContext) {
		response, err := svcClient.UpdatePlayer(ctx, connect.NewRequest(
			&api.UpdatePlayerRequest{
				Id:         player.GetId(),
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"phone", "email"}},
				Phone:      "+420987654",
				Email:      "john.brown.2@game.com",
			},
		))
		Expect(err).NotTo(HaveOccurred())
		Expect(response.Msg.GetPhone()).To(Equal("+420987654"))
		Expect(response.Msg.GetEmail()).To(Equal("john.brown.2@game.com"))
	})

	It("cannot update player IC to a invalid value", func(ctx SpecContext) {
		_, err := svcClient.UpdatePlayer(ctx, connect.NewRequest(
			&api.UpdatePlayerRequest{
				Id:         player.GetId(),
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"ic"}},
				Ic:         "12345678",
			},
		))
		Expect(err).To(MatchError(ContainSubstring("ic: value is not valid IČO")))
	})

	It("can update player IC to a valid value", func(ctx SpecContext) {
		response, err := svcClient.UpdatePlayer(ctx, connect.NewRequest(
			&api.UpdatePlayerRequest{
				Id:         player.GetId(),
				UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"ic"}},
				Ic:         "00023272",
			},
		))
		Expect(err).NotTo(HaveOccurred())
		Expect(response.Msg.GetIc()).To(Equal("00023272"))
		Expect(response.Msg.GetVatId()).To(Equal("CZ00023272"))
		Expect(response.Msg.GetAddress()).To(Equal("Václavské náměstí 1700/68\nNové Město\n11000 Praha 1"))
	})

	It("can get player by id", func(ctx SpecContext) {
		response, err := svcClient.GetPlayer(ctx, connect.NewRequest(
			&api.GetPlayerRequest{
				Id: player.GetId(),
			},
		))
		Expect(err).NotTo(HaveOccurred())
		Expect(response.Msg).Should(PointTo(MatchFields(IgnoreExtras, Fields{
			"FirstName": Equal("John"),
			"LastName":  Equal("Brown"),
			"Phone":     Equal("+420987654"),
			"Email":     Equal("john.brown.2@game.com"),
			"Ic":        Equal("00023272"),
			"VatId":     Equal("CZ00023272"),
			"Address":   Equal("Václavské náměstí 1700/68\nNové Město\n11000 Praha 1"),
		})))
	})

	It("cannot delete player, if not exists", func(ctx SpecContext) {
		_, err := svcClient.DeletePlayer(ctx, connect.NewRequest(
			&api.DeletePlayerRequest{
				Id: "00000000-0000-0000-0000-000000000001",
			},
		))
		Expect(err).To(MatchError(ContainSubstring("player '00000000-0000-0000-0000-000000000001' not found")))
	})

	It("can delete player", func(ctx SpecContext) {
		_, err := svcClient.DeletePlayer(ctx, connect.NewRequest(
			&api.DeletePlayerRequest{
				Id: player.GetId(),
			},
		))
		Expect(err).NotTo(HaveOccurred())
		player = nil
	})
})
