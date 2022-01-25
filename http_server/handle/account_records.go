package handle

import (
	"das_register_server/http_server/api_code"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqAccountRecords struct {
	Account string `json:"account"`
}

type RespAccountRecords struct {
	Records []RespAccountRecordsData `json:"records"`
}

type RespAccountRecordsData struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Value string `json:"value"`
	Ttl   string `json:"ttl"`
}

func (h *HttpHandle) RpcAccountRecords(p json.RawMessage, apiResp *api_code.ApiResp) {
	var req []ReqAccountRecords
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Error("json.Unmarshal err:", err.Error())
		apiResp.ApiRespErr(api_code.ApiCodeParamsInvalid, "params invalid")
		return
	} else if len(req) == 0 {
		log.Error("len(req) is 0")
		apiResp.ApiRespErr(api_code.ApiCodeParamsInvalid, "params invalid")
		return
	}

	if err = h.doAccountRecords(&req[0], apiResp); err != nil {
		log.Error("doVersion err:", err.Error())
	}
}

func (h *HttpHandle) AccountRecords(ctx *gin.Context) {
	var (
		funcName = "AccountRecords"
		clientIp = GetClientIp(ctx)
		req      ReqAccountRecords
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

	if err = h.doAccountRecords(&req, &apiResp); err != nil {
		log.Error("doAccountRecords err:", err.Error(), funcName, clientIp)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAccountRecords(req *ReqAccountRecords, apiResp *api_code.ApiResp) error {
	var resp RespAccountRecords

	resp.Records = make([]RespAccountRecordsData, 0)
	list, err := h.dbDao.SearchRecordsByAccount(req.Account)
	if err != nil {
		apiResp.ApiRespErr(api_code.ApiCodeDbError, "search records err")
		return fmt.Errorf("SearchRecordsByAccount err: %s", err.Error())
	}
	for _, v := range list {
		resp.Records = append(resp.Records, RespAccountRecordsData{
			Key:   v.Key,
			Type:  v.Type,
			Label: v.Label,
			Value: v.Value,
			Ttl:   v.Ttl,
		})
	}

	apiResp.ApiRespOK(resp)
	return nil
}
