package orderbook

import "fmt"

type RestOrder struct {
	order *Order
	next  *RestOrder
}

type Orderbook struct {
	Bids *RestOrder
	Asks *RestOrder
}

func New() *Orderbook {
	return &Orderbook{}
}

func (order *RestOrder) append(new *RestOrder) chan *RestOrder {

}

func (order *RestOrder) iterator() chan *RestOrder {
	res := make(chan *RestOrder)
	go func() {
		node := order
		for {
			if node == nil {
				close(res)
				return
			}
			res <- node
			node = node.next
		}
	}()
	return res
}

func (orderbook *Orderbook) matchMarket(order *Order) ([]*Trade, *Order) {
	return nil, nil
}

func (orderbook *Orderbook) matchLimit(order *Order) []*Trade {
	list := orderbook.Bids.iterator()
	for v := range list {
		fmt.Println("Val:", v)
	}
	fmt.Println("Finish")
	return nil
}

func (orderbook *Orderbook) Match(order *Order) ([]*Trade, *Order) {
	switch order.Kind {
	case KindMarket:
		return orderbook.matchMarket(order)
	case KindLimit:
		return orderbook.matchLimit(order), nil
	default:
		return nil, nil
	}
}
