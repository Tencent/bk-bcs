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

package types

import "fmt"

type CpusetNode struct {
	//cpuset node id
	//show available: 2 nodes (0-1)
	Id string
	//the node include cpusets
	//show node 0 cpus: 0 1 2 3 4 5 6 7 8 9 10 11 24 25 26 27 28 29 30 31 32 33 34 35
	Cpuset []string
	//allocated cpuset of container
	//the cpuset belongs to only one container
	AllocatedCpuset []string
}

func (c *CpusetNode) Capacity() int {
	return len(c.Cpuset) - len(c.AllocatedCpuset)
}

func (c *CpusetNode) AllocateCpuset(number int) ([]string, error) {
	if c.Capacity() < number {
		return nil, fmt.Errorf("Cpuset node %s Capacity %d, then can't allocate %d cpuset", c.Id, c.Capacity(), number)
	}

	cpuset := make([]string, 0, number)
	for _, o := range c.Cpuset {
		allocated := false
		for _, set := range c.AllocatedCpuset {
			if o == set {
				allocated = true
				break
			}
		}
		if allocated {
			continue
		}
		cpuset = append(cpuset, o)
		c.AllocatedCpuset = append(c.AllocatedCpuset, o)
		if len(cpuset) == number {
			break
		}
	}

	return cpuset, nil
}
