package models_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudfoundry-incubator/bbs/models"
	"github.com/gogo/protobuf/proto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func defaultCrashedActual(crashCount int32, lastCrashed int64) *models.ActualLRP {
	return &models.ActualLRP{
		ActualLRPKey: models.NewActualLRPKey("p-guid", 0, "domain"),
		State:        proto.String(models.ActualLRPStateCrashed),
		CrashCount:   proto.Int32(crashCount),
		Since:        proto.Int64(lastCrashed),
	}
}

type crashInfoTest interface {
	Test()
}

type crashInfoTests []crashInfoTest

func (tests crashInfoTests) Test() {
	for _, test := range tests {
		test.Test()
	}
}

type crashInfoBackoffTest struct {
	*models.ActualLRP
	WaitTime time.Duration
}

func newCrashInfoBackoffTest(crashCount int32, lastCrashed int64, waitTime time.Duration) crashInfoTest {
	return crashInfoBackoffTest{
		ActualLRP: defaultCrashedActual(crashCount, lastCrashed),
		WaitTime:  waitTime,
	}
}

func (test crashInfoBackoffTest) Test() {
	Context(fmt.Sprintf("when the crashCount is %d and the wait time is %s", test.CrashCount, test.WaitTime), func() {
		It("should NOT restart before the expected wait time", func() {
			calc := models.NewDefaultRestartCalculator()
			currentTimestamp := test.GetSince() + test.WaitTime.Nanoseconds() - time.Second.Nanoseconds()
			Expect(test.ShouldRestartCrash(time.Unix(0, currentTimestamp), calc)).To(BeFalse())
		})

		It("should restart after the expected wait time", func() {
			calc := models.NewDefaultRestartCalculator()
			currentTimestamp := test.GetSince() + test.WaitTime.Nanoseconds()
			Expect(test.ShouldRestartCrash(time.Unix(0, currentTimestamp), calc)).To(BeTrue())
		})
	})
}

type crashInfoNeverStartTest struct {
	*models.ActualLRP
}

func newCrashInfoNeverStartTest(crashCount int32, lastCrashed int64) crashInfoTest {
	return crashInfoNeverStartTest{
		ActualLRP: defaultCrashedActual(crashCount, lastCrashed),
	}
}

func (test crashInfoNeverStartTest) Test() {
	Context(fmt.Sprintf("when the crashCount is %d", test.CrashCount), func() {
		It("should never restart regardless of the wait time", func() {
			calc := models.NewDefaultRestartCalculator()
			theFuture := test.GetSince() + time.Hour.Nanoseconds()
			Expect(test.ShouldRestartCrash(time.Unix(0, 0), calc)).To(BeFalse())
			Expect(test.ShouldRestartCrash(time.Unix(0, test.GetSince()), calc)).To(BeFalse())
			Expect(test.ShouldRestartCrash(time.Unix(0, theFuture), calc)).To(BeFalse())
		})
	})
}

type crashInfoAlwaysStartTest struct {
	*models.ActualLRP
}

func newCrashInfoAlwaysStartTest(crashCount int32, lastCrashed int64) crashInfoTest {
	return crashInfoAlwaysStartTest{
		ActualLRP: defaultCrashedActual(crashCount, lastCrashed),
	}
}

func (test crashInfoAlwaysStartTest) Test() {
	Context(fmt.Sprintf("when the crashCount is %d", test.CrashCount), func() {
		It("should restart regardless of the wait time", func() {
			calc := models.NewDefaultRestartCalculator()
			theFuture := test.GetSince() + time.Hour.Nanoseconds()
			Expect(test.ShouldRestartCrash(time.Unix(0, 0), calc)).To(BeTrue())
			Expect(test.ShouldRestartCrash(time.Unix(0, test.GetSince()), calc)).To(BeTrue())
			Expect(test.ShouldRestartCrash(time.Unix(0, theFuture), calc)).To(BeTrue())
		})
	})
}

