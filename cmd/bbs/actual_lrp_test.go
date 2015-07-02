package main_test

import (
	"time"

	"github.com/cloudfoundry-incubator/bbs/db/etcd/internal/test_helpers"
	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/gogo/protobuf/proto"
	"github.com/tedsuo/ifrit/ginkgomon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ActualLRP API", func() {
	var testHelper *test_helpers.TestHelper

	BeforeEach(func() {
		bbsProcess = ginkgomon.Invoke(bbsRunner)
		testHelper = test_helpers.NewTestHelper(etcdClient)
	})

	AfterEach(func() {
		ginkgomon.Kill(bbsProcess)
	})

	Describe("GET /v1/actual_lrps_groups", func() {
		const (
			cellID          = "cell-id"
			noExpirationTTL = 0

			baseProcessGuid  = "base-process-guid"
			baseDomain       = "base-domain"
			baseInstanceGuid = "base-instance-guid"

			otherProcessGuid  = "other-process-guid"
			otherDomain       = "other-domain"
			otherInstanceGuid = "other-instance-guid"

			baseIndex  = 1
			otherIndex = 1

			evacuatingInstanceGuid = "evacuating-instance-guid"
		)

		var (
			expectedActualLRPGroups []*models.ActualLRPGroup
			actualActualLRPGroups   []*models.ActualLRPGroup

			baseLRP       models.ActualLRP
			otherLRP      models.ActualLRP
			evacuatingLRP models.ActualLRP

			baseLRPKey          models.ActualLRPKey
			baseLRPInstanceKey  models.ActualLRPInstanceKey
			otherLRPKey         models.ActualLRPKey
			otherLRPInstanceKey models.ActualLRPInstanceKey
			netInfo             models.ActualLRPNetInfo

			filter models.ActualLRPFilter

			getErr error
		)

		BeforeEach(func() {
			baseLRPKey = models.NewActualLRPKey(baseProcessGuid, baseIndex, baseDomain)
			baseLRPInstanceKey = models.NewActualLRPInstanceKey(baseInstanceGuid, cellID)

			otherLRPKey = models.NewActualLRPKey(otherProcessGuid, otherIndex, otherDomain)
			otherLRPInstanceKey = models.NewActualLRPInstanceKey(otherInstanceGuid, cellID)

			netInfo = models.NewActualLRPNetInfo("127.0.0.1", []*models.PortMapping{{proto.Uint32(8080), proto.Uint32(80)}})

			baseLRP = models.ActualLRP{
				ActualLRPKey:         baseLRPKey,
				ActualLRPInstanceKey: baseLRPInstanceKey,
				ActualLRPNetInfo:     netInfo,
				State:                proto.String(models.ActualLRPStateRunning),
				Since:                proto.Int64(time.Now().UnixNano()),
			}
			evacuatingLRP = models.ActualLRP{
				ActualLRPKey:         baseLRPKey,
				ActualLRPInstanceKey: models.NewActualLRPInstanceKey(evacuatingInstanceGuid, cellID),
				ActualLRPNetInfo:     netInfo,
				State:                proto.String(models.ActualLRPStateRunning),
				Since:                proto.Int64(time.Now().UnixNano() - 1000),
			}

			otherLRP = models.ActualLRP{
				ActualLRPKey:         otherLRPKey,
				ActualLRPInstanceKey: otherLRPInstanceKey,
				ActualLRPNetInfo:     netInfo,
				State:                proto.String(models.ActualLRPStateRunning),
				Since:                proto.Int64(time.Now().UnixNano()),
			}

			testHelper.SetRawActualLRP(baseLRP)
			testHelper.SetRawActualLRP(otherLRP)
			testHelper.SetRawEvacuatingActualLRP(evacuatingLRP, noExpirationTTL)
		})

		JustBeforeEach(func() {
			actualActualLRPGroups, getErr = client.ActualLRPGroups(filter)
		})

		It("responds without error", func() {
			Expect(getErr).NotTo(HaveOccurred())
		})

		Context("when not filtering", func() {
			It("has the correct number of responses", func() {
				Expect(actualActualLRPGroups).To(HaveLen(2))
			})

			It("returns all actual lrps from the bbs", func() {
				expectedActualLRPGroups = []*models.ActualLRPGroup{{Instance: &baseLRP, Evacuating: &evacuatingLRP}, {Instance: &otherLRP}}
				Expect(actualActualLRPGroups).To(ConsistOf(expectedActualLRPGroups))
			})
		})

		Context("when filtering by domain", func() {
			BeforeEach(func() {
				filter = models.ActualLRPFilter{Domain: baseDomain}
			})

			It("returns all actual lrps from the bbs", func() {
				expectedActualLRPGroups = []*models.ActualLRPGroup{{Instance: &baseLRP, Evacuating: &evacuatingLRP}}
				Expect(actualActualLRPGroups).To(ConsistOf(expectedActualLRPGroups))
			})
		})
	})
})
