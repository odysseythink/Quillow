package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/anthropics/firefly-iii-go/internal/adapter/repository/model"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	"github.com/anthropics/firefly-iii-go/internal/port"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// TransactionGroup

func (r *TransactionRepository) FindGroupByID(ctx context.Context, id uint) (*entity.TransactionGroup, error) {
	var m model.TransactionGroupModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("transaction group not found: %w", err)
	}
	return txGroupModelToEntity(&m), nil
}

func (r *TransactionRepository) ListGroups(ctx context.Context, userGroupID uint, transactionType string, start, end *time.Time, limit, offset int) ([]entity.TransactionGroup, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.TransactionGroupModel{}).Where("transaction_groups.deleted_at IS NULL")
	if userGroupID > 0 {
		query = query.Where("transaction_groups.user_group_id = ?", userGroupID)
	}

	if transactionType != "" {
		var typeID uint
		r.db.WithContext(ctx).Model(&model.TransactionTypeModel{}).Where("type = ?", transactionType).Pluck("id", &typeID)
		if typeID > 0 {
			query = query.Where("transaction_groups.id IN (SELECT transaction_group_id FROM transaction_journals WHERE transaction_type_id = ? AND deleted_at IS NULL)", typeID)
		}
	}

	if start != nil {
		query = query.Where("transaction_groups.id IN (SELECT transaction_group_id FROM transaction_journals WHERE date >= ? AND deleted_at IS NULL)", start)
	}
	if end != nil {
		query = query.Where("transaction_groups.id IN (SELECT transaction_group_id FROM transaction_journals WHERE date <= ? AND deleted_at IS NULL)", end)
	}

	var total int64
	query.Count(&total)

	var models []model.TransactionGroupModel
	if err := query.Order("transaction_groups.id DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	groups := make([]entity.TransactionGroup, len(models))
	for i, m := range models {
		groups[i] = *txGroupModelToEntity(&m)
	}
	return groups, total, nil
}

func (r *TransactionRepository) CreateGroup(ctx context.Context, group *entity.TransactionGroup) error {
	m := txGroupEntityToModel(group)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	group.ID = m.ID
	group.CreatedAt = m.CreatedAt
	group.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *TransactionRepository) UpdateGroup(ctx context.Context, group *entity.TransactionGroup) error {
	m := txGroupEntityToModel(group)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *TransactionRepository) DeleteGroup(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.TransactionGroupModel{}, id).Error
}

// TransactionJournal

func (r *TransactionRepository) FindJournalByID(ctx context.Context, id uint) (*entity.TransactionJournal, error) {
	var m model.TransactionJournalModel
	if err := r.db.WithContext(ctx).Where("deleted_at IS NULL").First(&m, id).Error; err != nil {
		return nil, fmt.Errorf("transaction journal not found: %w", err)
	}
	return txJournalModelToEntity(&m), nil
}

func (r *TransactionRepository) ListJournalsByGroupID(ctx context.Context, groupID uint) ([]entity.TransactionJournal, error) {
	var models []model.TransactionJournalModel
	if err := r.db.WithContext(ctx).Where("transaction_group_id = ? AND deleted_at IS NULL", groupID).Order("`order` ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	journals := make([]entity.TransactionJournal, len(models))
	for i, m := range models {
		journals[i] = *txJournalModelToEntity(&m)
	}
	return journals, nil
}

func (r *TransactionRepository) CreateJournal(ctx context.Context, journal *entity.TransactionJournal) error {
	m := txJournalEntityToModel(journal)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	journal.ID = m.ID
	journal.CreatedAt = m.CreatedAt
	journal.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *TransactionRepository) UpdateJournal(ctx context.Context, journal *entity.TransactionJournal) error {
	m := txJournalEntityToModel(journal)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *TransactionRepository) DeleteJournal(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.TransactionJournalModel{}, id).Error
}

// Transaction line items