func testBackoffCount(maxBackoffDuration time.Duration, expectedBackoffCount int32) {
	It(fmt.Sprintf("sets the MaxBackoffCount to %d based on the MaxBackoffDuration %s and the CrashBackoffMinDuration", expectedBackoffCount, maxBackoffDuration), func() {
		calc := models.NewRestartCalculator(models.DefaultImmediateRestarts, maxBackoffDuration, models.DefaultMaxRestarts)
		Expect(calc.MaxBackoffCount).To(Equal(expectedBackoffCount))
	})
}

var _ = Describe("RestartCalculator", func() {

	Describe("NewRestartCalculator", func() {
		testBackoffCount(20*time.Minute, 5)
		testBackoffCount(16*time.Minute, 5)
		testBackoffCount(8*time.Minute, 4)
		testBackoffCount(119*time.Second, 2)
		testBackoffCount(120*time.Second, 2)
		testBackoffCount(models.CrashBackoffMinDuration, 0)

		It("should work...", func() {
			nanoseconds := func(seconds int) int64 {
				return int64(seconds * 1000000000)
			}

			calc := models.NewRestartCalculator(3, 119*time.Second, 200)
			Expect(calc.ShouldRestart(0, 0, 0)).To(BeTrue())
			Expect(calc.ShouldRestart(0, 0, 1)).To(BeTrue())
			Expect(calc.ShouldRestart(0, 0, 2)).To(BeTrue())

			Expect(calc.ShouldRestart(0, 0, 3)).To(BeFalse())
			Expect(calc.ShouldRestart(nanoseconds(30), 0, 3)).To(BeTrue())

			Expect(calc.ShouldRestart(nanoseconds(30), 0, 4)).To(BeFalse())
			Expect(calc.ShouldRestart(nanoseconds(59), 0, 4)).To(BeFalse())
			Expect(calc.ShouldRestart(nanoseconds(60), 0, 4)).To(BeTrue())
			Expect(calc.ShouldRestart(nanoseconds(60), 0, 5)).To(BeFalse())
			Expect(calc.ShouldRestart(nanoseconds(118), 0, 5)).To(BeFalse())
			Expect(calc.ShouldRestart(nanoseconds(119), 0, 5)).To(BeTrue())
		})
	})

	Describe("Validate", func() {
		It("the default values are valid", func() {
			calc := models.NewDefaultRestartCalculator()
			Expect(calc.Validate()).NotTo(HaveOccurred())
		})

		It("invalid when MaxBackoffDuration is lower than the CrashBackoffMinDuration", func() {
			calc := models.NewRestartCalculator(models.DefaultImmediateRestarts, models.CrashBackoffMinDuration-time.Second, models.DefaultMaxRestarts)
			Expect(calc.Validate()).To(HaveOccurred())
		})
	})
})

