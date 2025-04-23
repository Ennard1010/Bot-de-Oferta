package mavismktgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/libs/defigo"
	"go/libs/utils"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

type WalletData struct {
	PrivateKey string
	PublicKey  string
	AuthToken  string
	Client     defigo.Web3Client
}
type CollectionValues struct {
	CollectionName     string
	CollectionDiscount float64
	CollectionAddress  string
	CollectionCriteria []map[string]interface{}
	FloorPrice         float64
	Wallet             *WalletData
}

var (
	Wallet1 = WalletData{
		PrivateKey: "Key",
		PublicKey:  "Key",
		Client:     *defigo.NewClient("Key", "https://api.roninchain.com/rpc"),
		AuthToken:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjFlZjM0ZTczLTQxMGMtNjAwNi1iMTg2LTJjYTk2ZGQyM2EwZSIsInNpZCI6MTc5NzI1ODU1LCJyb2xlcyI6WyJ1c2VyIl0sInNjcCI6WyJhbGwiXSwiYWN0aXZhdGVkIjp0cnVlLCJhY3QiOnRydWUsInJvbmluQWRkcmVzcyI6IjB4MzZiYTUwNTcwNmNkZmFkYzkzNGQ2OGEwNTJiYzRlODA5ZDkyYjBjOSIsImV4cCI6MTcyNjQzODE3MywiaWF0IjoxNzI1MjI4NTczLCJpc3MiOiJBeGllSW5maW5pdHkiLCJzdWIiOiIxZWYzNGU3My00MTBjLTYwMDYtYjE4Ni0yY2E5NmRkMjNhMGUifQ.Zkgcmr7m7_FMMngAD5YCkKhSkcQQ6m621sY_a_3cO8A",
	}
	Wallet2 = WalletData{
		PrivateKey: "Key",
		PublicKey:  "Key",
		Client:     *defigo.NewClient("Key", "https://api.roninchain.com/rpc"),
		AuthToken:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjFlZjY2MGQyLWIzYzYtNmQ3Zi1hMDIyLTlkMWIyN2E4MjBiNiIsInNpZCI6MTgwNjgxNzUyLCJyb2xlcyI6WyJ1c2VyIl0sInNjcCI6WyJhbGwiXSwiYWN0aXZhdGVkIjp0cnVlLCJhY3QiOnRydWUsInJvbmluQWRkcmVzcyI6IjB4OGNlNmI2NTc1ZTRjNTY5YTM5M2Q1YzQwNDM2ZDk5MmEzOTk0OTRiNCIsImV4cCI6MTcyNzEyMDQ5MSwiaWF0IjoxNzI1OTEwODkxLCJpc3MiOiJBeGllSW5maW5pdHkiLCJzdWIiOiIxZWY2NjBkMi1iM2M2LTZkN2YtYTAyMi05ZDFiMjdhODIwYjYifQ._S5ONi5-_SHK0eUQY1w4QmZyBlb2VINA5mb-kg6a-r4",
	}
	Wallet3 = WalletData{
		PrivateKey: "Key",
		PublicKey:  "Key",
		Client:     *defigo.NewClient("Key", "https://api.roninchain.com/rpc"),
		AuthToken:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjFlZjY2NzNiLTE3YWQtNjkxYS1hMDIyLTUwODQ3ODU4MmFmNSIsInNpZCI6MTgwNjgyMDU0LCJyb2xlcyI6WyJ1c2VyIl0sInNjcCI6WyJhbGwiXSwiYWN0aXZhdGVkIjp0cnVlLCJhY3QiOnRydWUsInJvbmluQWRkcmVzcyI6IjB4ODdkNjllMDFmY2QxMjY5NmM5NDNjMThiMzE5MmQ2OWE0NjAzNjU3OSIsImV4cCI6MTcyNzEyMDg0MCwiaWF0IjoxNzI1OTExMjQwLCJpc3MiOiJBeGllSW5maW5pdHkiLCJzdWIiOiIxZWY2NjczYi0xN2FkLTY5MWEtYTAyMi01MDg0Nzg1ODJhZjUifQ.5DEPE70KIFG1Rl6tCaa-cn1EP0ZLsgEmjosLvMDy5zM",
	}
	Wallet4 = WalletData{
		PrivateKey: "Key",
		PublicKey:  "Key",
		Client:     *defigo.NewClient("Key", "https://api.roninchain.com/rpc"),
		AuthToken:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjFlZjY2NzQ2LWM0NTEtNmRjNS1hZGI1LTIxZGI4ZDA1ZTU1YiIsInNpZCI6MTgwNjgyMDk1LCJyb2xlcyI6WyJ1c2VyIl0sInNjcCI6WyJhbGwiXSwiYWN0aXZhdGVkIjp0cnVlLCJhY3QiOnRydWUsInJvbmluQWRkcmVzcyI6IjB4ZDE1ODYwMTczNzRlNmZiMjQ5MjMwNzZkZDYzMWRiZmIzYWZlYjU4ZiIsImV4cCI6MTcyNzEyMDg4NywiaWF0IjoxNzI1OTExMjg3LCJpc3MiOiJBeGllSW5maW5pdHkiLCJzdWIiOiIxZWY2Njc0Ni1jNDUxLTZkYzUtYWRiNS0yMWRiOGQwNWU1NWIifQ.PGkAMHUQLS3-IkZuAz8gos-Oh5k6m6FP_Qw5J5cx7vU",
	}
	Wallet5 = WalletData{
		PrivateKey: "Key",
		PublicKey:  "Key",
		Client:     *defigo.NewClient("Key", "https://api.roninchain.com/rpc"),
		AuthToken:  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjFlZjY2NzJmLThhNjgtNmM2MC1iOWFhLTQ0MDU1YmE3YmUyNyIsInNpZCI6MTgwNjgxNzkzLCJyb2xlcyI6WyJ1c2VyIl0sInNjcCI6WyJhbGwiXSwiYWN0aXZhdGVkIjp0cnVlLCJhY3QiOnRydWUsInJvbmluQWRkcmVzcyI6IjB4OGM5M2FlYTAwMGQ2ZmI2ZDRmNWNhMGEzNTYwN2QzZjQxNTMzZmZiMCIsImV4cCI6MTcyNzEyMDU0NSwiaWF0IjoxNzI1OTEwOTQ1LCJpc3MiOiJBeGllSW5maW5pdHkiLCJzdWIiOiIxZWY2NjcyZi04YTY4LTZjNjAtYjlhYS00NDA1NWJhN2JlMjcifQ.j3PJpo9CCCosCBURWRxRY5WzKKXWijoz-obYcmWpBCA",
	}
)

