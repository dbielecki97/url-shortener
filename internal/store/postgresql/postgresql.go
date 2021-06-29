package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/dbielecki97/url-shortener/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

func (p Postgresql) Save(entry *domain.ShortURL) (*domain.ShortURL, error) {
	sqlInsert := "INSERT into urls (code, url, created_at) values ($1,$2,$3)"

	_, err := p.db.Exec(sqlInsert, entry.Code, entry.URL, entry.CreatedAt)
	if err != nil {
		return nil, errors.Errorf("unexpected database error: %v", err)
	}

	sqlSelect := "SELECT * from urls where code = $1"
	row := p.db.QueryRowx(sqlSelect, entry.Code)

	var res domain.ShortURL
	err = row.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NotFoundError{Err: errors.New("incorrect code")}
		}
		return nil, errors.Errorf("unexpected database error: %v", err)
	}

	return &res, nil
}

func (p Postgresql) Find(code string) (*domain.ShortURL, error) {
	sqlSelect := "SELECT * from urls where code = $1"
	row := p.db.QueryRowx(sqlSelect, code)

	var res domain.ShortURL
	err := row.StructScan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NotFoundError{Err: errors.New("incorrect code")}
		}
		return nil, errors.Errorf("unexpected database error: %v", err)
	}

	return &res, nil
}
