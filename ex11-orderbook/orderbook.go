package orderbook

import (
	"fmt"
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
	return &Orderbook{nil, nil}
}

func (orderbook *Orderbook) matchMarket(order *Order) ([]*Trade, *Order) {
	return nil, nil
}

func (order *Order) String() string {
	return strconv.Itoa(int(order.Price))
}

func (trade *Trade) String() string {
	return strconv.Itoa(int(trade.Price))
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

func (orderbook *Orderbook) match(order *Order) (trades []*Trade) {
	fmt.Println("FIRST:", trades, order.Side, order.Price)

	var mySideList *[]*Order
	var targetSideList *[]*Order
	var accept func(*Order, *Order) bool

	switch order.Side {
	case SideBid:
		mySideList = &orderbook.Bids
		targetSideList = &orderbook.Asks
		accept = func(bid *Order, ask *Order) bool {
			return (order.Kind == KindMarket) || (bid.Price >= ask.Price)
		}
	case SideAsk:
		mySideList = &orderbook.Asks
		targetSideList = &orderbook.Bids
		accept = func(ask *Order, bid *Order) bool {
			return (order.Kind == KindMarket) || (ask.Price <= bid.Price)
		}
	}

	fmt.Printf("Target side: %v, nil:%v\n", targetSideList, targetSideList == nil)
	for _, restOrder := range *targetSideList {
		if accept(order, restOrder) {
			var bid *Order
			var ask *Order
			if order.Side == SideBid {
				bid = order
				ask = restOrder
			} else {
				bid = restOrder
				ask = order
			}

			if trade := makeTrade(bid, ask, order.Side); trade != nil {
				trades = append(trades, trade)
				if order.Volume == 0 {
					return
				}
			}
			fmt.Println("trades:", trades)
		}
	}
	*mySideList = append(*mySideList, order)
	if order.Side == SideAsk {
		sort.Slice(*mySideList, func(i, j int) bool {
			return (*mySideList)[i].Price < (*mySideList)[j].Price
		})
	} else if order.Side == SideBid {
		sort.Slice(*mySideList, func(i, j int) bool {
			return (*mySideList)[i].Price > (*mySideList)[j].Price
		})
	}

	return trades
}

func (orderbook *Orderbook) Match(order *Order) ([]*Trade, *Order) {
	return orderbook.match(order), nil
}
