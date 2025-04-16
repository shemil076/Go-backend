package squareclient

import (
	"github.com/square/square-go-sdk"
	client "github.com/square/square-go-sdk/client"
	"github.com/square/square-go-sdk/option"
)




var NewClient *client.Client

func init(){
	// accessToken := os.Getenv("SQUARE_ACCESS_TOKEN")
	NewClient = client.NewClient(
		option.WithBaseURL(
            square.Environments.Sandbox,
        ),
		option.WithToken("EAAAl60M4ryddhhr1ymbq1I1eRGmhIyhNA1G2HgythnfHyflWcIZXGQjWszwKEak"),
	)
}