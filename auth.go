package goa

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	firebase "firebase.google.com/go"

	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (config *Config) Authenticate(req *http.Request, accountIDs []string) (authorizedAccountIDs []string, err error) {
	urlQuery := req.URL.Query()

	// Get API Key from headers or from URL parameters
	apiKey := req.Header.Get("X-API-KEY")
	if apiKey == "" {
		apiKey = urlQuery.Get("api_key")
	}

	if apiKey != "" {
		authorizedAccountIDs, err = config.AuthorizeApiKey(apiKey, accountIDs)
		if err != nil {
			return
		}
	} else {
		var userEmail string

		// Verify GIP ID Token and get the current user ID
		userEmail, err = VerifyToken(req.Header.Get("Authorization"))
		if err != nil {
			return
		}

		// Find account IDs and permissions for current user
		authorizedAccountIDs, err = config.AuthorizeToken(userEmail, accountIDs)
	}

	if len(authorizedAccountIDs) == 0 {
		err = errors.New("accounts not found")
	}

	return
}

// AuthorizeApiKey ...
func (config *Config) AuthorizeApiKey(apiKey string, accountIDs []string) (authorizedAccountIDs []string, err error) {
	// query := `
	// 	SELECT
	// 		api_keys.account_id as account_id
	// 	FROM api_keys
	// 	LEFT OUTER JOIN roles ON roles.role_name = api_keys.role_name
	// 	WHERE api_keys.api_key = $1
	// 		AND api_keys.account_id IN UNNEST ($2)
	// 		AND LOWER(roles.service_path) = $3
	// 		AND UPPER(roles.service_method) = $4
	// `

	query := psql.Select("account_id").From("api_keys").
		Join("roles USING (role_name)").
		Where(sq.Eq{"api_key": apiKey}).
		Where(sq.Eq{"LOWER(roles.service_path)": config.Service.Path}).
		Where(sq.Eq{"UPPER(roles.service_method)": config.Service.Method})

	// Find specific account ids when accountIDs present
	if len(accountIDs) > 0 {
		query = query.Where(sq.Eq{"account_id": accountIDs})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return
	}

	rows, err := DB.Query(context.Background(), sql, args)
	if err != nil {
		return
	}

	for rows.Next() {
		var accountID string

		err = rows.Scan(&accountID)
		if err != nil {
			return
		}

		authorizedAccountIDs = append(authorizedAccountIDs, accountID)

		if rows.Err() != nil {
			err = rows.Err()
			return
		}
	}

	return
}

func (config *Config) AuthorizeToken(userEmail string, accountIDs []string) (authorizedAccountIDs []string, err error) {
	// query := `
	// 	SELECT
	// 		Users.AccountId as account_id
	// 	FROM Users
	// 	LEFT OUTER JOIN Roles ON Roles.RoleName = Users.RoleName
	// 	WHERE Users.UserEmail = @user_email
	// 		AND Users.AccountId IN UNNEST (@account_ids)
	// 		AND LOWER(Roles.ServicePath) = @service_path
	// 		AND UPPER(Roles.ServiceMethod) = @service_method
	// `

	query := psql.Select("account_id").From("users").
		Join("LEFT OUTER JOIN roles ON roles.role_name = users.role_name").
		Where(sq.Eq{"user_email": userEmail}).
		Where(sq.Eq{"LOWER(roles.service_path)": config.Service.Path}).
		Where(sq.Eq{"UPPER(roles.service_method)": config.Service.Method})

	// Find specific account ids when accountIDs present
	if len(accountIDs) > 0 {
		query = query.Where(sq.Eq{"account_id": accountIDs})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return
	}

	rows, err := DB.Query(context.Background(), sql, args)
	if err != nil {
		return
	}

	for rows.Next() {
		var accountID string

		err = rows.Scan(&accountID)
		if err != nil {
			return
		}

		authorizedAccountIDs = append(authorizedAccountIDs, accountID)

		if rows.Err() != nil {
			err = rows.Err()
			return
		}
	}

	return
}

// VerifyToken ...
func VerifyToken(jwtToken string) (userID string, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, "  Execption in VerifyToken: ", string(debug.Stack()))
		}
	}()

	tokenParts := strings.Split(jwtToken, " ")
	if len(tokenParts) > 1 {
		jwtToken = strings.TrimSpace(tokenParts[1])
	} else {
		jwtToken = strings.TrimSpace(tokenParts[0])
	}

	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Println("VerifyToken firebase.NewApp err: ", err)
		return "", err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Println("VerifyToken app.Auth err: ", err)
		return "", err
	}

	token, err := client.VerifyIDToken(ctx, jwtToken)
	if err != nil {
		log.Println("VerifyToken client.VerifyIDToken err: ", err)
		return "", err
	}

	// token.Claims : map[
	// 	auth_time:1.610686521e+09
	// 	email:someone@gmail.com
	// 	email_verified:false
	// 	firebase:map[
	// 		identities:map[
	// 			email:[someone@gmail.com]
	// 			phone:[+919293949596]
	// 		]
	// 		sign_in_provider:password
	// 	]
	// 	name:giri giri
	// 	phone_number:+919293949596
	// 	user_id:8SyEz1pSizc1IrIK36zigfy6Ou12
	// ]

	return fmt.Sprintf("%v", token.Claims["email"]), nil
}
