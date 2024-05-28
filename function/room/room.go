package room

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type Getter interface {
	Get(Requirement) ([]Room, error)
}

type Requirement struct {
	RentLow        int
	RentHigh       int
	FloorSpaceLow  int
	FloorSpaceHigh int
	Room           []string
	Prefecture     PrefectureCode
	Area           AreaCode
}

type Room struct {
	Name      string `json:"name"`
	Location  string `json:"skcs"`
	RoomCount int    `json:"roomCount"`
	URL       string `json:"bukkenUrl"`
}

const PREFECTURE_TOKYO = "13"

type PrefectureCode string

const (
	Tokyo = PrefectureCode("13")
)

type AreaCode string

const (
	East23 = AreaCode("02")
)

func New() Getter {
	return &getterImpl{}
}

const URL = "https://chintai.r6.ur-net.go.jp/chintai/api/bukken/search/list_bukken/"

type getterImpl struct{}

func (g *getterImpl) Get(requirement Requirement) ([]Room, error) {
	values := url.Values{}
	if requirement.RentLow != 0 {
		values.Add("rent_low", strconv.Itoa(requirement.RentLow))
	}
	if requirement.RentHigh != 0 {
		values.Add("rent_high", strconv.Itoa(requirement.RentHigh))
	}
	if requirement.FloorSpaceLow != 0 {
		values.Add("floorspace_low", strconv.Itoa(requirement.FloorSpaceLow))
	}
	if requirement.FloorSpaceHigh != 0 {
		values.Add("floorspace_high", strconv.Itoa(requirement.FloorSpaceHigh))
	}
	for _, r := range requirement.Room {
		values.Add("room", r)
	}
	if requirement.Prefecture != "" {
		values.Add("tdfk", string(requirement.Prefecture))
	}
	if requirement.Area != "" {
		values.Add("area", string(requirement.Area))
	}

	req, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to make new HTTP request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, fmt.Errorf("failed to dump HTTP request: %w", err)
	}
	fmt.Println(string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}
	fmt.Println(string(content))

	rooms := []Room{}
	if err := json.Unmarshal(content, &rooms); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return rooms, nil
}
