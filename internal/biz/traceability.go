package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// ==========================================
// 1. Domain Models (Go Structs for Business Logic)
// ==========================================

type Enterprise struct {
	ID          int32
	Name        string
	Description string
}

type Site struct {
	ID           int32
	EnterpriseID int32
	Name         string
	Location     string
	Description  string
}

type Area struct {
	ID          int32
	SiteID      int32
	Name        string
	Description string
}

type ProductionLine struct {
	ID          int32
	AreaID      int32
	Name        string
	Description string
}

type ProductionUnit struct {
	ID               int32
	ProductionLineID int32
	Name             string
	Description      string
}

type EquipmentClass struct {
	ID          int32
	ClassName   string
	Version     string
	Description string
}

type EquipmentMaster struct {
	ID                string // String ID like "ROBOT-99"
	PhysicalAssetID   string
	ProductionUnitID  int32
	EquipmentClassID  int32
	OperationalStatus string
	ParentEquipmentID string // For Sub-Components
}

type EquipmentCapability struct {
	ID             int32
	EquipmentID    string
	CapabilityName string
	Value          string
	UOM            string
	Description    string
}

type EquipmentProperty struct {
	ID           int32
	EquipmentID  string
	PropertyName string
	CurrentValue string
	LastUpdated  time.Time
}

type TraceabilityLog struct {
	ID            int32
	EquipmentID   string
	EventType     string
	WorkOrderID   string
	MaterialLotID string
	OperatorID    string
	EventTime     time.Time
}

type MaterialDefinition struct {
	ID            int32
	Name          string
	MaterialType  string
	UnitOfMeasure string
	Description   string
}

