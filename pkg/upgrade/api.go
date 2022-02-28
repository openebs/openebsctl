/*
Copyright 2020-2022 The OpenEBS Authors

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

package upgrade

import (
	"fmt"

	corebuilder "github.com/openebs/api/v2/pkg/kubernetes/core"
	batchV1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

type Job struct {
	*batchV1.Job
}

// NewJob returns an empty instance of BatchJob
func NewJob() *Job {
	return &Job{
		&batchV1.Job{},
	}
}

// WithName sets the name of the field of Job
func (b *Job) WithName(name string) *Job {
	b.Name = name
	return b
}

// WithGeneratedName Creates a job with auto-generated name
func (b *Job) WithGeneratedName(name string) *Job {
	b.GenerateName = fmt.Sprintf("%s-", name)
	return b
}

// WithLabel sets label for the job
func (b *Job) WithLabel(label map[string]string) *Job {
	b.Labels = label
	return b
}

// WithNamespace sets the namespace of the Job
func (b *Job) WithNamespace(namespace string) *Job {
	b.Namespace = namespace
	return b
}

// BuildJobSpec builds an empty Job Spec
func (b *Job) BuildJobSpec() *Job {
	b.Spec = batchV1.JobSpec{}
	return b
}

// WithBackOffLimit sets the backOffLimit for pods in the Job with given value
func (b *Job) WithBackOffLimit(limit int32) *Job {
	b.Spec.BackoffLimit = &limit
	return b
}

// WithPodTemplateSpec sets the template Field for Job
func (b *Job) WithPodTemplateSpec(pts *corebuilder.PodTemplateSpec) *Job {
	templateSpecObj := pts.Build()
	b.Spec.Template = *templateSpecObj
	return b
}

// Temporary code until PR into openebs/api is not merged----
func (b *Job) WithRestartPolicy(policy corev1.RestartPolicy) *Job {
	b.Spec.Template.Spec.RestartPolicy = policy
	return b
}
