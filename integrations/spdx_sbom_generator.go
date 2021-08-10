package integrations

import (
	"encoding/json"
	"os"
)

type SpdxSBOMModel struct {
	Version    string `json:"Version,omitempty"`
	Name       string
	Path       string `json:"Path,omitempty"`
	LocalPath  string
	Supplier   SupplierContact
	PackageURL string
	// CheckSum                *CheckSum
	PackageHomePage         string
	PackageDownloadLocation string
	LicenseConcluded        string
	LicenseDeclared         string
	CommentsLicense         string
	// OtherLicense            []*License
	Copyright      string
	PackageComment string
	Root           bool
	// Modules                 map[string]*Module
}

// SupplierContact ...
type SupplierContact struct {
	Name  string
	Email string
}

func ParseSBOM(filename string) ([]SpdxSBOMModel, error) {
	var spdxModel []SpdxSBOMModel

	spdxFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(spdxFile).Decode(&spdxModel); err != nil {
		return nil, err
	}

	return spdxModel, nil
}
