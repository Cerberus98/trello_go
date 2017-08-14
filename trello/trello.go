package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var (
	ErrCannotReadBody = errors.New("Cannot read response body")
)

type TrelloBoardPrefs struct {
	Voting                string
	Comments              string
	BackgroundImageScaled string
	BackgroundTile        bool
	BackgroundBrightness  string
	CanBeOrg              bool
	PermissionLevel       string
	CardAging             string
	Background            string
	BackgroundImage       string
	CanBePublic           bool
	Invitations           string
	SelfJoin              bool
	BackgroundColor       string
	CardCovers            bool
	CalendarFeedEnabled   bool
	CanBePrivate          bool
	CanInvite             bool
}

type TrelloBoard struct {
	Client         TrelloApi
	Id             string
	Name           string
	IdOrganization string
	Closed         bool
	URL            string
	Desc           string
	DescData       string
	Pinned         bool
	ShortURL       string
	LabelNames     map[string]string
	Prefs          TrelloBoardPrefs `json:"prefs"`
}

type TrelloCardBadge struct {
	Description        bool
	Due                string
	DueComplete        bool
	ViewingMemberVoted bool
	Attachments        float64
	Fogbugz            struct {
		CheckItems        float64
		CheckItemsChecked float64
		Comments          float64
		Votes             float64
		Subscribed        float64
	}
}

type TrelloCard struct {
	IdLabels              []string
	Id                    string
	CheckItemStates       []string
	DateLastActivity      string
	Desc                  string
	DescData              string
	IdAttachmentCover     string
	Limits                int
	Pos                   float64
	Subscribed            bool
	IdBoard               string
	IdList                string
	ShortLink             string
	DueComplete           bool
	Due                   string
	IdChecklists          []string
	Labels                []string
	Badges                []TrelloCardBadge
	IdMembers             []string
	ShortUrl              string
	Closed                bool
	IdMembersVoted        []string
	IdShort               float64
	ManualCoverAttachment bool
	Name                  string
	Url                   string
}

type TrelloTimeZone struct {
	TimezoneNext    string
	DateNext        string
	OffsetNext      int
	TimezoneCurrent string
	OffsetCurrent   int
}

type TrelloMember struct {
	Id                       string
	AvatarHash               string
	Bio                      string
	BioData                  string
	Confirmed                bool
	FullName                 bool
	IdPremOrgsAdmin          []string
	Initials                 string
	MemberType               string
	Products                 string
	Status                   string
	Url                      string
	Username                 string
	AvatarSource             string
	Email                    string
	GravatarHash             string
	IdBoards                 []string
	IdEnterprise             string
	IdOrganizations          []string
	IdEnterprisesAdmin       []string
	LoginTypes               string
	OneTimeMessagesDismissed string
	Prefs                    struct {
		SendSummaries                 bool
		MinutesBetweenSummaries       int64
		MinutesBeforeDeadlineToNotify int64
		ColorBlind                    bool
		Locale                        string
		TimezoneInfo                  TrelloTimeZone
		TwoFactor                     struct {
			Enabled         bool
			NeedsNewBackups bool
		}
	}
	Trophies           []string
	UploadedAvatarHash string
	PremiumFeatures    string
	IdBoardsPinned     string
}

type HttpClient struct {
}

type TrelloApi interface {
	GetVersion() int
	getClient() *HttpClient
	UrlFor(apiKey string, token string, route ...string) string
	GetBoard(id string) *TrelloBoard
	GetCards(boardId string) []TrelloCard
	GetMembers(boardId string) []TrelloMember
}

type TrelloApiV1 struct {
	Client  HttpClient
	Key     string
	Token   string
	BaseUrl string
}

func (trello *TrelloApiV1) GetVersion() int {
	return 1
}

func (tapi *TrelloApiV1) getClient() *HttpClient {
	return &tapi.Client
}

func (tapi *TrelloApiV1) get(url string, trelloObj interface{}) error {
	if body, err := tapi.getClient().Get(url); err != nil {
		return err
	} else {
		json.Unmarshal(body, &trelloObj)
	}
	return nil
}

func (tapi *TrelloApiV1) GetBoard(id string) *TrelloBoard {
	var trelloBoard TrelloBoard
	url := tapi.UrlFor(tapi.Key, tapi.Token, "boards", id)
	if err := tapi.get(url, &trelloBoard); err != nil {
		log.Fatalf("Fetching boards failed with %s", err)
	}
	trelloBoard.Client = tapi
	return &trelloBoard
}

func (board *TrelloBoard) GetCards() []TrelloCard {
	return board.Client.GetCards(board.Id)
}

func (board *TrelloBoard) GetMembers() []TrelloMember {
	return board.Client.GetMembers(board.Id)
}

func (tapi *TrelloApiV1) GetCards(boardId string) []TrelloCard {
	var trelloCards []TrelloCard
	url := tapi.UrlFor(tapi.Key, tapi.Token, "boards", boardId, "cards")
	if err := tapi.get(url, &trelloCards); err != nil {
		log.Fatalf("Fetching cards failed with %s", err)
	}
	return trelloCards
}

func (tapi *TrelloApiV1) GetMembers(boardId string) []TrelloMember {
	var trelloMembers []TrelloMember
	url := tapi.UrlFor(tapi.Key, tapi.Token, "boards", boardId, "members")
	if err := tapi.get(url, &trelloMembers); err != nil {
		log.Fatalf("Fetching members failed with %s", err)
	}
	return trelloMembers
}

func (h *HttpClient) Get(url string) ([]byte, error) {
	if resp, err := http.Get(url); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			printResponseInfo(resp)
		}

		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return []byte{}, ErrCannotReadBody
		} else {
			return body, nil
		}
	} else {
		return []byte{}, errors.New(fmt.Sprintf("Failed to GET %s", url))
	}
}

type RequestParams struct {
	arguments map[string][]string
}

func (r *RequestParams) AddParam(key, value string) {
	if r.arguments == nil {
		r.arguments = make(map[string][]string)
	}

	r.arguments[key] = append(r.arguments[key], value)
}

func (r *RequestParams) ToValues() *url.Values {
	values := &url.Values{}
	for k, v := range r.arguments {
		values.Set(k, strings.Join(v, ","))
	}
	return values
}

func (tapi *TrelloApiV1) UrlFor(apiKey string, token string, route ...string) string {
	params := RequestParams{}
	params.AddParam("key", apiKey)
	params.AddParam("token", token)
	encoded := params.ToValues().Encode()

	return fmt.Sprintf("%s/%d/%s?%s", tapi.BaseUrl, tapi.GetVersion(), path.Join(route...), encoded)
}

func printResponseInfo(resp *http.Response) {
	fmt.Println(resp.Status)
	fmt.Println("Headers")

	for k, values := range resp.Header {
		fmt.Printf("\t%s: ", k)
		var headerbuf bytes.Buffer
		for _, v := range values {
			headerbuf.WriteString(fmt.Sprintf("\"%s\"", v))
		}
		fmt.Printf("\t%s\n", headerbuf.String())
	}
	fmt.Printf("Content-Length: %d\n", resp.ContentLength)
}