var (
	Discount    = 0.6
	OurDiscount = 0.75
	MarketFee   = 0.025
)

var CollectionMap = map[string]CollectionValues{
	"moki": {
		CollectionName:     "moki",
		CollectionDiscount: (OurDiscount - MarketFee - 0.05),
		CollectionAddress:  "0x47b5a7c2e4f07772696bbf8c8c32fe2b9eabd550",
		CollectionCriteria: []map[string]interface{}{},
		Wallet:             &Wallet2,
	},
	"ragnarokGenesisTamer": {
		CollectionName:     "ragnarokGenesisTamer",
		CollectionDiscount: (OurDiscount - MarketFee - 0.025),
		CollectionAddress:  "0x6dcafe91533bdd733152ea30d029ec29280d7e4b",
		CollectionCriteria: []map[string]interface{}{},
		Wallet:             &Wallet4,
	},
	// "forgottenRuniverseSettlement": {
	// 	CollectionName:     "forgottenRuniverse",
	// 	CollectionDiscount: (Discount - MarketFee - 0.05),
	// 	CollectionAddress:  "0x775f0a0bb8258501d0862df38a7f7ad8f8f7423d",
	// 	CollectionCriteria: []map[string]interface{}{{"name": "type", "values": []string{"settlement"}}},
	// 	Wallet:             &WalletBGBK2,
	// },
	"forgottenRuniverseHomestead": {
		CollectionName:     "forgottenRuniverse",
		CollectionDiscount: (OurDiscount - MarketFee - 0.05),
		CollectionAddress:  "0x775f0a0bb8258501d0862df38a7f7ad8f8f7423d",
		CollectionCriteria: []map[string]interface{}{},
		Wallet:             &Wallet2,
	},
	"WildForestUnits": {
		CollectionName:     "WildForestUnits",
		CollectionDiscount: (OurDiscount - MarketFee - 0.05),
		CollectionAddress:  "0xa038c593115f6fcd673f6833e15462b475994879",
		CollectionCriteria: []map[string]interface{}{{"name": "rarity", "values": []string{"legendary"}}},
		Wallet:             &Wallet5,
	},
	"pixelsPets": {
		CollectionName:     "pixelsPets",
		CollectionDiscount: (OurDiscount - MarketFee - 0.05),
		CollectionAddress:  "0xb806028b6ebc35926442770a8a8a7aeab6e2ce5c",
		CollectionCriteria: []map[string]interface{}{},
		Wallet:             &Wallet5,
	},
	"pixelsLands": {
		CollectionName:     "pixelsLands",
		CollectionDiscount: (OurDiscount - MarketFee - 0.03),
		CollectionAddress:  "0xf083289535052e8449d69e6dc41c0ae064d8e3f6",
		CollectionCriteria: []map[string]interface{}{},
		Wallet:             &Wallet3,
	},
	"wildForestLords": {
		CollectionName:     "wildForestLords",
		CollectionDiscount: (OurDiscount - MarketFee - 0.05),
		CollectionAddress:  "0xa1ce53b661be73bf9a5edd3f0087484f0e3e7363",
		CollectionCriteria: []map[string]interface{}{},
		Wallet:             &Wallet3,
	},
}

