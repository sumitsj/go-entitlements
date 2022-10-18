package contract

type UpdateRequest struct {
	Entitlements map[string]bool `json:"entitlements" binding:"required"`
	Reason       string          `json:"reason" binding:"required"`
}
