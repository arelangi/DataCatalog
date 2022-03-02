package main

type DataHubRequestType struct {
	Entity DataHubEntityType `json:"entity"`
}

type DataHubEntityType struct {
	Value DataHubEntityValueType `json:"value"`
}

type DataHubEntityValueType struct {
	DatasetSnapshot DatasetSnapshotType `json:"com.linkedin.metadata.snapshot.DatasetSnapshot"`
}

type DatasetSnapshotType struct {
	Urn     string       `json:"urn"`
	Aspects []AspectType `json:"aspects"`
}

type AspectType struct {
	Ownership           *OwnershipType           `json:"com.linkedin.common.Ownership,omitempty"`
	InstitutionalMemory *InstitutionalMemoryType `json:"com.linkedin.common.InstitutionalMemory,omitempty"`
	UpstreamLineage     *UpstreamLineageType     `json:"com.linkedin.dataset.UpstreamLineage,omitempty"`
	DatasetProperties   *DatasetPropertiesType   `json:"com.linkedin.dataset.DatasetProperties,omitempty"`
	SchemaMetaData      *SchemaMetadataType      `json:"com.linkedin.schema.SchemaMetadata,omitempty"`
}

type OwnershipType struct {
	Owners []OwnerType `json:"owners,omitempty"`
}

type OwnerType struct {
	Owner string `json:"owner,omitempty"`
	Type  string `json:"type,omitempty"`
}

type InstitutionalMemoryType struct {
	Elements []InstitutionalMemoryElementType `json:"elements,omitempty"`
}

type InstitutionalMemoryElementType struct {
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	CreateStamp AuditStamp `json:"createStamp,omitempty"`
}

type UpstreamLineageType struct {
	Upstreams []struct {
		Dataset string `json:"dataset,omitempty"`
		Type    string `json:"type,omitempty"`
	} `json:"upstreams,omitempty"`
}

type DatasetPropertiesType struct {
	CustomProperties map[string]string `json:"customProperties,omitempty"`
	ExternalURL string        `json:"externalUrl,omitempty"`
	Description string        `json:"description,omitempty"`
	Tags        []string `json:"tags"`
}

type SchemaMetadataType struct {
	SchemaName string `json:"schemaName,omitempty"`
	Platform   string `json:"platform,omitempty"`
	Version    int    `json:"version"`
	Created    AuditStamp `json:"created,omitempty"`
	Hash           string `json:"hash"`
	PlatformSchema struct {
		ComLinkedinSchemaKafkaSchema struct {
			DocumentSchema string `json:"documentSchema"`
		} `json:"com.linkedin.schema.KafkaSchema"`
	} `json:"platformSchema"`
	Fields []DatahubFieldType `json:"fields,omitempty"`
}

type DatahubFieldType struct{
	FieldPath      string `json:"fieldPath,omitempty"`
	Description    string `json:"description,omitempty"`
	NativeDataType string `json:"nativeDataType,omitempty"`
	Type           SchemaFieldDataType `json:"type,omitempty"`
}

type SchemaFieldDataType struct{
	Type LinkedinDataType`json:"type"`
}

type LinkedinDataType struct{
	StringType *struct{} `json:"com.linkedin.schema.StringType,omitempty"`
	NumberType *struct{} `json:"com.linkedin.schema.NumberType,omitempty"`
	DateType *struct{} `json:"com.linkedin.schema.DateType,omitempty"`
	TimeType *struct{} `json:"com.linkedin.schema.TimeType,omitempty"`
}

type AuditStamp struct{
	Time  int64    `json:"time,omitempty"`
	Actor string `json:"actor,omitempty"`
}


type DatahubAPIFailureResponse struct {
	ExceptionClass string `json:"exceptionClass,omitempty"`
	StackTrace     string `json:"stackTrace,omitempty"`
	Message        string `json:"message,omitempty"`
	Status         int    `json:"status,omitempty"`
}

