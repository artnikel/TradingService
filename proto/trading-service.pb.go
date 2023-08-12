// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: trading-service.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Deal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DealID        string                 `protobuf:"bytes,1,opt,name=DealID,proto3" json:"DealID,omitempty"`
	SharesCount   float64                `protobuf:"fixed64,2,opt,name=SharesCount,proto3" json:"SharesCount,omitempty"`
	ProfileID     string                 `protobuf:"bytes,3,opt,name=ProfileID,proto3" json:"ProfileID,omitempty"`
	Company       string                 `protobuf:"bytes,4,opt,name=Company,proto3" json:"Company,omitempty"`
	PurchasePrice float64                `protobuf:"fixed64,5,opt,name=PurchasePrice,proto3" json:"PurchasePrice,omitempty"`
	StopLoss      float64                `protobuf:"fixed64,6,opt,name=StopLoss,proto3" json:"StopLoss,omitempty"`
	TakeProfit    float64                `protobuf:"fixed64,7,opt,name=TakeProfit,proto3" json:"TakeProfit,omitempty"`
	DealTime      *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=DealTime,proto3" json:"DealTime,omitempty"`
	EndDealTime   *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=EndDealTime,proto3" json:"EndDealTime,omitempty"`
	Profit        float64                `protobuf:"fixed64,10,opt,name=Profit,proto3" json:"Profit,omitempty"`
}

func (x *Deal) Reset() {
	*x = Deal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deal) ProtoMessage() {}

func (x *Deal) ProtoReflect() protoreflect.Message {
	mi := &file_trading_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deal.ProtoReflect.Descriptor instead.
func (*Deal) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{0}
}

func (x *Deal) GetDealID() string {
	if x != nil {
		return x.DealID
	}
	return ""
}

func (x *Deal) GetSharesCount() float64 {
	if x != nil {
		return x.SharesCount
	}
	return 0
}

func (x *Deal) GetProfileID() string {
	if x != nil {
		return x.ProfileID
	}
	return ""
}

func (x *Deal) GetCompany() string {
	if x != nil {
		return x.Company
	}
	return ""
}

func (x *Deal) GetPurchasePrice() float64 {
	if x != nil {
		return x.PurchasePrice
	}
	return 0
}

func (x *Deal) GetStopLoss() float64 {
	if x != nil {
		return x.StopLoss
	}
	return 0
}

func (x *Deal) GetTakeProfit() float64 {
	if x != nil {
		return x.TakeProfit
	}
	return 0
}

func (x *Deal) GetDealTime() *timestamppb.Timestamp {
	if x != nil {
		return x.DealTime
	}
	return nil
}

func (x *Deal) GetEndDealTime() *timestamppb.Timestamp {
	if x != nil {
		return x.EndDealTime
	}
	return nil
}

func (x *Deal) GetProfit() float64 {
	if x != nil {
		return x.Profit
	}
	return 0
}

type StrategiesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Strategy string `protobuf:"bytes,1,opt,name=strategy,proto3" json:"strategy,omitempty"`
	Deal     *Deal  `protobuf:"bytes,2,opt,name=deal,proto3" json:"deal,omitempty"`
}

func (x *StrategiesRequest) Reset() {
	*x = StrategiesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StrategiesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StrategiesRequest) ProtoMessage() {}

func (x *StrategiesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trading_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StrategiesRequest.ProtoReflect.Descriptor instead.
func (*StrategiesRequest) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{1}
}

func (x *StrategiesRequest) GetStrategy() string {
	if x != nil {
		return x.Strategy
	}
	return ""
}

func (x *StrategiesRequest) GetDeal() *Deal {
	if x != nil {
		return x.Deal
	}
	return nil
}

type StrategiesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Profit float64 `protobuf:"fixed64,1,opt,name=Profit,proto3" json:"Profit,omitempty"`
}

func (x *StrategiesResponse) Reset() {
	*x = StrategiesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StrategiesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StrategiesResponse) ProtoMessage() {}

func (x *StrategiesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trading_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StrategiesResponse.ProtoReflect.Descriptor instead.
func (*StrategiesResponse) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{2}
}

