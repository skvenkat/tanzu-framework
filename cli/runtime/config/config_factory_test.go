// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	configapi "github.com/vmware-tanzu/tanzu-framework/cli/runtime/apis/config/v1alpha1"
	"gopkg.in/yaml.v3"
)

func setupCfgAndCfgV2Data() (string, string) {
	cfg := `apiVersion: config.tanzu.vmware.com/v1alpha1
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
servers:
  - name: test-mc
    type: managementcluster
    managementClusterOpts:
      endpoint: test-endpoint
      path: test-path
      context: test-context
      annotation: one
      required: true
    discoverySources:
      - gcp:
          name: test
          bucket: test-bucket
          manifestPath: test-manifest-path
          annotation: one
          required: true
        contextType: tmc
current: test-mc
`
	cfgV2 := `
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
      endpoint: test-endpoint
      path: test-path
      context: test-context
    discoverySources:
      - gcp:
          name: test
          bucket: test-bucket
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
`
	return cfg, cfgV2
}

func TestGetClientConfigWithLockAndWithoutLock(t *testing.T) {
	// Setup data
	cfg, cfgV2 := setupCfgAndCfgV2Data()
	// Setup config data
	f1, err := os.CreateTemp("", "tanzu_config")
	assert.Nil(t, err)
	err = os.WriteFile(f1.Name(), []byte(cfg), 0644)
	assert.Nil(t, err)
	err = os.Setenv(EnvConfigKey, f1.Name())
	assert.NoError(t, err)

	// Setup config v2 data
	f2, err := os.CreateTemp("", "tanzu_config_v2")
	assert.Nil(t, err)
	err = os.WriteFile(f2.Name(), []byte(cfgV2), 0644)
	assert.Nil(t, err)
	err = os.Setenv(EnvConfigV2Key, f2.Name())
	assert.NoError(t, err)

	//Setup metadata
	f3, err := os.CreateTemp("", "tanzu_config_metadata")
	assert.Nil(t, err)
	err = os.WriteFile(f3.Name(), []byte(""), 0644)
	assert.Nil(t, err)
	err = os.Setenv(EnvConfigMetadataKey, f3.Name())
	assert.NoError(t, err)

	// Cleanup
	defer func(name string) {
		err = os.Remove(name)
		assert.NoError(t, err)
	}(f1.Name())

	defer func(name string) {
		err = os.Remove(name)
		assert.NoError(t, err)
	}(f2.Name())

	defer func(name string) {
		err = os.Remove(name)
		assert.NoError(t, err)
	}(f3.Name())

	//Actions
	nodeWithLock, err := getClientConfig()

	//Actions
	nodeWithoutLocK, err := getClientConfigNoLock()

	nodes := []*yaml.Node{nodeWithLock, nodeWithoutLocK}
	for _, node := range nodes {
		// Assertions
		assert.NotNil(t, node)
		assert.NoError(t, err)

		expectedCtx := &configapi.Context{
			Name: "test-mc",
			Type: configapi.CtxTypeK8s,
			ClusterOpts: &configapi.ClusterServer{
				Endpoint:            "test-endpoint",
				Path:                "test-path",
				Context:             "test-context",
				IsManagementCluster: true,
			},
			DiscoverySources: []configapi.PluginDiscovery{
				{
					GCP: &configapi.GCPDiscovery{
						Name:         "test",
						Bucket:       "test-bucket",
						ManifestPath: "test-manifest-path",
					},
					ContextType: configapi.CtxTypeTMC,
				},
				{
					GCP: &configapi.GCPDiscovery{
						Name:         "test-two",
						Bucket:       "test-bucket",
						ManifestPath: "test-manifest-path",
					},
					ContextType: configapi.CtxTypeTMC,
				},
			},
		}

		ctx, err := getContext(node, "test-mc")
		assert.NoError(t, err)
		assert.Equal(t, expectedCtx, ctx)

		expectedServer := &configapi.Server{
			Name: "test-mc",
			Type: "managementcluster",
			ManagementClusterOpts: &configapi.ManagementClusterServer{
				Endpoint: "test-endpoint",
				Path:     "test-path",
				Context:  "test-context",
			},
			DiscoverySources: []configapi.PluginDiscovery{
				{
					GCP: &configapi.GCPDiscovery{
						Name:         "test",
						Bucket:       "test-bucket",
						ManifestPath: "test-manifest-path",
					},
					ContextType: configapi.CtxTypeTMC,
				},
			},
		}

		server, err := getServer(node, "test-mc")
		assert.NoError(t, err)
		assert.Equal(t, expectedServer, server)
	}
}

func TestPersistConfig(t *testing.T) {
	// Setup data
	cfg, cfgV2 := setupCfgAndCfgV2Data()
	// Setup config data
	f1, err := os.CreateTemp("", "tanzu_config")
	assert.Nil(t, err)
	err = os.WriteFile(f1.Name(), []byte(cfg), 0644)
	assert.Nil(t, err)
	err = os.Setenv(EnvConfigKey, f1.Name())
	assert.NoError(t, err)

	// Setup config v2 data
	f2, err := os.CreateTemp("", "tanzu_config_v2")
	assert.Nil(t, err)
	err = os.WriteFile(f2.Name(), []byte(cfgV2), 0644)
	assert.Nil(t, err)
	err = os.Setenv(EnvConfigV2Key, f2.Name())
	assert.NoError(t, err)

	//Setup metadata
	f3, err := os.CreateTemp("", "tanzu_config_metadata")
	assert.Nil(t, err)
	err = os.WriteFile(f3.Name(), []byte(""), 0644)
	assert.Nil(t, err)
	err = os.Setenv(EnvConfigMetadataKey, f3.Name())
	assert.NoError(t, err)

	// Cleanup
	defer func(name string) {
		err = os.Remove(name)
		assert.NoError(t, err)
	}(f1.Name())

	defer func(name string) {
		err = os.Remove(name)
		assert.NoError(t, err)
	}(f2.Name())

	defer func(name string) {
		err = os.Remove(name)
		assert.NoError(t, err)
	}(f3.Name())

	// Actions
	node, err := getClientConfig()
	assert.NotNil(t, node)
	assert.NoError(t, err)

	err = persistConfig(node)
	assert.NoError(t, err)
}
