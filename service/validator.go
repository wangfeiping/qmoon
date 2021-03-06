// Copyright 2018 The QOS Authors

package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/QOSGroup/qmoon/lib"
	stake_types "github.com/QOSGroup/qmoon/lib/qos/stake/types"
	"github.com/QOSGroup/qmoon/models"
	"github.com/QOSGroup/qmoon/types"
	"github.com/QOSGroup/qmoon/utils"
	qostypes "github.com/QOSGroup/qos/module/stake/types"
)

func ConvertToValidator(bv *models.Validator, latestHeight int64) *types.Validator {
	statusStr := "Active"
	if bv.Status != 0 {
		statusStr = "Inactive"
	}

	return &types.Validator{
		Name:             bv.Name,
		Identity:         bv.Identity,
		Logo:             bv.Logo,
		Website:          bv.Website,
		Details:          bv.Details,
		Owner:            bv.Owner,
		ChainID:          bv.ChainId,
		Address:          bv.Address,
		StakeAddress:     bv.StakeAddress,
		ConsPubKey:       "",
		PubKeyType:       bv.PubKeyType,
		PubKeyValue:      bv.PubKeyValue,
		VotingPower:      bv.VotingPower,
		Accum:            bv.Accum,
		Commission:       bv.Commission,
		FirstBlockHeight: bv.FirstBlockHeight,
		FirstBlockTime:   bv.FirstBlockTime,
		Status:           int8(bv.Status),
		StatusStr:        statusStr,
		InactiveCode:     qostypes.InactiveCode(bv.InactiveCode),
		InactiveTime:     bv.InactiveTime,
		InactiveHeight:   bv.InactiveHeight,
		BondHeight:       bv.BondHeight,
		PrecommitNum:     bv.PrecommitNum,
		BondedTokens:     bv.BondedTokens,
		SelfBond:         bv.SelfBond,
	}
}

// Validators 查询链所有的validator
func (n Node) Validators(height int64) (types.Validators, error) {
	// latest, err := n.LatestBlock()
	//if err != nil {
	//	return nil, err
	//}
	mvs, err := models.Validators(n.ChainID)
	if err != nil {
		return nil, err
	}

	var total int64
	var res types.Validators
	for _, v := range mvs {
		if int8(v.Status) == types.Active {
			total += v.VotingPower
		}
		vv := ConvertToValidator(v, height)
		uptimePercent, _ := models.QueryValidatorUptime(n.ChainID, v.Address, 0,1)
		vv.Uptime = "0"
		if uptimePercent != nil && len(uptimePercent)>0 {
			vv.UptimeFloat, _= strconv.ParseFloat(uptimePercent[0].Y, 64)
			vv.Uptime = uptimePercent[0].Y
		}
		// fmt.Println("before final convert ", v.Address, v.BondedTokens, v.SelfBond)
		res = append(res, *vv)
	}

	for i := 0; i < len(res); i++ {
		if res[i].Status == types.Active {
			res[i].Percent = fmt.Sprintf("%.5f", utils.Percent(uint64(res[i].VotingPower), uint64(total))*100)
		} else {
			res[i].Percent = "0"
		}
	}

	sort.Sort(res)

	return res, err
}

// RetrieveValidator 单个查询
func (n Node) RetrieveValidator(address string) (*types.Validator, error) {
	latestheight, err := n.LatestBlockHeight()
	if err != nil || latestheight == 0 {
		return nil, err
	}

	mv, err := models.ValidatorByAddress(n.ChainID, address)
	if err != nil {
		return nil, err
	}

	return ConvertToValidator(mv, latestheight), nil
}

func (n Node) RetrieveValidatorByStakingAddress(address string) (*types.Validator, error) {
	latestheight, err := n.LatestBlockHeight()
	if err != nil || latestheight == 0 {
		return nil, err
	}

	mv, err := models.ValidatorByStakeAddress(n.ChainID, address)
	if err != nil {
		return nil, err
	}

	return ConvertToValidator(mv, latestheight), nil
}

