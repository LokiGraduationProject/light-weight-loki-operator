/*
Copyright 2024.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type LokiStackSpec struct {
	// Size defines one of the support Loki deployment scale out sizes.
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:select:1x.extra-small","urn:alm:descriptor:com.tectonic.ui:select:1x.small","urn:alm:descriptor:com.tectonic.ui:select:1x.medium"},displayName="LokiStack Size"
	Size LokiStackSizeType `json:"size"`

	// Storage defines the spec for the object storage endpoint to store logs.
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Object Storage"
	Storage ObjectStorageSpec `json:"storage"`

	// Storage class name defines the storage class for ingester/querier PVCs.
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:io.kubernetes:StorageClass",displayName="Storage Class Name"
	StorageClassName string `json:"storageClassName"`

	// Template defines the resource per component.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:advanced",displayName="Node Placement"
	Template *LokiTemplateSpec `json:"template,omitempty"`

	// Limits defines the limits to be applied to log stream processing.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:advanced",displayName="Rate Limiting"
	Limits *LimitsSpec `json:"limits,omitempty"`
}

// LokiStackSizeType declares the type for loki cluster scale outs.
//
// +kubebuilder:validation:Enum="1x.demo";"1x.extra-small";"1x.small";"1x.medium"
type LokiStackSizeType string

const (
	// SizeOneXDemo defines the size of a single Loki deployment
	// with tiny resource requirements and without HA support.
	// This size is intended to run in single-node clusters on laptops,
	// it is only useful for very light testing, demonstrations, or prototypes.
	// There are no ingestion/query performance guarantees.
	// DO NOT USE THIS IN PRODUCTION!
	SizeOneXDemo LokiStackSizeType = "1x.demo"

	// SizeOneXExtraSmall defines the size of a single Loki deployment
	// with extra small resources/limits requirements and without HA support.
	// This size is ultimately dedicated for development and demo purposes.
	// DO NOT USE THIS IN PRODUCTION!
	//
	// FIXME: Add clear description of ingestion/query performance expectations.
	SizeOneXExtraSmall LokiStackSizeType = "1x.extra-small"

	// SizeOneXSmall defines the size of a single Loki deployment
	// with small resources/limits requirements and HA support for all
	// Loki components. This size is dedicated for setup **without** the
	// requirement for single replication factor and auto-compaction.
	//
	// FIXME: Add clear description of ingestion/query performance expectations.
	SizeOneXSmall LokiStackSizeType = "1x.small"

	// SizeOneXMedium defines the size of a single Loki deployment
	// with small resources/limits requirements and HA support for all
	// Loki components. This size is dedicated for setup **with** the
	// requirement for single replication factor and auto-compaction.
	//
	// FIXME: Add clear description of ingestion/query performance expectations.
	SizeOneXMedium LokiStackSizeType = "1x.medium"
)

// ObjectStorageSpec defines the requirements to access the object
// storage bucket to persist logs by the ingester component.
type ObjectStorageSpec struct {
	// Schemas for reading and writing logs.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinItems:=1
	// +kubebuilder:default:={{version:v11,effectiveDate:"2020-10-11"}}
	Schemas []ObjectStorageSchema `json:"schemas"`

	// Secret for object storage authentication.
	// Name of a secret in the same namespace as the LokiStack custom resource.
	//
	// +required
	// +kubebuilder:validation:Required
	Secret ObjectStorageSecretSpec `json:"secret"`
}

// ObjectStorageSchema defines a schema version and the date when it will become effective.
type ObjectStorageSchema struct {
	// Version for writing and reading logs.
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:select:v11","urn:alm:descriptor:com.tectonic.ui:select:v12","urn:alm:descriptor:com.tectonic.ui:select:v13"},displayName="Version"
	Version ObjectStorageSchemaVersion `json:"version"`

	// EffectiveDate contains a date in YYYY-MM-DD format which is interpreted in the UTC time zone.
	//
	// The configuration always needs at least one schema that is currently valid. This means that when creating a new
	// LokiStack it is recommended to add a schema with the latest available version and an effective date of "yesterday".
	// New schema versions added to the configuration always needs to be placed "in the future", so that Loki can start
	// using it once the day rolls over.
	//
	// +required
	// +kubebuilder:validation:Required
	EffectiveDate StorageSchemaEffectiveDate `json:"effectiveDate"`
}

// ObjectStorageSchemaVersion defines the storage schema version which will be
// used with the Loki cluster.
//
// +kubebuilder:validation:Enum=v11;v12;v13
type ObjectStorageSchemaVersion string

const (
	// ObjectStorageSchemaV11 when using v11 for the storage schema
	ObjectStorageSchemaV11 ObjectStorageSchemaVersion = "v11"

	// ObjectStorageSchemaV12 when using v12 for the storage schema
	ObjectStorageSchemaV12 ObjectStorageSchemaVersion = "v12"

	// ObjectStorageSchemaV13 when using v13 for the storage schema
	ObjectStorageSchemaV13 ObjectStorageSchemaVersion = "v13"
)

// ObjectStorageSecretSpec is a secret reference containing name only, no namespace.
type ObjectStorageSecretSpec struct {
	// Type of object storage that should be used
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:select:azure","urn:alm:descriptor:com.tectonic.ui:select:gcs","urn:alm:descriptor:com.tectonic.ui:select:s3","urn:alm:descriptor:com.tectonic.ui:select:swift","urn:alm:descriptor:com.tectonic.ui:select:alibabacloud"},displayName="Object Storage Secret Type"
	Type ObjectStorageSecretType `json:"type"`

	// Name of a secret in the namespace configured for object storage secrets.
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:io.kubernetes:Secret",displayName="Object Storage Secret Name"
	Name string `json:"name"`
}

// ObjectStorageSecretType defines the type of storage which can be used with the Loki cluster.
//
// +kubebuilder:validation:Enum=azure;gcs;s3;swift;alibabacloud;
type ObjectStorageSecretType string

const (
	// ObjectStorageSecretS3 when using S3 for Loki storage
	ObjectStorageSecretS3 ObjectStorageSecretType = "s3"
)

// LokiStackStatus defines the observed state of LokiStack
type LokiStackStatus struct {
	// Components provides summary of all Loki pod status grouped
	// per component.
	//
	// +optional
	// +kubebuilder:validation:Optional
	Components LokiStackComponentStatus `json:"components,omitempty"`

	// Storage provides summary of all changes that have occurred
	// to the storage configuration.
	//
	// +optional
	// +kubebuilder:validation:Optional
	Storage LokiStackStorageStatus `json:"storage,omitempty"`

	// Conditions of the Loki deployment health.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:io.kubernetes.conditions"
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// LokiStackComponentStatus defines the map of per pod status per LokiStack component.
// Each component is represented by a separate map of v1.Phase to a list of pods.
type LokiStackComponentStatus struct {
	// Compactor is a map to the pod status of the compactor pod.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses",displayName="Compactor",order=5
	Compactor PodStatusMap `json:"compactor,omitempty"`

	// Distributor is a map to the per pod status of the distributor deployment
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses",displayName="Distributor",order=1
	Distributor PodStatusMap `json:"distributor,omitempty"`

	// IndexGateway is a map to the per pod status of the index gateway statefulset
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses",displayName="IndexGateway",order=6
	IndexGateway PodStatusMap `json:"indexGateway,omitempty"`

	// Ingester is a map to the per pod status of the ingester statefulset
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses",displayName="Ingester",order=2
	Ingester PodStatusMap `json:"ingester,omitempty"`

	// Querier is a map to the per pod status of the querier deployment
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses",displayName="Querier",order=3
	Querier PodStatusMap `json:"querier,omitempty"`

	// QueryFrontend is a map to the per pod status of the query frontend deployment
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses",displayName="Query Frontend",order=4
	QueryFrontend PodStatusMap `json:"queryFrontend,omitempty"`

	// Gateway is a map to the per pod status of the lokistack gateway deployment.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses",displayName="Gateway",order=5
	Gateway PodStatusMap `json:"gateway,omitempty"`

	// Ruler is a map to the per pod status of the lokistack ruler statefulset.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:com.tectonic.ui:podStatuses",displayName="Ruler",order=6
	Ruler PodStatusMap `json:"ruler,omitempty"`
}

// PodStatusMap defines the type for mapping pod status to pod name.
type PodStatusMap map[PodStatus][]string

// PodStatus is a short description of the status a Pod can be in.
type PodStatus string

const (
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodStatus = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodStatus = "Running"
	// PodReady means the pod has been started and the readiness probe reports a successful status.
	PodReady PodStatus = "Ready"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodStatus = "Failed"
	// PodStatusUnknown is used when none of the other statuses apply or the information is not ready yet.
	PodStatusUnknown PodStatus = "Unknown"
)

// LokiComponentSpec defines the requirements to configure scheduling
// of each loki component individually.
type LokiComponentSpec struct {
	// Replicas defines the number of replica pods of the component.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:hidden"
	Replicas int32 `json:"replicas,omitempty"`

	// NodeSelector defines the labels required by a node to schedule
	// the component onto it.
	//
	// +optional
	// +kubebuilder:validation:Optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Tolerations defines the tolerations required by a node to schedule
	// the component onto it.
	//
	// +optional
	// +kubebuilder:validation:Optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// PodAntiAffinity defines the pod anti affinity scheduling rules to schedule pods
	// of a component.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:podAntiAffinity",displayName="PodAntiAffinity"
	PodAntiAffinity *corev1.PodAntiAffinity `json:"podAntiAffinity,omitempty"`
}

// LokiTemplateSpec defines the template of all requirements to configure
// scheduling of all Loki components to be deployed.
type LokiTemplateSpec struct {
	// Compactor defines the compaction component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Compactor pods"
	Compactor *LokiComponentSpec `json:"compactor,omitempty"`

	// Distributor defines the distributor component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Distributor pods"
	Distributor *LokiComponentSpec `json:"distributor,omitempty"`

	// Ingester defines the ingester component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Ingester pods"
	Ingester *LokiComponentSpec `json:"ingester,omitempty"`

	// Querier defines the querier component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Querier pods"
	Querier *LokiComponentSpec `json:"querier,omitempty"`

	// QueryFrontend defines the query frontend component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Query Frontend pods"
	QueryFrontend *LokiComponentSpec `json:"queryFrontend,omitempty"`

	// Gateway defines the lokistack gateway component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Gateway pods"
	Gateway *LokiComponentSpec `json:"gateway,omitempty"`

	// IndexGateway defines the index gateway component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Index Gateway pods"
	IndexGateway *LokiComponentSpec `json:"indexGateway,omitempty"`

	// Ruler defines the ruler component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Ruler pods"
	Ruler *LokiComponentSpec `json:"ruler,omitempty"`
}

// LokiStackConditionType deifnes the type of condition types of a Loki deployment.
type LokiStackConditionType string

const (
	// ConditionReady defines the condition that all components in the Loki deployment are ready.
	ConditionReady LokiStackConditionType = "Ready"

	// ConditionPending defines the condition that some or all components are in pending state.
	ConditionPending LokiStackConditionType = "Pending"

	// ConditionFailed defines the condition that components in the Loki deployment failed to roll out.
	ConditionFailed LokiStackConditionType = "Failed"

	// ConditionDegraded defines the condition that some or all components in the Loki deployment
	// are degraded or the cluster cannot connect to object storage.
	ConditionDegraded LokiStackConditionType = "Degraded"

	// ConditionWarning is used for configurations that are not recommended, but don't currently cause
	// issues. There can be multiple warning conditions active at a time.
	ConditionWarning LokiStackConditionType = "Warning"
)

// LokiStackConditionReason defines the type for valid reasons of a Loki deployment conditions.
type LokiStackConditionReason string

const (
	// ReasonFailedComponents when all/some LokiStack components fail to roll out.
	ReasonFailedComponents LokiStackConditionReason = "FailedComponents"
	// ReasonPendingComponents when all/some LokiStack components pending dependencies
	ReasonPendingComponents LokiStackConditionReason = "PendingComponents"
	// ReasonReadyComponents when all LokiStack components are ready to serve traffic.
	ReasonReadyComponents LokiStackConditionReason = "ReadyComponents"
	// ReasonMissingObjectStorageSecret when the required secret to store logs to object
	// storage is missing.
	ReasonMissingObjectStorageSecret LokiStackConditionReason = "MissingObjectStorageSecret"
	// ReasonInvalidObjectStorageSchema when the spec contains an invalid schema(s).
	ReasonInvalidObjectStorageSchema LokiStackConditionReason = "InvalidObjectStorageSchema"
	// ReasonInvalidObjectStorageSecret when the format of the secret is invalid.
	ReasonInvalidObjectStorageSecret LokiStackConditionReason = "InvalidObjectStorageSecret"
	// ReasonZoneAwareNodesMissing when the cluster does not contain any nodes with the labels needed for zone-awareness.
	ReasonZoneAwareNodesMissing LokiStackConditionReason = "ReasonZoneAwareNodesMissing"
	// ReasonZoneAwareEmptyLabel when the node-label used for zone-awareness has an empty value.
	ReasonZoneAwareEmptyLabel LokiStackConditionReason = "ReasonZoneAwareEmptyLabel"
	// ReasonStorageNeedsSchemaUpdate when the object storage schema version is older than V13
	ReasonStorageNeedsSchemaUpdate LokiStackConditionReason = "StorageNeedsSchemaUpdate"
)

// LokiStackStorageStatus defines the observed state of
// the Loki storage configuration.
type LokiStackStorageStatus struct {
	// Schemas is a list of schemas which have been applied
	// to the LokiStack.
	//
	// +optional
	// +kubebuilder:validation:Optional
	Schemas []ObjectStorageSchema `json:"schemas,omitempty"`

	// CredentialMode contains the authentication mode used for accessing the object storage.
	//
	// +optional
	// +kubebuilder:validation:Optional
	CredentialMode CredentialMode `json:"credentialMode,omitempty"`
}

// LimitsSpec defines the spec for limits applied at ingestion or query
// path across the cluster or per tenant.
type LimitsSpec struct {
	// Global defines the limits applied globally across the cluster.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Global Limits"
	Global *LimitsTemplateSpec `json:"global,omitempty"`
}

// LimitsTemplateSpec defines the limits  applied at ingestion or query path.
type LimitsTemplateSpec struct {
	// IngestionLimits defines the limits applied on ingested log streams.
	//
	// +optional
	// +kubebuilder:validation:Optional
	IngestionLimits *IngestionLimitSpec `json:"ingestion,omitempty"`

	// QueryLimits defines the limit applied on querying log streams.
	//
	// +optional
	// +kubebuilder:validation:Optional
	QueryLimits *QueryLimitSpec `json:"queries,omitempty"`
}

// IngestionLimitSpec defines the limits applied at the ingestion path.
type IngestionLimitSpec struct {
	// IngestionRate defines the sample size per second. Units MB.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Ingestion Rate (in MB)"
	IngestionRate int32 `json:"ingestionRate,omitempty"`

	// IngestionBurstSize defines the local rate-limited sample size per
	// distributor replica. It should be set to the set at least to the
	// maximum logs size expected in a single push request.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Ingestion Burst Size (in MB)"
	IngestionBurstSize int32 `json:"ingestionBurstSize,omitempty"`

	// MaxLabelNameLength defines the maximum number of characters allowed
	// for label keys in log streams.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Label Name Length"
	MaxLabelNameLength int32 `json:"maxLabelNameLength,omitempty"`

	// MaxLabelValueLength defines the maximum number of characters allowed
	// for label values in log streams.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Label Value Length"
	MaxLabelValueLength int32 `json:"maxLabelValueLength,omitempty"`

	// MaxLabelNamesPerSeries defines the maximum number of label names per series
	// in each log stream.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Labels Names per Series"
	MaxLabelNamesPerSeries int32 `json:"maxLabelNamesPerSeries,omitempty"`

	// MaxGlobalStreamsPerTenant defines the maximum number of active streams
	// per tenant, across the cluster.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Global Streams per  Tenant"
	MaxGlobalStreamsPerTenant int32 `json:"maxGlobalStreamsPerTenant,omitempty"`

	// MaxLineSize defines the maximum line size on ingestion path. Units in Bytes.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Line Size"
	MaxLineSize int32 `json:"maxLineSize,omitempty"`

	// PerStreamDesiredRate defines the desired ingestion rate per second that LokiStack should
	// target applying automatic stream sharding. Units MB.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Per Stream Desired Rate (in MB)"
	PerStreamDesiredRate int32 `json:"perStreamDesiredRate,omitempty"`

	// PerStreamRateLimit defines the maximum byte rate per second per stream. Units MB.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Maximum byte rate per second per stream (in MB)"
	PerStreamRateLimit int32 `json:"perStreamRateLimit,omitempty"`

	// PerStreamRateLimitBurst defines the maximum burst bytes per stream. Units MB.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Maximum burst bytes per stream (in MB)"
	PerStreamRateLimitBurst int32 `json:"perStreamRateLimitBurst,omitempty"`
}

// QueryLimitSpec defines the limits applies at the query path.
type QueryLimitSpec struct {
	// MaxEntriesLimitsPerQuery defines the maximum number of log entries
	// that will be returned for a query.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Entries Limit per Query"
	MaxEntriesLimitPerQuery int32 `json:"maxEntriesLimitPerQuery,omitempty"`

	// MaxChunksPerQuery defines the maximum number of chunks
	// that can be fetched by a single query.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Chunk per Query"
	MaxChunksPerQuery int32 `json:"maxChunksPerQuery,omitempty"`

	// MaxQuerySeries defines the maximum of unique series
	// that is returned by a metric query.
	//
	// + optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Query Series"
	MaxQuerySeries int32 `json:"maxQuerySeries,omitempty"`

	// Timeout when querying ingesters or storage during the execution of a query request.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="3m"
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Query Timeout"
	QueryTimeout string `json:"queryTimeout,omitempty"`

	// CardinalityLimit defines the cardinality limit for index queries.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Cardinality Limit"
	CardinalityLimit int32 `json:"cardinalityLimit,omitempty"`

	// MaxVolumeSeries defines the maximum number of aggregated series in a log-volume response
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Volume Series"
	MaxVolumeSeries int32 `json:"maxVolumeSeries,omitempty"`
}

// CredentialMode represents the type of authentication used for accessing the object storage.
//
// +kubebuilder:validation:Enum=static;token;token-cco
type CredentialMode string

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LokiStack is the Schema for the LokiStacks API
type LokiStack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LokiStackSpec   `json:"spec,omitempty"`
	Status LokiStackStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LokiStackList contains a list of LokiStack
type LokiStackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LokiStack `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LokiStack{}, &LokiStackList{})
}
