package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"regonline-backend/internal/repository"
	"regonline-backend/internal/util"
)

type OCRProvider string

const (
	OCRProviderBaidu    OCRProvider = "baidu"
	OCRProviderAlibaba  OCRProvider = "alibaba"
)

type OCRConfig struct {
	Provider   string `json:"provider"`
	APIKey     string `json:"api_key"`
	SecretKey  string `json:"secret_key"`
}

type OCRNumberCandidate struct {
	Value      string `json:"value"`
	Label      string `json:"label"`
	IsIDNumber bool   `json:"is_id_number"`
	Length     int    `json:"length"`
}

type OCRService struct {
	settingRepo *repository.SiteSettingRepository
	tokenCache  string
	tokenExpiry time.Time
	mu          sync.Mutex
}

func NewOCRService(settingRepo *repository.SiteSettingRepository) *OCRService {
	return &OCRService{
		settingRepo: settingRepo,
	}
}

func (s *OCRService) IsAvailable() bool {
	config := s.getConfig()
	return config.APIKey != "" && config.SecretKey != ""
}

func (s *OCRService) getConfig() OCRConfig {
	settings, err := s.settingRepo.GetAll()
	if err != nil {
		return OCRConfig{}
	}
	result := make(map[string]string, len(settings))
	for _, st := range settings {
		result[st.Key] = st.Value
	}
	return OCRConfig{
		Provider:  result["ocr_provider"],
		APIKey:    result["ocr_api_key"],
		SecretKey: result["ocr_secret_key"],
	}
}

func (s *OCRService) RecognizeAllNumbers(photoData []byte) ([]OCRNumberCandidate, error) {
	config := s.getConfig()
	if config.APIKey == "" || config.SecretKey == "" {
		return nil, fmt.Errorf("未配置 OCR 云识别服务，请在管理后台设置")
	}

	provider := OCRProvider(config.Provider)
	if provider == "" {
		provider = OCRProviderBaidu
	}

	switch provider {
	case OCRProviderAlibaba:
		return s.recognizeAlibaba(photoData, config)
	default:
		return s.recognizeBaidu(photoData, config)
	}
}

func (s *OCRService) recognizeBaidu(photoData []byte, config OCRConfig) ([]OCRNumberCandidate, error) {
	token, err := s.getBaiduToken(config)
	if err != nil {
		return nil, fmt.Errorf("获取百度 OCR Token 失败: %w", err)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(photoData)

	apiURL := fmt.Sprintf("https://aip.baidubce.com/rest/2.0/ocr/v1/idcard?access_token=%s", token)

	body := url.Values{}
	body.Set("image", imageBase64)
	body.Set("id_card_side", "front")
	body.Set("detect_direction", "true")
	body.Set("detect_risk", "false")

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("百度 OCR 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取百度 OCR 响应失败: %w", err)
	}

	log.Printf("Baidu OCR response: %s", string(respBody))

	var result struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
		WordsResult struct {
			IDNumber struct {
				Words string `json:"words"`
			} `json:"公民身份号码"`
			Name struct {
				Words string `json:"words"`
			} `json:"姓名"`
			Gender struct {
				Words string `json:"words"`
			} `json:"性别"`
			BirthDate struct {
				Words string `json:"words"`
			} `json:"出生"`
			Address struct {
				Words string `json:"words"`
			} `json:"住址"`
		} `json:"words_result"`
		WordsResultNum int `json:"words_result_num"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析百度 OCR 响应失败: %w", err)
	}

	if result.ErrorCode != 0 {
		if result.ErrorCode == 110 || result.ErrorCode == 111 {
			s.mu.Lock()
			s.tokenCache = ""
			s.tokenExpiry = time.Time{}
			s.mu.Unlock()
		}
		return nil, fmt.Errorf("百度 OCR 错误 [%d]: %s", result.ErrorCode, result.ErrorMsg)
	}

	var candidates []OCRNumberCandidate

	if idNum := strings.TrimSpace(result.WordsResult.IDNumber.Words); idNum != "" {
		normalized := util.NormalizeIDNumber(idNum)
		candidates = append(candidates, OCRNumberCandidate{
			Value:      normalized,
			Label:      "身份证号",
			IsIDNumber: true,
			Length:     len(normalized),
		})
	}

	if birth := strings.TrimSpace(result.WordsResult.BirthDate.Words); birth != "" {
		candidates = append(candidates, OCRNumberCandidate{
			Value:      birth,
			Label:      "出生日期",
			IsIDNumber: false,
			Length:     len(birth),
		})
	}

	if len(candidates) == 0 && result.WordsResultNum == 0 {
		return nil, fmt.Errorf("未识别到身份证信息，请确认照片清晰且包含身份证正面")
	}

	return candidates, nil
}

func (s *OCRService) getBaiduToken(config OCRConfig) (string, error) {
	s.mu.Lock()
	if s.tokenCache != "" && time.Now().Before(s.tokenExpiry) {
		token := s.tokenCache
		s.mu.Unlock()
		return token, nil
	}
	s.mu.Unlock()

	authURL := fmt.Sprintf(
		"https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s",
		config.APIKey,
		config.SecretKey,
	)

	resp, err := http.Post(authURL, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return "", fmt.Errorf("请求百度 Token 失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取百度 Token 响应失败: %w", err)
	}

	var tokenResult struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}

	if err := json.Unmarshal(respBody, &tokenResult); err != nil {
		return "", fmt.Errorf("解析百度 Token 响应失败: %w", err)
	}

	if tokenResult.Error != "" {
		return "", fmt.Errorf("百度 Token 错误: %s - %s", tokenResult.Error, tokenResult.ErrorDesc)
	}

	s.mu.Lock()
	s.tokenCache = tokenResult.AccessToken
	expireSeconds := tokenResult.ExpiresIn
	if expireSeconds <= 0 {
		expireSeconds = 2592000
	}
	s.tokenExpiry = time.Now().Add(time.Duration(expireSeconds) * time.Second).Add(-5 * time.Minute)
	s.mu.Unlock()

	log.Println("Baidu OCR access_token acquired successfully")
	return tokenResult.AccessToken, nil
}

func (s *OCRService) recognizeAlibaba(photoData []byte, config OCRConfig) ([]OCRNumberCandidate, error) {
	return nil, fmt.Errorf("阿里云 OCR 暂未实现，请使用百度云 OCR")
}