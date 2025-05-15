package usecase

import (
	"fmt"

	"github.com/devfullcycle/20-CleanArch/internal/entity"
	"github.com/devfullcycle/20-CleanArch/internal/event"
	"github.com/devfullcycle/20-CleanArch/pkg/events"
)

type ListOrdersDTO struct {
	OrderListDTO []OrderOutputDTO
}

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
		EventDispatcher: EventDispatcher,
	}
}

func (c *ListOrderUseCase) Execute() (ListOrdersDTO, error) {

	orders, err := c.OrderRepository.FindAll()
	if err != nil {
		return ListOrdersDTO{}, err
	}

	dto := ListOrdersDTO{
		OrderListDTO: make([]OrderOutputDTO, len(orders)),
	}

	for i, order := range orders {
		dto.OrderListDTO[i] = OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		}
	}

	var orderListedEvent = event.NewOrderListed()
	orderListedEvent.SetPayload(fmt.Sprintf("Total Orders listed: %v", len(dto.OrderListDTO)))
	c.EventDispatcher.Dispatch(orderListedEvent)

	return dto, nil
}
