# DataCatalog


Initiating the Datahub data catalog

```
cd /Users/aditya.relangi/Code/datahub/docker
```

```
./quickstart.sh
```

The following are the commands to interact with the Acryl Datahub Catalog.

### Create a new dataset

```
curl 'http://localhost:8080/entities?action=ingest' -X POST --data '
{
   "entity":{
      "value":{
         "com.linkedin.metadata.snapshot.DatasetSnapshot":{
            "urn":"urn:li:dataset:(urn:li:dataPlatform:fooe,UserReferenceExample,PROD)",
            "aspects":[
               {
                  "com.linkedin.common.Ownership":{
                     "owners":[
                        {
                           "owner":"urn:li:corpuser:Zarina Malik",
                           "type":"DATAOWNER"
                        },
                        {
                           "owner":"urn:li:corpuser:Aravind Rammoorthy",
                           "type":"PRODUCER"
                        },
                        {
                           "owner":"urn:li:corpuser:Prasanth Kanaujia",
                           "type":"DEVELOPER"
                        }
                     ]
                  }
               },
               {
                  "com.linkedin.common.InstitutionalMemory":{
                     "elements":[
                        {
                           "url":"https://www.linkedin.com",
                           "description":"User object represents a user of the survey funky platform",
                           "createStamp":{
                              "time":0,
                              "actor":"urn:li:corpuser:Zarina Malik"
                           }
                        }
                     ]
                  }
               },
               {
                 "com.linkedin.dataset.UpstreamLineage": {
                   "upstreams": [
                     {
                       "dataset": "urn:li:dataset:(urn:li:dataPlatform:kafka,user,PROD)",
                       "type": "TRANSFORMED"
                     }
                   ]
                 }
               },
               {
                  "com.linkedin.dataset.DatasetProperties":{
                      "customProperties": {
                            "html_wrapper": "<b> | </b>"
                        },
                        "externalUrl": "noll",
                        "description": "Kaun User7 description",
                        "tags": []
                  }
               },
               {
                  "com.linkedin.schema.SchemaMetadata":{
                     "schemaName":"FooEvent",
                     "platform":"urn:li:dataPlatform:foo",
                     "version":0,
                     "created":{
                        "time":0,
                        "actor":"urn:li:corpuser:Zarina Malik"
                     },
                     "hash":"",
                     "platformSchema":{
                        "com.linkedin.schema.KafkaSchema":{
                           "documentSchema":"{\"type\":\"record\",\"namespace\":\"com.surveyfunky\",\"name\":\"User\",\"fields\":[{\"name\":\"user_id\",\"doc\":\"Uniqueidentifieroftheuser\"},{\"name\":\"first_name\",\"type\":\"string\",\"doc\":\"FirstNameofUser\"},{\"name\":\"last_name\",\"type\":\"string\",\"doc\":\"LastNameofUser\"},{\"name\":\"created_date\",\"type\":[\"null\",{\"type\":\"string\",\"logicalType\":\"date\"}],\"doc\":\"Dateonwhichthisrecordiscreated\"},{\"name\":\"last_updated_time\",\"type\":[\"null\",{\"type\":\"string\",\"logicalType\":\"timestamp-millis\"}],\"doc\":\"Timestampwhenthisrecordwasmostrecentlyupdated\"},{\"name\":\"last_login\",\"type\":[\"null\",{\"type\":\"string\",\"logicalType\":\"timestamp-millis\"}],\"doc\":\"Timestampofthelatesttimetheuserhasloggedin\"}]}"
                        }
                     },
                     "primaryKeys": [
                         "user_id"
                     ],
                     "foreignKeys":[{
                        "name": "sample foreign key",
                        "foreignFields": [
                           "urn:li:schemaField:(urn:li:dataset:(urn:li:dataPlatform:kafka,user,PROD),user_id)"
                         ],
                         "sourceFields": [
                           "urn:li:schemaField:(urn:li:dataset:(urn:li:dataPlatform:fooe,UserReferenceExample,PROD),user_id)"
                         ],
                         "foreignDataset": "urn:li:dataset:(urn:li:dataPlatform:kafka,user,PROD)"
                     }],
                     "fields":[
                        {
                           "fieldPath":"user_id",
                           "description":"Unique identifier of the user",
                           "nativeDataType":"int",
                           "type":{
                              "type":{
                                 "com.linkedin.schema.NumberType":{
                                    
                                 }
                              }
                           }
                        },
                        {
                           "fieldPath":"first_name",
                           "description":"First name of the user",
                           "nativeDataType":"string",
                           "type":{
                              "type":{
                                 "com.linkedin.schema.StringType":{
                                    
                                 }
                              }
                           }
                        },
                        {
                           "fieldPath":"last_name",
                           "description":"Last name of the user",
                           "nativeDataType":"string",
                           "type":{
                              "type":{
                                 "com.linkedin.schema.StringType":{
                                    
                                 }
                              }
                           }
                        },
                        {
                           "fieldPath":"created_date",
                           "description":"Date on which this record is created",
                           "nativeDataType":"string",
                           "type":{
                              "type":{
                                 "com.linkedin.schema.DateType":{
                                    
                                 }
                              }
                           }
                        },
                        {
                           "fieldPath":"last_updated_time",
                           "description":"Timestamp when this record was most recently updated",
                           "nativeDataType":"string",
                           "type":{
                              "type":{
                                 "com.linkedin.schema.TimeType":{
                                    
                                 }
                              }
                           }
                        },
                        {
                           "fieldPath":"last_login",
                           "description":"Timestamp of the latest time the user has logged in",
                           "nativeDataType":"string",
                           "type":{
                              "type":{
                                 "com.linkedin.schema.TimeType":{
                                    
                                 }
                              }
                           }
                        }

                     ]
                  }
               }
            ]
         }
      }
   }
}'
```


### Retrieve the dataset

```
curl  'http://localhost:8080/entities/urn:li:dataset:(urn:li:dataPlatform:fooe,User,PROD)'


```



### ElasticSearch commands

[Kafka-connect with Elasticsearch](https://github.com/confluentinc/demo-scene/tree/master/kafka-to-elasticsearch)


```

curl -X POST http://localhost:8083/connectors \
-H "Content-type:application/json" \
--data-raw '{
  "name": "SINK_ELASTIC_TEST_05",
  "config":
  {
    "connector.class"                     : "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
    "connection.url"                      : "http://elasticsearch:9200",
    "value.converter"                     : "io.confluent.connect.avro.AvroConverter",
    "value.converter.schema.registry.url" : "http://schemaregistry:8082",
    "type.name"                           : "_doc",
    "topics"                              : "user",
    "key.ignore"                          : "true",
    "schema.ignore"                       : "true"
  }
}'
```

```
curl -s http://localhost:9200/user/_search \
    -H 'content-type: application/json' \
    -d '{ "size": 42  }' | jq -c '.hits.hits[]'
```


