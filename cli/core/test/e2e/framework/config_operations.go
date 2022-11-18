// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package framework

// ConfOps performs "tanzu config" command operations
type ConfOps interface {
	ConfigSetFeature(path, value string) error
	ConfigFeatureFlagExists(flagName, value string) (bool, error)
	ConfigFeatureFlagNotExists(flagName string) (bool, error)
	ConfigUnsetFeature(path string) error
	ConfigInit() error
	ConfigServerList() error
	ConfigServerDelete(serverName string) error
}

// confOps is the implementation of ConfOps interface
type confOps struct {
	CmdOps
}

func NewConfOps() ConfOps {
	return &confOps{
		CmdOps: NewCmdOps(),
	}
}

// ConfigSetFeature sets the tanzu config feature flag
func (co *confOps) ConfigSetFeature(path, value string) (err error) {
	confSetComd := ConfigSet + path + " " + value
	_, _, err = co.Exec(confSetComd)
	return err
}

// ConfigFeatureFlagExists validates the existence of feature flag in tanzu config
func (co *confOps) ConfigFeatureFlagExists(flagName, value string) (bool, error) {
	getFeatureCmd := ConfigGet
	featureFlagStr := flagName + ": " + "\"" + value + "\""
	err := co.ExecContainsString(getFeatureCmd, featureFlagStr)
	if err != nil {
		return false, err
	}
	return true, err
}

// ConfigFeatureFlagNotExists validates the existence of feature flag in tanzu config
func (co *confOps) ConfigFeatureFlagNotExists(flagName string) (bool, error) {
	getConfCmd := ConfigGet
	err := co.ExecNotContainsString(getConfCmd, flagName+": ")
	if err != nil {
		return false, err
	}
	return true, err
}

// ConfigUnsetFeature un-sets the tanzu config feature flag
func (co *confOps) ConfigUnsetFeature(path string) error {
	unsetFeatureCmd := ConfigUnset + path
	_, _, err := co.Exec(unsetFeatureCmd)
	return err
}

// ConfigInit performs "tanzu config init"
func (co *confOps) ConfigInit() error {
	_, _, err := co.Exec(ConfigInit)
	return err
}

// ConfigServerList returns the server list
// TODO: should return the servers info in proper format
func (co *confOps) ConfigServerList() (err error) {
	_, _, err = co.Exec(ConfigServerList)
	return nil
}

// ConfigServerDelete deletes a server from tanzu config
func (co *confOps) ConfigServerDelete(serverName string) error {
	_, _, err := co.Exec(ConfigServerDelete + serverName)
	return err
}
