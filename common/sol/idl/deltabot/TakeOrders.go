// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package deltabot

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// TakeOrders is the `takeOrders` instruction.
type TakeOrders struct {
	TakeOrderParam *TakeOrdersParam

	// [0] = [] gridBotState
	//
	// [1] = [] pair
	//
	// [2] = [] clock
	//
	// [3] = [WRITE] takerTokenAccount
	//
	// [4] = [WRITE] takerBuyTokenAccount
	//
	// [5] = [WRITE] makerGridBot
	//
	// [6] = [WRITE] makerForwardOrder
	//
	// [7] = [WRITE] makerReverseOrder
	//
	// [8] = [] globalBalanceBaseUser
	//
	// [9] = [WRITE] globalBalanceBase
	//
	// [10] = [] globalBalanceQuoteUser
	//
	// [11] = [WRITE] globalBalanceQuote
	//
	// [12] = [WRITE] protocolBalanceBaseRecord
	//
	// [13] = [WRITE] protocolBalanceQuoteRecord
	//
	// [14] = [] makerUsers
	//
	// [15] = [] takerSellLimit
	//
	// [16] = [] referralRecord
	//
	// [17] = [WRITE] referralBaseFee
	//
	// [18] = [WRITE] referralQuoteFee
	//
	// [19] = [] tokenProgram
	//
	// [20] = [] associatedTokenProgram
	//
	// [21] = [WRITE, SIGNER] user
	//
	// [22] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewTakeOrdersInstructionBuilder creates a new `TakeOrders` instruction builder.
func NewTakeOrdersInstructionBuilder() *TakeOrders {
	nd := &TakeOrders{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 23),
	}
	return nd
}

// SetTakeOrderParam sets the "takeOrderParam" parameter.
func (inst *TakeOrders) SetTakeOrderParam(takeOrderParam TakeOrdersParam) *TakeOrders {
	inst.TakeOrderParam = &takeOrderParam
	return inst
}

// SetGridBotStateAccount sets the "gridBotState" account.
func (inst *TakeOrders) SetGridBotStateAccount(gridBotState ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(gridBotState)
	return inst
}

// GetGridBotStateAccount gets the "gridBotState" account.
func (inst *TakeOrders) GetGridBotStateAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetPairAccount sets the "pair" account.
func (inst *TakeOrders) SetPairAccount(pair ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(pair)
	return inst
}

// GetPairAccount gets the "pair" account.
func (inst *TakeOrders) GetPairAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetClockAccount sets the "clock" account.
func (inst *TakeOrders) SetClockAccount(clock ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(clock)
	return inst
}

// GetClockAccount gets the "clock" account.
func (inst *TakeOrders) GetClockAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetTakerTokenAccountAccount sets the "takerTokenAccount" account.
func (inst *TakeOrders) SetTakerTokenAccountAccount(takerTokenAccount ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(takerTokenAccount).WRITE()
	return inst
}

// GetTakerTokenAccountAccount gets the "takerTokenAccount" account.
func (inst *TakeOrders) GetTakerTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetTakerBuyTokenAccountAccount sets the "takerBuyTokenAccount" account.
func (inst *TakeOrders) SetTakerBuyTokenAccountAccount(takerBuyTokenAccount ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(takerBuyTokenAccount).WRITE()
	return inst
}

// GetTakerBuyTokenAccountAccount gets the "takerBuyTokenAccount" account.
func (inst *TakeOrders) GetTakerBuyTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetMakerGridBotAccount sets the "makerGridBot" account.
func (inst *TakeOrders) SetMakerGridBotAccount(makerGridBot ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(makerGridBot).WRITE()
	return inst
}

// GetMakerGridBotAccount gets the "makerGridBot" account.
func (inst *TakeOrders) GetMakerGridBotAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetMakerForwardOrderAccount sets the "makerForwardOrder" account.
func (inst *TakeOrders) SetMakerForwardOrderAccount(makerForwardOrder ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(makerForwardOrder).WRITE()
	return inst
}

// GetMakerForwardOrderAccount gets the "makerForwardOrder" account.
func (inst *TakeOrders) GetMakerForwardOrderAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetMakerReverseOrderAccount sets the "makerReverseOrder" account.
func (inst *TakeOrders) SetMakerReverseOrderAccount(makerReverseOrder ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(makerReverseOrder).WRITE()
	return inst
}

