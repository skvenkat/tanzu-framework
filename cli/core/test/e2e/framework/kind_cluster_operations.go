// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package framework

import (
	"fmt"
)

// KindClusterOps performs KIND cluster operations
type KindClusterOps interface {
	CreateKindCluster(name string) (string, error)
	DockerStatus() (info string, err error)
}

type kindClusterOps struct {
	CmdOps
}

func NewKindClusterOps() KindClusterOps {
	return &kindClusterOps{
		CmdOps: NewCmdOps(),
	}
}

// CreateKindCluster create kind cluster with given name and returns stdout info
// if docker not running or any error then returns stdout and error info
func (kc *kindClusterOps) CreateKindCluster(name string) (string, error) {
	stdOut, err := kc.DockerStatus()
	if err != nil {
		return stdOut, err
	}
	stdOutBuffer, stdErrBuffer, err := kc.Exec(KindCreateCluster + name)
	if err != nil {
		return stdOutBuffer.String(), fmt.Errorf(stdErrBuffer.String(), err)
	}
	return stdOutBuffer.String(), err
}

// DockerStatus runs 'docker info' command and returns the stdout or error if any
func (kc *kindClusterOps) DockerStatus() (string, error) {
	stdOut, stdErr, err := kc.Exec(DockerInfo)
	if err != nil {
		return stdOut.String(), fmt.Errorf(stdErr.String(), err)
	}
	return stdOut.String(), err
}
