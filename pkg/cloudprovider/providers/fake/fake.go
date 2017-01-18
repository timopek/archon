/*
Copyright 2016 The Archon Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	"github.com/golang/glog"
	"io"

	archoncloudprovider "kubeup.com/archon/pkg/cloudprovider"
	"kubeup.com/archon/pkg/cluster"
	"kubeup.com/archon/pkg/userdata"

	"k8s.io/kubernetes/pkg/cloudprovider/providers/fake"
)

const ProviderName = "fake"

type FakeInstance struct {
	name     string
	spec     cluster.InstanceSpec
	userdata []byte
}

type FakeCloud struct {
	fake.FakeCloud
	FakeInstances map[string]FakeInstance
}

var (
	_ archoncloudprovider.Interface = &FakeCloud{}
)

func (f *FakeCloud) addCall(desc string) {
	f.FakeCloud.Calls = append(f.FakeCloud.Calls, desc)
}

func (f *FakeCloud) Archon() (archoncloudprovider.ArchonInterface, bool) {
	return f, true
}

func (f *FakeCloud) PrivateIP() (archoncloudprovider.PrivateIPInterface, bool) {
	return nil, false
}

func (f *FakeCloud) PublicIP() (archoncloudprovider.PublicIPInterface, bool) {
	return nil, false
}

func (f *FakeCloud) AddNetworkAnnotation(clustername string, instance *cluster.Instance, network *cluster.Network) error {
	return nil
}

func (f *FakeCloud) EnsureNetwork(clusterName string, network *cluster.Network) (status *cluster.NetworkStatus, err error) {
	return &cluster.NetworkStatus{Phase: cluster.NetworkRunning}, nil
}

func (f *FakeCloud) EnsureNetworkDeleted(clusterName string, network *cluster.Network) (err error) {
	return nil
}

func (f *FakeCloud) GetInstance(clusterName string, instance *cluster.Instance) (*cluster.InstanceStatus, error) {
	status := &cluster.InstanceStatus{}
	status.Phase = cluster.InstanceRunning

	return status, f.Err
}

func (f *FakeCloud) EnsureInstanceDependency(clusterName string, instance *cluster.Instance) (*cluster.InstanceStatus, error) {
	f.addCall("createDependency")
	return &instance.Status, nil
}

func (f *FakeCloud) EnsureInstance(clusterName string, instance *cluster.Instance) (*cluster.InstanceStatus, error) {
	f.addCall("create")
	if f.FakeInstances == nil {
		f.FakeInstances = make(map[string]FakeInstance)
	}

	name := instance.Name
	spec := instance.Spec
	userdata, err := userdata.Generate(instance)
	if err != nil {
		return nil, err
	}

	glog.V(2).Infof("Create instance with userdata: %s", string(userdata))

	f.FakeInstances[name] = FakeInstance{name, spec, userdata}

	status := &cluster.InstanceStatus{}
	status.Phase = cluster.InstanceRunning

	return status, f.Err
}

func (f *FakeCloud) EnsureInstanceDeleted(clusterName string, instance *cluster.Instance) error {
	f.addCall("delete")
	return f.Err
}

func (f *FakeCloud) EnsureInstanceDependencyDeleted(clusterName string, instance *cluster.Instance) (err error) {
	f.addCall("deleteDependency")
	return nil
}

func (f *FakeCloud) ListInstances(clusterName string, network *cluster.Network, selector map[string]string) ([]string, []*cluster.InstanceStatus, error) {
	f.addCall("list")
	result := make([]string, 0)
	instances := make([]*cluster.InstanceStatus, 0)
	return result, instances, f.Err
}

func init() {
	archoncloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (archoncloudprovider.Interface, error) {
		return &FakeCloud{}, nil
	})
}
