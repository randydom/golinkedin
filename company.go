package linkedin

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type CompanyNode struct {
	Metadata Metadata             `json:"metadata,omitempty"`
	Elements []CompanyNodeElement `json:"elements,omitempty"`
	Paging   Paging               `json:"paging,omitempty"`
	Keywords string               `json:"keywords,omitempty"`

	err        error
	ln         *Linkedin
	stopCursor bool
}

type CompanyNodeElement struct {
	ExtendedElements []interface{} `json:"extendedElements,omitempty"`
	Elements         []Company     `json:"elements,omitempty"`
	Type             string        `json:"type,omitempty"`
}

// Company contain information about a company
type Company struct {
	Industry             map[string]Industry   `json:"industry,omitempty"`
	IndustryUrns         []string              `json:"industryUrns,omitempty"`
	AntiAbuseAnnotations []AntiAbuseAnnotation `json:"$anti_abuse_annotations,omitempty"`
	EntityUrn            string                `json:"entityUrn,omitempty"`
	MiniCompany          *MiniCompany          `json:"miniCompany,omitempty"`
	EmployeeCountRange   *EmployeeCountRange   `json:"employeeCountRange,omitempty"`
	Industries           []string              `json:"industries,omitempty"`
	Name                 string                `json:"name,omitempty"`
	Logo                 *Logo                 `json:"logo,omitempty"`
	RecipeType           string                `json:"$recipeType,omitempty"`
	UniversalName        string                `json:"universalName,omitempty"`
	URL                  string                `json:"url,omitempty"`
	ObjectUrn            string                `json:"objectUrn,omitempty"`
	Showcase             bool                  `json:"showcase,omitempty"`
	Active               bool                  `json:"active,omitempty"`
	TrackingID           string                `json:"trackingId,omitempty"`
	Image                *Image                `json:"image,omitempty"`
	Subtext              *Text                 `json:"subtext,omitempty"`
	TargetUrn            string                `json:"targetUrn,omitempty"`
	Text                 *Text                 `json:"text,omitempty"`
	DashTargetUrn        string                `json:"dashTargetUrn,omitempty"`
	Type                 string                `json:"type,omitempty"`

	TrackingUrn string    `json:"trackingUrn,omitempty"`
	Title       *Title    `json:"title,omitempty"`
	Headline    *Headline `json:"headline,omitempty"`
	Subline     *Headline `json:"subline,omitempty"`
}

type Headline struct {
	TextDirection string `json:"textDirection,omitempty"`
	Text          string `json:"text,omitempty"`
}

type MiniCompany struct {
	ObjectUrn     string `json:"objectUrn,omitempty"`
	EntityUrn     string `json:"entityUrn,omitempty"`
	Name          string `json:"name,omitempty"`
	Showcase      bool   `json:"showcase,omitempty"`
	Active        bool   `json:"active,omitempty"`
	Logo          *Logo  `json:"logo,omitempty"`
	UniversalName string `json:"universalName,omitempty"`
	TrackingID    string `json:"trackingId,omitempty"`
}

func (comp *CompanyNode) SetLinkedin(ln *Linkedin) {
	comp.ln = ln
}

func (comp *CompanyNode) Next() bool {
	if comp.stopCursor {
		return false
	}

	start := strconv.Itoa(comp.Paging.Start)
	count := strconv.Itoa(comp.Paging.Count)
	raw, err := comp.ln.get("/search/blended", url.Values{
		"keywords":     {comp.Keywords},
		"origin":       {OriginSwitchSearchVertical},
		"q":            {QAll},
		"filters":      {composeFilter(DefaultSearchCompanyFilter)},
		"queryContext": {composeFilter(DefaultSearchCompanyQueryContext)},
		"start":        {start},
		"count":        {count},
	})

	if err != nil {
		comp.err = err
		return false
	}

	compNode := new(CompanyNode)
	if err := json.Unmarshal(raw, compNode); err != nil {
		comp.err = err
		return false
	}

	comp.Elements = compNode.Elements
	comp.Paging.Start = compNode.Paging.Start + compNode.Paging.Count

	if len(comp.Elements) == 0 {
		return false
	}

	if len(comp.Elements[0].Elements) < comp.Paging.Count {
		comp.stopCursor = true
	}

	return true
}

func (comp *CompanyNode) Error() error {
	return comp.err
}
