package main

import (
	"math"
	"math/big"
	"sync"

	"k2go/libs/defigo"
	"k2go/libs/mavismktgo"
	"k2go/libs/utils"
)

func CheckCancelOffers(collection *mavismktgo.CollectionValues) {
	for i := 0; i < 5; i++ {
		offerLs := mavismktgo.GetSentOffers(collection.Wallet.AuthToken, collection.Wallet.PublicKey, collection.CollectionAddress)

		for _, offers := range offerLs {
			ronAmtDecBig, _ := new(big.Int).SetString(offers.RonAmt, 10)
			ronAmtHexStr := mavismktgo.PadHexString(ronAmtDecBig.Text(16))

			expectedStateDecBig, _ := new(big.Int).SetString(offers.ExpectedStateStr, 10)
			expectedStateHexStr := mavismktgo.PadHexString(expectedStateDecBig.Text(16))

			nftNumberDecBig, _ := new(big.Int).SetString(offers.NftNumberStr, 10)
			nftNumberHexStr := mavismktgo.PadHexString(nftNumberDecBig.Text(16))

			unixStartDecBig, _ := new(big.Int).SetString(offers.UnixStartStr, 10)
			unixStartHexStr := mavismktgo.PadHexString(unixStartDecBig.Text(16))

			unixExpirationDecBig, _ := new(big.Int).SetString(offers.UnixExpirationStr, 10)
			unixExpirationHexStr := mavismktgo.PadHexString(unixExpirationDecBig.Text(16))

			ownerAddressStr := mavismktgo.PadHexString(offers.OwnerAddressStr)
			collectionAddressStr := mavismktgo.PadHexString(offers.CollectionAddressStr)
			paymentTokenAddressStr := mavismktgo.PadHexString(offers.PaymentTokenAddressStr)

			ronAmtCancell := big.NewInt(0)
			mavismktgo.SendCancellOrder(ronAmtCancell, ronAmtHexStr, ownerAddressStr, unixExpirationHexStr, expectedStateHexStr, unixStartHexStr, paymentTokenAddressStr, nftNumberHexStr, collectionAddressStr, &collection.Wallet.Client)
			utils.TimeAfter(1000)
		}
		utils.TimeAfter(20000)
	}
}

var (
	FloorPriceMap = map[string]float64{}
	OwnNodeClient *defigo.Web3Client
	getNFTLock    = sync.RWMutex{}
	Collections   = mavismktgo.CollectionMap
)

func main() {
	utils.Init(true, true) // Start TsNow and RefreshRateLimit
	utils.AsyncPrint("cyan", "Starting bot... ")
	// ownNodeUrl := "https://api.roninchain.com/rpc"
	// privateK := "Private Key"
	// OwnNodeClient = defigo.NewClient(privateK, ownNodeUrl)
	// mavismktgo.SetClient(OwnNodeClient)
	for _, collection := range Collections {
		go OfferToNFT(&collection)
		utils.TimeAfter(2 * 1000)
	}
	select {}
}

