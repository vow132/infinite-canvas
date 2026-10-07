package repository

import (
	"github.com/vow132/infinite-canvas/model"
	"gorm.io/gorm"
)

// SaveInviteCode 保存邀请码。
func SaveInviteCode(code model.InviteCode) error {
	db, err := DB()
	if err != nil {
		return err
	}
	return db.Save(&code).Error
}

// GetInviteCode 按码读取邀请码。
func GetInviteCode(code string) (model.InviteCode, bool, error) {
	db, err := DB()
	if err != nil {
		return model.InviteCode{}, false, err
	}
	var item model.InviteCode
	if err := db.Where("code = ?", code).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.InviteCode{}, false, nil
		}
		return model.InviteCode{}, false, err
	}
	return item, true, nil
}

// ListInviteCodes 分页返回邀请码。
func ListInviteCodes(q model.Query) ([]model.InviteCode, int64, error) {
	db, err := DB()
	if err != nil {
		return nil, 0, err
	}
	tx := db.Model(&model.InviteCode{})
	if q.Type != "" {
		tx = tx.Where("status = ?", q.Type)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	q.Normalize()
	var items []model.InviteCode
	if err := tx.Order("created_at desc").Offset(q.Offset()).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// DeleteInviteCodes 批量删除邀请码。
func DeleteInviteCodes(codes []string) error {
	if len(codes) == 0 {
		return nil
	}
	db, err := DB()
	if err != nil {
		return err
	}
	return db.Where("code IN ?", codes).Delete(&model.InviteCode{}).Error
}

// ClaimInviteCode 原子认领一个未使用的邀请码。
func ClaimInviteCode(code string, userID string, usedAt string) (bool, error) {
	db, err := DB()
	if err != nil {
		return false, err
	}
	result := db.Model(&model.InviteCode{}).
		Where("code = ? AND status = ?", code, "unused").
		Updates(map[string]any{"status": "used", "used_by": userID, "used_at": usedAt})
	return result.RowsAffected > 0, result.Error
}

// SaveRedeemCode 保存兑换码。
func SaveRedeemCode(code model.RedeemCode) error {
	db, err := DB()
	if err != nil {
		return err
	}
	return db.Save(&code).Error
}

// GetRedeemCode 按码读取兑换码。
func GetRedeemCode(code string) (model.RedeemCode, bool, error) {
	db, err := DB()
	if err != nil {
		return model.RedeemCode{}, false, err
	}
	var item model.RedeemCode
	if err := db.Where("code = ?", code).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.RedeemCode{}, false, nil
		}
		return model.RedeemCode{}, false, err
	}
	return item, true, nil
}

// ListRedeemCodes 分页返回兑换码。
func ListRedeemCodes(q model.Query) ([]model.RedeemCode, int64, error) {
	db, err := DB()
	if err != nil {
		return nil, 0, err
	}
	tx := db.Model(&model.RedeemCode{})
	if q.Type != "" {
		tx = tx.Where("status = ?", q.Type)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	q.Normalize()
	var items []model.RedeemCode
	if err := tx.Order("created_at desc").Offset(q.Offset()).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// DeleteRedeemCodes 批量删除兑换码。
func DeleteRedeemCodes(codes []string) error {
	if len(codes) == 0 {
		return nil
	}
	db, err := DB()
	if err != nil {
		return err
	}
	return db.Where("code IN ?", codes).Delete(&model.RedeemCode{}).Error
}

// ClaimRedeemCode 原子认领一个未使用的兑换码。
func ClaimRedeemCode(code string, userID string, usedAt string) (model.RedeemCode, bool, error) {
	db, err := DB()
	if err != nil {
		return model.RedeemCode{}, false, err
	}
	result := db.Model(&model.RedeemCode{}).
		Where("code = ? AND status = ?", code, "unused").
		Updates(map[string]any{"status": "used", "used_by": userID, "used_at": usedAt})
	if result.Error != nil {
		return model.RedeemCode{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return model.RedeemCode{}, false, nil
	}
	item, ok, err := GetRedeemCode(code)
	if err != nil || !ok {
		return model.RedeemCode{}, false, err
	}
	return item, true, nil
}

// ListUserCreditLogs 分页返回指定用户的算力点流水。
func ListUserCreditLogs(userID string, q model.Query) ([]model.CreditLog, int64, error) {
	db, err := DB()
	if err != nil {
		return nil, 0, err
	}
	tx := db.Model(&model.CreditLog{}).Where("user_id = ?", userID)
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	q.Normalize()
	var items []model.CreditLog
	if err := tx.Order("created_at desc").Offset(q.Offset()).Limit(q.PageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
