package models

type Provider struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

type ProviderVersions struct {
	Namespace string      `json:"namespace"`
	Name      string      `json:"name"`
	Version   string      `json:"version"`
	Platforms []Platforms `json:"platforms"`
}

type Platforms struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type GPGKeys struct {
	KeyID          string `json:"key_id"`
	ASCIIArmor     string `json:"ascii_armor"`
	TrustSignature string `json:"trust_signature"`
	Source         string `json:"source"`
	SourceURL      string `json:"source_url"`
}

type ProviderDownload struct {
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	Filename      string `json:"filename"`
	DownloadURL   string `json:"download_url"`
	SHASumsURL    string `json:"shasums_url"`
	SHASumsSigURL string `json:"shasums_signature_url"`
	SHASum        string `json:"shasum"`
	SigningKeys   struct {
		GPGPublicKeys []GPGKeys `json:"gpg_public_keys"`
	} `json:"signing_keys"`
}
