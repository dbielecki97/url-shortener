package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"github.com/dbielecki97/url-shortener/pkg/logger"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

type Postgresql struct {
	db *sqlx.DB
}

func New() *Postgresql {
	host := os.Getenv("POSTGRES_HOST")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, 5432, "postgresuser", "postgrespass", "url-shortener")

	var db *sqlx.DB
	var err error
	if db, err = sqlx.Open("postgres", psqlInfo); err != nil {
		logger.Fatal("Could not open connection to postgres", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal("Could not ping Postgresql: %v", err)
	}

	log.Println("Connected to Postgresql...")
	return &Postgresql{db: db}
}

func (p Postgresql) Save(entry *domain.ShortURL) (*domain.ShortURL, errs.RestErr) {
	sqlInsert := "INSERT into urls (code, url, created_at) values ($1,$2,$3)"

	_, err := p.db.Exec(sqlInsert, entry, entry.URL, entry.CreatedAt)
	if err != nil {
		logger.Error("Could not save ShortURL to client: %v", err)
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	sqlSelect := "SELECT * from urls where code = $1"
	row := p.db.QueryRowx(sqlSelect, entry.Code)

	var res domain.ShortURL
	err = row.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Could not find saved ShortURL: %v", err)
			return nil, errs.NewUnexpectedError("unexpected database error")
		}

		logger.Error("Could not scan ShortURL: %v", err)
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return &res, nil
}

func (p Postgresql) Find(code string) (*domain.ShortURL, errs.RestErr) {
	sqlSelect := "SELECT * from urls where code = $1"
	row := p.db.QueryRowx(sqlSelect, code)

	var res domain.ShortURL
	err := row.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("incorrect code")
		}

		logger.Error("Could not scan ShortURL: %v", err)
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return &res, nil
}
