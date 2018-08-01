package orderbook

import (
	"sort"
	"strconv"
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
	bids := make([]*Order, 0, 100)
	asks := make([]*Order, 0, 100)
	return &Orderbook{bids, asks}
	//return &Orderbook{nil, nil}
}

func (orderbook *Orderbook) matchMarket(order *Order) ([]*Trade, *Order) {
	return nil, nil
}

func (order *Order) String() string {
	return strconv.Itoa(int(order.Price)) + "vol:" + strconv.Itoa(int(order.Volume))
}

func (trade *Trade) String() string {
	return "Trade: " + strconv.Itoa(int(trade.Price)) + " vol: " + strconv.Itoa(int(trade.Volume))
}

func makeTrade(bid *Order, ask *Order, main Side) *Trade {
	var trade *Trade = nil

	var price uint64
	if main == SideAsk {
		price = bid.Price
	} else {
		price = ask.Price
	}

	if bid.Volume <= ask.Volume {
		trade = &Trade{bid, ask, bid.Volume, price}
		ask.Volume -= bid.Volume
		bid.Volume = 0
	} else {
		trade = &Trade{bid, ask, ask.Volume, price}
		bid.Volume -= ask.Volume
		ask.Volume = 0
	}

	return trade
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

func (orderbook *Orderbook) match(order *Order) (trades []*Trade, reject *Order) {
	var mySideList *[]*Order
	var targetSideList *[]*Order
	//var accept func(*Order, *Order) bool

	switch order.Side {
	case SideBid:
		mySideList = &orderbook.Bids
		targetSideList = &orderbook.Asks
		/*accept = func(bid *Order, ask *Order) bool {
			return (order.Kind == KindMarket) || (bid.Price >= ask.Price)
		}*/
	case SideAsk:
		mySideList = &orderbook.Asks
		targetSideList = &orderbook.Bids
		/*accept = func(ask *Order, bid *Order) bool {
			return (order.Kind == KindMarket) || (ask.Price <= bid.Price)
		}*/
	}

	for i := 0; i < len(*targetSideList); i++ {
		var bid *Order
		var ask *Order
		if order.Side == SideBid {
			bid = order
			ask = (*targetSideList)[i]
		} else {
			bid = (*targetSideList)[i]
			ask = order
		}

		//		if accept(order, (*targetSideList)[i]) {
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
