// Copyright 2018-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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

package vm

import (
	"context"
	"syscall"

	"github.com/containerd/fifo"
	"github.com/sirupsen/logrus"
)

// FIFOConnector adapts containerd's fifo package to the IOConnector interface
func FIFOConnector(path string) IOConnector {
	return func(procCtx context.Context, logger *logrus.Entry) <-chan IOConnectorResult {
		returnCh := make(chan IOConnectorResult, 1)
		defer close(returnCh)

		// We open the FIFO synchronously to ensure that the FIFO is created (via O_CREAT) before
		// it is passed to any task service. O_RDWR ensures that we don't block on the syscall
		// level (as documented in the fifo pkg).
		fifo, err := fifo.OpenFifo(procCtx, path, syscall.O_CREAT|syscall.O_RDWR, 0300)
		returnCh <- IOConnectorResult{
			ReadWriteCloser: fifo,
			Err:             err,
		}

		return returnCh
	}
}
