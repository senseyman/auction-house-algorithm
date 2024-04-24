package report

import (
	"fmt"

	"github.com/senseyman/auction-house/model"
)

const template = "%d|%s|%s|%s|%.2f|%d|%.2f|%.2f"

// Service reports auction results to stdout console.
type Service struct {
}

func New() *Service {
	return &Service{}
}

// Report prints results to stdout console by the template.
func (s *Service) Report(fos []model.ActionResult) error {
	for _, el := range fos {
		res := fmt.Sprintf(template,
			el.CloseTime,
			el.Item,
			digitOrEmpty(el.UserID),
			el.Status,
			el.PricePaid,
			el.Statistics.TotalBidCount,
			el.Statistics.HighestBid,
			el.Statistics.LowestBid,
		)

		fmt.Println(res)
	}

	return nil
}

func digitOrEmpty(el int) string {
	if el == 0 {
		return ""
	}
	return fmt.Sprintf("%d", el)
}
