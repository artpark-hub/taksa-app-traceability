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

// 1. Enterprise Management

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

// 2. Site Management

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

// 3. Area Management

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

// 4. Production Line Management

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

// 5. Production Unit Management

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
    err := s.uc.UpdateProductionUnit(ctx, &biz.ProductionUnit{
        ID:          req.Id,
        Name:        req.Name,
        Description: req.Description,
    })
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

// 6. Equipment Class Management

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
	err := s.uc.UpdateEquipmentClass(ctx, &biz.EquipmentClass{ID: req.Id, Version: req.Version, Description: req.Description})
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

// 7. Equipment Master Management

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
    res = append(res, &pb.Equipment{
        Id:                x.ID,
        PhysicalAssetId:   x.PhysicalAssetID,
        ProductionUnitId:  x.ProductionUnitID,
        EquipmentClassId:  x.EquipmentClassID,
        OperationalStatus: x.OperationalStatus,
        ParentEquipmentId: x.ParentEquipmentID,
    })
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

// 8. Equipment Capability

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
	err := s.uc.UpdateCapability(ctx, &biz.EquipmentCapability{ID: req.Id, EquipmentID: req.EquipmentId, Value: req.Value})
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

// 9. Equipment Properties

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
	err := s.uc.UpdateProperty(ctx, &biz.EquipmentProperty{ID: req.Id, EquipmentID: req.EquipmentId, CurrentValue: req.CurrentValue})
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

// 10. Traceability Logs

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

// 11. Material Definition

