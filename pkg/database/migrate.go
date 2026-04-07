package database

import (
	"fmt"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"gorm.io/gorm"
)

// AutoMigrate creates or updates all database tables based on GORM models.
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		// Users & Auth
		&model.UserModel{},
		&model.RoleModel{},
		&model.RoleUserModel{},
		&model.GroupMembershipModel{},
		&model.InvitedUserModel{},
		&model.PersonalAccessTokenModel{},
		&model.UserGroupModel{},

		// Account
		&model.AccountTypeModel{},
		&model.AccountModel{},
		&model.AccountMetaModel{},

		// Transaction
		&model.TransactionTypeModel{},
		&model.TransactionGroupModel{},
		&model.TransactionJournalModel{},
		&model.TransactionModel{},
		&model.TransactionJournalMetaModel{},
		&model.TransactionJournalLinkModel{},

		// Currency
		&model.TransactionCurrencyModel{},
		&model.CurrencyExchangeRateModel{},

		// Budget
		&model.BudgetModel{},
		&model.BudgetLimitModel{},
		&model.AutoBudgetModel{},
		&model.AvailableBudgetModel{},

		// Bill
		&model.BillModel{},

		// Category
		&model.CategoryModel{},

		// Tag
		&model.TagModel{},

		// PiggyBank
		&model.PiggyBankModel{},
		&model.PiggyBankEventModel{},
		&model.PiggyBankRepetitionModel{},

		// Rule
		&model.RuleGroupModel{},
		&model.RuleModel{},
		&model.RuleTriggerModel{},
		&model.RuleActionModel{},

		// Recurrence
		&model.RecurrenceModel{},
		&model.RecurrenceRepetitionModel{},
		&model.RecurrenceTransactionModel{},
		&model.RecurrenceMetaModel{},

		// Link Type
		&model.LinkTypeModel{},

		// Object Group
		&model.ObjectGroupModel{},

		// Attachment
		&model.AttachmentModel{},

		// AI Classification
		&model.ClassificationPatternModel{},

		// Preference & Configuration
		&model.PreferenceModel{},
		&model.ConfigurationModel{},

		// Webhook
		&model.WebhookModel{},
		&model.WebhookMessageModel{},
		&model.WebhookAttemptModel{},

		// Audit & Notes
		&model.AuditLogEntryModel{},
		&model.NoteModel{},
		&model.LocationModel{},
		&model.PeriodStatisticModel{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate failed: %w", err)
	}
	return nil
}
