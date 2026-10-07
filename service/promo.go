package service

import (
	"errors"
	"crypto/rand"
	"math/big"
	"strings"

	"github.com/vow132/infinite-canvas/model"
	"github.com/vow132/infinite-canvas/repository"
)

const promoCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func newPromoCode(prefix string) string {
	buf := make([]byte, 12)
	for i := range buf {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(promoCodeAlphabet))))
		if err != nil {
			buf[i] = promoCodeAlphabet[i%len(promoCodeAlphabet)]
			continue
		}
		buf[i] = promoCodeAlphabet[n.Int64()]
	}
	return prefix + "-" + string(buf)
}

func enrichInviteCodeNames(items []model.InviteCode) []model.InviteCode {
	ids := map[string]bool{}
	for _, item := range items {
		if item.UsedBy != "" {
			ids[item.UsedBy] = true
		}
	}
	names := map[string]string{}
	for id := range ids {
		if user, ok, err := repository.GetUserByID(id); err == nil && ok {
			names[id] = user.Username
		}
	}
	for i := range items {
		items[i].UsedByName = names[items[i].UsedBy]
	}
	return items
}

func enrichRedeemCodeNames(items []model.RedeemCode) []model.RedeemCode {
	ids := map[string]bool{}
	for _, item := range items {
		if item.UsedBy != "" {
			ids[item.UsedBy] = true
		}
	}
	names := map[string]string{}
	for id := range ids {
		if user, ok, err := repository.GetUserByID(id); err == nil && ok {
			names[id] = user.Username
		}
	}
	for i := range items {
		items[i].UsedByName = names[items[i].UsedBy]
	}
	return items
}

// AdminGenerateInviteCodes 批量生成注册邀请码。
func AdminGenerateInviteCodes(count int) ([]string, error) {
	if count < 1 {
		count = 1
	}
	if count > 500 {
		count = 500
	}
	codes := make([]string, 0, count)
	for i := 0; i < count; i++ {
		code := newPromoCode("INV")
		_, ok, err := repository.GetInviteCode(code)
		if err != nil {
			return nil, err
		}
		if ok {
			i--
			continue
		}
		if err := repository.SaveInviteCode(model.InviteCode{Code: code, Status: "unused", CreatedAt: now()}); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}

// AdminListInviteCodes 分页返回邀请码。
func AdminListInviteCodes(q model.Query) (model.InviteCodeList, error) {
	items, total, err := repository.ListInviteCodes(q)
	if err != nil {
		return model.InviteCodeList{}, err
	}
	return model.InviteCodeList{Items: enrichInviteCodeNames(items), Total: int(total)}, nil
}

// AdminDeleteInviteCodes 批量删除邀请码。
func AdminDeleteInviteCodes(codes []string) error {
	return repository.DeleteInviteCodes(codes)
}

// AdminUnusedInviteCodes 返回所有未使用的邀请码，用于导出。
func AdminUnusedInviteCodes() ([]string, int, error) {
	items, total, err := repository.ListInviteCodes(model.Query{Page: 1, PageSize: model.MaxPageSize, Type: "unused"})
	if err != nil {
		return nil, 0, err
	}
	if total > model.MaxPageSize {
		return nil, 0, errors.New("未使用邀请码超过单次导出上限")
	}
	codes := make([]string, 0, len(items))
	for _, item := range items {
		codes = append(codes, item.Code)
	}
	return codes, int(total), nil
}

// AdminGenerateRedeemCodes 批量生成指定额度的算力点兑换码。
func AdminGenerateRedeemCodes(count int, credits float64) ([]string, error) {
	credits = normalizeCredits(credits)
	if credits <= 0 {
		return nil, safeMessageError{message: "兑换码额度必须大于 0"}
	}
	if count < 1 {
		count = 1
	}
	if count > 500 {
		count = 500
	}
	codes := make([]string, 0, count)
	for i := 0; i < count; i++ {
		code := newPromoCode("RDM")
		_, ok, err := repository.GetRedeemCode(code)
		if err != nil {
			return nil, err
		}
		if ok {
			i--
			continue
		}
		if err := repository.SaveRedeemCode(model.RedeemCode{Code: code, Credits: credits, Status: "unused", CreatedAt: now()}); err != nil {
			return nil, err
		}
		codes = append(codes, code)
	}
	return codes, nil
}

// AdminListRedeemCodes 分页返回兑换码。
func AdminListRedeemCodes(q model.Query) (model.RedeemCodeList, error) {
	items, total, err := repository.ListRedeemCodes(q)
	if err != nil {
		return model.RedeemCodeList{}, err
	}
	return model.RedeemCodeList{Items: enrichRedeemCodeNames(items), Total: int(total)}, nil
}

// AdminDeleteRedeemCodes 批量删除兑换码。
func AdminDeleteRedeemCodes(codes []string) error {
	return repository.DeleteRedeemCodes(codes)
}

// AdminUnusedRedeemCodes 返回所有未使用的兑换码，用于导出。
func AdminUnusedRedeemCodes() ([]model.RedeemCode, error) {
	items, total, err := repository.ListRedeemCodes(model.Query{Page: 1, PageSize: model.MaxPageSize, Type: "unused"})
	if err != nil {
		return nil, err
	}
	if total > model.MaxPageSize {
		return nil, errors.New("未使用兑换码超过单次导出上限")
	}
	return items, nil
}

// RedeemCodeForUser 用户兑换算力点。
func RedeemCodeForUser(userID string, code string) (model.User, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return model.User{}, safeMessageError{message: "请填写兑换码"}
	}
	claimed, claimedOK, err := repository.ClaimRedeemCode(code, userID, now())
	if err != nil {
		return model.User{}, err
	}
	if !claimedOK {
		return model.User{}, safeMessageError{message: "兑换码无效或已被使用"}
	}
	user, ok, err := repository.GetUserByID(userID)
	if err != nil || !ok {
		if err == nil {
			// 释放误占用的兑换码
			_ = repository.SaveRedeemCode(model.RedeemCode{Code: claimed.Code, Credits: claimed.Credits, Status: "unused", CreatedAt: claimed.CreatedAt})
		}
		if err != nil {
			return model.User{}, err
		}
		return model.User{}, safeMessageError{message: "用户不存在"}
	}
	user.Credits = normalizeCredits(user.Credits + claimed.Credits)
	user.UpdatedAt = now()
	user, err = repository.SaveUser(user)
	if err != nil {
		// 保存失败时回滚兑换码占用，避免用户丢失额度
		_ = repository.SaveRedeemCode(model.RedeemCode{Code: claimed.Code, Credits: claimed.Credits, Status: "unused", CreatedAt: claimed.CreatedAt})
		return model.User{}, err
	}
	user.Password = ""
	if _, err := repository.SaveCreditLog(model.CreditLog{
		ID:        newID("credit"),
		UserID:    user.ID,
		Type:      model.CreditLogTypeRedeem,
		Amount:    claimed.Credits,
		Balance:   user.Credits,
		RelatedID: claimed.Code,
		Remark:    "兑换码兑换",
		CreatedAt: now(),
	}); err != nil {
		return user, err
	}
	return user, nil
}

// ListUserCreditLogs 返回指定用户的算力点流水。
func ListUserCreditLogs(userID string, q model.Query) (model.CreditLogList, error) {
	items, total, err := repository.ListUserCreditLogs(userID, q)
	if err != nil {
		return model.CreditLogList{}, err
	}
	return model.CreditLogList{Items: items, Total: int(total)}, nil
}
