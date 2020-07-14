/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package k8s

import (
	"net/http"
	"testing"

	"github.com/Tencent/bk-bcs/bcs-common/common/websocketDialer"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/tunnel"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
)

func testAuthorize(req *http.Request) (string, bool, error) {
	return "", true, nil
}

func TestGetTransport(t *testing.T) {
	clusterId := "k8s-001"
	wsCred := models.BcsWsClusterCredentials{
		ServerKey:     "k8s-001",
		ServerAddress: "https://127.0.0.1:443",
	}

	tunnel.DefaultTunnelServer = websocketDialer.New(testAuthorize, websocketDialer.DefaultErrorWriter, testCleanCredentials)
	tp := getTransport(clusterId, &wsCred)
	if tp != nil {
		t.Error("should have no tunnel session and return nil transport")
	}
}

func testCleanCredentials(serverKey string) {
	return
}