// GetMakerReverseOrderAccount gets the "makerReverseOrder" account.
func (inst *TakeOrders) GetMakerReverseOrderAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetGlobalBalanceBaseUserAccount sets the "globalBalanceBaseUser" account.
func (inst *TakeOrders) SetGlobalBalanceBaseUserAccount(globalBalanceBaseUser ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(globalBalanceBaseUser)
	return inst
}

// GetGlobalBalanceBaseUserAccount gets the "globalBalanceBaseUser" account.
func (inst *TakeOrders) GetGlobalBalanceBaseUserAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetGlobalBalanceBaseAccount sets the "globalBalanceBase" account.
func (inst *TakeOrders) SetGlobalBalanceBaseAccount(globalBalanceBase ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(globalBalanceBase).WRITE()
	return inst
}

// GetGlobalBalanceBaseAccount gets the "globalBalanceBase" account.
func (inst *TakeOrders) GetGlobalBalanceBaseAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

// SetGlobalBalanceQuoteUserAccount sets the "globalBalanceQuoteUser" account.
func (inst *TakeOrders) SetGlobalBalanceQuoteUserAccount(globalBalanceQuoteUser ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[10] = ag_solanago.Meta(globalBalanceQuoteUser)
	return inst
}

// GetGlobalBalanceQuoteUserAccount gets the "globalBalanceQuoteUser" account.
func (inst *TakeOrders) GetGlobalBalanceQuoteUserAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(10)
}

// SetGlobalBalanceQuoteAccount sets the "globalBalanceQuote" account.
func (inst *TakeOrders) SetGlobalBalanceQuoteAccount(globalBalanceQuote ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[11] = ag_solanago.Meta(globalBalanceQuote).WRITE()
	return inst
}

// GetGlobalBalanceQuoteAccount gets the "globalBalanceQuote" account.
func (inst *TakeOrders) GetGlobalBalanceQuoteAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(11)
}

// SetProtocolBalanceBaseRecordAccount sets the "protocolBalanceBaseRecord" account.
func (inst *TakeOrders) SetProtocolBalanceBaseRecordAccount(protocolBalanceBaseRecord ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[12] = ag_solanago.Meta(protocolBalanceBaseRecord).WRITE()
	return inst
}

// GetProtocolBalanceBaseRecordAccount gets the "protocolBalanceBaseRecord" account.
func (inst *TakeOrders) GetProtocolBalanceBaseRecordAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(12)
}

// SetProtocolBalanceQuoteRecordAccount sets the "protocolBalanceQuoteRecord" account.
func (inst *TakeOrders) SetProtocolBalanceQuoteRecordAccount(protocolBalanceQuoteRecord ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[13] = ag_solanago.Meta(protocolBalanceQuoteRecord).WRITE()
	return inst
}

// GetProtocolBalanceQuoteRecordAccount gets the "protocolBalanceQuoteRecord" account.
func (inst *TakeOrders) GetProtocolBalanceQuoteRecordAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(13)
}

// SetMakerUsersAccount sets the "makerUsers" account.
func (inst *TakeOrders) SetMakerUsersAccount(makerUsers ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[14] = ag_solanago.Meta(makerUsers)
	return inst
}

// GetMakerUsersAccount gets the "makerUsers" account.
func (inst *TakeOrders) GetMakerUsersAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(14)
}

// SetTakerSellLimitAccount sets the "takerSellLimit" account.
func (inst *TakeOrders) SetTakerSellLimitAccount(takerSellLimit ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[15] = ag_solanago.Meta(takerSellLimit)
	return inst
}

// GetTakerSellLimitAccount gets the "takerSellLimit" account.
func (inst *TakeOrders) GetTakerSellLimitAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(15)
}

// SetReferralRecordAccount sets the "referralRecord" account.
func (inst *TakeOrders) SetReferralRecordAccount(referralRecord ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[16] = ag_solanago.Meta(referralRecord)
	return inst
}

// GetReferralRecordAccount gets the "referralRecord" account.
func (inst *TakeOrders) GetReferralRecordAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(16)
}

// SetReferralBaseFeeAccount sets the "referralBaseFee" account.
func (inst *TakeOrders) SetReferralBaseFeeAccount(referralBaseFee ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[17] = ag_solanago.Meta(referralBaseFee).WRITE()
	return inst
}

