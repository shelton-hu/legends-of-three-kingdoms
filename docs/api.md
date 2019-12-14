# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [0.proto](#0.proto)
  
  
  
  

- [game.proto](#game.proto)
    - [StartGameRequets](#ltk.StartGameRequets)
    - [StartGameResponse](#ltk.StartGameResponse)
  
  
  
    - [GameService](#ltk.GameService)
  

- [iam.proto](#iam.proto)
    - [SignInOrSignUpRequest](#ltk.SignInOrSignUpRequest)
    - [SignInOrSignUpResponse](#ltk.SignInOrSignUpResponse)
    - [SignOutRequest](#ltk.SignOutRequest)
    - [SignOutResponse](#ltk.SignOutResponse)
  
  
  
    - [TokenService](#ltk.TokenService)
  

- [room.proto](#room.proto)
    - [ComeIntoRoomRequest](#ltk.ComeIntoRoomRequest)
    - [ComeIntoRoomResponse](#ltk.ComeIntoRoomResponse)
    - [CreateRoomRequest](#ltk.CreateRoomRequest)
    - [CreateRoomResponse](#ltk.CreateRoomResponse)
    - [DescribeRoomsRequest](#ltk.DescribeRoomsRequest)
    - [DescribeRoomsResponse](#ltk.DescribeRoomsResponse)
    - [Room](#ltk.Room)
  
  
  
    - [RoomService](#ltk.RoomService)
  

- [types.proto](#types.proto)
    - [ErrorDetail](#ltk.ErrorDetail)
  
  
  
  

- [Scalar Value Types](#scalar-value-types)



<a name="0.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## 0.proto


 

 

 

 



<a name="game.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## game.proto



<a name="ltk.StartGameRequets"></a>

### StartGameRequets







<a name="ltk.StartGameResponse"></a>

### StartGameResponse






 

 

 


<a name="ltk.GameService"></a>

### GameService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| StartGame | [StartGameRequets](#ltk.StartGameRequets) stream | [StartGameResponse](#ltk.StartGameResponse) |  |

 



<a name="iam.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iam.proto



<a name="ltk.SignInOrSignUpRequest"></a>

### SignInOrSignUpRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| nick_name | [google.protobuf.StringValue](#google.protobuf.StringValue) |  |  |
| password | [google.protobuf.StringValue](#google.protobuf.StringValue) |  |  |






<a name="ltk.SignInOrSignUpResponse"></a>

### SignInOrSignUpResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [google.protobuf.StringValue](#google.protobuf.StringValue) |  |  |






<a name="ltk.SignOutRequest"></a>

### SignOutRequest







<a name="ltk.SignOutResponse"></a>

### SignOutResponse






 

 

 


<a name="ltk.TokenService"></a>

### TokenService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SignInOrSignUp | [SignInOrSignUpRequest](#ltk.SignInOrSignUpRequest) | [SignInOrSignUpResponse](#ltk.SignInOrSignUpResponse) |  |
| SignOut | [SignOutRequest](#ltk.SignOutRequest) | [SignOutResponse](#ltk.SignOutResponse) |  |

 



<a name="room.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## room.proto



<a name="ltk.ComeIntoRoomRequest"></a>

### ComeIntoRoomRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| room_id | [string](#string) |  |  |






<a name="ltk.ComeIntoRoomResponse"></a>

### ComeIntoRoomResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| room_id | [google.protobuf.StringValue](#google.protobuf.StringValue) |  |  |






<a name="ltk.CreateRoomRequest"></a>

### CreateRoomRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| room_nick_name | [google.protobuf.StringValue](#google.protobuf.StringValue) |  |  |






<a name="ltk.CreateRoomResponse"></a>

### CreateRoomResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| room_id | [google.protobuf.StringValue](#google.protobuf.StringValue) |  |  |






<a name="ltk.DescribeRoomsRequest"></a>

### DescribeRoomsRequest







<a name="ltk.DescribeRoomsResponse"></a>

### DescribeRoomsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| rooms | [Room](#ltk.Room) | repeated |  |






<a name="ltk.Room"></a>

### Room



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| room_id | [google.protobuf.StringValue](#google.protobuf.StringValue) |  |  |
| room_nick_name | [google.protobuf.StringValue](#google.protobuf.StringValue) |  |  |





 

 

 


<a name="ltk.RoomService"></a>

### RoomService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateRoom | [CreateRoomRequest](#ltk.CreateRoomRequest) | [CreateRoomResponse](#ltk.CreateRoomResponse) |  |
| ComeIntoRoom | [ComeIntoRoomRequest](#ltk.ComeIntoRoomRequest) | [ComeIntoRoomResponse](#ltk.ComeIntoRoomResponse) |  |
| DescribeRooms | [DescribeRoomsRequest](#ltk.DescribeRoomsRequest) | [DescribeRoomsResponse](#ltk.DescribeRoomsResponse) |  |

 



<a name="types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## types.proto



<a name="ltk.ErrorDetail"></a>

### ErrorDetail



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| error_name | [string](#string) |  |  |
| cause | [string](#string) |  |  |





 

 

 

 



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

