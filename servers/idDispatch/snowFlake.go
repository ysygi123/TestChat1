package idDispatch

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits  uint8  = 10                      //10bit工作机器的id，如果你发现1024台机器不够那就调大次值
	numberBits  uint8  = 12                      //12bit 工作序号，如果你发现1毫秒并发生成4096个唯一id不够请调大次值
	workerMax   uint64 = -1 ^ (-1 << workerBits) // 1024
	numberMax   uint64 = -1 ^ (-1 << numberBits) // 4096
	timeShift   uint8  = workerBits + numberBits
	workerShift uint8  = numberBits
	// 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID，
	//这个值请自行设置为你系统准备上线前的精确到毫秒级别的时间戳，因为雪花时间戳保证唯一的部分最多管69年（2的41次方），
	//所以此值设置为你当前时间戳能够保证你的系统是从当前时间开始往后推69年
	startTime uint64 = 1626399550000
)

var SnowFlakeWorker *Worker

type Worker struct {
	mu        sync.Mutex
	timestamp uint64
	workerId  uint64
	number    uint64
}

func NewWorker(workerId uint64) error {
	if workerId < 0 || workerId > workerMax {
		return errors.New("Worker ID excess of quantity")
	}
	// 生成一个新节点
	SnowFlakeWorker = &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}
	return nil
}

func (w *Worker) GetId() uint64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := uint64(time.Now().UnixNano() / 1e6)
	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = uint64(time.Now().UnixNano() / 1e6)
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	//以下表达式才是主菜
	//  (now - startTime) << timeShift   产生了 41 + （10 + 12）的效应但却并不保证唯一
	//  | (w.workerId << workerShift)  保证了与其他机器不重复
	//  | (w.number))  保证了自己这台机不会重复
	ID := uint64((now-startTime)<<timeShift | (w.workerId << workerShift) | (w.number))
	return ID
}
