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

package v1http

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Tencent/bk-bcs/bcs-common/common"
	"github.com/Tencent/bk-bcs/bcs-common/common/blog"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/metrics"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/models"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/storages/sqlstore"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-user-manager/app/user-manager/utils"

	"github.com/dchest/uniuri"
	"github.com/emicklei/go-restful"
)

//DefaultTokenLength user token default length
// token is consisted of digital and alphabet(case sensetive)
// we can refer to http://coolaf.com/tool/rd when testing
const DefaultTokenLength = 32

// CreateAdminUser create a admin user
func CreateAdminUser(request *restful.Request, response *restful.Response) {
	start := time.Now()

	userName := request.PathParameter("user_name")
	user := &models.BcsUser{
		Name:     userName,
		UserType: sqlstore.AdminUser,
	}
	// if this user already exist
	userInDb := sqlstore.GetUserByCondition(user)
	if userInDb != nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		message := fmt.Sprintf("errcode: %d, user [%s] already exist", common.BcsErrApiBadRequest, userName)
		utils.WriteClientError(response, common.BcsErrApiBadRequest, message)
		return
	}
	user.UserToken = uniuri.NewLen(DefaultTokenLength)
	user.ExpiresAt = time.Now().Add(sqlstore.AdminSaasUserExpiredTime)

	// create this user and save to db
	err := sqlstore.CreateUser(user)
	if err != nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Errorf("failed to create user [%s]: %s", user.Name, err.Error())
		message := fmt.Sprintf("errcode: %d, creating user [%s] failed, error: %s", common.BcsErrApiInternalDbError, user.Name, err)
		utils.WriteServerError(response, common.BcsErrApiInternalDbError, message)
		return
	}

	data := utils.CreateResponeData(nil, "success", *user)
	_, _ = response.Write([]byte(data))

	metrics.RequestCount.WithLabelValues("user", request.Request.Method).Inc()
	metrics.RequestLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
}

// GetAdminUser get an admin user and usertoken information
func GetAdminUser(request *restful.Request, response *restful.Response) {
	start := time.Now()

	userName := request.PathParameter("user_name")
	user := sqlstore.GetUserByCondition(&models.BcsUser{Name: userName, UserType: sqlstore.AdminUser})
	if user == nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Warnf("user [%s] not found in db", userName)
		message := fmt.Sprintf("errcode: %d, user with user_name=%s not found", common.BcsErrApiBadRequest, userName)
		utils.WriteNotFoundError(response, common.BcsErrApiBadRequest, message)
		return
	}

	data := utils.CreateResponeData(nil, "success", *user)
	_, _ = response.Write([]byte(data))

	metrics.RequestCount.WithLabelValues("user", request.Request.Method).Inc()
	metrics.RequestLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
}

// CreateSaasUser create a saas user
func CreateSaasUser(request *restful.Request, response *restful.Response) {
	start := time.Now()

	userName := request.PathParameter("user_name")
	user := &models.BcsUser{
		Name:     userName,
		UserType: sqlstore.SaasUser,
	}
	// if this user already exist
	userInDb := sqlstore.GetUserByCondition(user)
	if userInDb != nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		message := fmt.Sprintf("errcode: %d, user [%s] already exist", common.BcsErrApiBadRequest, userName)
		utils.WriteClientError(response, common.BcsErrApiBadRequest, message)
		return
	}

	user.UserToken = uniuri.NewLen(DefaultTokenLength)
	user.ExpiresAt = time.Now().Add(sqlstore.AdminSaasUserExpiredTime)

	// create this user and save to db
	err := sqlstore.CreateUser(user)
	if err != nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Errorf("failed to create user [%s]: %s", user.Name, err.Error())
		message := fmt.Sprintf("errcode: %d, creating user [%s] failed, error: %s", common.BcsErrApiInternalDbError, user.Name, err)
		utils.WriteServerError(response, common.BcsErrApiInternalDbError, message)
		return
	}

	data := utils.CreateResponeData(nil, "success", *user)
	_, _ = response.Write([]byte(data))

	metrics.RequestCount.WithLabelValues("user", request.Request.Method).Inc()
	metrics.RequestLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
}

