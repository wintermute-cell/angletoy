package util

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawTextBoxed(font rl.Font, text string, rec rl.Rectangle, fontSize float32, spacing float32, lineSpacingOffset int32, wordWrap bool, tint rl.Color) {
	DrawTextBoxedSelectable(font, text, rec, fontSize, spacing, lineSpacingOffset, wordWrap, tint, 0, 0, rl.White, rl.White)
}

func DrawTextBoxedSelectable(font rl.Font, text string, rec rl.Rectangle, fontSize float32, spacing float32, lineSpacingOffset int32, wordWrap bool, tint rl.Color, selectStart int, selectLength int, selectTint rl.Color, selectBackTint rl.Color) {
	length := len(text)
	textOffsetY := float32(0)
	textOffsetX := float32(0)
	scaleFactor := fontSize / float32(font.BaseSize)

	const (
		MEASURE_STATE = 0
		DRAW_STATE    = 1
	)
	state := DRAW_STATE
	if wordWrap {
		state = MEASURE_STATE
	}

	startLine := -1
	endLine := -1
	lastk := -1
	lastWhitespace := -1
	lastWhitespaceWidth := float32(0)

	i := 0
	for k := 0; i < length; k++ {
		codepoint := rune(text[i])
		codepointByteCount := 1

		var glyphWidth float32
		if codepoint != '\n' {
			glyphWidth = rl.MeasureTextEx(font, string(codepoint), fontSize, spacing).X
		}

		if state == MEASURE_STATE {
			if codepoint == ' ' || codepoint == '\t' || codepoint == '\n' {
				lastWhitespace = i
				lastWhitespaceWidth = textOffsetX + glyphWidth
			}

			if textOffsetX+glyphWidth > rec.Width {
				if lastWhitespace >= 0 {
					endLine = lastWhitespace
					textOffsetX = lastWhitespaceWidth
				} else {
					endLine = i
				}
				state = DRAW_STATE
			} else if i+codepointByteCount >= length {
				endLine = length
				state = DRAW_STATE
			} else if codepoint == '\n' {
				state = DRAW_STATE
			}

			if state == DRAW_STATE {
				textOffsetX = 0
				i = startLine
				glyphWidth = 0

				tmp := lastk
				lastk = k - 1
				k = tmp
			}
		} else {
			if codepoint == '\n' {
				if !wordWrap {
					textOffsetY += (float32(font.BaseSize) + float32(font.BaseSize)/2) * scaleFactor
					textOffsetX = 0
				}
			} else {
				if !wordWrap && (textOffsetX+glyphWidth > rec.Width) {
					textOffsetY += (float32(font.BaseSize) + float32(font.BaseSize)/2) * scaleFactor
					textOffsetX = 0
				}

				if textOffsetY+float32(font.BaseSize)*scaleFactor > rec.Height {
					break
				}

				isGlyphSelected := selectStart >= 0 && k >= selectStart && k < selectStart+selectLength
				if isGlyphSelected {
					rl.DrawRectangleRec(rl.NewRectangle(rec.X+textOffsetX-1, rec.Y+textOffsetY, glyphWidth, float32(font.BaseSize)*scaleFactor), selectBackTint)
				}

				if codepoint != ' ' && codepoint != '\t' {
					currentTint := tint
					if isGlyphSelected {
						currentTint = selectTint
					}
					rl.DrawTextEx(font, string(codepoint), rl.NewVector2(rec.X+textOffsetX, rec.Y+textOffsetY), fontSize, 100, currentTint)
				}
			}

			if wordWrap && i == endLine {
				textOffsetY += (float32(font.BaseSize) + float32(lineSpacingOffset)) * scaleFactor
				textOffsetX = 0
				startLine = endLine
				endLine = -1
				glyphWidth = 0
				selectStart += lastk - k
				k = lastk

				state = MEASURE_STATE
			}
		}

		if textOffsetX != 0 || codepoint != ' ' {
			// FIXED: spacing has to be applied here, since the rl.DrawTextEx call only draws a single character.
			textOffsetX += glyphWidth + spacing
		}

		i += codepointByteCount
	}
}
