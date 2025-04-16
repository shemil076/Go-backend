package services

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/shemil076/loyalty-backend/internal/squareclient"
	"github.com/square/square-go-sdk"
	"github.com/square/square-go-sdk/loyalty"
)


func CreateLoyalAccount(phoneNumber string) (string, error){
	programID := os.Getenv("PROGRAMID")
	if programID == "" {
		return "", errors.New("PROGRAMID environment variable is not set")
	}

	// resp ,err := squareclient.NewClient.Loyalty.Accounts.Create(
	// 	context.TODO(),
	// 	&loyalty.CreateLoyaltyAccountRequest{
	// 		IdempotencyKey: uuid.NewString(),
	// 		LoyaltyAccount: &square.LoyaltyAccount{
	// 			ProgramID: programID,
	// 			Mapping: &square.LoyaltyAccountMapping{
	// 				PhoneNumber: square.String(phoneNumber),
	// 			},
	// 		},
	// 	},
	// )
	resp, err := squareclient.NewClient.Loyalty.Accounts.Create(
		context.TODO(),
		&loyalty.CreateLoyaltyAccountRequest{
			IdempotencyKey: uuid.NewString(),
			LoyaltyAccount: &square.LoyaltyAccount{
				ProgramID: programID,
				Mapping: &square.LoyaltyAccountMapping{
					PhoneNumber: square.String(phoneNumber),
				},
			},
		},
	)
	
	if err != nil {
		fmt.Printf("CreateLoyaltyAccount error: %v\n", err)
		return "", fmt.Errorf("square API error: %w", err)
	}
	
	fmt.Printf("resp: %+v\n", resp)
	

	fmt.Printf("resp is %v", resp)
	if err != nil{
		return "", fmt.Errorf("square API error: %w", err)
	}

	loyaltyAccount := resp.GetLoyaltyAccount()

	if loyaltyAccount == nil {
		return "", errors.New("received nil loyalty account from Square API")
	}

	return  *loyaltyAccount.ID, nil
}