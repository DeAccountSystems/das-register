package handle

import (
	"das_database/dao"
	"das_register_server/http_server/api_code"
	"das_register_server/tables"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

const (
	ReserveStatusNeverSet  = 0
	ReverseStatusOldOnly   = 1
	ReverseStatusOldAndNew = 2
	ReverseStatusNewOnly   = 3
)

type ReqReverseInfo struct {
	core.ChainTypeAddress
}

type RespReverseInfo struct {
	Account       string `json:"account"`
	IsValid       bool   `json:"is_valid"`
	ReserveStatus uint32 `json:"reserve_status"`
}

func (h *HttpHandle) ReverseInfo(ctx *gin.Context) {
	var (
		funcName = "ReverseInfo"
		clientIp = GetClientIp(ctx)
		req      ReqReverseInfo
		apiResp  api_code.ApiResp
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, clientIp)
		apiResp.ApiRespErr(api_code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, clientIp, toolib.JsonString(req))

	if err = h.doReverseInfo(&req, &apiResp); err != nil {
		log.Errorf("doReverseInfo err: %s funcName: %s clientIp: %s", err, funcName, clientIp)
	}
	ctx.JSON(http.StatusOK, apiResp)
}

// doReverseStatus
func (h *HttpHandle) doReverseInfo(req *ReqReverseInfo, apiResp *api_code.ApiResp) error {
	res := checkReqKeyInfo(h.dasCore.Daf(), &req.ChainTypeAddress, apiResp)
	reverse, err := h.dbDao.SearchLatestReverse(res.ChainType, res.AddressHex)
	if err != nil {
		return fmt.Errorf("SearchLatestReverse err: %s", err)
	}

	resp := &RespReverseInfo{}
	if reverse.Id == 0 {
		apiResp.ApiRespOK(resp)
		return nil
	}
	resp.Account = reverse.Account

	reverseOld, err := h.dbDao.SearchLatestReverseByType(res.ChainType, res.AddressHex, dao.ReverseTypeOld)
	if err != nil {
		return fmt.Errorf("SearchLatestReverseByType err: %s", err)
	}
	reverseNew, err := h.dbDao.SearchLatestReverseByType(res.ChainType, res.AddressHex, dao.ReverseTypeSmt)
	if err != nil {
		return fmt.Errorf("SearchLatestReverseByType err: %s", err)
	}
	if reverseOld.Id > 0 && reverseNew.Id == 0 {
		resp.ReserveStatus = ReverseStatusOldOnly
	} else if reverseOld.Id > 0 && reverseNew.Id > 0 {
		resp.ReserveStatus = ReverseStatusOldAndNew
	} else if reverseOld.Id == 0 && reverseNew.Id > 0 {
		resp.ReserveStatus = ReverseStatusNewOnly
	}

	// account
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(reverse.Account))
	acc, err := h.dbDao.GetAccountInfoByAccountId(accountId)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			apiResp.ApiRespErr(api_code.ApiCodeDbError, "search account err")
			return fmt.Errorf("SearchAccount err: %s", err.Error())
		}
	}
	if acc.Id == 0 || acc.Status == tables.AccountStatusOnCross {
		apiResp.ApiRespOK(resp)
		return nil
	}

	if strings.EqualFold(res.AddressHex, acc.Owner) || strings.EqualFold(res.AddressHex, acc.Manager) {
		resp.IsValid = true
		apiResp.ApiRespOK(resp)
		return nil
	}

	// records
	record, err := h.dbDao.SearchAccountReverseRecords(acc.Account, res.AddressHex)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			*apiResp = api_code.ApiRespErr(api_code.ApiCodeDbError, "search account err")
			return fmt.Errorf("SearchAccount err: %s", err.Error())
		}
	}
	if record.Id > 0 {
		resp.IsValid = true
	}
	apiResp.ApiRespOK(resp)
	return nil
}
