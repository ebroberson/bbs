package handlers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/bbs"
	"github.com/cloudfoundry-incubator/bbs/db/fakes"
	"github.com/cloudfoundry-incubator/bbs/handlers"
	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager"
)

var _ = Describe("ActualLRP Handlers", func() {
	var (
		logger           lager.Logger
		fakeActualLRPDB  *fakes.FakeActualLRPDB
		responseRecorder *httptest.ResponseRecorder
		handler          *handlers.ActualLRPHandler

		actualLRP1     models.ActualLRP
		actualLRP2     models.ActualLRP
		evacuatingLRP2 models.ActualLRP
	)

	BeforeEach(func() {
		fakeActualLRPDB = new(fakes.FakeActualLRPDB)
		logger = lager.NewLogger("test")
		logger.RegisterSink(lager.NewWriterSink(GinkgoWriter, lager.DEBUG))
		responseRecorder = httptest.NewRecorder()
		handler = handlers.NewActualLRPHandler(fakeActualLRPDB, logger)
	})

	Describe("ActualLRPGroups", func() {

		BeforeEach(func() {
			actualLRP1 = models.ActualLRP{
				ActualLRPKey: models.NewActualLRPKey(
					"process-guid-0",
					1,
					"domain-0",
				),
				ActualLRPInstanceKey: models.NewActualLRPInstanceKey(
					"instance-guid-0",
					"cell-id-0",
				),
				State: proto.String(models.ActualLRPStateRunning),
				Since: proto.Int64(1138),
			}

			actualLRP2 = models.ActualLRP{
				ActualLRPKey: models.NewActualLRPKey(
					"process-guid-1",
					2,
					"domain-1",
				),
				ActualLRPInstanceKey: models.NewActualLRPInstanceKey(
					"instance-guid-1",
					"cell-id-1",
				),
				State: proto.String(models.ActualLRPStateClaimed),
				Since: proto.Int64(4444),
			}

			evacuatingLRP2 = actualLRP2
			evacuatingLRP2.State = proto.String(models.ActualLRPStateRunning)
			evacuatingLRP2.Since = proto.Int64(3417)
		})

		JustBeforeEach(func() {
			handler.ActualLRPGroups(responseRecorder, newTestRequest(""))
		})

		Context("when reading actuallrps from DB succeeds", func() {
			var actualLRPGroups *models.ActualLRPGroups

			BeforeEach(func() {
				actualLRPGroups = &models.ActualLRPGroups{
					[]*models.ActualLRPGroup{
						{Instance: &actualLRP1},
						{Instance: &actualLRP2, Evacuating: &evacuatingLRP2},
					},
				}
				fakeActualLRPDB.ActualLRPGroupsReturns(actualLRPGroups, nil)
			})

			It("call the DB to retrieve the actual lrpGroups", func() {
				Expect(fakeActualLRPDB.ActualLRPGroupsCallCount()).To(Equal(1))
			})

			It("responds with 200 Status OK", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			})

			It("returns a list of actual lrpGroups", func() {
				response := &models.ActualLRPGroups{}
				err := response.Unmarshal(responseRecorder.Body.Bytes())
				Expect(err).NotTo(HaveOccurred())

				Expect(response).To(Equal(actualLRPGroups))
			})
		})

		Context("when the DB returns no actual lrpGroups", func() {
			BeforeEach(func() {
				fakeActualLRPDB.ActualLRPGroupsReturns(&models.ActualLRPGroups{}, nil)
			})

			It("responds with 200 Status OK", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			})

			It("returns an empty list", func() {
				response := &models.ActualLRPGroups{}
				err := response.Unmarshal(responseRecorder.Body.Bytes())
				Expect(err).NotTo(HaveOccurred())

				Expect(response).To(Equal(&models.ActualLRPGroups{}))
			})
		})

		Context("when the DB errors out", func() {
			BeforeEach(func() {
				fakeActualLRPDB.ActualLRPGroupsReturns(&models.ActualLRPGroups{}, errors.New("Something went wrong"))
			})

			It("responds with an error", func() {
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			})

			It("provides relevant error information", func() {
				var bbsError bbs.Error
				err := bbsError.Unmarshal(responseRecorder.Body.Bytes())
				Expect(err).NotTo(HaveOccurred())

				Expect(bbsError).To(Equal(bbs.Error{
					Type:    proto.String(bbs.UnknownError),
					Message: proto.String("Something went wrong"),
				}))
			})
		})
	})
})
