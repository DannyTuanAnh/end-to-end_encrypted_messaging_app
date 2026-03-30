package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/grpc/auth"
	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/grpc/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/idtoken"
)

type authService struct {
	auth_proto.UnimplementedAuthServiceServer
	auth_repo    repository.AuthRepository
	user_client  *client.UserClient
	redis_memory *redis.Client
}

func NewAuthService(auth_repo repository.AuthRepository, user_client *client.UserClient, rdb *redis.Client) *authService {
	return &authService{
		auth_repo:    auth_repo,
		user_client:  user_client,
		redis_memory: rdb,
	}
}

func (s *authService) LoginGoogle(ctx context.Context, req *auth_proto.LoginRequest) (*auth_proto.LoginResponse, error) {
	tokenResp, err := s.ExchangeGoogleCode(req.AuthorCode)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to exchange Google code", utils.ErrCodeInternal)
	}

	userInfo, err := s.VerifyIDGoogleToken(ctx, tokenResp.IdToken, tokenResp.AccessToken)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to verify Google ID token", utils.ErrCodeInternal)
	}

	deviceID := uuid.New()

	name := fmt.Sprintf("%s %s", userInfo.Claims["family_name"].(string), userInfo.Claims["given_name"].(string))

	resp, err := s.auth_repo.Login(ctx, sqlc.OAuthLoginParams{
		PProvider:       "google",
		PProviderUserID: userInfo.Claims["sub"].(string),
		PEmail:          userInfo.Claims["email"].(string),
		PDisplayName:    name,
		PDeviceID:       deviceID,
	})
	if err != nil {
		return nil, utils.WrapError(err, "Failed to login with Google", utils.ErrCodeInternal)
	}

	if !resp.ProfileExists {
		_, err := s.user_client.Client.CreateProfile(ctx, &user_proto.CreateProfileRequest{
			UserId:    resp.UserId,
			Name:      name,
			Email:     userInfo.Claims["email"].(string),
			Birthday:  userInfo.Claims["birthday"].(string),
			AvatarUrl: userInfo.Claims["picture"].(string),
		})
		if err != nil {
			return nil, utils.WrapError(err, "Failed to create user profile", utils.ErrCodeInternal)
		}

	}
	resp.DeviceID = deviceID

	version, err := utils.GetKeyRedisAndConvertToInt(ctx, fmt.Sprintf("user:%d:session_version", resp.UserId), s.redis_memory)
	if err != nil {
		log.Println("Error in get session_version (in service layer): ", err)
	}

	if version == 0 {
		if err := s.redis_memory.SetNX(ctx, fmt.Sprintf("user:%d:session_version", resp.UserId), 1, 0).Err(); err != nil {
			log.Println("Error in create session_version if redis didn't exist session_version before (in service layer): ", err)
		}
	}

	session := models.SessionRedis{
		UserID:         resp.UserId,
		DeviceID:       resp.DeviceID,
		SessionVersion: version,
		Valid:          true,
	}

	sessionJson, err := json.Marshal(session)
	if err != nil {
		log.Println("Error in marshal session data (in service layer): ", err)
	}

	err = s.redis_memory.Set(ctx, fmt.Sprintf("session:%s", resp.SessionId.String()), sessionJson, 24*7*time.Hour).Err()
	if err != nil {
		log.Println("Error in set session with marshal data in Redis (in service layer): ", err)
	}

	return &auth_proto.LoginResponse{
		Session:  resp.SessionId.String(),
		UserId:   resp.UserId,
		DeviceId: resp.DeviceID.String(),
	}, nil
}

func (s *authService) ExchangeGoogleCode(code string) (*models.GoogleTokenResponse, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", utils.GetEnv("GOOGLE_CLIENT_ID", ""))
	data.Set("client_secret", utils.GetEnv("GOOGLE_CLIENT_SECRET", ""))
	data.Set("redirect_uri", utils.GetEnv("GOOGLE_REDIRECT_URI", ""))
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(utils.GetEnv("GOOGLE_TOKEN_URL", ""), data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result models.GoogleTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *authService) VerifyIDGoogleToken(ctx context.Context, idToken string, accessToken string) (*idtoken.Payload, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to verify ID token: %v", err)
	}
	defer resp.Body.Close()

	payload, err := idtoken.Validate(ctx, idToken, utils.GetEnv("GOOGLE_CLIENT_ID", ""))
	if err != nil {
		return nil, fmt.Errorf("Failed to validate ID token: %v", err)
	}

	err = s.GetBirthday(payload, accessToken)
	if err != nil {
		log.Printf("Failed to get birthday: %v", err)
		return nil, err
	}

	return payload, nil
}

func (s *authService) GetBirthday(userInfo *idtoken.Payload, accessToken string) error {
	req, err := http.NewRequest("GET", "https://people.googleapis.com/v1/people/me?personFields=birthdays", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var birthdayResp models.GoogleBirthdayResponse
	err = json.NewDecoder(resp.Body).Decode(&birthdayResp)
	if err != nil {
		return err
	}

	if len(birthdayResp.Birthdays) == 0 {
		return nil
	}

	b := birthdayResp.Birthdays[0].Date
	userInfo.Claims["birthday"] = fmt.Sprintf("%04d-%02d-%02d", b.Year, b.Month, b.Day)

	return nil
}

func (s *authService) Logout(ctx context.Context, req *auth_proto.LogoutRequest) (*auth_proto.LogoutResponse, error) {
	sessionId, err := uuid.Parse(req.SessionId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to parse session ID", utils.ErrCodeInternal)
	}

	deviceId, err := uuid.Parse(req.DeviceId)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to parse device ID", utils.ErrCodeInternal)
	}

	params := sqlc.RevokeSessionParams{
		SessionID: sessionId,
		DeviceID:  deviceId,
	}

	err = s.auth_repo.Logout(ctx, params)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to logout", utils.ErrCodeInternal)
	}

	if err := s.redis_memory.Del(ctx, fmt.Sprintf("session:%s", req.SessionId)).Err(); err != nil {
		log.Println("Error in delete session in Redis (in service layer): ", err)
	}

	return &auth_proto.LogoutResponse{
		Success: true,
		Message: "Logout successfully",
	}, nil
}

func (s *authService) LogoutAll(ctx context.Context, req *auth_proto.LogoutAllRequest) (*auth_proto.LogoutAllResponse, error) {
	if req.UserId == "" {
		return nil, utils.NewError("Invalid user ID", utils.ErrCodeBadRequest)
	}

	userIdInt, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to parse user ID", utils.ErrCodeInternal)
	}

	err = s.auth_repo.LogoutAll(ctx, userIdInt)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to logout all sessions", utils.ErrCodeInternal)
	}

	if err := s.redis_memory.Incr(ctx, fmt.Sprintf("user:%d:session_version", userIdInt)).Err(); err != nil {
		log.Println("Error in increment session_version in Redis (in service layer): ", err)
	}

	return &auth_proto.LogoutAllResponse{
		Success: true,
		Message: "Logout all sessions successfully",
	}, nil
}
