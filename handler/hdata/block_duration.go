// Copyright 2018 The QOS Authors

package hdata

import (
	"net/http"
	"strconv"
	"time"

	"github.com/QOSGroup/qmoon/handler/middleware"
	"github.com/QOSGroup/qmoon/lib"
	"github.com/QOSGroup/qmoon/service/block"
	"github.com/QOSGroup/qmoon/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const blockDurationUrl = "/blockDuration"

func init() {
	hdataHander[blockDurationUrl] = BlockDurationGinRegister
}

// BlockDurationGinRegister 注册blockDuration
func BlockDurationGinRegister(r *gin.Engine) {
	r.GET(nodeProxy+blockDurationUrl, middleware.ApiAuthGin(), blockDurationGin())
}

type blockDurationQuery struct {
	MinHeight int64 `pa:"minHeight"`
	MaxHeight int64 `json:"maxHeight"`
}

func (q blockDurationQuery) Validator() error {

	return nil
}

type blockDurationResp struct {
	LastHeight int64                    `json:"last_height"`
	Blocks     []*types.ResultBlockBase `json:"blocks"`
}

func blockDurationGin() gin.HandlerFunc {
	return func(c *gin.Context) {
		nt, err := getNodeFromUrl(c)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCMethodNotFoundError(""))
			return
		}

		var minHeight, maxHeight int64
		maxHeightStr := c.Query("maxHeight")
		maxHeight, err = strconv.ParseInt(maxHeightStr, 10, 64)
		if err != nil {
			maxHeight = 0
		}

		minHeightStr := c.Query("minHeight")
		minHeight, err = strconv.ParseInt(minHeightStr, 10, 64)
		if err != nil {
			minHeight = 1
		}

		lb, err := block.Latest(nt.ChanID)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCServerError("", err))
			return
		}

		if maxHeight == 0 {
			maxHeight = lb.Height
		}

		if maxHeight >= 50 {
			minHeight = maxHeight - 49
		} else {
			minHeight = 1
		}

		bs, err := block.Search(nt.ChanID, minHeight, maxHeight)
		if err != nil {
			c.JSON(http.StatusOK, types.RPCServerError("", err))
			return
		}
		var res []types.ResultBlockDuration

		for i := len(bs) - 1; i > 0; i-- {
			var d types.ResultBlockDuration
			d.Height = bs[i].Height
			d.Duration = int64(bs[i-1].Time.Sub(bs[i].Time).Nanoseconds()) / int64(time.Millisecond)
			res = append(res, d)
		}

		c.JSON(http.StatusOK, types.NewRPCSuccessResponse(lib.Cdc, "", res))
	}
}