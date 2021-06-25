package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/dbielecki97/url-shortener/pkg/errs"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"os"
)

type Postgresql struct {
	db  *sqlx.DB
	log *logrus.Logger
}

func New(log *logrus.Logger) (*Postgresql, func()) {
	host := os.Getenv("POSTGRES_HOST")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, 5432, "postgresuser", "postgrespass", "url-shortener")

	var db *sqlx.DB
	var err error
	if db, err = sqlx.Open("postgres", psqlInfo); err != nil {
		log.Fatalf("Could not open connection to postgres")
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not ping Postgresql: %v", err)
	}

	closeFn := func() {
		err := db.Close()
		if err != nil {
			log.Printf("Could not close Postgresql: %v", err)
		}
	}

	log.Println("Connected to Postgresql...")
	return &Postgresql{db: db, log: log}, closeFn
}

func (p Postgresql) Save(entry *domain.ShortURL) (*domain.ShortURL, *errs.AppError) {
	sqlInsert := "INSERT into urls (code, url, created_at) values ($1,$2,$3)"

	_, err := p.db.Exec(sqlInsert, entry.Code, entry.URL, entry.CreatedAt)
	if err != nil {
		p.log.Errorf("Could not save ShortURL to store: %v", err)
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	sqlSelect := "SELECT * from urls where code = $1"
	row := p.db.QueryRowx(sqlSelect, entry.Code)

	var res domain.ShortURL
	err = row.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			p.log.Errorf("Could not find saved ShortURL: %v", err)
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
		p.log.Errorf("Could not scan ShortURL: %v", err)
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return &res, nil
}

func (p Postgresql) Find(code string) (*domain.ShortURL, *errs.AppError) {
	sqlSelect := "SELECT * from urls where code = $1"
	row := p.db.QueryRowx(sqlSelect, code)

	var res domain.ShortURL
	err := row.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("incorrect code")
		}
		p.log.Errorf("Could not scan ShortURL: %v", err)
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	return &res, nil
}
