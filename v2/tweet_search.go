package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TweetSearchSortOrder specifies the order the tweets are returned
type TweetSearchSortOrder string

const (
	// TweetSearchSortOrderRecency will return the tweets in order of recency
	TweetSearchSortOrderRecency TweetSearchSortOrder = "recency"
	// TweetSearchSortOrderRelevancy will return the tweets in order of relevancy
	TweetSearchSortOrderRelevancy TweetSearchSortOrder = "relevancy"
)

// TweetRecentSearchOpts are the optional parameters for the recent search API
type TweetRecentSearchOpts struct {
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
	StartTime   time.Time
	EndTime     time.Time
	SortOrder   TweetSearchSortOrder
	MaxResults  int
	NextToken   string
	SinceID     string
	UntilID     string
}

func (t TweetRecentSearchOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(t.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(t.Expansions), ","))
	}
	if len(t.MediaFields) > 0 {
		q.Add("media.fields", strings.Join(mediaFieldStringArray(t.MediaFields), ","))
	}
	if len(t.PlaceFields) > 0 {
		q.Add("place.fields", strings.Join(placeFieldStringArray(t.PlaceFields), ","))
	}
	if len(t.PollFields) > 0 {
		q.Add("poll.fields", strings.Join(pollFieldStringArray(t.PollFields), ","))
	}
	if len(t.TweetFields) > 0 {
		q.Add("tweet.fields", strings.Join(tweetFieldStringArray(t.TweetFields), ","))
	}
	if len(t.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(t.UserFields), ","))
	}
	if !t.StartTime.IsZero() {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if !t.EndTime.IsZero() {
		q.Add("end_time", t.EndTime.Format(time.RFC3339))
	}
	if t.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(t.MaxResults))
	}
	if len(t.NextToken) > 0 {
		q.Add("next_token", t.NextToken)
	}
	if len(t.SinceID) > 0 {
		q.Add("since_id", t.SinceID)
	}
	if len(t.UntilID) > 0 {
		q.Add("until_id", t.UntilID)
	}
	if len(t.SortOrder) > 0 {
		q.Add("sort_order", string(t.SortOrder))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

type TweetRecentSearchAsyncResponse struct {
	ID string `json:"job_id"`
}

func (*TweetRecentSearchAsyncResponse) Build(statusCode int, headers http.Header, body io.Reader) (*TweetRecentSearchResponse, error) {
	buffer, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("tweet recent search response read: %w", err)
	}

	rl := rateFromHeader(headers)
	darl := dailyAppRateFromHeader(headers)
	durl := dailyUserRateFromHeader(headers)

	if statusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.Unmarshal(buffer, e); err != nil {
			return nil, &HTTPError{
				StatusCode:         statusCode,
				RateLimit:          rl,
				DailyAppRateLimit:  darl,
				DailyUserRateLimit: durl,
			}
		}
		e.StatusCode = statusCode
		e.RateLimit = rl
		e.DailyAppRateLimit = darl
		e.DailyUserRateLimit = durl
		return nil, e
	}

	recentSearch := &TweetRecentSearchResponse{
		Raw:       &TweetRaw{},
		Meta:      &TweetRecentSearchMeta{},
		RateLimit: rl,
	}

	if err := json.Unmarshal(buffer, recentSearch.Raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet recent search",
			Err:       err,
			RateLimit: rl,
		}
	}

	if err := json.Unmarshal(buffer, recentSearch); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet recent search",
			Err:       err,
			RateLimit: rl,
		}
	}

	return recentSearch, nil
}

// TweetRecentSearchResponse contains all of the information from a tweet recent search
type TweetRecentSearchResponse struct {
	Raw       *TweetRaw
	Meta      *TweetRecentSearchMeta `json:"meta"`
	RateLimit *RateLimit
}

// TweetRecentSearchMeta contains the recent search information
type TweetRecentSearchMeta struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

// TweetSearchOpts are the tweet search options
type TweetSearchOpts struct {
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
	StartTime   time.Time
	EndTime     time.Time
	SortOrder   TweetSearchSortOrder
	MaxResults  int
	NextToken   string
	SinceID     string
	UntilID     string
}

