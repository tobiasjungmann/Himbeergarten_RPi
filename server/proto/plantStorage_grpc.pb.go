// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: plantStorage.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PlantStorageClient is the client API for PlantStorage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlantStorageClient interface {
	GetOverviewAllPlants(ctx context.Context, in *GetAllPlantsRequest, opts ...grpc.CallOption) (*AllPlantsReply, error)
	GetAdditionalDataPlant(ctx context.Context, in *GetAdditionalDataPlantRequest, opts ...grpc.CallOption) (*GetAdditionalDataPlantReply, error)
	// Also used to update a plant with the same id if it already exists
	AddNewPlant(ctx context.Context, in *AddPlantRequest, opts ...grpc.CallOption) (*PlantOverviewMsg, error)
	DeletePlant(ctx context.Context, in *PlantRequest, opts ...grpc.CallOption) (*DeletePlantReply, error)
	// Get an overview of all Devices given by a mac address and the sensorlots which are avilable
	GetConnectedSensorOverview(ctx context.Context, in *GetSensorOverviewRequest, opts ...grpc.CallOption) (*GetSensorOverviewReply, error)
	// get the data for a sensor given by its sensorslot and the mac address of the connected device
	GetDataForSensor(ctx context.Context, in *GetDataForSensorRequest, opts ...grpc.CallOption) (*GetDataForSensorReply, error)
	// Get the list of all Sensor Ids which are available for a device given by its mac address
	GetSensorsForDevice(ctx context.Context, in *GetSensorsForDeviceRequest, opts ...grpc.CallOption) (*GetSensorsForDeviceReply, error)
	// Set the list of sensors which are available at a device which should measure their values
	SetActiveSensorsForDevice(ctx context.Context, in *SetActiveSensorsForDeviceRequest, opts ...grpc.CallOption) (*SetActiveSensorsForDeviceReply, error)
}

type plantStorageClient struct {
	cc grpc.ClientConnInterface
}

func NewPlantStorageClient(cc grpc.ClientConnInterface) PlantStorageClient {
	return &plantStorageClient{cc}
}

