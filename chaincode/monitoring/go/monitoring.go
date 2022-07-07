/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//External Chaincode Server Config on main function
type serverConfig struct {
	CCID    string
	Address string
}

//Smart Contract for Monitoring Aloptama(Baku) + AlatOto(ARG&AWS)
type MonitoringContract struct {
	contractapi.Contract
}

// Aloptama declare data Alat Operasional
type Aloptama struct {
	NamaAlat     string `json:"namaalat"`
	KodeAlat     string `json:"kodealat"` // for legal notes
	Bulan        string `json:"bulan"`
	MerekAlat    string `json:"merekalat"`
	JumlahAlat   int64  `json:"jumlahalat"`
	ThnPengadaan string `json:"thnpengadaan"`
	Kondisi_B    int64  `json:"kondisi_b"`
	Kondisi_R    int64  `json:"kondisi_r"`
	Kondisi_RS   int64  `json:"kondisi_rs"`
	Keterangan   string `json:"keterangan"`
}

// AlatOto declare data Alat Otomatis
type AlatOto struct {
	NamaSite       string  `json:"namasite"`
	KodeSite       string  `json:"kodesite"` // for legal no
	Bulan          string  `json:"bulan"`
	JenisAlat      string  `json:"jenisalat"`
	LokasiSite     string  `json:"lokasisite"`
	MerekLogger    string  `json:"mereklogger"`
	MerekSensor    string  `json:"mereksensor"`
	ResSensor      float64 `json:"ressensor"`
	KapBaterai     int64   `json:"kapbaterai"`
	MerekRegulator string  `json:"merekregulator"`
	KapSolar       int64   `json:"kapsolar"`
	CorrMT         string  `json:"corrmt"`
	PrevMT         string  `json:"prevmt"`
}

//HistoryAloptamaResult structure used for returning result of history query
type HistoryAloptamaResult struct {
	Record    *Aloptama `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

//HistoryAlatOtoResult structure used for returning result of history query
type HistoryAlatOtoResult struct {
	Record    *AlatOto  `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

//QueryAloptamaResult structure used for handling result of query
type QueryAloptamaResult struct {
	Key    string `json:"key"`
	Record *Aloptama
}

//QueryAlatOtoResult structure used for handling result of query
type QueryAlatOtoResult struct {
	Key    string `json:"key"`
	Record *AlatOto
}

// InitAloptama adds a base set of testing asset to Aloptama
func (s *MonitoringContract) InitAloptama(ctx contractapi.TransactionContextInterface) error {
	assets := []Aloptama{
		{NamaAlat: "Automatic Rain Gauge (ARG)", KodeAlat: "230026062022", Bulan: "September", MerekAlat: "CASELLA", JumlahAlat: 16, ThnPengadaan: "2009-2019", Kondisi_B: 16, Kondisi_R: 0, Kondisi_RS: 0, Keterangan: ""},
	}

	for _, aloptama := range assets {
		aloptamaJSON, err := json.Marshal(aloptama)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(aloptama.KodeAlat, aloptamaJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)

		}
	}

	return nil
}

// InitAlatoto adds a base set of testing asset to Aloptama
func (s *MonitoringContract) InitAlatoto(ctx contractapi.TransactionContextInterface) error {
	assets := []AlatOto{
		{NamaSite: "ARG Banjar Irigasi", KodeSite: "150278", Bulan: "September", JenisAlat: "ARG", LokasiSite: "Desa Banjar Irigasi, Kec. Lebak Ledong, Kab. Lebak", MerekLogger: "LSI", MerekSensor: "LSI", ResSensor: 0.2, KapBaterai: 28, MerekRegulator: "Sundaya", CorrMT: "", PrevMT: ""},
	}

	for _, alatoto := range assets {
		alatotoJSON, err := json.Marshal(alatoto)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(alatoto.KodeSite, alatotoJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state: %v", err)

		}
	}

	return nil
}

// Create Aloptama Asset to world state
func (s *MonitoringContract) CreateAloptama(ctx contractapi.TransactionContextInterface, namaalat, kodealat, bulan, merekalat string, jumlahalat int64, thnpengadaan string, kondisi_b, kondisi_r, kondisi_rs int64, keterangan string) error {
	exists, err := s.AloptamaExists(ctx, kodealat)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("aloptama asset %s already exist", kodealat)
	}

	aloptama := Aloptama{
		NamaAlat:     namaalat,
		KodeAlat:     kodealat,
		Bulan:        bulan,
		MerekAlat:    merekalat,
		JumlahAlat:   jumlahalat,
		ThnPengadaan: thnpengadaan,
		Kondisi_B:    kondisi_b,
		Kondisi_R:    kondisi_r,
		Kondisi_RS:   kondisi_rs,
		Keterangan:   keterangan,
	}

	aloptamaJSON, err := json.Marshal(aloptama)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(kodealat, aloptamaJSON)
}