// GetReferralBaseFeeAccount gets the "referralBaseFee" account.
func (inst *TakeOrders) GetReferralBaseFeeAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(17)
}

// SetReferralQuoteFeeAccount sets the "referralQuoteFee" account.
func (inst *TakeOrders) SetReferralQuoteFeeAccount(referralQuoteFee ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[18] = ag_solanago.Meta(referralQuoteFee).WRITE()
	return inst
}

// GetReferralQuoteFeeAccount gets the "referralQuoteFee" account.
func (inst *TakeOrders) GetReferralQuoteFeeAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(18)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *TakeOrders) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[19] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *TakeOrders) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(19)
}

// SetAssociatedTokenProgramAccount sets the "associatedTokenProgram" account.
func (inst *TakeOrders) SetAssociatedTokenProgramAccount(associatedTokenProgram ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[20] = ag_solanago.Meta(associatedTokenProgram)
	return inst
}

// GetAssociatedTokenProgramAccount gets the "associatedTokenProgram" account.
func (inst *TakeOrders) GetAssociatedTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(20)
}

// SetUserAccount sets the "user" account.
func (inst *TakeOrders) SetUserAccount(user ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[21] = ag_solanago.Meta(user).WRITE().SIGNER()
	return inst
}

// GetUserAccount gets the "user" account.
func (inst *TakeOrders) GetUserAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(21)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *TakeOrders) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *TakeOrders {
	inst.AccountMetaSlice[22] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *TakeOrders) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(22)
}

func (inst TakeOrders) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_TakeOrders,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst TakeOrders) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *TakeOrders) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.TakeOrderParam == nil {
			return errors.New("TakeOrderParam parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.GridBotState is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Pair is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Clock is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.TakerTokenAccount is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.TakerBuyTokenAccount is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.MakerGridBot is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.MakerForwardOrder is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.MakerReverseOrder is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.GlobalBalanceBaseUser is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.GlobalBalanceBase is not set")
		}
		if inst.AccountMetaSlice[10] == nil {
			return errors.New("accounts.GlobalBalanceQuoteUser is not set")
		}
		if inst.AccountMetaSlice[11] == nil {
			return errors.New("accounts.GlobalBalanceQuote is not set")
		}
		if inst.AccountMetaSlice[12] == nil {
			return errors.New("accounts.ProtocolBalanceBaseRecord is not set")
		}
		if inst.AccountMetaSlice[13] == nil {
			return errors.New("accounts.ProtocolBalanceQuoteRecord is not set")
		}
		if inst.AccountMetaSlice[14] == nil {
			return errors.New("accounts.MakerUsers is not set")
		}
		if inst.AccountMetaSlice[15] == nil {
			return errors.New("accounts.TakerSellLimit is not set")
		}
		if inst.AccountMetaSlice[16] == nil {
			return errors.New("accounts.ReferralRecord is not set")
		}
		if inst.AccountMetaSlice[17] == nil {
			return errors.New("accounts.ReferralBaseFee is not set")
		}
		if inst.AccountMetaSlice[18] == nil {
			return errors.New("accounts.ReferralQuoteFee is not set")
		}
		if inst.AccountMetaSlice[19] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[20] == nil {
			return errors.New("accounts.AssociatedTokenProgram is not set")
		}
		if inst.AccountMetaSlice[21] == nil {
			return errors.New("accounts.User is not set")
		}
		if inst.AccountMetaSlice[22] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *TakeOrders) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("TakeOrders")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("TakeOrderParam", *inst.TakeOrderParam))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=23]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("              gridBotState", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("                      pair", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("                     clock", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("                takerToken", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("             takerBuyToken", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("              makerGridBot", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("         makerForwardOrder", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("         makerReverseOrder", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("     globalBalanceBaseUser", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("         globalBalanceBase", inst.AccountMetaSlice.Get(9)))
						accountsBranch.Child(ag_format.Meta("    globalBalanceQuoteUser", inst.AccountMetaSlice.Get(10)))
						accountsBranch.Child(ag_format.Meta("        globalBalanceQuote", inst.AccountMetaSlice.Get(11)))
						accountsBranch.Child(ag_format.Meta(" protocolBalanceBaseRecord", inst.AccountMetaSlice.Get(12)))
						accountsBranch.Child(ag_format.Meta("protocolBalanceQuoteRecord", inst.AccountMetaSlice.Get(13)))
						accountsBranch.Child(ag_format.Meta("                makerUsers", inst.AccountMetaSlice.Get(14)))
						accountsBranch.Child(ag_format.Meta("            takerSellLimit", inst.AccountMetaSlice.Get(15)))
						accountsBranch.Child(ag_format.Meta("            referralRecord", inst.AccountMetaSlice.Get(16)))
						accountsBranch.Child(ag_format.Meta("           referralBaseFee", inst.AccountMetaSlice.Get(17)))
						accountsBranch.Child(ag_format.Meta("          referralQuoteFee", inst.AccountMetaSlice.Get(18)))
						accountsBranch.Child(ag_format.Meta("              tokenProgram", inst.AccountMetaSlice.Get(19)))
						accountsBranch.Child(ag_format.Meta("    associatedTokenProgram", inst.AccountMetaSlice.Get(20)))
						accountsBranch.Child(ag_format.Meta("                      user", inst.AccountMetaSlice.Get(21)))
						accountsBranch.Child(ag_format.Meta("             systemProgram", inst.AccountMetaSlice.Get(22)))
					})
				})
		})
}

func (obj TakeOrders) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `TakeOrderParam` param:
	err = encoder.Encode(obj.TakeOrderParam)
	if err != nil {
		return err
	}
	return nil
}
func (obj *TakeOrders) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `TakeOrderParam`:
	err = decoder.Decode(&obj.TakeOrderParam)
	if err != nil {
		return err
	}
	return nil
}