func (t TweetSearchOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(t.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(t.Expansions), ","))
	}
	if len(t.MediaFields) > 0 {
		q.Add("media.fields", strings.Join(mediaFieldStringArray(t.MediaFields), ","))
	}
	if len(t.PlaceFields) > 0 {
		q.Add("place.fields", strings.Join(placeFieldStringArray(t.PlaceFields), ","))
	}
	if len(t.PollFields) > 0 {
		q.Add("poll.fields", strings.Join(pollFieldStringArray(t.PollFields), ","))
	}
	if len(t.TweetFields) > 0 {
		q.Add("tweet.fields", strings.Join(tweetFieldStringArray(t.TweetFields), ","))
	}
	if len(t.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(t.UserFields), ","))
	}
	if !t.StartTime.IsZero() {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if !t.EndTime.IsZero() {
		q.Add("end_time", t.EndTime.Format(time.RFC3339))
	}
	if t.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(t.MaxResults))
	}
	if len(t.NextToken) > 0 {
		q.Add("next_token", t.NextToken)
	}
	if len(t.SinceID) > 0 {
		q.Add("since_id", t.SinceID)
	}
	if len(t.UntilID) > 0 {
		q.Add("until_id", t.UntilID)
	}
	if len(t.SortOrder) > 0 {
		q.Add("sort_order", string(t.SortOrder))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// TweetSearchResponse is the tweet search response
type TweetSearchResponse struct {
	Raw       *TweetRaw
	Meta      *TweetSearchMeta `json:"meta"`
	RateLimit *RateLimit
}

// TweetSearchMeta is the tweet search meta data
type TweetSearchMeta struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

// TweetSearchStreamRule is the search stream filter rule
type TweetSearchStreamRule struct {
	Value string `json:"value"`
	Tag   string `json:"tag,omitempty"`
}

func (t TweetSearchStreamRule) validate() error {
	if len(t.Value) == 0 {
		return fmt.Errorf("tweet search stream rule value is required: %w", ErrParameter)
	}
	return nil
}

type tweetSearchStreamRules []TweetSearchStreamRule

func (t tweetSearchStreamRules) validate() error {
	for _, rule := range t {
		if err := rule.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TweetSearchStreamRuleID is the filter rule id
type TweetSearchStreamRuleID string

func (t TweetSearchStreamRuleID) validate() error {
	if len(t) == 0 {
		return fmt.Errorf("tweet search rule id is required %w", ErrParameter)
	}
	return nil
}

type tweetSearchStreamRuleIDs []TweetSearchStreamRuleID

func (t tweetSearchStreamRuleIDs) validate() error {
	for _, id := range t {
		if err := id.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (t tweetSearchStreamRuleIDs) toStringArray() []string {
	ids := make([]string, len(t))
	for i, id := range t {
		ids[i] = string(id)
	}
	return ids
}

// TweetSearchStreamRuleEntity is the search filter rule entity
type TweetSearchStreamRuleEntity struct {
	ID TweetSearchStreamRuleID `json:"id"`
	TweetSearchStreamRule
}

// TweetSearchStreamRulesResponse is the response to getting the search rules
type TweetSearchStreamRulesResponse struct {
	Rules     []*TweetSearchStreamRuleEntity `json:"data"`
	Meta      *TweetSearchStreamRuleMeta     `json:"meta"`
	Errors    []*ErrorObj                    `json:"errors,omitempty"`
	RateLimit *RateLimit
}

// TweetSearchStreamAddRuleResponse is the response from adding rules
type TweetSearchStreamAddRuleResponse struct {
	Rules     []*TweetSearchStreamRuleEntity `json:"data"`
	Meta      *TweetSearchStreamRuleMeta     `json:"meta"`
	Errors    []*ErrorObj                    `json:"errors,omitempty"`
	RateLimit *RateLimit
}

// TweetSearchStreamDeleteRuleResponse is the response from deleting rules
type TweetSearchStreamDeleteRuleResponse struct {
	Meta      *TweetSearchStreamRuleMeta `json:"meta"`
	Errors    []*ErrorObj                `json:"errors,omitempty"`
	RateLimit *RateLimit
}

// TweetSearchStreamRuleMeta is the meta data object from the request
type TweetSearchStreamRuleMeta struct {
	Sent    time.Time                    `json:"sent"`
	Summary TweetSearchStreamRuleSummary `json:"summary"`
}

// TweetSearchStreamRuleSummary is the summary of the search filters
type TweetSearchStreamRuleSummary struct {
	Created    int `json:"created"`
	NotCreated int `json:"not_created"`
	Deleted    int `json:"deleted"`
	NotDeleted int `json:"not_deleted"`
}

type TweetSearchApifyAsyncRequest struct {
	StartURLs          []string `json:"startUrls,omitempty"`          // Twitter (X) URLs. Paste the URLs and get the results immediately. Tweet, Profile, Search or List URLs are supported.
	SearchTerms        []string `json:"searchTerms,omitempty"`        // Search terms you want to search from Twitter (X). You can refer to https://github.com/igorbrigadir/twitter-advanced-search.
	TwitterHandles     []string `json:"twitterHandles,omitempty"`     // Twitter handles that you want to search on Twitter (X)
	ConversationIDs    []string `json:"conversationIds,omitempty"`    // Conversation IDs that you want to search on Twitter (X)
	MaxItems           int      `json:"maxItems,omitempty"`           // Maximum number of items that you want as output.
	Sort               string   `json:"sort,omitempty"`               // Sorts search results by the given option. Only works with search terms and search URLs. Value options: "Top", "Latest", "Media"
	TweetLanguage      string   `json:"tweetLanguage,omitempty"`      // Restricts tweets to the given language, given by an ISO 639-1 code.
	OnlyVerifiedUsers  bool     `json:"onlyVerifiedUsers,omitempty"`  // If selected, only returns tweets by users who are verified.
	OnlyTwitterBlue    bool     `json:"onlyTwitterBlue,omitempty"`    // If selected, only returns tweets by users who are Twitter Blue subscribers.
	OnlyImage          bool     `json:"onlyImage,omitempty"`          // If selected, only returns tweets that contain images.
	OnlyVideo          bool     `json:"onlyVideo,omitempty"`          // If selected, only returns tweets that contain videos.
	OnlyQuote          bool     `json:"onlyQuote,omitempty"`          // If selected, only returns tweets that are quotes.
	Author             string   `json:"author,omitempty"`             // Returns tweets sent by the given user. It should be a Twitter (X) Handle.
	InReplyTo          string   `json:"inReplyTo,omitempty"`          // Returns tweets that are replies to the given user. It should be a Twitter (X) Handle.
	Mentioning         string   `json:"mentioning,omitempty"`         // Returns tweets mentioning the given user. It should be a Twitter (X) Handle.
	GeotaggedNear      string   `json:"geotaggedNear,omitempty"`      // Returns tweets sent near the given location.
	WithinRadius       string   `json:"withinRadius,omitempty"`       // Returns tweets sent within the given radius of the given location.
	Geocode            string   `json:"geocode,omitempty"`            // Returns tweets sent by users located within a given radius of the given latitude/longitude.
	PlaceObjectId      string   `json:"placeObjectId,omitempty"`      // Returns tweets tagged with the given place.
	MinimumRetweets    int      `json:"minimumRetweets,omitempty"`    // Returns tweets with at least the given number of retweets.
	MinimumFavorites   int      `json:"minimumFavorites,omitempty"`   // Returns tweets with at least the given number of favorites.
	MinimumReplies     int      `json:"minimumReplies,omitempty"`     // Returns tweets with at least the given number of replies.
	Start              string   `json:"start,omitempty"`              // Returns tweets sent after the given date.
	End                string   `json:"end,omitempty"`                // Returns tweets sent before the given date.
	IncludeSearchTerms bool     `json:"includeSearchTerms,omitempty"` // If selected, a field will be added to each tweets about the search term that was used to find it.
	CustomMapFunction  string   `json:"customMapFunction,omitempty"`  // Function that takes each of the objects as argument and returns data that will be mapped by the function itself. This function is not intended for filtering, please don't use it for filtering purposes or you will get banned automatically.
}

type TweetSearchApifyAsyncResponse struct {
	ID string `json:"job_id"`
}

type TweetSearchApifyResponse struct {
	Type          string `json:"type"`
	ID            string `json:"id"`
	URL           string `json:"url"`
	TwitterURL    string `json:"twitterUrl"`
	Text          string `json:"text"`
	RetweetCount  int    `json:"retweetCount"`
	ReplyCount    int    `json:"replyCount"`
	LikeCount     int    `json:"likeCount"`
	QuoteCount    int    `json:"quoteCount"`
	CreatedAt     string `json:"createdAt"`
	Lang          string `json:"lang"`
	QuoteID       string `json:"quoteId"`
	BookmarkCount int    `json:"bookmarkCount"`
	IsReply       bool   `json:"isReply"`
	Source        string `json:"source"`
	IsRetweet     bool   `json:"isRetweet"`
	IsQuote       bool   `json:"isQuote"`
	Author        struct {
		Type           string `json:"type"`
		UserName       string `json:"userName"`
		URL            string `json:"url"`
		TwitterURL     string `json:"twitterUrl"`
		ID             string `json:"id"`
		Name           string `json:"name"`
		IsVerified     bool   `json:"isVerified"`
		IsBlueVerified bool   `json:"isBlueVerified"`
		ProfilePicture string `json:"profilePicture"`
		CoverPicture   string `json:"coverPicture"`
		Description    string `json:"description"`
		Location       string `json:"location"`
		Followers      int    `json:"followers"`
		Following      int    `json:"following"`
		Protected      bool   `json:"protected"`
		Status         string `json:"status"`
	}
}

func (*TweetSearchApifyAsyncResponse) Array(body io.Reader) ([]*TweetSearchApifyResponse, error) {
	r := make([]*TweetSearchApifyResponse, 0)
	if err := json.NewDecoder(body).Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}
