// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: proto/query-manager/server.proto

package query_manager

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MetricArg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueryType    string            `protobuf:"bytes,1,opt,name=query_type,json=queryType,proto3" json:"query_type,omitempty"`
	Metric       string            `protobuf:"bytes,2,opt,name=metric,proto3" json:"metric,omitempty"`
	StartTime    string            `protobuf:"bytes,3,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime      string            `protobuf:"bytes,4,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	TimeInterval string            `protobuf:"bytes,5,opt,name=time_interval,json=timeInterval,proto3" json:"time_interval,omitempty"`
	Filters      map[string]string `protobuf:"bytes,6,rep,name=filters,proto3" json:"filters,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *MetricArg) Reset() {
	*x = MetricArg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_manager_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetricArg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetricArg) ProtoMessage() {}

func (x *MetricArg) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_manager_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetricArg.ProtoReflect.Descriptor instead.
func (*MetricArg) Descriptor() ([]byte, []int) {
	return file_proto_query_manager_server_proto_rawDescGZIP(), []int{0}
}

func (x *MetricArg) GetQueryType() string {
	if x != nil {
		return x.QueryType
	}
	return ""
}

func (x *MetricArg) GetMetric() string {
	if x != nil {
		return x.Metric
	}
	return ""
}

func (x *MetricArg) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *MetricArg) GetEndTime() string {
	if x != nil {
		return x.EndTime
	}
	return ""
}

func (x *MetricArg) GetTimeInterval() string {
	if x != nil {
		return x.TimeInterval
	}
	return ""
}

func (x *MetricArg) GetFilters() map[string]string {
	if x != nil {
		return x.Filters
	}
	return nil
}

type EmptyArg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EmptyArg) Reset() {
	*x = EmptyArg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_manager_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmptyArg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmptyArg) ProtoMessage() {}

func (x *EmptyArg) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_manager_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmptyArg.ProtoReflect.Descriptor instead.
func (*EmptyArg) Descriptor() ([]byte, []int) {
	return file_proto_query_manager_server_proto_rawDescGZIP(), []int{1}
}

// Used in federating a service from one cluster to another
type TimeSeriesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// name of the group
	Code         uint32                               `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message      string                               `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Data         []*TimeSeriesResponse_TimeSeriesItem `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
	Last         string                               `protobuf:"bytes,4,opt,name=last,proto3" json:"last,omitempty"`
	TotalRecords uint32                               `protobuf:"varint,5,opt,name=total_records,json=totalRecords,proto3" json:"total_records,omitempty"`
}

func (x *TimeSeriesResponse) Reset() {
	*x = TimeSeriesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_manager_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeSeriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeSeriesResponse) ProtoMessage() {}

func (x *TimeSeriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_manager_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeSeriesResponse.ProtoReflect.Descriptor instead.
func (*TimeSeriesResponse) Descriptor() ([]byte, []int) {
	return file_proto_query_manager_server_proto_rawDescGZIP(), []int{2}
}

func (x *TimeSeriesResponse) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *TimeSeriesResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *TimeSeriesResponse) GetData() []*TimeSeriesResponse_TimeSeriesItem {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *TimeSeriesResponse) GetLast() string {
	if x != nil {
		return x.Last
	}
	return ""
}

func (x *TimeSeriesResponse) GetTotalRecords() uint32 {
	if x != nil {
		return x.TotalRecords
	}
	return 0
}

