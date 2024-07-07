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
package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/vistara-labs/firecracker-containerd/internal"
	"github.com/vistara-labs/firecracker-containerd/internal/integtest"
)

func init() {
	flag, err := internal.SupportCPUTemplate()
	if err != nil {
		panic(err)
	}

	if flag {
		integtest.DefaultRuntimeConfig.CPUTemplate = "T2"
	}
}

// devmapper is the only snapshotter we can use with Firecracker
const defaultSnapshotterName = "devmapper"

var testNameToVMIDReplacer = strings.NewReplacer("/", "-", "_", "-")

func testNameToVMID(s string) string {
	return testNameToVMIDReplacer.Replace(s)
}

func cgroupExists(name string) bool {
	// cgroups v1
	_, err := os.Stat(filepath.Join("/sys/fs/cgroup/cpu", name))
	if err == nil {
		return true
	}

	// cgroups v2
	_, err = os.Stat(filepath.Join("/sys/fs/cgroup", name))
	return err == nil
}
