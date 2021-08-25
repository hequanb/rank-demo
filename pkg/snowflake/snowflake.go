package snowflake

import (
	"time"
	
	"github.com/bwmarrin/snowflake"
)

// +--------------------------------------------------------------------------+
// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
// +--------------------------------------------------------------------------+

var node *snowflake.Node

// Init 由项目上线的日期开始计算分布式ID，由于存储41bit个时间戳，所以可以使用的最大时间是69年
func Init(startTime string, machineId int64) (err error) {
	var t time.Time
	t, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	// 缩减为毫秒级时间戳
	tx := t.UnixNano() / 1000000
	snowflake.Epoch = tx
	node, err = snowflake.NewNode(machineId)
	return
}

func GenId() int64 {
	return node.Generate().Int64()
}
