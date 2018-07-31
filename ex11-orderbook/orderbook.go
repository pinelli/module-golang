package orderbook

import "fmt"

type RestOrder struct {
	order *Order
	next  *RestOrder
}

type RestOrderList struct{
	head *RestOrder
}

type Orderbook struct {
	Bids *RestOrder
	Asks *RestOrder
}

func New() *Orderbook {
	return &Orderbook{}
}

func (order *RestOrder) append(newOrd *RestOrder) *RestOrder {
	if order == nil{
		order = newOrd
		return order
	}
	next := order.next
	order.next = newOrd
	newOrd.next = next
	return order
}

func (order *RestOrder) iterator() (chan *RestOrder, chan bool) {
	res := make(chan *RestOrder)
	stop := make(chan bool, 1)

	go func() {
		node := order
		for {
			if node == nil {
				close(res)
				return
			}
			select{
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

func (orderbook *Orderbook) matchLimit(order *Order)  []*Trade{
	var trades []*Trade = nil
	fmt.Println("FIRST:", trades, order.Side, order.Price)

	var mySideList **RestOrder
	var targetSideList **RestOrder
	var accept func (*Order,*Order) bool

	switch order.Side{
	case SideBid:
		mySideList = &orderbook.Bids
		targetSideList = &orderbook.Asks
		accept = func (bid *Order, ask *Order) bool {
			return bid.Price >= ask.Price
		}
	case SideAsk:
		mySideList = &orderbook.Asks
		targetSideList = &orderbook.Bids
		accept = func (ask *Order, bid *Order) bool {
			return ask.Price <= bid.Price
		}
	}

	list, stopIter := (*targetSideList).iterator()
	for restOrder := range list {
			fmt.Println("OK:", trades, order.Side, order.Price, restOrder.order.Side,  restOrder.order.Price)
		if accept(order, restOrder.order){
			fmt.Println("ACCEPTED:", trades, order.Side, order.Price, restOrder.order.Side,  restOrder.order.Price)
		
			trades = append(trades, &Trade{Volume: 10000, Price:  60000000})
			fmt.Println("trades:", trades[0])
		}
	}
	stopIter<-true
	//orderbook.Bids = orderbook.Bids.append(&RestOrder{order, nil})
	*mySideList = (*mySideList).append(&RestOrder{order, nil})

	return trades
}

func (orderbook *Orderbook) Match(order *Order) ([]*Trade, *Order) {
	// var trades []*Trade = make([]*Trade, 1)
	// trades = append(trades, &Trade{Bid: order, Ask: order, Volume: 10000, Price:  60000000})
	// fmt.Println("trades:", trades)
	// return trades, nil

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
