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

package msgqueue

import (
	"context"
	"errors"
	"fmt"
	glog "github.com/Tencent/bk-bcs/bcs-common/common/blog"
	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-plugins/broker/rabbitmq/v2"
	"github.com/micro/go-plugins/broker/stan/v2"
	natstan "github.com/nats-io/stan.go"
)

// MessageQueue is an interface used for asynchronous messaging.
type MessageQueue interface {
	// publish pub data to queue
	// data.Header  map[string]string id: cluster_id; resourceType:Pod; namespace:default; resourceName: name; event: Update
	Publish(data *broker.Message) error

	// subscribe all data
	Subscribe(handler Handler, filters []Filter) (UnSub, error)

	// subscribe specified data type
	SubscribePod(handler Handler, filters []Filter) (UnSub, error)
	SubscribeEvent(handler Handler, filters []Filter) (UnSub, error)
	SubscribeDeployment(handler Handler, filters []Filter) (UnSub, error)
	SubscribeStatefulSet(handler Handler, filters []Filter) (UnSub, error)

	// String return queue name
	String() (string, error)
	// Stop the message queue
	Stop()
}

// UnSub for unSubscribe topic
type UnSub interface {
	Unsubscribe() error
}

// MsgQueue struct
type MsgQueue struct {
	queueOptions *QueueOptions
	broker       broker.Broker
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewMsgQueue init queue for rabbitmq/nats
func NewMsgQueue(flag bool, kind QueueKind, address string, resourceToQueue map[string]string,
	opts ...QueueOption) (MessageQueue, error) {
	var (
		queueOptions = &QueueOptions{
			QueueFlag:       flag,
			QueueKind:       kind,
			Address:         address,
			ResourceToQueue: resourceToQueue,
		}
		err error
	)

	for _, o := range opts {
		o(queueOptions)
	}

	messageQueue := &MsgQueue{
		queueOptions: queueOptions,
	}

	messageQueue.broker, err = NewQueueBroker(queueOptions)
	if err != nil {
		errMsg := fmt.Errorf("messageQueue init broker failed: %v", err)
		return nil, errMsg
	}

	messageQueue.ctx, messageQueue.cancel = context.WithCancel(context.Background())

	return messageQueue, nil
}

// Publish pub data into the topic
// data.headers map[string]string   id cluster_id; namespace default; resourceType Pod; resourceName name; event: update
// data.Body    data
func (mq *MsgQueue) Publish(data *broker.Message) error {
	if !mq.queueOptions.QueueFlag {
		return errors.New("queue flag is off")
	}

	resourceType, ok := data.Header["resourceType"]
	if !ok {
		return errors.New("message not exist resourceType, please input correct data")
	}

	queueName, err := mq.isExistResourceQueue(resourceType)
	if err != nil {
		errMsg := fmt.Errorf("resourceType to queue failed: %v", err)
		glog.Errorf("resourceType to queue failed: %v", err)
		return errMsg
	}

	switch mq.queueOptions.QueueKind {
	case RABBITMQ:
		err = mq.broker.Publish(queueName, data, rabbitmq.DeliveryMode(mq.queueOptions.PublishOptions.DeliveryMode))
	case NATSTREAMING:
		err = mq.broker.Publish(queueName, data)
	default:
		return errors.New("unsupported queue kind")
	}
	if err != nil {
		glog.Errorf("[pub] message failed: [messageType: %s], [messageQueue: %s], [cluster_id: %s], [namespace: %s], [resourceName: %s]",
			data.Header["resourceType"], queueName, data.Header["id"], data.Header["namespace"], data.Header["resourceName"])
	}

	glog.Infof("[pub] message successful: [messageType: %s], [messageQueue: %s], [cluster_id: %s], [namespace: %s], [resourceName: %s]",
		data.Header["resourceType"], queueName, data.Header["id"], data.Header["namespace"], data.Header["resourceName"])

	return nil
}

// Subscribe handle all cluster data
func (mq *MsgQueue) Subscribe(handler Handler, filters []Filter) (UnSub, error) {
	return nil, nil
}

// SubscribePod subscribe pod data with specific handler and filters
func (mq *MsgQueue) SubscribePod(handler Handler, filters []Filter) (UnSub, error) {
	if !mq.queueOptions.QueueFlag {
		return nil, errors.New("queue flag is off")
	}

	var (
		queueName  = "Pod"
		podHandler = &objectHandler{
			resourceType: "Pod",
			handler:      handler,
			filter:       filters,
		}
	)
	subscribeOptions, err := mq.getSubOptions(queueName)
	if err != nil {
		return nil, err
	}

	subscribe, err := mq.broker.Subscribe(queueName, podHandler.selfHandler, subscribeOptions...)
	if err != nil {
		glog.Errorf("subscribe failed: %v", err)
		return nil, err
	}

	glog.Infof("subscribe [%s:%s] successful", queueName, subscribe.Topic())

	return subscribe, nil
}

// SubscribeEvent subscribe event data with specific handler and filters
func (mq *MsgQueue) SubscribeEvent(handler Handler, filters []Filter) (UnSub, error) {
	if !mq.queueOptions.QueueFlag {
		return nil, errors.New("queue flag is off")
	}

	var (
		queueName    = "Event"
		eventHandler = &objectHandler{
			resourceType: "Event",
			handler:      handler,
			filter:       filters,
		}
	)
	subscribeOptions, err := mq.getSubOptions(queueName)
	if err != nil {
		return nil, err
	}

	subscribe, err := mq.broker.Subscribe(queueName, eventHandler.selfHandler, subscribeOptions...)
	if err != nil {
		glog.Errorf("subscribe failed: %v", err)
		return nil, err
	}

	glog.Infof("subscribe [%s:%s] successful", queueName, subscribe.Topic())

	return subscribe, nil
}

// SubscribeDeployment subscribe deployment data with specific handler and filters
func (mq *MsgQueue) SubscribeDeployment(handler Handler, filters []Filter) (UnSub, error) {
	if !mq.queueOptions.QueueFlag {
		return nil, errors.New("queue flag is off")
	}

	var (
		queueName         = "Deployment"
		deploymentHandler = &objectHandler{
			resourceType: "Deployment",
			handler:      handler,
			filter:       filters,
		}
	)
	subscribeOptions, err := mq.getSubOptions(queueName)
	if err != nil {
		return nil, err
	}

	subscribe, err := mq.broker.Subscribe(queueName, deploymentHandler.selfHandler, subscribeOptions...)
	if err != nil {
		glog.Errorf("subscribe failed: %v", err)
		return nil, err
	}

	glog.Infof("subscribe [%s:%s] successful", queueName, subscribe.Topic())

	return subscribe, nil
}

// SubscribeStatefulSet subscribe statefulSet data with specific handler and filters
func (mq *MsgQueue) SubscribeStatefulSet(handler Handler, filters []Filter) (UnSub, error) {
	if !mq.queueOptions.QueueFlag {
		return nil, errors.New("queue flag is off")
	}

	var (
		queueName          = "StatefulSet"
		statefulSetHandler = &objectHandler{
			resourceType: "Deployment",
			handler:      handler,
			filter:       filters,
		}
	)
	subscribeOptions, err := mq.getSubOptions(queueName)
	if err != nil {
		return nil, err
	}

	subscribe, err := mq.broker.Subscribe(queueName, statefulSetHandler.selfHandler, subscribeOptions...)
	if err != nil {
		glog.Errorf("subscribe failed: %v", err)
		return nil, err
	}

	glog.Infof("subscribe [%s:%s] successful", queueName, subscribe.Topic())

	return subscribe, nil
}

// Handlers of all topics
func (mq *MsgQueue) isExistResourceQueue(resourceType string) (string, error) {
	q, ok := mq.queueOptions.ResourceToQueue[resourceType]
	if !ok {
		return "", fmt.Errorf("resourceType[%s] not on subscribe", resourceType)
	}
	glog.Infof("resourceType %s sub queue[%s]", resourceType, q)
	return q, nil
}

// String declare queue kind
func (mq *MsgQueue) String() (string, error) {
	if !mq.queueOptions.QueueFlag {
		return "", errors.New("queue flag is off")
	}

	return string(mq.queueOptions.QueueKind), nil
}

// Stop the message queue
func (mq *MsgQueue) Stop() {
	mq.cancel()
	mq.broker.Disconnect()
}

func (mq *MsgQueue) getSubOptions(queueName string) ([]broker.SubscribeOption, error) {
	var subOptions []broker.SubscribeOption
	subOptions = append(subOptions, broker.Queue(queueName))
	if mq.queueOptions.SubscribeOptions.DisableAutoAck {
		subOptions = append(subOptions, broker.DisableAutoAck())
	}
	if mq.queueOptions.SubscribeOptions.AckOnSuccess {
		subOptions = append(subOptions, rabbitmq.AckOnSuccess())
	}

	switch mq.queueOptions.QueueKind {
	case RABBITMQ:
		if mq.queueOptions.SubscribeOptions.Durable {
			subOptions = append(subOptions, rabbitmq.DurableQueue())
		}
		if mq.queueOptions.SubscribeOptions.RequeueOnError {
			subOptions = append(subOptions, rabbitmq.RequeueOnError())
		}
		if len(mq.queueOptions.SubscribeOptions.QueueArguments) > 0 {
			subOptions = append(subOptions, rabbitmq.QueueArguments(mq.queueOptions.SubscribeOptions.QueueArguments))
		}
	case NATSTREAMING:
		var natsopts []natstan.SubscriptionOption
		if mq.queueOptions.SubscribeOptions.Durable {
			natsopts = append(natsopts, natstan.DurableName(queueName))
		}
		if mq.queueOptions.SubscribeOptions.DeliverAllMessage {
			natsopts = append(natsopts, natstan.DeliverAllAvailable())
		}
		if mq.queueOptions.SubscribeOptions.ManualAckMode {
			natsopts = append(natsopts, natstan.SetManualAckMode())
		}
		if mq.queueOptions.SubscribeOptions.EnableAckWait {
			natsopts = append(natsopts, natstan.AckWait(mq.queueOptions.SubscribeOptions.AckWaitDuration))
		}
		if mq.queueOptions.SubscribeOptions.MaxInFlight != 0 {
			natsopts = append(natsopts, natstan.MaxInflight(mq.queueOptions.SubscribeOptions.MaxInFlight))
		}
		subOptions = append(subOptions, stan.SubscribeOption(natsopts...))
	default:
		return nil, errors.New("unsupported queue kind")
	}

	return subOptions, nil
}