//CreateAlatOto create AlatOto asset to ledger
func (s *MonitoringContract) CreateAlatOto(ctx contractapi.TransactionContextInterface, namasite, kodesite, bulan, jenisalat, lokasisite, mereklogger, mereksensor string, ressensor float64, kapbaterai int64, merekregulator string, kapsolar int64, corrmt, prevmt string) error {
	exists, err := s.AlatOtoExists(ctx, kodesite)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("alat otomatis asset %s already exist", kodesite)
	}

	alatoto := AlatOto{
		NamaSite:       namasite,
		KodeSite:       kodesite,
		Bulan:          bulan,
		JenisAlat:      jenisalat,
		LokasiSite:     lokasisite,
		MerekLogger:    mereklogger,
		MerekSensor:    mereksensor,
		ResSensor:      ressensor,
		KapBaterai:     kapbaterai,
		MerekRegulator: merekregulator,
		KapSolar:       kapsolar,
		CorrMT:         corrmt,
		PrevMT:         prevmt,
	}

	alatotoJSON, err := json.Marshal(alatoto)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(kodesite, alatotoJSON)
}

//ReadAloptama return the aloptama asset stored in world state with given kodealat (id)
func (s *MonitoringContract) ReadAloptama(ctx contractapi.TransactionContextInterface, kodealat string) (*Aloptama, error) {
	aloptamaJSON, err := ctx.GetStub().GetState(kodealat)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}
	if aloptamaJSON == nil {
		return nil, fmt.Errorf("the aloptama %s asset does not exist", kodealat)
	}

	var aloptama Aloptama
	err = json.Unmarshal(aloptamaJSON, &aloptama)

	if err != nil {
		return nil, err
	}

	return &aloptama, nil
}

//ReadAlatOto return the alatoto asset stored in world state with given kodealat (id)
func (s *MonitoringContract) ReadAlatOto(ctx contractapi.TransactionContextInterface, kodesite string) (*AlatOto, error) {
	alatotoJSON, err := ctx.GetStub().GetState(kodesite)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}
	if alatotoJSON == nil {
		return nil, fmt.Errorf("the alatoto %s asset does not exist", kodesite)
	}

	var alatoto AlatOto
	err = json.Unmarshal(alatotoJSON, &alatoto)

	if err != nil {
		return nil, err
	}

	return &alatoto, nil
}

//--Check Existence in world state--//

//AloptamaExists check Aloptama existence in world state
func (s *MonitoringContract) AloptamaExists(ctx contractapi.TransactionContextInterface, kodealat string) (bool, error) {
	aloptamaJSON, err := ctx.GetStub().GetState(kodealat)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	return aloptamaJSON != nil, nil
}

//AlatOtoExists check AlatOto existence in world state
func (s *MonitoringContract) AlatOtoExists(ctx contractapi.TransactionContextInterface, kodesite string) (bool, error) {
	alatotoJSON, err := ctx.GetStub().GetState(kodesite)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	return alatotoJSON != nil, nil
}

// --Update Kondisi+Keterangan Field in Aloptama data-- //
func (s *MonitoringContract) UpdateKondisiAloptama(ctx contractapi.TransactionContextInterface, kodealat, newBulan string, newKondisi_B, newKondisi_R, newKondisi_RS int64, newKeterangan string) error {
	aloptama, err := s.ReadAloptama(ctx, kodealat)
	if err != nil {
		return err
	}

	aloptama.Bulan = newBulan
	aloptama.Kondisi_B = newKondisi_B
	aloptama.Kondisi_R = newKondisi_R
	aloptama.Kondisi_RS = newKondisi_RS
	aloptama.Keterangan = newKeterangan

	aloptamaJSON, err := json.Marshal(aloptama)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(kodealat, aloptamaJSON)
}

//--Update PM+CM Field in AlatOto--//
func (s *MonitoringContract) UpdatePMCM(ctx contractapi.TransactionContextInterface, kodesite, newBulan string, newPrevMT string, newCorrMT string) error {
	alatoto, err := s.ReadAlatOto(ctx, kodesite)
	if err != nil {
		return err
	}

	alatoto.Bulan = newBulan
	alatoto.PrevMT = newPrevMT
	alatoto.CorrMT = newCorrMT

	alatotoJSON, err := json.Marshal(alatoto)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(kodesite, alatotoJSON)
}

