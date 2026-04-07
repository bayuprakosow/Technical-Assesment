package handlers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zentara/technical_assesment/internal/config"
	"github.com/zentara/technical_assesment/internal/middleware"
	"github.com/zentara/technical_assesment/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type Handlers struct {
	cfg      *config.Config
	db       *sql.DB
	users    *repository.UserRepository
	findings *repository.FindingRepository
}

func New(cfg *config.Config, db *sql.DB) *Handlers {
	return &Handlers{
		cfg:      cfg,
		db:       db,
		users:    repository.NewUserRepository(db),
		findings: repository.NewFindingRepository(db),
	}
}

func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handlers) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := h.db.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "error": "database unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

func (h *Handlers) PublicInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service": h.cfg.ServiceName,
		"version": h.cfg.ServiceVersion,
	})
}

type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *Handlers) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}
	u, err := h.users.Create(c.Request.Context(), strings.ToLower(strings.TrimSpace(req.Email)), string(hash))
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "email": u.Email})
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handlers) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.users.GetByEmail(c.Request.Context(), strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := h.issueToken(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "token_type": "Bearer", "expires_in": 86400})
}

func (h *Handlers) issueToken(userID int64) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(h.cfg.JWTSecret))
}

func (h *Handlers) Me(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	u, err := h.users.GetByID(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": u.ID, "email": u.Email})
}

// UserID helper in middleware - I need to add UserID function in middleware package
type findingDTO struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Severity  string `json:"severity"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func (h *Handlers) ListFindings(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	list, err := h.findings.ListByUser(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list findings"})
		return
	}
	out := make([]findingDTO, 0, len(list))
	for _, f := range list {
		out = append(out, findingDTO{
			ID:        f.ID,
			Title:     f.Title,
			Severity:  f.Severity,
			Status:    f.Status,
			CreatedAt: f.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, gin.H{"findings": out})
}

type createFindingReq struct {
	Title    string `json:"title" binding:"required"`
	Severity string `json:"severity" binding:"required"`
	Status   string `json:"status"`
}

func (h *Handlers) CreateFinding(c *gin.Context) {
	uid, ok := middleware.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req createFindingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f, err := h.findings.Create(c.Request.Context(), uid, strings.TrimSpace(req.Title), strings.TrimSpace(req.Severity), strings.TrimSpace(req.Status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create finding"})
		return
	}
	c.JSON(http.StatusCreated, findingDTO{
		ID:        f.ID,
		Title:     f.Title,
		Severity:  f.Severity,
		Status:    f.Status,
		CreatedAt: f.CreatedAt.UTC().Format(time.RFC3339),
	})
}

func (h *Handlers) InternalMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"uptime_stub": "ok",
		"note":        "placeholder metrics for assessment",
	})
}

func (h *Handlers) InternalPurge(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "cache purge simulated"})
}