type TimeSeriesResponse_TimeSeriesItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timeslot string            `protobuf:"bytes,1,opt,name=timeslot,proto3" json:"timeslot,omitempty"`
	Data     map[string]string `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *TimeSeriesResponse_TimeSeriesItem) Reset() {
	*x = TimeSeriesResponse_TimeSeriesItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_query_manager_server_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeSeriesResponse_TimeSeriesItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeSeriesResponse_TimeSeriesItem) ProtoMessage() {}

func (x *TimeSeriesResponse_TimeSeriesItem) ProtoReflect() protoreflect.Message {
	mi := &file_proto_query_manager_server_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeSeriesResponse_TimeSeriesItem.ProtoReflect.Descriptor instead.
func (*TimeSeriesResponse_TimeSeriesItem) Descriptor() ([]byte, []int) {
	return file_proto_query_manager_server_proto_rawDescGZIP(), []int{2, 0}
}

func (x *TimeSeriesResponse_TimeSeriesItem) GetTimeslot() string {
	if x != nil {
		return x.Timeslot
	}
	return ""
}

func (x *TimeSeriesResponse_TimeSeriesItem) GetData() map[string]string {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_query_manager_server_proto protoreflect.FileDescriptor

var file_proto_query_manager_server_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2d, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0c, 0x71, 0x75, 0x65, 0x72, 0x79, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72,
	0x22, 0x9d, 0x02, 0x0a, 0x09, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x41, 0x72, 0x67, 0x12, 0x1d,
	0x0a, 0x0a, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x71, 0x75, 0x65, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x23, 0x0a, 0x0d, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x76, 0x61, 0x6c, 0x12, 0x3e, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x18,
	0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x4d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x41, 0x72, 0x67, 0x2e, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x73, 0x1a, 0x3a, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x0a, 0x0a, 0x08, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x41, 0x72, 0x67, 0x22, 0xf7, 0x02, 0x0a,
	0x12, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x43, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2f, 0x2e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x49, 0x74, 0x65, 0x6d,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x61, 0x73, 0x74, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6c, 0x61, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x5f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x1a,
	0xb4, 0x01, 0x0a, 0x0e, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x49, 0x74,
	0x65, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x6c, 0x6f, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x6c, 0x6f, 0x74, 0x12, 0x4d,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x39, 0x2e, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x49, 0x74, 0x65, 0x6d, 0x2e, 0x44, 0x61,
	0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x37, 0x0a,
	0x09, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x53, 0x0a, 0x06, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x12, 0x49, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x17,
	0x2e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x41, 0x72, 0x67, 0x1a, 0x20, 0x2e, 0x71, 0x75, 0x65, 0x72, 0x79, 0x4d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x11, 0x5a, 0x0f, 0x2e,
	0x2f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2d, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_query_manager_server_proto_rawDescOnce sync.Once
	file_proto_query_manager_server_proto_rawDescData = file_proto_query_manager_server_proto_rawDesc
)

func file_proto_query_manager_server_proto_rawDescGZIP() []byte {
	file_proto_query_manager_server_proto_rawDescOnce.Do(func() {
		file_proto_query_manager_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_query_manager_server_proto_rawDescData)
	})
	return file_proto_query_manager_server_proto_rawDescData
}

var file_proto_query_manager_server_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_query_manager_server_proto_goTypes = []interface{}{
	(*MetricArg)(nil),          // 0: queryManager.MetricArg
	(*EmptyArg)(nil),           // 1: queryManager.EmptyArg
	(*TimeSeriesResponse)(nil), // 2: queryManager.TimeSeriesResponse
	nil,                        // 3: queryManager.MetricArg.FiltersEntry
	(*TimeSeriesResponse_TimeSeriesItem)(nil), // 4: queryManager.TimeSeriesResponse.TimeSeriesItem
	nil, // 5: queryManager.TimeSeriesResponse.TimeSeriesItem.DataEntry
}
var file_proto_query_manager_server_proto_depIdxs = []int32{
	3, // 0: queryManager.MetricArg.filters:type_name -> queryManager.MetricArg.FiltersEntry
	4, // 1: queryManager.TimeSeriesResponse.data:type_name -> queryManager.TimeSeriesResponse.TimeSeriesItem
	5, // 2: queryManager.TimeSeriesResponse.TimeSeriesItem.data:type_name -> queryManager.TimeSeriesResponse.TimeSeriesItem.DataEntry
	0, // 3: queryManager.Server.GetMetrics:input_type -> queryManager.MetricArg
	2, // 4: queryManager.Server.GetMetrics:output_type -> queryManager.TimeSeriesResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_query_manager_server_proto_init() }
func file_proto_query_manager_server_proto_init() {
	if File_proto_query_manager_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_query_manager_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetricArg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_query_manager_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmptyArg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_query_manager_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeSeriesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_query_manager_server_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeSeriesResponse_TimeSeriesItem); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_query_manager_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_query_manager_server_proto_goTypes,
		DependencyIndexes: file_proto_query_manager_server_proto_depIdxs,
		MessageInfos:      file_proto_query_manager_server_proto_msgTypes,
	}.Build()
	File_proto_query_manager_server_proto = out.File
	file_proto_query_manager_server_proto_rawDesc = nil
	file_proto_query_manager_server_proto_goTypes = nil
	file_proto_query_manager_server_proto_depIdxs = nil
}
