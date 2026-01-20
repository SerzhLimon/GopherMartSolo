package handler

import (
	"github.com/SerzhLimon/GopherMartSolo/internal/config"
	"github.com/SerzhLimon/GopherMartSolo/pkg/storage"
)

type Handler struct {
	Storage storage.GofferStorage
	JwtKey  []byte
}

func NewHandler(cfg *config.ServerConfig) *Handler {
	storage := storage.NewGofferStorage(cfg.DataBaseDsn)
	return &Handler{
		Storage: *storage,
		JwtKey:  []byte(cfg.SecretKey),
	}
}
