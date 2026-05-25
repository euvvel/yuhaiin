package config

type SniSpoofConfig struct {
    Enabled    bool   `json:"enabled"`
    TargetIP   string `json:"target_ip"`
    TargetPort int    `json:"target_port"`
    FakeSNI    string `json:"fake_sni"`
}

func DefaultSniSpoofConfig() *SniSpoofConfig {
    return &SniSpoofConfig{
        Enabled:    true,
        TargetIP:   "104.19.229.21",
        TargetPort: 443,
        FakeSNI:    "www.hcaptcha.com",
    }
}