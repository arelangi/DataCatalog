package main

import (
	"fmt"
	"strings"
	"time"
)

func (m *MetadataRequest) getSchemaMetadataAspect(a *App) (s SchemaMetadataType) {
	developerName := "Eminem"

	s.SchemaName = m.DatasetName
	s.Platform = fmt.Sprintf("urn:li:dataPlatform:%s", strings.ToLower(m.PlatformName))
	s.Version = 0
	s.Created = AuditStamp{
		Time:  time.Now().Unix(),
		Actor: fmt.Sprintf("urn:li:corpuser:%s", developerName),
	}
	//To-do: Pull the schema from the metadata service? or maybe save it to a database table

	dataset, _ := a.getCompleteDatasetByID(m.DatasetID)

	fmt.Println("This dataset the following number of fields ", len(dataset.Fields))

	for _, v := range dataset.Fields {
		fmt.Println("Looping on field ", v.Name)
		field := DatahubFieldType{
			FieldPath:      v.Name,
			Description:    v.Doc,
			NativeDataType: strings.ToLower(v.Type),
			Type:           SchemaFieldDataType{Type: nativeTypeToLinkedinTypeMapping(v.Type)},
		}

		s.Fields = append(s.Fields, field)
	}

	fmt.Println("Return object now has these many fields ", len(s.Fields))

	return
}

func nativeTypeToLinkedinTypeMapping(fieldType string) (res LinkedinDataType) {
	fieldType = strings.ToLower(fieldType)
	if strings.Contains(fieldType, "int") {
		res = LinkedinDataType{
			NumberType: &struct{}{},
		}
	} else if strings.Contains(fieldType, "string") {
		res = LinkedinDataType{
			StringType: &struct{}{},
		}
	}
	return
}

func (m *MetadataRequest) getDatasetPropertiesAspect() (d DatasetPropertiesType) {
	d.ExternalURL = fmt.Sprintf("http://datacatalog:3000/register/review/%d", m.DatasetID)
	d.Description = m.DatasetDescription
	d.Tags = []string{}
	d.CustomProperties = map[string]string{
		"Testing for html":   "<h1>Heading 1</h1><br/><h2>Heading 2</h2>",
		"Just regular stuff": "dassadsadsd",
	}

	//CustomProperties is going to be populated when we have the custom data saved

	return
}

//To-do finish this off
func (m *MetadataRequest) getUpstreamLineageAspect() (d UpstreamLineageType) {
	return
}

func (m *MetadataRequest) getInstitutionalMemoryAspect() (a InstitutionalMemoryType) {
	var element InstitutionalMemoryElementType
	developerName := "Eminem"

	/*
		URL is referring back to data catalog review page
		Description is the dataset description
		AuditStamp is based on the developer creating this
	*/
	element.URL = fmt.Sprintf("http://datacatalog:3000/register/review/%d", m.DatasetID)
	element.Description = m.DatasetDescription
	element.CreateStamp = AuditStamp{
		Time:  time.Now().Unix(),
		Actor: fmt.Sprintf("urn:li:corpuser:%s", developerName),
	}

	a.Elements = append(a.Elements, element)
	return
}

func (m *MetadataRequest) getOwnerAspect() (a OwnershipType) {

	//Build Ownership aspect
	/*
		Producer is the team
		DataOwner is the steward
		Developer is whoever is registering the dataset
	*/

	developerName := "Eminem"

	//Team is the data producer
	producer := OwnerType{
		Type:  "PRODUCER",
		Owner: fmt.Sprintf("urn:li:corpuser:%s", m.Team),
	}

	dataOwner := OwnerType{
		Type:  "DATAOWNER",
		Owner: fmt.Sprintf("urn:li:corpuser:%s", m.DataSteward),
	}

	developer := OwnerType{
		Type:  "DEVELOPER",
		Owner: fmt.Sprintf("urn:li:corpuser:%s", developerName),
	}

	a.Owners = append(a.Owners, producer, dataOwner, developer)

	return
}
