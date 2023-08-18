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

type TradingBalance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BalanceID string  `protobuf:"bytes,1,opt,name=BalanceID,proto3" json:"BalanceID,omitempty"`
	ProfileID string  `protobuf:"bytes,2,opt,name=ProfileID,proto3" json:"ProfileID,omitempty"`
	Operation float64 `protobuf:"fixed64,3,opt,name=Operation,proto3" json:"Operation,omitempty"`
}

func (x *TradingBalance) Reset() {
	*x = TradingBalance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TradingBalance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradingBalance) ProtoMessage() {}

func (x *TradingBalance) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TradingBalance.ProtoReflect.Descriptor instead.
func (*TradingBalance) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{1}
}

func (x *TradingBalance) GetBalanceID() string {
	if x != nil {
		return x.BalanceID
	}
	return ""
}

func (x *TradingBalance) GetProfileID() string {
	if x != nil {
		return x.ProfileID
	}
	return ""
}

func (x *TradingBalance) GetOperation() float64 {
	if x != nil {
		return x.Operation
	}
	return 0
}

type TradingShare struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Company string  `protobuf:"bytes,1,opt,name=Company,proto3" json:"Company,omitempty"`
	Price   float64 `protobuf:"fixed64,2,opt,name=Price,proto3" json:"Price,omitempty"`
}

func (x *TradingShare) Reset() {
	*x = TradingShare{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TradingShare) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradingShare) ProtoMessage() {}

func (x *TradingShare) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TradingShare.ProtoReflect.Descriptor instead.
func (*TradingShare) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{2}
}

func (x *TradingShare) GetCompany() string {
	if x != nil {
		return x.Company
	}
	return ""
}

func (x *TradingShare) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

type GetProfitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Strategy string `protobuf:"bytes,1,opt,name=strategy,proto3" json:"strategy,omitempty"`
	Deal     *Deal  `protobuf:"bytes,2,opt,name=deal,proto3" json:"deal,omitempty"`
}

func (x *GetProfitRequest) Reset() {
	*x = GetProfitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetProfitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProfitRequest) ProtoMessage() {}

func (x *GetProfitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trading_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProfitRequest.ProtoReflect.Descriptor instead.
func (*GetProfitRequest) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetProfitRequest) GetStrategy() string {
	if x != nil {
		return x.Strategy
	}
	return ""
}

func (x *GetProfitRequest) GetDeal() *Deal {
	if x != nil {
		return x.Deal
	}
	return nil
}

type GetProfitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Profit float64 `protobuf:"fixed64,1,opt,name=Profit,proto3" json:"Profit,omitempty"`
}

func (x *GetProfitResponse) Reset() {
	*x = GetProfitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetProfitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProfitResponse) ProtoMessage() {}

func (x *GetProfitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trading_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProfitResponse.ProtoReflect.Descriptor instead.
func (*GetProfitResponse) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetProfitResponse) GetProfit() float64 {
	if x != nil {
		return x.Profit
	}
	return 0
}

type ClosePositionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Strategy string          `protobuf:"bytes,1,opt,name=strategy,proto3" json:"strategy,omitempty"`
	Share    *TradingShare   `protobuf:"bytes,2,opt,name=share,proto3" json:"share,omitempty"`
	Deal     *Deal           `protobuf:"bytes,3,opt,name=deal,proto3" json:"deal,omitempty"`
	Balance  *TradingBalance `protobuf:"bytes,4,opt,name=balance,proto3" json:"balance,omitempty"`
}

func (x *ClosePositionRequest) Reset() {
	*x = ClosePositionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClosePositionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClosePositionRequest) ProtoMessage() {}

func (x *ClosePositionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trading_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClosePositionRequest.ProtoReflect.Descriptor instead.
func (*ClosePositionRequest) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{5}
}

func (x *ClosePositionRequest) GetStrategy() string {
	if x != nil {
		return x.Strategy
	}
	return ""
}

func (x *ClosePositionRequest) GetShare() *TradingShare {
	if x != nil {
		return x.Share
	}
	return nil
}

func (x *ClosePositionRequest) GetDeal() *Deal {
	if x != nil {
		return x.Deal
	}
	return nil
}

func (x *ClosePositionRequest) GetBalance() *TradingBalance {
	if x != nil {
		return x.Balance
	}
	return nil
}

type ClosePositionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ClosePositionResponse) Reset() {
	*x = ClosePositionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClosePositionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClosePositionResponse) ProtoMessage() {}

func (x *ClosePositionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trading_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClosePositionResponse.ProtoReflect.Descriptor instead.
func (*ClosePositionResponse) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{6}
}

type GetUnclosedPositionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Profileid string `protobuf:"bytes,1,opt,name=profileid,proto3" json:"profileid,omitempty"`
}

func (x *GetUnclosedPositionsRequest) Reset() {
	*x = GetUnclosedPositionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUnclosedPositionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUnclosedPositionsRequest) ProtoMessage() {}

func (x *GetUnclosedPositionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trading_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUnclosedPositionsRequest.ProtoReflect.Descriptor instead.
func (*GetUnclosedPositionsRequest) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{7}
}

func (x *GetUnclosedPositionsRequest) GetProfileid() string {
	if x != nil {
		return x.Profileid
	}
	return ""
}

type GetUnclosedPositionsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Deal []*Deal `protobuf:"bytes,1,rep,name=deal,proto3" json:"deal,omitempty"`
}

