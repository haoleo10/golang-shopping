package order

import (
	"shopping/domain/cart"
	"shopping/domain/product"
	"shopping/utils/pagination"
	"time"
)

var day14ToHours float64 = 336

type Service struct {
	orderRepository       Repository
	orderedItemRepository OrderedItemRepository
	productRepository     product.Repository
	cartRepository        cart.Repository
	cartItemRepository    cart.ItemRepository
}

// 实例化
func NewService(
	orderRepository Repository,
	orderedItemRepository OrderedItemRepository,
	productRepository product.Repository,
	cartRepository cart.Repository,
	cartItemRepository cart.ItemRepository,
) *Service {
	orderRepository.Migration()
	orderedItemRepository.Migration()
	return &Service{
		orderRepository:       orderRepository,
		orderedItemRepository: orderedItemRepository,
		productRepository:     productRepository,
		cartRepository:        cartRepository,
		cartItemRepository:    cartItemRepository,
	}

}

// 完成订单
func (c *Service) CompleteOrder(userId uint) error {
	//拿到当前的购物车
	currentCart, err := c.cartRepository.FindOrCreateByUserID(userId)
	if err != nil {
		return err
	}
	//从购物车拿到当前的item
	cartItems, err := c.cartItemRepository.GetItems(currentCart.UserID)
	if err != nil {
		return err
	}
	if len(cartItems) == 0 {
		return ErrEmptyCartFound
	}
	//通过购物车item创建一个订单的item
	orderedItems := make([]OrderedItem, 0)
	for _, item := range cartItems {
		orderedItems = append(orderedItems, *NewOrderedItem(item.Count, item.ProductID))
	}
	//新的order
	err = c.orderRepository.Create(NewOrder(userId, orderedItems))
	return err
}

// 取消订单，就是把订单里面的IsCanceled变为true
func (c *Service) CancelOrder(uid, oid uint) error {
	currentOrder, err := c.orderRepository.FindByOrderID(oid)
	if err != nil {
		return err
	}
	if currentOrder.UserID != uid {
		return ErrInvalidOrderID
	}
	if currentOrder.CreatedAt.Sub(time.Now()).Hours() > day14ToHours {
		return ErrCancelDurationPassed
	}
	currentOrder.IsCanceled = true
	err = c.orderRepository.Update(*currentOrder)

	return err
}

// 获得订单，分页查询所有订单
func (c *Service) GetAll(page *pagination.Pages, uid uint) *pagination.Pages {
	orders, count := c.orderRepository.GetAll(page.Page, page.PageSize, uid)
	page.Items = orders
	page.TotalCount = count
	return page
}
