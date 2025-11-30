package vo

// OIDCMetadata OIDC提供商元数据结构
type OIDCMetadata struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	RevocationEndpoint                string   `json:"revocation_endpoint,omitempty"`
	IntrospectionEndpoint             string   `json:"introspection_endpoint,omitempty"`
	JwksURI                           string   `json:"jwks_uri"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	ScopesSupported                   []string `json:"scopes_supported,omitempty"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported,omitempty"`
	ClaimsSupported                   []string `json:"claims_supported,omitempty"`
}

// OIDCTokenResponse OIDC token响应结构
type OIDCTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token,omitempty"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

// OIDCUserInfo OIDC用户信息结构
type OIDCUserInfo struct {
	Sub               string `json:"sub"`  // 用户唯一标识
	Name              string `json:"name"` // 用户名称
	NickName          string `json:"nickname,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Email             string `json:"email,omitempty"`
	Picture           string `json:"picture,omitempty"`
}

// OIDCConfigResponse OIDC配置响应结构
type OIDCConfigResponse struct {
	Enabled     bool   `json:"enabled"`     // OIDC是否启用
	AuthURL     string `json:"authURL"`     // 授权URL
	State       string `json:"state"`       // 状态参数
	ConfigValid bool   `json:"configValid"` // 配置是否有效
}
