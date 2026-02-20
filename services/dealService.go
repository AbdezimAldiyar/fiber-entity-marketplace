package services

import (
	"database/sql"
	"errors"
	"time"

	"marketplace_entity/database"
	"marketplace_entity/models"
)

var ErrDealNotFound = errors.New("deal not found")

// ----------------- GET -----------------

func GetDealByID(dealID int) (models.Deal, error) {
	var deal models.Deal

	query := `
		SELECT deal_id, request_id, executor_id, agreed_price, status, created_at, closed_at
		FROM deals
		WHERE deal_id = $1
	`

	err := database.DB.QueryRow(query, dealID).Scan(
		&deal.DealID,
		&deal.RequestID,
		&deal.ExecutorID,
		&deal.AgreedPrice,
		&deal.Status,
		&deal.CreatedAt,
		&deal.ClosedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Deal{}, ErrDealNotFound
	}
	if err != nil {
		return models.Deal{}, err
	}

	return deal, nil
}

func GetDeals(limit int) ([]models.Deal, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	query := `
		SELECT deal_id, request_id, executor_id, agreed_price, status, created_at, closed_at
		FROM deals
		ORDER BY deal_id DESC
		LIMIT $1
	`

	rows, err := database.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deals := make([]models.Deal, 0, limit)

	for rows.Next() {
		var deal models.Deal
		if err := rows.Scan(
			&deal.DealID,
			&deal.RequestID,
			&deal.ExecutorID,
			&deal.AgreedPrice,
			&deal.Status,
			&deal.CreatedAt,
			&deal.ClosedAt,
		); err != nil {
			return nil, err
		}
		deals = append(deals, deal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deals, nil
}

// ----------------- POST -----------------

func CreateDeal(newDeal models.Deal) (models.Deal, error) {
	query := `
		INSERT INTO deals (request_id, executor_id, agreed_price, status)
		VALUES ($1, $2, $3, $4)
		RETURNING deal_id, request_id, executor_id, agreed_price, status, created_at, closed_at
	`

	var createdDeal models.Deal
	err := database.DB.QueryRow(
		query,
		newDeal.RequestID,
		newDeal.ExecutorID,
		newDeal.AgreedPrice,
		newDeal.Status,
	).Scan(
		&createdDeal.DealID,
		&createdDeal.RequestID,
		&createdDeal.ExecutorID,
		&createdDeal.AgreedPrice,
		&createdDeal.Status,
		&createdDeal.CreatedAt,
		&createdDeal.ClosedAt,
	)
	if err != nil {
		return models.Deal{}, err
	}

	return createdDeal, nil
}

// ----------------- PUT -----------------
func UpdateDealFull(dealID int, updatedDeal models.Deal) (models.Deal, error) {
	currentDeal, err := GetDealByID(dealID)
	if err != nil {
		return models.Deal{}, err
	}

	// Если статус = done и закрытие сделки ещё не стоит — ставим время закрытой сделки
	if updatedDeal.Status == "done" {
		if currentDeal.ClosedAt == nil {
			now := time.Now()
			updatedDeal.ClosedAt = &now
		} else {
			// если уже было закрыто — сохраняем старое время
			updatedDeal.ClosedAt = currentDeal.ClosedAt
		}
	} else {
		// если статус НЕ done — закрытие NULL
		updatedDeal.ClosedAt = nil
	}

	query := `
		UPDATE deals
		SET request_id = $1,
		    executor_id = $2,
		    agreed_price = $3,
		    status = $4,
		    closed_at = $5
		WHERE deal_id = $6
		RETURNING deal_id, request_id, executor_id, agreed_price, status, created_at, closed_at
	`

	var savedDeal models.Deal
	err = database.DB.QueryRow(
		query,
		updatedDeal.RequestID,
		updatedDeal.ExecutorID,
		updatedDeal.AgreedPrice,
		updatedDeal.Status,
		updatedDeal.ClosedAt,
		dealID,
	).Scan(
		&savedDeal.DealID,
		&savedDeal.RequestID,
		&savedDeal.ExecutorID,
		&savedDeal.AgreedPrice,
		&savedDeal.Status,
		&savedDeal.CreatedAt,
		&savedDeal.ClosedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return models.Deal{}, ErrDealNotFound
	}
	if err != nil {
		return models.Deal{}, err
	}

	return savedDeal, nil
}

// ----------------- PATCH -----------------

func UpdateDeal(dealID int, patchData map[string]interface{}) (models.Deal, error) {
	existingDeal, err := GetDealByID(dealID)
	if err != nil {
		return models.Deal{}, err
	}

	// request_id
	if requestIDValue, exists := patchData["request_id"]; exists {
		requestIDFloat, ok := requestIDValue.(float64)
		if !ok {
			return models.Deal{}, errors.New("request_id must be a number")
		}
		existingDeal.RequestID = int(requestIDFloat)
	}

	// executor_id
	if executorIDValue, exists := patchData["executor_id"]; exists {
		executorIDFloat, ok := executorIDValue.(float64)
		if !ok {
			return models.Deal{}, errors.New("executor_id must be a number")
		}
		existingDeal.ExecutorID = int(executorIDFloat)
	}

	// agreed_price
	if agreedPriceValue, exists := patchData["agreed_price"]; exists {
		agreedPriceFloat, ok := agreedPriceValue.(float64)
		if !ok {
			return models.Deal{}, errors.New("agreed_price must be a number")
		}
		existingDeal.AgreedPrice = agreedPriceFloat
	}

	// status
	if statusValue, exists := patchData["status"]; exists {
		statusString, ok := statusValue.(string)
		if !ok {
			return models.Deal{}, errors.New("status must be a string")
		}
		existingDeal.Status = statusString
	}

	// closed_at — руками менять запрещаем (оно выставляется логикой статуса)
	if _, exists := patchData["closed_at"]; exists {
		return models.Deal{}, errors.New("closed_at cannot be set manually")
	}

	return UpdateDealFull(dealID, existingDeal)
}

// ----------------- DELETE -----------------

func DeleteDeal(dealID int) error {
	deleteQuery := `DELETE FROM deals WHERE deal_id = $1`

	result, err := database.DB.Exec(deleteQuery, dealID)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return ErrDealNotFound
	}

	return nil
}
