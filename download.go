package playstore

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func (p *Playstore) DownloadPackage(pkg string, versionCode int) ([]byte, error) {
	body := url.Values{}
	body.Add("ot", "1")
	body.Add("doc", pkg)
	body.Add("vc", strconv.Itoa(versionCode))

	d("Will download %s(%d)", pkg, versionCode)

	resp, err := p.makeApiRequest("/purchase", nil, []byte(body.Encode()))

	if err != nil {
		return nil, err
	}

	buyResp := resp.GetPayload().GetBuyResponse()

	if buyResp == nil {
		return nil, errors.New("Failed to obtain download details for package")
	}

	if buyResp.GetCheckoutinfo().GetItem().GetAmount().GetMicros() != 0 {
		return nil, errors.New("Cant download apps that arent free")
	}

	downloadData := buyResp.GetPurchaseStatusResponse().GetAppDeliveryData()

	if downloadData == nil {
		return nil, errors.New("Failed to obtain download details")
	}

	authCookies := downloadData.GetDownloadAuthCookie()

	req, err := http.NewRequest("GET", downloadData.GetDownloadUrl(), nil)

	if err != nil {
		return nil, err
	}

	for _, c := range authCookies {
		req.AddCookie(&http.Cookie{
			Name:  c.GetName(),
			Value: c.GetValue(),
		})
	}

	res, err := p.client.Do(req)

	if err != nil {
		d("Failed to execute download request for %s %e", pkg, err)
		return nil, err
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
