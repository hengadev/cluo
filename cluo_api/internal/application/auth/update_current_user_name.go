package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

func (s *AuthService) UpdateCurrentUserName(ctx context.Context, userID uuid.UUID, name string) error {
	userEncx, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	dek, err := s.crypto.DecryptDEKWithVersion(ctx, userEncx.DEKEncrypted, userEncx.KeyVersion)
	if err != nil {
		return err
	}

	nameBytes, err := encx.SerializeValue(name)
	if err != nil {
		return err
	}

	nameEncrypted, err := s.crypto.EncryptData(ctx, nameBytes, dek)
	if err != nil {
		return err
	}

	return s.userRepo.UpdateUserName(ctx, userID, nameEncrypted)
}
