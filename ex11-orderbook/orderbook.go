package orderbook

import (
	"sort"
)

type RestOrder struct {
	order *Order
	next  *RestOrder
}

type Orderbook struct {
	Bids []*Order
	Asks []*Order
}

func New() *Orderbook {
	return &Orderbook{nil, nil}
}

func min(a uint64, b uint64) uint64 {
	if b <= a {
		return b
	}
	return a
}

func makeTrade(bid *Order, ask *Order, initiator Side) *Trade {
	var price uint64

	switch initiator{
	case SideAsk:
		price = bid.Price
	case SideBid:
		price = ask.Price
	}

	amount := min(bid.Volume, ask.Volume)
	bid.Volume -= amount
	ask.Volume -= amount

	return &Trade{bid, ask, amount, price}
}

func sortRestOrders(list *[]*Order, initiator Side){
	switch initiator{ //sort in ascending order
	case SideAsk:
		sort.Slice(*list, func(i, j int) bool {return (*list)[i].Price < (*list)[j].Price})
	case SideBid://sort in descnding order
		sort.Slice(*list, func(i, j int) bool {return (*list)[i].Price > (*list)[j].Price})
	}
}

func processLeftVolume(myRestOrders *[]*Order, initiator *Order) (reject *Order) {
	if initiator.Kind == KindMarket {
		return initiator
	}

	newOrder := *initiator
	*myRestOrders = append(*myRestOrders, &newOrder)

	sortRestOrders(myRestOrders, initiator.Side)
	return nil
}

func restOrders(orderbook *Orderbook, order *Order) (myRestOrders *[]*Order, targetRestOrders *[]*Order) {
	switch order.Side {
	case SideBid:
		myRestOrders = &orderbook.Bids
		targetRestOrders = &orderbook.Asks
	case SideAsk:
		myRestOrders = &orderbook.Asks
		targetRestOrders = &orderbook.Bids
	}
	return
}

func sides(initiator *Order, target *Order) (bid *Order, ask *Order) {
	switch initiator.Side {
	case SideBid:
		bid = initiator
		ask = target
	case SideAsk:
		bid = target
		ask = initiator
	}
	return
}

func delRestOrdIfFulfilled(target *[]*Order,  i *int){
	if (*target)[*i].Volume == 0 {
		(*target) = (*target)[:*i+copy((*target)[*i:], (*target)[*i+1:])]
		(*i)--
	}
}

func (orderbook *Orderbook) Match(order *Order) (trades []*Trade, reject *Order) {
	myRestOrders, targetRestOrders := restOrders(orderbook, order)

	for i := 0; i < len(*targetRestOrders); i++ {
		bid, ask := sides(order, (*targetRestOrders)[i])

		if (order.Kind == KindMarket) || (ask.Price <= bid.Price) { // price is suitable
			trade := makeTrade(bid, ask, order.Side);
			trades = append(trades, trade)

			delRestOrdIfFulfilled(targetRestOrders, &i)

			if order.Volume == 0 { //order fulfilled
				return
			}
		}else{ //there is no suitable price among resting orders
			break
		}
	}
	reject = processLeftVolume(myRestOrders, order) //reject or add to resting orders
	return
}
