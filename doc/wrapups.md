# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [pkg/wrapups/wrapups.proto](#pkg/wrapups/wrapups.proto)
    - [CreateWrapupRequest](#wrapups.CreateWrapupRequest)
    - [GetWrapupRequest](#wrapups.GetWrapupRequest)
    - [ListWrapupsRequest](#wrapups.ListWrapupsRequest)
    - [ListWrapupsResponse](#wrapups.ListWrapupsResponse)
    - [Wrapup](#wrapups.Wrapup)
  
  
  
    - [Wrapups](#wrapups.Wrapups)
  

- [Scalar Value Types](#scalar-value-types)



<a name="pkg/wrapups/wrapups.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## pkg/wrapups/wrapups.proto
wrapups.proto

Define Wrapups service and related messaages.


<a name="wrapups.CreateWrapupRequest"></a>

### CreateWrapupRequest
CreateWrapupRequest represents the request message for Create operation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title | [string](#string) |  | title of paper. |
| wrapup | [string](#string) |  | wrapup of paper. |
| comment | [string](#string) |  | comment of paper. |
| note | [string](#string) |  | note of paper. |






<a name="wrapups.GetWrapupRequest"></a>

### GetWrapupRequest
GetWrapupRequest represents the request message for Get operation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id to fetch specific wrapup object from Elasticsearch server. |






<a name="wrapups.ListWrapupsRequest"></a>

### ListWrapupsRequest
ListWrapupsRequest represents the request message for List operation.






<a name="wrapups.ListWrapupsResponse"></a>

### ListWrapupsResponse
ListWrapupsResponse represents the response of List operation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| count | [int32](#int32) |  | number of wrapup objects included in this response. |
| wrapups | [Wrapup](#wrapups.Wrapup) | repeated | list of wrapup object. |






<a name="wrapups.Wrapup"></a>

### Wrapup
Wrapup represents one wrapup object.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | ID of the wrapup object assigned by Elasticsearch. |
| title | [string](#string) |  | title of the paper. |
| wrapup | [string](#string) |  | wrapup of the paper. |
| comment | [string](#string) |  | comment of the paper. |
| note | [string](#string) |  | notes of the paper. |
| create_time | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  | timestamp which indicates when this wrapup object is created. |





 

 

 


<a name="wrapups.Wrapups"></a>

### Wrapups
Service for handling wrapup management.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListWrapups | [ListWrapupsRequest](#wrapups.ListWrapupsRequest) | [ListWrapupsResponse](#wrapups.ListWrapupsResponse) | ListWrapups returns the list of wrapup document stored in Elasticsearch. |
| GetWrapup | [GetWrapupRequest](#wrapups.GetWrapupRequest) | [Wrapup](#wrapups.Wrapup) | GetWrapup returns a wrapup document matched to request. |
| CreateWrapup | [CreateWrapupRequest](#wrapups.CreateWrapupRequest) | [Wrapup](#wrapups.Wrapup) | CreateWrapup creates new wrapup document and stores it in Elasticsearch. |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

