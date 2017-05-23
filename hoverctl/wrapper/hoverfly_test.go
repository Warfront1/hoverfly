package wrapper

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/util"
	"github.com/SpectoLabs/hoverfly/hoverctl/configuration"
	. "github.com/onsi/gomega"
)

func Test_isLocal_WhenLocalhost(t *testing.T) {
	RegisterTestingT(t)

	Expect(IsLocal("localhost")).To(BeTrue())
}

func Test_isLocal_WhenLocalhost_WithHttp(t *testing.T) {
	RegisterTestingT(t)

	Expect(IsLocal("http://localhost")).To(BeTrue())
}

func Test_isLocal_WhenLocalhostIP(t *testing.T) {
	RegisterTestingT(t)

	Expect(IsLocal("127.0.0.1")).To(BeTrue())
}

func Test_isLocal_WhenLocalhostIP_WithHttp(t *testing.T) {
	RegisterTestingT(t)

	Expect(IsLocal("http://127.0.0.1")).To(BeTrue())
}

func Test_isLocal_WhenAnotherDNS(t *testing.T) {
	RegisterTestingT(t)

	Expect(IsLocal("specto.io")).To(BeFalse())
}

func Test_BuildUrl_AddsHostAdminPortAndPath(t *testing.T) {
	RegisterTestingT(t)

	target := configuration.Target{
		Host:      "http://localhost",
		AdminPort: 1234,
	}

	Expect(BuildURL(target, "/something")).To(Equal("http://localhost:1234/something"))
}

func Test_BuildUrl_AddsHostAdminPortAndPath_Https(t *testing.T) {
	RegisterTestingT(t)

	target := configuration.Target{
		Host:      "https://localhost",
		AdminPort: 1234,
	}

	Expect(BuildURL(target, "/something")).To(Equal("https://localhost:1234/something"))
}

func Test_BuildUrl_AddsHttpIfHostIsLocalhost(t *testing.T) {
	RegisterTestingT(t)

	target := configuration.Target{
		Host:      "localhost",
		AdminPort: 1234,
	}

	Expect(BuildURL(target, "/something")).To(Equal("http://localhost:1234/something"))
}

func Test_BuildUrl_AddsHttpIfHostIsExternal(t *testing.T) {
	RegisterTestingT(t)

	target := configuration.Target{
		Host:      "test-instance.hoverfly.io",
		AdminPort: 1234,
	}

	Expect(BuildURL(target, "/something")).To(Equal("https://test-instance.hoverfly.io:1234/something"))
}

func Test_Shutdown_SendsCorrectHTTPRequest(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("DELETE"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/shutdown"),
						},
					},
					Response: v2.ResponseDetailsView{
						Status: 200,
						Body:   ``,
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := Shutdown(target)
	Expect(err).To(BeNil())
}

func Test_Shutdown_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	err := Shutdown(inaccessibleTarget)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not connect to Hoverfly at something:1234"))
}

func Test_Shutdown_ErrorsWhen_HoverflyReturnsNon200(t *testing.T) {
	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("DELETE"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/v2/shutdown"),
						},
					},
					Response: v2.ResponseDetailsView{
						Status: 400,
						Body:   "{\"error\":\"test error\"}",
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := Shutdown(target)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Could not stop Hoverfly\n\ntest error"))
}

func Test_CheckIfRunning_ReturnsNilWhen_HoverflyAccessible(t *testing.T) {

	RegisterTestingT(t)

	hoverfly.DeleteSimulation()
	hoverfly.PutSimulation(v2.SimulationViewV2{
		v2.DataViewV2{
			RequestResponsePairs: []v2.RequestResponsePairViewV2{
				v2.RequestResponsePairViewV2{
					Request: v2.RequestDetailsViewV2{
						Method: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("GET"),
						},
						Path: &v2.RequestFieldMatchersView{
							ExactMatch: util.StringToPointer("/api/public"),
						},
					},
					Response: v2.ResponseDetailsView{
						Status: 200,
						Body:   "",
					},
				},
			},
		},
		v2.MetaView{
			SchemaVersion: "v2",
		},
	})

	err := CheckIfRunning(target)

	Expect(err).To(BeNil())
}

func Test_CheckIfRunning_ErrorsWhen_HoverflyNotAccessible(t *testing.T) {
	RegisterTestingT(t)

	err := CheckIfRunning(inaccessibleTarget)

	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("Target Hoverfly is not running\n\nRun `hoverctl start -t ` to start it"))
}
