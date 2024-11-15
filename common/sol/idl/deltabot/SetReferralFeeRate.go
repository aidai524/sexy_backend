// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package deltabot

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// SetReferralFeeRate is the `setReferralFeeRate` instruction.
type SetReferralFeeRate struct {
	NewReferralFeeRate *uint32

	// [0] = [WRITE] gridBotState
	//
	// [1] = [WRITE, SIGNER] user
	//
	// [2] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewSetReferralFeeRateInstructionBuilder creates a new `SetReferralFeeRate` instruction builder.
func NewSetReferralFeeRateInstructionBuilder() *SetReferralFeeRate {
	nd := &SetReferralFeeRate{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	return nd
}

// SetNewReferralFeeRate sets the "newReferralFeeRate" parameter.
func (inst *SetReferralFeeRate) SetNewReferralFeeRate(newReferralFeeRate uint32) *SetReferralFeeRate {
	inst.NewReferralFeeRate = &newReferralFeeRate
	return inst
}

// SetGridBotStateAccount sets the "gridBotState" account.
func (inst *SetReferralFeeRate) SetGridBotStateAccount(gridBotState ag_solanago.PublicKey) *SetReferralFeeRate {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(gridBotState).WRITE()
	return inst
}

// GetGridBotStateAccount gets the "gridBotState" account.
func (inst *SetReferralFeeRate) GetGridBotStateAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetUserAccount sets the "user" account.
func (inst *SetReferralFeeRate) SetUserAccount(user ag_solanago.PublicKey) *SetReferralFeeRate {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(user).WRITE().SIGNER()
	return inst
}

// GetUserAccount gets the "user" account.
func (inst *SetReferralFeeRate) GetUserAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *SetReferralFeeRate) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *SetReferralFeeRate {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *SetReferralFeeRate) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

func (inst SetReferralFeeRate) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_SetReferralFeeRate,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst SetReferralFeeRate) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SetReferralFeeRate) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.NewReferralFeeRate == nil {
			return errors.New("NewReferralFeeRate parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.GridBotState is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.User is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *SetReferralFeeRate) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("SetReferralFeeRate")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("NewReferralFeeRate", *inst.NewReferralFeeRate))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=3]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta(" gridBotState", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("         user", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(2)))
					})
				})
		})
}

func (obj SetReferralFeeRate) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `NewReferralFeeRate` param:
	err = encoder.Encode(obj.NewReferralFeeRate)
	if err != nil {
		return err
	}
	return nil
}
func (obj *SetReferralFeeRate) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `NewReferralFeeRate`:
	err = decoder.Decode(&obj.NewReferralFeeRate)
	if err != nil {
		return err
	}
	return nil
}

// NewSetReferralFeeRateInstruction declares a new SetReferralFeeRate instruction with the provided parameters and accounts.
func NewSetReferralFeeRateInstruction(
	// Parameters:
	newReferralFeeRate uint32,
	// Accounts:
	gridBotState ag_solanago.PublicKey,
	user ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *SetReferralFeeRate {
	return NewSetReferralFeeRateInstructionBuilder().
		SetNewReferralFeeRate(newReferralFeeRate).
		SetGridBotStateAccount(gridBotState).
		SetUserAccount(user).
		SetSystemProgramAccount(systemProgram)
}
