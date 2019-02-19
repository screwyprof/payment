package gin

import (
	"github.com/screwyprof/payment/internal/pkg/delivery/gin/request"
	"github.com/screwyprof/payment/internal/pkg/delivery/gin/response"
)

type GopherPayTestClient struct {
	client *HttpTestClient
}

func (c *GopherPayTestClient) ShowAccount(number string) (response.AccountInfo, error) {
	resp := response.AccountInfo{}
	err := c.client.SendGetRequest("/api/v1/accounts/"+number, &resp)

	return resp, err
}

func (c *GopherPayTestClient) OpenAccount(r request.OpenAccount) (response.ShortAccountInfo, error) {
	resp := response.ShortAccountInfo{}
	err := c.client.SendPostRequest("/api/v1/accounts", r, &resp)

	return resp, err
}

func (c *GopherPayTestClient) TransferMoney(r request.TransferMoney) (response.Message, error) {
	resp := response.Message{}
	err := c.client.SendPostRequest("/api/v1/accounts/"+r.From+"/transfer", r, &resp)

	return resp, err
}
