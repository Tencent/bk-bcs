package bkdata

import (
	"encoding/json"

	"github.com/Tencent/bk-bcs/bcs-common/pkg/esb"
)

const (
	BKDataAPIURL           = "http://bk-data.apigw.o.oa.com/prod/"
	NewAccessDeployPlanApi = "v3/access/deploy_plan/"
	NewCleanConfigApi      = "v3/databus/cleans/"
)

type BKDataClient struct {
	esb    *esb.EsbClient
	config BKDataClientConfig
}

func NewBKDataClient(conf BKDataClientConfig) (*BKDataClient, error) {
	esbclient, err := esb.NewEsbClient(conf.BkAppCode, conf.BkAppSecret, conf.BkUsername, BKDataAPIURL)
	if err != nil {
		return nil, err
	}
	return &BKDataClient{
		esb:    esbclient,
		config: conf,
	}, nil
}

// ObtainDataId obtain a new dataid from bk-data
// dataid != -1 : access deploy plan succ with new dataid returned
// error != nil : Some error occured while obtain dataid. The situation that
//                error != nil and  dataid != -1 is possible
func (c *BKDataClient) ObtainDataId(conf CustomAccessDeployPlanConfig) (int64, error) {
	conf.BkAppCode = c.config.BkAppCode
	conf.BkAppSecret = c.config.BkAppSecret
	conf.BkUsername = c.config.BkUsername
	jsonstr, err := json.Marshal(conf)
	if err != nil {
		return -1, err
	}
	var payload map[string]interface{}
	err = json.Unmarshal(jsonstr, &payload)
	if err != nil {
		return -1, err
	}
	data, err := c.esb.RequestEsb("POST", NewAccessDeployPlanApi, payload)
	if err != nil {
		return -1, err
	}
	var res map[string]interface{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return -1, err
	}
	var dataid int64
	dataid = res["raw_data_id"].(int64)
	return dataid, nil
}

func (c *BKDataClient) SetCleanStrategy(strategy DataCleanStrategy) error {
	strategy.BkAppCode = c.config.BkAppCode
	strategy.BkAppSecret = c.config.BkAppSecret
	strategy.BkUsername = c.config.BkUsername
	payload := map[string]interface{}{}
	jsonstr, err := json.Marshal(strategy)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonstr, &payload)
	if err != nil {
		return err
	}
	_, err = c.esb.RequestEsb("POST", NewCleanConfigApi, payload)
	if err != nil {
		return err
	}
	return nil
}