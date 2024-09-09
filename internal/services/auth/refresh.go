package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"test-auth/internal/models"
)

// RefreshHandler         godoc
// @Summary      		  Refresh access token
// @Description  		  Refresh access token endpoint for (mobile, TV)
// @Tags         		  Auth
// @Accept       		  json
// @Produce      		  json
// @Param 				  req             body       models.RefreshTokenRequest true "refresh token request"
// @Param                 X-Forwarded-For header     string                     true  "IP address"
// @Success      		  200             {object}   models.RefreshTokenResponse
// @Failure               400             {object}   models.ErrorResponse
// @Failure               401             {object}   models.ErrorResponse
// @Failure               500             {object}   models.ErrorResponse
// @Router       	      /auth/refresh-mobile [post]
func (s *Service) RefreshHandler(ctx *gin.Context) {
	var req models.RefreshTokenRequest
	if err := ctx.BindJSON(&req); err != nil {
		s.Logger.Error("auth.RefreshHandler: unmarshal request body", zap.Error(err))
		sendErrorResponse(ctx, "Invalid request body", http.StatusBadRequest)
		return
	}

	accessClaims, err := s.tokenManager.ParseAccessToken(req.AccessToken)
	if err != nil {
		s.Logger.Error("auth.RefreshHandler: parse access token", zap.Error(err))
		sendErrorResponse(ctx, "Invalid access token", http.StatusBadRequest)
		return
	}

	refreshClaims, err := s.tokenManager.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		s.Logger.Error("auth.RefreshHandler: parse refresh token", zap.Error(err))
		sendErrorResponse(ctx, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	if accessClaims.ID != refreshClaims.ID {
		s.Logger.Error("auth.RefreshHandler: invalid token pair")
		sendErrorResponse(ctx, "Invalid token pair", http.StatusBadRequest)
		return
	}

	refreshTokenHash, err := hashRefreshToken(req.RefreshToken)
	if err != nil {
		s.Logger.Error("auth.RefreshHandler: hash refresh token", zap.Error(err))
		sendErrorResponse(ctx, "Internal server error", http.StatusInternalServerError)
		return
	}

	hashIsValid, err := s.PgRepo.CheckRefreshTokenHash(ctx, accessClaims.UserID, refreshTokenHash)
	if err != nil {
		s.Logger.Error("auth.RefreshHandler: check refresh token hash", zap.Error(err))
		sendErrorResponse(ctx, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !hasIsValid {
		s.Logger.Error("auth.RefreshHandler: invalid refresh token")
		sendErrorResponse(ctx, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	ips := strings.Split(ctx.GetHeader("X-Forwarded-For"), ",")
	userIP := ips[0]
	if userIP == "" {
		s.Logger.Error("auth.RefreshHandler: empty user ip")
		sendErrorResponse(ctx, "Invalid user ip", http.StatusBadRequest)
		return
	}

	tokenId := uuid.NewString()

	newAccessToken, err := s.tokenManager.NewAccessToken(accessClaims.UserID, tokenId, userIP, s.config.Token.AccessTTL)
	if err != nil {
		s.Logger.Error("auth.RefreshHandler: create new access token", zap.Error(err))
		sendErrorResponse(ctx, "Internal server error", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := s.tokenManager.NewRefreshToken(tokenId, userIP, s.config.Token.RefreshTTL)
	if err != nil {
		s.Logger.Error("auth.RefreshHandler: create new access token", zap.Error(err))
		sendErrorResponse(ctx, "Internal server error", http.StatusInternalServerError)
		return
	}

	newRefreshTokenHash, err = hashRefreshToken(newRefreshToken)
	if err != nil {
		s.Logger.Error("auth.RefreshHandler: hash refresh token failed", zap.Error(err))
		sendErrorResponse(ctx, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = s.PgRepo.UpdateTokenHash(ctx, accessClaims.UserID, newRefreshTokenHash)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("auth.RefreshHandler: update refresh token failed, for user id: %s", userId), zap.Error(err))
		if err.Error() == "user not found" {
			sendErrorResponse(ctx, "User not found", http.StatusNotFound)
			return
		}

		sendErrorResponse(ctx, "Internal server error", http.StatusInternalServerError)
		return
	}

	//TODO need to add check err for sql.No rows
	if accessClaims.IP != userIP {
		email := ""
		email, err = s.PgRepo.GetEmailByUserId(ctx, accessClaims.UserID)
		if err != nil {
			s.Logger.Error("auth.RefreshHandler: get email by user id", zap.Error(err))
			sendErrorResponse(ctx, "Internal server error", http.StatusInternalServerError)
			return
		}

		text := fmt.Sprintf("Previous IP address was - %s, now IP address is - %s", accessClaims.IP, userIP)
		err = s.smtp.SendNotification(ctx, email, text)
		if err != nil {
			s.Logger.Error("auth.RefreshHandler: send email notification", zap.Error(err))
			sendErrorResponse(ctx, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	resp := models.RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	sendSuccessResponse(ctx, resp, http.StatusOK)
}
