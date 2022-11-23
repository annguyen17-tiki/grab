package service

import "github.com/annguyen17-tiki/grab/internal/model"

func (svc *service) SaveToken(accountID, token string) error {
	return svc.store.Token().Save(&model.Token{AccountID: accountID, Token: token})
}