func OfferToNFT(collection *mavismktgo.CollectionValues) {
	for {
		getNFTLock.Lock()
		nftDataLs := mavismktgo.GetNFTListAll(collection.CollectionAddress, collection.CollectionCriteria)
		getNFTLock.Unlock()

		bestPriceWithDisStr, floorPriceStr, floorMaker := mavismktgo.GetTokenBestPrice(collection.CollectionAddress, collection.CollectionCriteria, collection.CollectionDiscount)
		floatFloorPrice := mavismktgo.ConvertStringtoFloat(floorPriceStr)
		if collection.FloorPrice == 0 {
			collection.FloorPrice = floatFloorPrice
		}
		utils.AsyncPrint("red", "FloorPriceAtual:%v", floatFloorPrice)
		utils.AsyncPrint("red", "FloorPriceAntigo:%v", collection.FloorPrice)
		floorDif := math.Abs((collection.FloorPrice / floatFloorPrice) - 1)
		if floorDif >= 0.05 && floorMaker != collection.Wallet.PublicKey {
			utils.AsyncPrint("cyan", "Cancelling offers because floor changed")
			CheckCancelOffers(collection)
		}
		for _, nftData := range nftDataLs {

			offerPrice := bestPriceWithDisStr
			basePriceHOStr, makerStr := mavismktgo.GetNFTTokenHOP(nftData.NftNumberStr, collection.CollectionAddress)
			floatHighestOfferPrice := mavismktgo.ConvertStringtoFloat(basePriceHOStr)
			floatWithDis := mavismktgo.ConvertStringtoFloat(bestPriceWithDisStr)

			if makerStr == collection.Wallet.PublicKey || floatHighestOfferPrice > floatFloorPrice*(collection.CollectionDiscount+0.1) || nftData.OwnerAddressStr == collection.Wallet.PublicKey {
				continue
			}
			basePricePO, nftNumberPO, tokenAddressPO, makerPO, expectedStatePO, unixStartPO, unixExpirationPO := mavismktgo.GetNFTTokenPersonalInfo(nftData.NftNumberStr, collection.CollectionAddress, collection.Wallet.PublicKey)
			floatBasePO := mavismktgo.ConvertStringtoFloat(basePricePO)
			if floatHighestOfferPrice > floatBasePO && makerPO != "" {
				ronAmtDecBig, _ := new(big.Int).SetString(basePricePO, 10)
				ronAmtHexStr := mavismktgo.PadHexString(ronAmtDecBig.Text(16))

				expectedStateDecBig, _ := new(big.Int).SetString(expectedStatePO, 10)
				expectedStateHexStr := mavismktgo.PadHexString(expectedStateDecBig.Text(16))

				nftNumberDecBig, _ := new(big.Int).SetString(nftNumberPO, 10)
				nftNumberHexStr := mavismktgo.PadHexString(nftNumberDecBig.Text(16))

				unixStartDecBig, _ := new(big.Int).SetString(unixStartPO, 10)
				unixStartHexStr := mavismktgo.PadHexString(unixStartDecBig.Text(16))

				unixExpirationDecBig, _ := new(big.Int).SetString(unixExpirationPO, 10)
				unixExpirationHexStr := mavismktgo.PadHexString(unixExpirationDecBig.Text(16))

				ownerAddressStr := mavismktgo.PadHexString(makerPO)
				collectionAddressStr := mavismktgo.PadHexString(collection.CollectionAddress)
				paymentTokenAddressStr := mavismktgo.PadHexString(tokenAddressPO)

				ronAmtCancell := big.NewInt(0)
				mavismktgo.SendCancellOrder(ronAmtCancell, ronAmtHexStr, ownerAddressStr, unixExpirationHexStr, expectedStateHexStr, unixStartHexStr, paymentTokenAddressStr, nftNumberHexStr, collectionAddressStr, &collection.Wallet.Client)
				basePriceHOStr = ""
				utils.TimeAfter(18 * 1000)
			}

			if makerStr == "" || floatHighestOfferPrice < floatWithDis {
				offerPrice = bestPriceWithDisStr
			} else {
				newOfferPrice := math.Min(floatHighestOfferPrice+0.01, floatFloorPrice*(collection.CollectionDiscount+0.1))
				newOffer := mavismktgo.FloatToStringWithoutDot(newOfferPrice)
				offerPrice = newOffer
			}
			floatOfferPrice := mavismktgo.ConvertStringtoFloat(offerPrice)
			utils.AsyncPrint("cyan", "Nova oferta para NFT:%v", nftData.NftNumberStr)
			utils.AsyncPrint("cyan", "Preço da Oferta:%v", floatOfferPrice)
			if floatOfferPrice > nftData.RonAmt*0.9 {
				utils.AsyncPrint("red", "OFFER PRICE HIGHER THAN NFT PRICE \n Our Price:%v \n NFT price: %v \n Wallet: %v", floatOfferPrice, nftData.RonAmt, collection.Wallet.PublicKey)
				break
			}
			mavismktgo.SendNFTOfferOrder(collection, offerPrice, nftData.ExpectedStateStr, nftData.OwnerAddressStr, nftData.UnixStartStr, nftData.UnixExpirationStr, collection.CollectionAddress, nftData.NftNumberStr, nftData.PaymentTokenAddressStr)
			basePriceHOStr = ""
			utils.TimeAfter(4 * 1000)
		}
		utils.AsyncPrint("yellow", "NFT OFFER PROCESSED")
		utils.TimeAfter(3 * 60 * 1000)
	}
}