func (r *TransactionRepository) ListByJournalID(ctx context.Context, journalID uint) ([]entity.Transaction, error) {
	var models []model.TransactionModel
	if err := r.db.WithContext(ctx).Where("transaction_journal_id = ? AND deleted_at IS NULL", journalID).Find(&models).Error; err != nil {
		return nil, err
	}
	txns := make([]entity.Transaction, len(models))
	for i, m := range models {
		txns[i] = *txModelToEntity(&m)
	}
	return txns, nil
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, txn *entity.Transaction) error {
	m := txEntityToModel(txn)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	txn.ID = m.ID
	return nil
}

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, txn *entity.Transaction) error {
	m := txEntityToModel(txn)
	return r.db.WithContext(ctx).Save(m).Error
}

// Meta

func (r *TransactionRepository) GetJournalMeta(ctx context.Context, journalID uint, name string) (string, error) {
	var m model.TransactionJournalMetaModel
	if err := r.db.WithContext(ctx).Where("transaction_journal_id = ? AND name = ?", journalID, name).First(&m).Error; err != nil {
		return "", nil
	}
	return m.Data, nil
}

func (r *TransactionRepository) SetJournalMeta(ctx context.Context, journalID uint, name, value string) error {
	var existing model.TransactionJournalMetaModel
	err := r.db.WithContext(ctx).Where("transaction_journal_id = ? AND name = ?", journalID, name).First(&existing).Error
	if err != nil {
		return r.db.WithContext(ctx).Create(&model.TransactionJournalMetaModel{
			TransactionJournalID: journalID,
			Name:                 name,
			Data:                 value,
		}).Error
	}
	existing.Data = value
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *TransactionRepository) GetJournalNotes(ctx context.Context, journalID uint) (string, error) {
	var note model.NoteModel
	if err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type LIKE ? AND deleted_at IS NULL", journalID, "%TransactionJournal%").First(&note).Error; err != nil {
		return "", nil
	}
	if note.Text != nil {
		return *note.Text, nil
	}
	return "", nil
}

func (r *TransactionRepository) SetJournalNotes(ctx context.Context, journalID uint, text string) error {
	noteableType := "FireflyIII\\Models\\TransactionJournal"
	var existing model.NoteModel
	err := r.db.WithContext(ctx).Where("noteable_id = ? AND noteable_type = ?", journalID, noteableType).First(&existing).Error
	if err != nil {
		return r.db.WithContext(ctx).Create(&model.NoteModel{
			NoteableID:   journalID,
			NoteableType: noteableType,
			Text:         &text,
		}).Error
	}
	existing.Text = &text
	return r.db.WithContext(ctx).Save(&existing).Error
}

func (r *TransactionRepository) GetJournalTags(ctx context.Context, journalID uint) ([]string, error) {
	// Tags are stored in tag_transaction_journal pivot table joined with tags table
	var tags []string
	r.db.WithContext(ctx).Table("tags").
		Joins("JOIN tag_transaction_journal ttj ON ttj.tag_id = tags.id").
		Where("ttj.transaction_journal_id = ?", journalID).
		Pluck("tag", &tags)
	return tags, nil
}

func (r *TransactionRepository) SetJournalTags(ctx context.Context, journalID uint, tags []string) error {
	// Clear existing
	r.db.WithContext(ctx).Exec("DELETE FROM tag_transaction_journal WHERE transaction_journal_id = ?", journalID)

	for _, tagName := range tags {
		var tagModel model.TagModel
		err := r.db.WithContext(ctx).Where("tag = ?", tagName).First(&tagModel).Error
		if err != nil {
			// Create new tag
			tagModel = model.TagModel{Tag: tagName, TagMode: "nothing"}
			r.db.WithContext(ctx).Create(&tagModel)
		}
		r.db.WithContext(ctx).Exec("INSERT INTO tag_transaction_journal (tag_id, transaction_journal_id) VALUES (?, ?)", tagModel.ID, journalID)
	}
	return nil
}

// Transaction type

func (r *TransactionRepository) GetTransactionType(ctx context.Context, typeName string) (*entity.TransactionType, error) {
	var m model.TransactionTypeModel
	if err := r.db.WithContext(ctx).Where("type = ?", typeName).First(&m).Error; err != nil {
		return nil, fmt.Errorf("transaction type not found: %s", typeName)
	}
	return &entity.TransactionType{ID: m.ID, Type: m.Type}, nil
}