type CoinData struct {
	LastTradeTs              int64
	SkipCheck                bool
	Asset                    string
	Market                   string
	Price                    float64
	BinanceFundingRate       float64
	BinanceRemainingOrderUsd float64
	BinanceSendingOrders     bool
	PosFoundInBinance        bool
	MarketSkew               *big.Int
	MarketKey                [32]byte
	MarketSkewAdj            float64
	MarketSkewUsd            float64
	MarketSkewBps            float64
	SafetyTrigger            bool
	AvgOpportunity           float64
	AvgOpportunityAbs        float64
	KeeperFeeBps             float64
	MarketSkewDiff           *big.Int
	KwentaFundingRate        *big.Int
	KwentaFundingRateAdj     float64
	FinalFundingRate         float64
	SkewScale                float64
	SkewScaleUsd             float64
	BpsToUsd                 float64
	MarketContract           *defigo.Contract
	SizeDelta                *big.Int
	MaxPositionUsd           float64
	MaxOrderSizeUsd          float64
	SkewScaleLevOffset       float64
	CurrentPosLevOffset      float64
	PositionSize             big.Int
	PositionFound            bool
	PositionSizeAdj          float64
	PositionSizeUsd          float64
	PositionMarginAdj        float64
	PositionLastPriceAdj     float64
	PositionPriceDiff        float64
	PositionExpectedMargin   float64
	PositionLeverage         float64
	PositionLeverageAdj      float64
	PositionAvgEntryPrice    big.Int
	PositionInitialMargin    big.Int
	PositionMargin           big.Int
	PositionLastPrice        big.Int
	PositionBinDiff          float64
	PositionBinDiffUsd       float64

	Settings struct {
		//DelayedOrderConfirmWindow    float64
		//LiquidationBufferRatio       float64
		//LiquidationPremiumMultiplier float64
		MakerFee       float64
		TakerFee       float64
		MaxLeverage    float64
		MaxMarketValue float64
		// SkewScale      float64
		//MakerFeeDelayedOrder         float64
		//MakerFeeOffchainDelayedOrder float64
		//MaxDelayTimeDelta            float64
		//MaxFundingVelocity           float64
		//MaxLiquidationDelta          float64
		//MaxPD                        float64
		//MinDelayTimeDelta            float64
		//NextPriceConfirmWindow       float64
		//OffchainDelayedOrderMaxAge   float64
		//OffchainDelayedOrderMinAge   float64
		//OffchainPriceDivergence      float64
		//TakerFeeDelayedOrder         float64
		//TakerFeeOffchainDelayedOrder float64
	}
}
type NftData struct {
	RonAmt                 float64
	RonAmtDecStr           string
	SignatureStr           string
	NftNumberStr           string
	UnixExpirationStr      string
	UnixStartStr           string
	OwnerAddressStr        string
	ExpectedStateStr       string
	PaymentTokenAddressStr string
}

var marketGatewayProxyContract *defigo.Contract
var routerContract *defigo.Contract
var FreeMargin float64 = 0.0
var SMIdleMargin float64 = 0.0
var MakerFeeBps float64 = 2.0
var FundingWeight float64 = 0.4
var CurrentPosLevMaxOffset float64 = 3.0

var InitialLev float64 = 10.0
var MaxLev float64 = 14.0
var MinLev float64 = 7.0