// NewTakeOrdersInstruction declares a new TakeOrders instruction with the provided parameters and accounts.
func NewTakeOrdersInstruction(
	// Parameters:
	takeOrderParam TakeOrdersParam,
	// Accounts:
	gridBotState ag_solanago.PublicKey,
	pair ag_solanago.PublicKey,
	clock ag_solanago.PublicKey,
	takerTokenAccount ag_solanago.PublicKey,
	takerBuyTokenAccount ag_solanago.PublicKey,
	makerGridBot ag_solanago.PublicKey,
	makerForwardOrder ag_solanago.PublicKey,
	makerReverseOrder ag_solanago.PublicKey,
	globalBalanceBaseUser ag_solanago.PublicKey,
	globalBalanceBase ag_solanago.PublicKey,
	globalBalanceQuoteUser ag_solanago.PublicKey,
	globalBalanceQuote ag_solanago.PublicKey,
	protocolBalanceBaseRecord ag_solanago.PublicKey,
	protocolBalanceQuoteRecord ag_solanago.PublicKey,
	makerUsers ag_solanago.PublicKey,
	takerSellLimit ag_solanago.PublicKey,
	referralRecord ag_solanago.PublicKey,
	referralBaseFee ag_solanago.PublicKey,
	referralQuoteFee ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	associatedTokenProgram ag_solanago.PublicKey,
	user ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *TakeOrders {
	return NewTakeOrdersInstructionBuilder().
		SetTakeOrderParam(takeOrderParam).
		SetGridBotStateAccount(gridBotState).
		SetPairAccount(pair).
		SetClockAccount(clock).
		SetTakerTokenAccountAccount(takerTokenAccount).
		SetTakerBuyTokenAccountAccount(takerBuyTokenAccount).
		SetMakerGridBotAccount(makerGridBot).
		SetMakerForwardOrderAccount(makerForwardOrder).
		SetMakerReverseOrderAccount(makerReverseOrder).
		SetGlobalBalanceBaseUserAccount(globalBalanceBaseUser).
		SetGlobalBalanceBaseAccount(globalBalanceBase).
		SetGlobalBalanceQuoteUserAccount(globalBalanceQuoteUser).
		SetGlobalBalanceQuoteAccount(globalBalanceQuote).
		SetProtocolBalanceBaseRecordAccount(protocolBalanceBaseRecord).
		SetProtocolBalanceQuoteRecordAccount(protocolBalanceQuoteRecord).
		SetMakerUsersAccount(makerUsers).
		SetTakerSellLimitAccount(takerSellLimit).
		SetReferralRecordAccount(referralRecord).
		SetReferralBaseFeeAccount(referralBaseFee).
		SetReferralQuoteFeeAccount(referralQuoteFee).
		SetTokenProgramAccount(tokenProgram).
		SetAssociatedTokenProgramAccount(associatedTokenProgram).
		SetUserAccount(user).
		SetSystemProgramAccount(systemProgram)
}
