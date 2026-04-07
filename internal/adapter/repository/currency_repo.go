package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"github.com/anthropics/quillow/internal/entity"
	"gorm.io/gorm"
)

type CurrencyRepository struct {
	db *gorm.DB
}

func NewCurrencyRepository(db *gorm.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

func (r *CurrencyRepository) FindByID(ctx context.Context, id uint) (*entity.TransactionCurrency, error) {
	var m model.TransactionCurrencyModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("currency not found: %w", err)
	}
	return currencyModelToEntity(&m), nil
}

func (r *CurrencyRepository) FindByCode(ctx context.Context, code string) (*entity.TransactionCurrency, error) {
	var m model.TransactionCurrencyModel
	if err := r.db.WithContext(ctx).Where("code = ? AND deleted_at IS NULL", code).First(&m).Error; err != nil {
		return nil, fmt.Errorf("currency not found: %w", err)
	}
	return currencyModelToEntity(&m), nil
}

func (r *CurrencyRepository) List(ctx context.Context, limit, offset int) ([]entity.TransactionCurrency, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&model.TransactionCurrencyModel{}).Where("deleted_at IS NULL").Count(&total)

	var models []model.TransactionCurrencyModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Order("code ASC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	currencies := make([]entity.TransactionCurrency, len(models))
	for i, m := range models {
		currencies[i] = *currencyModelToEntity(&m)
	}
	return currencies, total, nil
}

func (r *CurrencyRepository) Create(ctx context.Context, currency *entity.TransactionCurrency) error {
	m := currencyEntityToModel(currency)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	currency.ID = m.ID
	currency.CreatedAt = m.CreatedAt
	currency.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *CurrencyRepository) Update(ctx context.Context, currency *entity.TransactionCurrency) error {
	m := currencyEntityToModel(currency)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *CurrencyRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.TransactionCurrencyModel{}, id).Error
}

func (r *CurrencyRepository) GetPrimary(ctx context.Context) (*entity.TransactionCurrency, error) {
	var m model.TransactionCurrencyModel
	if err := r.db.WithContext(ctx).Where("enabled = ? AND deleted_at IS NULL", true).Order("id ASC").First(&m).Error; err != nil {
		return nil, fmt.Errorf("no primary currency found: %w", err)
	}
	return currencyModelToEntity(&m), nil
}

func currencyModelToEntity(m *model.TransactionCurrencyModel) *entity.TransactionCurrency {
	return &entity.TransactionCurrency{
		ID:            m.ID,
		Code:          m.Code,
		Name:          m.Name,
		Symbol:        m.Symbol,
		DecimalPlaces: m.DecimalPlaces,
		Enabled:       m.Enabled,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
	}
}

func currencyEntityToModel(c *entity.TransactionCurrency) *model.TransactionCurrencyModel {
	return &model.TransactionCurrencyModel{
		ID:            c.ID,
		Code:          c.Code,
		Name:          c.Name,
		Symbol:        c.Symbol,
		DecimalPlaces: c.DecimalPlaces,
		Enabled:       c.Enabled,
	}
}

// ExchangeRate Repository

type ExchangeRateRepository struct {
	db *gorm.DB
}

func NewExchangeRateRepository(db *gorm.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{db: db}
}

func (r *ExchangeRateRepository) FindByID(ctx context.Context, id uint) (*entity.CurrencyExchangeRate, error) {
	var m model.CurrencyExchangeRateModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("exchange rate not found: %w", err)
	}
	return exchangeRateModelToEntity(&m), nil
}

func (r *ExchangeRateRepository) FindByPair(ctx context.Context, fromID, toID uint, date time.Time) (*entity.CurrencyExchangeRate, error) {
	var m model.CurrencyExchangeRateModel
	dateStr := date.Format("2006-01-02")
	if err := r.db.WithContext(ctx).
		Where("from_currency_id = ? AND to_currency_id = ? AND date = ? AND deleted_at IS NULL", fromID, toID, dateStr).
		First(&m).Error; err != nil {
		return nil, fmt.Errorf("exchange rate not found: %w", err)
	}
	return exchangeRateModelToEntity(&m), nil
}

func (r *ExchangeRateRepository) List(ctx context.Context, userGroupID uint, limit, offset int) ([]entity.CurrencyExchangeRate, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.CurrencyExchangeRateModel{}).Where("deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	query.Count(&total)

	var models []model.CurrencyExchangeRateModel
	if err := query.Order("date DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	rates := make([]entity.CurrencyExchangeRate, len(models))
	for i, m := range models {
		rates[i] = *exchangeRateModelToEntity(&m)
	}
	return rates, total, nil
}

func (r *ExchangeRateRepository) ListByPair(ctx context.Context, fromCode, toCode string) ([]entity.CurrencyExchangeRate, error) {
	var models []model.CurrencyExchangeRateModel
	if err := r.db.WithContext(ctx).
		Joins("JOIN transaction_currencies fc ON fc.id = currency_exchange_rates.from_currency_id").
		Joins("JOIN transaction_currencies tc ON tc.id = currency_exchange_rates.to_currency_id").
		Where("fc.code = ? AND tc.code = ? AND currency_exchange_rates.deleted_at IS NULL", fromCode, toCode).
		Order("date DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	rates := make([]entity.CurrencyExchangeRate, len(models))
	for i, m := range models {
		rates[i] = *exchangeRateModelToEntity(&m)
	}
	return rates, nil
}

func (r *ExchangeRateRepository) Create(ctx context.Context, rate *entity.CurrencyExchangeRate) error {
	m := exchangeRateEntityToModel(rate)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	rate.ID = m.ID
	rate.CreatedAt = m.CreatedAt
	rate.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *ExchangeRateRepository) Update(ctx context.Context, rate *entity.CurrencyExchangeRate) error {
	m := exchangeRateEntityToModel(rate)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *ExchangeRateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.CurrencyExchangeRateModel{}, id).Error
}

func (r *ExchangeRateRepository) DeleteByPair(ctx context.Context, fromID, toID uint) error {
	return r.db.WithContext(ctx).Where("from_currency_id = ? AND to_currency_id = ?", fromID, toID).
		Delete(&model.CurrencyExchangeRateModel{}).Error
}

func exchangeRateModelToEntity(m *model.CurrencyExchangeRateModel) *entity.CurrencyExchangeRate {
	e := &entity.CurrencyExchangeRate{
		ID:             m.ID,
		UserGroupID:    m.UserGroupID,
		FromCurrencyID: m.FromCurrencyID,
		ToCurrencyID:   m.ToCurrencyID,
		Date:           m.Date,
		Rate:           m.Rate,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
	}
	if m.UserRate != nil {
		e.UserRate = *m.UserRate
	}
	return e
}

func exchangeRateEntityToModel(e *entity.CurrencyExchangeRate) *model.CurrencyExchangeRateModel {
	m := &model.CurrencyExchangeRateModel{
		ID:             e.ID,
		UserGroupID:    e.UserGroupID,
		FromCurrencyID: e.FromCurrencyID,
		ToCurrencyID:   e.ToCurrencyID,
		Date:           e.Date,
		Rate:           e.Rate,
	}
	if e.UserRate != "" {
		m.UserRate = &e.UserRate
	}
	return m
}