func (c *plantStorageClient) GetOverviewAllPlants(ctx context.Context, in *GetAllPlantsRequest, opts ...grpc.CallOption) (*AllPlantsReply, error) {
	out := new(AllPlantsReply)
	err := c.cc.Invoke(ctx, "/smart_home.PlantStorage/getOverviewAllPlants", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plantStorageClient) GetAdditionalDataPlant(ctx context.Context, in *GetAdditionalDataPlantRequest, opts ...grpc.CallOption) (*GetAdditionalDataPlantReply, error) {
	out := new(GetAdditionalDataPlantReply)
	err := c.cc.Invoke(ctx, "/smart_home.PlantStorage/getAdditionalDataPlant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plantStorageClient) AddNewPlant(ctx context.Context, in *AddPlantRequest, opts ...grpc.CallOption) (*PlantOverviewMsg, error) {
	out := new(PlantOverviewMsg)
	err := c.cc.Invoke(ctx, "/smart_home.PlantStorage/addNewPlant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plantStorageClient) DeletePlant(ctx context.Context, in *PlantRequest, opts ...grpc.CallOption) (*DeletePlantReply, error) {
	out := new(DeletePlantReply)
	err := c.cc.Invoke(ctx, "/smart_home.PlantStorage/deletePlant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plantStorageClient) GetConnectedSensorOverview(ctx context.Context, in *GetSensorOverviewRequest, opts ...grpc.CallOption) (*GetSensorOverviewReply, error) {
	out := new(GetSensorOverviewReply)
	err := c.cc.Invoke(ctx, "/smart_home.PlantStorage/getConnectedSensorOverview", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plantStorageClient) GetDataForSensor(ctx context.Context, in *GetDataForSensorRequest, opts ...grpc.CallOption) (*GetDataForSensorReply, error) {
	out := new(GetDataForSensorReply)
	err := c.cc.Invoke(ctx, "/smart_home.PlantStorage/GetDataForSensor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plantStorageClient) GetSensorsForDevice(ctx context.Context, in *GetSensorsForDeviceRequest, opts ...grpc.CallOption) (*GetSensorsForDeviceReply, error) {
	out := new(GetSensorsForDeviceReply)
	err := c.cc.Invoke(ctx, "/smart_home.PlantStorage/GetSensorsForDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plantStorageClient) SetActiveSensorsForDevice(ctx context.Context, in *SetActiveSensorsForDeviceRequest, opts ...grpc.CallOption) (*SetActiveSensorsForDeviceReply, error) {
	out := new(SetActiveSensorsForDeviceReply)
	err := c.cc.Invoke(ctx, "/smart_home.PlantStorage/SetActiveSensorsForDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlantStorageServer is the server API for PlantStorage service.
// All implementations must embed UnimplementedPlantStorageServer
// for forward compatibility
type PlantStorageServer interface {
	GetOverviewAllPlants(context.Context, *GetAllPlantsRequest) (*AllPlantsReply, error)
	GetAdditionalDataPlant(context.Context, *GetAdditionalDataPlantRequest) (*GetAdditionalDataPlantReply, error)
	// Also used to update a plant with the same id if it already exists
	AddNewPlant(context.Context, *AddPlantRequest) (*PlantOverviewMsg, error)
	DeletePlant(context.Context, *PlantRequest) (*DeletePlantReply, error)
	// Get an overview of all Devices given by a mac address and the sensorlots which are avilable
	GetConnectedSensorOverview(context.Context, *GetSensorOverviewRequest) (*GetSensorOverviewReply, error)
	// get the data for a sensor given by its sensorslot and the mac address of the connected device
	GetDataForSensor(context.Context, *GetDataForSensorRequest) (*GetDataForSensorReply, error)
	// Get the list of all Sensor Ids which are available for a device given by its mac address
	GetSensorsForDevice(context.Context, *GetSensorsForDeviceRequest) (*GetSensorsForDeviceReply, error)
	// Set the list of sensors which are available at a device which should measure their values
	SetActiveSensorsForDevice(context.Context, *SetActiveSensorsForDeviceRequest) (*SetActiveSensorsForDeviceReply, error)
	mustEmbedUnimplementedPlantStorageServer()
}

// UnimplementedPlantStorageServer must be embedded to have forward compatible implementations.
type UnimplementedPlantStorageServer struct {
}

func (UnimplementedPlantStorageServer) GetOverviewAllPlants(context.Context, *GetAllPlantsRequest) (*AllPlantsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOverviewAllPlants not implemented")
}
func (UnimplementedPlantStorageServer) GetAdditionalDataPlant(context.Context, *GetAdditionalDataPlantRequest) (*GetAdditionalDataPlantReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAdditionalDataPlant not implemented")
}
func (UnimplementedPlantStorageServer) AddNewPlant(context.Context, *AddPlantRequest) (*PlantOverviewMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddNewPlant not implemented")
}
func (UnimplementedPlantStorageServer) DeletePlant(context.Context, *PlantRequest) (*DeletePlantReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePlant not implemented")
}
func (UnimplementedPlantStorageServer) GetConnectedSensorOverview(context.Context, *GetSensorOverviewRequest) (*GetSensorOverviewReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConnectedSensorOverview not implemented")
}
func (UnimplementedPlantStorageServer) GetDataForSensor(context.Context, *GetDataForSensorRequest) (*GetDataForSensorReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDataForSensor not implemented")
}
func (UnimplementedPlantStorageServer) GetSensorsForDevice(context.Context, *GetSensorsForDeviceRequest) (*GetSensorsForDeviceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSensorsForDevice not implemented")
}
func (UnimplementedPlantStorageServer) SetActiveSensorsForDevice(context.Context, *SetActiveSensorsForDeviceRequest) (*SetActiveSensorsForDeviceReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetActiveSensorsForDevice not implemented")
}
func (UnimplementedPlantStorageServer) mustEmbedUnimplementedPlantStorageServer() {}

// UnsafePlantStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlantStorageServer will
// result in compilation errors.
type UnsafePlantStorageServer interface {
	mustEmbedUnimplementedPlantStorageServer()
}

func RegisterPlantStorageServer(s grpc.ServiceRegistrar, srv PlantStorageServer) {
	s.RegisterService(&PlantStorage_ServiceDesc, srv)
}

func _PlantStorage_GetOverviewAllPlants_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllPlantsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlantStorageServer).GetOverviewAllPlants(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smart_home.PlantStorage/getOverviewAllPlants",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlantStorageServer).GetOverviewAllPlants(ctx, req.(*GetAllPlantsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlantStorage_GetAdditionalDataPlant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAdditionalDataPlantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlantStorageServer).GetAdditionalDataPlant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smart_home.PlantStorage/getAdditionalDataPlant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlantStorageServer).GetAdditionalDataPlant(ctx, req.(*GetAdditionalDataPlantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlantStorage_AddNewPlant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPlantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlantStorageServer).AddNewPlant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smart_home.PlantStorage/addNewPlant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlantStorageServer).AddNewPlant(ctx, req.(*AddPlantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlantStorage_DeletePlant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlantStorageServer).DeletePlant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smart_home.PlantStorage/deletePlant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlantStorageServer).DeletePlant(ctx, req.(*PlantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlantStorage_GetConnectedSensorOverview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSensorOverviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlantStorageServer).GetConnectedSensorOverview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smart_home.PlantStorage/getConnectedSensorOverview",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlantStorageServer).GetConnectedSensorOverview(ctx, req.(*GetSensorOverviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlantStorage_GetDataForSensor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDataForSensorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlantStorageServer).GetDataForSensor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smart_home.PlantStorage/GetDataForSensor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlantStorageServer).GetDataForSensor(ctx, req.(*GetDataForSensorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlantStorage_GetSensorsForDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSensorsForDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlantStorageServer).GetSensorsForDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smart_home.PlantStorage/GetSensorsForDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlantStorageServer).GetSensorsForDevice(ctx, req.(*GetSensorsForDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlantStorage_SetActiveSensorsForDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetActiveSensorsForDeviceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlantStorageServer).SetActiveSensorsForDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smart_home.PlantStorage/SetActiveSensorsForDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlantStorageServer).SetActiveSensorsForDevice(ctx, req.(*SetActiveSensorsForDeviceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PlantStorage_ServiceDesc is the grpc.ServiceDesc for PlantStorage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PlantStorage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smart_home.PlantStorage",
	HandlerType: (*PlantStorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getOverviewAllPlants",
			Handler:    _PlantStorage_GetOverviewAllPlants_Handler,
		},
		{
			MethodName: "getAdditionalDataPlant",
			Handler:    _PlantStorage_GetAdditionalDataPlant_Handler,
		},
		{
			MethodName: "addNewPlant",
			Handler:    _PlantStorage_AddNewPlant_Handler,
		},
		{
			MethodName: "deletePlant",
			Handler:    _PlantStorage_DeletePlant_Handler,
		},
		{
			MethodName: "getConnectedSensorOverview",
			Handler:    _PlantStorage_GetConnectedSensorOverview_Handler,
		},
		{
			MethodName: "GetDataForSensor",
			Handler:    _PlantStorage_GetDataForSensor_Handler,
		},
		{
			MethodName: "GetSensorsForDevice",
			Handler:    _PlantStorage_GetSensorsForDevice_Handler,
		},
		{
			MethodName: "SetActiveSensorsForDevice",
			Handler:    _PlantStorage_SetActiveSensorsForDevice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plantStorage.proto",
}
