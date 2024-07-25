package main

import (
	"fmt"
	"strings"

	"github.com/your-org/marker/ocr/tesseract"
	"github.com/your-org/marker/settings"
	"github.com/your-org/surya/languages"
	"github.com/your-org/surya/model/recognition/tokenizer"
)

func langsToIds(langs []string) []int {
	uniqueLangs := make(map[string]bool)
	for _, lang := range langs {
		uniqueLangs[lang] = true
	}

	uniqueLangsList := make([]string, 0, len(uniqueLangs))
	for lang := range uniqueLangs {
		uniqueLangsList = append(uniqueLangsList, lang)
	}

	_, langTokens := tokenizer.LangTokenize("", uniqueLangsList)
	return langTokens
}

func replacelangsWithCodes(langs []string) []string {
	if settings.OCREngine == "surya" {
		for i, lang := range langs {
			if code, ok := languages.LanguageToCode[strings.Title(lang)]; ok {
				langs[i] = code
			}
		}
	} else {
		for i, lang := range langs {
			if code, ok := tesseract.LanguageToTesseractCode[lang]; ok {
				langs[i] = code
			}
		}
	}
	return langs
}

func validateLangs(langs []string) error {
	if settings.OCREngine == "surya" {
		for _, lang := range langs {
			if _, ok := languages.CodeToLanguage[lang]; !ok {
				return fmt.Errorf("invalid language code %s for Surya OCR", lang)
			}
		}
	} else {
		for _, lang := range langs {
			if _, ok := tesseract.TesseractCodeToLanguage[lang]; !ok {
				return fmt.Errorf("invalid language code %s for Tesseract", lang)
			}
		}
	}
	return nil
}
