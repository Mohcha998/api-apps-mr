package db

import (
	"apps/config"
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MysqlDBInterface interface {
	Create(ctx context.Context, doc any) error
	Update(ctx context.Context, data any) error
	Find(ctx context.Context, data any, opts ...FindOption) error
	FindOne(ctx context.Context, data any, opts ...FindOption) error
	Count(ctx context.Context, model any, total *int64, opts ...FindOption) error
	CreateInBatches(ctx context.Context, data any, batchSize int) error
	WithTransaction(function func() error) error
}

type MysqlDB struct {
	db *gorm.DB
}

type Query struct {
	Query string
	Args  []any
}

func NewQuery(query string, args ...any) Query {
	return Query{
		Query: query,
		Args:  args,
	}
}

func NewMysqlConnection(cfg *config.Config) (*MysqlDB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.DbName,
	)

	conn, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Warn),
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	// if err = conn.AutoMigrate(
	// 	&domain.Customer{},
	// 	&domain.Category{},
	// 	&domain.Product{},
	// 	&domain.Cart{},
	// 	&domain.Order{},
	// 	&domain.OrderItem{},
	// 	&domain.Payment{},
	// ); err != nil {
	// 	return nil, err
	// }

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.Mysql.MaxIdleConnection)
	sqlDB.SetMaxOpenConns(cfg.Mysql.MaxIdleConnection)

	return &MysqlDB{
		db: conn,
	}, nil
}

func (d *MysqlDB) Create(ctx context.Context, data any) error {
	return d.db.WithContext(ctx).Create(data).Error
}

func (d *MysqlDB) Update(ctx context.Context, data any) error {
	return d.db.WithContext(ctx).Save(data).Error
}

func (d *MysqlDB) Find(ctx context.Context, data any, opts ...FindOption) error {
	query := d.applyOptions(opts...)
	if err := query.WithContext(ctx).Find(data).Error; err != nil {
		return err
	}

	return nil
}

func (d *MysqlDB) FindOne(ctx context.Context, data any, opts ...FindOption) error {
	query := d.applyOptions(opts...)
	if err := query.WithContext(ctx).First(data).Error; err != nil {
		return err
	}

	return nil
}

func (d *MysqlDB) Count(ctx context.Context, model any, total *int64, opts ...FindOption) error {
	query := d.applyOptions(opts...)
	if err := query.Model(model).WithContext(ctx).Count(total).Error; err != nil {
		return err
	}

	return nil
}

func (d *MysqlDB) CreateInBatches(ctx context.Context, data any, batchSize int) error {
	return d.db.WithContext(ctx).CreateInBatches(data, batchSize).Error
}

func (d *MysqlDB) WithTransaction(function func() error) error {
	callback := func(db *gorm.DB) error {
		return function()
	}

	tx := d.db.Begin()
	if err := callback(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (d *MysqlDB) applyOptions(opts ...FindOption) *gorm.DB {
	query := d.db

	opt := getOption(opts...)

	if len(opt.preloads) != 0 {
		for _, preload := range opt.preloads {
			query = query.Preload(preload)
		}
	}

	if opt.query != nil {
		for _, q := range opt.query {
			query = query.Where(q.Query, q.Args)
		}
	}

	if opt.order != "" {
		query = query.Order(opt.order)
	}

	if opt.offset != 0 {
		query = query.Offset(opt.offset)
	}

	if opt.limit != 0 {
		query = query.Limit(opt.limit)
	}

	return query
}