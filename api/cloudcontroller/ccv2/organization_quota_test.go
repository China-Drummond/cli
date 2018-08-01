package ccv2_test

import (
	"net/http"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"
	. "code.cloudfoundry.org/cli/api/cloudcontroller/ccv2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

var _ = Describe("OrganizationQuota", func() {
	var client *Client

	BeforeEach(func() {
		client = NewTestClient()
	})

	Describe("GetOrganizationQuota", func() {

		Context("when getting the organization quota does not return an error", func() {
			BeforeEach(func() {
				response := `{
				"metadata": {
					"guid": "some-org-quota-guid"
				},
				"entity": {
					"name": "some-org-quota"
				}
			}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v2/quota_definitions/some-org-quota-guid"),
						RespondWith(http.StatusOK, response, http.Header{"X-Cf-Warnings": {"warning-1"}}),
					),
				)
			})

			It("returns the organization quota", func() {
				orgQuota, warnings, err := client.GetOrganizationQuota("some-org-quota-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(warnings).To(Equal(Warnings{"warning-1"}))
				Expect(orgQuota).To(Equal(OrganizationQuota{
					GUID: "some-org-quota-guid",
					Name: "some-org-quota",
				}))
			})
		})

		Context("when the organization quota returns an error", func() {
			BeforeEach(func() {
				response := `{
				  "description": "Quota Definition could not be found: some-org-quota-guid",
				  "error_code": "CF-QuotaDefinitionNotFound",
				  "code": 240001
				}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v2/quota_definitions/some-org-quota-guid"),
						RespondWith(http.StatusNotFound, response, http.Header{"X-Cf-Warnings": {"warning-1"}}),
					),
				)
			})

			It("returns the error", func() {
				_, warnings, err := client.GetOrganizationQuota("some-org-quota-guid")
				Expect(err).To(MatchError(ccerror.ResourceNotFoundError{
					Message: "Quota Definition could not be found: some-org-quota-guid",
				}))
				Expect(warnings).To(Equal(Warnings{"warning-1"}))
			})
		})

	})

	Describe("GetOrganizationQuotaByName", func() {
		var (
			quota      OrganizationQuota
			warnings   Warnings
			executeErr error
		)

		JustBeforeEach(func() {
			quota, warnings, executeErr = client.GetOrganizationQuotaByName("some-quota-name")
		})

		Context("when quota exists", func() {
			BeforeEach(func() {
				response := `{
				"metadata": {
					"guid": "some-org-quota-guid"
				},
				"entity": {
					"name": "some-quota-name"
				}
			}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v2/quota_definitions", "q=name:some-quota-name"),
						RespondWith(http.StatusOK, response, http.Header{"X-Cf-Warnings": {"warning-1"}}),
					),
				)
			})

			It("returns the organization quota", func() {
				Expect(executeErr).NotTo(HaveOccurred())
				Expect(warnings).To(Equal(Warnings{"warning-1"}))
				Expect(quota).To(Equal(OrganizationQuota{
					GUID: "some-org-quota-guid",
					Name: "some-quota-name",
				}))

			})
		})

		Context("when quota doesn't exist", func() {
			BeforeEach(func() {
				response := `{
					"code": 10001,
					"description": "Some Error",
					"error_code": "CF-SomeError"
				}`
				server.AppendHandlers(
					CombineHandlers(
						VerifyRequest(http.MethodGet, "/v2/quota_definitions", "q=name:some-quota-name"),
						RespondWith(http.StatusNotFound, response, http.Header{"X-Cf-Warnings": {"warning-1"}}),
					),
				)
			})

			It("returns warnings and errors", func() {
				Expect(executeErr).To(MatchError(ccerror.ResourceNotFoundError{
					Message: "Some Error",
				}))
				Expect(warnings).To(Equal(Warnings{"warning-1"}))
			})
		})
	})
})
