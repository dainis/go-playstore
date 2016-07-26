package playstore

import (
	"bytes"
	"errors"
	play "github.com/dainis/go-playstore/protobuf"
	proto "github.com/golang/protobuf/proto"
	"github.com/tj/go-debug"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var d = debug.Debug("playstore")

type Playstore struct {
	authToken string
	deviceId  string
	client    *http.Client
}

func New(email string, password string, deviceId string) (*Playstore, error) {
	postBody := url.Values{}

	postBody.Add("Email", email)
	postBody.Add("Passwd", password)
	postBody.Add("service", "androidmarket")
	postBody.Add("accountType", "HOSTED_OR_GOOGLE")
	postBody.Add("has_permission", "1")
	postBody.Add("source", "android")
	postBody.Add("androidId", deviceId)
	postBody.Add("app", "com.android.vending")
	postBody.Add("device_country", DEVICE_COUNTRY)
	postBody.Add("operatorCountry", DEVICE_COUNTRY)
	postBody.Add("lang", DEVICE_COUNTRY)
	postBody.Add("sdk_version", SDK_VERSION)

	resp, err := http.PostForm(LOGIN_URL, postBody)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	loginResult := parseBodyResponse(respBody)

	auth, ok := loginResult["auth"]

	if !ok {
		d("Could not log in in playstore")
		return nil, errors.New("Could not log in")
	}

	return &Playstore{authToken: auth, deviceId: deviceId, client: &http.Client{}}, nil
}

func (p *Playstore) getBasicHeaders() *http.Header {
	headers := &http.Header{}

	headers.Add("Accept-Language", "en_US")
	headers.Add("Authorization", "GoogleLogin auth="+p.authToken)
	headers.Add("X-DFE-Enabled-Experiments", "cl:billing.select_add_instrument_by_default")
	headers.Add("X-DFE-Unsupported-Experiments", "nocache:billing.use_charging_poller,market_emails,buyer_currency,prod_baseline,checkin.set_asset_paid_app_field,shekel_test,content_ratings,buyer_currency_in_app,nocache:encrypted_apk,recent_changes")
	headers.Add("X-DFE-Device-Id", p.deviceId)
	headers.Add("X-DFE-Client-Id", "am-android-google")
	headers.Add("User-Agent", USER_AGENT)
	headers.Add("X-DFE-SmallestScreenWidthDp", "320")
	headers.Add("X-DFE-Filter-Level", "3")
	headers.Add("Host", "android.clients.google.com")

	return headers
}

//Poor man's body parser
func parseBodyResponse(response []byte) map[string]string {
	s := string(response)

	lines := strings.Split(s, "\n")

	result := make(map[string]string)

	for _, line := range lines {
		values := strings.Split(line, "=")

		if len(values) != 2 {
			continue
		}

		result[strings.ToLower(strings.TrimSpace(values[0]))] = strings.TrimSpace(values[1])
	}

	return result
}

func (p *Playstore) makeApiRequest(endpoint string, query url.Values, body []byte) (*play.ResponseWrapper, error) {
	var method string

	headers := p.getBasicHeaders()

	if body != nil {
		method = "POST"
		headers.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	} else {
		method = "GET"
	}

	u := API_BASE_PATH + endpoint

	if query != nil && query.Encode() != "" {
		u += "?" + query.Encode()
	}

	d("Will make request to %s : %s", method, u)

	req, err := http.NewRequest(method, u, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header = *headers
	res, err := p.client.Do(req)

	if err != nil {
		d("Failed to make request to %s %s", u, err)
		return nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		d("Failed to read response from %s %s", u, err)
		return nil, err
	}

	respWrapper := &play.ResponseWrapper{}

	if err := proto.Unmarshal(b, respWrapper); err != nil {
		d("Failed to unmarshall response %s %s", u, err)
		return nil, err
	}

	e := respWrapper.GetCommands().GetDisplayErrorMessage()

	if e != "" {
		d("Playstore returned its own error for %s %s", u, e)
		return nil, errors.New(e)
	}

	return respWrapper, nil
}
