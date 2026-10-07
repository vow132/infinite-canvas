package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/vow132/infinite-canvas/model"
	"github.com/vow132/infinite-canvas/service"
)

type adminGenerateInviteRequest struct {
	Count int `json:"count"`
}

type adminGenerateRedeemRequest struct {
	Count   int     `json:"count"`
	Credits float64 `json:"credits"`
}

type adminPromoDeleteRequest struct {
	Codes []string `json:"codes"`
}

type userRedeemRequest struct {
	Code string `json:"code"`
}

func parseIntParam(r *http.Request, name string, fallback int) int {
	value := r.URL.Query().Get(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func AdminGenerateInviteCodes(w http.ResponseWriter, r *http.Request) {
	var request adminGenerateInviteRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	codes, err := service.AdminGenerateInviteCodes(request.Count)
	if err != nil {
		FailError(w, err)
		return
	}
	OK(w, codes)
}

func AdminListInviteCodes(w http.ResponseWriter, r *http.Request) {
	query := model.Query{
		Page:     parseIntParam(r, "page", 1),
		PageSize: parseIntParam(r, "pageSize", 20),
		Type:     strings.TrimSpace(r.URL.Query().Get("status")),
	}
	result, err := service.AdminListInviteCodes(query)
	if err != nil {
		FailError(w, err)
		return
	}
	OK(w, result)
}

func AdminDeleteInviteCodes(w http.ResponseWriter, r *http.Request) {
	var request adminPromoDeleteRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	if err := service.AdminDeleteInviteCodes(request.Codes); err != nil {
		FailError(w, err)
		return
	}
	OK(w, true)
}

func AdminExportInviteCodes(w http.ResponseWriter, r *http.Request) {
	codes, _, err := service.AdminUnusedInviteCodes()
	if err != nil {
		FailError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=invite-codes.txt")
	_, _ = w.Write([]byte(strings.Join(codes, "\n")))
}

func AdminGenerateRedeemCodes(w http.ResponseWriter, r *http.Request) {
	var request adminGenerateRedeemRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	codes, err := service.AdminGenerateRedeemCodes(request.Count, request.Credits)
	if err != nil {
		FailError(w, err)
		return
	}
	OK(w, codes)
}

func AdminListRedeemCodes(w http.ResponseWriter, r *http.Request) {
	query := model.Query{
		Page:     parseIntParam(r, "page", 1),
		PageSize: parseIntParam(r, "pageSize", 20),
		Type:     strings.TrimSpace(r.URL.Query().Get("status")),
	}
	result, err := service.AdminListRedeemCodes(query)
	if err != nil {
		FailError(w, err)
		return
	}
	OK(w, result)
}

func AdminDeleteRedeemCodes(w http.ResponseWriter, r *http.Request) {
	var request adminPromoDeleteRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	if err := service.AdminDeleteRedeemCodes(request.Codes); err != nil {
		FailError(w, err)
		return
	}
	OK(w, true)
}

func AdminExportRedeemCodes(w http.ResponseWriter, r *http.Request) {
	codes, err := service.AdminUnusedRedeemCodes()
	if err != nil {
		FailError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=redeem-codes.txt")
	lines := make([]string, 0, len(codes))
	for _, item := range codes {
		lines = append(lines, item.Code+"  "+strconv.FormatFloat(item.Credits, 'f', -1, 64))
	}
	_, _ = w.Write([]byte(strings.Join(lines, "\n")))
}

func UserRedeemCode(w http.ResponseWriter, r *http.Request) {
	user, ok := service.UserFromContext(r.Context())
	if !ok {
		FailWithStatus(w, http.StatusUnauthorized, "未登录或权限不足")
		return
	}
	var request userRedeemRequest
	_ = json.NewDecoder(r.Body).Decode(&request)
	result, err := service.RedeemCodeForUser(user.ID, request.Code)
	if err != nil {
		FailError(w, err)
		return
	}
	OK(w, result)
}

func UserCreditLogs(w http.ResponseWriter, r *http.Request) {
	user, ok := service.UserFromContext(r.Context())
	if !ok {
		FailWithStatus(w, http.StatusUnauthorized, "未登录或权限不足")
		return
	}
	query := model.Query{
		Page:     parseIntParam(r, "page", 1),
		PageSize: parseIntParam(r, "pageSize", 20),
	}
	result, err := service.ListUserCreditLogs(user.ID, query)
	if err != nil {
		FailError(w, err)
		return
	}
	OK(w, result)
}
