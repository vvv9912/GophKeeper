package sqllite

import "context"

func (db *Database) GetJWTToken(ctx context.Context) (string, error) {
	q := "Select jwt_Token from info"
	var jwt string
	err := db.db.QueryRowxContext(ctx, q).Scan(&jwt)
	if err != nil {
		return "", err
	}
	return jwt, nil
}
func (db *Database) SetJWTToken(ctx context.Context, JWTToken string) error {
	// проверка существует ли поле
	q := "UPDATE info SET jwt_Token = ?"
	_, err := db.db.QueryContext(ctx, q, JWTToken)
	if err != nil {
		return err
	}
	return nil
}
