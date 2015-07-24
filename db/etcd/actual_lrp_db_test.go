package etcd_test

import (
	"github.com/cloudfoundry-incubator/bbs/db"
	. "github.com/cloudfoundry-incubator/bbs/db/etcd"
	"github.com/cloudfoundry-incubator/bbs/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ActualLRPDB", func() {
	const (
		cellID          = "cell-id"
		noExpirationTTL = 0

		baseProcessGuid   = "base-process-guid"
		baseDomain        = "base-domain"
		baseInstanceGuid  = "base-instance-guid"
		otherInstanceGuid = "other-instance-guid"

		baseIndex  = 1
		otherIndex = 2

		evacuatingInstanceGuid = "evacuating-instance-guid"

		otherDomainProcessGuid = "other-domain-process-guid"

		otherDomain = "other-domain"
		otherCellID = "other-cell-id"
	)

	var (
		etcdDB db.ActualLRPDB

		baseLRP        *models.ActualLRP
		otherIndexLRP  *models.ActualLRP
		evacuatingLRP  *models.ActualLRP
		otherDomainLRP *models.ActualLRP
		otherCellIdLRP *models.ActualLRP

		baseLRPKey          models.ActualLRPKey
		baseLRPInstanceKey  models.ActualLRPInstanceKey
		otherLRPInstanceKey models.ActualLRPInstanceKey
		netInfo             models.ActualLRPNetInfo
	)

	BeforeEach(func() {
		baseLRPKey = models.NewActualLRPKey(baseProcessGuid, baseIndex, baseDomain)
		baseLRPInstanceKey = models.NewActualLRPInstanceKey(baseInstanceGuid, cellID)
		otherLRPInstanceKey = models.NewActualLRPInstanceKey(otherInstanceGuid, otherCellID)

		netInfo = models.NewActualLRPNetInfo("127.0.0.1", models.NewPortMapping(8080, 80))

		baseLRP = &models.ActualLRP{
			ActualLRPKey:         baseLRPKey,
			ActualLRPInstanceKey: baseLRPInstanceKey,
			ActualLRPNetInfo:     netInfo,
			State:                models.ActualLRPStateRunning,
			Since:                clock.Now().UnixNano(),
		}

		evacuatingLRP = &models.ActualLRP{
			ActualLRPKey:         baseLRPKey,
			ActualLRPInstanceKey: models.NewActualLRPInstanceKey(evacuatingInstanceGuid, cellID),
			ActualLRPNetInfo:     netInfo,
			State:                models.ActualLRPStateRunning,
			Since:                clock.Now().UnixNano() - 1000,
		}

		otherIndexLRP = &models.ActualLRP{
			ActualLRPKey:         models.NewActualLRPKey(baseProcessGuid, otherIndex, baseDomain),
			ActualLRPInstanceKey: baseLRPInstanceKey,
			State:                models.ActualLRPStateClaimed,
			Since:                clock.Now().UnixNano(),
		}

		otherDomainLRP = &models.ActualLRP{
			ActualLRPKey:         models.NewActualLRPKey(otherDomainProcessGuid, baseIndex, otherDomain),
			ActualLRPInstanceKey: baseLRPInstanceKey,
			ActualLRPNetInfo:     netInfo,
			State:                models.ActualLRPStateRunning,
			Since:                clock.Now().UnixNano(),
		}

		otherCellIdLRP = &models.ActualLRP{
			ActualLRPKey:         models.NewActualLRPKey(otherDomainProcessGuid, otherIndex, otherDomain),
			ActualLRPInstanceKey: otherLRPInstanceKey,
			ActualLRPNetInfo:     netInfo,
			State:                models.ActualLRPStateRunning,
			Since:                clock.Now().UnixNano(),
		}
		etcdDB = NewETCD(etcdClient, clock)
	})

	Describe("ActualLRPGroups", func() {
		var filter models.ActualLRPFilter

		Context("when there are both /instance and /evacuating LRPs", func() {
			BeforeEach(func() {
				filter = models.ActualLRPFilter{}
				testHelper.SetRawActualLRP(baseLRP)
				testHelper.SetRawEvacuatingActualLRP(evacuatingLRP, noExpirationTTL)
				testHelper.SetRawActualLRP(otherDomainLRP)
				testHelper.SetRawEvacuatingActualLRP(otherIndexLRP, noExpirationTTL)
				testHelper.SetRawActualLRP(otherCellIdLRP)
			})

			It("returns all the /instance LRPs and /evacuating LRPs in groups", func() {
				actualLRPGroups, err := etcdDB.ActualLRPGroups(logger, filter)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualLRPGroups.GetActualLrpGroups()).To(ConsistOf(
					&models.ActualLRPGroup{Instance: baseLRP, Evacuating: evacuatingLRP},
					&models.ActualLRPGroup{Instance: otherDomainLRP, Evacuating: nil},
					&models.ActualLRPGroup{Instance: nil, Evacuating: otherIndexLRP},
					&models.ActualLRPGroup{Instance: otherCellIdLRP, Evacuating: nil},
				))
			})

			It("can filter by domain", func() {
				filter.Domain = otherDomain
				actualLRPGroups, err := etcdDB.ActualLRPGroups(logger, filter)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualLRPGroups.GetActualLrpGroups()).To(ConsistOf(
					&models.ActualLRPGroup{Instance: otherDomainLRP, Evacuating: nil},
					&models.ActualLRPGroup{Instance: otherCellIdLRP, Evacuating: nil},
				))
			})

			It("can filter by cell id", func() {
				filter.CellID = otherCellID
				actualLRPGroups, err := etcdDB.ActualLRPGroups(logger, filter)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualLRPGroups.GetActualLrpGroups()).To(ConsistOf(
					&models.ActualLRPGroup{Instance: otherCellIdLRP, Evacuating: nil},
				))
			})
		})

		Context("when there are no LRPs", func() {
			It("returns an empty list", func() {
				actualLRPGroups, err := etcdDB.ActualLRPGroups(logger, filter)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualLRPGroups).NotTo(BeNil())
				Expect(actualLRPGroups.GetActualLrpGroups()).To(BeEmpty())
			})
		})

		Context("when the root node exists with no child nodes", func() {
			BeforeEach(func() {
				testHelper.SetRawActualLRP(baseLRP)

				processGuid := baseLRP.ActualLRPKey.GetProcessGuid()
				_, err := etcdClient.Delete(ActualLRPProcessDir(processGuid), true)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an empty list", func() {
				actualLRPGroups, err := etcdDB.ActualLRPGroups(logger, filter)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualLRPGroups).NotTo(BeNil())
				Expect(actualLRPGroups.GetActualLrpGroups()).To(BeEmpty())
			})
		})

		Context("when there is invalid data", func() {
			BeforeEach(func() {
				testHelper.CreateValidActualLRP("some-guid", 0)
				testHelper.CreateMalformedActualLRP("some-other-guid", 0)
				testHelper.CreateValidActualLRP("some-third-guid", 0)
			})

			It("errors", func() {
				_, err := etcdDB.ActualLRPGroups(logger, filter)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when etcd is not there", func() {
			BeforeEach(func() {
				etcdRunner.Stop()
			})

			AfterEach(func() {
				etcdRunner.Start()
			})

			It("errors", func() {
				_, err := etcdDB.ActualLRPGroups(logger, filter)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ActualLRPGroupsByProcessGuid", func() {
		Context("when there are both /instance and /evacuating LRPs", func() {
			BeforeEach(func() {
				testHelper.SetRawActualLRP(baseLRP)
				testHelper.SetRawEvacuatingActualLRP(evacuatingLRP, noExpirationTTL)
				testHelper.SetRawActualLRP(otherDomainLRP)
				testHelper.SetRawEvacuatingActualLRP(otherIndexLRP, noExpirationTTL)
				testHelper.SetRawActualLRP(otherCellIdLRP)
			})

			It("returns all the /instance LRPs and /evacuating LRPs in groups", func() {
				actualLRPGroups, err := etcdDB.ActualLRPGroupsByProcessGuid(logger, baseProcessGuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualLRPGroups.GetActualLrpGroups()).To(ConsistOf(
					&models.ActualLRPGroup{Instance: baseLRP, Evacuating: evacuatingLRP},
					&models.ActualLRPGroup{Instance: nil, Evacuating: otherIndexLRP},
				))
			})
		})

		Context("when there are no LRPs", func() {
			It("returns an empty list", func() {
				actualLRPGroups, err := etcdDB.ActualLRPGroupsByProcessGuid(logger, baseProcessGuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualLRPGroups).NotTo(BeNil())
				Expect(actualLRPGroups.GetActualLrpGroups()).To(BeEmpty())
			})
		})

		Context("when the root node exists with no child nodes", func() {
			BeforeEach(func() {
				testHelper.SetRawActualLRP(baseLRP)

				processGuid := baseLRP.ActualLRPKey.GetProcessGuid()
				_, err := etcdClient.Delete(ActualLRPProcessDir(processGuid), true)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an empty list", func() {
				actualLRPGroups, err := etcdDB.ActualLRPGroupsByProcessGuid(logger, baseProcessGuid)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualLRPGroups).NotTo(BeNil())
				Expect(actualLRPGroups.GetActualLrpGroups()).To(BeEmpty())
			})
		})

		Context("when there is invalid data", func() {
			BeforeEach(func() {
				testHelper.CreateValidActualLRP("some-guid", 0)
				testHelper.CreateMalformedActualLRP("some-other-guid", 0)
				testHelper.CreateValidActualLRP("some-third-guid", 0)
			})

			It("errors", func() {
				_, err := etcdDB.ActualLRPGroupsByProcessGuid(logger, "some-other-guid")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when etcd is not there", func() {
			BeforeEach(func() {
				etcdRunner.Stop()
			})

			AfterEach(func() {
				etcdRunner.Start()
			})

			It("errors", func() {
				_, err := etcdDB.ActualLRPGroupsByProcessGuid(logger, "some-other-guid")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ActualLRPGroupsByProcessGuidAndIndex", func() {
		Context("when there are both /instance and /evacuating LRPs", func() {
			BeforeEach(func() {
				testHelper.SetRawActualLRP(baseLRP)
				testHelper.SetRawEvacuatingActualLRP(evacuatingLRP, noExpirationTTL)
				testHelper.SetRawActualLRP(otherDomainLRP)
				testHelper.SetRawEvacuatingActualLRP(otherIndexLRP, noExpirationTTL)
				testHelper.SetRawActualLRP(otherCellIdLRP)
			})

			It("returns the /instance LRPs and /evacuating LRPs in the group", func() {
				actualLRPGroup, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, baseProcessGuid, baseIndex)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualLRPGroup).To(Equal(&models.ActualLRPGroup{Instance: baseLRP, Evacuating: evacuatingLRP}))
			})
		})

		Context("when there are no LRPs", func() {
			It("returns an error", func() {
				_, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, baseProcessGuid, baseIndex)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(models.ErrResourceNotFound))
			})
		})

		Context("when the root node exists with no child nodes", func() {
			BeforeEach(func() {
				testHelper.SetRawActualLRP(baseLRP)

				processGuid := baseLRP.ActualLRPKey.GetProcessGuid()
				_, err := etcdClient.Delete(ActualLRPSchemaPath(processGuid, baseIndex), true)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error", func() {
				_, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, baseProcessGuid, baseIndex)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(models.ErrResourceNotFound))
			})
		})

		Context("when there is invalid data", func() {
			BeforeEach(func() {
				testHelper.CreateValidActualLRP("some-guid", 0)
				testHelper.CreateMalformedActualLRP("some-other-guid", 0)
				testHelper.CreateValidActualLRP("some-third-guid", 0)
			})

			It("errors", func() {
				_, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, "some-other-guid", 0)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("ClaimActualLRP", func() {
		var (
			actualLRP        *models.ActualLRP
			claimedActualLRP *models.ActualLRP
			claimErr         *models.Error

			lrpKey      models.ActualLRPKey
			instanceKey models.ActualLRPInstanceKey

			processGuid string
			index       int32
			domain      string
		)

		JustBeforeEach(func() {
			claimedActualLRP, claimErr = etcdDB.ClaimActualLRP(logger, processGuid, index, &instanceKey)
		})

		Context("when the actual LRP exists", func() {
			BeforeEach(func() {
				processGuid = "some-process-guid"
				index = 1
				domain = "some-domain"

				lrpKey = models.NewActualLRPKey(processGuid, index, domain)
				actualLRP = &models.ActualLRP{
					ActualLRPKey: lrpKey,
					State:        models.ActualLRPStateUnclaimed,
					Since:        clock.Now().UnixNano(),
				}

				testHelper.SetRawActualLRP(actualLRP)
			})

			Context("when the instance key is invalid", func() {
				BeforeEach(func() {
					instanceKey = models.NewActualLRPInstanceKey(
						"", // invalid InstanceGuid
						cellID,
					)
				})

				It("returns a validation error", func() {
					Expect(claimErr.Type).To(Equal(models.InvalidRecord))
				})

				It("does not modify the persisted actual LRP", func() {
					lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
					Expect(err).NotTo(HaveOccurred())

					Expect(lrpGroupInBBS.Instance.State).To(Equal(models.ActualLRPStateUnclaimed))
				})
			})

			Context("when the existing ActualLRP is Unclaimed", func() {
				BeforeEach(func() {
					instanceKey = models.NewActualLRPInstanceKey("some-instance-guid", cellID)
				})

				It("does not error", func() {
					Expect(claimErr).NotTo(HaveOccurred())
				})

				It("claims the actual LRP", func() {
					lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
					Expect(err).NotTo(HaveOccurred())

					Expect(lrpGroupInBBS.Instance.State).To(Equal(models.ActualLRPStateClaimed))
				})

				It("updates the ModificationIndex", func() {
					lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
					Expect(err).NotTo(HaveOccurred())

					Expect(lrpGroupInBBS.Instance.ModificationTag.Index).To(Equal(actualLRP.ModificationTag.Index + 1))
				})
			})

			Context("when the existing ActualLRP is Claimed", func() {
				var instanceGuid string

				BeforeEach(func() {
					instanceGuid = "some-instance-guid"
					instanceKey := models.NewActualLRPInstanceKey(instanceGuid, cellID)
					_, err := etcdDB.ClaimActualLRP(logger, processGuid, index, &instanceKey)
					Expect(err).NotTo(HaveOccurred())
				})

				Context("with the same cell and instance guid", func() {
					var previousTime int64

					BeforeEach(func() {
						instanceKey = models.NewActualLRPInstanceKey(instanceGuid, cellID)

						previousTime = clock.Now().UnixNano()
						clock.IncrementBySeconds(1)
					})

					It("does not return an error", func() {
						Expect(claimErr).NotTo(HaveOccurred())
					})

					It("does not alter the state of the persisted LRP", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.State).To(Equal(models.ActualLRPStateClaimed))
					})

					It("does not update the timestamp of the persisted actual lrp", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.Since).To(Equal(previousTime))
					})
				})

				Context("with a different cell", func() {
					BeforeEach(func() {
						instanceKey = models.NewActualLRPInstanceKey(instanceGuid, "another-cell-id")
					})

					It("returns an error", func() {
						Expect(claimErr).To(Equal(models.ErrActualLRPCannotBeClaimed))
					})

					It("does not alter the existing LRP", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.CellId).To(Equal(cellID))
					})
				})

				Context("when the instance guid differs", func() {
					BeforeEach(func() {
						instanceKey = models.NewActualLRPInstanceKey("another-instance-guid", cellID)
					})

					It("returns an error", func() {
						Expect(claimErr).To(Equal(models.ErrActualLRPCannotBeClaimed))
					})

					It("does not alter the existing actual", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.InstanceGuid).To(Equal(instanceGuid))
					})
				})
			})

			Context("when the existing ActualLRP is Running", func() {
				var instanceGuid string

				BeforeEach(func() {
					instanceGuid = "some-instance-guid"
					instanceKey = models.NewActualLRPInstanceKey(instanceGuid, cellID)

					actualLRP.State = models.ActualLRPStateRunning
					actualLRP.ActualLRPInstanceKey = instanceKey
					actualLRP.ActualLRPNetInfo = models.ActualLRPNetInfo{Address: "1"}

					Expect(actualLRP.Validate()).NotTo(HaveOccurred())

					testHelper.SetRawActualLRP(actualLRP)
				})

				Context("with the same cell and instance guid", func() {
					BeforeEach(func() {
						instanceKey = models.NewActualLRPInstanceKey(instanceGuid, cellID)
					})

					It("does not return an error", func() {
						Expect(claimErr).NotTo(HaveOccurred())
					})

					It("reverts the persisted LRP to the CLAIMED state", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.State).To(Equal(models.ActualLRPStateClaimed))
					})

					It("clears the net info", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.Address).To(BeEmpty())
						Expect(lrpGroupInBBS.Instance.Ports).To(BeEmpty())
					})
				})

				Context("with a different cell", func() {
					BeforeEach(func() {
						instanceKey = models.NewActualLRPInstanceKey(instanceGuid, "another-cell-id")
					})

					It("returns an error", func() {
						Expect(claimErr).To(Equal(models.ErrActualLRPCannotBeClaimed))
					})

					It("does not alter the existing LRP", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.CellId).To(Equal(cellID))
					})
				})

				Context("when the instance guid differs", func() {
					BeforeEach(func() {
						instanceKey = models.NewActualLRPInstanceKey("another-instance-guid", cellID)
					})

					It("returns an error", func() {
						Expect(claimErr).To(Equal(models.ErrActualLRPCannotBeClaimed))
					})

					It("does not alter the existing actual", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.InstanceGuid).To(Equal(instanceGuid))
					})
				})
			})

			Context("when there is a placement error", func() {
				BeforeEach(func() {
					instanceKey = models.NewActualLRPInstanceKey("some-instance-guid", cellID)
					actualLRP.PlacementError = "insufficient resources"
					testHelper.SetRawActualLRP(actualLRP)
				})

				It("should clear placement error", func() {
					Expect(claimErr).NotTo(HaveOccurred())
					Expect(claimedActualLRP.PlacementError).To(BeEmpty())
					lrp, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
					Expect(err).NotTo(HaveOccurred())
					Expect(lrp.Instance.PlacementError).To(BeEmpty())
				})
			})
		})

		Context("when the actual LRP does not exist", func() {
			BeforeEach(func() {
				// Do not make a lrp.
			})

			It("cannot claim the LRP", func() {
				Expect(claimErr).To(Equal(models.ErrResourceNotFound))
			})

			It("does not create an actual LRP", func() {
				_, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, "process-guid", 1)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("StartActualLRP", func() {
		var (
			startErr         *models.Error
			request          models.StartActualLRPRequest
			startedActualLRP *models.ActualLRP

			lrpKey      models.ActualLRPKey
			instanceKey models.ActualLRPInstanceKey
			netInfo     models.ActualLRPNetInfo
		)

		JustBeforeEach(func() {
			request.ActualLrpKey = &lrpKey
			request.ActualLrpInstanceKey = &instanceKey
			request.ActualLrpNetInfo = &netInfo
			startedActualLRP, startErr = etcdDB.StartActualLRP(logger, &request)
		})

		Context("when the actual LRP exists", func() {
			var (
				processGuid string
				index       int32
				actualLRP   *models.ActualLRP
			)

			BeforeEach(func() {
				index = 1
				processGuid = "some-process-guid"
				key := models.NewActualLRPKey(processGuid, index, "domain")
				actualLRP = &models.ActualLRP{
					ActualLRPKey: key,
					State:        models.ActualLRPStateUnclaimed,
					Since:        123,
				}

				testHelper.SetRawActualLRP(actualLRP)
			})

			Context("when the existing ActualLRP is Unclaimed", func() {
				BeforeEach(func() {
					lrpKey = actualLRP.ActualLRPKey
					instanceKey = models.NewActualLRPInstanceKey("some-guid", cellID)
					netInfo = models.NewActualLRPNetInfo("1.2.3.4", models.NewPortMapping(5678, 1234))
				})

				It("does not error", func() {
					Expect(startErr).NotTo(HaveOccurred())
				})

				It("starts the actual LRP", func() {
					lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
					Expect(err).NotTo(HaveOccurred())

					Expect(lrpGroupInBBS.Instance.State).To(Equal(models.ActualLRPStateRunning))
				})

				Context("when there is a placement error", func() {
					BeforeEach(func() {
						actualLRP.PlacementError = "insufficient resources"
						testHelper.SetRawActualLRP(actualLRP)
					})

					It("should clear placement error", func() {
						Expect(startErr).NotTo(HaveOccurred())
						Expect(startedActualLRP.PlacementError).To(BeEmpty())
						lrp, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())
						Expect(lrp.Instance.PlacementError).To(BeEmpty())
					})
				})
			})

			Context("when the domain differs", func() {
				BeforeEach(func() {
					lrpKey = actualLRP.ActualLRPKey
					lrpKey.Domain = "some-other-domain"
					instanceKey = models.NewActualLRPInstanceKey("some-guid", cellID)
					netInfo = models.NewActualLRPNetInfo("1.2.3.4", models.NewPortMapping(5678, 1234))
				})

				It("returns an error", func() {
					Expect(startErr).To(Equal(models.ErrActualLRPCannotBeStarted))
				})

				It("does not modify the persisted actual LRP", func() {
					lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
					Expect(err).NotTo(HaveOccurred())

					Expect(lrpGroupInBBS.Instance.State).To(Equal(models.ActualLRPStateUnclaimed))
				})
			})

			Context("when the existing ActualLRP is Claimed", func() {
				var instanceGuid string

				BeforeEach(func() {
					instanceGuid = "some-instance-guid"
					instanceKey := models.NewActualLRPInstanceKey(instanceGuid, cellID)
					_, err := etcdDB.ClaimActualLRP(
						logger,
						actualLRP.ProcessGuid,
						actualLRP.Index,
						&instanceKey,
					)
					Expect(err).NotTo(HaveOccurred())
				})

				Context("with the same cell and instance guid", func() {
					BeforeEach(func() {
						lrpKey = actualLRP.ActualLRPKey
						instanceKey = models.NewActualLRPInstanceKey(instanceGuid, cellID)
						netInfo = models.NewActualLRPNetInfo("1.2.3.4", models.NewPortMapping(5678, 1234))
					})

					It("does not return an error", func() {
						Expect(startErr).NotTo(HaveOccurred())
					})

					It("promotes the persisted LRP to RUNNING", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.State).To(Equal(models.ActualLRPStateRunning))
					})
				})

				Context("with a different cell", func() {
					BeforeEach(func() {
						lrpKey = actualLRP.ActualLRPKey
						instanceKey = models.NewActualLRPInstanceKey(instanceGuid, "another-cell-id")
						netInfo = models.NewActualLRPNetInfo("1.2.3.4", models.NewPortMapping(5678, 1234))
					})

					It("does not return an error", func() {
						Expect(startErr).NotTo(HaveOccurred())
					})

					It("promotes the persisted LRP to RUNNING", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.State).To(Equal(models.ActualLRPStateRunning))
					})
				})

				Context("when the instance guid differs", func() {
					BeforeEach(func() {
						lrpKey = actualLRP.ActualLRPKey
						instanceKey = models.NewActualLRPInstanceKey("another-instance-guid", cellID)
						netInfo = models.NewActualLRPNetInfo("1.2.3.4", models.NewPortMapping(5678, 1234))
					})

					It("does not return an error", func() {
						Expect(startErr).NotTo(HaveOccurred())
					})

					It("promotes the persisted LRP to RUNNING", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.State).To(Equal(models.ActualLRPStateRunning))
					})
				})
			})

			Context("when the existing ActualLRP is Running", func() {
				var instanceGuid string

				BeforeEach(func() {
					instanceGuid = "some-instance-guid"

					existingLRPRequest := &models.StartActualLRPRequest{}
					existingLRPRequest.ActualLrpKey = &actualLRP.ActualLRPKey
					existingInstanceKey := models.NewActualLRPInstanceKey(instanceGuid, cellID)
					existingLRPRequest.ActualLrpInstanceKey = &existingInstanceKey
					existingNetInfo := models.NewActualLRPNetInfo("1.2.3.4", models.NewPortMapping(5678, 1234))
					existingLRPRequest.ActualLrpNetInfo = &existingNetInfo

					_, err := etcdDB.StartActualLRP(logger, existingLRPRequest)
					Expect(err).NotTo(HaveOccurred())
				})

				Context("with the same cell and instance guid", func() {
					BeforeEach(func() {
						lrpKey = actualLRP.ActualLRPKey
						instanceKey = models.NewActualLRPInstanceKey(instanceGuid, cellID)
						netInfo = models.NewActualLRPNetInfo("5.6.7.8", models.NewPortMapping(4321, 4567))
					})

					It("does not return an error", func() {
						Expect(startErr).NotTo(HaveOccurred())
					})

					It("does not alter the state of the persisted LRP", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.State).To(Equal(models.ActualLRPStateRunning))
					})

					It("updates the net info", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.ActualLRPNetInfo).To(Equal(netInfo))
					})

					Context("and the same net info", func() {
						var previousTime int64
						BeforeEach(func() {
							netInfo = models.NewActualLRPNetInfo("1.2.3.4", models.NewPortMapping(5678, 1234))

							previousTime = clock.Now().UnixNano()
							clock.IncrementBySeconds(1)
						})

						It("does not update the timestamp of the persisted actual lrp", func() {
							lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
							Expect(err).NotTo(HaveOccurred())

							Expect(lrpGroupInBBS.Instance.Since).To(Equal(previousTime))
						})
					})
				})

				Context("with a different cell", func() {
					BeforeEach(func() {
						lrpKey = actualLRP.ActualLRPKey
						instanceKey = models.NewActualLRPInstanceKey(instanceGuid, "another-cell-id")
						netInfo = models.NewActualLRPNetInfo("5.6.7.8", models.NewPortMapping(4321, 4567))
					})

					It("returns an error", func() {
						Expect(startErr).To(Equal(models.ErrActualLRPCannotBeStarted))
					})

					It("does not alter the existing LRP", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.CellId).To(Equal(cellID))
					})
				})

				Context("when the instance guid differs", func() {
					BeforeEach(func() {
						lrpKey = actualLRP.ActualLRPKey
						instanceKey = models.NewActualLRPInstanceKey("another-instance-guid", cellID)
						netInfo = models.NewActualLRPNetInfo("5.6.7.8", models.NewPortMapping(4321, 4567))
					})

					It("returns an error", func() {
						Expect(startErr).To(Equal(models.ErrActualLRPCannotBeStarted))
					})

					It("does not alter the existing actual", func() {
						lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
						Expect(err).NotTo(HaveOccurred())

						Expect(lrpGroupInBBS.Instance.InstanceGuid).To(Equal(instanceGuid))
					})
				})
			})
		})

		Context("when the actual LRP does not exist", func() {
			BeforeEach(func() {
				lrpKey = models.NewActualLRPKey("process-guid", 1, "domain")
				instanceKey = models.NewActualLRPInstanceKey("instance-guid", cellID)
				netInfo = models.NewActualLRPNetInfo("1.2.3.4", models.NewPortMapping(5678, 1234))
			})

			It("starts the LRP", func() {
				Expect(startErr).NotTo(HaveOccurred())
			})

			It("sets the State", func() {
				lrpGroup, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, "process-guid", 1)
				Expect(err).NotTo(HaveOccurred())

				Expect(lrpGroup.Instance.State).To(Equal(models.ActualLRPStateRunning))
			})

			It("sets the ModificationTag", func() {
				lrpGroup, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, "process-guid", 1)
				Expect(err).NotTo(HaveOccurred())

				Expect(lrpGroup.Instance.ModificationTag.Epoch).NotTo(BeEmpty())
				Expect(lrpGroup.Instance.ModificationTag.Index).To(BeEquivalentTo(0))
			})
		})
	})

	Describe("FailActualLRP", func() {
		var (
			request   models.FailActualLRPRequest
			failErr   *models.Error
			actualLRP *models.ActualLRP

			lrpKey       models.ActualLRPKey
			instanceKey  models.ActualLRPInstanceKey
			errorMessage string

			processGuid string
			index       int32
			domain      string
		)

		JustBeforeEach(func() {
			request.ActualLrpKey = &lrpKey
			request.ErrorMessage = errorMessage
			failErr = etcdDB.FailActualLRP(logger, &request)
		})

		Context("when the actual LRP exists", func() {
			BeforeEach(func() {
				processGuid = "some-process-guid"
				index = 1
				domain = "some-domain"
				errorMessage = "some-error"

				lrpKey = models.NewActualLRPKey(processGuid, index, domain)
				actualLRP = &models.ActualLRP{
					ActualLRPKey: lrpKey,
					State:        models.ActualLRPStateUnclaimed,
					Since:        clock.Now().UnixNano(),
				}

				testHelper.SetRawActualLRP(actualLRP)
			})

			Context("when the existing ActualLRP is Unclaimed", func() {
				It("does not error", func() {
					Expect(failErr).NotTo(HaveOccurred())
				})

				It("fails the actual LRP", func() {
					lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
					Expect(err).NotTo(HaveOccurred())

					Expect(lrpGroupInBBS.Instance.PlacementError).To(Equal(errorMessage))
				})

				It("updates the ModificationIndex", func() {
					lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
					Expect(err).NotTo(HaveOccurred())

					Expect(lrpGroupInBBS.Instance.ModificationTag.Index).To(Equal(actualLRP.ModificationTag.Index + 1))
				})
			})

			Context("when the existing ActualLRP is not Unclaimed", func() {
				var instanceGuid string

				BeforeEach(func() {
					instanceGuid = "some-instance-guid"
					instanceKey = models.NewActualLRPInstanceKey(instanceGuid, cellID)

					actualLRP.State = models.ActualLRPStateRunning
					actualLRP.ActualLRPInstanceKey = instanceKey
					actualLRP.ActualLRPNetInfo = models.ActualLRPNetInfo{Address: "1"}

					Expect(actualLRP.Validate()).NotTo(HaveOccurred())

					testHelper.SetRawActualLRP(actualLRP)
				})

				It("returns an error", func() {
					Expect(failErr).To(Equal(models.ErrActualLRPCannotBeFailed))
				})
			})
		})

		Context("when the actual LRP does not exist", func() {
			BeforeEach(func() {
				// Do not make a lrp.
			})

			It("cannot fail the LRP", func() {
				Expect(failErr).To(Equal(models.ErrResourceNotFound))
			})
		})
	})

	Describe("RemoveActualLRP", func() {
		var (
			actualLRP *models.ActualLRP
			removeErr *models.Error

			lrpKey models.ActualLRPKey

			processGuid string
			index       int32
			domain      string
		)

		JustBeforeEach(func() {
			removeErr = etcdDB.RemoveActualLRP(logger, processGuid, index)
		})

		Context("when the actual LRP exists", func() {
			BeforeEach(func() {
				processGuid = "some-process-guid"
				index = 1
				domain = "some-domain"

				lrpKey = models.NewActualLRPKey(processGuid, index, domain)
				actualLRP = &models.ActualLRP{
					ActualLRPKey: lrpKey,
					State:        models.ActualLRPStateUnclaimed,
					Since:        clock.Now().UnixNano(),
				}

				testHelper.SetRawActualLRP(actualLRP)
			})

			It("does not error", func() {
				Expect(removeErr).NotTo(HaveOccurred())
			})

			It("removes the actual LRP", func() {
				lrpGroupInBBS, err := etcdDB.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(models.ErrResourceNotFound))
				Expect(lrpGroupInBBS).To(BeNil())
			})
		})

		Context("when the actual LRP does not exist", func() {
			BeforeEach(func() {
				// Do not make a lrp.
			})

			It("cannot remove the LRP", func() {
				Expect(removeErr).To(Equal(models.ErrResourceNotFound))
			})
		})
	})
})