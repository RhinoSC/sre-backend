package donation_helper

import (
	"github.com/RhinoSC/sre-backend/internal"
)

type BidDetails struct {
	BidID            string           `json:"bid_id,omitempty"`
	Bidname          string           `json:"bidname,omitempty"`
	Goal             float64          `json:"goal,omitempty"`
	CurrentAmount    float64          `json:"current_amount,omitempty"`
	BidDescription   string           `json:"bid_description,omitempty"`
	Type             internal.BidType `json:"type,omitempty"`
	CreateNewOptions bool             `json:"create_new_options,omitempty"`
	RunID            string           `json:"run_id,omitempty"`
	OptionID         string           `json:"option_id,omitempty"`
	OptionName       string           `json:"option_name,omitempty"`
	OptionAmount     float64          `json:"option_amount,omitempty"`
}

type BidDetailsAsBodyJSON struct {
	BidID            string `json:"bid_id,omitempty"`
	CreateNewOptions bool   `json:"create_new_options,omitempty"`
	OptionID         string `json:"option_id,omitempty"`
	OptionName       string `json:"option_name,omitempty"`
}

type DonationAsJSON struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	TimeMili    int64      `json:"time_mili"`
	Amount      float64    `json:"amount"`
	Description string     `json:"description"`
	ToBid       *bool      `json:"to_bid"`
	EventID     string     `json:"event_id"`
	BidDetails  BidDetails `json:"bid_details"`
}

type DonationAsBodyJSON struct {
	Name        string               `json:"name" validate:"required"`
	Email       string               `json:"email" validate:"required"`
	TimeMili    int64                `json:"time_mili" validate:"required"`
	Amount      float64              `json:"amount" validate:"required"`
	Description string               `json:"description"`
	ToBid       *bool                `json:"to_bid" validate:"required"`
	EventID     string               `json:"event_id" validate:"required"`
	BidDetails  BidDetailsAsBodyJSON `json:"bid_details,omitempty"`
}

func ConvertDonationToJSON(donation internal.Donation) (donationJSON DonationAsJSON) {
	donationJSON = DonationAsJSON{
		ID:          donation.ID,
		Name:        donation.Name,
		Email:       donation.Email,
		TimeMili:    donation.TimeMili,
		Amount:      donation.Amount,
		Description: donation.Description,
		ToBid:       &donation.ToBid,
		EventID:     donation.EventID,
	}
	return
}

func ConvertDonationsToJSON(donations []internal.Donation) []DonationAsJSON {
	result := make([]DonationAsJSON, len(donations))
	for i, donation := range donations {
		result[i] = ConvertDonationToJSON(donation)
	}
	return result
}

func ConvertDonationWithBidDetailsToJSON(donationWithBidDetails internal.DonationWithBidDetails) (donationJSON DonationAsJSON) {
	donationJSON = DonationAsJSON{
		ID:          donationWithBidDetails.Donation.ID,
		Name:        donationWithBidDetails.Donation.Name,
		Email:       donationWithBidDetails.Donation.Email,
		TimeMili:    donationWithBidDetails.Donation.TimeMili,
		Amount:      donationWithBidDetails.Donation.Amount,
		Description: donationWithBidDetails.Donation.Description,
		ToBid:       &donationWithBidDetails.Donation.ToBid,
		EventID:     donationWithBidDetails.Donation.EventID,
		BidDetails: BidDetails{
			BidID:            donationWithBidDetails.BidID,
			Bidname:          donationWithBidDetails.Bidname,
			Goal:             donationWithBidDetails.Goal,
			CurrentAmount:    donationWithBidDetails.CurrentAmount,
			BidDescription:   donationWithBidDetails.BidDescription,
			Type:             donationWithBidDetails.Type,
			CreateNewOptions: donationWithBidDetails.CreateNewOptions,
			RunID:            donationWithBidDetails.RunID,
			OptionID:         donationWithBidDetails.OptionID,
			OptionName:       donationWithBidDetails.OptionName,
			OptionAmount:     donationWithBidDetails.OptionAmount,
		},
	}
	return
}

func ConvertDonationsWithBidDetailsToJSON(donationsWithBidDetails []internal.DonationWithBidDetails) []DonationAsJSON {
	result := make([]DonationAsJSON, len(donationsWithBidDetails))
	for i, donation := range donationsWithBidDetails {
		result[i] = ConvertDonationWithBidDetailsToJSON(donation)
	}
	return result
}
