package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type crypto []struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type stock []struct {
	Symbol   string  `json:"symbol"`
	Name     string  `json:"name,omitempty"`
	Price    float64 `json:"price"`
	Exchange string  `json:"exchange,omitempty"`
}

type usdrates struct {
	Rates struct {
		CAD float64 `json:"CAD"`
		HKD float64 `json:"HKD"`
		ISK float64 `json:"ISK"`
		PHP float64 `json:"PHP"`
		DKK float64 `json:"DKK"`
		HUF float64 `json:"HUF"`
		CZK float64 `json:"CZK"`
		GBP float64 `json:"GBP"`
		RON float64 `json:"RON"`
		SEK float64 `json:"SEK"`
		IDR float64 `json:"IDR"`
		INR float64 `json:"INR"`
		BRL float64 `json:"BRL"`
		RUB float64 `json:"RUB"`
		HRK float64 `json:"HRK"`
		JPY float64 `json:"JPY"`
		THB float64 `json:"THB"`
		CHF float64 `json:"CHF"`
		EUR float64 `json:"EUR"`
		MYR float64 `json:"MYR"`
		BGN float64 `json:"BGN"`
		TRY float64 `json:"TRY"`
		CNY float64 `json:"CNY"`
		NOK float64 `json:"NOK"`
		NZD float64 `json:"NZD"`
		ZAR float64 `json:"ZAR"`
		USD float64 `json:"USD"`
		MXN float64 `json:"MXN"`
		SGD float64 `json:"SGD"`
		AUD float64 `json:"AUD"`
		ILS float64 `json:"ILS"`
		KRW float64 `json:"KRW"`
		PLN float64 `json:"PLN"`
	} `json:"rates"`
	Base string `json:"base"`
	Date string `json:"date"`
}

type priceData struct {
	Rates map[string]map[string]float64
	mux   sync.Mutex
}

var latestPriceData *priceData
var latestPriceDataBackup *priceData

const stockAPIKEY = "bcdc99463e7c246c97b9a91ec764012c"

var errorHappened bool

func startUpdatePriceInterval() {
	latestPriceData = new(priceData)
	latestPriceData.Rates = make(map[string]map[string]float64)
	latestPriceData.Rates["crypto"] = make(map[string]float64)
	latestPriceData.Rates["cash"] = make(map[string]float64)
	latestPriceData.Rates["stock"] = make(map[string]float64)

	updatePrice()
	for range time.Tick(time.Minute * 10) {
		updatePrice()
	}
}

func updatePrice() {
	latestPriceData.mux.Lock()
	defer latestPriceData.mux.Unlock()

	errorHappened = false

	var wg sync.WaitGroup

	wg.Add(1)
	go getCryptoPrices(&wg)
	wg.Add(1)
	go getFiat(&wg)
	wg.Add(1)
	go getsStock(&wg)
	wg.Wait()

	convertToEUR()

	if errorHappened {
		//retriveLatestPrices()
		latestPriceData = latestPriceDataBackup
	} else {
		insertPrice()
		latestPriceDataBackup = latestPriceData
	}
}

func getCryptoPrices(wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get("https://api.binance.com/api/v3/ticker/price")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		errorHappened = true
		return
	}

	var c crypto
	err = json.Unmarshal(data, &c)

	if err != nil {
		log.Println(err)
		errorHappened = true
		return
	}

	for _, v := range c {
		if strings.HasSuffix(v.Symbol, "USDT") {
			fValue, err := strconv.ParseFloat(v.Price, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			latestPriceData.Rates["crypto"][v.Symbol[:len(v.Symbol)-4]] = fValue
		}
	}
}

func getFiat(wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get("https://api.exchangeratesapi.io/latest?base=USD")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		errorHappened = true
		return
	}

	var c usdrates
	err = json.Unmarshal(data, &c)

	if err != nil {
		log.Println(err)
		errorHappened = true
		return
	}

	v := reflect.ValueOf(c.Rates)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		latestPriceData.Rates["cash"][typeOfS.Field(i).Name] = v.Field(i).Float()
	}
}

func convertToEUR() {
	convertCryptoToEUR()
	convertStockToEUR()
	convertCashToEUR()
}

func convertCryptoToEUR() {
	eurValue := latestPriceData.Rates["cash"]["EUR"]
	for k, v := range latestPriceData.Rates["crypto"] {
		latestPriceData.Rates["crypto"][k] = v * eurValue
	}
}

func convertStockToEUR() {
	eurValue := latestPriceData.Rates["cash"]["EUR"]
	for k, v := range latestPriceData.Rates["stock"] {
		latestPriceData.Rates["stock"][k] = v * eurValue
	}
}

