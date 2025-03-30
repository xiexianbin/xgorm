// Copyright 2025~time.Now xiexianbin<me@xiexianbin.cn>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package paginate implements GORM Pagination support for the Go language.
//
// Source code and other details for the project are available at GitHub:
//
//	https://github.com/xiexianbin/xgorm

package xgorm

import (
	"context"

	"gorm.io/gorm"
)

// IRepository general Repository interface
type IRepository[T any] interface {
	GetDB() *gorm.DB

	// create on record
	Create(ctx context.Context, entity *T) error

	// batch create records
	CreateBatch(ctx context.Context, entities []*T, batchSizes ...int) error

	// update record
	Update(ctx context.Context, entity *T) error

	// delete record by id
	Delete(ctx context.Context, id any) error

	// find record by id
	FindByID(ctx context.Context, id any) (*T, error)

	// find all records
	FindAll(ctx context.Context) ([]*T, error)

	// find records by condition
	FindByCondition(ctx context.Context, conds ...interface{}) ([]*T, error)

	// get records count by condition
	CountByCondition(ctx context.Context, conds ...interface{}) (int64, error)

	// check is record exist by condition
	ExistsByCondition(ctx context.Context, conds ...interface{}) (bool, error)

	// run Transaction
	Transaction(ctx context.Context, fn func(txRepo IRepository[T]) error) error
}

// Repository base on gorm implementation
type Repository[T any] struct {
	db *gorm.DB
}

// NewRepository creating a new repository record
func NewRepository[T any](db *gorm.DB) IRepository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) GetDB() *gorm.DB {
	return r.db
}

// Create create on record
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// CreateBatch batch create records
// fix: Prepared statement contains too many placeholders
// For KDF method pbkdf2_hmac iterations value less than 1000 or more than 65535 is not allowed due to security reasons. Please provide iterations >= 1000 and iterations < 65535
//
// https://dev.mysql.com/doc/mysql-errors/8.4/en/server-error-reference.html
func (r *Repository[T]) CreateBatch(ctx context.Context, entities []*T, batchSizes ...int) error {
	var batchSize int
	if len(batchSizes) > 0 && batchSizes[0] > 0 && batchSizes[0] <= 1000 {
		batchSize = batchSizes[0]
	} else {
		batchSize = 1000
	}

	for i := 0; i < len(entities); i += batchSize {
		end := i + batchSize
		if end > len(entities) {
			end = len(entities)
		}
		err := r.db.WithContext(ctx).CreateInBatches(entities[i:end], len(entities[i:end])).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// Update update record
func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Delete delete record by id
func (r *Repository[T]) Delete(ctx context.Context, id any) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

// FindByID find record by id
func (r *Repository[T]) FindByID(ctx context.Context, id any) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// FindAll find all records
func (r *Repository[T]) FindAll(ctx context.Context) ([]*T, error) {
	var entities []*T
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

// FindByCondition find records by condition
func (r *Repository[T]) FindByCondition(ctx context.Context, conds ...interface{}) ([]*T, error) {
	var entities []*T
	err := r.db.WithContext(ctx).Where(conds[0], conds[1:]...).Find(&entities).Error
	return entities, err
}

// CountByCondition get records count by condition
func (r *Repository[T]) CountByCondition(ctx context.Context, conds ...interface{}) (int64, error) {
	var count int64
	var entity T
	err := r.db.WithContext(ctx).Model(&entity).Where(conds[0], conds[1:]...).Count(&count).Error
	return count, err
}

// ExistsByCondition check is record exist by condition
func (r *Repository[T]) ExistsByCondition(ctx context.Context, conds ...interface{}) (bool, error) {
	count, err := r.CountByCondition(ctx, conds...)
	return count > 0, err
}

// Transaction run Transaction
func (r *Repository[T]) Transaction(ctx context.Context, fn func(txRepo IRepository[T]) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository[T](tx)
		return fn(txRepo)
	})
}
