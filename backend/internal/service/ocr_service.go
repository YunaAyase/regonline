package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
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
	OCRProviderBaidu   OCRProvider = "baidu"
	OCRProviderAlibaba OCRProvider = "alibaba"
)

type OCRConfig struct {
	Provider            string `json:"provider"`
	BaiduAPIKey       string `json:"baidu_api_key"`
	BaiduSecretKey    string `json:"baidu_secret_key"`
	AlibabaAccessKeyId     string `json:"alibaba_access_key_id"`
	AlibabaAccessKeySecret string `json:"alibaba_access_key_secret"`
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
	if config.Provider == "alibaba" {
		return config.AlibabaAccessKeyId != "" && config.AlibabaAccessKeySecret != ""
	}
	return config.BaiduAPIKey != "" && config.BaiduSecretKey != ""
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

	baiduKey := result["ocr_baidu_api_key"]
	baiduSecret := result["ocr_baidu_secret_key"]
	alibabaId := result["ocr_alibaba_access_key_id"]
	alibabaSecret := result["ocr_alibaba_access_key_secret"]

	// 兼容旧字段：新字段为空时回退到旧字段
	if baiduKey == "" {
		baiduKey = result["ocr_api_key"]
	}
	if baiduSecret == "" {
		baiduSecret = result["ocr_secret_key"]
	}

	return OCRConfig{
		Provider:              result["ocr_provider"],
		BaiduAPIKey:           baiduKey,
		BaiduSecretKey:        baiduSecret,
		AlibabaAccessKeyId:    alibabaId,
		AlibabaAccessKeySecret: alibabaSecret,
	}
}

func (s *OCRService) RecognizeAllNumbers(photoData []byte) ([]OCRNumberCandidate, error) {
	config := s.getConfig()
	provider := OCRProvider(config.Provider)
	if provider == "" {
		provider = OCRProviderBaidu
	}

	switch provider {
	case OCRProviderAlibaba:
		if config.AlibabaAccessKeyId == "" || config.AlibabaAccessKeySecret == "" {
			return nil, fmt.Errorf("未配置阿里云 OCR AccessKey，请在管理后台设置")
		}
		return s.recognizeAlibaba(photoData, config)
	default:
		if config.BaiduAPIKey == "" || config.BaiduSecretKey == "" {
			return nil, fmt.Errorf("未配置百度云 OCR API Key，请在管理后台设置")
		}
		return s.recognizeBaidu(photoData, config)
	}
}

// ======================== 百度云 OCR ========================

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
		config.BaiduAPIKey,
		config.BaiduSecretKey,
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

// ======================== 阿里云 OCR (V3 签名) ========================

