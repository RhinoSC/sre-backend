package donation_helper

import "github.com/RhinoSC/sre-backend/internal"

func ConvertDonationWithBidDetailsDBtoInternal(donationDB internal.DonationWithBidDetailsDB) (donation internal.DonationWithBidDetails) {
	donation = internal.DonationWithBidDetails{
		Donation: internal.Donation{
			ID:          donationDB.ID,
			Name:        donationDB.Name,
			Email:       donationDB.Email,
			TimeMili:    donationDB.TimeMili,
			Amount:      donationDB.Amount,
			Description: donationDB.Description,
			ToBid:       donationDB.ToBid,
			EventID:     donationDB.EventID,
		},
		BidID:            donationDB.BidID.String,
		Bidname:          donationDB.Bidname.String,
		Goal:             donationDB.Goal.Float64,
		CurrentAmount:    donationDB.CurrentAmount.Float64,
		BidDescription:   donationDB.BidDescription.String,
		Type:             internal.BidType(donationDB.Type.String),
		CreateNewOptions: donationDB.CreateNewOptions.Bool,
		RunID:            donationDB.RunID.String,
		OptionID:         donationDB.OptionID.String,
		OptionName:       donationDB.OptionName.String,
		OptionAmount:     donationDB.OptionAmount.Float64,
	}

	return
}
