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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	bkbcstencentcomv1 "github.com/Tencent/bk-bcs/bcs-services/bcs-log-manager/pkg/apis/bkbcs.tencent.com/v1"
)

// BKDataApiConfigReconciler reconciles a BKDataApiConfig object
type BKDataApiConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=bkbcs.tencent.com.bkbcs.tencent.com,resources=bkdataapiconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=bkbcs.tencent.com.bkbcs.tencent.com,resources=bkdataapiconfigs/status,verbs=get;update;patch

func (r *BKDataApiConfigReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("bkdataapiconfig", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *BKDataApiConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bkbcstencentcomv1.BKDataApiConfig{}).
		Complete(r)
}