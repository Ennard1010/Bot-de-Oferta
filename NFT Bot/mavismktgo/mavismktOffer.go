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
	"os/exec"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/params"
)

type AssetType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type EIP712Domain struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainId           string `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
}

type Order struct {
	Maker               string  `json:"maker"`
	Kind                string  `json:"kind"`
	Assets              []Asset `json:"assets"`
	ExpiredAt           string  `json:"expiredAt"`
	PaymentToken        string  `json:"paymentToken"`
	StartedAt           string  `json:"startedAt"`
	BasePrice           string  `json:"basePrice"`
	EndedAt             string  `json:"endedAt"`
	EndedPrice          string  `json:"endedPrice"`
	ExpectedState       string  `json:"expectedState"`
	Nonce               string  `json:"nonce"`
	MarketFeePercentage string  `json:"marketFeePercentage"`
}

type Asset struct {
	Erc      string `json:"erc"`
	Addr     string `json:"addr"`
	Id       string `json:"id"`
	Quantity string `json:"quantity"`
}

type Message struct {
	Maker               string  `json:"maker"`
	Kind                string  `json:"kind"`
	Assets              []Asset `json:"assets"`
	ExpiredAt           string  `json:"expiredAt"`
	PaymentToken        string  `json:"paymentToken"`
	StartedAt           string  `json:"startedAt"`
	BasePrice           string  `json:"basePrice"`
	EndedAt             string  `json:"endedAt"`
	EndedPrice          string  `json:"endedPrice"`
	ExpectedState       string  `json:"expectedState"`
	Nonce               string  `json:"nonce"`
	MarketFeePercentage string  `json:"marketFeePercentage"`
}

func MakeSignatureRequest(jsonData string, privateKey string) string {

	cmd := exec.Command("node", "teste.js", jsonData, privateKey)
	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running script: %v\n", err)
		return ""
	}

	return strings.TrimSpace(string(output))
}

type OfferInfo struct {
	NftId                  string
	Name                   string
	RonAmt                 string
	SignatureStr           string
	NftNumberStr           string
	UnixExpirationStr      string
	UnixStartStr           string
	OwnerAddressStr        string
	ExpectedStateStr       string
	PaymentTokenAddressStr string
	CollectionAddressStr   string
}

func GetSentOffers(authToken string, checkId string, collectionAddress string) []OfferInfo {
	url := "https://marketplace-graphql.skymavis.com/graphql"
	requestBody := fmt.Sprintf(
		"{\"operationName\":\"GetSentOffers\",\"variables\":{\"from\":0,\"size\":10,\"User\":\"%s\",\"sort\":\"ExpiredAtAsc\",\"collectibleFilters\":{\"tokenAddresses\":[\"%s\"]}},\"query\":\"query GetSentOffers($from: Int!, $size: Int!, $isValid: Boolean, $sort: OfferSortBy, $collectibleFilters: CollectibleFilter!) {\\n  sentOffers(\\n    collectibleFilters: $collectibleFilters\\n    from: $from\\n    size: $size\\n    sort: $sort\\n    isValid: $isValid\\n  ) {\\n    total\\n    data {\\n      ...OfferInfo\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\\nfragment OfferInfo on Order {\\n  id\\n  maker\\n  kind\\n  assets {\\n    ...OfferAssetInfo\\n    __typename\\n  }\\n  expiredAt\\n  paymentToken\\n  startedAt\\n  basePrice\\n  expectedState\\n  nonce\\n  marketFeePercentage\\n  signature\\n  hash\\n  duration\\n  timeLeft\\n  currentPrice\\n  suggestedPrice\\n  makerProfile {\\n    ...PublicProfileBrief\\n    __typename\\n  }\\n  orderStatus\\n  orderQuantity {\\n    orderId\\n    quantity\\n    remainingQuantity\\n    availableQuantity\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment OfferAssetInfo on Asset {\\n  erc\\n  address\\n  id\\n  quantity\\n  token {\\n    ... on Erc721 {\\n      tokenAddress\\n      tokenId\\n      slug\\n      image\\n      cdnImage\\n      name\\n      owner\\n     ownerProfile {\\n        ...PublicProfileBrief\\n        __typename\\n      }\\n      minPrice\\n      collectionMetadata\\n      __typename\\n    }\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment PublicProfileBrief on PublicProfile {\\n  accountId\\n  addresses {\\n    ...Addresses\\n    __typename\\n  }\\n  activated\\n  name\\n  __typename\\n}\\n\\nfragment Addresses on NetAddresses {\\n  ethereum\\n  ronin\\n  __typename\\n}\\n\"}",
		checkId, collectionAddress,
	)
	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", authToken)
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
	var offerData OfferResponseData
	// Unmarshal the response body to a Go struct
	err = json.Unmarshal(body, &offerData)
	if err != nil {
		fmt.Printf("Error:%v", err)
	}
	OfferLS := []OfferInfo{}
	for _, data := range offerData.Data.SentOffers.Data {
		ronAmt := data.BasePrice
		name := data.Assets[0].Token.Name
		nftNumberStr := data.Assets[0].Token.TokenID
		ownerAddress := data.Maker
		collectionAddress := data.Assets[0].Address
		paymentToken := data.PaymentToken
		expectedState := data.ExpectedState
		startedAt := utils.Int64ToString(data.StartedAt)
		expiredAt := utils.Int64ToString(data.ExpiredAt)
		NewOfferData := OfferInfo{
			Name:                   name,
			NftNumberStr:           nftNumberStr,
			OwnerAddressStr:        ownerAddress,
			PaymentTokenAddressStr: paymentToken,
			UnixExpirationStr:      expiredAt,
			UnixStartStr:           startedAt,
			ExpectedStateStr:       expectedState,
			CollectionAddressStr:   collectionAddress,
			RonAmt:                 ronAmt,
		}
		OfferLS = append(OfferLS, NewOfferData)
	}
	return OfferLS
}

type OfferResponseData struct {
	Data struct {
		SentOffers struct {
			Total int `json:"total"`
			Data  []struct {
				ID     int64  `json:"id"`
				Maker  string `json:"maker"`
				Kind   string `json:"kind"`
				Assets []struct {
					Erc      string `json:"erc"`
					Address  string `json:"address"`
					ID       string `json:"id"`
					Quantity string `json:"quantity"`
					Token    struct {
						TokenAddress string `json:"tokenAddress"`
						TokenID      string `json:"tokenId"`
						Name         string `json:"name"`
						Owner        string `json:"owner"`
						MinPrice     string `json:"minPrice"`
					} `json:"token"`
				} `json:"assets"`
				ExpiredAt      int64  `json:"expiredAt"`
				PaymentToken   string `json:"paymentToken"`
				StartedAt      int64  `json:"startedAt"`
				BasePrice      string `json:"basePrice"`
				CurrentPrice   string `json:"currentPrice"`
				SuggestedPrice string `json:"suggestedPrice"`
				ExpectedState  string `json:"expectedState"`
			} `json:"data"`
		} `json:"sentOffers"`
	} `json:"data"`
}

func GetFloorPrice(collectionAddress string) (floorPrice string) {
	url := "https://marketplace-graphql.skymavis.com/graphql"
	requestBody := fmt.Sprintf(
		"{\"operationName\":\"GetTokenData\",\"variables\":{\"tokenAddress\":\"%s\"},\"query\":\"query GetTokenData($tokenAddress: String, $slug: String) {\\n  tokenData(tokenAddress: $tokenAddress, slug: $slug) {\\n    ...TokenData\\n    allowedPaymentTokens\\n    __typename\\n  }\\n}\\n\\nfragment TokenData on TokenData {\\n  tokenAddress\\n  slug\\n  collectionMetadata\\n  volumeAllTime\\n  totalOwners\\n  totalItems\\n  totalListing\\n  minPrice\\n  erc\\n  groupTraits\\n  content\\n  __typename\\n}\\n\"}",
		collectionAddress,
	)
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
	var FloorData = Floor{}
	// Unmarshal the response body to a Go struct
	err = json.Unmarshal(body, &FloorData)
	if err != nil {
		fmt.Printf("Error:%v", err)
	}
	floorPrice = FloorData.Data.TokenData.MinimumPrice
	return

}
func GetNFTTokenPersonalInfo(nftNumber string, collectionAddress string, checkId string) (rontAmtPO string, nftNumberPO string, tokenAddressPO string, makerPO string, expectedStatePO string, unixStartPO string, unixExpirationPO string) {
	url := "https://marketplace-graphql.skymavis.com/graphql"

	requestBody := fmt.Sprintf(
		"{\"operationName\":\"GetERC721Offers\",\"variables\":{\"tokenAddress\":\"%s\",\"tokenId\":\"%s\",\"from\":0,\"size\":5},\"query\":\"query GetERC721Offers($tokenAddress: String, $tokenId: String!, $from: Int!, $size: Int!) {\\n  erc721Token(tokenAddress: $tokenAddress, tokenId: $tokenId) {\\n    numActiveOffers\\n    order {\\n      currentPrice\\n      paymentToken\\n      __typename\\n    }\\n    transferHistory(from: 0, size: 1) {\\n      total\\n      results {\\n        withPrice\\n        paymentToken\\n        __typename\\n      }\\n      __typename\\n    }\\n    offers(from: $from, size: $size) {\\n      ...OfferInfo\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\\nfragment OfferInfo on Order {\\n  id\\n  maker\\n  kind\\n  assets {\\n    ...OfferAssetInfo\\n    __typename\\n  }\\n  expiredAt\\n  paymentToken\\n  startedAt\\n  basePrice\\n  expectedState\\n  nonce\\n  marketFeePercentage\\n  signature\\n  hash\\n  duration\\n  timeLeft\\n  currentPrice\\n  suggestedPrice\\n  makerProfile {\\n    ...PublicProfileBrief\\n    __typename\\n  }\\n  orderStatus\\n  orderQuantity {\\n    orderId\\n    quantity\\n    remainingQuantity\\n    availableQuantity\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment OfferAssetInfo on Asset {\\n  erc\\n  address\\n  id\\n  quantity\\n  token {\\n    ... on Erc721 {\\n      tokenAddress\\n      tokenId\\n      slug\\n      image\\n      cdnImage\\n      name\\n      owner\\n      ownerProfile {\\n        ...PublicProfileBrief\\n        __typename\\n      }\\n      minPrice\\n      collectionMetadata\\n      __typename\\n    }\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment PublicProfileBrief on PublicProfile {\\n  accountId\\n  addresses {\\n    ...Addresses\\n    __typename\\n  }\\n  activated\\n  name\\n  __typename\\n}\\n\\nfragment Addresses on NetAddresses {\\n  ethereum\\n  ronin\\n  __typename\\n}\\n\"}",
		collectionAddress, nftNumber,
	)
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
	var PersonalData NFTpersonalResponseData
	// Unmarshal the response body to a Go struct
	err = json.Unmarshal(body, &PersonalData)
	if err != nil {
		fmt.Printf("Error:%v", err)
	}
	for _, offer := range PersonalData.Data.Erc721Token.Offers {
		if offer.Maker == checkId {
			rontAmtPO = offer.BasePrice
			nftNumberPO = offer.Assets[0].NFTNumber
			tokenAddressPO = offer.PaymentToken
			makerPO = offer.Maker
			expectedStatePO = offer.ExpectedState
			unixStartPO = utils.Int64ToString(offer.StartedAt)
			unixExpirationPO = utils.Int64ToString(offer.ExpiredAt)
		} else {
			continue
		}
	}
	return
}

type NFTpersonalResponseData struct {
	Data struct {
		Erc721Token struct {
			Offers []struct {
				BasePrice string `json:"basePrice"`
				Maker     string `json:"maker"`
				Assets    []struct {
					NFTNumber string `json:"id"`
					Address   string `json:"address"`
					Erc       string `json:"erc"`
					Quantity  string `json:"quantity"`
				} `json:"assets"`
				StartedAt     int64  `json:"startedAt"`
				ExpiredAt     int64  `json:"expiredAt"`
				ExpectedState string `json:"expectedState"`
				PaymentToken  string `json:"paymentToken"`
			} `json:"offers"`
		} `json:"erc721Token"`
	} `json:"data"`
}

func GetNFTTokenHOP(nftNumber string, tokenAddress string) (basePrice string, makerStr string) {
	url := "https://marketplace-graphql.skymavis.com/graphql"
	requestBody := fmt.Sprintf(
		"{\"operationName\":\"GetERC721TokenDetail\",\"variables\":{\"tokenId\":\"%s\",\"tokenAddress\":\"%s\"},\"query\":\"query GetERC721TokenDetail($tokenAddress: String, $slug: String, $tokenId: String!) {\\n  erc721Token(tokenAddress: $tokenAddress, slug: $slug, tokenId: $tokenId) {\\n    ...Erc721Token\\n    __typename\\n  }\\n}\\n\\nfragment Erc721Token on Erc721 {\\n  tokenAddress\\n  tokenId\\n  slug\\n  owner\\n  name\\n  order {\\n    ...OrderInfo\\n    __typename\\n  }\\n  transferHistory(from: 0, size: 1) {\\n    total\\n    results {\\n      ...TransferRecordBrief\\n      __typename\\n    }\\n    __typename\\n  }\\n  minPrice\\n  attributes\\n  image\\n  cdnImage\\n  video\\n  animationUrl\\n  ownerProfile {\\n    ...PublicProfileBrief\\n    __typename\\n  }\\n  traitDistribution {\\n    ...TokenTrait\\n    __typename\\n  }\\n  isLocked\\n  yourOffer {\\n    ...OfferInfo\\n    __typename\\n  }\\n  highestOffer {\\n    ...OfferInfo\\n    __typename\\n  }\\n  collectionMetadata\\n  __typename\\n}\\n\\nfragment OfferInfo on Order {\\n  id\\n  maker\\n  kind\\n  assets {\\n    ...OfferAssetInfo\\n    __typename\\n  }\\n  expiredAt\\n  paymentToken\\n  startedAt\\n  basePrice\\n  expectedState\\n  nonce\\n  marketFeePercentage\\n  signature\\n  hash\\n  duration\\n  timeLeft\\n  currentPrice\\n  suggestedPrice\\n  makerProfile {\\n    ...PublicProfileBrief\\n    __typename\\n  }\\n  orderStatus\\n  orderQuantity {\\n    orderId\\n    quantity\\n    remainingQuantity\\n    availableQuantity\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment OfferAssetInfo on Asset {\\n  erc\\n  address\\n  id\\n  quantity\\n  token {\\n    ... on Erc721 {\\n      tokenAddress\\n      tokenId\\n      slug\\n      image\\n      cdnImage\\n      name\\n      owner\\n      ownerProfile {\\n        ...PublicProfileBrief\\n        __typename\\n      }\\n      minPrice\\n      collectionMetadata\\n      __typename\\n    }\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment PublicProfileBrief on PublicProfile {\\n  accountId\\n  addresses {\\n    ...Addresses\\n    __typename\\n  }\\n  activated\\n  name\\n  __typename\\n}\\n\\nfragment Addresses on NetAddresses {\\n  ethereum\\n  ronin\\n  __typename\\n}\\n\\nfragment OrderInfo on Order {\\n  id\\n  maker\\n  kind\\n  assets {\\n    ...AssetInfo\\n    __typename\\n  }\\n  expiredAt\\n  paymentToken\\n  startedAt\\n  basePrice\\n  expectedState\\n  nonce\\n  marketFeePercentage\\n  signature\\n  hash\\n  duration\\n  timeLeft\\n  currentPrice\\n  suggestedPrice\\n  makerProfile {\\n    ...PublicProfileBrief\\n    __typename\\n  }\\n  orderStatus\\n  orderQuantity {\\n    orderId\\n    quantity\\n    remainingQuantity\\n    availableQuantity\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment AssetInfo on Asset {\\n  erc\\n  address\\n  id\\n  quantity\\n  __typename\\n}\\n\\nfragment TransferRecordBrief on TransferRecord {\\n  tokenId\\n  from\\n  to\\n  fromProfile {\\n    ...PublicProfileBrief\\n    __typename\\n  }\\n  toProfile {\\n    ...PublicProfileBrief\\n    __typename\\n  }\\n  timestamp\\n  txHash\\n  withPrice\\n  quantity\\n  paymentToken\\n  __typename\\n}\\n\\nfragment TokenTrait on TokenTrait {\\n  key\\n  value\\n  count\\n  percentage\\n  displayType\\n  maxValue\\n  __typename\\n}\\n\"}",
		nftNumber, tokenAddress,
	)

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// Define a local struct to hold the response data
	var localResponseData struct {
		Data struct {
			Erc721Token struct {
				HighestOffer struct {
					BasePrice string `json:"basePrice"`
					Maker     string `json:"maker"`
				} `json:"highestOffer"`
			} `json:"erc721Token"`
		} `json:"data"`
	}

	// Unmarshal the response body to the local struct
	err = json.Unmarshal(body, &localResponseData)
	if err != nil {
		return
	}

	basePrice = localResponseData.Data.Erc721Token.HighestOffer.BasePrice
	makerStr = localResponseData.Data.Erc721Token.HighestOffer.Maker

	return basePrice, makerStr
}

func ConvertStringtoFloat(convertionNumber string) float64 {
	ronAmtBigDec := new(big.Int)
	ronAmtBigDec.SetString(convertionNumber, 10)
	ronAmtBigFloat := new(big.Float).Quo(new(big.Float).SetInt(ronAmtBigDec), big.NewFloat(params.Ether))
	floatNumber, _ := ronAmtBigFloat.Float64()
	return floatNumber
}
func MultiplyNftPrice(basePriceMultiplyer string, multiplier *big.Float) string {
	if basePriceMultiplyer == "" {
		fmt.Print("Error multiplying")
	}

	// Convert the string to big.Float
	num, err := stringToBigFloat(basePriceMultiplyer)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Perform multiplication
	result := new(big.Float).Mul(num, multiplier)

	// Convert the result back to a string without decimal places
	formattedResult := BigFloatToStringWithoutDecimal(result, 18)
	newLength := len(formattedResult) - 18
	finalResult := formattedResult[:newLength]
	return finalResult
}

// Convert a string with 18 decimal places to a big.Float
func stringToBigFloat(s string) (*big.Float, error) {
	f, _, err := big.ParseFloat(s, 10, 0, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to big.Float: %v", err)
	}
	return f, nil
}
func FloatToStringWithoutDot(number float64) string {
	// Formata o número com 18 casas decimais
	formatted := fmt.Sprintf("%.18f", number)

	// Remove o ponto decimal
	withoutDot := strings.ReplaceAll(formatted, ".", "")

	return withoutDot
}

// Convert a big.Float to a string without decimal places
func BigFloatToStringWithoutDecimal(f *big.Float, decimalPlaces int) string {
	// Shift decimal places to the right
	multiplier := new(big.Float).SetFloat64(1)
	for i := 0; i < decimalPlaces; i++ {
		multiplier.Mul(multiplier, big.NewFloat(10))
	}

	// Multiply and convert to integer
	result := new(big.Float).Mul(f, multiplier)
	intValue, _ := result.Int(nil)

	// Convert integer to string
	return intValue.String()
}

type Floor struct {
	Data struct {
		TokenData struct {
			MinimumPrice string `json:"minPrice"`
		} `json:"tokenData"`
	} `json:"data"`
}

func GetTokenBestPrice(tokenAddress string, criteriaJson []map[string]interface{}, discountPerc float64) (string, string, string) {
	url := "https://marketplace-graphql.skymavis.com/graphql"
	// Convert criteriaJson to JSON string
	criteriaBytes, err := json.Marshal(criteriaJson)
	if err != nil {
		fmt.Println("Error marshaling criteriaJson:", err)
	}
	criteriaJsonStr := string(criteriaBytes)

	requestBody := fmt.Sprintf(
		"{\"operationName\":\"GetERC721TokensList\",\"variables\":{\"from\":0,\"auctionType\":\"All\",\"size\":50,\"sort\":\"PriceAsc\",\"criteria\":%v,\"rangeCriteria\":[],\"tokenAddress\":\"%s\"},\"query\":\"query GetERC721TokensList($tokenAddress: String, $slug: String, $owner: String, $auctionType: AuctionType, $criteria: [SearchCriteria!], $from: Int!, $size: Int!, $sort: SortBy, $name: String, $priceRange: InputRange, $rangeCriteria: [RangeSearchCriteria!], $excludeAddress: String) {\\n  erc721Tokens(\\n    tokenAddress: $tokenAddress\\n    slug: $slug\\n    owner: $owner\\n    auctionType: $auctionType\\n    criteria: $criteria\\n    from: $from\\n    size: $size\\n    sort: $sort\\n    name: $name\\n    priceRange: $priceRange\\n    rangeCriteria: $rangeCriteria\\n    excludeAddress: $excludeAddress\\n  ) {\\n    total\\n    results {\\n      ...Erc721TokenBrief\\n      __typename\\n    }\\n    __typename\\n  }\\n}\\n\\nfragment Erc721TokenBrief on Erc721 {\\n  tokenAddress\\n  tokenId\\n  slug\\n  owner\\n  name\\n  order {\\n    ...OrderInfo\\n    __typename\\n  }\\n  image\\n  cdnImage\\n  video\\n  isLocked\\n  attributes\\n  traitDistribution {\\n    ...TokenTrait\\n    __typename\\n  }\\n  collectionMetadata\\n  ownerProfile {\\n    name\\n    accountId\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment OrderInfo on Order {\\n  id\\n  maker\\n  kind\\n  assets {\\n    ...AssetInfo\\n    __typename\\n  }\\n  expiredAt\\n  paymentToken\\n  startedAt\\n  basePrice\\n  expectedState\\n  nonce\\n  marketFeePercentage\\n  signature\\n  hash\\n  duration\\n  timeLeft\\n  currentPrice\\n  suggestedPrice\\n  makerProfile {\\n    ...PublicProfileBrief\\n    __typename\\n  }\\n  orderStatus\\n  orderQuantity {\\n    orderId\\n    quantity\\n    remainingQuantity\\n    availableQuantity\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment AssetInfo on Asset {\\n  erc\\n  address\\n  id\\n  quantity\\n  __typename\\n}\\n\\nfragment PublicProfileBrief on PublicProfile {\\n  accountId\\n  addresses {\\n    ...Addresses\\n    __typename\\n  }\\n  activated\\n  name\\n  __typename\\n}\\n\\nfragment Addresses on NetAddresses {\\n  ethereum\\n  ronin\\n  __typename\\n}\\n\\nfragment TokenTrait on TokenTrait {\\n  key\\n  value\\n  count\\n  percentage\\n  displayType\\n  maxValue\\n  __typename\\n}\\n\"}",
		criteriaJsonStr, tokenAddress,
	)

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
	var PriceStruct NftPriceResponseData
	// Unmarshal the response body to a Go struct
	err = json.Unmarshal(body, &PriceStruct)
	if err != nil {
		fmt.Printf("Error:%v", err)
	}
	if len(PriceStruct.Data.Erc721Tokens.Results) == 0 {
		fmt.Println("No results found.")
	}
	maker := PriceStruct.Data.Erc721Tokens.Results[0].Owner
	order := PriceStruct.Data.Erc721Tokens.Results[0].Order
	multiplier := big.NewFloat(discountPerc)
	ronAmtStr := order.BasePrice
	cronAmtstr := MultiplyNftPrice(ronAmtStr, multiplier)

	return cronAmtstr, ronAmtStr, maker
}
func SendCancellOrder(ronAmtBigDec *big.Int, ronAmtHexStr string, maker string, unixExpiration string, expectedState string, unixStart string, spendingTokenAdress string, nftNumber string, collectionAddressStr string, client *defigo.Web3Client) {
	callDataStr := "95a4ec0000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000e4f524445525f45584348414e47450000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000284b910b6640000000000000000000000000000000000000000000000000000000000000020"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000240"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000020"
	callDataStr += maker
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000000"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000180"
	callDataStr += unixExpiration
	callDataStr += spendingTokenAdress
	callDataStr += unixStart
	callDataStr += ronAmtHexStr
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000000"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000000"
	callDataStr += expectedState
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000000"
	callDataStr += "00000000000000000000000000000000000000000000000000000000000001a9"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000001"
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000001"
	callDataStr += collectionAddressStr
	callDataStr += nftNumber
	callDataStr += "0000000000000000000000000000000000000000000000000000000000000000"
	callDataStr += "00000000000000000000000000000000000000000000000000000000"
	defigo.CallWriteFunctionRawData(marketGatewayProxyContract, client, ronAmtBigDec, callDataStr)
}
func SendNFTOfferOrder(collection *CollectionValues, offerAmount string, expectedState string, ownerAddress string, unixStart string, unixExpiration string, collectionAddress string, nftNumber string, spendingTokenAddress string) {
	utils.TimeAfter(5000)
	url := "https://marketplace-graphql.skymavis.com/graphql"

	// Define the request body
	// Define your variables for the request
	nonce := 0
	timeNow := time.Now().Unix()
	startedAt := int(timeNow)
	expiredAt := int(timeNow + 24*60*60)
	kind := "Offer"
	paymentToken := spendingTokenAddress
	jsonData := "{\"types\":{\"Asset\":[{\"name\":\"erc\",\"type\":\"uint8\"},{\"name\":\"addr\",\"type\":\"address\"},{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"quantity\",\"type\":\"uint256\"}],\"Order\":[{\"name\":\"maker\",\"type\":\"address\"},{\"name\":\"kind\",\"type\":\"uint8\"},{\"name\":\"assets\",\"type\":\"Asset[]\"},{\"name\":\"expiredAt\",\"type\":\"uint256\"},{\"name\":\"paymentToken\",\"type\":\"address\"},{\"name\":\"startedAt\",\"type\":\"uint256\"},{\"name\":\"basePrice\",\"type\":\"uint256\"},{\"name\":\"endedAt\",\"type\":\"uint256\"},{\"name\":\"endedPrice\",\"type\":\"uint256\"},{\"name\":\"expectedState\",\"type\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\"},{\"name\":\"marketFeePercentage\",\"type\":\"uint256\"}],\"EIP712Domain\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"version\",\"type\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\"}]},\"domain\":{\"name\":\"MarketGateway\",\"version\":\"1\",\"chainId\":\"2020\",\"verifyingContract\":\"0x3b3adf1422f84254b7fbb0e7ca62bd0865133fe3\"},\"primaryType\":\"Order\","
	jsonData += fmt.Sprintf("\"message\":{\"maker\":\"%v\",\"kind\":\"0\",\"assets\":[{\"erc\":\"1\",\"addr\":\"%s\",\"id\":\"%s\",\"quantity\":\"0\"}],\"expiredAt\":\"%d\",\"paymentToken\":\"%s\",\"startedAt\":\"%d\",\"basePrice\":\"%s\",\"endedAt\":\"0\",\"endedPrice\":\"0\",\"expectedState\":\"%s\",\"nonce\":\"0\",\"marketFeePercentage\":\"425\"}}", strings.ToLower(collection.Wallet.Client.PublicKeyStr), collectionAddress, nftNumber, expiredAt, paymentToken, startedAt, offerAmount, expectedState)
	signature := MakeSignatureRequest(jsonData, collection.Wallet.PrivateKey)

	// Construct the request body
	requestBody := fmt.Sprintf(
		"{\"operationName\":\"CreateOrder\",\"variables\":{\"order\":{\"nonce\":%v,\"assets\":[{\"id\":\"%s\",\"address\":\"%s\",\"erc\":\"Erc721\",\"quantity\":\"0\"}],\"basePrice\":\"%s\",\"startedAt\":%v,\"expiredAt\":%v,\"kind\":\"%s\",\"expectedState\":\"%s\",\"paymentToken\":\"%s\"},\"signature\":\"%s\"},\"query\":\"mutation CreateOrder($order: InputOrder!, $signature: String!) {\\n  createOrder(order: $order, signature: $signature) {\\n    ...OrderInfo\\n    __typename\\n  }\\n}\\n\\nfragment OrderInfo on Order {\\n  id\\n  maker\\n  kind\\n  assets {\\n    ...AssetInfo\\n    __typename\\n  }\\n  expiredAt\\n  paymentToken\\n  startedAt\\n  basePrice\\n  expectedState\\n  nonce\\n  marketFeePercentage\\n  signature\\n  hash\\n  duration\\n  timeLeft\\n  currentPrice\\n  suggestedPrice\\n  makerProfile {\\n    ...PublicProfileBrief\\n    __typename\\n  }\\n  orderStatus\\n  orderQuantity {\\n    orderId\\n    quantity\\n    remainingQuantity\\n    availableQuantity\\n    __typename\\n  }\\n  __typename\\n}\\n\\nfragment AssetInfo on Asset {\\n  erc\\n  address\\n  id\\n  quantity\\n  __typename\\n}\\n\\nfragment PublicProfileBrief on PublicProfile {\\n  accountId\\n  addresses {\\n    ...Addresses\\n    __typename\\n  }\\n  activated\\n  name\\n  __typename\\n}\\n\\nfragment Addresses on NetAddresses {\\n  ethereum\\n  ronin\\n  __typename\\n}\\n\"}",
		nonce, nftNumber, collectionAddress, offerAmount, startedAt, expiredAt, kind, expectedState, paymentToken, signature,
	)
	utils.AsyncPrint("red", "sending jsonBody:\n")
	utils.AsyncPrint("red", "%s", string(requestBody))

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", collection.Wallet.AuthToken)

	// Create an HTTP client and send the request
	HTTPclient := &http.Client{}
	resp, err := HTTPclient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body for debugging
	fmt.Println("Response Body:", string(body))
	var ResponseOfferData NftOfferRequestBody
	// Unmarshal the response body to a Go struct
	err = json.Unmarshal(body, &ResponseOfferData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v", err)
		return
	}
}

type NftOfferRequestBody struct {
	OperationName string `json:"operationName"`
	Variables     struct {
		Order struct {
			Nonce  int `json:"nonce"`
			Assets []struct {
				ID       string `json:"id"`
				Address  string `json:"address"`
				Erc      string `json:"erc"`
				Quantity string `json:"quantity"`
			} `json:"assets"`
			BasePrice     string `json:"basePrice"`
			StartedAt     int    `json:"startedAt"`
			ExpiredAt     int    `json:"expiredAt"`
			Kind          string `json:"kind"`
			ExpectedState string `json:"expectedState"`
			PaymentToken  string `json:"paymentToken"`
		} `json:"order"`
		Signature string `json:"signature"`
	} `json:"variables"`
	Query string `json:"query"`
}

type NftPriceResponseData struct {
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
