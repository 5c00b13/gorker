package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/sahilm/fuzzy"
)

const CHUNK_MIN_CHARS = 25

func chunkText(text string, chunkLen int) []string {
	var chunks []string
	for i := 0; i < len(text); i += chunkLen {
		end := i + chunkLen
		if end > len(text) {
			end = len(text)
		}
		chunk := strings.TrimSpace(text[i:end])
		if len(chunk) > CHUNK_MIN_CHARS {
			chunks = append(chunks, chunk)
		}
	}
	return chunks
}

func overlapScore(hypothesisChunks, referenceChunks []string) []float64 {
	lengthModifier := float64(len(hypothesisChunks)) / float64(len(referenceChunks))
	searchDistance := int(math.Max(float64(len(referenceChunks)/5), 10))
	chunkScores := make([]float64, 0, len(hypothesisChunks))

	for i, hypChunk := range hypothesisChunks {
		maxScore := 0.0
		iOffset := int(float64(i) * lengthModifier)
		chunkRangeStart := int(math.Max(0, float64(iOffset-searchDistance)))
		chunkRangeEnd := int(math.Min(float64(len(referenceChunks)), float64(iOffset+searchDistance)))

		for j := chunkRangeStart; j < chunkRangeEnd; j++ {
			refChunk := referenceChunks[j]
			score := fuzzy.RatioForStrings([]rune(hypChunk), []rune(refChunk), fuzzy.DefaultOptions)
			if score > 30 {
				score = float64(score) / 100
				if score > maxScore {
					maxScore = score
				}
			}
		}
		chunkScores = append(chunkScores, maxScore)
	}
	return chunkScores
}

func mean(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

func scoreText(hypothesis, reference string) float64 {
	hypothesisChunks := chunkText(hypothesis, 500)
	referenceChunks := chunkText(reference, 500)
	chunkScores := overlapScore(hypothesisChunks, referenceChunks)
	return mean(chunkScores)
}

func main() {
	hypothesis := "Your hypothesis text here"
	reference := "Your reference text here"
	score := scoreText(hypothesis, reference)
	fmt.Printf("Alignment score: %f\n", score)
}
