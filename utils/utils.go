// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Hack for re-exporting
package utils

import (
	"context"
	"time"

	"github.com/vistara-labs/firecracker-containerd/internal/vm"
	"github.com/vistara-labs/firecracker-containerd/proto"
	"github.com/sirupsen/logrus"
)

const defaultVSockConnectTimeout = 5 * time.Second

type TaskManager = vm.TaskManager

func NewTaskManager(shimCtx context.Context, logger *logrus.Entry) TaskManager {
	return vm.NewTaskManager(shimCtx, logger)
}

func NewIOProxy(logger *logrus.Entry, stdin, stdout, stderr, VSockPath string, extraData *proto.ExtraData) (vm.IOProxy, error) {
	var ioConnectorSet vm.IOProxy

	if vm.IsAgentOnlyIO(stdout, logger) {
		ioConnectorSet = vm.NewNullIOProxy()
	} else {
		var stdinConnectorPair *vm.IOConnectorPair
		if stdin != "" {
			stdinConnectorPair = &vm.IOConnectorPair{
				ReadConnector:  vm.ReadFIFOConnector(stdin),
				WriteConnector: vm.VSockDialConnector(defaultVSockConnectTimeout, VSockPath, extraData.StdinPort),
			}
		}

		var stdoutConnectorPair *vm.IOConnectorPair
		if stdout != "" {
			stdoutConnectorPair = &vm.IOConnectorPair{
				ReadConnector:  vm.VSockDialConnector(defaultVSockConnectTimeout, VSockPath, extraData.StdoutPort),
				WriteConnector: vm.WriteFIFOConnector(stdout),
			}
		}

		var stderrConnectorPair *vm.IOConnectorPair
		if stderr != "" {
			stderrConnectorPair = &vm.IOConnectorPair{
				ReadConnector:  vm.VSockDialConnector(defaultVSockConnectTimeout, VSockPath, extraData.StderrPort),
				WriteConnector: vm.WriteFIFOConnector(stderr),
			}
		}

		ioConnectorSet = vm.NewIOConnectorProxy(stdinConnectorPair, stdoutConnectorPair, stderrConnectorPair)
	}
	return ioConnectorSet, nil
}
