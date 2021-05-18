package model

import (
	"github.com/google/uuid"
)

type Client struct {
	ClientID     uuid.UUID
	ClientName   string
	ClientSecret string
}

func (c *Client) GetClientSecret() {


}

func (c *Client) Validate() bool {


	return true
}