var _ = Describe("ActualLRP", func() {
	Describe("ShouldRestartCrash", func() {
		Context("when the lpr is CRASHED", func() {
			const maxWaitTime = 16 * time.Minute
			var now = time.Now().UnixNano()
			var crashTests = crashInfoTests{
				newCrashInfoAlwaysStartTest(0, now),
				newCrashInfoAlwaysStartTest(1, now),
				newCrashInfoAlwaysStartTest(2, now),
				newCrashInfoBackoffTest(3, now, 30*time.Second),
				newCrashInfoBackoffTest(7, now, 8*time.Minute),
				newCrashInfoBackoffTest(8, now, maxWaitTime),
				newCrashInfoBackoffTest(199, now, maxWaitTime),
				newCrashInfoNeverStartTest(200, now),
				newCrashInfoNeverStartTest(201, now),
			}

			crashTests.Test()
		})

		Context("when the lrp is not CRASHED", func() {
			It("returns false", func() {
				now := time.Now()
				actual := defaultCrashedActual(0, now.UnixNano())
				calc := models.NewDefaultRestartCalculator()
				for _, state := range models.ActualLRPStates {
					actual.State = proto.String(state)
					if state == models.ActualLRPStateCrashed {
						Expect(actual.ShouldRestartCrash(now, calc)).To(BeTrue(), "should restart CRASHED lrp")
					} else {
						Expect(actual.ShouldRestartCrash(now, calc)).To(BeFalse(), fmt.Sprintf("should not restart %s lrp", state))
					}
				}
			})
		})
	})

	Describe("ActualLRPKey", func() {
		Describe("Validate", func() {
			var actualLRPKey *models.ActualLRPKey

			BeforeEach(func() {
				actualLRPKey = models.NewActualLRPKey("process-guid", 1, "domain")
			})

			Context("when valid", func() {
				It("returns nil", func() {
					Expect(actualLRPKey.Validate()).To(BeNil())
				})
			})

			Context("when the ProcessGuid is blank", func() {
				BeforeEach(func() {
					actualLRPKey.ProcessGuid = proto.String("")
				})

				It("returns a validation error", func() {
					Expect(actualLRPKey.Validate()).To(ConsistOf(models.ErrInvalidField{"process_guid"}))
				})
			})

			Context("when the Domain is blank", func() {
				BeforeEach(func() {
					actualLRPKey.Domain = proto.String("")
				})

				It("returns a validation error", func() {
					Expect(actualLRPKey.Validate()).To(ConsistOf(models.ErrInvalidField{"domain"}))
				})
			})

			Context("when the Index is negative", func() {
				BeforeEach(func() {
					actualLRPKey.Index = proto.Int32(-1)
				})

				It("returns a validation error", func() {
					Expect(actualLRPKey.Validate()).To(ConsistOf(models.ErrInvalidField{"index"}))
				})
			})
		})
	})

	Describe("ActualLRPInstanceKey", func() {
		Describe("Validate", func() {
			var actualLRPInstanceKey *models.ActualLRPInstanceKey

			Context("when both instance guid and cell id are specified", func() {
				It("returns nil", func() {
					actualLRPInstanceKey = models.NewActualLRPInstanceKey("instance-guid", "cell-id")
					Expect(actualLRPInstanceKey.Validate()).To(BeNil())
				})
			})

			Context("when both instance guid and cell id are empty", func() {
				It("returns a validation error", func() {
					actualLRPInstanceKey = models.NewActualLRPInstanceKey("", "")
					Expect(actualLRPInstanceKey.Validate()).To(ConsistOf(
						models.ErrInvalidField{"cell_id"},
						models.ErrInvalidField{"instance_guid"},
					))

				})
			})

			Context("when only the instance guid is specified", func() {
				It("returns a validation error", func() {
					actualLRPInstanceKey = models.NewActualLRPInstanceKey("instance-guid", "")
					Expect(actualLRPInstanceKey.Validate()).To(ConsistOf(models.ErrInvalidField{"cell_id"}))
				})
			})

			Context("when only the cell id is specified", func() {
				It("returns a validation error", func() {
					actualLRPInstanceKey = models.NewActualLRPInstanceKey("", "cell-id")
					Expect(actualLRPInstanceKey.Validate()).To(ConsistOf(models.ErrInvalidField{"instance_guid"}))
				})
			})
		})
	})

	Describe("ActualLRPGroup", func() {
		Describe("Resolve", func() {
			var (
				instanceLRP   *models.ActualLRP
				evacuatingLRP *models.ActualLRP

				group models.ActualLRPGroup

				resolvedLRP *models.ActualLRP
				evacuating  bool
				resolveErr  error
			)

			BeforeEach(func() {
				lrpKey := models.NewActualLRPKey("process-guid", 1, "domain")
				instanceLRP = &models.ActualLRP{
					ActualLRPKey: lrpKey,
					Since:        proto.Int64(1138),
				}
				evacuatingLRP = &models.ActualLRP{
					ActualLRPKey: lrpKey,
					Since:        proto.Int64(3417),
				}
			})

			JustBeforeEach(func() {
				resolvedLRP, evacuating, resolveErr = group.Resolve()
			})

			Context("When neither the Instance nor the Evacuating LRP is set", func() {
				BeforeEach(func() {
					group = models.ActualLRPGroup{}
				})

				It("returns ErrActualLRPGroupInvalid", func() {
					Expect(resolveErr).To(Equal(models.ErrActualLRPGroupInvalid))
				})
			})

			Context("When only the Instance LRP is set", func() {
				BeforeEach(func() {
					group = models.ActualLRPGroup{
						Instance: instanceLRP,
					}
				})

				It("returns the Instance LRP", func() {
					Expect(resolveErr).NotTo(HaveOccurred())
					Expect(resolvedLRP).To(Equal(instanceLRP))
					Expect(evacuating).To(BeFalse())
				})
			})

			Context("When only the Evacuating LRP is set", func() {
				BeforeEach(func() {
					group = models.ActualLRPGroup{
						Evacuating: evacuatingLRP,
					}
				})

				It("returns the Evacuating LRP", func() {
					Expect(resolveErr).NotTo(HaveOccurred())
					Expect(resolvedLRP).To(Equal(evacuatingLRP))
					Expect(evacuating).To(BeTrue())
				})
			})

			Context("When both the Instance and the Evacuating LRP are set", func() {
				BeforeEach(func() {
					group = models.ActualLRPGroup{
						Evacuating: evacuatingLRP,
						Instance:   instanceLRP,
					}
				})

				Context("When the Instance is UNCLAIMED", func() {
					BeforeEach(func() {
						instanceLRP.State = proto.String(models.ActualLRPStateUnclaimed)
					})

					It("returns the Evacuating LRP", func() {
						Expect(resolveErr).NotTo(HaveOccurred())
						Expect(resolvedLRP).To(Equal(evacuatingLRP))
						Expect(evacuating).To(BeTrue())
					})
				})

				Context("When the Instance is CLAIMED", func() {
					BeforeEach(func() {
						instanceLRP.State = proto.String(models.ActualLRPStateClaimed)
					})

					It("returns the Evacuating LRP", func() {
						Expect(resolveErr).NotTo(HaveOccurred())
						Expect(resolvedLRP).To(Equal(evacuatingLRP))
						Expect(evacuating).To(BeTrue())
					})
				})

				Context("When the Instance is RUNNING", func() {
					BeforeEach(func() {
						instanceLRP.State = proto.String(models.ActualLRPStateRunning)
					})

					It("returns the Instance LRP", func() {
						Expect(resolveErr).NotTo(HaveOccurred())
						Expect(resolvedLRP).To(Equal(instanceLRP))
						Expect(evacuating).To(BeFalse())
					})
				})

				Context("When the Instance is CRASHED", func() {
					BeforeEach(func() {
						instanceLRP.State = proto.String(models.ActualLRPStateCrashed)
					})

					It("returns the Instance LRP", func() {
						Expect(resolveErr).NotTo(HaveOccurred())
						Expect(resolvedLRP).To(Equal(instanceLRP))
						Expect(evacuating).To(BeFalse())
					})
				})
			})
		})
	})

	Describe("ActualLRP", func() {
		var lrp models.ActualLRP
		var lrpKey *models.ActualLRPKey
		var instanceKey *models.ActualLRPInstanceKey
		var netInfo *models.ActualLRPNetInfo

		lrpPayload := `{
    "process_guid":"some-guid",
    "instance_guid":"some-instance-guid",
    "address": "1.2.3.4",
    "ports": [
      { "container_port": 8080 },
      { "container_port": 8081, "host_port": 1234 }
    ],
    "index": 2,
    "state": "RUNNING",
    "since": 1138,
    "cell_id":"some-cell-id",
    "domain":"some-domain",
		"crash_count": 1,
		"modification_tag": {
			"epoch": "some-guid",
			"index": 50
		}
  }`

		BeforeEach(func() {
			lrpKey = models.NewActualLRPKey("some-guid", 2, "some-domain")
			instanceKey = models.NewActualLRPInstanceKey("some-instance-guid", "some-cell-id")
			netInfo = models.NewActualLRPNetInfo("1.2.3.4", []*models.PortMapping{
				{ContainerPort: proto.Uint32(8080)},
				{ContainerPort: proto.Uint32(8081), HostPort: proto.Uint32(1234)},
			})

			lrp = models.ActualLRP{
				ActualLRPKey:         lrpKey,
				ActualLRPInstanceKey: instanceKey,
				ActualLRPNetInfo:     netInfo,
				CrashCount:           proto.Int32(1),
				State:                proto.String(models.ActualLRPStateRunning),
				Since:                proto.Int64(1138),
				ModificationTag: &models.ModificationTag{
					Epoch: proto.String("some-guid"),
					Index: proto.Uint32(50),
				},
			}
		})

		Describe("To JSON", func() {
			It("should JSONify", func() {
				marshalled, err := json.Marshal(&lrp)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(marshalled)).To(MatchJSON(lrpPayload))
			})
		})

		Describe("FromJSON", func() {
			It("returns a LRP with correct fields", func() {
				aLRP := &models.ActualLRP{}
				err := models.FromJSON([]byte(lrpPayload), aLRP)
				Expect(err).NotTo(HaveOccurred())

				Expect(aLRP).To(Equal(&lrp))
			})

			Context("with an invalid payload", func() {
				It("returns the error", func() {
					aLRP := &models.ActualLRP{}
					err := models.FromJSON([]byte("something lol"), aLRP)
					Expect(err).To(HaveOccurred())
				})
			})

			for field, payload := range map[string]string{
				"process_guid":  `{"instance_guid": "instance_guid", "cell_id": "cell_id", "domain": "domain"}`,
				"instance_guid": `{"process_guid": "process-guid", "cell_id": "cell_id", "domain": "domain","state":"CLAIMED"}`,
				"cell_id":       `{"process_guid": "process-guid", "instance_guid": "instance_guid", "domain": "domain", "state":"RUNNING"}`,
				"domain":        `{"process_guid": "process-guid", "cell_id": "cell_id", "instance_guid": "instance_guid"}`,
			} {
				missingField := field
				jsonPayload := payload

				Context("when the json is missing a "+missingField, func() {
					It("returns an error indicating so", func() {
						aLRP := &models.ActualLRP{}
						err := models.FromJSON([]byte(jsonPayload), aLRP)
						Expect(err.Error()).To(ContainSubstring(missingField))
					})
				})
			}
		})

		Describe("AllowsTransitionTo", func() {
			var (
				before   *models.ActualLRP
				afterKey *models.ActualLRPKey
			)

			BeforeEach(func() {
				before = &models.ActualLRP{
					ActualLRPKey: models.NewActualLRPKey("fake-process-guid", 1, "fake-domain"),
				}
				afterKey = &models.ActualLRPKey{}
				*afterKey = *before.ActualLRPKey
			})

			Context("when the ProcessGuid fields differ", func() {
				BeforeEach(func() {
					before.ProcessGuid = proto.String("some-process-guid")
					afterKey.ProcessGuid = proto.String("another-process-guid")
				})

				It("is not allowed", func() {
					Expect(before.AllowsTransitionTo(afterKey, before.ActualLRPInstanceKey, before.GetState())).To(BeFalse())
				})
			})

			Context("when the Index fields differ", func() {
				BeforeEach(func() {
					before.Index = proto.Int32(1138)
					afterKey.Index = proto.Int32(3417)
				})

				It("is not allowed", func() {
					Expect(before.AllowsTransitionTo(afterKey, before.ActualLRPInstanceKey, before.GetState())).To(BeFalse())
				})
			})

			Context("when the Domain fields differ", func() {
				BeforeEach(func() {
					before.Domain = proto.String("some-domain")
					afterKey.Domain = proto.String("another-domain")
				})

				It("is not allowed", func() {
					Expect(before.AllowsTransitionTo(afterKey, before.ActualLRPInstanceKey, before.GetState())).To(BeFalse())
				})
			})

			Context("when the ProcessGuid, Index, and Domain are equivalent", func() {
				var (
					emptyKey                 = models.NewActualLRPInstanceKey("", "")
					claimedKey               = models.NewActualLRPInstanceKey("some-instance-guid", "some-cell-id")
					differentInstanceGuidKey = models.NewActualLRPInstanceKey("some-other-instance-guid", "some-cell-id")
					differentCellIDKey       = models.NewActualLRPInstanceKey("some-instance-guid", "some-other-cell-id")
				)

				type stateTableEntry struct {
					BeforeState       string
					AfterState        string
					BeforeInstanceKey *models.ActualLRPInstanceKey
					AfterInstanceKey  *models.ActualLRPInstanceKey
					Allowed           bool
				}

				var EntryToString = func(entry stateTableEntry) string {
					return fmt.Sprintf("is %t when the before has state %s and instance guid '%s' and cell id '%s' and the after has state %s and instance guid '%s' and cell id '%s'",
						entry.Allowed,
						entry.BeforeState,
						entry.BeforeInstanceKey.InstanceGuid,
						entry.BeforeInstanceKey.CellId,
						entry.AfterState,
						entry.AfterInstanceKey.InstanceGuid,
						entry.AfterInstanceKey.CellId,
					)
				}

				stateTable := []stateTableEntry{
					{models.ActualLRPStateUnclaimed, models.ActualLRPStateUnclaimed, emptyKey, emptyKey, true},
					{models.ActualLRPStateUnclaimed, models.ActualLRPStateClaimed, emptyKey, claimedKey, true},
					{models.ActualLRPStateUnclaimed, models.ActualLRPStateRunning, emptyKey, claimedKey, true},
					{models.ActualLRPStateClaimed, models.ActualLRPStateUnclaimed, claimedKey, emptyKey, true},
					{models.ActualLRPStateClaimed, models.ActualLRPStateClaimed, claimedKey, claimedKey, true},
					{models.ActualLRPStateClaimed, models.ActualLRPStateClaimed, claimedKey, differentInstanceGuidKey, false},
					{models.ActualLRPStateClaimed, models.ActualLRPStateClaimed, claimedKey, differentCellIDKey, false},
					{models.ActualLRPStateClaimed, models.ActualLRPStateRunning, claimedKey, claimedKey, true},
					{models.ActualLRPStateClaimed, models.ActualLRPStateRunning, claimedKey, differentInstanceGuidKey, true},
					{models.ActualLRPStateClaimed, models.ActualLRPStateRunning, claimedKey, differentCellIDKey, true},
					{models.ActualLRPStateRunning, models.ActualLRPStateUnclaimed, claimedKey, emptyKey, true},
					{models.ActualLRPStateRunning, models.ActualLRPStateClaimed, claimedKey, claimedKey, true},
					{models.ActualLRPStateRunning, models.ActualLRPStateClaimed, claimedKey, differentInstanceGuidKey, false},
					{models.ActualLRPStateRunning, models.ActualLRPStateClaimed, claimedKey, differentCellIDKey, false},
					{models.ActualLRPStateRunning, models.ActualLRPStateRunning, claimedKey, claimedKey, true},
					{models.ActualLRPStateRunning, models.ActualLRPStateClaimed, claimedKey, differentInstanceGuidKey, false},
					{models.ActualLRPStateRunning, models.ActualLRPStateClaimed, claimedKey, differentCellIDKey, false},
				}

				for _, entry := range stateTable {
					entry := entry
					It(EntryToString(entry), func() {
						before.State = proto.String(entry.BeforeState)
						before.ActualLRPInstanceKey = entry.BeforeInstanceKey
						Expect(before.AllowsTransitionTo(before.ActualLRPKey, entry.AfterInstanceKey, entry.AfterState)).To(Equal(entry.Allowed))
					})
				}
			})
		})

		Describe("Validate", func() {

			Context("when state is unclaimed", func() {
				BeforeEach(func() {
					lrp = models.ActualLRP{
						ActualLRPKey: lrpKey,
						State:        proto.String(models.ActualLRPStateUnclaimed),
						Since:        proto.Int64(1138),
					}
				})

				itValidatesPresenceOfTheLRPKey(&lrp)
				itValidatesAbsenceOfTheInstanceKey(&lrp)
				itValidatesAbsenceOfNetInfo(&lrp)
				itValidatesPresenceOfPlacementError(&lrp)
			})

			Context("when state is claimed", func() {
				BeforeEach(func() {
					lrp = models.ActualLRP{
						ActualLRPKey:         lrpKey,
						ActualLRPInstanceKey: instanceKey,
						State:                proto.String(models.ActualLRPStateClaimed),
						Since:                proto.Int64(1138),
					}
				})

				itValidatesPresenceOfTheLRPKey(&lrp)
				itValidatesPresenceOfTheInstanceKey(&lrp)
				itValidatesAbsenceOfNetInfo(&lrp)
				itValidatesAbsenceOfPlacementError(&lrp)
			})

			Context("when state is running", func() {
				BeforeEach(func() {
					lrp = models.ActualLRP{
						ActualLRPKey:         lrpKey,
						ActualLRPInstanceKey: instanceKey,
						ActualLRPNetInfo:     netInfo,
						State:                proto.String(models.ActualLRPStateRunning),
						Since:                proto.Int64(1138),
					}
				})

				itValidatesPresenceOfTheLRPKey(&lrp)
				itValidatesPresenceOfTheInstanceKey(&lrp)
				itValidatesPresenceOfNetInfo(&lrp)
				itValidatesAbsenceOfPlacementError(&lrp)
			})

			Context("when state is not set", func() {
				BeforeEach(func() {
					lrp = models.ActualLRP{
						ActualLRPKey: lrpKey,
						State:        proto.String(""),
						Since:        proto.Int64(1138),
					}
				})

				It("validate returns an error", func() {
					err := lrp.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("state"))
				})

			})

			Context("when since is not set", func() {
				BeforeEach(func() {
					lrp = models.ActualLRP{
						ActualLRPKey: lrpKey,
						State:        proto.String(models.ActualLRPStateUnclaimed),
						Since:        proto.Int64(0),
					}
				})

				It("validate returns an error", func() {
					err := lrp.Validate()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("since"))
				})
			})

			Context("when state is crashed", func() {
				BeforeEach(func() {
					lrp = models.ActualLRP{
						ActualLRPKey: lrpKey,
						State:        proto.String(models.ActualLRPStateCrashed),
						Since:        proto.Int64(1138),
					}
				})

				itValidatesPresenceOfTheLRPKey(&lrp)
				itValidatesAbsenceOfTheInstanceKey(&lrp)
				itValidatesAbsenceOfNetInfo(&lrp)
				itValidatesAbsenceOfPlacementError(&lrp)
			})
		})
	})
})

