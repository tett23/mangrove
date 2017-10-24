package name_resolver

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/unicode/norm"

	"github.com/pkg/errors"
	"github.com/tett23/mangrove/models"
)

type Program struct {
	Name          string
	Title         string
	EpisodeName   string
	EpisodeNumber int
	Digest        string
	queryData     queryData
}

type queryData struct {
	Network string
	StartAt time.Time
	EndAt   time.Time
}

const md5Length = 128 / 8
const extension = ".mp4"
const offset = 1
const maxLength = 256

var nameLenght = maxLength - md5Length - offset - len(extension)

var loc, _ = time.LoadLocation("Asia/Tokyo")

// CreateInstance 名前を出す
func CreateInstance(name string) Program {
	n := strings.TrimSuffix(name, filepath.Ext(name))

	return Program{Name: n}
}

// GetName 名前を出す
func (p *Program) GetName() (string, error) {
	return p.resolve()
}

func (p *Program) resolve() (string, error) {
	if p.Digest == "" {
		// return "", errors.Errorf("Program.resolve not resolved name=%s", p.Name)
	}

	var mp models.Program
	if mp.Search(p.queryData.Network, p.queryData.StartAt, p.queryData.EndAt) {
		p.Title = mp.Title
		p.EpisodeName = mp.EpisodeName
		p.EpisodeNumber = mp.Count
	}

	n := ""
	if p.Title != "" {
		n += p.Title

		if p.EpisodeName != "" {
			n += "「" + p.EpisodeName + "」"
		}

		if p.EpisodeNumber != 0 {
			n += fmt.Sprintf("#%d", p.EpisodeNumber)
		}
	} else {
		n += p.Name
	}

	return fmt.Sprintf("%s_%s.mp4", n, p.Digest), nil
}

// AddProgramData 番組情報を追加
func (p *Program) AddProgramData(text string) (queryData, error) {
	ret := queryData{}

	if len(text) == 0 {
		return ret, errors.Errorf("name_resolver.Program.AddProgramData text not found")
	}

	lines := strings.Split(text, "\n")

	startAt, endAt, err := getStartAndEndDate(lines[0])
	if err != nil {
		return ret, errors.Errorf("name_resolver.Program.AddProgramData datetime parse fail %s", lines[0])
	}

	normalizedChannelName := normalizeChannelName(lines[1])

	ret.Network = normalizedChannelName
	ret.StartAt = startAt
	ret.EndAt = endAt

	p.queryData = ret

	return ret, nil
}

const datetimeFormt = "2006/01/0215:04"

func getStartAndEndDate(str string) (time.Time, time.Time, error) {
	datetime := trim(str)
	date := strings.Split(datetime, "(")[0]
	startEnd := strings.Split(datetime, " ")[1]
	r := strings.Split(startEnd, "～")
	startTime, endTime := r[0], r[1]
	startHourString := strings.Split(startTime, ":")[0]
	endHourString := strings.Split(endTime, ":")[0]
	startHour, err := strconv.Atoi(startHourString)
	endHour, err := strconv.Atoi(endHourString)

	startAt, err := time.Parse(datetimeFormt, date+startTime)
	if err != nil {
		return startAt, startAt, errors.Errorf("getStartAndEndDate datetime parse error %s", datetime)
	}
	endAt, err := time.Parse(datetimeFormt, date+endTime)
	if err != nil {
		return startAt, endAt, errors.Errorf("getStartAndEndDate datetime parse error %s", datetime)
	}

	if startHour < endHour {
		endAt = endAt.AddDate(0, 0, 1)
	}

	return startAt, endAt, nil
}

var channelNameTable = map[string]string{
	"TOKYO MX1": "TOKYO MX",
	"テレビ東京1":    "テレビ東京",
	"TBS1":      "TBS",
	"日テレ1":      "日テレ",
	"NHK総合1・東京": "NHK総合",
	"NHKEテレ1東京": "NHK Eテレ",
}

func normalizeChannelName(str string) string {
	normalized := string(norm.NFKC.Bytes([]byte(trim(str))))

	match, ok := channelNameTable[normalized]
	if ok {
		return match
	}

	return normalized
}

func (p *Program) UpdateMD5(bytes []byte) (string, error) {
	return "", nil
}

func trim(s string) string {
	s = strings.Trim(s, "\r")
	s = strings.Trim(s, "\n")
	s = strings.TrimSpace(s)

	return s
}

func truncate(str string) string {
	if len(str) >= nameLenght {
		return str[0 : nameLenght-1]

	}

	return str
}

// 実ファイル、/年/月くらいで分けて、ファイル名にMD5をいれたほうがいいのかな。サブタイとかは抜きにして、何のアニメ、何話、MD5でいいかな"
