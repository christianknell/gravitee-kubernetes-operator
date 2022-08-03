/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// +kubebuilder:docs-gen:collapse=Apache License
package controllers

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gio "github.com/gravitee-io/gravitee-kubernetes-operator/api/v1alpha1"
	"github.com/gravitee-io/gravitee-kubernetes-operator/test"
)

// +kubebuilder:docs-gen:collapse=Imports

var _ = Describe("API Definition Controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		namespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	ctx := context.Background()
	httpClient := http.Client{Timeout: 5 * time.Second}

	AfterEach(func() {
		// Delete the API definition
		Eventually(func() error {
			return k8sClient.DeleteAllOf(ctx, new(gio.ApiDefinition), &client.DeleteAllOfOptions{
				ListOptions:   client.ListOptions{Namespace: namespace},
				DeleteOptions: client.DeleteOptions{},
			})
		}).ShouldNot(HaveOccurred())

		// Delete the ManagementContext
		Eventually(func() error {
			return k8sClient.DeleteAllOf(ctx, new(gio.ManagementContext), &client.DeleteAllOfOptions{
				ListOptions:   client.ListOptions{Namespace: namespace},
				DeleteOptions: client.DeleteOptions{},
			})
		}).ShouldNot(HaveOccurred())
	})

	Context("API definition Resource", func() {

		It("Should update an API Definition", func() {
			By("Create an API definition resource without a management context")
			const apiDefinitionSample = "../config/samples/apim/basic-example.yml"

			apiDefinitionFixture, err := test.NewApiDefinition(apiDefinitionSample)
			Expect(err).ToNot(HaveOccurred())

			Expect(k8sClient.Create(ctx, apiDefinitionFixture)).Should(Succeed())

			apiLookupKey := types.NamespacedName{Name: apiDefinitionFixture.Name, Namespace: namespace}
			createdApiDefinition := new(gio.ApiDefinition)

			Eventually(func() error {
				return k8sClient.Get(ctx, apiLookupKey, createdApiDefinition)
			}, timeout, interval).ShouldNot(HaveOccurred())

			By("Call initial API definition URL and expect no error")

			// Check created api is callable
			var endpointInitial = test.GatewayUrl + createdApiDefinition.Spec.Proxy.VirtualHosts[0].Path

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpointInitial)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())

			By("Update the context path in API definition and expect no error")

			updatedApiDefinition := createdApiDefinition.DeepCopy()

			expectedPath := updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path + "-updated"
			updatedApiDefinition.Spec.Proxy.VirtualHosts[0].Path = expectedPath

			err = k8sClient.Update(ctx, updatedApiDefinition)
			Expect(err).ToNot(HaveOccurred())

			By("Call updated API definition URL and expect no error")

			var endpointUpdated = test.GatewayUrl + expectedPath

			Eventually(func() bool {
				res, callErr := httpClient.Get(endpointUpdated)
				return callErr == nil && res.StatusCode == 200
			}, timeout, interval).Should(BeTrue())
		})
	})
})