//--Delete Feature --//

//DeleteAloptama deletes a given aloptama asset from the world state
func (s *MonitoringContract) DeleteAloptama(ctx contractapi.TransactionContextInterface, kodealat string) error {
	exists, err := s.AloptamaExists(ctx, kodealat)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the aloptama %s asset does not exist", kodealat)
	}

	return ctx.GetStub().DelState(kodealat)
}

//DeleteAloptama deletes a given aloptama asset from the world state
func (s *MonitoringContract) DeleteAlatOto(ctx contractapi.TransactionContextInterface, kodesite string) error {
	exists, err := s.AlatOtoExists(ctx, kodesite)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the alatoto %s asset does not exist", kodesite)
	}

	return ctx.GetStub().DelState(kodesite)
}

//GetAllAloptama returns all aloptama found in world state
func (s *MonitoringContract) GetAllAloptama(ctx contractapi.TransactionContextInterface) ([]QueryAloptamaResult, error) {
	// range query with empty string for startKey and endKey does an open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryAloptamaResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var aloptama Aloptama

		err = json.Unmarshal(queryResponse.Value, &aloptama)

		if err != nil {
			return nil, err
		}

		queryaloptamaresult := QueryAloptamaResult{Key: queryResponse.Key, Record: &aloptama}
		results = append(results, queryaloptamaresult)
	}

	return results, nil
}

//GetAllAlatOto returns all alatoto asset in world state
func (s *MonitoringContract) GetAllAlatOto(ctx contractapi.TransactionContextInterface) ([]QueryAlatOtoResult, error) {
	// range query with empty string for startKey and endKey does an open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var results []QueryAlatOtoResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		var alatoto AlatOto

		err = json.Unmarshal(queryResponse.Value, &alatoto)

		if err != nil {
			return nil, err
		}

		queryalatotoresult := QueryAlatOtoResult{Key: queryResponse.Key, Record: &alatoto}
		results = append(results, queryalatotoresult)
	}

	return results, nil
}

// GetAloptamaHistory returns the chain of custody for an asset since issuance.
func (t *MonitoringContract) GetAloptamaHistory(ctx contractapi.TransactionContextInterface, kodealat string) ([]HistoryAloptamaResult, error) {
	log.Printf("GetAloptamaHistory: Kode Alat %v", kodealat)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(kodealat)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var aloptamarecords []HistoryAloptamaResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var aloptama Aloptama
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &aloptama)
			if err != nil {
				return nil, err
			}
		} else {
			aloptama = Aloptama{
				KodeAlat: kodealat,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryAloptamaResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &aloptama,
			IsDelete:  response.IsDelete,
		}
		aloptamarecords = append(aloptamarecords, record)
	}

	return aloptamarecords, nil
}

// GetAlatOtoHistory returns the chain of custody for an asset since issuance.
func (t *MonitoringContract) GetAlatOtoHistory(ctx contractapi.TransactionContextInterface, kodesite string) ([]HistoryAlatOtoResult, error) {
	log.Printf("GetAlatOtoHistory: Kode Site %v", kodesite)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(kodesite)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	var alatotorecords []HistoryAlatOtoResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var alatoto AlatOto
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &alatoto)
			if err != nil {
				return nil, err
			}
		} else {
			alatoto = AlatOto{
				KodeSite: kodesite,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryAlatOtoResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &alatoto,
			IsDelete:  response.IsDelete,
		}
		alatotorecords = append(alatotorecords, record)
	}

	return alatotorecords, nil
}

func main() {
	// See chaincode.env.example
	config := serverConfig{
		CCID:    os.Getenv("CHAINCODE_ID"),
		Address: os.Getenv("CHAINCODE_SERVER_ADDRESS"),
	}

	chaincode, err := contractapi.NewChaincode(new(MonitoringContract))

	if err != nil {
		log.Panicf("Error create Aloptama+AlatOto Chaincode: %s", err)
		return
	}

	server := &shim.ChaincodeServer{
		CCID:    config.CCID,
		Address: config.Address,
		CC:      chaincode,
		TLSProps: shim.TLSProperties{
			Disabled: true,
		},
	}

	if err := server.Start(); err != nil {
		log.Panicf("Error starting chaincodes: %s", err)
	}

}
