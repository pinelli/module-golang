package orderbook

import "fmt"

type RestOrder struct {
	order *Order
	next  *RestOrder
}

type RestOrderList struct {
	head *RestOrder
}

type Orderbook struct {
	Bids RestOrderList
	Asks RestOrderList
}

func New() *Orderbook {
	return &Orderbook{}
}

func (list *RestOrderList) appendFront(newOrd *RestOrder) {
	if list.head == nil {
		fmt.Println("YES")
		list.head = newOrd
		return
	}
	fmt.Println("NO")
	newOrd.next = list.head
	list.head = newOrd
	return
}

func (list *RestOrderList) pushSort(newOrd *RestOrder) {
	if list.head == nil {
		list.head = newOrd
		return
	}
	phony := &RestOrder{nil, list.head}
	nd := phony
	for nd.next != nil {
		if nd.next.order.Price > newOrd.order.Price {
			break
		}
		nd = nd.next
	}
	nd.append(newOrd)
	//push after nd

	if newOrd.next == list.head {
		list.head = newOrd
	}
	return
}

func (order *RestOrder) append(newOrd *RestOrder) *RestOrder {
	if order == nil {
		order = newOrd
		return order
	}
	next := order.next
	order.next = newOrd
	newOrd.next = next
	return order
}

func (list *RestOrderList) iterator() (chan *RestOrder, chan bool) {
	var order *RestOrder = list.head

	res := make(chan *RestOrder)
	stop := make(chan bool, 1)

	go func() {
		node := order
		for {
			if node == nil {
				close(res)
				return
			}
			select {
			case res <- node:
			case <-stop:
				fmt.Println("list killed")
				close(res)
				return
			}
			node = node.next
		}
	}()
	return res, stop
}

func (orderbook *Orderbook) matchMarket(order *Order) ([]*Trade, *Order) {
	return nil, nil
}

func putOrderToRest(list **RestOrder, order *RestOrder) {

}

func makeTrade(bid *Order, ask *Order, main Side) *Trade {
	var trade *Trade = nil

	if bid.Volume <= ask.Volume {
		trade = &Trade{bid, ask, bid.Volume, ask.Price}
		ask.Volume -= bid.Volume
		bid.Volume = 0
	} else {
		trade = &Trade{bid, ask, ask.Volume, ask.Price}
		bid.Volume -= ask.Volume
		ask.Volume = 0
	}

	return trade
}

func (orderbook *Orderbook) matchLimit(order *Order) []*Trade {
	var trades []*Trade = nil
	fmt.Println("FIRST:", trades, order.Side, order.Price)

	var mySideList *RestOrderList
	var targetSideList *RestOrderList
	var accept func(*Order, *Order) bool

	switch order.Side {
	case SideBid:
		mySideList = &orderbook.Bids
		targetSideList = &orderbook.Asks
		accept = func(bid *Order, ask *Order) bool {
			return bid.Price >= ask.Price
		}
	case SideAsk:
		mySideList = &orderbook.Asks
		targetSideList = &orderbook.Bids
		accept = func(ask *Order, bid *Order) bool {
			return ask.Price <= bid.Price
		}
	}

	list, stopIter := targetSideList.iterator()

	for restOrder := range list {
		if accept(order, restOrder.order) {
			var bid *Order
			var ask *Order
			if order.Side == SideBid {
				bid = order
				ask = restOrder.order
			} else {
				bid = restOrder.order
				ask = order
			}

			if trade := makeTrade(bid, ask, order.Side); trade != nil {
				trades = append(trades, trade)
			}
			fmt.Println("trades:", trades)
		}
	}
	stopIter <- true

	mySideList.pushSort(&RestOrder{order, nil})
	//mySideList.appendFront(&RestOrder{order, nil})

	return trades
}

func (orderbook *Orderbook) Match(order *Order) ([]*Trade, *Order) {
	return orderbook.matchLimit(order), nil

	// switch order.Kind {
	// case KindMarket:
	// 	return orderbook.matchMarket(order)
	// case KindLimit:
	// 	return orderbook.matchLimit(order), nil
	// default:
	// 	return nil, nil
	// }
}
