package player_test

import (
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/michaljurecko/ch-demo/api/gen/go/demo/v1/apiconnect"
	"github.com/michaljurecko/ch-demo/internal/pkg/app/demo/service/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRace(t *testing.T) {
	t.Parallel()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Race")
}

var _ = Describe("Race Operations", Ordered, func() {
	// Run server on a random port and create generated client
	var svcClient apiconnect.ApiServiceClient
	BeforeAll(func(ctx SpecContext) {
		baseURL := test.StartTestServer(ctx)
		svcClient = apiconnect.NewApiServiceClient(http.DefaultClient, baseURL)
	})

	It("can list races", func(ctx SpecContext) {
		response, err := svcClient.ListRaces(ctx, connect.NewRequest(&emptypb.Empty{}))
		Expect(err).NotTo(HaveOccurred())
		Expect(response.Msg.Races).NotTo(BeEmpty())
	})
})
