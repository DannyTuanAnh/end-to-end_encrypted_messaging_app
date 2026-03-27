package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/client"
	// "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/db/sqlc"
	auth_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/grpc/auth"
	user_proto "github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/grpc/user"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/models"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/repository"
	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
	// "github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type authService struct {
	auth_proto.UnimplementedAuthServiceServer
	auth_repo   repository.AuthRepository
	user_client *client.UserClient
}

func NewAuthService(auth_repo repository.AuthRepository, user_client *client.UserClient) *authService {
	return &authService{
		auth_repo:   auth_repo,
		user_client: user_client,
	}
}

func (s *authService) LoginGoogle(ctx context.Context, req *auth_proto.LoginRequest) (*auth_proto.LoginResponse, error) {
	// tokenResp, err := s.ExchangeGoogleCode(req.AuthorCode)
	// if err != nil {
	// 	return nil, err
	// }

	// userInfo, err := s.VerifyIDGoogleToken(ctx, tokenResp.IdToken, tokenResp.AccessToken)
	// if err != nil {
	// 	return nil, err
	// }

	// deviceID, err := uuid.NewUUID()
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to generate device ID: %v", err)
	// }

	// name := fmt.Sprintf("%s %s", userInfo.Claims["family_name"].(string), userInfo.Claims["given_name"].(string))

	// resp, err := s.auth_repo.Login(ctx, sqlc.OAuthLoginParams{
	// 	PProvider:       "google",
	// 	PProviderUserID: userInfo.Claims["sub"].(string),
	// 	PEmail:          userInfo.Claims["email"].(string),
	// 	PDisplayName:    name,
	// 	PDeviceID:       deviceID,
	// })
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to login user: %v", err)
	// }

	// if !resp.ProfileExists {
	// 	_, err := s.user_client.Client.CreateProfile(ctx, &user_proto.CreateProfileRequest{
	// 		UserId:    resp.UserId,
	// 		Name:      name,
	// 		Email:     userInfo.Claims["email"].(string),
	// 		Birthday:  userInfo.Claims["birthday"].(string),
	// 		AvatarUrl: userInfo.Claims["picture"].(string),
	// 	})
	// 	if err != nil {
	// 		return nil, fmt.Errorf("Failed to create user profile: %v", err)
	// 	}

	// 	resp.DeviceID = deviceID

	// }
	// return &auth_proto.LoginResponse{
	// 	Session:  resp.SessionId.String(),
	// 	UserId:   resp.UserId,
	// 	DeviceId: resp.DeviceID.String(),
	// }, nil
	_, err := s.user_client.Client.CreateProfile(ctx, &user_proto.CreateProfileRequest{})
	if err != nil {
		log.Printf("Failed to create profile: %v", err)
		return nil, err
	}

	return nil, nil
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
		log.Printf("Failed to exchange auth code for token: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result models.GoogleTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Printf("Failed to decode token response: %v", err)
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
		log.Printf("Failed to create request for birthday: %v", err)
		return err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Failed to get birthday: %v", err)
		return err
	}
	defer resp.Body.Close()

	var birthdayResp models.GoogleBirthdayResponse
	err = json.NewDecoder(resp.Body).Decode(&birthdayResp)
	if err != nil {
		log.Printf("Failed to decode birthday response: %v", err)
		return err
	}

	if len(birthdayResp.Birthdays) == 0 {
		log.Println("No birthday information found for user")
		return nil
	}

	b := birthdayResp.Birthdays[0].Date
	userInfo.Claims["birthday"] = fmt.Sprintf("%04d-%02d-%02d", b.Year, b.Month, b.Day)

	return nil
}
