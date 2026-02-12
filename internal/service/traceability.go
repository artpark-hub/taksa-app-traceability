package service

import (
	"context"
	"time"

	pb "traceability/api/traceability/v1"
	"traceability/internal/biz"
)

type TraceabilityService struct {
	pb.UnimplementedTraceabilityServer
	uc *biz.TraceabilityUsecase
}

func NewTraceabilityService(uc *biz.TraceabilityUsecase) *TraceabilityService {
	return &TraceabilityService{uc: uc}
}

// ==========================================
// 1. Enterprise Management
// ==========================================

func (s *TraceabilityService) CreateEnterprise(ctx context.Context, req *pb.CreateEnterpriseRequest) (*pb.CreateEnterpriseReply, error) {
	id, err := s.uc.CreateEnterprise(ctx, &biz.Enterprise{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateEnterpriseReply{Id: id}, nil
}

func (s *TraceabilityService) ListEnterprises(ctx context.Context, req *pb.ListEnterprisesRequest) (*pb.ListEnterprisesReply, error) {
	list, err := s.uc.ListEnterprises(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Enterprise, 0)
	for _, x := range list {
		res = append(res, &pb.Enterprise{Id: x.ID, Name: x.Name, Description: x.Description})
	}
	return &pb.ListEnterprisesReply{Enterprises: res}, nil
}

func (s *TraceabilityService) UpdateEnterprise(ctx context.Context, req *pb.UpdateEnterpriseRequest) (*pb.UpdateEnterpriseReply, error) {
	err := s.uc.UpdateEnterprise(ctx, &biz.Enterprise{ID: req.Id, Name: req.Name, Description: req.Description})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateEnterpriseReply{Success: true}, nil
}

func (s *TraceabilityService) DeleteEnterprise(ctx context.Context, req *pb.DeleteEnterpriseRequest) (*pb.DeleteEnterpriseReply, error) {
	err := s.uc.DeleteEnterprise(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteEnterpriseReply{Success: true}, nil
}

// ==========================================
// 2. Site Management
// ==========================================

func (s *TraceabilityService) CreateSite(ctx context.Context, req *pb.CreateSiteRequest) (*pb.CreateSiteReply, error) {
	id, err := s.uc.CreateSite(ctx, &biz.Site{
		EnterpriseID: req.EnterpriseId,
		Name:         req.Name,
		Location:     req.Location,
		Description:  req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateSiteReply{Id: id}, nil
}

func (s *TraceabilityService) ListSites(ctx context.Context, req *pb.ListSitesRequest) (*pb.ListSitesReply, error) {
	list, err := s.uc.ListSites(ctx, req.EnterpriseId)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Site, 0)
	for _, x := range list {
		res = append(res, &pb.Site{Id: x.ID, Name: x.Name, Location: x.Location, Description: x.Description})
	}
	return &pb.ListSitesReply{Sites: res}, nil
}

func (s *TraceabilityService) UpdateSite(ctx context.Context, req *pb.UpdateSiteRequest) (*pb.UpdateSiteReply, error) {
	err := s.uc.UpdateSite(ctx, &biz.Site{ID: req.Id, Location: req.Location})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateSiteReply{Success: true}, nil
}

func (s *TraceabilityService) DeleteSite(ctx context.Context, req *pb.DeleteSiteRequest) (*pb.DeleteSiteReply, error) {
	err := s.uc.DeleteSite(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteSiteReply{Success: true}, nil
}

// ==========================================
// 3. Area Management
// ==========================================

func (s *TraceabilityService) CreateArea(ctx context.Context, req *pb.CreateAreaRequest) (*pb.CreateAreaReply, error) {
	id, err := s.uc.CreateArea(ctx, &biz.Area{
		SiteID:      req.SiteId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateAreaReply{Id: id}, nil
}

func (s *TraceabilityService) ListAreas(ctx context.Context, req *pb.ListAreasRequest) (*pb.ListAreasReply, error) {
	list, err := s.uc.ListAreas(ctx, req.SiteId)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Area, 0)
	for _, x := range list {
		res = append(res, &pb.Area{Id: x.ID, Name: x.Name, Description: x.Description})
	}
	return &pb.ListAreasReply{Areas: res}, nil
}

func (s *TraceabilityService) UpdateArea(ctx context.Context, req *pb.UpdateAreaRequest) (*pb.UpdateAreaReply, error) {
	err := s.uc.UpdateArea(ctx, &biz.Area{ID: req.Id, Name: req.Name, Description: req.Description})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateAreaReply{Success: true}, nil
}

func (s *TraceabilityService) DeleteArea(ctx context.Context, req *pb.DeleteAreaRequest) (*pb.DeleteAreaReply, error) {
	err := s.uc.DeleteArea(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteAreaReply{Success: true}, nil
}

// ==========================================
// 4. Production Line Management
// ==========================================

func (s *TraceabilityService) CreateLine(ctx context.Context, req *pb.CreateLineRequest) (*pb.CreateLineReply, error) {
	id, err := s.uc.CreateLine(ctx, &biz.ProductionLine{
		AreaID:      req.AreaId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateLineReply{Id: id}, nil
}

func (s *TraceabilityService) ListLines(ctx context.Context, req *pb.ListLinesRequest) (*pb.ListLinesReply, error) {
	list, err := s.uc.ListLines(ctx, req.AreaId)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Line, 0)
	for _, x := range list {
		res = append(res, &pb.Line{Id: x.ID, Name: x.Name, Description: x.Description})
	}
	return &pb.ListLinesReply{Lines: res}, nil
}

func (s *TraceabilityService) UpdateLine(ctx context.Context, req *pb.UpdateLineRequest) (*pb.UpdateLineReply, error) {
	err := s.uc.UpdateLine(ctx, &biz.ProductionLine{ID: req.Id, Name: req.Name, Description: req.Description})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateLineReply{Success: true}, nil
}

func (s *TraceabilityService) DeleteLine(ctx context.Context, req *pb.DeleteLineRequest) (*pb.DeleteLineReply, error) {
	err := s.uc.DeleteLine(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteLineReply{Success: true}, nil
}

// ==========================================
// 5. Production Unit Management
// ==========================================

func (s *TraceabilityService) CreateProductionUnit(ctx context.Context, req *pb.CreateProductionUnitRequest) (*pb.CreateProductionUnitReply, error) {
	id, err := s.uc.CreateProductionUnit(ctx, &biz.ProductionUnit{
		ProductionLineID: req.ProductionLineId,
		Name:             req.Name,
		Description:      req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateProductionUnitReply{Id: id}, nil
}

func (s *TraceabilityService) ListProductionUnits(ctx context.Context, req *pb.ListProductionUnitsRequest) (*pb.ListProductionUnitsReply, error) {
	list, err := s.uc.ListProductionUnits(ctx, req.LineId)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.ProductionUnit, 0)
	for _, x := range list {
		res = append(res, &pb.ProductionUnit{Id: x.ID, Name: x.Name, Description: x.Description})
	}
	return &pb.ListProductionUnitsReply{Units: res}, nil
}

func (s *TraceabilityService) UpdateProductionUnit(ctx context.Context, req *pb.UpdateProductionUnitRequest) (*pb.UpdateProductionUnitReply, error) {
	err := s.uc.UpdateProductionUnit(ctx, &biz.ProductionUnit{ID: req.Id, Name: req.Name})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateProductionUnitReply{Success: true}, nil
}

func (s *TraceabilityService) DeleteProductionUnit(ctx context.Context, req *pb.DeleteProductionUnitRequest) (*pb.DeleteProductionUnitReply, error) {
	err := s.uc.DeleteProductionUnit(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductionUnitReply{Success: true}, nil
}

// ==========================================
// 6. Equipment Class Management
// ==========================================

func (s *TraceabilityService) CreateEquipmentClass(ctx context.Context, req *pb.CreateEquipmentClassRequest) (*pb.CreateEquipmentClassReply, error) {
	id, err := s.uc.CreateEquipmentClass(ctx, &biz.EquipmentClass{
		ClassName:   req.ClassName,
		Version:     req.Version,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateEquipmentClassReply{Id: id}, nil
}

func (s *TraceabilityService) ListEquipmentClasses(ctx context.Context, req *pb.ListEquipmentClassesRequest) (*pb.ListEquipmentClassesReply, error) {
	list, err := s.uc.ListEquipmentClasses(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.EquipmentClass, 0)
	for _, x := range list {
		res = append(res, &pb.EquipmentClass{Id: x.ID, ClassName: x.ClassName, Version: x.Version, Description: x.Description})
	}
	return &pb.ListEquipmentClassesReply{Classes: res}, nil
}

func (s *TraceabilityService) UpdateEquipmentClass(ctx context.Context, req *pb.UpdateEquipmentClassRequest) (*pb.UpdateEquipmentClassReply, error) {
	err := s.uc.UpdateEquipmentClass(ctx, &biz.EquipmentClass{ID: req.Id, Version: req.Version})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateEquipmentClassReply{Success: true}, nil
}

func (s *TraceabilityService) DeleteEquipmentClass(ctx context.Context, req *pb.DeleteEquipmentClassRequest) (*pb.DeleteEquipmentClassReply, error) {
	err := s.uc.DeleteEquipmentClass(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteEquipmentClassReply{Success: true}, nil
}

// ==========================================
// 7. Equipment Master Management
// ==========================================

func (s *TraceabilityService) RegisterEquipment(ctx context.Context, req *pb.RegisterEquipmentRequest) (*pb.RegisterEquipmentReply, error) {
	id, err := s.uc.RegisterEquipment(ctx, &biz.EquipmentMaster{
		ID:                req.Id,
		PhysicalAssetID:   req.PhysicalAssetId,
		ProductionUnitID:  req.ProductionUnitId,
		EquipmentClassID:  req.EquipmentClassId,
		OperationalStatus: req.OperationalStatus,
		ParentEquipmentID: req.ParentEquipmentId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.RegisterEquipmentReply{Id: id}, nil
}

func (s *TraceabilityService) ListEquipment(ctx context.Context, req *pb.ListEquipmentRequest) (*pb.ListEquipmentReply, error) {
	list, err := s.uc.ListEquipment(ctx, req.LineId, req.ParentEquipmentId)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Equipment, 0)
	for _, x := range list {
		res = append(res, &pb.Equipment{Id: x.ID, OperationalStatus: x.OperationalStatus})
	}
	return &pb.ListEquipmentReply{Equipment: res}, nil
}

func (s *TraceabilityService) UpdateEquipment(ctx context.Context, req *pb.UpdateEquipmentRequest) (*pb.UpdateEquipmentReply, error) {
	err := s.uc.UpdateEquipment(ctx, &biz.EquipmentMaster{
		ID:                req.Id,
		OperationalStatus: req.OperationalStatus,
		ProductionUnitID:  req.ProductionUnitId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateEquipmentReply{Success: true}, nil
}

func (s *TraceabilityService) DeleteEquipment(ctx context.Context, req *pb.DeleteEquipmentRequest) (*pb.DeleteEquipmentReply, error) {
	err := s.uc.DeleteEquipment(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteEquipmentReply{Success: true}, nil
}

// ==========================================
// 8. Equipment Capability
// ==========================================

func (s *TraceabilityService) AddCapability(ctx context.Context, req *pb.AddCapabilityRequest) (*pb.AddCapabilityReply, error) {
	id, err := s.uc.AddCapability(ctx, &biz.EquipmentCapability{
		EquipmentID:    req.EquipmentId,
		CapabilityName: req.CapabilityName,
		Value:          req.Value,
		UOM:            req.Uom,
		Description:    req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AddCapabilityReply{Id: id}, nil
}

func (s *TraceabilityService) ListCapabilities(ctx context.Context, req *pb.ListCapabilitiesRequest) (*pb.ListCapabilitiesReply, error) {
	list, err := s.uc.ListCapabilities(ctx, req.EquipmentId)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Capability, 0)
	for _, x := range list {
		res = append(res, &pb.Capability{Id: x.ID, Name: x.CapabilityName, Value: x.Value, Uom: x.UOM, Description: x.Description})
	}
	return &pb.ListCapabilitiesReply{Capabilities: res}, nil
}

func (s *TraceabilityService) UpdateCapability(ctx context.Context, req *pb.UpdateCapabilityRequest) (*pb.UpdateCapabilityReply, error) {
	err := s.uc.UpdateCapability(ctx, &biz.EquipmentCapability{ID: req.Id, Value: req.Value})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateCapabilityReply{Success: true}, nil
}

func (s *TraceabilityService) DeleteCapability(ctx context.Context, req *pb.DeleteCapabilityRequest) (*pb.DeleteCapabilityReply, error) {
	err := s.uc.DeleteCapability(ctx, req.EquipmentId, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteCapabilityReply{Success: true}, nil
}

// ==========================================
// 9. Equipment Properties
// ==========================================

func (s *TraceabilityService) SetProperty(ctx context.Context, req *pb.SetPropertyRequest) (*pb.SetPropertyReply, error) {
	id, err := s.uc.SetProperty(ctx, &biz.EquipmentProperty{
		EquipmentID:  req.EquipmentId,
		PropertyName: req.PropertyName,
		CurrentValue: req.CurrentValue,
	})
	if err != nil {
		return nil, err
	}
	return &pb.SetPropertyReply{Id: id}, nil
}

func (s *TraceabilityService) ListProperties(ctx context.Context, req *pb.ListPropertiesRequest) (*pb.ListPropertiesReply, error) {
	list, err := s.uc.ListProperties(ctx, req.EquipmentId)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Property, 0)
	for _, x := range list {
		res = append(res, &pb.Property{Id: x.ID, Name: x.PropertyName, Value: x.CurrentValue})
	}
	return &pb.ListPropertiesReply{Properties: res}, nil
}

func (s *TraceabilityService) UpdateProperty(ctx context.Context, req *pb.UpdatePropertyRequest) (*pb.UpdatePropertyReply, error) {
	err := s.uc.UpdateProperty(ctx, &biz.EquipmentProperty{ID: req.Id, CurrentValue: req.CurrentValue})
	if err != nil {
		return nil, err
	}
	return &pb.UpdatePropertyReply{Success: true}, nil
}

func (s *TraceabilityService) DeleteProperty(ctx context.Context, req *pb.DeletePropertyRequest) (*pb.DeletePropertyReply, error) {
	err := s.uc.DeleteProperty(ctx, req.EquipmentId, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeletePropertyReply{Success: true}, nil
}

// ==========================================
// 10. Traceability Logs
// ==========================================

func (s *TraceabilityService) LogEvent(ctx context.Context, req *pb.LogEventRequest) (*pb.LogEventReply, error) {
	id, err := s.uc.LogEvent(ctx, &biz.TraceabilityLog{
		EquipmentID:   req.EquipmentId,
		EventType:     req.EventType,
		WorkOrderID:   req.WorkOrderId,
		MaterialLotID: req.MaterialLotId,
		OperatorID:    req.OperatorId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.LogEventReply{Id: id}, nil
}

func (s *TraceabilityService) ListLogs(ctx context.Context, req *pb.ListLogsRequest) (*pb.ListLogsReply, error) {
	list, err := s.uc.ListLogs(ctx, req.WorkOrderId)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.LogEntry, 0)
	for _, x := range list {
		res = append(res, &pb.LogEntry{
			EquipmentId: x.EquipmentID,
			EventType:   x.EventType,
			EventTime:   x.EventTime.Format(time.RFC3339),
			OperatorId:  x.OperatorID,
		})
	}
	return &pb.ListLogsReply{Logs: res}, nil
}