func (s *OCRService) recognizeAlibaba(photoData []byte, config OCRConfig) ([]OCRNumberCandidate, error) {
	endpoint := "ocr-api.cn-hangzhou.aliyuncs.com"
	action := "RecognizeIdcard"
	apiVersion := "2021-07-07"

	// 1. 准备请求体（二进制图片数据）
	payload := photoData
	contentType := "application/octet-stream"

	// 2. 准备请求头
	now := time.Now().UTC()
	acsDate := now.Format("2006-01-02T15:04:05Z")
	nonce := generateACSNonce()
	payloadHash := computeSHA256Hex(payload)

	headers := map[string]string{
		"host":                      endpoint,
		"content-type":              contentType,
		"x-acs-action":              action,
		"x-acs-content-sha256":     payloadHash,
		"x-acs-date":                acsDate,
		"x-acs-signature-nonce":     nonce,
		"x-acs-version":             apiVersion,
	}

	// 3. 构造规范化请求 CanonicalRequest
	signedHeaders := []string{"content-type", "host", "x-acs-action", "x-acs-content-sha256", "x-acs-date", "x-acs-signature-nonce", "x-acs-version"}

	canonicalHeaders := ""
	for _, name := range signedHeaders {
		canonicalHeaders += fmt.Sprintf("%s:%s\n", name, strings.TrimSpace(headers[name]))
	}

	canonicalQueryString := ""
	canonicalURI := "/"

	canonicalRequest := "POST\n" +
		canonicalURI + "\n" +
		canonicalQueryString + "\n" +
		canonicalHeaders + "\n" +
		strings.Join(signedHeaders, ";") + "\n" +
		payloadHash

	// 4. 构造待签名字符串
	hashedCanonical := computeSHA256Hex([]byte(canonicalRequest))
	stringToSign := "ACS3-HMAC-SHA256\n" + hashedCanonical

	// 5. 计算签名
	signature := computeHMACSHA256Hex(config.AlibabaAccessKeySecret, stringToSign)

	// 6. 构造 Authorization
	authorization := fmt.Sprintf(
		"ACS3-HMAC-SHA256 Credential=%s,SignedHeaders=%s,Signature=%s",
		config.AlibabaAccessKeyId,
		strings.Join(signedHeaders, ";"),
		signature,
	)

	// 7. 发起请求
	reqURL := fmt.Sprintf("https://%s/", endpoint)
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("创建阿里云 OCR 请求失败: %w", err)
	}

	req.Host = endpoint
	for k, v := range headers {
		if k == "host" {
			continue // 由 req.Host 控制，避免 Go 发送重复/冲突的 Host 头
		}
		req.Header.Set(k, v)
	}
	req.Header.Set("Authorization", authorization)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("阿里云 OCR 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取阿里云 OCR 响应失败: %w", err)
	}

	log.Printf("Alibaba OCR canonical request:\n%s", canonicalRequest)
	log.Printf("Alibaba OCR stringToSign:\n%s", stringToSign)
	log.Printf("Alibaba OCR signature: %s", signature)
	log.Printf("Alibaba OCR authorization: %s", authorization)
	log.Printf("Alibaba OCR response (status=%d): %s", resp.StatusCode, string(respBody))

	// 8. 检查 HTTP 状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if len(respBody) > 0 && respBody[0] == '<' {
			return nil, fmt.Errorf("阿里云网关错误 (HTTP %d): %s", resp.StatusCode, extractAliyunXMLError(respBody))
		}
		return nil, fmt.Errorf("阿里云 OCR HTTP 错误 (status=%d): %s", resp.StatusCode, string(respBody))
	}

	// 9. 解析响应
	var apiResp struct {
		RequestId string `json:"RequestId"`
		Data      string `json:"Data"`
		Code      string `json:"Code"`
		Message   string `json:"Message"`
	}

	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("解析阿里云 OCR 响应失败: %w", err)
	}

	if apiResp.Code != "" {
		return nil, fmt.Errorf("阿里云 OCR 错误: %s - %s", apiResp.Code, apiResp.Message)
	}

	if apiResp.Data == "" {
		return nil, fmt.Errorf("阿里云 OCR 返回数据为空")
	}

	// 10. 解析 Data 字段内的 JSON 内容
	var dataContent struct {
		Data struct {
			Face struct {
				Data struct {
					IDNumber  string `json:"idNumber"`
					BirthDate string `json:"birthDate"`
				} `json:"data"`
			} `json:"face"`
			Back struct {
				Data struct {
					IssueAuthority string `json:"issueAuthority"`
					ValidPeriod  string `json:"validPeriod"`
				} `json:"data"`
			} `json:"back"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(apiResp.Data), &dataContent); err != nil {
		// 有些时候 Data 已经是一个 JSON 字符串包含 data 字段
		var altContent struct {
			Face struct {
				Data struct {
					IDNumber  string `json:"idNumber"`
					BirthDate string `json:"birthDate"`
				} `json:"data"`
			} `json:"face"`
		}
		if err2 := json.Unmarshal([]byte(apiResp.Data), &altContent); err2 != nil {
			return nil, fmt.Errorf("解析阿里云 OCR Data 字段失败: %w (原始内容: %s)", err, apiResp.Data)
		}
		dataContent.Data.Face.Data = altContent.Face.Data
	}

	var candidates []OCRNumberCandidate

	// 优先从 face 提取
	if idNum := strings.TrimSpace(dataContent.Data.Face.Data.IDNumber); idNum != "" {
		normalized := util.NormalizeIDNumber(idNum)
		candidates = append(candidates, OCRNumberCandidate{
			Value:      normalized,
			Label:      "身份证号",
			IsIDNumber: true,
			Length:     len(normalized),
		})
	}

	if birth := strings.TrimSpace(dataContent.Data.Face.Data.BirthDate); birth != "" {
		candidates = append(candidates, OCRNumberCandidate{
			Value:      birth,
			Label:      "出生日期",
			IsIDNumber: false,
			Length:     len(birth),
		})
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("未识别到身份证信息，请确认照片清晰且包含身份证正面")
	}

	return candidates, nil
}

// ======================== 工具函数 ========================

func extractAliyunXMLError(body []byte) string {
	bodyStr := string(body)
	// 尽量从 XML 中提取 Code 和 Message
	var code, msg string
	if idx := strings.Index(bodyStr, "<Code>"); idx >= 0 {
		end := strings.Index(bodyStr[idx:], "</Code>")
		if end > 0 {
			code = strings.TrimSpace(bodyStr[idx+6 : idx+end])
		}
	}
	if idx := strings.Index(bodyStr, "<Message>"); idx >= 0 {
		end := strings.Index(bodyStr[idx:], "</Message>")
		if end > 0 {
			msg = strings.TrimSpace(bodyStr[idx+9 : idx+end])
		}
	}
	if code != "" || msg != "" {
		return fmt.Sprintf("%s: %s", code, msg)
	}
	return strings.TrimSpace(bodyStr)
}

func computeSHA256Hex(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func computeHMACSHA256Hex(secret, data string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func generateACSNonce() string {
	const letters = "abcdef0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
