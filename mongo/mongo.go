package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

const MAX_CONNECTION_TIME = 10

type MClient struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func MongoConnect(ctx context.Context, uri string, option bool) *MClient {
	var client *mongo.Client
	var err error

	mongoConn := options.Client().ApplyURI(uri)

	if option {
		// 1. 로그기록 사용 여부
		// 2. 기록 저장용 노드 최대
		// 3. 쓰기에 사용 될 수 잇는 최대 시간
		jmajority := writeconcern.New(writeconcern.J(true))
		wmajority := writeconcern.New(writeconcern.W(1))
		tmajority := writeconcern.New(writeconcern.WTimeout(1000 * time.Microsecond))
		readConcert := readconcern.New(readconcern.Level("majority"))

		mongoConn.
			SetConnectTimeout(MAX_CONNECTION_TIME * time.Second).
			SetMaxPoolSize(50).SetMinPoolSize(5).
			SetWriteConcern(jmajority).
			SetWriteConcern(wmajority).
			SetWriteConcern(tmajority).
			SetReadConcern(readConcert)
	}

	if client, err = mongo.Connect(ctx, mongoConn); err != nil {
		panic(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		panic(err)
	}

	return &MClient{Client: client}
}

func (m *MClient) SetMongoDataBase(dbName string) {
	if m.DB != nil {
		panic("Already Set Mongo DB")
	}

	m.DB = m.Client.Database(dbName)
}

func (m *MClient) GetSession(collection string) *mongo.Collection {
	return m.DB.Collection(collection)
}
