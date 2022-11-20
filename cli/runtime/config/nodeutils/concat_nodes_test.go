// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package nodeutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConcatNode(t *testing.T) {
	tests := []struct {
		name   string
		cfg    string
		cfg2   string
		output string
	}{
		{
			name: "success concat src into empty dst node",
			cfg: `apiVersion: config.tanzu.vmware.com/v1alpha1
clientOptions:
  cli:
    bomRepo: projects.registry.vmware.com/tkg
    compatibilityFilePath: tkg-compatibility
    discoverySources:
      - contextType: k8s
        local:
          name: default-local
          path: standalone
      - local:
          name: admin-local
          path: admin
    edition: tkg
  features:
    cluster:
      custom-nameservers: 'false'
      dual-stack-ipv4-primary: 'false'
      dual-stack-ipv6-primary: 'false'
    global:
      context-aware-cli-for-plugins: 'true'
      context-target: 'false'
      tkr-version-v1alpha3-beta: 'false'
    management-cluster:
      aws-instance-types-exclude-arm: 'true'
      custom-nameservers: 'false'
      dual-stack-ipv4-primary: 'false'
      dual-stack-ipv6-primary: 'false'
      export-from-confirm: 'true'
      import: 'false'
      standalone-cluster-mode: 'false'
    package:
      kctrl-package-command-tree: 'true'
contexts:
  - name: test-mc
    type: k8s
    group: one
    clusterOpts:
      isManagementCluster: true
      annotation: one
      required: true
      annotationStruct:
        one: one
      endpoint: cfg-test-endpoint
      path: cfg-test-path
      context: cfg-test-context
    discoverySources:
      - gcp:
          name: test
          bucket: cfg-test-bucket
          manifestPath: cfg-test-manifest-path
          annotation: one
          required: true
        contextType: tmc
      - gcp:
          name: test-two
          bucket: cfg-test-bucket
          manifestPath: cfg-test-manifest-path
          annotation: two
          required: true
        contextType: tmc
currentContext:
  k8s: test-mc
kind: ClientConfig
metadata:
  creationTimestamp: null`,
			cfg2: ``,
			output: `apiVersion: config.tanzu.vmware.com/v1alpha1
clientOptions:
    cli:
        bomRepo: projects.registry.vmware.com/tkg
        compatibilityFilePath: tkg-compatibility
        discoverySources:
            - contextType: k8s
              local:
                name: default-local
                path: standalone
            - local:
                name: admin-local
                path: admin
        edition: tkg
    features:
        cluster:
            custom-nameservers: 'false'
            dual-stack-ipv4-primary: 'false'
            dual-stack-ipv6-primary: 'false'
        global:
            context-aware-cli-for-plugins: 'true'
            context-target: 'false'
            tkr-version-v1alpha3-beta: 'false'
        management-cluster:
            aws-instance-types-exclude-arm: 'true'
            custom-nameservers: 'false'
            dual-stack-ipv4-primary: 'false'
            dual-stack-ipv6-primary: 'false'
            export-from-confirm: 'true'
            import: 'false'
            standalone-cluster-mode: 'false'
        package:
            kctrl-package-command-tree: 'true'
contexts:
    - name: test-mc
      type: k8s
      group: one
      clusterOpts:
        isManagementCluster: true
        annotation: one
        required: true
        annotationStruct:
            one: one
        endpoint: cfg-test-endpoint
        path: cfg-test-path
        context: cfg-test-context
      discoverySources:
        - gcp:
            name: test
            bucket: cfg-test-bucket
            manifestPath: cfg-test-manifest-path
            annotation: one
            required: true
          contextType: tmc
        - gcp:
            name: test-two
            bucket: cfg-test-bucket
            manifestPath: cfg-test-manifest-path
            annotation: two
            required: true
          contextType: tmc
currentContext:
    k8s: test-mc
kind: ClientConfig
metadata:
    creationTimestamp: null
`,
		},
		{
			name: "success concat src into dst node",
			cfg: `apiVersion: config.tanzu.vmware.com/v1alpha1
clientOptions:
  cli:
    bomRepo: projects.registry.vmware.com/tkg
    compatibilityFilePath: tkg-compatibility
    discoverySources:
      - contextType: k8s
        local:
          name: default-local
          path: standalone
      - local:
          name: admin-local
          path: admin
    edition: tkg
  features:
    cluster:
      custom-nameservers: 'false'
      dual-stack-ipv4-primary: 'false'
      dual-stack-ipv6-primary: 'false'
    global:
      context-aware-cli-for-plugins: 'true'
      context-target: 'false'
      tkr-version-v1alpha3-beta: 'false'
    management-cluster:
      aws-instance-types-exclude-arm: 'true'
      custom-nameservers: 'false'
      dual-stack-ipv4-primary: 'false'
      dual-stack-ipv6-primary: 'false'
      export-from-confirm: 'true'
      import: 'false'
      standalone-cluster-mode: 'false'
    package:
      kctrl-package-command-tree: 'true'
contexts:
  - name: test-mc
    type: k8s
    group: one
    clusterOpts:
      isManagementCluster: true
      annotation: one
      required: true
      annotationStruct:
        one: one
      endpoint: cfg-test-endpoint
      path: cfg-test-path
      context: cfg-test-context
    discoverySources:
      - gcp:
          name: test
          bucket: cfg-test-bucket
          manifestPath: cfg-test-manifest-path
          annotation: one
          required: true
        contextType: tmc
      - gcp:
          name: test-two
          bucket: cfg-test-bucket
          manifestPath: cfg-test-manifest-path
          annotation: two
          required: true
        contextType: tmc
currentContext:
  k8s: test-mc
kind: ClientConfig
metadata:
  creationTimestamp: null`,
			cfg2: `
contexts:
  - name: test-mc
    type: k8s
    group: one
    clusterOpts:
      isManagementCluster: true
      annotation: one
      required: true
      annotationStruct:
        one: one
      endpoint: cfg2-test-endpoint
      path: cfg2-test-path
      context: cfg2-test-context
    discoverySources:
      - gcp:
          name: test
          bucket: cfg2-test-bucket
          manifestPath: test-manifest-path
          annotation: one
          required: true
        contextType: tmc
      - gcp:
          name: test-two
          bucket: test-bucket
          manifestPath: test-manifest-path
          annotation: two
          required: true
        contextType: tmc
currentContext:
  k8s: test-mc`,
			output: `contexts:
    - name: test-mc
      type: k8s
      group: one
      clusterOpts:
        isManagementCluster: true
        annotation: one
        required: true
        annotationStruct:
            one: one
        endpoint: cfg2-test-endpoint
        path: cfg2-test-path
        context: cfg2-test-context
      discoverySources:
        - gcp:
            name: test
            bucket: cfg2-test-bucket
            manifestPath: test-manifest-path
            annotation: one
            required: true
          contextType: tmc
        - gcp:
            name: test-two
            bucket: test-bucket
            manifestPath: test-manifest-path
            annotation: two
            required: true
          contextType: tmc
currentContext:
    k8s: test-mc
apiVersion: config.tanzu.vmware.com/v1alpha1
clientOptions:
    cli:
        bomRepo: projects.registry.vmware.com/tkg
        compatibilityFilePath: tkg-compatibility
        discoverySources:
            - contextType: k8s
              local:
                name: default-local
                path: standalone
            - local:
                name: admin-local
                path: admin
        edition: tkg
    features:
        cluster:
            custom-nameservers: 'false'
            dual-stack-ipv4-primary: 'false'
            dual-stack-ipv6-primary: 'false'
        global:
            context-aware-cli-for-plugins: 'true'
            context-target: 'false'
            tkr-version-v1alpha3-beta: 'false'
        management-cluster:
            aws-instance-types-exclude-arm: 'true'
            custom-nameservers: 'false'
            dual-stack-ipv4-primary: 'false'
            dual-stack-ipv6-primary: 'false'
            export-from-confirm: 'true'
            import: 'false'
            standalone-cluster-mode: 'false'
        package:
            kctrl-package-command-tree: 'true'
kind: ClientConfig
metadata:
    creationTimestamp: null
`,
		},
	}

	for _, spec := range tests {
		t.Run(spec.name, func(t *testing.T) {
			// Setup data
			var src yaml.Node
			var dst yaml.Node
			var err error
			err = yaml.Unmarshal([]byte(spec.cfg), &src)
			assert.NoError(t, err)
			err = yaml.Unmarshal([]byte(spec.cfg2), &dst)
			assert.NoError(t, err)

			// Perform action
			err = ConcatNodes(&src, &dst)
			assert.NoError(t, err)
			data, err := yaml.Marshal(&dst)
			assert.NoError(t, err)
			assert.Equal(t, spec.output, string(data))
		})
	}
}
