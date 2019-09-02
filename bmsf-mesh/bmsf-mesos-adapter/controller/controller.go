/*
Copyright (C) 2019 The BlueKing Authors. All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package controller

import (
	"bk-bcs/bmsf-mesh/bmsf-mesos-adapter/controller/appnode"
	"bk-bcs/bmsf-mesh/bmsf-mesos-adapter/controller/appsvc"
	"bk-bcs/bmsf-mesh/bmsf-mesos-adapter/pkg/queue"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager, svcQ queue.Queue, nodeQ queue.Queue) error {
	if err := appsvc.Add(m, svcQ); err != nil {
		return err
	}
	if err := appnode.Add(m, nodeQ); err != nil {
		return err
	}
	return nil
}