func itValidatesPresenceOfTheLRPKey(lrp *models.ActualLRP) {
	Context("when the lrp key is set", func() {
		BeforeEach(func() {
			lrp.ActualLRPKey = models.NewActualLRPKey("some-guid", 1, "domain")
		})

		It("validate does not return an error", func() {
			Expect(lrp.Validate()).NotTo(HaveOccurred())
		})
	})

	Context("when the lrp key is not set", func() {
		BeforeEach(func() {
			lrp.ActualLRPKey = &models.ActualLRPKey{}
		})

		It("validate returns an error", func() {
			err := lrp.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("process_guid"))
		})
	})
}

func itValidatesPresenceOfTheInstanceKey(lrp *models.ActualLRP) {
	Context("when the instance key is set", func() {
		BeforeEach(func() {
			lrp.ActualLRPInstanceKey = models.NewActualLRPInstanceKey("some-instance", "some-cell")
		})

		It("validate does not return an error", func() {
			Expect(lrp.Validate()).NotTo(HaveOccurred())
		})
	})

	Context("when the instance key is not set", func() {
		BeforeEach(func() {
			lrp.ActualLRPInstanceKey = &models.ActualLRPInstanceKey{}
		})

		It("validate returns an error", func() {
			err := lrp.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("instance_guid"))
		})
	})
}

