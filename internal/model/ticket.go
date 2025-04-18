package model

/**
 * @File: ticket.go
 * @Description:
 *
 * @Author: Timmy
 * @Create: 2025/4/18 下午10:07
 * @Software: GoLand
 * @Version:  1.0
 */

type TicketRequest struct {
	EventID string `json:"eventId"`
	UserID  string `json:"userId"`
	Zone    string `json:"zone"`
}
