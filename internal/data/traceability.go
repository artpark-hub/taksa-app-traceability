package data

import (
	"context"
	"time"

	"traceability/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// ==========================================
// 1. ORM Definitions (Database Tables)
// ==========================================

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
	ParentEquipmentID *string // Pointer because it can be null
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

// ==========================================
// 2. Repository Implementation
// ==========================================

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
	r.data.db.WithContext(ctx).Find(&dbList)
	var list []*biz.Enterprise
	for _, x := range dbList {
		list = append(list, &biz.Enterprise{ID: x.ID, Name: x.Name, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateEnterprise(ctx context.Context, e *biz.Enterprise) error {
	return r.data.db.WithContext(ctx).Model(&EnterpriseORM{}).Where("id = ?", e.ID).Updates(EnterpriseORM{Name: e.Name}).Error
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
	query.Find(&dbList)
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
	query.Find(&dbList)
	var list []*biz.Area
	for _, x := range dbList {
		list = append(list, &biz.Area{ID: x.ID, SiteID: x.SiteID, Name: x.Name, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateArea(ctx context.Context, a *biz.Area) error {
	return r.data.db.WithContext(ctx).Model(&AreaORM{}).Where("id = ?", a.ID).Updates(AreaORM{Name: a.Name}).Error
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
	query.Find(&dbList)
	var list []*biz.ProductionLine
	for _, x := range dbList {
		list = append(list, &biz.ProductionLine{ID: x.ID, AreaID: x.AreaID, Name: x.Name, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateLine(ctx context.Context, l *biz.ProductionLine) error {
	return r.data.db.WithContext(ctx).Model(&LineORM{}).Where("id = ?", l.ID).Updates(LineORM{Name: l.Name}).Error
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
	query.Find(&dbList)
	var list []*biz.ProductionUnit
	for _, x := range dbList {
		list = append(list, &biz.ProductionUnit{ID: x.ID, ProductionLineID: x.ProductionLineID, Name: x.Name, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateProductionUnit(ctx context.Context, u *biz.ProductionUnit) error {
	return r.data.db.WithContext(ctx).Model(&UnitORM{}).Where("id = ?", u.ID).Updates(UnitORM{Name: u.Name}).Error
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
	r.data.db.WithContext(ctx).Find(&dbList)
	var list []*biz.EquipmentClass
	for _, x := range dbList {
		list = append(list, &biz.EquipmentClass{ID: x.ID, ClassName: x.ClassName, Version: x.Version, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateEquipmentClass(ctx context.Context, c *biz.EquipmentClass) error {
	return r.data.db.WithContext(ctx).Model(&ClassORM{}).Where("id = ?", c.ID).Updates(ClassORM{Version: c.Version}).Error
}
func (r *traceabilityRepo) DeleteEquipmentClass(ctx context.Context, id int32) error {
	return r.data.db.WithContext(ctx).Delete(&ClassORM{}, id).Error
}

// Equipment
func (r *traceabilityRepo) RegisterEquipment(ctx context.Context, e *biz.EquipmentMaster) (string, error) {
	// Handle nil parent ID logic
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

	// Join with ProductionUnit to filter by LineID if provided
	if lid > 0 {
		query = query.Joins("JOIN production_unit ON production_unit.id = equipment_master.production_unit_id").
			Where("production_unit.production_line_id = ?", lid)
	}
	// Filter by Parent ID
	if pid != "" {
		query = query.Where("parent_equipment_id = ?", pid)
	}

	query.Find(&dbList)
	var list []*biz.EquipmentMaster
	for _, x := range dbList {
		// Handle potential nil ParentID
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
	r.data.db.WithContext(ctx).Where("equipment_id = ?", eid).Find(&dbList)
	var list []*biz.EquipmentCapability
	for _, x := range dbList {
		list = append(list, &biz.EquipmentCapability{ID: x.ID, EquipmentID: x.EquipmentID, CapabilityName: x.CapabilityName, Value: x.Value, UOM: x.UOM, Description: x.Description})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateCapability(ctx context.Context, c *biz.EquipmentCapability) error {
	return r.data.db.WithContext(ctx).Model(&CapabilityORM{}).Where("id = ?", c.ID).Updates(CapabilityORM{Value: c.Value}).Error
}
func (r *traceabilityRepo) DeleteCapability(ctx context.Context, eid string, id int32) error {
	return r.data.db.WithContext(ctx).Delete(&CapabilityORM{}, id).Error
}

// Property
func (r *traceabilityRepo) SetProperty(ctx context.Context, p *biz.EquipmentProperty) (int32, error) {
	orm := PropertyORM{EquipmentID: p.EquipmentID, PropertyName: p.PropertyName, CurrentValue: p.CurrentValue, LastUpdated: time.Now()}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListProperties(ctx context.Context, eid string) ([]*biz.EquipmentProperty, error) {
	var dbList []PropertyORM
	r.data.db.WithContext(ctx).Where("equipment_id = ?", eid).Find(&dbList)
	var list []*biz.EquipmentProperty
	for _, x := range dbList {
		list = append(list, &biz.EquipmentProperty{ID: x.ID, EquipmentID: x.EquipmentID, PropertyName: x.PropertyName, CurrentValue: x.CurrentValue, LastUpdated: x.LastUpdated})
	}
	return list, nil
}
func (r *traceabilityRepo) UpdateProperty(ctx context.Context, p *biz.EquipmentProperty) error {
	return r.data.db.WithContext(ctx).Model(&PropertyORM{}).Where("id = ?", p.ID).Updates(PropertyORM{CurrentValue: p.CurrentValue, LastUpdated: time.Now()}).Error
}
func (r *traceabilityRepo) DeleteProperty(ctx context.Context, eid string, id int32) error {
	return r.data.db.WithContext(ctx).Delete(&PropertyORM{}, id).Error
}

// Logs
func (r *traceabilityRepo) LogEvent(ctx context.Context, l *biz.TraceabilityLog) (int32, error) {
	orm := LogORM{EquipmentID: l.EquipmentID, EventType: l.EventType, WorkOrderID: l.WorkOrderID, MaterialLotID: l.MaterialLotID, OperatorID: l.OperatorID, EventTime: time.Now()}
	res := r.data.db.WithContext(ctx).Create(&orm)
	return orm.ID, res.Error
}
func (r *traceabilityRepo) ListLogs(ctx context.Context, wid string) ([]*biz.TraceabilityLog, error) {
	var dbList []LogORM
	r.data.db.WithContext(ctx).Where("work_order_id = ?", wid).Order("event_time desc").Find(&dbList)
	var list []*biz.TraceabilityLog
	for _, x := range dbList {
		list = append(list, &biz.TraceabilityLog{ID: x.ID, EquipmentID: x.EquipmentID, EventType: x.EventType, WorkOrderID: x.WorkOrderID, MaterialLotID: x.MaterialLotID, OperatorID: x.OperatorID, EventTime: x.EventTime})
	}
	return list, nil
}