// GetSaasUser get an saas user and usertoken information
func GetSaasUser(request *restful.Request, response *restful.Response) {
	start := time.Now()

	userName := request.PathParameter("user_name")
	user := sqlstore.GetUserByCondition(&models.BcsUser{Name: userName, UserType: sqlstore.SaasUser})
	if user == nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Warnf("user [%s] not found in db", userName)
		message := fmt.Sprintf("errcode: %d, user with user_name=%s not found", common.BcsErrApiBadRequest, userName)
		utils.WriteNotFoundError(response, common.BcsErrApiBadRequest, message)
		return
	}

	data := utils.CreateResponeData(nil, "success", *user)
	_, _ = response.Write([]byte(data))

	metrics.RequestCount.WithLabelValues("user", request.Request.Method).Inc()
	metrics.RequestLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
}

// CreatePlainUser create a plain user
func CreatePlainUser(request *restful.Request, response *restful.Response) {
	start := time.Now()

	userName := request.PathParameter("user_name")
	user := &models.BcsUser{
		Name:     userName,
		UserType: sqlstore.PlainUser,
	}
	// if this user already exist
	userInDb := sqlstore.GetUserByCondition(user)
	if userInDb != nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		message := fmt.Sprintf("errcode: %d, user [%s] already exist", common.BcsErrApiBadRequest, userName)
		utils.WriteClientError(response, common.BcsErrApiBadRequest, message)
		return
	}

	user.UserToken = uniuri.NewLen(DefaultTokenLength)
	user.ExpiresAt = time.Now().Add(sqlstore.PlainUserExpiredTime)

	// create this user and save to db
	err := sqlstore.CreateUser(user)
	if err != nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Errorf("failed to create user [%s]: %s", user.Name, err.Error())
		message := fmt.Sprintf("errcode: %d, creating user [%s] failed, error: %s", common.BcsErrApiInternalDbError, user.Name, err)
		utils.WriteServerError(response, common.BcsErrApiInternalDbError, message)
		return
	}

	data := utils.CreateResponeData(nil, "success", *user)
	_, _ = response.Write([]byte(data))

	metrics.RequestCount.WithLabelValues("user", request.Request.Method).Inc()
	metrics.RequestLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
}

// GetPlainUser get an plain user and usertoken information
func GetPlainUser(request *restful.Request, response *restful.Response) {
	start := time.Now()

	userName := request.PathParameter("user_name")
	user := sqlstore.GetUserByCondition(&models.BcsUser{Name: userName, UserType: sqlstore.PlainUser})
	if user == nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Warnf("failed to get user, user [%s] not found in db", userName)
		message := fmt.Sprintf("errcode: %d, user with user_name=%s not found", common.BcsErrApiBadRequest, userName)
		utils.WriteNotFoundError(response, common.BcsErrApiBadRequest, message)
		return
	}

	data := utils.CreateResponeData(nil, "success", *user)
	_, _ = response.Write([]byte(data))

	metrics.RequestCount.WithLabelValues("user", request.Request.Method).Inc()
	metrics.RequestLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
}

