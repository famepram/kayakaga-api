package transaction

import (
	"errors"
	"kayakaga-api/domain/mysql"
	"time"

	"gorm.io/gorm"
)

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) ListTransactions(userID uint, filters *ListFilters) ([]mysql.Transaction, error) {
	var txns []mysql.Transaction
	query := r.db.Preload("Category").Preload("Source").Preload("Account").
		Where("user_id = ?", userID)

	query = r.applyPeriodFilter(query, filters.Period)

	if filters.AccountID != nil {
		query = query.Where("account_id = ?", *filters.AccountID)
	}
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.Merchant != "" {
		query = query.Where("merchant LIKE ?", "%"+filters.Merchant+"%")
	}
	if filters.IsRecurring != nil {
		query = query.Where("is_recurring = ?", *filters.IsRecurring)
	}

	err := query.Order("date DESC, time DESC").Find(&txns).Error
	return txns, err
}

func (r *repo) applyPeriodFilter(query *gorm.DB, period string) *gorm.DB {
	now := time.Now().UTC()

	switch period {
	case "today":
		return query.Where("DATE(date) = ?", now.Format("2006-01-02"))
	case "week":
		weekAgo := now.AddDate(0, 0, -7)
		return query.Where("date >= ?", weekAgo.Format("2006-01-02"))
	case "month":
		return query.Where("YEAR(date) = ? AND MONTH(date) = ?", now.Year(), int(now.Month()))
	case "last_month":
		lastMonth := now.AddDate(0, -1, 0)
		return query.Where("YEAR(date) = ? AND MONTH(date) = ?", lastMonth.Year(), int(lastMonth.Month()))
	case "year":
		return query.Where("YEAR(date) = ?", now.Year())
	default:
		return query
	}
}

func (r *repo) GetTransactionByID(id, userID uint) (*mysql.Transaction, error) {
	var txn mysql.Transaction
	err := r.db.Preload("Category").Preload("Source").Preload("Account").
		Where("id = ? AND user_id = ?", id, userID).
		First(&txn).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &txn, nil
}

func (r *repo) CreateTransaction(txn *mysql.Transaction) error {
	return r.db.Create(txn).Error
}

func (r *repo) UpdateTransaction(txn *mysql.Transaction) error {
	return r.db.Save(txn).Error
}

func (r *repo) DeleteTransaction(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&mysql.Transaction{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("transaction not found")
	}
	return nil
}

func (r *repo) BulkInsert(txns []mysql.Transaction) error {
	if len(txns) == 0 {
		return nil
	}
	return r.db.Create(&txns).Error
}

func (r *repo) CheckDuplicates(userID, accountID uint, dates []time.Time, amounts []int64, merchants []string) (map[int]bool, error) {
	if len(dates) == 0 {
		return make(map[int]bool), nil
	}

	duplicates := make(map[int]bool)

	for i := 0; i < len(dates); i++ {
		var count int64
		err := r.db.Model(&mysql.Transaction{}).
			Where("user_id = ? AND account_id = ? AND DATE(date) = ? AND amount = ? AND LOWER(merchant) = ?",
				userID, accountID, dates[i].Format("2006-01-02"), amounts[i], lower(merchants[i])).
			Count(&count).Error

		if err != nil {
			return nil, err
		}

		if count > 0 {
			duplicates[i] = true
		}
	}

	return duplicates, nil
}

func lower(s string) string {
	return s
}