func (x *StrategiesResponse) GetProfit() float64 {
	if x != nil {
		return x.Profit
	}
	return 0
}

var File_trading_service_proto protoreflect.FileDescriptor

var file_trading_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x74, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe8, 0x02, 0x0a, 0x04, 0x44, 0x65, 0x61,
	0x6c, 0x12, 0x16, 0x0a, 0x06, 0x44, 0x65, 0x61, 0x6c, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x44, 0x65, 0x61, 0x6c, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x68, 0x61,
	0x72, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0b,
	0x53, 0x68, 0x61, 0x72, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x50,
	0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6d,
	0x70, 0x61, 0x6e, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x6d, 0x70,
	0x61, 0x6e, 0x79, 0x12, 0x24, 0x0a, 0x0d, 0x50, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0d, 0x50, 0x75, 0x72, 0x63,
	0x68, 0x61, 0x73, 0x65, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x74, 0x6f,
	0x70, 0x4c, 0x6f, 0x73, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x53, 0x74, 0x6f,
	0x70, 0x4c, 0x6f, 0x73, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x61, 0x6b, 0x65, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x54, 0x61, 0x6b, 0x65, 0x50,
	0x72, 0x6f, 0x66, 0x69, 0x74, 0x12, 0x36, 0x0a, 0x08, 0x44, 0x65, 0x61, 0x6c, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x08, 0x44, 0x65, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3c, 0x0a,
	0x0b, 0x45, 0x6e, 0x64, 0x44, 0x65, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b,
	0x45, 0x6e, 0x64, 0x44, 0x65, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x50,
	0x72, 0x6f, 0x66, 0x69, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x74, 0x22, 0x4a, 0x0a, 0x11, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x69, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x65, 0x67, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x65, 0x67, 0x79, 0x12, 0x19, 0x0a, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x05, 0x2e, 0x44, 0x65, 0x61, 0x6c, 0x52, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x22,
	0x2c, 0x0a, 0x12, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x32, 0x47, 0x0a,
	0x0e, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x35, 0x0a, 0x0a, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x69, 0x65, 0x73, 0x12, 0x12, 0x2e,
	0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x13, 0x2e, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x69, 0x65, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x74, 0x6e, 0x69, 0x6b, 0x65, 0x6c, 0x2f, 0x54, 0x72,
	0x61, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_trading_service_proto_rawDescOnce sync.Once
	file_trading_service_proto_rawDescData = file_trading_service_proto_rawDesc
)

func file_trading_service_proto_rawDescGZIP() []byte {
	file_trading_service_proto_rawDescOnce.Do(func() {
		file_trading_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_trading_service_proto_rawDescData)
	})
	return file_trading_service_proto_rawDescData
}

var file_trading_service_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_trading_service_proto_goTypes = []interface{}{
	(*Deal)(nil),                  // 0: Deal
	(*StrategiesRequest)(nil),     // 1: StrategiesRequest
	(*StrategiesResponse)(nil),    // 2: StrategiesResponse
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_trading_service_proto_depIdxs = []int32{
	3, // 0: Deal.DealTime:type_name -> google.protobuf.Timestamp
	3, // 1: Deal.EndDealTime:type_name -> google.protobuf.Timestamp
	0, // 2: StrategiesRequest.deal:type_name -> Deal
	1, // 3: TradingService.Strategies:input_type -> StrategiesRequest
	2, // 4: TradingService.Strategies:output_type -> StrategiesResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_trading_service_proto_init() }
func file_trading_service_proto_init() {
	if File_trading_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_trading_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Deal); i {
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
		file_trading_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StrategiesRequest); i {
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
		file_trading_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StrategiesResponse); i {
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
			RawDescriptor: file_trading_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_trading_service_proto_goTypes,
		DependencyIndexes: file_trading_service_proto_depIdxs,
		MessageInfos:      file_trading_service_proto_msgTypes,
	}.Build()
	File_trading_service_proto = out.File
	file_trading_service_proto_rawDesc = nil
	file_trading_service_proto_goTypes = nil
	file_trading_service_proto_depIdxs = nil
}
