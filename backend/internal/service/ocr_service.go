package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

	"regonline-backend/internal/util"
)

type OCRService struct {
	available bool
}

type OCRNumberCandidate struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	IsIDNumber  bool   `json:"is_id_number"`
	Length      int    `json:"length"`
}

func NewOCRService() *OCRService {
	s := &OCRService{}
	s.checkAvailability()
	return s
}

func (s *OCRService) checkAvailability() {
	s.available = false

	log.Println("OCR service initialization...")

	_, err := exec.LookPath("tesseract")
	if err != nil {
		log.Println("Tesseract not found in PATH, OCR will fall back to manual input")
		return
	}

	log.Println("Tesseract found, OCR service available")
	s.available = true
}

func (s *OCRService) IsAvailable() bool {
	return s.available
}

func (s *OCRService) RecognizeIDNumber(photoData []byte) (string, error) {
	if !s.available {
		return "", fmt.Errorf("OCR service not available")
	}

	tmpFile, err := os.CreateTemp("", "ocr_*.png")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(photoData); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	text, err := s.runTesseract(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("tesseract execution failed: %w", err)
	}

	idNumber := s.extractIDNumber(text)
	if idNumber == "" {
		return "", fmt.Errorf("未识别到身份证号")
	}

	return util.NormalizeIDNumber(idNumber), nil
}

func (s *OCRService) RecognizeAllNumbers(photoData []byte) ([]OCRNumberCandidate, error) {
	if !s.available {
		return nil, fmt.Errorf("OCR service not available")
	}

	tmpFile, err := os.CreateTemp("", "ocr_*.png")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(photoData); err != nil {
		return nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	text, err := s.runTesseract(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("tesseract execution failed: %w", err)
	}

	return s.extractAllNumbers(text), nil
}

func (s *OCRService) runTesseract(imagePath string) (string, error) {
	cmd := exec.Command("tesseract", imagePath, "stdout", "-l", "chi_sim+eng", "--oem", "3", "--psm", "6")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Printf("Tesseract exit error: %s, stderr: %s", exitErr.Error(), string(exitErr.Stderr))
		}
		return "", fmt.Errorf("tesseract command failed: %w", err)
	}
	return string(output), nil
}

func (s *OCRService) extractIDNumber(text string) string {
	candidates := s.extractAllNumbers(text)
	for _, c := range candidates {
		if c.IsIDNumber {
			return c.Value
		}
	}
	return ""
}

func (s *OCRService) extractAllNumbers(text string) []OCRNumberCandidate {
	var candidates []OCRNumberCandidate
	seen := make(map[string]bool)

	cleaned := regexp.MustCompile(`[\s\-_—–·]`).ReplaceAllString(text, "")
	cleaned = regexp.MustCompile(`[Oo]`).ReplaceAllString(cleaned, "0")
	cleaned = regexp.MustCompile(`[Ii]`).ReplaceAllString(cleaned, "1")

	patterns := []struct {
		regex  *regexp.Regexp
		label  string
		isID   bool
	}{
		{regexp.MustCompile(`\d{17}[\dXx]`), "身份证号", true},
		{regexp.MustCompile(`\d{6}\d{8}\d{3}[\dXx]`), "身份证号", true},
		{regexp.MustCompile(`(\d{4})[\-/.年](\d{1,2})[\-/.月](\d{1,2})`), "日期", false},
		{regexp.MustCompile(`1[3-9]\d{9}`), "手机号", false},
		{regexp.MustCompile(`\d{11,}`), "数字串", false},
		{regexp.MustCompile(`\d{8,}`), "数字串", false},
		{regexp.MustCompile(`\d{6,}`), "数字串", false},
	}

	for _, p := range patterns {
		matches := p.regex.FindAllString(cleaned, -1)
		for _, m := range matches {
			normalized := util.NormalizeIDNumber(m)
			if normalized == "" {
				normalized = m
			}
			if len(normalized) < 4 {
				continue
			}
			if seen[normalized] {
				continue
			}
			seen[normalized] = true

			label := p.label
			if !p.isID {
				switch {
				case len(normalized) == 18:
					label = "18位数字"
				case len(normalized) >= 8 && len(normalized) <= 10:
					label = "日期"
				case len(normalized) == 11:
					label = "11位数字"
				default:
					label = fmt.Sprintf("%d位数字", len(normalized))
				}
			}

			candidates = append(candidates, OCRNumberCandidate{
				Value:      normalized,
				Label:      label,
				IsIDNumber: p.isID,
				Length:     len(normalized),
			})
		}
	}

	return candidates
}
