package invindex

func isHangul(r rune) bool {
	return '가' <= r && r <= '힣'
}

func isConsonant(r rune) bool {
	return 'ㄱ' <= r && r < 'ㅎ'
}

func hasBatchim(hangul rune) bool {
	_, _, jongseongIndex := disassembleIntoIndexes(hangul)
	return jongseongIndex > 0
}

func disassembleIntoIndexes(hangul rune) (int, int, int) {
	hangul -= '가'
	return int(hangul / 28 / 21), int(hangul / 28 % 21), int(hangul % 28)
}

func assembleFromIndexes(choseongIndex int, jungseongIndex int, jongseongIndex int) rune {
	return rune('가' + (choseongIndex*21+jungseongIndex)*28 + jongseongIndex)
}

func getFirstVowelPart(jungseong rune) rune {
	switch jungseong {
	case 'ㅘ':
		return 'ㅗ'
	case 'ㅙ':
		return 'ㅗ'
	case 'ㅚ':
		return 'ㅗ'
	case 'ㅝ':
		return 'ㅜ'
	case 'ㅞ':
		return 'ㅜ'
	case 'ㅟ':
		return 'ㅜ'
	case 'ㅢ':
		return 'ㅡ'
	default:
		return -1
	}
}

func getFirstConsonantPart(jongseong rune) rune {
	switch jongseong {
	case 'ㄳ':
		return 'ㄱ'
	case 'ㄵ':
		return 'ㄴ'
	case 'ㄶ':
		return 'ㄴ'
	case 'ㄺ':
		return 'ㄹ'
	case 'ㄻ':
		return 'ㄹ'
	case 'ㄼ':
		return 'ㄹ'
	case 'ㄽ':
		return 'ㄹ'
	case 'ㄾ':
		return 'ㄹ'
	case 'ㄿ':
		return 'ㄹ'
	case 'ㅀ':
		return 'ㄹ'
	case 'ㅄ':
		return 'ㅂ'
	default:
		return -1
	}
}

func splitConsonant(consonant rune) (rune, rune) {
	switch consonant {
	case 'ㄳ':
		return 'ㄱ', 'ㅅ'
	case 'ㄵ':
		return 'ㄴ', 'ㅈ'
	case 'ㄶ':
		return 'ㄴ', 'ㅎ'
	case 'ㄺ':
		return 'ㄹ', 'ㄱ'
	case 'ㄻ':
		return 'ㄹ', 'ㅁ'
	case 'ㄼ':
		return 'ㄹ', 'ㅂ'
	case 'ㄽ':
		return 'ㄹ', 'ㅅ'
	case 'ㄾ':
		return 'ㄹ', 'ㅌ'
	case 'ㄿ':
		return 'ㄹ', 'ㅍ'
	case 'ㅀ':
		return 'ㄹ', 'ㅎ'
	case 'ㅄ':
		return 'ㅂ', 'ㅅ'
	default:
		return consonant, -1
	}
}

func getChoseongIndex(choseong rune) int {
	switch choseong {
	case 'ㄱ':
		return 0
	case 'ㄲ':
		return 1
	case 'ㄴ':
		return 2
	case 'ㄷ':
		return 3
	case 'ㄸ':
		return 4
	case 'ㄹ':
		return 5
	case 'ㅁ':
		return 6
	case 'ㅂ':
		return 7
	case 'ㅃ':
		return 8
	case 'ㅅ':
		return 9
	case 'ㅆ':
		return 10
	case 'ㅇ':
		return 11
	case 'ㅈ':
		return 12
	case 'ㅉ':
		return 13
	case 'ㅊ':
		return 14
	case 'ㅋ':
		return 15
	case 'ㅌ':
		return 16
	case 'ㅍ':
		return 17
	case 'ㅎ':
		return 18
	default:
		return -1
	}
}

func getJungseongIndex(jungseong rune) int {
	switch jungseong {
	case 'ㅏ':
		return 0
	case 'ㅐ':
		return 1
	case 'ㅑ':
		return 2
	case 'ㅒ':
		return 3
	case 'ㅓ':
		return 4
	case 'ㅔ':
		return 5
	case 'ㅕ':
		return 6
	case 'ㅖ':
		return 7
	case 'ㅗ':
		return 8
	case 'ㅘ':
		return 9
	case 'ㅙ':
		return 10
	case 'ㅚ':
		return 11
	case 'ㅛ':
		return 12
	case 'ㅜ':
		return 13
	case 'ㅝ':
		return 14
	case 'ㅞ':
		return 15
	case 'ㅟ':
		return 16
	case 'ㅠ':
		return 17
	case 'ㅡ':
		return 18
	case 'ㅢ':
		return 19
	case 'ㅣ':
		return 20
	default:
		return -1
	}
}

func getJongseongIndex(jongseong rune) int {
	switch jongseong {
	case -1:
		return 0
	case 'ㄱ':
		return 1
	case 'ㄲ':
		return 2
	case 'ㄳ':
		return 3
	case 'ㄴ':
		return 4
	case 'ㄵ':
		return 5
	case 'ㄶ':
		return 6
	case 'ㄷ':
		return 7
	case 'ㄹ':
		return 8
	case 'ㄺ':
		return 9
	case 'ㄻ':
		return 10
	case 'ㄼ':
		return 11
	case 'ㄽ':
		return 12
	case 'ㄾ':
		return 13
	case 'ㄿ':
		return 14
	case 'ㅀ':
		return 15
	case 'ㅁ':
		return 16
	case 'ㅂ':
		return 17
	case 'ㅄ':
		return 18
	case 'ㅅ':
		return 19
	case 'ㅆ':
		return 20
	case 'ㅇ':
		return 21
	case 'ㅈ':
		return 22
	case 'ㅊ':
		return 23
	case 'ㅋ':
		return 24
	case 'ㅌ':
		return 25
	case 'ㅍ':
		return 26
	case 'ㅎ':
		return 27
	default:
		return -1
	}
}

var choseongs = [...]rune{'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}
var jungseongs = [...]rune{'ㅏ', 'ㅐ', 'ㅑ', 'ㅒ', 'ㅓ', 'ㅔ', 'ㅕ', 'ㅖ', 'ㅗ', 'ㅘ', 'ㅙ', 'ㅚ', 'ㅛ', 'ㅜ', 'ㅝ', 'ㅞ', 'ㅟ', 'ㅠ', 'ㅡ', 'ㅢ', 'ㅣ'}
var jongseongs = [...]rune{0, 'ㄱ', 'ㄲ', 'ㄳ', 'ㄴ', 'ㄵ', 'ㄶ', 'ㄷ', 'ㄹ', 'ㄺ', 'ㄻ', 'ㄼ', 'ㄽ', 'ㄾ', 'ㄿ', 'ㅀ', 'ㅁ', 'ㅂ', 'ㅄ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}
