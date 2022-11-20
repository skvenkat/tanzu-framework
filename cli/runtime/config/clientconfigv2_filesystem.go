// Copyright 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"
)

const (
	// EnvConfigV2Key is the environment variable that points to a tanzu config.
	EnvConfigV2Key = "TANZU_CONFIG_V2"

	// CfgV2Name is the name of the config metadata
	CfgV2Name = "config-v2.yaml"
)

// clientConfigV2Path constructs the full config path, checking for environment overrides.
func clientConfigV2Path(localDirGetter func() (string, error)) (path string, err error) {
	localDir, err := localDirGetter()
	if err != nil {
		return path, err
	}
	var ok bool
	path, ok = os.LookupEnv(EnvConfigV2Key)
	if !ok {
		path = filepath.Join(localDir, CfgV2Name)
		return
	}
	return
}

// ClientConfigV2Path retrieved config-alt file path
func ClientConfigV2Path() (path string, err error) {
	return clientConfigV2Path(LocalDir)
}
