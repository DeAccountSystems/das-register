package handle

import (
	"das_register_server/http_server/api_code"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"time"
)

func (h *HttpHandle) Query(ctx *gin.Context) {
	var (
		req       api_code.JsonRequest
		resp      api_code.JsonResponse
		apiResp   api_code.ApiResp
		clientIp  = GetClientIp(ctx)
		startTime = time.Now()
	)
	resp.Result = &apiResp

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Error("ShouldBindJSON err:", err.Error())
		apiResp.ApiRespErr(api_code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp.ID, resp.JsonRpc = req.ID, req.JsonRpc
	log.Info("Query:", req.Method, clientIp, toolib.JsonString(req))

	switch req.Method {
	case api_code.MethodTokenList:
		h.RpcTokenList(req.Params, &apiResp)
	case api_code.MethodConfigInfo:
		h.RpcConfigInfo(req.Params, &apiResp)
	case api_code.MethodAccountList:
		h.RpcAccountList(req.Params, &apiResp)
	case api_code.MethodAccountMine:
		h.RpcAccountMine(req.Params, &apiResp)
	case api_code.MethodAccountDetail:
		h.RpcAccountDetail(req.Params, &apiResp)
	case api_code.MethodAccountRecords:
		h.RpcAccountRecords(req.Params, &apiResp)
	case api_code.MethodReverseLatest:
		h.RpcReverseLatest(req.Params, &apiResp)
	case api_code.MethodReverseList:
		h.RpcReverseList(req.Params, &apiResp)
	case api_code.MethodTransactionStatus:
		h.RpcTransactionStatus(req.Params, &apiResp)
	case api_code.MethodBalanceInfo:
		h.RpcBalanceInfo(req.Params, &apiResp)
	case api_code.MethodTransactionList:
		h.RpcTransactionList(req.Params, &apiResp)
	case api_code.MethodRewardsMine:
		h.RpcRewardsMine(req.Params, &apiResp)
	case api_code.MethodWithdrawList:
		h.RpcWithdrawList(req.Params, &apiResp)
	case api_code.MethodAccountSearch:
		h.RpcAccountSearch(req.Params, &apiResp)
	case api_code.MethodRegisteringList:
		h.RpcRegisteringList(req.Params, &apiResp)
	case api_code.MethodOrderDetail:
		h.RpcOrderDetail(req.Params, &apiResp)
	default:
		log.Error("method not exist:", req.Method)
		apiResp.ApiRespErr(api_code.ApiCodeMethodNotExist, fmt.Sprintf("method [%s] not exits", req.Method))
	}

	api_code.DoMonitorLogRpc(&apiResp, req.Method, clientIp, startTime)

	ctx.JSON(http.StatusOK, resp)
	return
}