// RefreshPlainToken refresh usertoken for a plain user
func RefreshPlainToken(request *restful.Request, response *restful.Response) {
	start := time.Now()

	userName := request.PathParameter("user_name")
	expireDays := request.PathParameter("expire_time")
	expireDaysInt, err := strconv.Atoi(expireDays)
	if err != nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Warnf("invalid expire_time, failed to atoi: %s", err.Error())
		message := fmt.Sprintf("errcode: %d, invalid expire_time, failed to atoi: %s", common.BcsErrApiBadRequest, err.Error())
		utils.WriteClientError(response, common.BcsErrApiBadRequest, message)
		return
	}
	if expireDaysInt < 0 {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Warnf("invalid expire_time: %d", expireDaysInt)
		message := fmt.Sprintf("errcode: %d, invalid expire_time", common.BcsErrApiBadRequest)
		utils.WriteClientError(response, common.BcsErrApiBadRequest, message)
		return
	}

	user := sqlstore.GetUserByCondition(&models.BcsUser{Name: userName, UserType: sqlstore.PlainUser})
	if user == nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Warnf("failed to refresh token, user [%s] not found in db", userName)
		message := fmt.Sprintf("errcode: %d, user with user_name=%s not found", common.BcsErrApiBadRequest, userName)
		utils.WriteNotFoundError(response, common.BcsErrApiBadRequest, message)
		return
	}

	expireTime := time.Duration(expireDaysInt) * sqlstore.PlainUserExpiredTime
	updatedUser := user
	// if usertoken has been expired, refresh the usertoken
	// or just refresh the expiresTime and return the same token
	if time.Now().After(user.ExpiresAt) {
		updatedUser.UserToken = uniuri.NewLen(DefaultTokenLength)
		updatedUser.ExpiresAt = time.Now().Add(expireTime)
	} else {
		updatedUser.ExpiresAt = time.Now().Add(expireTime)
	}

	// update and save to db
	// if update failed, it's better to refresh by client
	err = sqlstore.UpdateUser(user, updatedUser)
	if err != nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Errorf("failed to refresh usertoken [%s]: %s", user.Name, err.Error())
		message := fmt.Sprintf("errcode: %d, failed to refresh usertoken [%s], error: %s", common.BcsErrApiInternalDbError, userName, err)
		utils.WriteServerError(response, common.BcsErrApiInternalDbError, message)
		return
	}

	data := utils.CreateResponeData(nil, "success", *user)
	_, _ = response.Write([]byte(data))

	metrics.RequestCount.WithLabelValues("user", request.Request.Method).Inc()
	metrics.RequestLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
}

// RefreshSaasToken refresh usertoken for a saas user
func RefreshSaasToken(request *restful.Request, response *restful.Response) {
	start := time.Now()

	userName := request.PathParameter("user_name")
	user := sqlstore.GetUserByCondition(&models.BcsUser{Name: userName, UserType: sqlstore.SaasUser})
	if user == nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Warnf("failed to refresh token, user [%s] not found in db", userName)
		message := fmt.Sprintf("errcode: %d, user with user_name=%s not found", common.BcsErrApiBadRequest, userName)
		utils.WriteNotFoundError(response, common.BcsErrApiBadRequest, message)
		return
	}

	// refresh the usertoken
	updatedUser := user
	updatedUser.UserToken = uniuri.NewLen(DefaultTokenLength)
	updatedUser.ExpiresAt = time.Now().Add(sqlstore.AdminSaasUserExpiredTime)

	// update and save to db
	// if update failed, it's better to refresh by client
	err := sqlstore.UpdateUser(user, updatedUser)
	if err != nil {
		metrics.RequestErrorCount.WithLabelValues("user", request.Request.Method).Inc()
		metrics.RequestErrorLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
		blog.Errorf("failed to refresh usertoken [%s]: %s", user.Name, err.Error())
		message := fmt.Sprintf("errcode: %d, failed to refresh usertoken [%s], error: %s", common.BcsErrApiInternalDbError, userName, err)
		utils.WriteServerError(response, common.BcsErrApiInternalDbError, message)
		return
	}

	data := utils.CreateResponeData(nil, "success", *user)
	_, _ = response.Write([]byte(data))

	metrics.RequestCount.WithLabelValues("user", request.Request.Method).Inc()
	metrics.RequestLatency.WithLabelValues("user", request.Request.Method).Observe(time.Since(start).Seconds())
}
