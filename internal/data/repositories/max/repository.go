package max_repository

// import (
// 	"context"
// 	db "daos_core/internal/bootstrap/db"
// )

// var _ repository.MaxRepository = (*maxRepositoryImpl)(nil)

// type maxRepositoryImpl struct {
// 	// Db Postgres
// 	db *db.Postgres
// }

// func NewMaxRepositoryImpl(db *db.Postgres) repository.MaxRepository {
// 	return &maxRepositoryImpl{db: db}
// }

// func (r *maxRepositoryImpl) GetBotToken(id string) (*string, error) {
// 	const getBotQuery = `
// 		SELECT token
// 		FROM bot_token
// 		WHERE id == $1
// 	`

// 	//var model sqlmodels.UsersSqlModel
// 	err := r.db.Pool.QueryRow(
// 		context.Background(),
// 		getBotQuery,
// 		id,
// 	).Scan(
// 	// &model.Id,
// 	// &model.Referer,
// 	// &model.AccessToken,
// 	// &model.RefreshToken,
// 	// &model.ExpiresAt,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return nil, nil
// }
