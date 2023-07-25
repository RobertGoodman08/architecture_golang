package group

import (
	"go.uber.org/zap"

	log "architecture_go/pkg/type/logger"
	"architecture_go/services/contact/internal/useCase/adapters/storage"
)

type UseCase struct {
	adapterStorage storage.Group
	options        Options
}

type Options struct{}

func New(storage storage.Group, options Options) *UseCase {
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
