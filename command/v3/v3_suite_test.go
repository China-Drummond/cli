package v3_test

import (
	"code.cloudfoundry.org/cli/command"
	"code.cloudfoundry.org/cli/command/v3/shared"
	"code.cloudfoundry.org/cli/command/v3/shared/sharedfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"

	"testing"
)

func TestV3(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "V3 Command Suite")
}

var _ = BeforeEach(func() {
	log.SetLevel(log.PanicLevel)
})

func NewTestAppSummaryDisplayer(appName string, ui command.UI, config command.Config) (*shared.AppSummaryDisplayer, *sharedfakes.FakeV2AppActor, *sharedfakes.FakeV3AppSummaryActor) {
	v2Fake := new(sharedfakes.FakeV2AppActor)
	v3Fake := new(sharedfakes.FakeV3AppSummaryActor)
	return shared.NewAppSummaryDisplayer(appName, ui, config, v2Fake, v3Fake), v2Fake, v3Fake
}
