// Package config_e2e_test provides config command specific E2E test cases
package config_e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/cluster-api/util"

	"github.com/vmware-tanzu/tanzu-framework/cli/core/test/e2e/framework"
)

var _ = framework.CLICoreDescribe("[Tests:E2E][Feature:Command-Config]", func() {
	var (
		tf *framework.Framework
	)
	Context("config feature flag operations", func() {
		BeforeEach(func() {
			tf = framework.NewFramework()
		})
		When("new config flag set with value", func() {
			It("should set flag and unset flag successfully", func() {
				randomFlagName := RandomNameWithE2ETestString()
				randomFeatureFlagPath := "features.global." + randomFlagName
				flagVal := "true"
				err := tf.ConfigSetFeature(randomFeatureFlagPath, flagVal)
				Expect(err).To(BeNil())

				val, err := tf.ConfigFeatureFlagExists(randomFlagName, flagVal)
				Expect(err).To(BeNil())
				Expect(val).To(BeTrue())

				err = tf.ConfigUnsetFeature(randomFeatureFlagPath)
				Expect(err).To(BeNil())

				val, err = tf.ConfigFeatureFlagNotExists(randomFlagName)
				Expect(err).To(BeNil())
				Expect(val).To(BeTrue())
			})
		})
		When("config init called", func() {
			It("should initialize configuration successfully", func() {
				err := tf.ConfigInit()
				Expect(err).To(BeNil())
			})
		})

		// TODO: test list and delete use cases as part of context or login commands
		When("config server list called", func() {
			It("should list all available servers", func() {
				err := tf.ConfigServerList()
				Expect(err).To(BeNil())
			})
		})
	})
})

func RandomNameWithE2ETestString() string {
	return "e2e-test-" + util.RandomString(4)
}
