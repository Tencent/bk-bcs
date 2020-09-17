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

package app

import (
	"context"
	"strings"

	"github.com/Tencent/bk-bcs/bcs-common/common/ssl"

	"github.com/Tencent/bk-bcs/bcs-common/common"
	"github.com/Tencent/bk-bcs/bcs-common/common/blog"
	commonconf "github.com/Tencent/bk-bcs/bcs-common/common/conf"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-log-manager/app/api"
	bkdata "github.com/Tencent/bk-bcs/bcs-services/bcs-log-manager/app/bkdataapi"
	bkdataCli "github.com/Tencent/bk-bcs/bcs-common/pkg/esb/apigateway/bkdata"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-log-manager/app/k8s"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-log-manager/app/options"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-log-manager/config"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-log-manager/pkg/util"
)

// Run run the bcs-log-manager module
func Run(ctx context.Context, stopCh chan struct{}, op *options.LogManagerOption) error {
	if err := common.SavePid(op.ProcessConfig); err != nil {
		blog.Error("fail to save pid: err:%s", err.Error())
	}

	// controller := bkdata.NewBKDataController(stopCh, op.KubeConfig, op.BKDataAPIHost)
	controller := &bkdata.BKDataController{
		StopCh:        stopCh,
		KubeConfig:    op.KubeConfig,
		ApiHost:       op.BKDataAPIHost,
		ClientCreator: bkdataCli.NewClientCreator(),
	}
	err := controller.Start()
	if err != nil {
		blog.Errorf("BKDataApi controller start failed: %s", err.Error())
		util.SendTermSignal()
	}
	blog.Info("BKDataApiConfig controller started")

	conf := &config.ManagerConfig{}
	setManagerConfig(op, conf)
	conf.StopCh = stopCh
	manager := k8s.NewManager(conf)
	manager.Start()
	blog.Info("Log Manager started")

	apiconf := &config.APIServerConfig{}
	setAPIServerConfig(op, apiconf)
	server := api.NewAPIServer(ctx, apiconf, manager)
	err = server.Run()
	if err != nil {
		blog.Errorf("APIServer start failed: %s", err.Error())
		util.SendTermSignal()
	}
	blog.Info("APIServer started")
	return nil
}

func setManagerConfig(op *options.LogManagerOption, conf *config.ManagerConfig) {
	conf.CollectionConfigs = op.CollectionConfigs
	conf.BcsApiConfig.Host = op.BcsAPIHost
	conf.BcsApiConfig.AuthToken = op.AuthToken
	conf.BcsApiConfig.Gateway = op.Gateway
	// TODO
	if op.BcsAPICAFile != "" {
		var err error
		conf.BcsApiConfig.TLSConfig, err = ssl.ClientTslConfVerity(op.BcsAPICAFile, op.BcsAPICertFile, op.BcsAPIKeyFile, "")
		if err != nil {
			blog.Errorf("ClientTslConfVerity of bcsapi failed: %s", err.Error())
			conf.BcsApiConfig.TLSConfig = nil
		}
	}
	conf.CAFile = op.ClientCertFile
	conf.SystemDataID = op.SystemDataID
	conf.BkAppCode = op.BkAppCode
	conf.BkUsername = op.BkUsername
	conf.BkAppSecret = op.BkAppSecret
	conf.BkBizID = op.BkBizID
	conf.KubeConfig = op.KubeConfig
}

func setAPIServerConfig(op *options.LogManagerOption, conf *config.APIServerConfig) {
	conf.Host = op.ServiceConfig.Address
	conf.Port = op.ServiceConfig.Port
	conf.BKDataAPIHost = op.BKDataAPIHost
	conf.EtcdHosts = strings.Split(op.EtcdHosts, ",")
	conf.EtcdCerts = commonconf.CertConfig{
		CAFile:         op.EtcdCAFile,
		ClientCertFile: op.EtcdCertFile,
		ClientKeyFile:  op.EtcdKeyFile,
	}
	conf.APICerts = commonconf.CertConfig{
		CAFile:         op.LogManagerCAFile,
		ServerCertFile: op.LogManagerCertFile,
		ServerKeyFile:  op.LogManagerKeyFile,
	}
	conf.ZkConfig = op.ZkConfig
}