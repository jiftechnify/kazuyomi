package kazuyomi

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ErrNotNumericString is returned when input is not a numeric string.
	ErrNotNumericString = errors.New("input string is not numeric")
)

var regexpNumeric = regexp.MustCompile(`^(-|\+)?[[:digit:],_.]+$`)

func cleanNumericStr(numStr string) string {
	return strings.NewReplacer(",", "", "_", "").Replace(numStr)
}

var basicDigitReadings = map[rune]string{
	'0': "ゼロ",
	'1': "イチ",
	'2': "ニ",
	'3': "サン",
	'4': "ヨン",
	'5': "ゴ",
	'6': "ロク",
	'7': "ナナ",
	'8': "ハチ",
	'9': "キュウ",
	'.': "テン",
}

type smallNumDigitIdx int

const (
	thousands smallNumDigitIdx = 0
	hundreds  smallNumDigitIdx = 1
	tens      smallNumDigitIdx = 2
)

type digitWithIdx struct {
	digit rune
	idx   smallNumDigitIdx
}

var specialDigitReadings = map[digitWithIdx]string{
	{digit: '1', idx: thousands}: "セン",
	{digit: '1', idx: hundreds}:  "ヒャク",
	{digit: '1', idx: tens}:      "ジュウ",
	{digit: '3', idx: thousands}: "サンゼン",
	{digit: '3', idx: hundreds}:  "サンビャク",
	{digit: '6', idx: hundreds}:  "ロッピャク",
	{digit: '8', idx: thousands}: "ハッセン",
	{digit: '8', idx: hundreds}:  "ハッピャク",
}

func appendLiteralReading(sb *strings.Builder, strNum string) {
	for _, r := range strNum {
		sb.WriteString(basicDigitReadings[r])
	}
}

func smallIntReading(runes []rune) string {
	var sb strings.Builder
	bias := 4 - len(runes)
	for i, r := range runes {
		if r == '0' {
			continue
		}
		if read, ok := specialDigitReadings[digitWithIdx{digit: r, idx: smallNumDigitIdx(i + bias)}]; ok {
			sb.WriteString(read)
			continue
		}
		sb.WriteString(basicDigitReadings[r])
		switch i + bias {
		case 0:
			sb.WriteString("セン")
		case 1:
			sb.WriteString("ヒャク")
		case 2:
			sb.WriteString("ジュウ")
		}
	}
	return sb.String()
}

// nasal sound change = 促音便
func applyNasalSoundChange(r string) string {
	if cut, found := strings.CutSuffix(r, "イチ"); found {
		return cut + "イッ"
	}
	if cut, found := strings.CutSuffix(r, "ハチ"); found {
		return cut + "ハッ"
	}
	if cut, found := strings.CutSuffix(r, "ジュウ"); found {
		return cut + "ジッ"
	}
	return r
}

func intPartReading(strInt string) string {
	if strInt == "" {
		return ""
	}
	if strInt == "0" {
		return "ゼロ"
	}

	var sb strings.Builder
	runes := []rune(strInt)
	for i := 5; i >= 1; i-- {
		if len(runes) <= 4*(i-1) {
			continue
		}
		smallRead := smallIntReading(runes[max(len(runes)-4*i, 0) : len(runes)-4*(i-1)])
		if smallRead == "" {
			continue
		}
		switch i {
		case 5:
			sb.WriteString(applyNasalSoundChange(smallRead))
			sb.WriteString("ケイ")
		case 4:
			sb.WriteString(applyNasalSoundChange(smallRead))
			sb.WriteString("チョウ")
		case 3:
			sb.WriteString(smallRead)
			sb.WriteString("オク")
		case 2:
			sb.WriteString(smallRead)
			sb.WriteString("マン")
		case 1:
			sb.WriteString(smallRead)
		}
	}
	return sb.String()
}

func consumePrefixedSign(s string) (string, string) {
	if cut, found := strings.CutPrefix(s, "-"); found {
		return cut, "マイナス"
	}
	if cut, found := strings.CutPrefix(s, "+"); found {
		return cut, "プラス"
	}
	return s, ""
}

func numericStrReading(numStr string) string {
	var sb strings.Builder

	// read prefixed sign
	s, signRead := consumePrefixedSign(numStr)
	sb.WriteString(signRead)

	parts := strings.Split(s, ".")

	intPart := parts[0]
	if len(intPart) >= 2 && intPart[0] == '0' || len(intPart) > 20 {
		appendLiteralReading(&sb, s)
		return sb.String()
	}

	// read integer part
	ipr := intPartReading(intPart)

	// input is integer or decimal part is empty
	if len(parts) == 1 || parts[1] == "" {
		sb.WriteString(ipr)
		return sb.String()
	}

	// read decimal part
	if len(parts) == 2 && parts[1] != "" {
		if ipr == "ゼロ" {
			ipr = "レイ"
		}
		sb.WriteString(applyNasalSoundChange(ipr))
		sb.WriteString("テン")
		appendLiteralReading(&sb, parts[1])
	}
	return sb.String()
}

// ReadString returns the Japanese reading (読み仮名) of the given numeric string.
// The result is given as a string of katakanas (カタカナ).
//
// If the given string is not a numeric string, ErrNotNumericString will be returned.
// "Numeric string" is defined by the following regexp:
//
//	^(-|\+)?[0-9,_.]+$
//
// As shown above, the input string can contain "," and "_" as separators, and they will be ignored.
//
// If the input numeric string satisfies the conditions below, the result will be the enumeration of the "literal" reading of each digit:
//   - has multiple dots (e.g. "1.2.3")
//   - its integer part has more than 20 digits (N >= 10^20(1垓))
//   - its integer part starts with "0", except it is exactly single digit "0"
func ReadString(s string) (string, error) {
	if !regexpNumeric.MatchString(s) {
		return "", ErrNotNumericString
	}

	dots := strings.Count(s, ".")
	if dots >= 2 {
		var sb strings.Builder
		appendLiteralReading(&sb, s)
		return sb.String(), nil
	}

	return numericStrReading(cleanNumericStr(s)), nil
}

// ReadInt returns the Japanese reading (読み仮名) of the given signed integer.
// The result is given as a string of katakanas (カタカナ).
func ReadInt(i int) string {
	res, _ := ReadString(strconv.FormatInt(int64(i), 10))
	return res
}

// ReadUint returns the Japanese reading (読み仮名) of the given unsigned integer.
// The result is given as a string of katakanas (カタカナ).
func ReadUint(i uint) string {
	res, _ := ReadString(strconv.FormatUint(uint64(i), 10))
	return res
}

// ReadFloat64 returns the Japanese reading (読み仮名) of the given 64-bits floating point number.
// The result is given as a string of katakanas (カタカナ).
func ReadFloat64(f float64) string {
	res, _ := ReadString(strconv.FormatFloat(f, 'f', -1, 64))
	return res
}
