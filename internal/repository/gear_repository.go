package repository

import (
	"encoding/json"
	"github.com/BurntSushi/rure-go"
	"github.com/Lechros/hangul_regexp"
	"itemsim-server/internal"
	"log"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

var gears map[string]json.RawMessage
var names map[string]string
var gearOrigins map[string][2]int
var gearRawOrigins map[string][2]int
var invertedIndex util.HangulInvertedIndex[string]

func SearchGearByName(search string, size int) []json.RawMessage {
	search = strings.ToLower(search)
	pattern, _ := hangul_regexp.GetPattern(search, false, true, true)
	regex := rure.MustCompile("(?i)" + pattern) // (?i): Case insensitive

	exactIds := make([]string, 0, 10)
	regexIds := make([]string, 0, 10)

	useIndex := false
	if utf8.RuneCountInString(search) > 1 {
		foundIds := invertedIndex.FindAll(search)
		if foundIds != nil {
			useIndex = true
			for _, id := range foundIds {
				name := names[id]
				if strings.Contains(name, search) {
					exactIds = append(exactIds, id)
				} else if len(exactIds) < size && regex.IsMatch(name) {
					regexIds = append(regexIds, id)
				}
			}
		}
	}
	if !useIndex {
		for id, name := range names {
			if strings.Contains(name, search) {
				exactIds = append(exactIds, id)
			} else if len(exactIds) < size && regex.IsMatch(name) {
				regexIds = append(regexIds, id)
			}
		}
	}

	if len(exactIds) >= size {
		sortExactMatches(exactIds, search)

		result := make([]json.RawMessage, size)
		for i, id := range exactIds[:size] {
			result[i] = gears[id]
		}
		return result
	} else {
		sortExactMatches(exactIds, search)
		sortRegexMatches(regexIds, regex)

		resultSize := min(len(exactIds)+len(regexIds), size)
		result := make([]json.RawMessage, 0, resultSize)
		for _, id := range exactIds {
			result = append(result, gears[id])
		}
		for _, id := range regexIds {
			if len(result) == resultSize {
				break
			}
			result = append(result, gears[id])
		}
		return result
	}
}

func GetGearById(id string) (json.RawMessage, bool) {
	gear, ok := gears[id]
	return gear, ok
}

func GetGearIconOriginById(id string) ([2]int, bool) {
	origin, ok := gearOrigins[id]
	return origin, ok
}

func GetGearRawIconOriginById(id string) ([2]int, bool) {
	origin, ok := gearRawOrigins[id]
	return origin, ok
}

func sortExactMatches(exactIds []string, search string) {
	type IdInfo struct {
		intId   int
		nameLen int
		index   int
	}
	infos := make([]IdInfo, len(exactIds))
	for i, id := range exactIds {
		intId, _ := strconv.Atoi(id)
		name := names[id]
		index := strings.Index(name, search)
		infos[i] = IdInfo{intId, len(name), index}
	}
	slices.SortFunc(infos, func(aInfo, bInfo IdInfo) int {
		// Contains substring at different index
		if aInfo.index != bInfo.index {
			return aInfo.index - bInfo.index
		}
		// Contains substring at same index
		// Shorter name should be more exact match
		if aInfo.nameLen != bInfo.nameLen {
			return aInfo.nameLen - bInfo.nameLen
		}
		return aInfo.intId - bInfo.intId
	})
	for i, info := range infos {
		exactIds[i] = strconv.Itoa(info.intId)
	}
}

func sortRegexMatches(regexIds []string, regex *rure.Regex) {
	type IdInfo struct {
		intId int
		index [2]int
	}
	infos := make([]IdInfo, len(regexIds))
	for i, id := range regexIds {
		intId, _ := strconv.Atoi(id)
		name := names[id]
		start, end, _ := regex.Find(name)
		infos[i] = IdInfo{intId, [2]int{start, end}}
	}
	slices.SortFunc(infos, func(aInfo, bInfo IdInfo) int {
		if aInfo.index[0] != bInfo.index[0] {
			return aInfo.index[0] - bInfo.index[0]
		}
		if aInfo.index[1] != bInfo.index[1] {
			return aInfo.index[1] - bInfo.index[1]
		}
		return aInfo.intId - bInfo.intId
	})
	for i, info := range infos {
		regexIds[i] = strconv.Itoa(info.intId)
	}
}

func init() {
	util.ReadJson("resources/gear-data.json", &gears)
	util.ReadJson("resources/gear-origin.json", &gearOrigins)
	util.ReadJson("resources/gear-raw-origin.json", &gearRawOrigins)

	names = make(map[string]string)
	invertedIndex = util.HangulInvertedIndex[string]{}
	invertedIndex.Init()
	for id, gearData := range gears {
		var gear map[string]interface{}
		if err := json.Unmarshal(gearData, &gear); err != nil {
			log.Fatal(err)
		}
		name := strings.ToLower(gear["name"].(string))
		names[id] = name
		invertedIndex.Add(id, name)
	}
}
