package usecase

import (
	"github.com/FabioWGermano/clean-architecture/internal/entity"
	"github.com/FabioWGermano/clean-architecture/pkg/events"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderList       events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderList events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
		OrderList:       OrderList,
		EventDispatcher: EventDispatcher,
	}
}

func (c *ListOrdersUseCase) Execute() ([]*OrderOutputDTO, error) {
	orders, err := c.OrderRepository.ListOrders()
	if err != nil {
		return []*OrderOutputDTO{}, err
	}

	var orderDTOs []*OrderOutputDTO
	for _, order := range orders {
		orderDTO := &OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		}
		orderDTOs = append(orderDTOs, orderDTO)
	}

	c.OrderList.SetPayload(orderDTOs)
	c.EventDispatcher.Dispatch(c.OrderList)

	return orderDTOs, nil
}