func itValidatesAbsenceOfTheInstanceKey(lrp *models.ActualLRP) {
	Context("when the instance key is set", func() {
		BeforeEach(func() {
			lrp.ActualLRPInstanceKey = models.NewActualLRPInstanceKey("some-instance", "some-cell")
		})

		It("validate returns an error", func() {
			err := lrp.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("instance key"))
		})
	})

	Context("when the instance key is not set", func() {
		BeforeEach(func() {
			lrp.ActualLRPInstanceKey = nil
		})

		It("validate does not return an error", func() {
			Expect(lrp.Validate()).NotTo(HaveOccurred())
		})
	})
}

func itValidatesPresenceOfNetInfo(lrp *models.ActualLRP) {
	Context("when net info is set", func() {
		BeforeEach(func() {
			lrp.ActualLRPNetInfo = models.NewActualLRPNetInfo("1.2.3.4", []*models.PortMapping{})
		})

		It("validate does not return an error", func() {
			Expect(lrp.Validate()).NotTo(HaveOccurred())
		})
	})

	Context("when net info is not set", func() {
		BeforeEach(func() {
			lrp.ActualLRPNetInfo = &models.ActualLRPNetInfo{}
		})

		It("validate returns an error", func() {
			err := lrp.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("address"))
		})
	})
}