type KwentaCoinDataHandler func(*CoinData)

func init() {
	routerContract = defigo.NewContract(uniswapRouterContractAddress, uniswapRouterContractABI)
	marketGatewayProxyContract = defigo.NewContract(marketGatewayProxyAddress, marketGatewayProxyABI)
}
func SendNFTBuyOrder(ronAmtDecStr string, expectedState string, ownerAddress string, unixStart string, unixExpiration string, collectionAddress string, nftNumber string, signature string, spendingTokenAddress string, client *defigo.Web3Client) {
	ronAmtDecBig, _ := new(big.Int).SetString(ronAmtDecStr, 10)
	ronAmtHexStr := PadHexString(ronAmtDecBig.Text(16))

	callDataStr := "95a4ec0000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000e4f524445525f45584348414e474500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003e40abe93d80000000000000000000000000000000000000000000000000000000000000040"
	callDataStr += ronAmtHexStr
	//callDataStr += "valorEmRons(dec->hex)"
	callDataStr += "00000000000000000000000000000000000000000000000000000000000000c0"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000320"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000000"
	callDataStr += expectedState
	//callDataStr += "expectedState (dec->hex)"
	callDataStr += PadHexString(client.PublicKeyStr)
	//callDataStr += "minhaWallet (64bytes com 22 zeros)"
	callDataStr += PadHexString(client.PublicKeyStr)
	//callDataStr += "minhaWallet (64bytes com 22 zeros)"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000240"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000020"
	callDataStr += ownerAddress
	//callDataStr += "wallet do vendedor (64bytes com 22 zeros)"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000001"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000180"
	callDataStr += unixExpiration
	//callDataStr += "unix expiration"
	callDataStr += spendingTokenAddress
	//callDataStr += "000000000000000000000000e514d9deb7966c8be0ca922de8a064264ea6bcd4" // endereço WRON
	callDataStr += unixStart
	//callDataStr += "unix start"
	callDataStr += ronAmtHexStr
	//callDataStr += "valorEmRons(dec->hex)"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000000"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000000"
	callDataStr += expectedState
	//callDataStr += "expectedState (dec->hex)"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000000"
	callDataStr += "00000000000000000000000000000000000000000000000000000000000001a9"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000001"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000001"
	callDataStr += collectionAddress
	//callDataStr += "address da coleção (64bytes com 22 zeros)"
	callDataStr += nftNumber
	//callDataStr += "numero da NFT (dec->hex)"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000000"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000041"
	callDataStr += signature
	//callDataStr += "signature (130bytes)"
	callDataStr += "00000000000000000000000000000000000000000000000000000000000000"
	defigo.CallWriteFunctionRawData(marketGatewayProxyContract, client, ronAmtDecBig, callDataStr)
}
func GetTokenPrice(amount float64, AssetA string, AssetB string, client *defigo.Web3Client) float64 {
	AddressA := AssetToAddressMap[AssetA]
	AddressB := AssetToAddressMap[AssetB]

	ArgA := common.HexToAddress(AddressA)
	ArgB := common.HexToAddress(AddressB)
	amount *= 1000000
	amountInt := int(amount)
	amountStr := utils.IntToString(amountInt)
	if AssetA != "USDC" {
		amountStr += "000000000000"
	}

	amountOutMin := new(big.Int)
	amountOutMin.SetString(amountStr, 10)

	result := defigo.CallReadFunction(routerContract, client, "getAmountsOut", amountOutMin, []common.Address{ArgA, ArgB}).([]interface{})
	balanceArr := result[0].([]*big.Int)
	balance := balanceArr[1]
	balanceBigFloat := defigo.WeiToEther(balance)
	balanceFloat, _ := balanceBigFloat.Float64()
	if AssetB == "USDC" {
		balanceFloat *= 1000000000000
	}
	return balanceFloat
}

type NFTDetailsVariables struct {
	TokenID      string `json:"tokenId"`
	TokenAddress string `json:"tokenAddress"`
}

type NFTDetailsRequestBody struct {
	OperationName string              `json:"operationName"`
	Variables     NFTDetailsVariables `json:"variables"`
	Query         string              `json:"query"`
}

