package sqllite

import "context"

func (db *Database) GetJWTToken(ctx context.Context) (string, error) {
	query := "Select jwt_Token from info"
	var jwt string
	err := db.db.QueryRowxContext(ctx, query).Scan(&jwt)
	if err != nil {
		return "", err
	}
	return jwt, nil
}
func (db *Database) SetJWTToken(ctx context.Context, JWTToken string) error {
	//todo проверка существует ли поле
	query := "UPDATE info SET jwt_Token = ?"
	_, err := db.db.QueryContext(ctx, query, JWTToken)
	if err != nil {
		return err
	}
	return nil
}