func (n Node) UpdateValidatorBlock(address string, height int64, t time.Time) error {
	mv, err := models.ValidatorByAddress(n.ChainID, address)
	if err != nil {
		mv = &models.Validator{
			Address:          address,
			FirstBlockHeight: height,
			FirstBlockTime:   t,
			PrecommitNum:     1,
		}

		if err := mv.Insert(n.ChainID); err != nil {
			return err
		}
	} else {
		if mv.FirstBlockHeight == 0 {
			mv.FirstBlockHeight = height
			mv.FirstBlockTime = t
			mv.PrecommitNum = 1
		} else {
			mv.PrecommitNum = mv.PrecommitNum + 1
		}

		if err := mv.Update(n.ChainID, "first_block_height", "first_block_time", "precommit_num"); err != nil {
			return err
		}
	}

	return nil
}

func (n Node) InactiveValidator(address string, status int, inactiveHeight int64, inactiveTime time.Time) error {
	mv, err := models.ValidatorByAddress(n.ChainID, address)
	if err == nil {
		if mv.Status != status {
			mv.Status = status
			mv.InactiveCode = int(types.Inactive)
			mv.InactiveTime = inactiveTime
			mv.InactiveHeight = inactiveHeight

			if err := mv.UpdateStatus(n.ChainID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (n Node) CloseValidator(address string, inactiveHeight int64, inactiveTime time.Time) error {
	mv, err := models.ValidatorByAddress(n.ChainID, address)
	if err == nil {
		if mv.Status != 2 {
			mv.Status = 2
			mv.InactiveCode = int(types.Inactive)
			mv.InactiveTime = inactiveTime
			mv.InactiveHeight = inactiveHeight

			if err := mv.UpdateStatus(n.ChainID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (n Node) CreateValidator(vl types.Validator) error {
	mv, err := models.ValidatorByAddress(n.ChainID, vl.Address)
	if err != nil {
		mv = &models.Validator{
			Address:        vl.Address,
			StakeAddress:   vl.StakeAddress,
			PubKeyType:     vl.PubKeyType,
			PubKeyValue:    vl.PubKeyValue,
			VotingPower:    vl.VotingPower,
			Accum:          vl.Accum,
			Status:         int(vl.Status),
			InactiveCode:   int(vl.InactiveCode),
			InactiveTime:   vl.InactiveTime,
			InactiveHeight: vl.InactiveHeight,
			BondHeight:     vl.BondHeight,
			Name:           vl.Name,
			Details:        vl.Details,
			Identity:       vl.Identity,
			Logo:           vl.Logo,
			Website:        vl.Website,
			Owner:          vl.Owner,
			Commission:     vl.Commission,
			BondedTokens:   vl.BondedTokens,
			SelfBond:       vl.SelfBond,
		}
		fmt.Println("new insert ", mv.Address, mv.Status)
		if err := mv.Insert(n.ChainID); err != nil {
			return err
		}
		if err := models.UpdateBondtime(n.ChainID); err != nil {
			return err
		}
	} else {
		if mv.StakeAddress != vl.StakeAddress {
			return types.NewValidatorAddressUnmatched(mv.StakeAddress, vl.Address)
		}
		cols := make([]string, 0)
		mv.PubKeyType = vl.PubKeyType
		cols = append(cols, "pub_key_type")

		if vl.PubKeyValue != "" {
			mv.PubKeyValue = vl.PubKeyValue
			cols = append(cols, "pub_key_value")
		}
		mv.VotingPower = vl.VotingPower
		cols = append(cols, "voting_power")
		mv.Accum = vl.Accum
		cols = append(cols, "accum")
		mv.Status = int(vl.Status)
		cols = append(cols, "status")
		mv.InactiveCode = int(vl.InactiveCode)
		cols = append(cols, "inactive_code")
		if !vl.InactiveTime.IsZero() {
			mv.InactiveTime = vl.InactiveTime
			cols = append(cols, "inactive_time")
		}
		if vl.InactiveHeight >= 0 {
			mv.InactiveHeight = vl.InactiveHeight
			cols = append(cols, "inactive_height")
		}
		if vl.BondHeight >= 0 {
			mv.BondHeight = vl.BondHeight
			cols = append(cols, "bond_height")
		}
		mv.Name = vl.Name
		cols = append(cols, "name")
		mv.Logo = vl.Logo
		cols = append(cols, "logo")
		mv.Details = vl.Details
		cols = append(cols, "details")
		mv.Identity = vl.Identity
		cols = append(cols, "identity")
		mv.Website = vl.Website
		cols = append(cols, "website")
		mv.Owner = vl.Owner
		cols = append(cols, "owner")
		mv.Commission = vl.Commission
		cols = append(cols, "commission")
		mv.BondedTokens = vl.BondedTokens
		cols = append(cols, "bonded_tokens")
		mv.SelfBond = vl.SelfBond
		cols = append(cols, "self_bond")

		if err := mv.Update(n.ChainID, cols...); err != nil {
			return err
		}
	}

	return nil
}

func (n Node) ConvertDisplayValidators(val stake_types.ValidatorDisplayInfo) (types.Validator, error) {
	bondTokens_int64, err := strconv.ParseInt(val.BondedTokens, 10, 64)
	if err != nil {
		err = types.NewInvalidTypeError("val.BondedTokens "+val.BondedTokens, "int64")
		return types.Validator{}, err
	}
	selfBond_int64, err := strconv.ParseInt(val.SelfBond, 10, 64)
	if err != nil {
		err = types.NewInvalidTypeError("val.SelfBond "+val.SelfBond, "int64")
		return types.Validator{}, err
	}

	status_int8 := types.Active
	if !strings.EqualFold(val.Status, "active") {
		status_int8 = types.Inactive
	}
	inactive_int8 := int64(0)
	if val.InactiveDesc != "" && utils.IsDigit(val.InactiveDesc) {
		inactive_int8, err = strconv.ParseInt(val.InactiveDesc, 10, 8)
		if err != nil {
			return types.Validator{}, types.NewInvalidTypeError("val.InactiveDesc "+val.InactiveDesc, "int64")
		}
	}

	hexAddress := lib.Bech32AddressToHex(val.ConsPubKey)
	percent := "0.0"
	vh, err := models.ValidatorHistoryByAddress(n.ChainID, hexAddress, 0, 0, 1)
	if err == nil && vh != nil && len(vh) > 0 {
		percent = strconv.FormatFloat(float64(vh[0].VotingPower)/float64(vh[0].TotalPower)*100, 'f', -2, 64)
	}
	uptime := float64(0)
	uptimepercent := "0.0"
	uptimePercent, _ := models.QueryValidatorUptime(n.ChainID, hexAddress, 0,1)
	if uptimePercent != nil && len(uptimePercent)>0 {
		uptime, _ = strconv.ParseFloat(uptimePercent[0].Y, 64)
		uptimepercent = uptimePercent[0].Y
	}

	vall := types.Validator{
		Name:    val.Description.Moniker,
		Details: val.Description.Details,
		Logo:    val.Description.Logo,
		Website: val.Description.Website,
		Owner:   val.Owner,
		ChainID: n.Name,
		// Address:        lib.PubkeyToBech32Address(n.Bech32PrefixConsPub(), "tendermint/PubKeyEd25519", val.ConsPubKey),
		Address:        hexAddress,
		StakeAddress:   val.OperatorAddress,
		PubKeyType:     "tendermint/PubKeyEd25519",
		PubKeyValue:    val.ConsPubKey,
		VotingPower:    bondTokens_int64,
		Status:         int8(status_int8),
		InactiveCode:   qostypes.InactiveCode(inactive_int8),
		InactiveTime:   val.InactiveTime,
		InactiveHeight: val.InactiveHeight,
		BondHeight:     val.BondHeight,
		Commission:     val.Commission.Rate,
		BondedTokens:   bondTokens_int64,
		SelfBond:       selfBond_int64,

		Percent:	percent,
		UptimeFloat: uptime,
		Uptime: uptimepercent,
	}
	if vall.Status != 0 && vall.Uptime == "" {
		vall.Uptime = "0.0"
	}
	return vall, nil
}
