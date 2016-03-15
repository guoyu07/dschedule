package api

import (
	"github.com/weibocom/dschedule/structs"
	"github.com/weibocom/dschedule/strategy"
	"github.com/weibocom/dschedule/crontab"
	"testing"
)

func TEST_SERVICE_ADD(t *testing.T) {
	srv := MakeHTTPServer(t)
	defer srv.Shutdown()

	service := &structs.Service{
		ServiceId:    "feed-1",
		ServiceType:  structs.ServiceTypeProd,
		StrategyName: structs.ServiceStrategyCrontab,
		StrategyConfig : []*strategy[{
			Time: "10:10",
			InstanceNum:3,
		},{

		}],
	}

}
