package service

import "context"

func (a *AuthService) Register(ctx context.Context, name, password string) error {
	return nil
}

func (a *AuthService) Login(ctx context.Context, name, password string) (int, error) {
	return 0, nil
}