type NFTDetailsResponseData struct {
	Data struct {
		Erc721Token struct {
			TokenAddress string              `json:"tokenAddress"`
			TokenId      string              `json:"tokenId"`
			Slug         string              `json:"slug"`
			Owner        string              `json:"owner"`
			Name         string              `json:"name"`
			Order        NFTDetailsOrderData `json:"order"`
		} `json:"erc721Token"`
	} `json:"data"`
}
type NFTDetailsOrderData struct {
	ExpiredAt     int64  `json:"expiredAt"`
	StartedAt     int64  `json:"startedAt"`
	BasePrice     string `json:"basePrice"`
	ExpectedState string `json:"expectedState"`
	Signature     string `json:"signature"`
	PaymentToken  string `json:"paymentToken"`
}

func GetNFTDetails(nftNumber int, collectionAddress string) (string, string, string, string, string, string, string, string, string) {
	url := "https://marketplace-graphql.skymavis.com/graphql"

	// Define the request body
	requestBody := NFTDetailsRequestBody{
		OperationName: "GetERC721TokensList",
		Variables: NFTDetailsVariables{
			TokenID:      utils.IntToString(nftNumber),
			TokenAddress: collectionAddress,
		},
		Query: `query GetERC721TokenDetail($tokenAddress: String, $slug: String, $tokenId: String!) {
			erc721Token(tokenAddress: $tokenAddress, slug: $slug, tokenId: $tokenId) {
				...Erc721Token
				__typename
			}
		}

		fragment Erc721Token on Erc721 {
			tokenAddress
			tokenId
			slug
			owner
			name
			order {
				...OrderInfo
				__typename
			}
			transferHistory(from: 0, size: 1) {
				total
				results {
					...TransferRecordBrief
					__typename
				}
				__typename
			}
			minPrice
			attributes
			image
			cdnImage
			video
			animationUrl
			ownerProfile {
				...PublicProfileBrief
				__typename
			}
			traitDistribution {
				...TokenTrait
				__typename
			}
			isLocked
			yourOffer {
				...OfferInfo
				__typename
			}
			highestOffer {
				...OfferInfo
				__typename
			}
			collectionMetadata
			__typename
		}

		fragment OfferInfo on Order {
			id
			maker
			kind
			assets {
				...OfferAssetInfo
				__typename
			}
			expiredAt
			paymentToken
			startedAt
			basePrice
			expectedState
			nonce
			marketFeePercentage
			signature
			hash
			duration
			timeLeft
			currentPrice
			suggestedPrice
			makerProfile {
				...PublicProfileBrief
				__typename
			}
			orderStatus
			orderQuantity {
				orderId
				quantity
				remainingQuantity
				availableQuantity
				__typename
			}
			__typename
		}

		fragment OfferAssetInfo on Asset {
			erc
			address
			id
			quantity
			token {
				... on Erc721 {
					tokenAddress
					tokenId
					slug
					image
					cdnImage
					name
					owner
					ownerProfile {
						...PublicProfileBrief
						__typename
					}
					minPrice
					collectionMetadata
					__typename
				}
				__typename
			}
			__typename
		}

		fragment PublicProfileBrief on PublicProfile {
			accountId
			addresses {
				...Addresses
				__typename
			}
			activated
			name
			__typename
		}

		fragment Addresses on NetAddresses {
			ethereum
			ronin
			__typename
		}

		fragment OrderInfo on Order {
			id
			maker
			kind
			assets {
				...AssetInfo
				__typename
			}
			expiredAt
			paymentToken
			startedAt
			basePrice
			expectedState
			nonce
			marketFeePercentage
			signature
			hash
			duration
			timeLeft
			currentPrice
			suggestedPrice
			makerProfile {
				...PublicProfileBrief
				__typename
			}
			orderStatus
			orderQuantity {
				orderId
				quantity
				remainingQuantity
				availableQuantity
				__typename
			}
			__typename
		}

		fragment AssetInfo on Asset {
			erc
			address
			id
			quantity
			__typename
		}

		fragment TransferRecordBrief on TransferRecord {
			tokenId
			from
			to
			fromProfile {
				...PublicProfileBrief
				__typename
			}
			toProfile {
				...PublicProfileBrief
				__typename
			}
			timestamp
			txHash
			withPrice
			quantity
			paymentToken
			__typename
		}

		fragment TokenTrait on TokenTrait {
			key
			value
			count
			percentage
			displayType
			maxValue
			__typename
		}`,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	// Unmarshal the response body to a Go struct
	var responseData NFTDetailsResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error unmarshalling response JSON:", err)
	}
	ownerAddress := responseData.Data.Erc721Token.Owner
	order := responseData.Data.Erc721Token.Order
	ronAmtStr := order.BasePrice
	startedAt := utils.Int64ToString(order.StartedAt)
	expiredAt := utils.Int64ToString(order.ExpiredAt)
	signature := order.Signature
	paymentToken := order.PaymentToken
	expectedState := order.ExpectedState
	nftNumberStr := utils.IntToString(nftNumber)
	return ronAmtStr, expectedState, ownerAddress, startedAt, expiredAt, collectionAddress, nftNumberStr, signature, paymentToken
	// Print the response data
}

