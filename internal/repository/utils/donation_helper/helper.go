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
		BidDetails: &internal.DonationBidDetails{
			BidID:            donationDB.BidDetails.BidID.String,
			Bidname:          donationDB.BidDetails.Bidname.String,
			Goal:             donationDB.BidDetails.Goal.Float64,
			CurrentAmount:    donationDB.BidDetails.CurrentAmount.Float64,
			BidDescription:   donationDB.BidDetails.BidDescription.String,
			Type:             internal.BidType(donationDB.BidDetails.Type.String),
			CreateNewOptions: donationDB.BidDetails.CreateNewOptions.Bool,
			RunID:            donationDB.BidDetails.RunID.String,
			OptionID:         donationDB.BidDetails.OptionID.String,
			OptionName:       donationDB.BidDetails.OptionName.String,
			OptionAmount:     donationDB.BidDetails.OptionAmount.Float64,
		},
	}

	return
}
