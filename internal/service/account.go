package service

import (
	"time"

	"github.com/annguyen17-tiki/grab/internal/dto"
	"github.com/annguyen17-tiki/grab/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (svc *service) CreateAccount(input *dto.CreateAccount) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), svc.cfg.BcryptCost)
	if err != nil {
		return err
	}

	input.Password = string(hash)

	var account model.Account
	if err := copier.Copy(&account, input); err != nil {
		return err
	}

	return svc.store.Account().Create(&account)
}

func (svc *service) UpdateAccount(input *dto.UpdateAccount) error {
	var account model.Account
	if err := copier.Copy(&account, input); err != nil {
		return err
	}

	return svc.store.Account().Update(&account)
}

func (svc *service) Login(input *dto.LoginInput) (string, error) {
	account, err := svc.store.Account().Get(&model.Account{Username: input.Username})
	if err != nil {
		return "", model.NewErrUnauthorized("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(input.Password)); err != nil {
		return "", model.NewErrUnauthorized("invalid username or password")
	}

	claim := jwt.StandardClaims{
		Subject:   account.ID,
		ExpiresAt: time.Now().Add(time.Duration(svc.cfg.ValidTokenHour) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(svc.cfg.JWTSecret))
}

func (svc *service) GetAccount(id string) (*model.Account, error) {
	return svc.store.Account().Get(&model.Account{ID: id}, "Driver")
}

func (svc *service) GetAccountByPhone(phone string) (*model.Account, error) {
	account, err := svc.store.Account().Get(&model.Account{Phone: phone})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewErrNotFound("not found account")
		}
		return nil, err
	}
	return account, nil
}
