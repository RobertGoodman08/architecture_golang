package contact

import (
	"go.uber.org/zap"

	"architecture_go/services/contact/internal/useCase/adapters/storage"

	log "architecture_go/pkg/type/logger"
)

type UseCase struct {
	adapterStorage storage.Contact
	options        Options
}

type Options struct{}

func New(storage storage.Contact, options Options) *UseCase {
	var uc = &UseCase{
		adapterStorage: storage,
	}
	uc.SetOptions(options)
	return uc
}

func (uc *UseCase) SetOptions(options Options) {
	if uc.options != options {
		uc.options = options
		log.Info("set new options", zap.Any("options", uc.options))
	}
}
