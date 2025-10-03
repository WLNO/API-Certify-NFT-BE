package service

import "api-certify-nft-be/internal/repository"

type AuthService interface {
	Login(wallet string) (isNew bool, role string, err error)
	RegisterUser(name, email, wallet string) (id int, created, updated string, err error)
	RegisterVendor(vendorName, email, contact, wallet string) (int, error)
	IsWalletRegistered(wallet string) (bool, error)
}

type authSvc struct{ repo repository.AuthRepo }

func NewAuthService(r repository.AuthRepo) AuthService { return &authSvc{repo: r} }

func (s *authSvc) Login(wallet string) (bool, string, error) {
	u, err := s.repo.IsUserWallet(wallet)
	if err != nil {
		return false, "", err
	}
	if u {
		return false, "users", nil
	}
	v, err := s.repo.IsVendorWallet(wallet)
	if err != nil {
		return false, "", err
	}
	if v {
		return false, "vendors", nil
	}
	return true, "", nil
}
func (s *authSvc) RegisterUser(name, email, wallet string) (int, string, string, error) {
	return s.repo.CreateUser(name, email, wallet)
}
func (s *authSvc) RegisterVendor(vn, email, contact, wallet string) (int, error) {
	return s.repo.CreateVendor(vn, email, contact, wallet)
}
func (s *authSvc) IsWalletRegistered(wallet string) (bool, error) {
	u, err := s.repo.IsUserWallet(wallet)
	if err != nil {
		return false, err
	}
	if u {
		return true, nil
	}
	v, err := s.repo.IsVendorWallet(wallet)
	if err != nil {
		return false, err
	}
	return v, nil
}
