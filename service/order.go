package service

import (
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/zyfdegh/fanach/coreserver/entity"
	"github.com/zyfdegh/fanach/coreserver/util"
)

// CreateOrder creates an order
func CreateOrder(req entity.ReqPostOrder) (order *entity.Order, err error) {
	fmt.Printf("%+v", req)
	order = &entity.Order{}
	order.ID = genOrderID()
	return
}

func genOrderID() string {
	return fmt.Sprintf("%s%s", util.FormatTime(time.Now(), util.TimeLayout2), uniuri.New())
}
