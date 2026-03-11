package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	"traceability/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// 1. ORM Definitions (Database Tables)

type EnterpriseORM struct {
	ID          int32 `gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string
}

func (EnterpriseORM) TableName() string { return "enterprise" }

type SiteORM struct {
	ID           int32 `gorm:"primaryKey;autoIncrement"`
	EnterpriseID int32
	Name         string
	Location     string
	Description  string
}

func (SiteORM) TableName() string { return "site" }

type AreaORM struct {
	ID          int32 `gorm:"primaryKey;autoIncrement"`
	SiteID      int32
	Name        string
	Description string
}

func (AreaORM) TableName() string { return "area" }

type LineORM struct {
	ID          int32 `gorm:"primaryKey;autoIncrement"`
	AreaID      int32
	Name        string
	Description string
}

func (LineORM) TableName() string { return "production_line" }

type UnitORM struct {
	ID               int32 `gorm:"primaryKey;autoIncrement"`
	ProductionLineID int32
	Name             string
	Description      string
}

func (UnitORM) TableName() string { return "production_unit" }

type ClassORM struct {
	ID          int32 `gorm:"primaryKey;autoIncrement"`
	ClassName   string
	Version     string
	Description string
}

func (ClassORM) TableName() string { return "equipment_class" }

type EquipmentORM struct {
	ID                string `gorm:"primaryKey"`
	PhysicalAssetID   string
	ProductionUnitID  int32
	EquipmentClassID  int32
	OperationalStatus string
	ParentEquipmentID *string
}

func (EquipmentORM) TableName() string { return "equipment_master" }

type CapabilityORM struct {
	ID             int32 `gorm:"primaryKey;autoIncrement"`
	EquipmentID    string
	CapabilityName string
	Value          string
	UOM            string
	Description    string
}

func (CapabilityORM) TableName() string { return "equipment_capability" }

type PropertyORM struct {
	ID           int32 `gorm:"primaryKey;autoIncrement"`
	EquipmentID  string
	PropertyName string
	CurrentValue string
	LastUpdated  time.Time
}

func (PropertyORM) TableName() string { return "equipment_property" }

type LogORM struct {
	ID            int32 `gorm:"primaryKey;autoIncrement"`
	EquipmentID   string
	EventType     string
	WorkOrderID   string
	MaterialLotID string
	OperatorID    string
	EventTime     time.Time
}

func (LogORM) TableName() string { return "traceability_log" }

type MaterialDefinitionORM struct {
	ID            int32  `gorm:"primaryKey;autoIncrement;column:id"`
	Name          string `gorm:"column:name"`
	MaterialType  string `gorm:"column:material_type"`
	UnitOfMeasure string `gorm:"column:unit_of_measure"`
	Description   string `gorm:"column:description"`
}

func (MaterialDefinitionORM) TableName() string { return "material_definition" }

type MaterialLotORM struct {
	LotID                string    `gorm:"primaryKey;column:lot_id"`
	MaterialDefinitionID int32     `gorm:"column:material_definition_id"`
	Quantity             float64   `gorm:"column:quantity"`
	UnitOfMeasure        string    `gorm:"column:unit_of_measure"`
	Status               string    `gorm:"column:status"`
	CreatedAt            time.Time `gorm:"column:created_at"`
	UpdatedAt            time.Time `gorm:"column:updated_at"`
}

func (MaterialLotORM) TableName() string { return "material_lot" }

type OperatorORM struct {
	OperatorID string `gorm:"primaryKey;column:operator_id"`
	Name       string `gorm:"column:name"`
	Role       string `gorm:"column:role"`
	Shift      string `gorm:"column:shift"`
	Status     string `gorm:"column:status"`
}

func (OperatorORM) TableName() string { return "operator" }

type WorkOrderORM struct {
	WorkOrderID     string     `gorm:"primaryKey;column:work_order_id"`
	Description     string     `gorm:"column:description"`
	Status          string     `gorm:"column:status"`
	EquipmentID     string     `gorm:"column:equipment_id"`
	OperatorID      string     `gorm:"column:operator_id"`
	OutputLotID     string     `gorm:"column:output_lot_id"`
	PlannedQuantity float64    `gorm:"column:planned_quantity"`
	ActualQuantity  *float64   `gorm:"column:actual_quantity"`
	UnitOfMeasure   string     `gorm:"column:unit_of_measure"`
	PlannedStart    time.Time  `gorm:"column:planned_start"`
	PlannedEnd      time.Time  `gorm:"column:planned_end"`
	ActualStart     *time.Time `gorm:"column:actual_start"`
	ActualEnd       *time.Time `gorm:"column:actual_end"`
}

func (WorkOrderORM) TableName() string { return "work_order" }

type LotGenealogyORM struct {
	ID               int32     `gorm:"primaryKey;autoIncrement;column:id"`
	ParentLotID      string    `gorm:"column:parent_lot_id"`
	ChildLotID       string    `gorm:"column:child_lot_id"`
	WorkOrderID      string    `gorm:"column:work_order_id"`
	EquipmentID      string    `gorm:"column:equipment_id"`
	QuantityConsumed float64   `gorm:"column:quantity_consumed"`
	QuantityProduced float64   `gorm:"column:quantity_produced"`
	EventTime        time.Time `gorm:"column:event_time"`
}

func (LotGenealogyORM) TableName() string { return "lot_genealogy" }

// 2. Repository Implementation

type traceabilityRepo struct {
	data *Data
	log  *log.Helper
}

func NewTraceabilityRepo(data *Data, logger log.Logger) biz.TraceabilityRepo {
	return &traceabilityRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Enterprise
func (r *traceabilityRepo) CreateEnterprise(ctx context.Context, e *biz.Enterprise) (int32, error) {
	orm := EnterpriseORM{Name: e.Name, Description: e.Description}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListEnterprises(ctx context.Context) ([]*biz.Enterprise, error) {
	var dbList []EnterpriseORM
	if err := r.data.db.WithContext(ctx).Find(&dbList).Error; err != nil {
		return nil, err
	}
	var list []*biz.Enterprise
	for _, x := range dbList {
		list = append(list, &biz.Enterprise{ID: x.ID, Name: x.Name, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateEnterprise(ctx context.Context, e *biz.Enterprise) error {
	return r.data.db.WithContext(ctx).Model(&EnterpriseORM{}).Where("id = ?", e.ID).Updates(EnterpriseORM{Name: e.Name, Description: e.Description}).Error
}
func (r *traceabilityRepo) DeleteEnterprise(ctx context.Context, id int32) error {
	return r.data.db.WithContext(ctx).Delete(&EnterpriseORM{}, id).Error
}

// Site
func (r *traceabilityRepo) CreateSite(ctx context.Context, s *biz.Site) (int32, error) {
	orm := SiteORM{EnterpriseID: s.EnterpriseID, Name: s.Name, Location: s.Location, Description: s.Description}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListSites(ctx context.Context, eid int32) ([]*biz.Site, error) {
	var dbList []SiteORM
	query := r.data.db.WithContext(ctx)
	if eid > 0 {
		query = query.Where("enterprise_id = ?", eid)
	}
	if err := query.Find(&dbList).Error; err != nil {
		return nil, err
	}
	var list []*biz.Site
	for _, x := range dbList {
		list = append(list, &biz.Site{ID: x.ID, EnterpriseID: x.EnterpriseID, Name: x.Name, Location: x.Location, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateSite(ctx context.Context, s *biz.Site) error {
	return r.data.db.WithContext(ctx).Model(&SiteORM{}).Where("id = ?", s.ID).Updates(SiteORM{Location: s.Location}).Error
}
func (r *traceabilityRepo) DeleteSite(ctx context.Context, id int32) error {
	return r.data.db.WithContext(ctx).Delete(&SiteORM{}, id).Error
}

// Area
func (r *traceabilityRepo) CreateArea(ctx context.Context, a *biz.Area) (int32, error) {
	orm := AreaORM{SiteID: a.SiteID, Name: a.Name, Description: a.Description}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListAreas(ctx context.Context, sid int32) ([]*biz.Area, error) {
	var dbList []AreaORM
	query := r.data.db.WithContext(ctx)
	if sid > 0 {
		query = query.Where("site_id = ?", sid)
	}
	if err := query.Find(&dbList).Error; err != nil {
		return nil, err
	}
	var list []*biz.Area
	for _, x := range dbList {
		list = append(list, &biz.Area{ID: x.ID, SiteID: x.SiteID, Name: x.Name, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateArea(ctx context.Context, a *biz.Area) error {
	return r.data.db.WithContext(ctx).Model(&AreaORM{}).Where("id = ?", a.ID).Updates(AreaORM{Name: a.Name, Description: a.Description}).Error
}
func (r *traceabilityRepo) DeleteArea(ctx context.Context, id int32) error {
	return r.data.db.WithContext(ctx).Delete(&AreaORM{}, id).Error
}

// Line
func (r *traceabilityRepo) CreateLine(ctx context.Context, l *biz.ProductionLine) (int32, error) {
	orm := LineORM{AreaID: l.AreaID, Name: l.Name, Description: l.Description}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListLines(ctx context.Context, aid int32) ([]*biz.ProductionLine, error) {
	var dbList []LineORM
	query := r.data.db.WithContext(ctx)
	if aid > 0 {
		query = query.Where("area_id = ?", aid)
	}
	if err := query.Find(&dbList).Error; err != nil {
		return nil, err
	}
	var list []*biz.ProductionLine
	for _, x := range dbList {
		list = append(list, &biz.ProductionLine{ID: x.ID, AreaID: x.AreaID, Name: x.Name, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateLine(ctx context.Context, l *biz.ProductionLine) error {
	return r.data.db.WithContext(ctx).Model(&LineORM{}).Where("id = ?", l.ID).Updates(LineORM{Name: l.Name, Description: l.Description}).Error
}
func (r *traceabilityRepo) DeleteLine(ctx context.Context, id int32) error {
	return r.data.db.WithContext(ctx).Delete(&LineORM{}, id).Error
}

// Unit
func (r *traceabilityRepo) CreateProductionUnit(ctx context.Context, u *biz.ProductionUnit) (int32, error) {
	orm := UnitORM{ProductionLineID: u.ProductionLineID, Name: u.Name, Description: u.Description}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListProductionUnits(ctx context.Context, lid int32) ([]*biz.ProductionUnit, error) {
	var dbList []UnitORM
	query := r.data.db.WithContext(ctx)
	if lid > 0 {
		query = query.Where("production_line_id = ?", lid)
	}
	if err := query.Find(&dbList).Error; err != nil {
		return nil, err
	}
	var list []*biz.ProductionUnit
	for _, x := range dbList {
		list = append(list, &biz.ProductionUnit{ID: x.ID, ProductionLineID: x.ProductionLineID, Name: x.Name, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateProductionUnit(ctx context.Context, u *biz.ProductionUnit) error {
    updates := map[string]interface{}{}
    if u.Name != "" {
        updates["name"] = u.Name
    }
    if u.Description != "" {
        updates["description"] = u.Description
    }
    return r.data.db.WithContext(ctx).Model(&UnitORM{}).Where("id = ?", u.ID).Updates(updates).Error
}
func (r *traceabilityRepo) DeleteProductionUnit(ctx context.Context, id int32) error {
	return r.data.db.WithContext(ctx).Delete(&UnitORM{}, id).Error
}

// Class
func (r *traceabilityRepo) CreateEquipmentClass(ctx context.Context, c *biz.EquipmentClass) (int32, error) {
	orm := ClassORM{ClassName: c.ClassName, Version: c.Version, Description: c.Description}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListEquipmentClasses(ctx context.Context) ([]*biz.EquipmentClass, error) {
	var dbList []ClassORM
	if err := r.data.db.WithContext(ctx).Find(&dbList).Error; err != nil {
		return nil, err
	}
	var list []*biz.EquipmentClass
	for _, x := range dbList {
		list = append(list, &biz.EquipmentClass{ID: x.ID, ClassName: x.ClassName, Version: x.Version, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateEquipmentClass(ctx context.Context, c *biz.EquipmentClass) error {
	updates := map[string]interface{}{"version": c.Version}
	if c.Description != "" {
		updates["description"] = c.Description
	}
	return r.data.db.WithContext(ctx).Model(&ClassORM{}).Where("id = ?", c.ID).Updates(updates).Error
}
func (r *traceabilityRepo) DeleteEquipmentClass(ctx context.Context, id int32) error {
	return r.data.db.WithContext(ctx).Delete(&ClassORM{}, id).Error
}

// Equipment
func (r *traceabilityRepo) RegisterEquipment(ctx context.Context, e *biz.EquipmentMaster) (string, error) {
	var parentID *string
	if e.ParentEquipmentID != "" {
		parentID = &e.ParentEquipmentID
	}

	orm := EquipmentORM{
		ID:                e.ID,
		PhysicalAssetID:   e.PhysicalAssetID,
		ProductionUnitID:  e.ProductionUnitID,
		EquipmentClassID:  e.EquipmentClassID,
		OperationalStatus: e.OperationalStatus,
		ParentEquipmentID: parentID,
	}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListEquipment(ctx context.Context, lid int32, pid string) ([]*biz.EquipmentMaster, error) {
	var dbList []EquipmentORM
	query := r.data.db.WithContext(ctx)

	if lid > 0 {
		query = query.Joins("JOIN production_unit ON production_unit.id = equipment_master.production_unit_id").
			Where("production_unit.production_line_id = ?", lid)
	}
	if pid != "" {
		query = query.Where("parent_equipment_id = ?", pid)
	}

	if err := query.Find(&dbList).Error; err != nil {
		return nil, err
	}
	var list []*biz.EquipmentMaster
	for _, x := range dbList {
		parentID := ""
		if x.ParentEquipmentID != nil {
			parentID = *x.ParentEquipmentID
		}

		list = append(list, &biz.EquipmentMaster{
			ID: x.ID, PhysicalAssetID: x.PhysicalAssetID, ProductionUnitID: x.ProductionUnitID,
			EquipmentClassID: x.EquipmentClassID, OperationalStatus: x.OperationalStatus, ParentEquipmentID: parentID,
		})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateEquipment(ctx context.Context, e *biz.EquipmentMaster) error {
	updates := map[string]interface{}{}
	if e.OperationalStatus != "" {
		updates["operational_status"] = e.OperationalStatus
	}
	if e.ProductionUnitID > 0 {
		updates["production_unit_id"] = e.ProductionUnitID
	}

	return r.data.db.WithContext(ctx).Model(&EquipmentORM{}).Where("id = ?", e.ID).Updates(updates).Error
}
func (r *traceabilityRepo) DeleteEquipment(ctx context.Context, id string) error {
	return r.data.db.WithContext(ctx).Where("id = ?", id).Delete(&EquipmentORM{}).Error
}

// Capability
func (r *traceabilityRepo) AddCapability(ctx context.Context, c *biz.EquipmentCapability) (int32, error) {
	orm := CapabilityORM{EquipmentID: c.EquipmentID, CapabilityName: c.CapabilityName, Value: c.Value, UOM: c.UOM, Description: c.Description}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListCapabilities(ctx context.Context, eid string) ([]*biz.EquipmentCapability, error) {
	var dbList []CapabilityORM
	if err := r.data.db.WithContext(ctx).Where("equipment_id = ?", eid).Find(&dbList).Error; err != nil {
		return nil, err
	}
	var list []*biz.EquipmentCapability
	for _, x := range dbList {
		list = append(list, &biz.EquipmentCapability{ID: x.ID, EquipmentID: x.EquipmentID, CapabilityName: x.CapabilityName, Value: x.Value, UOM: x.UOM, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateCapability(ctx context.Context, c *biz.EquipmentCapability) error {
	return r.data.db.WithContext(ctx).Model(&CapabilityORM{}).Where("equipment_id = ? AND id = ?", c.EquipmentID, c.ID).Updates(CapabilityORM{Value: c.Value}).Error
}
func (r *traceabilityRepo) DeleteCapability(ctx context.Context, eid string, id int32) error {
	return r.data.db.WithContext(ctx).Where("equipment_id = ? AND id = ?", eid, id).Delete(&CapabilityORM{}).Error
}

// Property
func (r *traceabilityRepo) SetProperty(ctx context.Context, p *biz.EquipmentProperty) (int32, error) {
	orm := PropertyORM{EquipmentID: p.EquipmentID, PropertyName: p.PropertyName, CurrentValue: p.CurrentValue, LastUpdated: time.Now()}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListProperties(ctx context.Context, eid string) ([]*biz.EquipmentProperty, error) {
	var dbList []PropertyORM
	if err := r.data.db.WithContext(ctx).Where("equipment_id = ?", eid).Find(&dbList).Error; err != nil {
		return nil, err
	}
	var list []*biz.EquipmentProperty
	for _, x := range dbList {
		list = append(list, &biz.EquipmentProperty{ID: x.ID, EquipmentID: x.EquipmentID, PropertyName: x.PropertyName, CurrentValue: x.CurrentValue, LastUpdated: x.LastUpdated})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateProperty(ctx context.Context, p *biz.EquipmentProperty) error {
	return r.data.db.WithContext(ctx).Model(&PropertyORM{}).Where("equipment_id = ? AND id = ?", p.EquipmentID, p.ID).Updates(PropertyORM{CurrentValue: p.CurrentValue, LastUpdated: time.Now()}).Error
}
func (r *traceabilityRepo) DeleteProperty(ctx context.Context, eid string, id int32) error {
	return r.data.db.WithContext(ctx).Where("equipment_id = ? AND id = ?", eid, id).Delete(&PropertyORM{}).Error
}

// Logs
func (r *traceabilityRepo) LogEvent(ctx context.Context, l *biz.TraceabilityLog) (int32, error) {
	orm := LogORM{EquipmentID: l.EquipmentID, EventType: l.EventType, WorkOrderID: l.WorkOrderID, MaterialLotID: l.MaterialLotID, OperatorID: l.OperatorID, EventTime: time.Now()}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListLogs(ctx context.Context, wid string) ([]*biz.TraceabilityLog, error) {
	var dbList []LogORM
	if err := r.data.db.WithContext(ctx).Where("work_order_id = ?", wid).Order("event_time desc").Find(&dbList).Error; err != nil {
		return nil, err
	}
	var list []*biz.TraceabilityLog
	for _, x := range dbList {
		list = append(list, &biz.TraceabilityLog{ID: x.ID, EquipmentID: x.EquipmentID, EventType: x.EventType, WorkOrderID: x.WorkOrderID, MaterialLotID: x.MaterialLotID, OperatorID: x.OperatorID, EventTime: x.EventTime})
	}
	return list, nil
}

// 11. Material Definition

func (r *traceabilityRepo) CreateMaterialDefinition(ctx context.Context, md *biz.MaterialDefinition) (int32, error) {
	orm := MaterialDefinitionORM{
		Name:          md.Name,
		MaterialType:  md.MaterialType,
		UnitOfMeasure: md.UnitOfMeasure,
		Description:   md.Description,
	}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}

func (r *traceabilityRepo) ListMaterialDefinitions(ctx context.Context, mType string) ([]*biz.MaterialDefinition, error) {
	var orms []MaterialDefinitionORM
	query := r.data.db.WithContext(ctx)
	if mType != "" {
		query = query.Where("material_type = ?", mType)
	}
	if err := query.Find(&orms).Error; err != nil {
		return nil, err
	}
	list := make([]*biz.MaterialDefinition, len(orms))
	for i, x := range orms {
		list[i] = &biz.MaterialDefinition{
			ID:            x.ID,
			Name:          x.Name,
			MaterialType:  x.MaterialType,
			UnitOfMeasure: x.UnitOfMeasure,
			Description:   x.Description,
		}
	}
	return list, nil
}

// 12. Material Lot

func (r *traceabilityRepo) CreateMaterialLot(ctx context.Context, ml *biz.MaterialLot) (string, error) {
	orm := MaterialLotORM{
		LotID:                ml.LotID,
		MaterialDefinitionID: ml.MaterialDefinitionID,
		Quantity:             ml.Quantity,
		UnitOfMeasure:        ml.UnitOfMeasure,
		Status:               ml.Status,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.LotID, res.Error
}

func (r *traceabilityRepo) GetMaterialLot(ctx context.Context, id string) (*biz.MaterialLotDetail, error) {
    sqlStr := `SELECT ml.lot_id, ml.status, ml.quantity, ml.unit_of_measure, ml.created_at, ml.updated_at,
        md.id AS material_definition_id, md.name AS material_name, md.material_type, md.description AS material_description
        FROM material_lot ml JOIN material_definition md ON md.id = ml.material_definition_id WHERE ml.lot_id = ?`
    var res biz.MaterialLotDetail
    if err := r.data.db.WithContext(ctx).Raw(sqlStr, id).Scan(&res).Error; err != nil {
        return nil, err
    }
    if res.LotID == "" {
        return nil, fmt.Errorf("material lot not found: %s", id)
    }
    return &res, nil
}

func (r *traceabilityRepo) ListMaterialLots(ctx context.Context, status, mType string) ([]*biz.MaterialLotSummary, error) {
	sqlStr := `SELECT ml.lot_id, ml.status, ml.quantity, ml.unit_of_measure, ml.created_at,
		md.name AS material_name, md.material_type
		FROM material_lot ml JOIN material_definition md ON md.id = ml.material_definition_id WHERE 1=1`
	var args []interface{}
	if status != "" {
		sqlStr += " AND ml.status = ?"
		args = append(args, status)
	}
	if mType != "" {
		sqlStr += " AND md.material_type = ?"
		args = append(args, mType)
	}
	sqlStr += " ORDER BY ml.created_at DESC"
	
	var list []*biz.MaterialLotSummary
	if err := r.data.db.WithContext(ctx).Raw(sqlStr, args...).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *traceabilityRepo) UpdateMaterialLotStatus(ctx context.Context, id, status string) error {
	return r.data.db.WithContext(ctx).Model(&MaterialLotORM{}).Where("lot_id = ?", id).Updates(map[string]interface{}{"status": status, "updated_at": time.Now()}).Error
}

// 13. Operator

func (r *traceabilityRepo) CreateOperator(ctx context.Context, op *biz.Operator) (string, error) {
	orm := OperatorORM{
		OperatorID: op.OperatorID,
		Name:       op.Name,
		Role:       op.Role,
		Shift:      op.Shift,
		Status:     op.Status,
	}
	if orm.Status == "" {
		orm.Status = "active"
	}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.OperatorID, res.Error
}

func (r *traceabilityRepo) ListOperators(ctx context.Context, shift, status string) ([]*biz.Operator, error) {
	var orms []OperatorORM
	query := r.data.db.WithContext(ctx)
	if shift != "" {
		query = query.Where("shift = ?", shift)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Find(&orms).Error; err != nil {
		return nil, err
	}
	var list []*biz.Operator
	for _, x := range orms {
		list = append(list, &biz.Operator{OperatorID: x.OperatorID, Name: x.Name, Role: x.Role, Shift: x.Shift, Status: x.Status})
	}
	return list, nil
}

func (r *traceabilityRepo) UpdateOperator(ctx context.Context, op *biz.Operator) error {
	updates := map[string]interface{}{}
	if op.Role != "" { updates["role"] = op.Role }
	if op.Shift != "" { updates["shift"] = op.Shift }
	if op.Status != "" { updates["status"] = op.Status }
	return r.data.db.WithContext(ctx).Model(&OperatorORM{}).Where("operator_id = ?", op.OperatorID).Updates(updates).Error
}

// 14. Work Order

func (r *traceabilityRepo) CreateWorkOrder(ctx context.Context, wo *biz.WorkOrder) (string, error) {
	orm := WorkOrderORM{
		WorkOrderID:     wo.WorkOrderID,
		Description:     wo.Description,
		Status:          wo.Status,
		EquipmentID:     wo.EquipmentID,
		OperatorID:      wo.OperatorID,
		OutputLotID:     wo.OutputLotID,
		PlannedQuantity: wo.PlannedQuantity,
		UnitOfMeasure:   wo.UnitOfMeasure,
		PlannedStart:    wo.PlannedStart,
		PlannedEnd:      wo.PlannedEnd,
	}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.WorkOrderID, res.Error
}

func (r *traceabilityRepo) GetWorkOrder(ctx context.Context, id string) (*biz.WorkOrderDetail, error) {
	type woScan struct {
		WorkOrderID        string     `gorm:"column:work_order_id"`
		Description        string     `gorm:"column:description"`
		Status             string     `gorm:"column:status"`
		PlannedQuantity    float64    `gorm:"column:planned_quantity"`
		ActualQuantity     *float64   `gorm:"column:actual_quantity"`
		UnitOfMeasure      string     `gorm:"column:unit_of_measure"`
		PlannedStart       time.Time  `gorm:"column:planned_start"`
		PlannedEnd         time.Time  `gorm:"column:planned_end"`
		ActualStart        *time.Time `gorm:"column:actual_start"`
		ActualEnd          *time.Time `gorm:"column:actual_end"`
		OutputLotID        string     `gorm:"column:output_lot_id"`
		EquipmentID        string     `gorm:"column:equipment_id"`
		EquipmentClassName string     `gorm:"column:equipment_class_name"`
		OperatorID         string     `gorm:"column:operator_id"`
		OperatorName       string     `gorm:"column:operator_name"`
		OperatorShift      string     `gorm:"column:operator_shift"`
	}
	sql1 := `SELECT wo.work_order_id, wo.description, wo.status, wo.planned_quantity, wo.actual_quantity,
		wo.unit_of_measure, wo.planned_start, wo.planned_end, wo.actual_start, wo.actual_end, wo.output_lot_id,
		wo.equipment_id, ec.class_name AS equipment_class_name, wo.operator_id, op.name AS operator_name, op.shift AS operator_shift
		FROM work_order wo
		LEFT JOIN equipment_master em ON em.id = wo.equipment_id
		LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
		LEFT JOIN operator op ON op.operator_id = wo.operator_id
		WHERE wo.work_order_id = ?`
	var scan woScan
	if err := r.data.db.WithContext(ctx).Raw(sql1, id).Scan(&scan).Error; err != nil {
		return nil, err
	}
	if scan.WorkOrderID == "" {
		return nil, fmt.Errorf("work order not found: %s", id)
	}
	sql2 := `SELECT lg.parent_lot_id AS lot_id, md.name AS material_name, md.material_type, lg.quantity_consumed, ml.unit_of_measure
		FROM lot_genealogy lg
		JOIN material_lot ml ON ml.lot_id = lg.parent_lot_id
		JOIN material_definition md ON md.id = ml.material_definition_id
		WHERE lg.work_order_id = ?`
	var inputs []*biz.WorkOrderInputLot
	if err := r.data.db.WithContext(ctx).Raw(sql2, id).Scan(&inputs).Error; err != nil {
		return nil, err
	}
	return &biz.WorkOrderDetail{
		WorkOrder: biz.WorkOrder{
			WorkOrderID:     scan.WorkOrderID,
			Description:     scan.Description,
			Status:          scan.Status,
			EquipmentID:     scan.EquipmentID,
			OperatorID:      scan.OperatorID,
			OutputLotID:     scan.OutputLotID,
			PlannedQuantity: scan.PlannedQuantity,
			ActualQuantity:  scan.ActualQuantity,
			UnitOfMeasure:   scan.UnitOfMeasure,
			PlannedStart:    scan.PlannedStart,
			PlannedEnd:      scan.PlannedEnd,
			ActualStart:     scan.ActualStart,
			ActualEnd:       scan.ActualEnd,
		},
		EquipmentClassName: scan.EquipmentClassName,
		OperatorName:       scan.OperatorName,
		OperatorShift:      scan.OperatorShift,
		InputLots:          inputs,
	}, nil
}

func (r *traceabilityRepo) ListWorkOrders(ctx context.Context, status, eqID string) ([]*biz.WorkOrderSummary, error) {
	query := r.data.db.WithContext(ctx).Model(&WorkOrderORM{})
	if status != "" { query = query.Where("status = ?", status) }
	if eqID != "" { query = query.Where("equipment_id = ?", eqID) }
	var orms []WorkOrderORM
	if err := query.Find(&orms).Error; err != nil {
		return nil, err
	}
	var list []*biz.WorkOrderSummary
	for _, x := range orms {
		list = append(list, &biz.WorkOrderSummary{
			WorkOrderID: x.WorkOrderID, Description: x.Description, Status: x.Status,
			EquipmentID: x.EquipmentID, OperatorID: x.OperatorID, OutputLotID: x.OutputLotID,
			ActualStart: x.ActualStart, ActualEnd: x.ActualEnd,
		})
	}
	return list, nil
}

func (r *traceabilityRepo) UpdateWorkOrderStatus(ctx context.Context, wo *biz.WorkOrder) error {
	updates := map[string]interface{}{}
	if wo.Status != "" { updates["status"] = wo.Status }
	if wo.ActualQuantity != nil { updates["actual_quantity"] = *wo.ActualQuantity }
	if wo.ActualStart != nil { updates["actual_start"] = *wo.ActualStart }
	if wo.ActualEnd != nil { updates["actual_end"] = *wo.ActualEnd }
	return r.data.db.WithContext(ctx).Model(&WorkOrderORM{}).Where("work_order_id = ?", wo.WorkOrderID).Updates(updates).Error
}

// 15. Genealogy
func (r *traceabilityRepo) RegisterGenealogyLink(ctx context.Context, lg *biz.LotGenealogy) (int32, error) {
	orm := LotGenealogyORM{
		ParentLotID:      lg.ParentLotID,
		ChildLotID:       lg.ChildLotID,
		WorkOrderID:      lg.WorkOrderID,
		EquipmentID:      lg.EquipmentID,
		QuantityConsumed: lg.QuantityConsumed,
		QuantityProduced: lg.QuantityProduced,
		EventTime:        time.Now(),
	}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}

// 16. Trace Queries

func (r *traceabilityRepo) TraceBackward(ctx context.Context, lotID string) ([]*biz.TraceNode, error) {
	sqlStr := `WITH RECURSIVE backward_trace AS (
		SELECT lg.parent_lot_id, lg.child_lot_id, lg.work_order_id, lg.equipment_id, lg.quantity_consumed AS qty_used, lg.event_time, 1 AS depth
		FROM lot_genealogy lg WHERE lg.child_lot_id = ?
		UNION ALL
		SELECT lg.parent_lot_id, lg.child_lot_id, lg.work_order_id, lg.equipment_id, lg.quantity_consumed AS qty_used, lg.event_time, bt.depth + 1
		FROM lot_genealogy lg JOIN backward_trace bt ON lg.child_lot_id = bt.parent_lot_id
	)
	SELECT bt.depth, bt.parent_lot_id AS lot_id, bt.child_lot_id AS related_lot_id, md.name AS material_name, md.material_type,
		ml.status AS lot_status, ml.quantity, ml.unit_of_measure, bt.qty_used AS quantity_used,
		bt.work_order_id, bt.equipment_id, ec.class_name AS equipment_class_name, wo.operator_id, op.name AS operator_name, bt.event_time
	FROM backward_trace bt
	JOIN material_lot ml ON ml.lot_id = bt.parent_lot_id JOIN material_definition md ON md.id = ml.material_definition_id
	LEFT JOIN equipment_master em ON em.id = bt.equipment_id LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
	LEFT JOIN work_order wo ON wo.work_order_id = bt.work_order_id LEFT JOIN operator op ON op.operator_id = wo.operator_id
	ORDER BY bt.depth, bt.parent_lot_id`
	var list []*biz.TraceNode
	if err := r.data.db.WithContext(ctx).Raw(sqlStr, lotID).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *traceabilityRepo) TraceForward(ctx context.Context, lotID string) ([]*biz.TraceNode, error) {
	sqlStr := `WITH RECURSIVE forward_trace AS (
		SELECT lg.parent_lot_id, lg.child_lot_id, lg.work_order_id, lg.equipment_id, lg.quantity_produced AS qty_used, lg.event_time, 1 AS depth
		FROM lot_genealogy lg WHERE lg.parent_lot_id = ?
		UNION ALL
		SELECT lg.parent_lot_id, lg.child_lot_id, lg.work_order_id, lg.equipment_id, lg.quantity_produced AS qty_used, lg.event_time, ft.depth + 1
		FROM lot_genealogy lg JOIN forward_trace ft ON lg.parent_lot_id = ft.child_lot_id
	)
	SELECT ft.depth, ft.child_lot_id AS lot_id, ft.parent_lot_id AS related_lot_id, md.name AS material_name, md.material_type,
		ml.status AS lot_status, ml.quantity, ml.unit_of_measure, ft.qty_used AS quantity_used,
		ft.work_order_id, ft.equipment_id, ec.class_name AS equipment_class_name, wo.operator_id, op.name AS operator_name, ft.event_time
	FROM forward_trace ft
	JOIN material_lot ml ON ml.lot_id = ft.child_lot_id JOIN material_definition md ON md.id = ml.material_definition_id
	LEFT JOIN equipment_master em ON em.id = ft.equipment_id LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
	LEFT JOIN work_order wo ON wo.work_order_id = ft.work_order_id LEFT JOIN operator op ON op.operator_id = wo.operator_id
	ORDER BY ft.depth, ft.child_lot_id`
	var list []*biz.TraceNode
	if err := r.data.db.WithContext(ctx).Raw(sqlStr, lotID).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *traceabilityRepo) TraceFullGenealogyNodes(ctx context.Context, lotID string) ([]*biz.GenealogyNode, error) {
	sqlStr := `WITH RECURSIVE backward AS (
		SELECT lg.parent_lot_id AS lot_id FROM lot_genealogy lg WHERE lg.child_lot_id = ?
		UNION SELECT lg.parent_lot_id FROM lot_genealogy lg JOIN backward b ON lg.child_lot_id = b.lot_id
	), forward AS (
		SELECT lg.child_lot_id AS lot_id FROM lot_genealogy lg WHERE lg.parent_lot_id = ?
		UNION SELECT lg.child_lot_id FROM lot_genealogy lg JOIN forward f ON lg.parent_lot_id = f.lot_id
	), all_lot_ids AS (
		SELECT lot_id FROM backward UNION SELECT ? AS lot_id UNION SELECT lot_id FROM forward
	)
	SELECT ml.lot_id, md.name AS material_name, md.material_type, ml.status, ml.quantity, ml.unit_of_measure
	FROM all_lot_ids a JOIN material_lot ml ON ml.lot_id = a.lot_id JOIN material_definition md ON md.id = ml.material_definition_id
	ORDER BY ml.created_at`
	var list []*biz.GenealogyNode
	if err := r.data.db.WithContext(ctx).Raw(sqlStr, lotID, lotID, lotID).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *traceabilityRepo) TraceFullGenealogyEdges(ctx context.Context, lotID string) ([]*biz.GenealogyEdge, error) {
	sqlStr := `WITH RECURSIVE backward AS (
		SELECT lg.parent_lot_id AS lot_id FROM lot_genealogy lg WHERE lg.child_lot_id = ?
		UNION SELECT lg.parent_lot_id FROM lot_genealogy lg JOIN backward b ON lg.child_lot_id = b.lot_id
	), forward AS (
		SELECT lg.child_lot_id AS lot_id FROM lot_genealogy lg WHERE lg.parent_lot_id = ?
		UNION SELECT lg.child_lot_id FROM lot_genealogy lg JOIN forward f ON lg.parent_lot_id = f.lot_id
	), all_lot_ids AS (
		SELECT lot_id FROM backward UNION SELECT ? AS lot_id UNION SELECT lot_id FROM forward
	)
	SELECT lg.parent_lot_id AS source_lot_id, lg.child_lot_id AS target_lot_id, lg.work_order_id, lg.equipment_id,
		lg.quantity_consumed, lg.quantity_produced, lg.event_time
	FROM lot_genealogy lg WHERE lg.parent_lot_id IN (SELECT lot_id FROM all_lot_ids) OR lg.child_lot_id IN (SELECT lot_id FROM all_lot_ids)
	ORDER BY lg.event_time`
	var list []*biz.GenealogyEdge
	if err := r.data.db.WithContext(ctx).Raw(sqlStr, lotID, lotID, lotID).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *traceabilityRepo) GetEquipmentProcessHistorySummary(ctx context.Context, lotID string) ([]*biz.EquipmentParameterSummary, error) {
	sqlStr := `WITH lot_processing_window AS (
		SELECT tl.equipment_id, MIN(tl.event_time) AS process_start, MAX(tl.event_time) AS process_end
		FROM traceability_log tl WHERE tl.material_lot_id = ? GROUP BY tl.equipment_id
	)
	SELECT et.equipment_id, ec.class_name AS equipment_class_name, et.parameter_name, et.unit_of_measure,
		MIN(et.value) AS min_value, MAX(et.value) AS max_value, AVG(et.value) AS avg_value, COUNT(et.value) AS reading_count,
		lpw.process_start, lpw.process_end
	FROM lot_processing_window lpw
	JOIN equipment_telemetry et ON et.equipment_id = lpw.equipment_id AND et.recorded_at >= lpw.process_start AND et.recorded_at <= lpw.process_end
	JOIN equipment_master em ON em.id = et.equipment_id JOIN equipment_class ec ON ec.id = em.equipment_class_id
	GROUP BY et.equipment_id, ec.class_name, et.parameter_name, et.unit_of_measure, lpw.process_start, lpw.process_end
	ORDER BY et.equipment_id, et.parameter_name`
	var list []*biz.EquipmentParameterSummary
	if err := r.data.db.WithContext(ctx).Raw(sqlStr, lotID).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *traceabilityRepo) GetEquipmentProcessHistoryReadings(ctx context.Context, lotID string) ([]*biz.TelemetryReading, error) {
	sqlStr := `WITH lot_processing_window AS (
		SELECT tl.equipment_id, MIN(tl.event_time) AS process_start, MAX(tl.event_time) AS process_end
		FROM traceability_log tl WHERE tl.material_lot_id = ? GROUP BY tl.equipment_id
	)
	SELECT et.equipment_id, et.parameter_name, et.value, et.unit_of_measure, et.recorded_at
	FROM lot_processing_window lpw
	JOIN equipment_telemetry et ON et.equipment_id = lpw.equipment_id AND et.recorded_at >= lpw.process_start AND et.recorded_at <= lpw.process_end
	ORDER BY et.equipment_id, et.parameter_name, et.recorded_at`
	var list []*biz.TelemetryReading
	if err := r.data.db.WithContext(ctx).Raw(sqlStr, lotID).Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// 17. Analytics

func (r *traceabilityRepo) GetMachinePerformance(ctx context.Context, req *biz.MachinePerformanceRequest) (*biz.MachinePerformanceResult, error) {
	type scan struct {
		EquipmentID        string  `gorm:"column:equipment_id"`
		EquipmentClassName string  `gorm:"column:equipment_class_name"`
		OperationalStatus  string  `gorm:"column:operational_status"`
		TotalWorkOrders    int32   `gorm:"column:total_work_orders"`
		TotalUnitsProduced float64 `gorm:"column:total_units_produced"`
		AvgCycleTimeHours  float64 `gorm:"column:avg_cycle_time_hours"`
		TotalActiveHours   float64 `gorm:"column:total_active_hours"`
		UtilizationPct     float64 `gorm:"column:utilization_pct"`
	}
	sqlStr := `
		SELECT em.id AS equipment_id, ec.class_name AS equipment_class_name, em.operational_status,
			COUNT(DISTINCT wo.work_order_id) AS total_work_orders,
			COALESCE(SUM(wo.actual_quantity), 0) AS total_units_produced,
			COALESCE(AVG(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0) AS avg_cycle_time_hours,
			COALESCE(SUM(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0) AS total_active_hours,
			CASE WHEN EXTRACT(EPOCH FROM (?::timestamptz - ?::timestamptz)) > 0 THEN
				100.0 * COALESCE(SUM(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0)
				/ (EXTRACT(EPOCH FROM (?::timestamptz - ?::timestamptz)) / 3600.0)
			ELSE 0 END AS utilization_pct
		FROM equipment_master em
		LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
		LEFT JOIN work_order wo ON wo.equipment_id = em.id
			AND wo.actual_start >= ?::timestamptz AND wo.actual_end <= ?::timestamptz
			AND wo.status = 'completed'
		WHERE em.id = ?
		GROUP BY em.id, ec.class_name, em.operational_status`

	var s scan
	if err := r.data.db.WithContext(ctx).Raw(sqlStr,
		req.To, req.From, req.To, req.From, req.From, req.To, req.EquipmentID,
	).Scan(&s).Error; err != nil {
		return nil, err
	}
	return &biz.MachinePerformanceResult{
		EquipmentID:        s.EquipmentID,
		EquipmentClassName: s.EquipmentClassName,
		OperationalStatus:  s.OperationalStatus,
		TotalWorkOrders:    s.TotalWorkOrders,
		TotalUnitsProduced: s.TotalUnitsProduced,
		AvgCycleTimeHours:  s.AvgCycleTimeHours,
		TotalActiveHours:   s.TotalActiveHours,
		UtilizationPct:     s.UtilizationPct,
	}, nil
}

func (r *traceabilityRepo) GetMachineEventSummary(ctx context.Context, equipmentID string, from, to time.Time) ([]*biz.MachineEventSummary, error) {
	type scan struct {
		EventType string `gorm:"column:event_type"`
		Count     int32  `gorm:"column:count"`
	}
	sqlStr := `SELECT event_type, COUNT(*) AS count FROM traceability_log
		WHERE equipment_id = ? AND event_time >= ?::timestamptz AND event_time <= ?::timestamptz
		GROUP BY event_type ORDER BY count DESC`
	var rows []scan
	if err := r.data.db.WithContext(ctx).Raw(sqlStr, equipmentID, from, to).Scan(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]*biz.MachineEventSummary, len(rows))
	for i, x := range rows {
		list[i] = &biz.MachineEventSummary{EventType: x.EventType, Count: x.Count}
	}
	return list, nil
}

func (r *traceabilityRepo) CompareMachinePerformance(ctx context.Context, req *biz.MachineComparisonRequest) ([]*biz.MachineMetrics, error) {
	if len(req.EquipmentIDs) == 0 {
		return []*biz.MachineMetrics{}, nil
	}

	placeholders := make([]string, len(req.EquipmentIDs))
	for i := range req.EquipmentIDs {
		placeholders[i] = "?"
	}
	args := make([]interface{}, 0, 8+len(req.EquipmentIDs))
	args = append(args, req.To, req.From, req.To, req.From, req.From, req.To, req.From, req.To)
	for _, id := range req.EquipmentIDs {
		args = append(args, id)
	}

	inClause := "(" + strings.Join(placeholders, ",") + ")"
	sqlStr := fmt.Sprintf(`
		SELECT em.id AS equipment_id, ec.class_name AS equipment_class_name, em.operational_status,
			COUNT(DISTINCT wo.work_order_id) AS total_work_orders,
			COALESCE(SUM(wo.actual_quantity), 0) AS total_units_produced,
			COALESCE(AVG(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0) AS avg_cycle_time_hours,
			CASE WHEN EXTRACT(EPOCH FROM (?::timestamptz - ?::timestamptz)) > 0 THEN
				100.0 * COALESCE(SUM(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0)
				/ (EXTRACT(EPOCH FROM (?::timestamptz - ?::timestamptz)) / 3600.0)
			ELSE 0 END AS utilization_pct,
			COALESCE(err.error_count, 0) AS error_count
		FROM equipment_master em
		LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
		LEFT JOIN work_order wo ON wo.equipment_id = em.id
			AND wo.actual_start >= ?::timestamptz AND wo.actual_end <= ?::timestamptz
			AND wo.status = 'completed'
		LEFT JOIN (
			SELECT equipment_id, COUNT(*) AS error_count FROM traceability_log
			WHERE event_time >= ?::timestamptz AND event_time <= ?::timestamptz
				AND event_type LIKE 'error%%'
			GROUP BY equipment_id
		) err ON err.equipment_id = em.id
		WHERE em.id IN %s
		GROUP BY em.id, ec.class_name, em.operational_status, err.error_count
		ORDER BY total_units_produced DESC`, inClause)

	type scan struct {
		EquipmentID        string  `gorm:"column:equipment_id"`
		EquipmentClassName string  `gorm:"column:equipment_class_name"`
		OperationalStatus  string  `gorm:"column:operational_status"`
		TotalWorkOrders    int32   `gorm:"column:total_work_orders"`
		TotalUnitsProduced float64 `gorm:"column:total_units_produced"`
		AvgCycleTimeHours  float64 `gorm:"column:avg_cycle_time_hours"`
		UtilizationPct     float64 `gorm:"column:utilization_pct"`
		ErrorCount         int32   `gorm:"column:error_count"`
	}
	var rows []scan
	if err := r.data.db.WithContext(ctx).Raw(sqlStr, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]*biz.MachineMetrics, len(rows))
	for i, x := range rows {
		list[i] = &biz.MachineMetrics{
			EquipmentID:        x.EquipmentID,
			EquipmentClassName: x.EquipmentClassName,
			OperationalStatus:  x.OperationalStatus,
			TotalWorkOrders:    x.TotalWorkOrders,
			TotalUnitsProduced: x.TotalUnitsProduced,
			AvgCycleTimeHours:  x.AvgCycleTimeHours,
			UtilizationPct:     x.UtilizationPct,
			ErrorCount:         x.ErrorCount,
		}
	}
	return list, nil
}

func (r *traceabilityRepo) GetProductionTrends(ctx context.Context, req *biz.ProductionTrendsRequest) ([]*biz.ProductionTrendPoint, error) {
	truncFunc := "day"
	switch req.Granularity {
	case "weekly":
		truncFunc = "week"
	case "monthly":
		truncFunc = "month"
	}

	args := []interface{}{req.From, req.To}
	filterClause := ""
	if req.EquipmentID != "" {
		filterClause = " AND equipment_id = ?"
		args = append(args, req.EquipmentID)
	}

	sqlStr := fmt.Sprintf(`
		SELECT DATE_TRUNC('%s', actual_end) AS period, equipment_id,
			COALESCE(SUM(actual_quantity), 0) AS units_produced,
			COUNT(*) AS work_orders_completed,
			COALESCE(AVG(EXTRACT(EPOCH FROM (actual_end - actual_start)) / 3600.0), 0) AS avg_cycle_time_hours
		FROM work_order
		WHERE actual_end >= ?::timestamptz AND actual_end <= ?::timestamptz
			AND status = 'completed'%s
		GROUP BY DATE_TRUNC('%s', actual_end), equipment_id
		ORDER BY period, equipment_id`, truncFunc, filterClause, truncFunc)

	type scan struct {
		Period              string  `gorm:"column:period"`
		EquipmentID         string  `gorm:"column:equipment_id"`
		UnitsProduced       float64 `gorm:"column:units_produced"`
		WorkOrdersCompleted int32   `gorm:"column:work_orders_completed"`
		AvgCycleTimeHours   float64 `gorm:"column:avg_cycle_time_hours"`
	}
	var rows []scan
	if err := r.data.db.WithContext(ctx).Raw(sqlStr, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	list := make([]*biz.ProductionTrendPoint, len(rows))
	for i, x := range rows {
		list[i] = &biz.ProductionTrendPoint{
			Period:              x.Period,
			EquipmentID:         x.EquipmentID,
			UnitsProduced:       x.UnitsProduced,
			WorkOrdersCompleted: x.WorkOrdersCompleted,
			AvgCycleTimeHours:   x.AvgCycleTimeHours,
		}
	}
	return list, nil
}

func (r *traceabilityRepo) GetDashboardSummary(ctx context.Context, req *biz.DashboardSummaryRequest) (*biz.DashboardSummary, error) {
	// 1. Production KPIs
	type prodScan struct {
		TotalUnitsProduced       float64 `gorm:"column:total_units_produced"`
		TotalWorkOrdersCompleted int32   `gorm:"column:total_work_orders_completed"`
		ActiveEquipmentCount     int32   `gorm:"column:active_equipment_count"`
	}
	var prod prodScan
	sqlProd := `SELECT
		COALESCE(SUM(actual_quantity), 0) AS total_units_produced,
		COUNT(CASE WHEN status = 'completed' THEN 1 END) AS total_work_orders_completed,
		COUNT(DISTINCT CASE WHEN status IN ('completed','in_progress') THEN equipment_id END) AS active_equipment_count
		FROM work_order
		WHERE created_at >= ?::timestamptz AND created_at <= ?::timestamptz`
	if err := r.data.db.WithContext(ctx).Raw(sqlProd, req.From, req.To).Scan(&prod).Error; err != nil {
		return nil, err
	}

	// 2. Quality KPIs
	type qualScan struct {
		LotsReleased    int32 `gorm:"column:lots_released"`
		LotsQuarantined int32 `gorm:"column:lots_quarantined"`
	}
	var qual qualScan
	sqlQual := `SELECT
		COUNT(CASE WHEN status = 'released' THEN 1 END) AS lots_released,
		COUNT(CASE WHEN status = 'quarantined' THEN 1 END) AS lots_quarantined
		FROM material_lot
		WHERE created_at >= ?::timestamptz AND created_at <= ?::timestamptz`
	if err := r.data.db.WithContext(ctx).Raw(sqlQual, req.From, req.To).Scan(&qual).Error; err != nil {
		return nil, err
	}

	// 3. Error events
	type errScan struct {
		TotalErrors int32 `gorm:"column:total_errors"`
	}
	var errs errScan
	sqlErr := `SELECT COUNT(*) AS total_errors FROM traceability_log
		WHERE event_time >= ?::timestamptz AND event_time <= ?::timestamptz
		AND event_type LIKE 'error%'`
	if err := r.data.db.WithContext(ctx).Raw(sqlErr, req.From, req.To).Scan(&errs).Error; err != nil {
		return nil, err
	}

	// 4. Top performers (top 5 by units produced)
	type topScan struct {
		EquipmentID        string  `gorm:"column:equipment_id"`
		EquipmentClassName string  `gorm:"column:equipment_class_name"`
		OperationalStatus  string  `gorm:"column:operational_status"`
		TotalWorkOrders    int32   `gorm:"column:total_work_orders"`
		TotalUnitsProduced float64 `gorm:"column:total_units_produced"`
		AvgCycleTimeHours  float64 `gorm:"column:avg_cycle_time_hours"`
		UtilizationPct     float64 `gorm:"column:utilization_pct"`
	}
	sqlTop := `
		SELECT em.id AS equipment_id, ec.class_name AS equipment_class_name, em.operational_status,
			COUNT(DISTINCT wo.work_order_id) AS total_work_orders,
			COALESCE(SUM(wo.actual_quantity), 0) AS total_units_produced,
			COALESCE(AVG(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0) AS avg_cycle_time_hours,
			CASE WHEN EXTRACT(EPOCH FROM (?::timestamptz - ?::timestamptz)) > 0 THEN
				100.0 * COALESCE(SUM(EXTRACT(EPOCH FROM (wo.actual_end - wo.actual_start)) / 3600.0), 0)
				/ (EXTRACT(EPOCH FROM (?::timestamptz - ?::timestamptz)) / 3600.0)
			ELSE 0 END AS utilization_pct
		FROM equipment_master em
		LEFT JOIN equipment_class ec ON ec.id = em.equipment_class_id
		LEFT JOIN work_order wo ON wo.equipment_id = em.id
			AND wo.actual_start >= ?::timestamptz AND wo.actual_end <= ?::timestamptz
			AND wo.status = 'completed'
		GROUP BY em.id, ec.class_name, em.operational_status
		ORDER BY total_units_produced DESC
		LIMIT 5`
	var topRows []topScan
	if err := r.data.db.WithContext(ctx).Raw(sqlTop, req.To, req.From, req.To, req.From, req.From, req.To).Scan(&topRows).Error; err != nil {
		return nil, err
	}

	topPerformers := make([]*biz.MachineMetrics, len(topRows))
	for i, x := range topRows {
		topPerformers[i] = &biz.MachineMetrics{
			EquipmentID:        x.EquipmentID,
			EquipmentClassName: x.EquipmentClassName,
			OperationalStatus:  x.OperationalStatus,
			TotalWorkOrders:    x.TotalWorkOrders,
			TotalUnitsProduced: x.TotalUnitsProduced,
			AvgCycleTimeHours:  x.AvgCycleTimeHours,
			UtilizationPct:     x.UtilizationPct,
		}
	}

	qualityRate := float64(0)
	total := qual.LotsReleased + qual.LotsQuarantined
	if total > 0 {
		qualityRate = 100.0 * float64(qual.LotsReleased) / float64(total)
	}

	return &biz.DashboardSummary{
		TotalUnitsProduced:       prod.TotalUnitsProduced,
		TotalWorkOrdersCompleted: prod.TotalWorkOrdersCompleted,
		ActiveEquipmentCount:     prod.ActiveEquipmentCount,
		LotsReleased:             qual.LotsReleased,
		LotsQuarantined:          qual.LotsQuarantined,
		QualityRatePct:           qualityRate,
		TotalErrorEvents:         errs.TotalErrors,
		TopPerformers:            topPerformers,
	}, nil
}