func (r *TransactionRepository) GetTransactionTypeByID(ctx context.Context, id uint) (*entity.TransactionType, error) {
	var m model.TransactionTypeModel
	if err := r.db.WithContext(ctx).First(&m, id).Error; err != nil {
		return nil, err
	}
	return &entity.TransactionType{ID: m.ID, Type: m.Type}, nil
}

// Search

func (r *TransactionRepository) SearchGroups(ctx context.Context, userGroupID uint, query string, limit, offset int) ([]entity.TransactionGroup, int64, error) {
	subQuery := r.db.WithContext(ctx).Model(&model.TransactionJournalModel{}).
		Select("DISTINCT transaction_group_id").
		Where("deleted_at IS NULL AND description LIKE ?", "%"+query+"%")
	if userGroupID > 0 {
		subQuery = subQuery.Where("user_group_id = ?", userGroupID)
	}

	var total int64
	r.db.WithContext(ctx).Model(&model.TransactionGroupModel{}).
		Where("id IN (?) AND deleted_at IS NULL", subQuery).Count(&total)

	var models []model.TransactionGroupModel
	if err := r.db.WithContext(ctx).
		Where("id IN (?) AND deleted_at IS NULL", subQuery).
		Order("id DESC").Limit(limit).Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	groups := make([]entity.TransactionGroup, len(models))
	for i, m := range models {
		groups[i] = *txGroupModelToEntity(&m)
	}
	return groups, total, nil
}

func (r *TransactionRepository) CountByQuery(ctx context.Context, userGroupID uint, query string) (int64, error) {
	subQuery := r.db.WithContext(ctx).Model(&model.TransactionJournalModel{}).
		Select("DISTINCT transaction_group_id").
		Where("deleted_at IS NULL AND description LIKE ?", "%"+query+"%")
	if userGroupID > 0 {
		subQuery = subQuery.Where("user_group_id = ?", userGroupID)
	}

	var count int64
	r.db.WithContext(ctx).Model(&model.TransactionGroupModel{}).
		Where("id IN (?) AND deleted_at IS NULL", subQuery).Count(&count)
	return count, nil
}

// Insight aggregations

func (r *TransactionRepository) SumByType(ctx context.Context, userGroupID uint, transactionType string, start, end time.Time) ([]port.InsightEntry, error) {
	var results []port.InsightEntry

	var typeID uint
	r.db.WithContext(ctx).Model(&model.TransactionTypeModel{}).Where("type = ?", transactionType).Pluck("id", &typeID)
	if typeID == 0 {
		return results, nil
	}

	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT a.id, a.name, SUM(t.amount) as difference, tc.id as currency_id, tc.code as currency_code
		FROM transactions t
		JOIN transaction_journals tj ON tj.id = t.transaction_journal_id
		JOIN accounts a ON a.id = t.account_id
		JOIN transaction_currencies tc ON tc.id = t.transaction_currency_id
		WHERE tj.transaction_type_id = ?
		AND tj.date >= ? AND tj.date <= ?
		AND tj.deleted_at IS NULL AND t.deleted_at IS NULL
		AND t.amount < 0
		GROUP BY a.id, a.name, tc.id, tc.code
	`, typeID, start, end).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry port.InsightEntry
		if err := rows.Scan(&entry.ID, &entry.Name, &entry.Difference, &entry.CurrencyID, &entry.CurrencyCode); err != nil {
			continue
		}
		results = append(results, entry)
	}

	return results, nil
}

func (r *TransactionRepository) SumByAccount(ctx context.Context, userGroupID uint, accountID uint, transactionType string, start, end time.Time) ([]port.InsightEntry, error) {
	var results []port.InsightEntry

	var typeID uint
	r.db.WithContext(ctx).Model(&model.TransactionTypeModel{}).Where("type = ?", transactionType).Pluck("id", &typeID)
	if typeID == 0 {
		return results, nil
	}

	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT a.id, a.name, SUM(t.amount) as difference, tc.id as currency_id, tc.code as currency_code
		FROM transactions t
		JOIN transaction_journals tj ON tj.id = t.transaction_journal_id
		JOIN accounts a ON a.id = t.account_id
		JOIN transaction_currencies tc ON tc.id = t.transaction_currency_id
		WHERE tj.transaction_type_id = ?
		AND t.account_id = ?
		AND tj.date >= ? AND tj.date <= ?
		AND tj.deleted_at IS NULL AND t.deleted_at IS NULL
		GROUP BY a.id, a.name, tc.id, tc.code
	`, typeID, accountID, start, end).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var entry port.InsightEntry
		if err := rows.Scan(&entry.ID, &entry.Name, &entry.Difference, &entry.CurrencyID, &entry.CurrencyCode); err != nil {
			continue
		}
		results = append(results, entry)
	}

	return results, nil
}

