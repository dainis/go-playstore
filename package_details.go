package playstore

import (
	"errors"
	"net/url"
)

type AppDetails struct {
	Title         string
	Id            string
	Creator       string
	Mature        bool
	Images        []AppImage
	VersionCode   int
	VersionString string
	Developer     *Developer
}

type AppImage struct {
	Url  string
	Type int
}

type Developer struct {
	Name    string
	Email   string
	Website string
}

func (p *Playstore) PackageDetails(pkg string) (*AppDetails, error) {
	query := url.Values{}
	query.Add("doc", pkg)

	d("Will fetch data for package %s", pkg)

	resp, err := p.makeApiRequest("/details", query, nil)

	if err != nil {
		return nil, err
	}

	doc := resp.GetPayload().GetDetailsResponse().GetDocV2()

	appDetails := doc.GetDetails().GetAppDetails()

	if appDetails == nil {
		return nil, errors.New("Could not extract app details from package")
	}

	images := make([]AppImage, len(doc.Image))

	for j, i := range doc.Image {
		images[j] = AppImage{Url: i.GetImageUrl(), Type: int(i.GetImageType())}
	}

	return &AppDetails{
		Title:         doc.GetTitle(),
		Id:            doc.GetDocid(),
		Creator:       doc.GetCreator(),
		Mature:        doc.GetMature(),
		Images:        images,
		VersionCode:   int(appDetails.GetVersionCode()),
		VersionString: appDetails.GetVersionString(),
		Developer: &Developer{
			Name:    appDetails.GetDeveloperName(),
			Email:   appDetails.GetDeveloperEmail(),
			Website: appDetails.GetDeveloperWebsite(),
		},
	}, nil
}