type NftListVariables struct {
	From          int                      `json:"from"`
	AuctionType   string                   `json:"auctionType"`
	Criteria      []map[string]interface{} `json:"criteria"`
	Size          int                      `json:"size"`
	Sort          string                   `json:"sort"`
	RangeCriteria []string                 `json:"rangeCriteria"`
	TokenAddress  string                   `json:"tokenAddress"`
}

type NftListRequestBody struct {
	OperationName string           `json:"operationName"`
	Variables     NftListVariables `json:"variables"`
	Query         string           `json:"query"`
}

type NftListResponseData struct {
	Data struct {
		Erc721Tokens struct {
			Total   int `json:"total"`
			Results []struct {
				TokenAddress string              `json:"tokenAddress"`
				TokenId      string              `json:"tokenId"`
				Slug         string              `json:"slug"`
				Owner        string              `json:"owner"`
				Name         string              `json:"name"`
				Order        NFTDetailsOrderData `json:"order"`
				// Add more fields as needed
			} `json:"results"`
		} `json:"erc721Tokens"`
	} `json:"data"`
}

// func GetNFTList(collection CollectionValues) (float64, string, string, string, string, string, string, string, string) {
// 	url := "https://marketplace-graphql.skymavis.com/graphql"
// 	// Define the request body
// 	requestBody := NftListRequestBody{
// 		OperationName: "GetERC721TokensList",
// 		Variables: NftListVariables{
// 			From:          0,
// 			AuctionType:   "All",
// 			Criteria:      collection.CollectionCriteria,
// 			Size:          50,
// 			Sort:          "PriceAsc",
// 			RangeCriteria: []string{},
// 			TokenAddress:  collection.CollectionAddress,
// 		},
// 		Query: `query GetERC721TokensList($tokenAddress: String, $slug: String, $owner: String, $auctionType: AuctionType, $criteria: [SearchCriteria!], $from: Int!, $size: Int!, $sort: SortBy, $name: String, $priceRange: InputRange, $rangeCriteria: [RangeSearchCriteria!], $excludeAddress: String) {
// 			erc721Tokens(
// 				tokenAddress: $tokenAddress
// 				slug: $slug
// 				owner: $owner
// 				auctionType: $auctionType
// 				criteria: $criteria
// 				from: $from
// 				size: $size
// 				sort: $sort
// 				name: $name
// 				priceRange: $priceRange
// 				rangeCriteria: $rangeCriteria
// 				excludeAddress: $excludeAddress
// 			) {
// 				total
// 				results {
// 					...Erc721TokenBrief
// 					__typename
// 				}
// 				__typename
// 			}
// 		}
// 		fragment Erc721TokenBrief on Erc721 {
// 			tokenAddress
// 			tokenId
// 			slug
// 			owner
// 			name
// 			order {
// 				...OrderInfo
// 				__typename
// 			}
// 			image
// 			cdnImage
// 			video
// 			isLocked
// 			attributes
// 			traitDistribution {
// 				...TokenTrait
// 				__typename
// 			}
// 			collectionMetadata
// 			ownerProfile {
// 				name
// 				accountId
// 				__typename
// 			}
// 			__typename
// 		}
// 		fragment OrderInfo on Order {
// 			id
// 			maker
// 			kind
// 			assets {
// 				...AssetInfo
// 				__typename
// 			}
// 			expiredAt
// 			paymentToken
// 			startedAt
// 			basePrice
// 			expectedState
// 			nonce
// 			marketFeePercentage
// 			signature
// 			hash
// 			duration
// 			timeLeft
// 			currentPrice
// 			suggestedPrice
// 			makerProfile {
// 				...PublicProfileBrief
// 				__typename
// 			}
// 			orderStatus
// 			orderQuantity {
// 				orderId
// 				quantity
// 				remainingQuantity
// 				availableQuantity
// 				__typename
// 			}
// 			__typename
// 		}
// 		fragment AssetInfo on Asset {
// 			erc
// 			address
// 			id
// 			quantity
// 			__typename
// 		}
// 		fragment PublicProfileBrief on PublicProfile {
// 			accountId
// 			addresses {
// 				...Addresses
// 				__typename
// 			}
// 			activated
// 			name
// 			__typename
// 		}
// 		fragment Addresses on NetAddresses {
// 			ethereum
// 			ronin
// 			__typename
// 		}
// 		fragment TokenTrait on TokenTrait {
// 			key
// 			value
// 			count
// 			percentage
// 			displayType
// 			maxValue
// 			__typename
// 		}`,
// 	}

