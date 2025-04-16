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
	

	loyaltyAccount := resp.GetLoyaltyAccount()

	if loyaltyAccount == nil {
		return "", errors.New("received nil loyalty account from Square API")
	}

	return  *loyaltyAccount.ID, nil
}

func calculatePoints(loyaltyId string, totalAmountInCent int64, programID string) (int, error){
	resp, err := squareclient.NewClient.Loyalty.Programs.Calculate(
		context.TODO(),
		&loyalty.CalculateLoyaltyPointsRequest{
			ProgramID: programID,
			TransactionAmountMoney: &square.Money{
				Currency: square.CurrencyUsd.Ptr(),
				Amount: square.Int64(
					totalAmountInCent,
				),
			},
			LoyaltyAccountID: square.String(
				loyaltyId,
			),
		},
	)

	if err != nil {
		return 0, fmt.Errorf("error occurred while counting the points %w", err)
	} 

	return *resp.Points, nil
}

func EarnPoints(loyaltyId string, totalAmountInCent int64) error {
	programID, locationID := os.Getenv("PROGRAMID"), os.Getenv("DEFAULT_LOCATIONID")
	if programID == "" || locationID == "" {
		return fmt.Errorf("environment variable are not set")
	}

	resp, err := calculatePoints(locationID, totalAmountInCent, programID)

	if err != nil {
		return fmt.Errorf("error occurred while calculating points %w", err)
	}

	_, err = squareclient.NewClient.Loyalty.Accounts.AccumulatePoints(
		context.TODO(),
		&loyalty.AccumulateLoyaltyPointsRequest{
			AccountID: loyaltyId,
			IdempotencyKey: uuid.NewString(),
			LocationID: locationID,
			AccumulatePoints: &square.LoyaltyEventAccumulatePoints{
				Points: square.Int(
					resp,
				),
			},
		},
	)

	return err
}