func (x *GetUnclosedPositionsResponse) Reset() {
	*x = GetUnclosedPositionsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_trading_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUnclosedPositionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUnclosedPositionsResponse) ProtoMessage() {}

func (x *GetUnclosedPositionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trading_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUnclosedPositionsResponse.ProtoReflect.Descriptor instead.
func (*GetUnclosedPositionsResponse) Descriptor() ([]byte, []int) {
	return file_trading_service_proto_rawDescGZIP(), []int{8}
}

func (x *GetUnclosedPositionsResponse) GetDeal() []*Deal {
	if x != nil {
		return x.Deal
	}
	return nil
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
	0x66, 0x69, 0x74, 0x22, 0x6a, 0x0a, 0x0e, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x42, 0x61,
	0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63,
	0x65, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x44,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x49,
	0x44, 0x12, 0x1c, 0x0a, 0x09, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x3e, 0x0a, 0x0c, 0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x68, 0x61, 0x72, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x43, 0x6f, 0x6d, 0x70, 0x61, 0x6e, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x69,
	0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x22,
	0x49, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12,
	0x19, 0x0a, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e,
	0x44, 0x65, 0x61, 0x6c, 0x52, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x22, 0x2b, 0x0a, 0x11, 0x47, 0x65,
	0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x06, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x22, 0x9d, 0x01, 0x0a, 0x14, 0x43, 0x6c, 0x6f, 0x73,
	0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x23, 0x0a, 0x05,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x54, 0x72,
	0x61, 0x64, 0x69, 0x6e, 0x67, 0x53, 0x68, 0x61, 0x72, 0x65, 0x52, 0x05, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x12, 0x19, 0x0a, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x05, 0x2e, 0x44, 0x65, 0x61, 0x6c, 0x52, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x12, 0x29, 0x0a, 0x07,
	0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x54, 0x72, 0x61, 0x64, 0x69, 0x6e, 0x67, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x07,
	0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x17, 0x0a, 0x15, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x3b, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x69, 0x64, 0x22, 0x39, 0x0a,
	0x1c, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a,
	0x04, 0x64, 0x65, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x44, 0x65,
	0x61, 0x6c, 0x52, 0x04, 0x64, 0x65, 0x61, 0x6c, 0x32, 0xd9, 0x01, 0x0a, 0x0e, 0x54, 0x72, 0x61,
	0x64, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x09, 0x47,
	0x65, 0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x12, 0x11, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x47, 0x65,
	0x74, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3e, 0x0a, 0x0d, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x15, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x50,
	0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x53, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x50, 0x6f,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1c, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x63,
	0x6c, 0x6f, 0x73, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x6e, 0x63, 0x6c, 0x6f,
	0x73, 0x65, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x74, 0x6e, 0x69, 0x6b, 0x65, 0x6c, 0x2f, 0x54, 0x72, 0x61, 0x64,
	0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_trading_service_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_trading_service_proto_goTypes = []interface{}{
	(*Deal)(nil),                         // 0: Deal
	(*TradingBalance)(nil),               // 1: TradingBalance
	(*TradingShare)(nil),                 // 2: TradingShare
	(*GetProfitRequest)(nil),             // 3: GetProfitRequest
	(*GetProfitResponse)(nil),            // 4: GetProfitResponse
	(*ClosePositionRequest)(nil),         // 5: ClosePositionRequest
	(*ClosePositionResponse)(nil),        // 6: ClosePositionResponse
	(*GetUnclosedPositionsRequest)(nil),  // 7: GetUnclosedPositionsRequest
	(*GetUnclosedPositionsResponse)(nil), // 8: GetUnclosedPositionsResponse
	(*timestamppb.Timestamp)(nil),        // 9: google.protobuf.Timestamp
}
var file_trading_service_proto_depIdxs = []int32{
	9,  // 0: Deal.DealTime:type_name -> google.protobuf.Timestamp
	9,  // 1: Deal.EndDealTime:type_name -> google.protobuf.Timestamp
	0,  // 2: GetProfitRequest.deal:type_name -> Deal
	2,  // 3: ClosePositionRequest.share:type_name -> TradingShare
	0,  // 4: ClosePositionRequest.deal:type_name -> Deal
	1,  // 5: ClosePositionRequest.balance:type_name -> TradingBalance
	0,  // 6: GetUnclosedPositionsResponse.deal:type_name -> Deal
	3,  // 7: TradingService.GetProfit:input_type -> GetProfitRequest
	5,  // 8: TradingService.ClosePosition:input_type -> ClosePositionRequest
	7,  // 9: TradingService.GetUnclosedPositions:input_type -> GetUnclosedPositionsRequest
	4,  // 10: TradingService.GetProfit:output_type -> GetProfitResponse
	6,  // 11: TradingService.ClosePosition:output_type -> ClosePositionResponse
	8,  // 12: TradingService.GetUnclosedPositions:output_type -> GetUnclosedPositionsResponse
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
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
			switch v := v.(*TradingBalance); i {
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
			switch v := v.(*TradingShare); i {
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
		file_trading_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetProfitRequest); i {
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
		file_trading_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetProfitResponse); i {
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
		file_trading_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClosePositionRequest); i {
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
		file_trading_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClosePositionResponse); i {
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
		file_trading_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUnclosedPositionsRequest); i {
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
		file_trading_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUnclosedPositionsResponse); i {
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
			NumMessages:   9,
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
