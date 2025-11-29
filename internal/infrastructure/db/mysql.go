package db

import (
	"apps/config"
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MysqlDBInterface mendefinisikan fungsi-fungsi umum yang bisa digunakan untuk operasi database
type MysqlDBInterface interface {
	Create(ctx context.Context, doc any) error
	Update(ctx context.Context, data any) error
	Find(ctx context.Context, data any, opts ...FindOption) error
	FindOne(ctx context.Context, data any, opts ...FindOption) error
	Count(ctx context.Context, model any, total *int64, opts ...FindOption) error
	CreateInBatches(ctx context.Context, data any, batchSize int) error
	WithTransaction(function func() error) error
	Raw(ctx context.Context, dest any, query string, args ...any) error

	Conn() *gorm.DB
}

// MysqlDB adalah implementasi MysqlDBInterface yang membungkus *gorm.DB
type MysqlDB struct {
	db *gorm.DB
}

// Query struct menyimpan raw query (opsional)
type Query struct {
	Query string
	Args  []any
}

// NewQuery membuat struct Query baru
func NewQuery(query string, args ...any) Query {
	return Query{
		Query: query,
		Args:  args,
	}
}

// NewMysqlConnection membuat koneksi baru ke MySQL
func NewMysqlConnection(cfg *config.Config) (*MysqlDB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
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

	sqlDB, err := conn.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.Mysql.MaxIdleConnection)
	sqlDB.SetMaxOpenConns(cfg.Mysql.MaxOpenConnection)

	return &MysqlDB{
		db: conn,
	}, nil
}

// Create menambahkan satu record baru ke database
func (d *MysqlDB) Create(ctx context.Context, data any) error {
	return d.db.WithContext(ctx).Create(data).Error
}

//Raw Ctx
func (d *MysqlDB) Raw(ctx context.Context, dest any, query string, args ...any) error {
    return d.db.WithContext(ctx).Raw(query, args...).Scan(dest).Error
}

// Update memperbarui record yang sudah ada
func (d *MysqlDB) Update(ctx context.Context, data any) error {
	return d.db.WithContext(ctx).Save(data).Error
}

// Find mencari banyak data berdasarkan option yang diberikan
func (d *MysqlDB) Find(ctx context.Context, data any, opts ...FindOption) error {
	query := d.applyOptions(opts...)
	if err := query.WithContext(ctx).Find(data).Error; err != nil {
		return err
	}
	return nil
}

// FindOne mencari satu data berdasarkan option yang diberikan
func (d *MysqlDB) FindOne(ctx context.Context, data any, opts ...FindOption) error {
	query := d.applyOptions(opts...)
	if err := query.WithContext(ctx).First(data).Error; err != nil {
		return err
	}
	return nil
}

// Count menghitung total data berdasarkan kriteria
func (d *MysqlDB) Count(ctx context.Context, model any, total *int64, opts ...FindOption) error {
	query := d.applyOptions(opts...)
	if err := query.Model(model).WithContext(ctx).Count(total).Error; err != nil {
		return err
	}
	return nil
}

// CreateInBatches menambahkan banyak record dalam sekali eksekusi
func (d *MysqlDB) CreateInBatches(ctx context.Context, data any, batchSize int) error {
	return d.db.WithContext(ctx).CreateInBatches(data, batchSize).Error
}

// WithTransaction menjalankan operasi dalam transaksi
func (d *MysqlDB) WithTransaction(function func() error) error {
	tx := d.db.Begin()
	if err := function(); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// Conn mengembalikan koneksi *gorm.DB asli
func (d *MysqlDB) Conn() *gorm.DB {
	return d.db
}

// applyOptions menggabungkan semua opsi query (where, preload, order, dll)
func (d *MysqlDB) applyOptions(opts ...FindOption) *gorm.DB {
    query := d.db
    opt := getOption(opts...)

    // Preload
    if len(opt.preloads) != 0 {
        for _, preload := range opt.preloads {
            query = query.Preload(preload)
        }
    }

    // Where
    if opt.query != nil {
        for _, q := range opt.query {
            query = query.Where(q.Query, q.Args...)
        }
    }

    // Jika user override order → pakai itu
    if opt.order != "" {
        query = query.Order(opt.order)
    } else if !opt.noOrder {
        // AUTO DETECT TABLE NAME
        stmt := &gorm.Statement{DB: d.db}
        stmt.Parse(query.Statement.Model)

        table := stmt.Table

        if table == "user" {
            query = query.Order("id_user")
        } else {
            query = query.Order("id")
        }
    }

    // Offset
    if opt.offset != 0 {
        query = query.Offset(opt.offset)
    }

    // Limit
    if opt.limit != 0 {
        query = query.Limit(opt.limit)
    }

    return query
}