// Model conversions

func txGroupModelToEntity(m *model.TransactionGroupModel) *entity.TransactionGroup {
	g := &entity.TransactionGroup{
		ID:          m.ID,
		UserID:      m.UserID,
		UserGroupID: m.UserGroupID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
	if m.Title != nil {
		g.Title = *m.Title
	}
	return g
}

func txGroupEntityToModel(g *entity.TransactionGroup) *model.TransactionGroupModel {
	m := &model.TransactionGroupModel{
		ID:          g.ID,
		UserID:      g.UserID,
		UserGroupID: g.UserGroupID,
	}
	if g.Title != "" {
		m.Title = &g.Title
	}
	return m
}

func txJournalModelToEntity(m *model.TransactionJournalModel) *entity.TransactionJournal {
	return &entity.TransactionJournal{
		ID:                    m.ID,
		UserID:                m.UserID,
		UserGroupID:           m.UserGroupID,
		TransactionTypeID:     m.TransactionTypeID,
		BillID:                m.BillID,
		TransactionCurrencyID: m.TransactionCurrencyID,
		Description:           m.Description,
		Date:                  m.Date,
		Order:                 m.Order,
		TagCount:              m.TagCount,
		Encrypted:             m.Encrypted,
		Completed:             m.Completed,
		TransactionGroupID:    m.TransactionGroupID,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		DeletedAt:             m.DeletedAt,
	}
}

func txJournalEntityToModel(j *entity.TransactionJournal) *model.TransactionJournalModel {
	return &model.TransactionJournalModel{
		ID:                    j.ID,
		UserID:                j.UserID,
		UserGroupID:           j.UserGroupID,
		TransactionTypeID:     j.TransactionTypeID,
		BillID:                j.BillID,
		TransactionCurrencyID: j.TransactionCurrencyID,
		Description:           j.Description,
		Date:                  j.Date,
		Order:                 j.Order,
		TagCount:              j.TagCount,
		Encrypted:             j.Encrypted,
		Completed:             j.Completed,
		TransactionGroupID:    j.TransactionGroupID,
	}
}

func txModelToEntity(m *model.TransactionModel) *entity.Transaction {
	t := &entity.Transaction{
		ID:                    m.ID,
		TransactionJournalID:  m.TransactionJournalID,
		AccountID:             m.AccountID,
		TransactionCurrencyID: m.TransactionCurrencyID,
		ForeignCurrencyID:     m.ForeignCurrencyID,
		Amount:                m.Amount,
		Reconciled:            m.Reconciled,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		DeletedAt:             m.DeletedAt,
	}
	if m.ForeignAmount != nil {
		t.ForeignAmount = *m.ForeignAmount
	}
	if m.Description != nil {
		t.Description = *m.Description
	}
	return t
}

func txEntityToModel(t *entity.Transaction) *model.TransactionModel {
	m := &model.TransactionModel{
		ID:                    t.ID,
		TransactionJournalID:  t.TransactionJournalID,
		AccountID:             t.AccountID,
		TransactionCurrencyID: t.TransactionCurrencyID,
		ForeignCurrencyID:     t.ForeignCurrencyID,
		Amount:                t.Amount,
		Reconciled:            t.Reconciled,
	}
	if t.ForeignAmount != "" {
		m.ForeignAmount = &t.ForeignAmount
	}
	if t.Description != "" {
		m.Description = &t.Description
	}
	return m
}