func itValidatesAbsenceOfNetInfo(lrp *models.ActualLRP) {
	Context("when net info is set", func() {
		BeforeEach(func() {
			lrp.ActualLRPNetInfo = models.NewActualLRPNetInfo("1.2.3.4", []*models.PortMapping{})
		})

		It("validate returns an error", func() {
			err := lrp.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("net info"))
		})
	})

	Context("when net info is not set", func() {
		BeforeEach(func() {
			lrp.ActualLRPNetInfo = nil
		})

		It("validate does not return an error", func() {
			Expect(lrp.Validate()).NotTo(HaveOccurred())
		})
	})
}

func itValidatesPresenceOfPlacementError(lrp *models.ActualLRP) {
	Context("when placement error is set", func() {
		BeforeEach(func() {
			lrp.PlacementError = proto.String("insufficient capacity")
		})

		It("validate does not return an error", func() {
			Expect(lrp.Validate()).NotTo(HaveOccurred())
		})
	})

	Context("when placement error is not set", func() {
		BeforeEach(func() {
			lrp.PlacementError = proto.String("")
		})

		It("validate does not return an error", func() {
			Expect(lrp.Validate()).NotTo(HaveOccurred())
		})
	})
}

func itValidatesAbsenceOfPlacementError(lrp *models.ActualLRP) {
	Context("when placement error is set", func() {
		BeforeEach(func() {
			lrp.PlacementError = proto.String("insufficient capacity")
		})

		It("validate returns an error", func() {
			err := lrp.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("placement error"))
		})
	})

	Context("when placement error is not set", func() {
		BeforeEach(func() {
			lrp.PlacementError = proto.String("")
		})

		It("validate does not return an error", func() {
			Expect(lrp.Validate()).NotTo(HaveOccurred())
		})
	})
}
