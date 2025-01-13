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
)

var gears map[int]json.RawMessage
var names map[int]string
var gearOrigins map[int][2]int
var gearRawOrigins map[int][2]int
var concatNames string
var concatIndex []int
var concatIds []int

func SearchGearByName(search string, size int) (data []json.RawMessage, highlight []string) {
	search = strings.ToLower(search)
	pattern, _ := hangul_regexp.GetPattern(search, false, true, true, false)
	regex := rure.MustCompile("(?i)" + pattern) // (?i): Case insensitive

	// Can't infer capacity from FindAll's length since single line can match multiple times.
	matchedIds := make([]int, 0, 10)

	matches := regex.FindAll(concatNames)
	lastIndex := -1
	for i, match := range matches {
		if i%2 == 0 {
			index, found := slices.BinarySearch(concatIndex, match)
			if !found {
				index--
			}
			if index > lastIndex {
				id := concatIds[index]
				matchedIds = append(matchedIds, id)
				lastIndex = index
			}
		}
	}

	sortRegexMatches(matchedIds, regex)

	capturingPattern, _ := hangul_regexp.GetPattern(search, false, true, true, true)
	capturingRegex := rure.MustCompile("(?i)" + capturingPattern)
	resultSize := min(len(matchedIds), size)
	data = make([]json.RawMessage, resultSize)
	highlight = make([]string, resultSize)
	for i, id := range matchedIds {
		if i == resultSize {
			break
		}
		data[i] = gears[id]
		highlight[i] = convertHighlight(names[id], capturingRegex)
	}
	return data, highlight
}

func GetGearById(id string) (json.RawMessage, bool) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil, false
	}
	gear, ok := gears[intId]
	return gear, ok
}

func GetGearIconOriginById(id string) ([2]int, bool) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return [2]int{}, false
	}
	origin, ok := gearOrigins[intId]
	return origin, ok
}

func GetGearRawIconOriginById(id string) ([2]int, bool) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		return [2]int{}, false
	}
	origin, ok := gearRawOrigins[intId]
	return origin, ok
}

func sortRegexMatches(ids []int, regex *rure.Regex) {
	type IdInfo struct {
		id    int
		name  string
		index [2]int
	}
	infos := make([]IdInfo, len(ids))
	for i, id := range ids {
		name := names[id]
		start, end, _ := regex.Find(name)
		infos[i] = IdInfo{id, name, [2]int{start, end}}
	}
	slices.SortFunc(infos, func(aInfo, bInfo IdInfo) int {
		if aInfo.index[0] != bInfo.index[0] {
			return aInfo.index[0] - bInfo.index[0]
		}
		if aInfo.index[1] != bInfo.index[1] {
			return aInfo.index[1] - bInfo.index[1]
		}
		if len(aInfo.name) != len(bInfo.name) {
			return len(aInfo.name) - len(bInfo.name)
		}
		return aInfo.id - bInfo.id
	})
	for i, info := range infos {
		ids[i] = info.id
	}
}

// Returns string where match rune is '1' and else is '0'
func convertHighlight(name string, regex *rure.Regex) string {
	captures := regex.NewCaptures()
	regex.Captures(captures, name)
	builder := strings.Builder{}
	builder.Grow(len(name)) // Generous amount of buffer is faster than utf8.RuneCountInString, or reallocation

	gi := 1
	start, end, _ := captures.Group(gi)
	for i := range name {
		for gi < captures.Len() && i >= end {
			gi++
			start, end, _ = captures.Group(gi)
		}
		if i >= start && i < end {
			builder.WriteRune('1')
		} else {
			builder.WriteRune('0')
		}
	}
	return builder.String()
}

func init() {
	util.ReadJson("resources/gear-data.json", &gears)
	util.ReadJson("resources/gear-origin.json", &gearOrigins)
	util.ReadJson("resources/gear-raw-origin.json", &gearRawOrigins)

	names = make(map[int]string)
	concatIndex = make([]int, 0, len(gears))
	concatIds = make([]int, 0, len(gears))
	builder := strings.Builder{}
	for id, gearData := range gears {
		var gear map[string]interface{}
		if err := json.Unmarshal(gearData, &gear); err != nil {
			log.Fatal(err)
		}
		name := strings.ToLower(gear["name"].(string))
		names[id] = name

		concatIndex = append(concatIndex, builder.Len())
		concatIds = append(concatIds, id)
		builder.WriteString(name)
		builder.WriteRune('\n')
	}
	concatNames = builder.String()
}
