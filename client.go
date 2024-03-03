package rinha2024q1crebito

import "fmt"

type ClientID struct {
	id int
}

func NewClientID(id int) (ClientID, error) {
	c := ClientID{id: id}
	if err := c.validate(); err != nil {
		return ClientID{}, err
	}
	return c, nil
}

func (c *ClientID) validate() error {
	if c.id < 1 {
		return fmt.Errorf(
			"id do cliente deve ser um número inteiro positivo: %w",
			ErrInvalidParameter)
	}
	return nil
}

func (c *ClientID) Value() int {
	return c.id
}