func convertCashToEUR() {
	eurValue := latestPriceData.Rates["cash"]["EUR"]
	for k, v := range latestPriceData.Rates["cash"] {
		if k != "EUR" {
			latestPriceData.Rates["cash"][k] = eurValue / v
		}
	}
	latestPriceData.Rates["cash"]["EUR"] = 1.0
}

func getsStock(wg *sync.WaitGroup) {
	defer wg.Done()
	urlWithRevolut := "https://financialmodelingprep.com/api/v3/quote-short/TSLA,AAPL,AAL,AMZN,MSFT,GOOGL,BA,DIS,DAL,KO,AMD,NFLX,BABA,FB,CCL,XOM,BYND,MA,T,SPCE,NCLH,V,GILD,F,UBER,RCL,MMM,NVDA,BKNG,GE,BAC,ATVI,MRO,NIO,PFE,GOOG,GRPN,AMRX,INTC,SBUX,ABEV,SHOP,GPRO,ADBE,WORK,CSCO,NKE,AXP,SNE,APA,LUV,OXY,ABBV,AMC,WFC,IVR,TWTR,MUR,BTG,AR,ANF,MFA,FIT,UAA,ABT,UAL,M,CVX,APPS,SNAP,QCOM,GM,HOME,PYPL,WMT,ZNGA,JBLU,NYMT,PINS,SM,AA,EB,ERJ,GOL,C,HLT,LXRX,PEP,SPOT,PLUG,SQ,ARNC,TEVA,PBR,CVE,CRM,GT,TRIP,GOLD,MPC,MGM,FCAU,GES,BB,MAR,EA,EBAY,BIDU,O,AL,DOCU,ZNH,MGI,GPS,TWOU,MO,EXPE,GS,LPL,CHS,HAL,ET,SPGI,ROKU,UMC,KHC,TXMD,LMT,NET,HPQ,LYFT,IMGN,ARI,TGI,JMIA,SPG,COTY,CLNY,ENLC,WM,NEE,AXL,FDX,VZ,DBX,MRK,MTDR,BSBR,SLB,AUY,MDRX,BEP,OPK,CAT,AGNC,SPWR,LEVI,XRX,DELL,HOG,CARS,TM,HMY,FSLR,TSM,JD,BBAR,AMAT,TTM,NRZ,CX,BX,FL,BLK,BMY,BBD,HPE,ALXN,IRBT,MDT,DDOG,CVS,AES,PM,FVRR,BAM,PTON,DPZ,BSX,OKE,IVZ,MU,COP,EGO,ADSK,REGI,CLDR,AEO,GME,DXC,ARR,WBA,ALLY,VRTX,NBL,NOV,COST,BRFS,HIMX,CZR,ANGI,CYH,ORCL,LULU,CIM,DVN,CMCSA,WDC,TWO,CAJ,KSS,HAS,DHC,PBF,QD,MUX,TWLO,CL,BIIB,CNDT,STNE,PSX,EPD,TME,KGC,HON,H,WU,NEM,OKTA,NBEV,AU,AKAM,ETSY,NDAQ,GFI,RL,AVGO,RY,PTEN,WPX,HMC,AMGN,CHL,LLY,SHAK,MELI,NVTA,HUBS,EL,BBY,ATUS,LTC,NMR,WDAY,QSR,FOLD,CLF,BIP,HL,PANW,W,UA,UNH,CTL,ALGN,CG,CRWD,WELL,VLO,UXIN,FWONK,AMT,SYY,FTNT,TRGP,TSN,RUN,EMR,XEC,DUK,BEN,WRK,INFY,FTI,LVS,MCO,PAA,APD,FSK,ZEN,YPF,JKS,FLR,MSI,STX,MIK,STZ,A,MTCH,OSTK,NYT,MCHP,WEN,PRU,EQIX,CBRE,NTNX,PAAS,MRVL,LRCX,JWN,ADP,MET,SYF,MOMO,SFIX,MYL,KR,ECL,CAH,PS,MFG,BHC,WWE,LB,GGB,ANSS,WUBA,VALE,GEO,HLF,KMB,FVE,HDB,ON,INTU,MNST,CFG,XLNX,PE,NOC,D,FCX,DLR,PPL,SYK,MAT,EOG,AAP,SO,CHWY,BHP,SNPS,HUYA,DHR,BMA,VIPS,KMI,HCM,GLUU,VER,PBI,SKX,GIS,ZS,BOX,FEYE,EXEL,BG,KKR,ZTO,IBN,SUPV,VIV,EW,YUM,CC,WIT,SBSW,TIF,TCOM,MPLX,GRUB,NBIX,TIGR,TPR,DLTR,SCHW,LNG,PXD,CNX,ANET,BAX,RMD,VEEV,HRL,BJ,VGR,KAR,HSIC,NKTR,PAM,CNC,NAVI,SHW,GDS,CGNX,IPG,BITA,ROK,HOLX,PHM,ATR,ESI,STAY,BAP,HTHT,TFX,MXIM,HBAN,COMM,MORN,HIG,CXO,JEF,ISBC,SSNC,MLCO,DRE,LKQ,TSU,FTV,TW,KDP,QLYS,LEN,NWSA,MMC,MAS,KT,CSX,GO,OTIS,PWR,FITB,NTRS,VRSK,CBOE,SLM,INVH,FDS,CIEN,ORLY,GPK,WING,PGR,PFPT,ZION,PCAR,NSC,QTWO,PEG,NUVA,HUN,NYCB,JBHT,KNX,PBH,PRI,NOAH,SBS,BAH,UGP,DFS,LVGO,NKLA,NVAX,TDOC,AG,WMB,SEB,FAST,BDC,FLT,ARCC,NVCR,BDX,X,SFM,Z,IBKR,CF,AVLR,LYV,ATHM,LTHM,CY,HRB,SPLK,CTAS,DD,DG,NTAP,ED,HST,ISRG,FE,UNM,NTCO,UNP,GD,HUM,BIO,GGAL,SKT,HD,AXTA,WTM,NTES,HWM,IP,BKR,IT,CABO,IDXX,ROST,BLL,MBT,UNIT,UCTT,CVNA,DRI,REAL,MDLZ,MDB,LX,ARMK,MS,SRE,ETFC,DVA,PD,EPAM,PH,CERN,BRX,FFIV,OLN,STT,ZBH,OMC,JNPR,VICI,RF,TEAM,SE,SU,SWK,SWN,ICE,MKL,TECK,TV,CROX,BWA,CTSH,SWKS,KIM,VG,TMUS,CTVA,WB,WY,ALLT,MOS,ELAN,IGT,FSLY,XEL,CARR,MPW,CRSP,YY,IGMS,ULTA,CTXS,XGN,ZBRA,IIPR,ETRN,PLAN,ANTM,TAK,MTG,TAL,EDU,CAG,EFX,BMRN,NPTN,CCI,SMAR,RAD,NLOK,TER,CDE,CWEN,GLW,TFC,VIACA,CHGG,BILI,SOGO,EIX,INCY,TGT,IRM,GNL,LSCC,WYNN,GNW,ZTS,BIPC,SMFG,VMW,RES,CHKP,PLNT,PAYX,AEP,AFL,RLGY,INFN,PDD,KTOS,TMO,CLVS,TROW,CLR,CLX,CMA,CME,CMG,EQR,EQT,PDCE,VST,LDOS,CNP,CSGP,ASML,VTR,COF,TPX,COG,BEPC,DISH,ALL,RNG,TRV,AME,AMH,TSG,CHTR,EVR,TTD,PLD,GNTX,EXC,MRNA,JKHY,APH,CDNS,VRSN,RRC,PNC,PBCT,KEYS,SEDG,LOMA,PPG,TXN,GDDY,RTX,NLY,DISCA,DISCK,FISV,HBI,FREQ,PSA,HCA,AVB,NOW,SCCO,HEI,SIRI,HES,COLD,NRG,HFC,GWPH,SGMO,LMND,AYX,AZO,HGV,NICE,WEX,NUE,DAN,STLD,AMTD,NVR,DNKN,EQNR,NWL,FHN,URBN,PSTG,BNTX,SBH,BRK.B,ZM,MCD,JPM,FOXA,JNJ,IBM,AIG,PG,AM,RACE,EQH,UPS,SLCA,REGN,OVV,BSMX,DE,ILMN,SID,ITUB,DOW,K,TTWO,USB,TD,PCG,ADM,LOW,TJX,FLEX,BF.B,IQ,ADI,QRTEA,DHI,ENIA,AOS,ITW,YUMC,WAB,FIS,BK,CI,MUFG,ODP,SQM,IAG,RH,KEY,BVN,MTD,VFC,CIG,SSSS,ALB,BUD,APPN,ARCT,ASAN,AZN,AAXN,BZUN,DXCM,FANG,DIN,ENPH,EXAS,GSK,GMED,ICPT,FROG,LOGI,MSGS,PLTR,PENN,RXT,RDFN,RKT,SNOW,SUMO,TIMB,SMG,U,VIR,VRM,WIX,XPEV,YETI?apikey="
	url := urlWithRevolut + stockAPIKEY
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		errorHappened = true
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		errorHappened = true
		return
	}

	var c stock
	err = json.Unmarshal(data, &c)

	if err != nil {
		log.Println(err)
		return
	}

	for _, stock := range c {
		latestPriceData.Rates["stock"][stock.Symbol] = stock.Price
	}
}
