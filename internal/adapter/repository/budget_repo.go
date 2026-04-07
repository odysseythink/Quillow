package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"gorm.io/gorm"
)

// ---------------------------------------------------------------------------
// BudgetRepository
// ---------------------------------------------------------------------------

type BudgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) FindByID(ctx context.Context, id uint) (*entity.Budget, error) {
	var m model.BudgetModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("budget not found: %w", err)
	}
	return budgetModelToEntity(&m), nil
}

func (r *BudgetRepository) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.Budget, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.BudgetModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	query.Count(&total)

	var models []model.BudgetModel
	if err := query.Order("`order` ASC, name ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	budgets := make([]entity.Budget, len(models))
	for i, m := range models {
		budgets[i] = *budgetModelToEntity(&m)
	}
	return budgets, total, nil
}

func (r *BudgetRepository) ListActive(ctx context.Context, userGroupID uint) ([]entity.Budget, error) {
	query := r.db.WithContext(ctx).Model(&model.BudgetModel{}).Where("deleted_at IS NULL AND active = ?", true)
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var models []model.BudgetModel
	if err := query.Order("`order` ASC, name ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	budgets := make([]entity.Budget, len(models))
	for i, m := range models {
		budgets[i] = *budgetModelToEntity(&m)
	}
	return budgets, nil
}

func (r *BudgetRepository) Create(ctx context.Context, budget *entity.Budget) error {
	m := budgetEntityToModel(budget)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	budget.ID = m.ID
	budget.CreatedAt = m.CreatedAt
	budget.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *BudgetRepository) Update(ctx context.Context, budget *entity.Budget) error {
	m := budgetEntityToModel(budget)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BudgetRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.BudgetModel{}, id).Error
}

func (r *BudgetRepository) GetNotes(ctx context.Context, budgetID uint) (string, error) {
	var note model.NoteModel
	if err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type LIKE ? AND deleted_at IS NULL", budgetID, "%Budget%").First(&note).Error; err != nil {
		return "", nil
	}
	if note.Text != nil {
		return *note.Text, nil
	}
	return "", nil
}

func (r *BudgetRepository) SetNotes(ctx context.Context, budgetID uint, text string) error {
	var existing model.NoteModel
	noteableType := "Quillow\\Models\\Budget"
	err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type = ?", budgetID, noteableType).First(&existing).Error
	if err != nil {
		return r.db.WithContext(ctx).Create(&model.NoteModel{
			NoteableID:   budgetID,
			NoteableType: noteableType,
			Text:         &text,
		}).Error
	}
	existing.Text = &text
	return r.db.WithContext(ctx).Save(&existing).Error
}

// ---------------------------------------------------------------------------
// BudgetLimitRepository
// ---------------------------------------------------------------------------

type BudgetLimitRepository struct {
	db *gorm.DB
}

func NewBudgetLimitRepository(db *gorm.DB) *BudgetLimitRepository {
	return &BudgetLimitRepository{db: db}
}

func (r *BudgetLimitRepository) FindByID(ctx context.Context, id uint) (*entity.BudgetLimit, error) {
	var m model.BudgetLimitModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("budget limit not found: %w", err)
	}
	return budgetLimitModelToEntity(&m), nil
}

func (r *BudgetLimitRepository) ListByBudget(ctx context.Context, budgetID uint) ([]entity.BudgetLimit, error) {
	var models []model.BudgetLimitModel
	if err := r.db.WithContext(ctx).Where("budget_id = ?", budgetID).Order("start_date ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	limits := make([]entity.BudgetLimit, len(models))
	for i, m := range models {
		limits[i] = *budgetLimitModelToEntity(&m)
	}
	return limits, nil
}

func (r *BudgetLimitRepository) ListByPeriod(ctx context.Context, budgetID uint, start, end time.Time) ([]entity.BudgetLimit, error) {
	var models []model.BudgetLimitModel
	if err := r.db.WithContext(ctx).
		Where("budget_id = ? AND start_date >= ? AND end_date <= ?", budgetID, start, end).
		Order("start_date ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	limits := make([]entity.BudgetLimit, len(models))
	for i, m := range models {
		limits[i] = *budgetLimitModelToEntity(&m)
	}
	return limits, nil
}

func (r *BudgetLimitRepository) Create(ctx context.Context, limit *entity.BudgetLimit) error {
	m := budgetLimitEntityToModel(limit)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	limit.ID = m.ID
	limit.CreatedAt = m.CreatedAt
	limit.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *BudgetLimitRepository) Update(ctx context.Context, limit *entity.BudgetLimit) error {
	m := budgetLimitEntityToModel(limit)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *BudgetLimitRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.BudgetLimitModel{}, id).Error
}

// ---------------------------------------------------------------------------
// Conversion helpers
// ---------------------------------------------------------------------------

func budgetModelToEntity(m *model.BudgetModel) *entity.Budget {
	return &entity.Budget{
		ID:          m.ID,
		UserID:      m.UserID,
		UserGroupID: m.UserGroupID,
		Name:        m.Name,
		Active:      m.Active,
		Encrypted:   m.Encrypted,
		Order:       m.Order,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}

func budgetEntityToModel(b *entity.Budget) *model.BudgetModel {
	return &model.BudgetModel{
		ID:          b.ID,
		UserID:      b.UserID,
		UserGroupID: b.UserGroupID,
		Name:        b.Name,
		Active:      b.Active,
		Encrypted:   b.Encrypted,
		Order:       b.Order,
	}
}

func budgetLimitModelToEntity(m *model.BudgetLimitModel) *entity.BudgetLimit {
	bl := &entity.BudgetLimit{
		ID:                    m.ID,
		BudgetID:              m.BudgetID,
		TransactionCurrencyID: m.TransactionCurrencyID,
		StartDate:             m.StartDate,
		EndDate:               m.EndDate,
		Amount:                m.Amount,
		Generated:             m.Generated,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
	}
	if m.Period != nil {
		bl.Period = *m.Period
	}
	return bl
}

func budgetLimitEntityToModel(bl *entity.BudgetLimit) *model.BudgetLimitModel {
	m := &model.BudgetLimitModel{
		ID:                    bl.ID,
		BudgetID:              bl.BudgetID,
		TransactionCurrencyID: bl.TransactionCurrencyID,
		StartDate:             bl.StartDate,
		EndDate:               bl.EndDate,
		Amount:                bl.Amount,
		Generated:             bl.Generated,
	}
	if bl.Period != "" {
		m.Period = &bl.Period
	}
	return m
}
