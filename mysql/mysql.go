package mysql

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func NewMySQLClient(sourceUri string, option bool) *Queries {
	dbInstance, err := sql.Open("mysql", sourceUri)
	if err != nil {
		panic(err.Error())
	}

	err = dbInstance.Ping()
	if err != nil {
		panic(err.Error())
	}

	if option {
		dbInstance.SetConnMaxLifetime(time.Minute * 1)
		dbInstance.SetMaxIdleConns(3)
		dbInstance.SetMaxOpenConns(6)
	}

	return New(dbInstance)
}

// !!!!!!!!!!!! EXAMPLE !!!!!!!!!!!!!!!!
// ---> sqlc를 활용해서 Query만 생성하면 이러한 형태로 생성되니 사용하면 된다.

// const createNewComment = `-- name: CreateNewComment :execresult
// INSERT INTO comment (
//     post_id,
//     comment_owner_account,
//     text
// ) VALUES (
//    ?, ?, ?
// )
// `

// type CreateNewCommentParams struct {
// 	PostID              int64  `json:"post_id"`
// 	CommentOwnerAccount string `json:"comment_owner_account"`
// 	Text                string `json:"text"`
// }

// func (q *Queries) CreateNewComment(ctx context.Context, arg CreateNewCommentParams) (sql.Result, error) {
// 	return q.db.ExecContext(ctx, createNewComment, arg.PostID, arg.CommentOwnerAccount, arg.Text)
// }
