// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package linkvalidator_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/gardener/docforge/pkg/nodeplugins/markdown/linkvalidator"
	"github.com/gardener/docforge/pkg/osfakes/httpclient/httpclientfakes"
	"github.com/gardener/docforge/pkg/registry/registryfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJobs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validator Suite")
}

var _ = Describe("Executing Validate", func() {
	var (
		err        error
		httpClient *httpclientfakes.FakeClient
		repository *registryfakes.FakeInterface
		worker     *linkvalidator.ValidatorWorker

		linkDestination   string
		contentSourcePath string
		ctx               context.Context

		hostToReport []string
	)
	BeforeEach(func() {
		httpClient = &httpclientfakes.FakeClient{}
		repository = &registryfakes.FakeInterface{}
		repository.ClientReturns(httpClient)

		ctx = context.Background()
		httpClient.DoReturns(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(""))),
		}, nil)
		linkDestination = "https://repoHost/fake_link"
		contentSourcePath = "fake_path"
		hostToReport = []string{}
	})

	JustBeforeEach(func() {
		worker, err = linkvalidator.NewValidatorWorker(repository, hostToReport)
		Expect(worker).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred())

		err = worker.Validate(ctx, linkDestination, contentSourcePath)
	})

	Context("localhost", func() {
		BeforeEach(func() {
			linkDestination = "https://127.0.0.1/fake_link"

		})
		It("skips link validation", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(0))
		})
	})
	// FContext("url is not valid", func() {
	// 	BeforeEach(func() {
	// 		Expect(err).NotTo(HaveOccurred())
	// 		linkDestination = "https://invalid_host"

	// 	})
	// 	It("fails", func() {
	// 		Expect(err).To(HaveOccurred())
	// 		Expect(err.Error()).To(ContainSubstring("no sutiable repository host for https://invalid_host"))
	// 		Expect(httpClient.DoCallCount()).To(Equal(0))
	// 	})
	// })
	Context("http client returns errors", func() {
		BeforeEach(func() {
			httpClient.DoReturnsOnCall(0, nil, errors.New("fake_error"))
		})
		It("succeeded", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(1))
		})
	})
	Context("http client returns StatusTooManyRequests", func() {
		BeforeEach(func() {
			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil)
		})
		It("retries on StatusTooManyRequests", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(2))
		})
	})
	Context("http client returns StatusUnauthorized", func() {
		BeforeEach(func() {
			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil)
		})
		It("returns on StatusUnauthorized", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(1))
		})
	})
	Context("http client returns error status code", func() {
		BeforeEach(func() {
			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil)
		})
		It("retries on error status code", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(2))
		})
	})
	Context("http client returns error on retry", func() {
		BeforeEach(func() {
			httpClient.DoReturns(nil, errors.New("fake_error"))
			httpClient.DoReturnsOnCall(0, &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil)
		})
		It("succeeded", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(2))
		})
	})
	Context("http client returns error code on retry", func() {
		BeforeEach(func() {
			httpClient.DoReturns(&http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil)
		})
		It("succeeded", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(2))
		})
	})
	When("resource handlers for the link is found", func() {
		var (
			handlerHttpClient *httpclientfakes.FakeClient
		)
		BeforeEach(func() {
			handlerHttpClient = &httpclientfakes.FakeClient{}
			handlerHttpClient.DoReturns(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			}, nil)
			repository.ClientReturns(handlerHttpClient)
		})
		It("uses handler's client", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(httpClient.DoCallCount()).To(Equal(0))
			Expect(handlerHttpClient.DoCallCount()).To(Equal(1))
		})
	})
	It("succeeded", func() {
		Expect(err).NotTo(HaveOccurred())
		Expect(httpClient.DoCallCount()).To(Equal(1))
		req := httpClient.DoArgsForCall(0)
		Expect(req).NotTo(BeNil())
		Expect(req.Host).To(Equal("repoHost"))
	})
})