// 	// Marshal the request body to JSON
// 	jsonBody, err := json.Marshal(requestBody)
// 	if err != nil {
// 		fmt.Println("Error marshalling JSON:", err)
// 	}

// 	// Create a new HTTP POST request
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
// 	if err != nil {
// 		fmt.Println("Error creating request:", err)
// 	}

// 	// Set the request headers
// 	req.Header.Set("Content-Type", "application/json")

// 	// Create an HTTP client and send the request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending request:", err)
// 	}
// 	defer resp.Body.Close()

// 	// Read and parse the response body
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 	}

// 	// Unmarshal the response body to a Go struct
// 	err = json.Unmarshal(body, &responseData)
// 	if err != nil {
// 		fmt.Printf("Error:%v", err)
// 	}
// 	nftNumberStr := responseData.Data.Erc721Tokens.Results[0].TokenId
// 	ownerAddress := responseData.Data.Erc721Tokens.Results[0].Owner
// 	order := responseData.Data.Erc721Tokens.Results[0].Order
// 	signature := order.Signature
// 	paymentToken := order.PaymentToken
// 	expectedState := order.ExpectedState
// 	startedAt := utils.Int64ToString(order.StartedAt)
// 	expiredAt := utils.Int64ToString(order.ExpiredAt)

// 	ronAmtStr := order.BasePrice
// 	ronAmtBigDec := new(big.Int)
// 	ronAmtBigDec.SetString(ronAmtStr, 10)
// 	ronAmtBigFloat := new(big.Float).Quo(new(big.Float).SetInt(ronAmtBigDec), big.NewFloat(params.Ether))
// 	ronAmt, _ := ronAmtBigFloat.Float64()

// 	return ronAmt, ronAmtStr, expectedState, ownerAddress, startedAt, expiredAt, nftNumberStr, signature, paymentToken
// }

