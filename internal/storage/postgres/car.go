package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/domain/models"
	"github.com/EwvwGeN/EffectiveMobile_assignment/internal/storage"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

func (pp *postgresProvider) SaveCars(ctx context.Context, carList []models.Car) error {
	tx, err := pp.dbConn.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return storage.ErrStartTx
	}
	batch := pgx.Batch{}
	for _, car := range carList {
		batch.Queue(fmt.Sprintf(`
			INSERT INTO "%s" (reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic)
			VALUES($1,$2,$3,$4,$5,$6,$7)
			ON CONFLICT (reg_num) DO NOTHING;`,
			pp.cfg.CarTable),
			car.RegisterNumber, car.Mark, car.Model,
			car.Year, car.Owner.Name, car.Owner.Surname,
			car.Owner.Patronymic,
		)
	}
	results := tx.SendBatch(ctx, &batch)
	for range carList {
		_, err = results.Exec()
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				continue
			}
			results.Close()
			if err := tx.Rollback(ctx); err != nil {
				return storage.ErrRollbackTx
			}
			return err
		}
	}
	results.Close()
	if err := tx.Commit(ctx); err != nil {
		return storage.ErrCommitTx
	}
	return nil
}

func (pp *postgresProvider) GetCarById(ctx context.Context, carId string) (models.Car, error) {
	row := pp.dbConn.QueryRow(ctx, fmt.Sprintf(`
		SELECT car_id, reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic
		FROM "%s"
		WHERE car_id = $1;`,
	pp.cfg.CarTable),
	carId)
	var (
		car models.Car
	)
	err := row.Scan(
		&car.Id,
		&car.RegisterNumber,
		&car.Mark,
		&car.Model,
		&car.Year,
		&car.Owner.Name,
		&car.Owner.Surname,
		&car.Owner.Patronymic)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Car{}, storage.ErrCarNotFound
		}
		return models.Car{}, err
	}
	return car, nil
}

func (pp *postgresProvider) GetCarsWithFilterAndPagination(ctx context.Context, pgOption models.PaginationOption, filter models.Filter) ([]models.Car, error) {
	var preparedQuery strings.Builder
	preparedQuery.WriteString(
		fmt.Sprintf(
			`SELECT car_id, reg_num, mark, model, year, owner_name, owner_surname, owner_patronymic FROM "%s" `,
		pp.cfg.CarTable))
	var (
		fieldCount int
		usedData []interface{}
	)
	filterCount := 0
	filters := filter.Fields
	if len(filters) != 0 {
		preparedQuery.WriteString("WHERE ")
	} 
	for _, field := range filters {
		if filterCount != 0 {
			preparedQuery.WriteString(field.UnionCondition + " ")
		}
		// SQL INJECTION, cant pass columns name via argument
		preparedQuery.WriteString(fmt.Sprintf("%s %s $%d ", field.Name, field.Operator, fieldCount+1))
		usedData = append(usedData, field.Value)
		fieldCount+=1
		filterCount++
	}
	if pgOption.Limit != 0 {
		preparedQuery.WriteString(fmt.Sprintf("LIMIT $%d OFFSET $%d", fieldCount+1, fieldCount+2))
		usedData = append(usedData, pgOption.Limit, pgOption.Offset)
	}
	rows, err := pp.dbConn.Query(ctx, preparedQuery.String(), usedData...)
	if err != nil {
		return nil, err
	}
	var outProducts []models.Car
	for rows.Next() {
		var car models.Car
		err := rows.Scan(
			&car.Id,
			&car.RegisterNumber,
			&car.Mark,
			&car.Model,
			&car.Year,
			&car.Owner.Name,
			&car.Owner.Surname,
			&car.Owner.Patronymic)
		if err != nil {
			return nil, err
		}
		outProducts = append(outProducts, car)
	}
	return outProducts, nil
}

func (pp *postgresProvider) UpdateCarById(ctx context.Context, carId string, newData models.CarForPatch) error {
	tx, err := pp.dbConn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return storage.ErrStartTx
	}
	var preparedQuery strings.Builder
	preparedQuery.WriteString(fmt.Sprintf("UPDATE \"%s\" SET ", pp.cfg.CarTable)) 
	fieldsCount := 0
	var usedData []interface{}
	if newData.RegisterNumber != nil {
		fieldsCount++
		preparedQuery.WriteString(fmt.Sprintf("\"reg_num\" = $%d, ", fieldsCount))
		usedData = append(usedData, *newData.RegisterNumber)
	}
	if newData.Mark != nil {
		fieldsCount++
		preparedQuery.WriteString(fmt.Sprintf("\"mark\" = $%d, ", fieldsCount))
		
		usedData = append(usedData, *newData.Mark)
	}
	if newData.Model != nil {
		fieldsCount++
		preparedQuery.WriteString(fmt.Sprintf("\"model\" = $%d, ", fieldsCount))
		
		usedData = append(usedData, *newData.Model)
	}
	if newData.Year != nil {
		fieldsCount++
		preparedQuery.WriteString(fmt.Sprintf("\"year\" = $%d, ", fieldsCount))
		
		usedData = append(usedData, *newData.Year)
	}
	if newData.Owner != nil {
		if newData.Owner.Name != nil {
			fieldsCount++
			preparedQuery.WriteString(fmt.Sprintf("\"owner_name\" = $%d, ", fieldsCount))
			
			usedData = append(usedData, *newData.Owner.Name)
		}
		if newData.Owner.Surname != nil {
			fieldsCount++
			preparedQuery.WriteString(fmt.Sprintf("\"owner_surname\" = $%d, ", fieldsCount))
			
			usedData = append(usedData, *newData.Owner.Surname)
		}
		if newData.Owner.Patronymic != nil {
			fieldsCount++
			preparedQuery.WriteString(fmt.Sprintf("\"owner_patronymic\" = $%d, ", fieldsCount))
			
			usedData = append(usedData, *newData.Owner.Patronymic)
		}
	}
	query :=  preparedQuery.String()
	if fieldsCount != 0 {
		query = query[:len(query)-2]
	}
	usedData = append(usedData, carId)
	_, err = tx.Exec(ctx, fmt.Sprintf("%s WHERE \"car_id\" = $%d", query, fieldsCount+1), usedData...)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return storage.ErrRollbackTx
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return storage.ErrCarExist
			}
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return storage.ErrCommitTx
	}
	return nil
}

func (pp *postgresProvider) DeleteCarById(ctx context.Context, carId string) error {
	tx, err := pp.dbConn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return storage.ErrStartTx
	}
	_, err = tx.Exec(ctx, fmt.Sprintf("DELETE FROM \"%s\" WHERE \"car_id\" = $1", pp.cfg.CarTable), carId)
	if err != nil {
		if err := tx.Rollback(ctx); err != nil {
			return storage.ErrRollbackTx
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return storage.ErrCommitTx
	}
	return nil
}