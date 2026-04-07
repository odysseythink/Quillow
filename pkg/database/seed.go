package database

import (
	"log"

	"github.com/anthropics/quillow/internal/adapter/repository/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed populates reference/lookup tables with initial data.
func Seed(db *gorm.DB) error {
	if err := seedAccountTypes(db); err != nil {
		return err
	}
	if err := seedTransactionTypes(db); err != nil {
		return err
	}
	if err := seedCurrencies(db); err != nil {
		return err
	}
	if err := seedLinkTypes(db); err != nil {
		return err
	}
	if err := seedRoles(db); err != nil {
		return err
	}
	if err := seedConfig(db); err != nil {
		return err
	}
	if err := seedDefaultUser(db); err != nil {
		return err
	}
	return nil
}

func seedAccountTypes(db *gorm.DB) error {
	types := []string{
		"Asset account",
		"Beneficiary account",
		"Cash account",
		"Credit card",
		"Debt",
		"Default account",
		"Expense account",
		"Import account",
		"Initial balance account",
		"Liability credit account",
		"Loan",
		"Mortgage",
		"Reconciliation account",
		"Revenue account",
	}
	for _, t := range types {
		m := model.AccountTypeModel{Type: t}
		result := db.Where("type = ?", t).FirstOrCreate(&m)
		if result.Error != nil {
			return result.Error
		}
	}
	log.Printf("Seeded %d account types", len(types))
	return nil
}

func seedTransactionTypes(db *gorm.DB) error {
	types := []string{
		"Deposit",
		"Invalid",
		"Liability credit",
		"Opening balance",
		"Reconciliation",
		"Transfer",
		"Withdrawal",
	}
	for _, t := range types {
		m := model.TransactionTypeModel{Type: t}
		result := db.Where("type = ?", t).FirstOrCreate(&m)
		if result.Error != nil {
			return result.Error
		}
	}
	log.Printf("Seeded %d transaction types", len(types))
	return nil
}

type currencyDef struct {
	Code          string
	Name          string
	Symbol        string
	DecimalPlaces int
	Enabled       bool
}

func seedCurrencies(db *gorm.DB) error {
	currencies := []currencyDef{
		{"EUR", "Euro", "€", 2, true},
		{"USD", "US Dollar", "$", 2, false},
		{"GBP", "British Pound", "£", 2, false},
		{"HUF", "Hungarian forint", "Ft", 2, false},
		{"UAH", "Ukrainian hryvnia", "₴", 2, false},
		{"PLN", "Polish złoty", "zł", 2, false},
		{"TRY", "Turkish lira", "₺", 2, false},
		{"DKK", "Dansk krone", "kr.", 2, false},
		{"ISK", "Íslensk króna", "kr.", 2, false},
		{"NOK", "Norsk krone", "kr.", 2, false},
		{"SEK", "Svensk krona", "kr.", 2, false},
		{"RON", "Romanian leu", "lei", 2, false},
		{"BRL", "Brazilian real", "R$", 2, false},
		{"CAD", "Canadian dollar", "C$", 2, false},
		{"MXN", "Mexican peso", "MX$", 2, false},
		{"PEN", "Peruvian Sol", "S/", 2, false},
		{"ARS", "Argentinian Peso", "$", 2, false},
		{"COP", "Colombian Peso", "$", 2, false},
		{"CLP", "Chilean Peso", "$", 2, false},
		{"UYU", "Uruguayan Peso", "$", 2, false},
		{"IDR", "Indonesian rupiah", "Rp", 2, false},
		{"AUD", "Australian dollar", "A$", 2, false},
		{"NZD", "New Zealand dollar", "NZ$", 2, false},
		{"EGP", "Egyptian pound", "E£", 2, false},
		{"MAD", "Moroccan dirham", "DH", 2, false},
		{"ZAR", "South African rand", "R", 2, false},
		{"JPY", "Japanese yen", "¥", 0, false},
		{"CNY", "Chinese yuan", "¥", 2, false},
		{"KRW", "South Korean won", "₩", 2, false},
		{"RUB", "Russian ruble", "₽", 2, false},
		{"INR", "Indian rupee", "₹", 2, false},
		{"ILS", "Israeli new shekel", "₪", 2, false},
		{"CHF", "Swiss franc", "CHF", 2, false},
		{"HRK", "Croatian kuna", "kn", 2, false},
		{"HKD", "Hong Kong dollar", "HK$", 2, false},
		{"CZK", "Czech koruna", "Kč", 2, false},
		{"KZT", "Kazakhstani tenge", "₸", 2, false},
		{"SAR", "Saudi Riyal", "SAR", 2, false},
		{"RSD", "Serbian Dinar", "RSD", 2, false},
		{"TWD", "New Taiwan Dollar", "NT$", 0, false},
		{"THB", "Thai baht", "฿", 2, false},
	}
	for _, c := range currencies {
		m := model.TransactionCurrencyModel{
			Code:          c.Code,
			Name:          c.Name,
			Symbol:        c.Symbol,
			DecimalPlaces: c.DecimalPlaces,
			Enabled:       c.Enabled,
		}
		result := db.Where("code = ?", c.Code).FirstOrCreate(&m)
		if result.Error != nil {
			return result.Error
		}
	}
	log.Printf("Seeded %d currencies", len(currencies))
	return nil
}

func seedLinkTypes(db *gorm.DB) error {
	type linkDef struct {
		Name     string
		Inward   string
		Outward  string
		Editable bool
	}
	links := []linkDef{
		{"Related", "relates to", "relates to", false},
		{"Refund", "is (partially) refunded by", "(partially) refunds", false},
		{"Paid", "is (partially) paid for by", "(partially) pays for", false},
		{"Reimbursement", "is (partially) reimbursed by", "(partially) reimburses", false},
	}
	for _, l := range links {
		m := model.LinkTypeModel{
			Name:     l.Name,
			Inward:   l.Inward,
			Outward:  l.Outward,
			Editable: l.Editable,
		}
		result := db.Where("name = ?", l.Name).FirstOrCreate(&m)
		if result.Error != nil {
			return result.Error
		}
	}
	log.Printf("Seeded %d link types", len(links))
	return nil
}

func seedRoles(db *gorm.DB) error {
	type roleDef struct {
		Name        string
		DisplayName string
		Description string
	}
	roles := []roleDef{
		{"owner", "Site Owner", "User runs this instance of FF3"},
		{"demo", "Demo User", "User is a demo user"},
	}
	for _, r := range roles {
		m := model.RoleModel{
			Name:        r.Name,
			DisplayName: &r.DisplayName,
			Description: &r.Description,
		}
		result := db.Where("name = ?", r.Name).FirstOrCreate(&m)
		if result.Error != nil {
			return result.Error
		}
	}
	log.Printf("Seeded %d roles", len(roles))
	return nil
}

func seedConfig(db *gorm.DB) error {
	m := model.ConfigurationModel{
		Name: "db_version",
		Data: "1",
	}
	result := db.Where("name = ?", "db_version").FirstOrCreate(&m)
	if result.Error != nil {
		return result.Error
	}
	log.Println("Seeded configuration")
	return nil
}

func seedDefaultUser(db *gorm.DB) error {
	var count int64
	db.Model(&model.UserModel{}).Count(&count)
	if count > 0 {
		return nil
	}

	// Create default user group
	group := model.UserGroupModel{Title: "Default"}
	if err := db.FirstOrCreate(&group, "title = ?", "Default").Error; err != nil {
		return err
	}

	// Default admin: admin@firefly.local / firefly
	hashed, err := bcrypt.GenerateFromPassword([]byte("firefly"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := model.UserModel{
		Email:       "admin@firefly.local",
		Password:    string(hashed),
		UserGroupID: &group.ID,
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	// Assign owner role
	var ownerRole model.RoleModel
	if err := db.Where("name = ?", "owner").First(&ownerRole).Error; err == nil {
		db.Create(&model.RoleUserModel{UserID: user.ID, RoleID: ownerRole.ID})
	}

	log.Printf("Seeded default admin user: admin@firefly.local / firefly")
	return nil
}
