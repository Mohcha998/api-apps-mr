package repository

import (
	"apps/internal/domain"
	"apps/internal/infrastructure/db"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type userRepository struct {
	db db.MysqlDBInterface
	httpClient *http.Client
}

func NewUserRepository(db db.MysqlDBInterface, ) *userRepository {
	return &userRepository{
		db: db,
		httpClient: &http.Client{Timeout: time.Second*10},
	}
}

func (r *userRepository) Login(ctx context.Context, user *domain.User) (*domain.User, error) {
	var result domain.User
	var query db.Query
	if user.Email != "" {
		query = db.NewQuery("email = ?", user.Email)
	} else if user.Mobile != "" {
		query = db.NewQuery("mobile = ?", user.Mobile)
	} else {
		return nil, errors.New("email or phone must be provided")
	}
	
	if err := r.db.FindOne(ctx, &result, db.WithQuery(query), db.WithOrder("date_created DESC")); err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	var user_token domain.UserToken
	
	if err := r.db.Create(ctx, user); err != nil {
		return err
	}

	user_token.Token = base64.StdEncoding.EncodeToString([]byte(r.RandomString(32)))
	user_token.Email = user.Email
	user_token.DateCreated = time.Now().Unix()
	

	if err := r.db.Create(ctx, &user_token); err != nil {
		return err
	}

	apiURL := "https://apps.mri.co.id/Api_email/sendMail"
	_, err := r.CURLEmail(ctx, apiURL, &user_token, user.Name)
	if err != nil {
		return err
	}

	
	
	return nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.Update(ctx, user)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	query := db.NewQuery("email = ?", email)
	if err := r.db.FindOne(ctx, &user, db.WithQuery(query), db.WithOrder("date_created DESC")); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByPhone(ctx context.Context, mobile string) (*domain.User, error) {
	var user domain.User
	query := db.NewQuery("mobile = ?", mobile)
	if err := r.db.FindOne(ctx, &user, db.WithQuery(query), db.WithOrder("date_created DESC")); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CURLEmail(ctx context.Context, reqURL string, user_token *domain.UserToken, name string) (string, error) {

	form := url.Values{}
	form.Add("name", name)
	form.Add("email", user_token.Email)
	form.Add("token", user_token.Token)
	form.Add("type", "verify")
	
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("youtube api error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Response body:", string(body))

	return string(body), nil
}

func (r *userRepository) RandomString(length int64) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())

	randomStr := make([]byte, length)
	for i := range randomStr {
		randomStr[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomStr)
}

