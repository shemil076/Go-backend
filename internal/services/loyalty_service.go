package services

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/shemil076/loyalty-backend/internal/models"
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

func EarnPoints(loyaltyId string, totalAmountInCent int64) (int, error) {
	programID, locationID := os.Getenv("PROGRAMID"), os.Getenv("DEFAULT_LOCATIONID")
	if programID == "" || locationID == "" {
		return 0,fmt.Errorf("environment variable are not set")
	}

	resp, err := calculatePoints(locationID, totalAmountInCent, programID)

	if err != nil {
		return 0,fmt.Errorf("error occurred while calculating points %w", err)
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
	if err != nil {
		return 0,fmt.Errorf("error occurred while accumulating points %w", err)
	}
	return resp, nil
}


func createLoyaltyReward(loyaltyId string, rewardTierId string)(string, error){

	fmt.Printf("\n reward id is, %v",loyaltyId)
	fmt.Printf("\n reward id is, %v",rewardTierId)
	resp, err := squareclient.NewClient.Loyalty.Rewards.Create(
		context.TODO(),
		&loyalty.CreateLoyaltyRewardRequest{
			IdempotencyKey: uuid.NewString(),
			Reward: &square.LoyaltyReward{
				LoyaltyAccountID: loyaltyId,
				RewardTierID: rewardTierId,
			},
		},
	)

	if err != nil {
		fmt.Print("\n Error occured when creating theloyalty reward")
		return "", fmt.Errorf("error occurred while creating reward %w", err)
	}

	return *resp.Reward.ID ,nil
}


func RedeemReward(loyaltyId string, totalAmountInCent int64) (float64, error) {
	locationID, rewardTierId :=  os.Getenv("DEFAULT_LOCATIONID"), os.Getenv("REWARDTIERID")
	if rewardTierId == "" || locationID == "" {
		return 0, fmt.Errorf("environment variable are not set")
	}

	rewardID, err := createLoyaltyReward(loyaltyId, rewardTierId);

	if err != nil {
		fmt.Print("reward creation failed")
		return 0, fmt.Errorf("reward creation failed: %w",err)
	}


	_, err =  squareclient.NewClient.Loyalty.Rewards.Redeem(
		context.TODO(),
		&loyalty.RedeemLoyaltyRewardRequest{
			IdempotencyKey: uuid.NewString(),
			RewardID: rewardID,
			LocationID: locationID,
		},
	)
	
	
	if (err != nil){
		return 0, fmt.Errorf("redemption failed: %w", err)
	}
	

	amountInUSD := float64(totalAmountInCent) / 100
	finalAmount := amountInUSD - (amountInUSD * 0.10)
	if finalAmount < 0 {
		finalAmount = 0 
	}

	resp, _ := retrieveReward(rewardID)
	
	if resp != "REDEEMED" {
		return 0, fmt.Errorf("reward not properly redeemed")
	}
	return finalAmount, nil

}

func retrieveReward(rewardId string) (string,error){
	resp, err :=  squareclient.NewClient.Loyalty.Rewards.Get(
		context.TODO(),
		&loyalty.GetRewardsRequest{
			RewardID: rewardId,
		},
	)

	if err != nil {
		return "", fmt.Errorf("reward retrieving failed %w", err)
	}

	return string(*resp.Reward.Status), nil
}

func GetBalance(loyaltyId string)(int, error){
	resp, err := squareclient.NewClient.Loyalty.Accounts.Get(
		context.TODO(),
		&loyalty.GetAccountsRequest{
			AccountID: loyaltyId,
		},
	)

	if err != nil {
		return 0, fmt.Errorf("account retrieving failed %w", err)
	}

	return *resp.LoyaltyAccount.Balance, nil
}


func getEvents(loyaltyId string)([]*square.LoyaltyEvent,error){
	resp, err := squareclient.NewClient.Loyalty.SearchEvents(
		context.TODO(),
		&square.SearchLoyaltyEventsRequest{
			Limit: square.Int(30),
			Query: &square.LoyaltyEventQuery{
				Filter: &square.LoyaltyEventFilter{
					LoyaltyAccountFilter: &square.LoyaltyEventLoyaltyAccountFilter{
						LoyaltyAccountID: loyaltyId,
					},
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("events retrieving failed %w", err)
	}

	return resp.Events, nil
}

func formatEvents(events []*square.LoyaltyEvent) []models.Transaction {
	var transactions []models.Transaction

	for _, event := range events {
		t := models.Transaction {
			Date: event.CreatedAt,
		}

		switch event.Type {
		case "ACCUMULATE_POINTS":
			t.Title = "ACCUMULATE_POINTS"
			t.Description = "Points are added to a loyalty account for a purchase that qualified for points based on an accrual rule."
			t.Points = int64(*event.AccumulatePoints.Points)

		case "CREATE_REWARD":
			t.Title = "CREATE_REWARD"
			t.Description = "A loyalty reward is created."
			t.Points = int64(event.CreateReward.Points)

		case "REDEEM_REWARD":
			t.Title = "REDEEM_REWARD"
			t.Description = "A loyalty reward is redeemed."
			// t.Points = int64(event.re)
		
		case "DELETE_REWARD":
			t.Title = "DELETE_REWARD"
			t.Description = "A loyalty reward is deleted."
			t.Points = int64(event.DeleteReward.Points)

		case "ADJUST_POINTS":
			t.Title = "ADJUST_POINTS"
			t.Description = "Loyalty points are manually adjusted."
			// t.Points = int64(event.)

		case "EXPIRE_POINTS":
			t.Title = "EXPIRE_POINTS"
			t.Description = "Loyalty points are expired according to the expiration policy of the loyalty program."
			t.Points = int64(event.ExpirePoints.Points)

		case "OTHER":
			t.Title = "OTHER"
			t.Description = "Some other loyalty event occurred."
			t.Points = int64(event.OtherEvent.Points)

		case "ACCUMULATE_PROMOTION_POINTS":
			t.Title = "ACCUMULATE_POINTS"
			t.Description = "Points are added to a loyalty account for a purchase that qualified for a loyalty promotion."
			t.Points = int64(event.AccumulatePromotionPoints.Points)
		}
		transactions = append(transactions, t)
	}

	return transactions
}


func ReturnTransactions(loyaltyId string) ([]models.Transaction, error) {
	events, err := getEvents(loyaltyId)
	
	if err != nil {
		return nil, fmt.Errorf("error occurred while fetching the events %w", err)
	}

	transactions := formatEvents(events)

	fmt.Printf("transactions => %v",transactions)

	if transactions == nil {
		return nil, fmt.Errorf("error occurred while formatting the events")
	}

	return transactions, nil
}
