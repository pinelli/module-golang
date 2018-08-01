package orderbook

import (
	"sort"
)

type RestOrder struct {
	order *Order
	next  *RestOrder
}

type RestOrderList struct {
	list []*Order
}

type Orderbook struct {
	Bids []*Order
	Asks []*Order
}

func New() *Orderbook {
	return &Orderbook{nil, nil}
}

func (orderbook *Orderbook) matchMarket(order *Order) ([]*Trade, *Order) {
	return nil, nil
}

func min(a uint64, b uint64) uint64 {
	if b <= a {
		return b
	}
	return a
}

func makeTrade(bid *Order, ask *Order, main Side) *Trade {
	var price uint64
	if main == SideAsk {
		price = bid.Price
	} else {
		price = ask.Price
	}

	val := min(bid.Volume, ask.Volume)

	bid.Volume -= val
	ask.Volume -= val

	return &Trade{bid, ask, val, price}
}

func processLeftVolume(list *[]*Order, order *Order) (reject *Order) {
	if order.Kind == KindMarket {
		return order
	}
	*list = append(*list, order)
	if order.Side == SideAsk {
		sort.Slice(*list, func(i, j int) bool {
			return (*list)[i].Price < (*list)[j].Price
		})
	} else if order.Side == SideBid {
		sort.Slice(*list, func(i, j int) bool {
			return (*list)[i].Price > (*list)[j].Price
		})
	}
	return nil
}

func lists(orderbook *Orderbook, order *Order) (mySideList *[]*Order, targetSideList *[]*Order) {
	switch order.Side {
	case SideBid:
		mySideList = &orderbook.Bids
		targetSideList = &orderbook.Asks
	case SideAsk:
		mySideList = &orderbook.Asks
		targetSideList = &orderbook.Bids
	}
	return
}

func sides(order *Order, target *[]*Order, i int) (bid *Order, ask *Order) {
	if order.Side == SideBid {
		bid = order
		ask = (*target)[i]
	} else {
		bid = (*target)[i]
		ask = order
	}
	return
}
func (orderbook *Orderbook) match(order *Order) (trades []*Trade, reject *Order) {
	mySideList, targetSideList := lists(orderbook, order)

	for i := 0; i < len(*targetSideList); i++ {
		bid, ask := sides(order, targetSideList, i)

		if (order.Kind == KindMarket) || (ask.Price <= bid.Price) {
			if trade := makeTrade(bid, ask, order.Side); trade != nil {
				trades = append(trades, trade)
				if (*targetSideList)[i].Volume == 0 {
					(*targetSideList) = (*targetSideList)[:i+copy((*targetSideList)[i:], (*targetSideList)[i+1:])]
					i--
				}
				if order.Volume == 0 {
					return
				}
			}
		}
	}
	reject = processLeftVolume(mySideList, order)
	return trades, reject
}

func (orderbook *Orderbook) Match(order *Order) ([]*Trade, *Order) {
	return orderbook.match(order)
}
