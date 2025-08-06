package utils

import (
	"fmt"
	"strings"
)

// HighLightString ..
func HighLightString(str string, subStr string, highLightMinecraftColor string) (result string) {
	subStr = strings.ToLower(subStr)
	newStr := strings.ReplaceAll(
		strings.ToLower(str),
		subStr,
		fmt.Sprintf("§r%s%s§r", highLightMinecraftColor, subStr),
	)

	strPtr, newStrPtr := 0, 0
	strLen, newStrLen := len(str), len(newStr)

	for {
		for {
			if newStrPtr >= newStrLen-1 {
				break
			}
			if newStr[newStrPtr:newStrPtr+2] == "§" {
				result += newStr[newStrPtr:min(newStrPtr+3, newStrLen)]
				newStrPtr += 3
			} else {
				break
			}
		}

		for {
			if strPtr >= strLen-1 {
				break
			}
			if str[strPtr:strPtr+2] == "§" {
				strPtr += 3
			} else {
				break
			}
		}

		if strPtr >= len(str) {
			break
		}

		result += str[strPtr : strPtr+1]
		strPtr++
		newStrPtr++
	}

	return result
}