func GetNFTListAll(tokenAddress string, criteriaJson []map[string]interface{}) []NftData {
	url := "https://marketplace-graphql.skymavis.com/graphql"
	// Convert criteriaJson to JSON string
	criteriaBytes, err := json.Marshal(criteriaJson)
	if err != nil {
		fmt.Println("Error marshaling criteriaJson:", err)
		return nil
	}
	criteriaJsonStr := string(criteriaBytes)
	// Define the request body
	requestBody := fmt.Sprintf(
		"{\"operationName\":\"GetERC721TokensList\",\"variables\":{\"from\":0,\"auctionType\":\"All\",\"size\":50,\"sort\":\"PriceAsc\",\"criteria\":%v,\"rangeCriteria\":[],\"tokenAddress\":\"%s\"},\"query\":\"query GetERC721TokensList($tokenAddress: String, $slug: String, $owner: String, $auctionType: AuctionType, $criteria: [SearchCriteria!], $from: Int!, $size: Int!, $sort: SortBy, $name: String, $priceRange: InputRange, $rangeCriteria: [RangeSearchCriteria!], $excludeAddress: String) {\\n  erc721Tokens(\\n    tokenAddress: $tokenAddress\\n    slug: $slug\\n    owner: $owner\\n    auctionType: $auctionType\\n    criteria: $criteria\\n    from: $from\\n    size: $size\\n    sort: $sort\\n    name: $name\\n    priceRange: $priceRange\\n    rangeCriteria: $rangeCriteria\\n    excludeAddress: $excludeAddress\\n  ) {\\n    total\\n    results {\\n      ...Erc721TokenBrief\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\\nfragment Erc721TokenBrief on Erc721 {\\n  tokenAddress\\n  tokenId\\n  slug\\n  owner\\n  name\\n  order {\\n    ...OrderInfo\\n    __typename\\n  }\\n  image\\n  cdnImage\\n  video\\n  isLocked\\n  attributes\\n  traitDistribution {\\n    ...TokenTrait\\n    __typename\\n  }\\n  collectionMetadata\\n  ownerProfile {\\n    name\\n    accountId\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment OrderInfo on Order {\\n  id\\n  maker\\n  kind\\n  assets {\\n    ...AssetInfo\\n    __typename\\n  }\\n  expiredAt\\n  paymentToken\\n  startedAt\\n  basePrice\\n  expectedState\\n  nonce\\n  marketFeePercentage\\n  signature\\n  hash\\n  duration\\n  timeLeft\\n  currentPrice\\n  suggestedPrice\\n  makerProfile {\\n    ...PublicProfileBrief\\n    __typename\\n  }\\n  orderStatus\\n  orderQuantity {\\n    orderId\\n    quantity\\n    remainingQuantity\\n    availableQuantity\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment AssetInfo on Asset {\\n  erc\\n  address\\n  id\\n  quantity\\n  __typename\\n}\\n\\nfragment PublicProfileBrief on PublicProfile {\\n  accountId\\n  addresses {\\n    ...Addresses\\n    __typename\\n  }\\n  activated\\n  name\\n  __typename\\n}\\n\\nfragment Addresses on NetAddresses {\\n  ethereum\\n  ronin\\n  __typename\\n}\\n\\nfragment TokenTrait on TokenTrait {\\n  key\\n  value\\n  count\\n  percentage\\n  displayType\\n  maxValue\\n  __typename\\n}\\n\"}",
		criteriaJsonStr, tokenAddress)

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	var responseData NftListResponseData
	// Unmarshal the response body to a Go struct
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Printf("Error:%v", err)
	}
	nftDataLs := []NftData{}
	for _, result := range responseData.Data.Erc721Tokens.Results {

		nftNumberStr := result.TokenId
		ownerAddress := result.Owner
		order := result.Order
		signature := order.Signature
		paymentToken := order.PaymentToken
		expectedState := order.ExpectedState
		startedAt := utils.Int64ToString(order.StartedAt)
		expiredAt := utils.Int64ToString(order.ExpiredAt)

		ronAmtStr := order.BasePrice
		ronAmtBigDec := new(big.Int)
		ronAmtBigDec.SetString(ronAmtStr, 10)
		ronAmtBigFloat := new(big.Float).Quo(new(big.Float).SetInt(ronAmtBigDec), big.NewFloat(params.Ether))
		ronAmt, _ := ronAmtBigFloat.Float64()
		if expectedState == "" && ronAmt == 0 {
			continue
		}
		newNftData := NftData{
			RonAmt:                 ronAmt,
			RonAmtDecStr:           ronAmtStr,
			SignatureStr:           signature,
			NftNumberStr:           nftNumberStr,
			UnixExpirationStr:      expiredAt,
			UnixStartStr:           startedAt,
			OwnerAddressStr:        ownerAddress,
			ExpectedStateStr:       expectedState,
			PaymentTokenAddressStr: paymentToken,
		}
		nftDataLs = append(nftDataLs, newNftData)
	}

	return nftDataLs
}

func PadHexString(hexStr string) string {
	// Remove the "0x" prefix if present
	hexStr = strings.TrimPrefix(hexStr, "0x")

	// Calculate the number of leading zeros needed
	padLength := 64 - len(hexStr)

	// Pad the string with leading zeros
	paddedHexStr := fmt.Sprintf("%0*s", padLength+len(hexStr), hexStr)

	return paddedHexStr
}

func FloatTo18PlacesString(origFloat float64) string {
	// Create a big float with the initial value
	f := new(big.Float).SetPrec(200).SetFloat64(origFloat)

	// Create a big int with the value 10^18
	multiplier := new(big.Int)
	multiplier.SetString("1000000000000000000", 10) // 10^18

	// Convert big int to big float
	bigMultiplier := new(big.Float).SetInt(multiplier)

	// Multiply the float by 10^18
	result := new(big.Float).Mul(f, bigMultiplier)

	// Convert the result to an integer
	intResult := new(big.Int)
	result.Int(intResult)

	// Convert the integer result to a string
	resultStr := intResult.Text(10)

	if len(resultStr) > 0 {
		resultStr = resultStr[:len(resultStr)-1] + "0"
	}

	return resultStr
}