type MaterialLot struct {
	LotID                string
	MaterialDefinitionID int32
	Quantity             float64
	UnitOfMeasure        string
	Status               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type MaterialLotSummary struct {
	LotID         string
	Status        string
	Quantity      float64
	UnitOfMeasure string
	CreatedAt     time.Time
	MaterialName  string
	MaterialType  string
}

type MaterialLotDetail struct {
	MaterialLot
	MaterialName        string
	MaterialType        string
	MaterialDescription string
}

type Operator struct {
	OperatorID string
	Name       string
	Role       string
	Shift      string
	Status     string
}

type WorkOrder struct {
	WorkOrderID     string
	Description     string
	Status          string
	EquipmentID     string
	OperatorID      string
	OutputLotID     string
	PlannedQuantity float64
	ActualQuantity  *float64
	UnitOfMeasure   string
	PlannedStart    time.Time
	PlannedEnd      time.Time
	ActualStart     *time.Time
	ActualEnd       *time.Time
}

type WorkOrderSummary struct {
	WorkOrderID  string
	Description  string
	Status       string
	EquipmentID  string
	OperatorID   string
	OutputLotID  string
	ActualStart  *time.Time
	ActualEnd    *time.Time
}

type WorkOrderInputLot struct {
	LotID            string
	MaterialName     string
	MaterialType     string
	QuantityConsumed float64
	UnitOfMeasure    string
}

type WorkOrderDetail struct {
	WorkOrder
	EquipmentClassName string
	OperatorName       string
	OperatorShift      string
	InputLots          []*WorkOrderInputLot
}

type LotGenealogy struct {
	ID               int32
	ParentLotID      string
	ChildLotID       string
	WorkOrderID      string
	EquipmentID      string
	QuantityConsumed float64
	QuantityProduced float64
	EventTime        time.Time
}

type TraceNode struct {
	Depth              int32
	LotID              string
	RelatedLotID       string
	MaterialName       string
	MaterialType       string
	LotStatus          string
	Quantity           float64
	UnitOfMeasure      string
	QuantityUsed       float64
	WorkOrderID        string
	EquipmentID        string
	EquipmentClassName string
	OperatorID         string
	OperatorName       string
	EventTime          *time.Time
}

type GenealogyNode struct {
	LotID         string
	MaterialName  string
	MaterialType  string
	Status        string
	Quantity      float64
	UnitOfMeasure string
}

type GenealogyEdge struct {
	SourceLotID      string
	TargetLotID      string
	WorkOrderID      string
	EquipmentID      string
	QuantityConsumed float64
	QuantityProduced float64
	EventTime        *time.Time
}

type EquipmentParameterSummary struct {
	EquipmentID        string
	EquipmentClassName string
	ParameterName      string
	UnitOfMeasure      string
	MinValue           float64
	MaxValue           float64
	AvgValue           float64
	ReadingCount       int32
	ProcessStart       *time.Time
	ProcessEnd         *time.Time
}

type TelemetryReading struct {
	EquipmentID   string
	ParameterName string
	Value         float64
	UnitOfMeasure string
	RecordedAt    *time.Time
}

// ==========================================
// 2. Repository Interface (Contract for Data Layer)
// ==========================================

type TraceabilityRepo interface {
	// Enterprise
	CreateEnterprise(ctx context.Context, e *Enterprise) (int32, error)
	ListEnterprises(ctx context.Context) ([]*Enterprise, error)
	UpdateEnterprise(ctx context.Context, e *Enterprise) error
	DeleteEnterprise(ctx context.Context, id int32) error

	// Site
	CreateSite(ctx context.Context, s *Site) (int32, error)
	ListSites(ctx context.Context, enterpriseID int32) ([]*Site, error)
	UpdateSite(ctx context.Context, s *Site) error
	DeleteSite(ctx context.Context, id int32) error

	// Area
	CreateArea(ctx context.Context, a *Area) (int32, error)
	ListAreas(ctx context.Context, siteID int32) ([]*Area, error)
	UpdateArea(ctx context.Context, a *Area) error
	DeleteArea(ctx context.Context, id int32) error

	// Line
	CreateLine(ctx context.Context, l *ProductionLine) (int32, error)
	ListLines(ctx context.Context, areaID int32) ([]*ProductionLine, error)
	UpdateLine(ctx context.Context, l *ProductionLine) error
	DeleteLine(ctx context.Context, id int32) error

	// Unit (Slot)
	CreateProductionUnit(ctx context.Context, u *ProductionUnit) (int32, error)
	ListProductionUnits(ctx context.Context, lineID int32) ([]*ProductionUnit, error)
	UpdateProductionUnit(ctx context.Context, u *ProductionUnit) error
	DeleteProductionUnit(ctx context.Context, id int32) error

	// Class
	CreateEquipmentClass(ctx context.Context, c *EquipmentClass) (int32, error)
	ListEquipmentClasses(ctx context.Context) ([]*EquipmentClass, error)
	UpdateEquipmentClass(ctx context.Context, c *EquipmentClass) error
	DeleteEquipmentClass(ctx context.Context, id int32) error

	// Equipment
	RegisterEquipment(ctx context.Context, e *EquipmentMaster) (string, error)
	ListEquipment(ctx context.Context, lineID int32, parentID string) ([]*EquipmentMaster, error)
	UpdateEquipment(ctx context.Context, e *EquipmentMaster) error
	DeleteEquipment(ctx context.Context, id string) error

	// Capability
	AddCapability(ctx context.Context, c *EquipmentCapability) (int32, error)
	ListCapabilities(ctx context.Context, equipmentID string) ([]*EquipmentCapability, error)
	UpdateCapability(ctx context.Context, c *EquipmentCapability) error
	DeleteCapability(ctx context.Context, equipmentID string, id int32) error

	// Property
	SetProperty(ctx context.Context, p *EquipmentProperty) (int32, error)
	ListProperties(ctx context.Context, equipmentID string) ([]*EquipmentProperty, error)
	UpdateProperty(ctx context.Context, p *EquipmentProperty) error
	DeleteProperty(ctx context.Context, equipmentID string, id int32) error

	// Logs
	LogEvent(ctx context.Context, l *TraceabilityLog) (int32, error)
	ListLogs(ctx context.Context, workOrderID string) ([]*TraceabilityLog, error)

	// Material Definition
	CreateMaterialDefinition(ctx context.Context, md *MaterialDefinition) (int32, error)
	ListMaterialDefinitions(ctx context.Context, mType string) ([]*MaterialDefinition, error)

	// Material Lot
	CreateMaterialLot(ctx context.Context, ml *MaterialLot) (string, error)
	GetMaterialLot(ctx context.Context, id string) (*MaterialLotDetail, error)
	ListMaterialLots(ctx context.Context, status, mType string) ([]*MaterialLotSummary, error)
	UpdateMaterialLotStatus(ctx context.Context, id, status string) error

	// Operator
	CreateOperator(ctx context.Context, op *Operator) (string, error)
	ListOperators(ctx context.Context, shift, status string) ([]*Operator, error)
	UpdateOperator(ctx context.Context, op *Operator) error

	// Work Order
	CreateWorkOrder(ctx context.Context, wo *WorkOrder) (string, error)
	GetWorkOrder(ctx context.Context, id string) (*WorkOrderDetail, error)
	ListWorkOrders(ctx context.Context, status, eqID string) ([]*WorkOrderSummary, error)
	UpdateWorkOrderStatus(ctx context.Context, wo *WorkOrder) error

	// Genealogy
	RegisterGenealogyLink(ctx context.Context, lg *LotGenealogy) (int32, error)

	// Trace
	TraceBackward(ctx context.Context, lotID string) ([]*TraceNode, error)
	TraceForward(ctx context.Context, lotID string) ([]*TraceNode, error)
	TraceFullGenealogyNodes(ctx context.Context, lotID string) ([]*GenealogyNode, error)
	TraceFullGenealogyEdges(ctx context.Context, lotID string) ([]*GenealogyEdge, error)
	GetEquipmentProcessHistorySummary(ctx context.Context, lotID string) ([]*EquipmentParameterSummary, error)
	GetEquipmentProcessHistoryReadings(ctx context.Context, lotID string) ([]*TelemetryReading, error)
}

// ==========================================
// 3. Usecase Implementation (Business Logic)
// ==========================================

type TraceabilityUsecase struct {
	repo TraceabilityRepo
	log  *log.Helper
}

func NewTraceabilityUsecase(repo TraceabilityRepo, logger log.Logger) *TraceabilityUsecase {
	return &TraceabilityUsecase{repo: repo, log: log.NewHelper(logger)}
}

// --- Methods (Pass-through to Repo for now) ---

func (uc *TraceabilityUsecase) CreateEnterprise(ctx context.Context, e *Enterprise) (int32, error) {
	return uc.repo.CreateEnterprise(ctx, e)
}
func (uc *TraceabilityUsecase) ListEnterprises(ctx context.Context) ([]*Enterprise, error) {
	return uc.repo.ListEnterprises(ctx)
}
func (uc *TraceabilityUsecase) UpdateEnterprise(ctx context.Context, e *Enterprise) error {
	return uc.repo.UpdateEnterprise(ctx, e)
}
func (uc *TraceabilityUsecase) DeleteEnterprise(ctx context.Context, id int32) error {
	return uc.repo.DeleteEnterprise(ctx, id)
}

func (uc *TraceabilityUsecase) CreateSite(ctx context.Context, s *Site) (int32, error) {
	return uc.repo.CreateSite(ctx, s)
}
func (uc *TraceabilityUsecase) ListSites(ctx context.Context, eid int32) ([]*Site, error) {
	return uc.repo.ListSites(ctx, eid)
}
func (uc *TraceabilityUsecase) UpdateSite(ctx context.Context, s *Site) error {
	return uc.repo.UpdateSite(ctx, s)
}
func (uc *TraceabilityUsecase) DeleteSite(ctx context.Context, id int32) error {
	return uc.repo.DeleteSite(ctx, id)
}

func (uc *TraceabilityUsecase) CreateArea(ctx context.Context, a *Area) (int32, error) {
	return uc.repo.CreateArea(ctx, a)
}
func (uc *TraceabilityUsecase) ListAreas(ctx context.Context, sid int32) ([]*Area, error) {
	return uc.repo.ListAreas(ctx, sid)
}
func (uc *TraceabilityUsecase) UpdateArea(ctx context.Context, a *Area) error {
	return uc.repo.UpdateArea(ctx, a)
}
func (uc *TraceabilityUsecase) DeleteArea(ctx context.Context, id int32) error {
	return uc.repo.DeleteArea(ctx, id)
}

func (uc *TraceabilityUsecase) CreateLine(ctx context.Context, l *ProductionLine) (int32, error) {
	return uc.repo.CreateLine(ctx, l)
}
func (uc *TraceabilityUsecase) ListLines(ctx context.Context, aid int32) ([]*ProductionLine, error) {
	return uc.repo.ListLines(ctx, aid)
}
func (uc *TraceabilityUsecase) UpdateLine(ctx context.Context, l *ProductionLine) error {
	return uc.repo.UpdateLine(ctx, l)
}
func (uc *TraceabilityUsecase) DeleteLine(ctx context.Context, id int32) error {
	return uc.repo.DeleteLine(ctx, id)
}

func (uc *TraceabilityUsecase) CreateProductionUnit(ctx context.Context, u *ProductionUnit) (int32, error) {
	return uc.repo.CreateProductionUnit(ctx, u)
}
func (uc *TraceabilityUsecase) ListProductionUnits(ctx context.Context, lid int32) ([]*ProductionUnit, error) {
	return uc.repo.ListProductionUnits(ctx, lid)
}
func (uc *TraceabilityUsecase) UpdateProductionUnit(ctx context.Context, u *ProductionUnit) error {
	return uc.repo.UpdateProductionUnit(ctx, u)
}
func (uc *TraceabilityUsecase) DeleteProductionUnit(ctx context.Context, id int32) error {
	return uc.repo.DeleteProductionUnit(ctx, id)
}

func (uc *TraceabilityUsecase) CreateEquipmentClass(ctx context.Context, c *EquipmentClass) (int32, error) {
	return uc.repo.CreateEquipmentClass(ctx, c)
}
func (uc *TraceabilityUsecase) ListEquipmentClasses(ctx context.Context) ([]*EquipmentClass, error) {
	return uc.repo.ListEquipmentClasses(ctx)
}
func (uc *TraceabilityUsecase) UpdateEquipmentClass(ctx context.Context, c *EquipmentClass) error {
	return uc.repo.UpdateEquipmentClass(ctx, c)
}
func (uc *TraceabilityUsecase) DeleteEquipmentClass(ctx context.Context, id int32) error {
	return uc.repo.DeleteEquipmentClass(ctx, id)
}

func (uc *TraceabilityUsecase) RegisterEquipment(ctx context.Context, e *EquipmentMaster) (string, error) {
	return uc.repo.RegisterEquipment(ctx, e)
}
func (uc *TraceabilityUsecase) ListEquipment(ctx context.Context, lid int32, pid string) ([]*EquipmentMaster, error) {
	return uc.repo.ListEquipment(ctx, lid, pid)
}
func (uc *TraceabilityUsecase) UpdateEquipment(ctx context.Context, e *EquipmentMaster) error {
	return uc.repo.UpdateEquipment(ctx, e)
}
func (uc *TraceabilityUsecase) DeleteEquipment(ctx context.Context, id string) error {
	return uc.repo.DeleteEquipment(ctx, id)
}

func (uc *TraceabilityUsecase) AddCapability(ctx context.Context, c *EquipmentCapability) (int32, error) {
	return uc.repo.AddCapability(ctx, c)
}
func (uc *TraceabilityUsecase) ListCapabilities(ctx context.Context, eid string) ([]*EquipmentCapability, error) {
	return uc.repo.ListCapabilities(ctx, eid)
}
func (uc *TraceabilityUsecase) UpdateCapability(ctx context.Context, c *EquipmentCapability) error {
	return uc.repo.UpdateCapability(ctx, c)
}
func (uc *TraceabilityUsecase) DeleteCapability(ctx context.Context, eid string, id int32) error {
	return uc.repo.DeleteCapability(ctx, eid, id)
}

func (uc *TraceabilityUsecase) SetProperty(ctx context.Context, p *EquipmentProperty) (int32, error) {
	return uc.repo.SetProperty(ctx, p)
}
func (uc *TraceabilityUsecase) ListProperties(ctx context.Context, eid string) ([]*EquipmentProperty, error) {
	return uc.repo.ListProperties(ctx, eid)
}
func (uc *TraceabilityUsecase) UpdateProperty(ctx context.Context, p *EquipmentProperty) error {
	return uc.repo.UpdateProperty(ctx, p)
}
func (uc *TraceabilityUsecase) DeleteProperty(ctx context.Context, eid string, id int32) error {
	return uc.repo.DeleteProperty(ctx, eid, id)
}

func (uc *TraceabilityUsecase) LogEvent(ctx context.Context, l *TraceabilityLog) (int32, error) {
	return uc.repo.LogEvent(ctx, l)
}
func (uc *TraceabilityUsecase) ListLogs(ctx context.Context, wid string) ([]*TraceabilityLog, error) {
	return uc.repo.ListLogs(ctx, wid)
}

func (uc *TraceabilityUsecase) CreateMaterialDefinition(ctx context.Context, md *MaterialDefinition) (int32, error) {
	return uc.repo.CreateMaterialDefinition(ctx, md)
}
func (uc *TraceabilityUsecase) ListMaterialDefinitions(ctx context.Context, mType string) ([]*MaterialDefinition, error) {
	return uc.repo.ListMaterialDefinitions(ctx, mType)
}

func (uc *TraceabilityUsecase) CreateMaterialLot(ctx context.Context, ml *MaterialLot) (string, error) {
	return uc.repo.CreateMaterialLot(ctx, ml)
}
func (uc *TraceabilityUsecase) GetMaterialLot(ctx context.Context, id string) (*MaterialLotDetail, error) {
	return uc.repo.GetMaterialLot(ctx, id)
}
func (uc *TraceabilityUsecase) ListMaterialLots(ctx context.Context, status, mType string) ([]*MaterialLotSummary, error) {
	return uc.repo.ListMaterialLots(ctx, status, mType)
}
func (uc *TraceabilityUsecase) UpdateMaterialLotStatus(ctx context.Context, id, status string) error {
	return uc.repo.UpdateMaterialLotStatus(ctx, id, status)
}

func (uc *TraceabilityUsecase) CreateOperator(ctx context.Context, op *Operator) (string, error) {
	return uc.repo.CreateOperator(ctx, op)
}
func (uc *TraceabilityUsecase) ListOperators(ctx context.Context, shift, status string) ([]*Operator, error) {
	return uc.repo.ListOperators(ctx, shift, status)
}
func (uc *TraceabilityUsecase) UpdateOperator(ctx context.Context, op *Operator) error {
	return uc.repo.UpdateOperator(ctx, op)
}

func (uc *TraceabilityUsecase) CreateWorkOrder(ctx context.Context, wo *WorkOrder) (string, error) {
	return uc.repo.CreateWorkOrder(ctx, wo)
}
func (uc *TraceabilityUsecase) GetWorkOrder(ctx context.Context, id string) (*WorkOrderDetail, error) {
	return uc.repo.GetWorkOrder(ctx, id)
}
func (uc *TraceabilityUsecase) ListWorkOrders(ctx context.Context, status, eqID string) ([]*WorkOrderSummary, error) {
	return uc.repo.ListWorkOrders(ctx, status, eqID)
}
func (uc *TraceabilityUsecase) UpdateWorkOrderStatus(ctx context.Context, wo *WorkOrder) error {
	return uc.repo.UpdateWorkOrderStatus(ctx, wo)
}

func (uc *TraceabilityUsecase) RegisterGenealogyLink(ctx context.Context, lg *LotGenealogy) (int32, error) {
	return uc.repo.RegisterGenealogyLink(ctx, lg)
}

func (uc *TraceabilityUsecase) TraceBackward(ctx context.Context, lotID string) ([]*TraceNode, error) {
	return uc.repo.TraceBackward(ctx, lotID)
}
func (uc *TraceabilityUsecase) TraceForward(ctx context.Context, lotID string) ([]*TraceNode, error) {
	return uc.repo.TraceForward(ctx, lotID)
}
func (uc *TraceabilityUsecase) TraceFullGenealogy(ctx context.Context, lotID string) ([]*GenealogyNode, []*GenealogyEdge, error) {
	nodes, err := uc.repo.TraceFullGenealogyNodes(ctx, lotID)
	if err != nil {
		return nil, nil, err
	}
	edges, err := uc.repo.TraceFullGenealogyEdges(ctx, lotID)
	if err != nil {
		return nil, nil, err
	}
	return nodes, edges, nil
}
func (uc *TraceabilityUsecase) GetEquipmentProcessHistory(ctx context.Context, lotID string) ([]*EquipmentParameterSummary, []*TelemetryReading, error) {
	summary, err := uc.repo.GetEquipmentProcessHistorySummary(ctx, lotID)
	if err != nil {
		return nil, nil, err
	}
	readings, err := uc.repo.GetEquipmentProcessHistoryReadings(ctx, lotID)
	if err != nil {
		return nil, nil, err
	}
	return summary, readings, nil
}
