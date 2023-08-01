package round_robin

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"math"
	"strconv"
	"sync"
)

type WeightBalancer struct {
	connections []*weightConn
}

func (w *WeightBalancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	if len(w.connections) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var totalWeight uint32
	var res *weightConn
	for _, c := range w.connections {
		c.mutex.Lock()
		totalWeight = totalWeight + c.effectiveWeight
		c.currentWeight = c.currentWeight + c.effectiveWeight
		if res == nil || res.currentWeight < c.currentWeight {
			res = c
		}
		c.mutex.Unlock()
	}
	res.mutex.Lock()
	res.currentWeight = res.currentWeight - totalWeight
	res.mutex.Unlock()
	return balancer.PickResult{
		SubConn: res.subConn,
		Done: func(info balancer.DoneInfo) {
			res.mutex.Lock()
			if info.Err != nil && res.effectiveWeight == 0 {
				return
			}
			if info.Err == nil && res.effectiveWeight == math.MaxUint32 {
				return
			}
			if info.Err != nil {
				res.effectiveWeight--
			} else {
				res.effectiveWeight++
			}
			res.mutex.Unlock()
			//if atomic.CompareAndSwapUint32(&res.effectiveWeight, effectiveWeight, weight) {
			//	return
			//}
		},
	}, nil
}

type WeightBalancerBuilder struct {
}

func (w *WeightBalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	cs := make([]*weightConn, 0, len(info.ReadySCs))
	for sub, subInfo := range info.ReadySCs {
		weightStr := subInfo.Address.Attributes.Value("weight").(string)
		//if !ok || weightStr == "" {
		//	panic(fmt.Sprintf("没有权重"))
		//}
		weight, err := strconv.ParseUint(weightStr, 10, 64)
		if err != nil {
			panic(err)
		}
		cs = append(cs, &weightConn{
			subConn:         sub,
			weight:          uint32(weight),
			currentWeight:   uint32(weight),
			effectiveWeight: uint32(weight),
		})
	}
	return &WeightBalancer{
		connections: cs,
	}
}

type weightConn struct {
	mutex           sync.Mutex
	subConn         balancer.SubConn
	weight          uint32
	currentWeight   uint32
	effectiveWeight uint32
}
