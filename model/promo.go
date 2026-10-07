package model

// InviteCode 注册邀请码。
type InviteCode struct {
	Code      string `json:"code" gorm:"primaryKey"`
	Status    string `json:"status"` // unused | used
	UsedBy    string `json:"usedBy" gorm:"index"`
	UsedByName string `json:"usedByName" gorm:"-"`
	UsedAt    string `json:"usedAt"`
	CreatedAt string `json:"createdAt"`
}

// InviteCodeList 邀请码分页结果。
type InviteCodeList struct {
	Items []InviteCode `json:"items"`
	Total int          `json:"total"`
}

// RedeemCode 算力点兑换码。
type RedeemCode struct {
	Code       string  `json:"code" gorm:"primaryKey"`
	Credits    float64 `json:"credits" gorm:"type:decimal(20,2)"`
	Status     string  `json:"status"` // unused | used
	UsedBy     string  `json:"usedBy" gorm:"index"`
	UsedByName string  `json:"usedByName" gorm:"-"`
	UsedAt     string  `json:"usedAt"`
	CreatedAt  string  `json:"createdAt"`
}

// RedeemCodeList 兑换码分页结果。
type RedeemCodeList struct {
	Items []RedeemCode `json:"items"`
	Total int          `json:"total"`
}
