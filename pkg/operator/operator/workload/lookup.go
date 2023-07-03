// Tencent is pleased to support the open source community by making
// 蓝鲸智云 - 监控平台 (BlueKing - Monitor) available.
// Copyright (C) 2022 THL A29 Limited, a Tencent company. All rights reserved.
// Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://opensource.org/licenses/MIT
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package workload

func Lookup(oid ObjectID, podObj *Objects, objs map[string]*Objects) *OwnerRef {
	refs, ok := podObj.GetRefs(oid)
	if !ok {
		return nil
	}

	for _, ref := range refs {
		parent := lookup(oid.Namespace, ref, objs)
		if parent != nil {
			return parent
		}
	}

	return &OwnerRef{
		Kind: kindPod,
		Name: oid.Name,
	}
}

func lookup(namespace string, ref OwnerRef, objsMap map[string]*Objects) *OwnerRef {
	parent := &OwnerRef{}
	recursiveLookup(namespace, ref, objsMap, parent)
	if parent.Kind == "" && parent.Name == "" {
		return nil
	}
	return parent
}

func recursiveLookup(namespace string, ref OwnerRef, objsMap map[string]*Objects, parent *OwnerRef) {
	objs, ok := objsMap[ref.Kind]
	if !ok {
		return
	}

	found, ok := objs.Get(ObjectID{
		Name:      ref.Name,
		Namespace: namespace,
	})
	if !ok {
		return
	}
	parent.Kind = objs.Kind()
	parent.Name = ref.Name
	for _, childRef := range found.OwnerRefs {
		recursiveLookup(namespace, childRef, objsMap, parent)
	}
}