func (s *TraceabilityService) CreateMaterialDefinition(ctx context.Context, req *pb.CreateMaterialDefinitionRequest) (*pb.CreateMaterialDefinitionReply, error) {
	id, err := s.uc.CreateMaterialDefinition(ctx, &biz.MaterialDefinition{
		Name:          req.Name,
		MaterialType:  req.MaterialType,
		UnitOfMeasure: req.UnitOfMeasure,
		Description:   req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateMaterialDefinitionReply{Id: id}, nil
}

func (s *TraceabilityService) ListMaterialDefinitions(ctx context.Context, req *pb.ListMaterialDefinitionsRequest) (*pb.ListMaterialDefinitionsReply, error) {
	list, err := s.uc.ListMaterialDefinitions(ctx, req.MaterialType)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.MaterialDefinition, 0)
	for _, x := range list {
		res = append(res, &pb.MaterialDefinition{
			Id:            x.ID,
			Name:          x.Name,
			MaterialType:  x.MaterialType,
			UnitOfMeasure: x.UnitOfMeasure,
			Description:   x.Description,
		})
	}
	return &pb.ListMaterialDefinitionsReply{Definitions: res}, nil
}

// 12. Material Lot

func (s *TraceabilityService) CreateMaterialLot(ctx context.Context, req *pb.CreateMaterialLotRequest) (*pb.CreateMaterialLotReply, error) {
	lotID, err := s.uc.CreateMaterialLot(ctx, &biz.MaterialLot{
		LotID:                req.LotId,
		MaterialDefinitionID: req.MaterialDefinitionId,
		Quantity:             req.Quantity,
		UnitOfMeasure:        req.UnitOfMeasure,
		Status:               "available",
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateMaterialLotReply{LotId: lotID}, nil
}

func (s *TraceabilityService) GetMaterialLot(ctx context.Context, req *pb.GetMaterialLotRequest) (*pb.GetMaterialLotReply, error) {
	det, err := s.uc.GetMaterialLot(ctx, req.LotId)
	if err != nil {
		return nil, err
	}
	return &pb.GetMaterialLotReply{
		LotId:                det.LotID,
		Status:               det.Status,
		Quantity:             det.Quantity,
		UnitOfMeasure:        det.UnitOfMeasure,
		CreatedAt:            det.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            det.UpdatedAt.Format(time.RFC3339),
		MaterialDefinitionId: det.MaterialDefinitionID,
		MaterialName:         det.MaterialName,
		MaterialType:         det.MaterialType,
		MaterialDescription:  det.MaterialDescription,
	}, nil
}

func (s *TraceabilityService) ListMaterialLots(ctx context.Context, req *pb.ListMaterialLotsRequest) (*pb.ListMaterialLotsReply, error) {
	list, err := s.uc.ListMaterialLots(ctx, req.Status, req.MaterialType)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.MaterialLotSummary, 0)
	for _, x := range list {
		res = append(res, &pb.MaterialLotSummary{
			LotId:         x.LotID,
			Status:        x.Status,
			Quantity:      x.Quantity,
			UnitOfMeasure: x.UnitOfMeasure,
			CreatedAt:     x.CreatedAt.Format(time.RFC3339),
			MaterialName:  x.MaterialName,
			MaterialType:  x.MaterialType,
		})
	}
	return &pb.ListMaterialLotsReply{Lots: res}, nil
}

func (s *TraceabilityService) UpdateMaterialLotStatus(ctx context.Context, req *pb.UpdateMaterialLotStatusRequest) (*pb.UpdateMaterialLotStatusReply, error) {
	err := s.uc.UpdateMaterialLotStatus(ctx, req.LotId, req.Status)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateMaterialLotStatusReply{Success: true}, nil
}

// 13. Operator

func (s *TraceabilityService) CreateOperator(ctx context.Context, req *pb.CreateOperatorRequest) (*pb.CreateOperatorReply, error) {
	opID, err := s.uc.CreateOperator(ctx, &biz.Operator{
		OperatorID: req.OperatorId,
		Name:       req.Name,
		Role:       req.Role,
		Shift:      req.Shift,
		Status:     "active",
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateOperatorReply{OperatorId: opID}, nil
}

func (s *TraceabilityService) ListOperators(ctx context.Context, req *pb.ListOperatorsRequest) (*pb.ListOperatorsReply, error) {
	list, err := s.uc.ListOperators(ctx, req.Shift, req.Status)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.Operator, 0)
	for _, x := range list {
		res = append(res, &pb.Operator{
			OperatorId: x.OperatorID,
			Name:       x.Name,
			Role:       x.Role,
			Shift:      x.Shift,
			Status:     x.Status,
		})
	}
	return &pb.ListOperatorsReply{Operators: res}, nil
}

func (s *TraceabilityService) UpdateOperator(ctx context.Context, req *pb.UpdateOperatorRequest) (*pb.UpdateOperatorReply, error) {
	err := s.uc.UpdateOperator(ctx, &biz.Operator{
		OperatorID: req.OperatorId,
		Role:       req.Role,
		Shift:      req.Shift,
		Status:     req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpdateOperatorReply{Success: true}, nil
}

// 14. Work Order

func (s *TraceabilityService) CreateWorkOrder(ctx context.Context, req *pb.CreateWorkOrderRequest) (*pb.CreateWorkOrderReply, error) {
	start, _ := time.Parse(time.RFC3339, req.PlannedStart)
	end, _ := time.Parse(time.RFC3339, req.PlannedEnd)
	woID, err := s.uc.CreateWorkOrder(ctx, &biz.WorkOrder{
		WorkOrderID:     req.WorkOrderId,
		Description:     req.Description,
		EquipmentID:     req.EquipmentId,
		OperatorID:      req.OperatorId,
		OutputLotID:     req.OutputLotId,
		PlannedQuantity: req.PlannedQuantity,
		UnitOfMeasure:   req.UnitOfMeasure,
		PlannedStart:    start,
		PlannedEnd:      end,
		Status:          "planned",
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateWorkOrderReply{WorkOrderId: woID}, nil
}

func (s *TraceabilityService) GetWorkOrder(ctx context.Context, req *pb.GetWorkOrderRequest) (*pb.GetWorkOrderReply, error) {
	x, err := s.uc.GetWorkOrder(ctx, req.WorkOrderId)
	if err != nil {
		return nil, err
	}
	
	inputs := make([]*pb.WorkOrderInputLot, len(x.InputLots))
	for i, in := range x.InputLots {
		inputs[i] = &pb.WorkOrderInputLot{
			LotId:            in.LotID,
			MaterialName:     in.MaterialName,
			MaterialType:     in.MaterialType,
			QuantityConsumed: in.QuantityConsumed,
			UnitOfMeasure:    in.UnitOfMeasure,
		}
	}
	
	res := &pb.GetWorkOrderReply{
		WorkOrderId:        x.WorkOrderID,
		Description:        x.Description,
		Status:             x.Status,
		PlannedQuantity:    x.PlannedQuantity,
		UnitOfMeasure:      x.UnitOfMeasure,
		PlannedStart:       x.PlannedStart.Format(time.RFC3339),
		PlannedEnd:         x.PlannedEnd.Format(time.RFC3339),
		OutputLotId:        x.OutputLotID,
		EquipmentId:        x.EquipmentID,
		EquipmentClassName: x.EquipmentClassName,
		OperatorId:         x.OperatorID,
		OperatorName:       x.OperatorName,
		OperatorShift:      x.OperatorShift,
		InputLots:          inputs,
	}
	if x.ActualQuantity != nil { res.ActualQuantity = *x.ActualQuantity }
	if x.ActualStart != nil { res.ActualStart = x.ActualStart.Format(time.RFC3339) }
	if x.ActualEnd != nil { res.ActualEnd = x.ActualEnd.Format(time.RFC3339) }
	
	return res, nil
}

func (s *TraceabilityService) ListWorkOrders(ctx context.Context, req *pb.ListWorkOrdersRequest) (*pb.ListWorkOrdersReply, error) {
	list, err := s.uc.ListWorkOrders(ctx, req.Status, req.EquipmentId)
	if err != nil {
		return nil, err
	}
	res := make([]*pb.WorkOrderSummary, 0)
	for _, x := range list {
		wo := &pb.WorkOrderSummary{
			WorkOrderId: x.WorkOrderID,
			Description: x.Description,
			Status:      x.Status,
			EquipmentId: x.EquipmentID,
			OperatorId:  x.OperatorID,
			OutputLotId: x.OutputLotID,
		}
		if x.ActualStart != nil { wo.ActualStart = x.ActualStart.Format(time.RFC3339) }
		if x.ActualEnd != nil { wo.ActualEnd = x.ActualEnd.Format(time.RFC3339) }
		res = append(res, wo)
	}
	return &pb.ListWorkOrdersReply{WorkOrders: res}, nil
}

func (s *TraceabilityService) UpdateWorkOrderStatus(ctx context.Context, req *pb.UpdateWorkOrderStatusRequest) (*pb.UpdateWorkOrderStatusReply, error) {
	wo := &biz.WorkOrder{
		WorkOrderID: req.WorkOrderId,
		Status:      req.Status,
	}
	if req.ActualQuantity > 0 { wo.ActualQuantity = &req.ActualQuantity }
	if req.ActualStart != "" { t, _ := time.Parse(time.RFC3339, req.ActualStart); wo.ActualStart = &t }
	if req.ActualEnd != "" { t, _ := time.Parse(time.RFC3339, req.ActualEnd); wo.ActualEnd = &t }
	
	err := s.uc.UpdateWorkOrderStatus(ctx, wo)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateWorkOrderStatusReply{Success: true}, nil
}

// 15. Genealogy

func (s *TraceabilityService) RegisterGenealogyLink(ctx context.Context, req *pb.RegisterGenealogyLinkRequest) (*pb.RegisterGenealogyLinkReply, error) {
	id, err := s.uc.RegisterGenealogyLink(ctx, &biz.LotGenealogy{
		ParentLotID:      req.ParentLotId,
		ChildLotID:       req.ChildLotId,
		WorkOrderID:      req.WorkOrderId,
		EquipmentID:      req.EquipmentId,
		QuantityConsumed: req.QuantityConsumed,
		QuantityProduced: req.QuantityProduced,
	})
	if err != nil {
		return nil, err
	}
	return &pb.RegisterGenealogyLinkReply{Id: id}, nil
}

// 16. Trace

func mapTraceNodes(bizNodes []*biz.TraceNode) []*pb.TraceNode {
	res := make([]*pb.TraceNode, len(bizNodes))
	for i, x := range bizNodes {
		res[i] = &pb.TraceNode{
			Depth:              x.Depth,
			LotId:              x.LotID,
			RelatedLotId:       x.RelatedLotID,
			MaterialName:       x.MaterialName,
			MaterialType:       x.MaterialType,
			LotStatus:          x.LotStatus,
			Quantity:           x.Quantity,
			UnitOfMeasure:      x.UnitOfMeasure,
			QuantityUsed:       x.QuantityUsed,
			WorkOrderId:        x.WorkOrderID,
			EquipmentId:        x.EquipmentID,
			EquipmentClassName: x.EquipmentClassName,
			OperatorId:         x.OperatorID,
			OperatorName:       x.OperatorName,
		}
		if x.EventTime != nil { res[i].EventTime = x.EventTime.Format(time.RFC3339) }
	}
	return res
}

func (s *TraceabilityService) TraceBackward(ctx context.Context, req *pb.TraceRequest) (*pb.TraceReply, error) {
	list, err := s.uc.TraceBackward(ctx, req.LotId)
	if err != nil {
		return nil, err
	}
	return &pb.TraceReply{
		QueriedLotId:   req.LotId,
		TraceDirection: "backward",
		Nodes:          mapTraceNodes(list),
	}, nil
}

func (s *TraceabilityService) TraceForward(ctx context.Context, req *pb.TraceRequest) (*pb.TraceReply, error) {
	list, err := s.uc.TraceForward(ctx, req.LotId)
	if err != nil {
		return nil, err
	}
	return &pb.TraceReply{
		QueriedLotId:   req.LotId,
		TraceDirection: "forward",
		Nodes:          mapTraceNodes(list),
	}, nil
}

func (s *TraceabilityService) TraceFullGenealogy(ctx context.Context, req *pb.TraceRequest) (*pb.GenealogyTreeReply, error) {
	bNodes, bEdges, err := s.uc.TraceFullGenealogy(ctx, req.LotId)
	if err != nil {
		return nil, err
	}
	
	nodes := make([]*pb.GenealogyNode, len(bNodes))
	for i, x := range bNodes {
		nodes[i] = &pb.GenealogyNode{
			LotId:         x.LotID,
			MaterialName:  x.MaterialName,
			MaterialType:  x.MaterialType,
			Status:        x.Status,
			Quantity:      x.Quantity,
			UnitOfMeasure: x.UnitOfMeasure,
		}
	}
	
	edges := make([]*pb.GenealogyEdge, len(bEdges))
	for i, x := range bEdges {
		edges[i] = &pb.GenealogyEdge{
			SourceLotId:      x.SourceLotID,
			TargetLotId:      x.TargetLotID,
			WorkOrderId:      x.WorkOrderID,
			EquipmentId:      x.EquipmentID,
			QuantityConsumed: x.QuantityConsumed,
			QuantityProduced: x.QuantityProduced,
		}
		if x.EventTime != nil { edges[i].EventTime = x.EventTime.Format(time.RFC3339) }
	}
	
	return &pb.GenealogyTreeReply{
		CenterLotId: req.LotId,
		Nodes:       nodes,
		Edges:       edges,
	}, nil
}

// 17. Historical Analytics 

func mapMachineMetrics(m *biz.MachineMetrics) *pb.MachineMetrics {
	return &pb.MachineMetrics{
		EquipmentId:        m.EquipmentID,
		EquipmentClassName: m.EquipmentClassName,
		OperationalStatus:  m.OperationalStatus,
		TotalWorkOrders:    m.TotalWorkOrders,
		TotalUnitsProduced: m.TotalUnitsProduced,
		AvgCycleTimeHours:  m.AvgCycleTimeHours,
		UtilizationPct:     m.UtilizationPct,
		ErrorCount:         m.ErrorCount,
	}
}

func (s *TraceabilityService) GetMachinePerformance(ctx context.Context, req *pb.MachinePerformanceRequest) (*pb.MachinePerformanceReply, error) {
	from, _ := time.Parse(time.RFC3339, req.FromTime)
	to, _ := time.Parse(time.RFC3339, req.ToTime)
	result, err := s.uc.GetMachinePerformance(ctx, &biz.MachinePerformanceRequest{
		EquipmentID: req.EquipmentId,
		From:        from,
		To:          to,
	})
	if err != nil {
		return nil, err
	}
	events := make([]*pb.MachineEventSummary, len(result.EventSummary))
	for i, e := range result.EventSummary {
		events[i] = &pb.MachineEventSummary{EventType: e.EventType, Count: e.Count}
	}
	return &pb.MachinePerformanceReply{
		EquipmentId:        result.EquipmentID,
		EquipmentClassName: result.EquipmentClassName,
		OperationalStatus:  result.OperationalStatus,
		FromTime:           req.FromTime,
		ToTime:             req.ToTime,
		TotalWorkOrders:    result.TotalWorkOrders,
		TotalUnitsProduced: result.TotalUnitsProduced,
		AvgCycleTimeHours:  result.AvgCycleTimeHours,
		TotalActiveHours:   result.TotalActiveHours,
		UtilizationPct:     result.UtilizationPct,
		EventSummary:       events,
	}, nil
}

func (s *TraceabilityService) CompareMachinePerformance(ctx context.Context, req *pb.MachineComparisonRequest) (*pb.MachineComparisonReply, error) {
	from, _ := time.Parse(time.RFC3339, req.FromTime)
	to, _ := time.Parse(time.RFC3339, req.ToTime)
	list, err := s.uc.CompareMachinePerformance(ctx, &biz.MachineComparisonRequest{
		EquipmentIDs: req.EquipmentIds,
		From:         from,
		To:           to,
	})
	if err != nil {
		return nil, err
	}
	machines := make([]*pb.MachineMetrics, len(list))
	for i, m := range list {
		machines[i] = mapMachineMetrics(m)
	}
	return &pb.MachineComparisonReply{
		FromTime: req.FromTime,
		ToTime:   req.ToTime,
		Machines: machines,
	}, nil
}

func (s *TraceabilityService) GetProductionTrends(ctx context.Context, req *pb.ProductionTrendsRequest) (*pb.ProductionTrendsReply, error) {
	from, _ := time.Parse(time.RFC3339, req.FromTime)
	to, _ := time.Parse(time.RFC3339, req.ToTime)
	granularity := req.Granularity
	if granularity == "" {
		granularity = "daily"
	}
	list, err := s.uc.GetProductionTrends(ctx, &biz.ProductionTrendsRequest{
		From:        from,
		To:          to,
		Granularity: granularity,
		EquipmentID: req.EquipmentId,
	})
	if err != nil {
		return nil, err
	}
	points := make([]*pb.ProductionTrendPoint, len(list))
	for i, p := range list {
		points[i] = &pb.ProductionTrendPoint{
			Period:              p.Period,
			EquipmentId:         p.EquipmentID,
			UnitsProduced:       p.UnitsProduced,
			WorkOrdersCompleted: p.WorkOrdersCompleted,
			AvgCycleTimeHours:   p.AvgCycleTimeHours,
		}
	}
	return &pb.ProductionTrendsReply{
		FromTime:   req.FromTime,
		ToTime:     req.ToTime,
		Granularity: granularity,
		DataPoints: points,
	}, nil
}

func (s *TraceabilityService) GetDashboardSummary(ctx context.Context, req *pb.DashboardSummaryRequest) (*pb.DashboardSummaryReply, error) {
	from, _ := time.Parse(time.RFC3339, req.FromTime)
	to, _ := time.Parse(time.RFC3339, req.ToTime)
	summary, err := s.uc.GetDashboardSummary(ctx, &biz.DashboardSummaryRequest{From: from, To: to})
	if err != nil {
		return nil, err
	}
	top := make([]*pb.MachineMetrics, len(summary.TopPerformers))
	for i, m := range summary.TopPerformers {
		top[i] = mapMachineMetrics(m)
	}
	return &pb.DashboardSummaryReply{
		FromTime:                 req.FromTime,
		ToTime:                   req.ToTime,
		TotalUnitsProduced:       summary.TotalUnitsProduced,
		TotalWorkOrdersCompleted: summary.TotalWorkOrdersCompleted,
		ActiveEquipmentCount:     summary.ActiveEquipmentCount,
		LotsReleased:             summary.LotsReleased,
		LotsQuarantined:          summary.LotsQuarantined,
		QualityRatePct:           summary.QualityRatePct,
		TotalErrorEvents:         summary.TotalErrorEvents,
		TopPerformers:            top,
	}, nil
}

func (s *TraceabilityService) GetEquipmentProcessHistory(ctx context.Context, req *pb.TraceRequest) (*pb.EquipmentProcessHistoryReply, error) {
	sParams, sReadings, err := s.uc.GetEquipmentProcessHistory(ctx, req.LotId)
	if err != nil {
		return nil, err
	}
	
	params := make([]*pb.EquipmentParameterSummary, len(sParams))
	for i, x := range sParams {
		params[i] = &pb.EquipmentParameterSummary{
			EquipmentId:        x.EquipmentID,
			EquipmentClassName: x.EquipmentClassName,
			ParameterName:      x.ParameterName,
			UnitOfMeasure:      x.UnitOfMeasure,
			MinValue:           x.MinValue,
			MaxValue:           x.MaxValue,
			AvgValue:           x.AvgValue,
			ReadingCount:       x.ReadingCount,
		}
		if x.ProcessStart != nil { params[i].ProcessStart = x.ProcessStart.Format(time.RFC3339) }
		if x.ProcessEnd != nil { params[i].ProcessEnd = x.ProcessEnd.Format(time.RFC3339) }
	}
	
	readings := make([]*pb.TelemetryReading, len(sReadings))
	for i, x := range sReadings {
		readings[i] = &pb.TelemetryReading{
			EquipmentId:   x.EquipmentID,
			ParameterName: x.ParameterName,
			Value:         x.Value,
			UnitOfMeasure: x.UnitOfMeasure,
		}
		if x.RecordedAt != nil { readings[i].RecordedAt = x.RecordedAt.Format(time.RFC3339) }
	}
	
	return &pb.EquipmentProcessHistoryReply{
		LotId:      req.LotId,
		Parameters: params,
		Readings:   readings,
	}, nil
}

