package netapp

import (
	"encoding/xml"
	"net/http"
)

// LUN Operations

const (
	LunOnlineOperation  = "lun-online"
	LunOfflineOperation = "lun-offline"
	LunDestroyOperation = "lun-destroy"
	LunCreateOperation  = "lun-create-by-size"
)

type LunOperation struct {
	Base
	Params struct {
		XMLName           xml.Name
		LunPath           *lunPath
		LunInitiatorGroup *lunInitiatorGroup
		LunCreateOptions
	}
}

type lunPath struct {
	XMLName xml.Name
	Path    string `xml:",innerxml"`
}

type lunInitiatorGroup struct {
	XMLName        xml.Name
	InitiatorGroup string `xml:",innerxml"`
}

type LunCreateOptions struct {
	Class                   string `xml:"class,omitempty"`
	Comment                 string `xml:"comment,omitempty"`
	ForeignDisk             string `xml:"foreign-disk,omitempty"`
	InitiatorGroup          string `xml:"initiator-group,omitempty"`
	OsType                  string `xml:"ostype,omitempty"`
	Path                    string `xml:"path,omitempty"`
	PrefixSize              string `xml:"prefix-size,omitempty"`
	QosAdaptivePolicyGroup  string `xml:"qos-adaptive-policy-group,omitempty"`
	QosPolicyGroup          string `xml:"qos-policy-group,omitempty"`
	Size                    int64  `xml:"size,omitempty"`
	SpaceAllocationEnabled  bool   `xml:"space-allocation-enabled,omitempty"`
	SpaceReservationEnabled bool   `xml:"space-reservation-enabled,omitempty"`
	UseExactSize            bool   `xml:"use-exact-size,omitempty"`
}

func (l LunOperation) Create(vserverName string, options *LunCreateOptions) (*SingleResultResponse, *http.Response, error) {
	l.Params.XMLName = xml.Name{Local: LunCreateOperation}
	l.Name = vserverName
	l.Params.LunCreateOptions = *options
	r := SingleResultResponse{}
	res, err := l.get(l, &r)
	return &r, res, err
}

func (l LunOperation) Map(vserverName string, lunPathName string, initiatorGroup string) (*SingleResultResponse, *http.Response, error) {
	l.Params.XMLName = xml.Name{Local: LunMapOperation}
	l.Name = vserverName
	l.Params.LunInitiatorGroup = &lunInitiatorGroup{
		XMLName: xml.Name{Local: "initiator-group"},
		Path:    initiatorGroup,
	}
	l.Params.LunPath = &lunPath{
		XMLName: xml.Name{Local: "path"},
		Path:    lunPathName,
	}
	l.Params.LunPath = lunPath
}

func (l LunOperation) Operation(vserverName string, lunPathName string, operation string) (*SingleResultResponse, *http.Response, error) {
	l.Params.XMLName = xml.Name{Local: operation}
	l.Name = vserverName
	elementName := "path"
	l.Params.LunPath = &lunPath{
		XMLName: xml.Name{Local: elementName},
		Path:    lunPathName,
	}
	r := SingleResultResponse{}
	res, err := l.get(l, &r)
	return &r, res, err
}
