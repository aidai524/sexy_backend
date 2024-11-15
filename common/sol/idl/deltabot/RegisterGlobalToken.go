// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package deltabot

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// RegisterGlobalToken is the `registerGlobalToken` instruction.
type RegisterGlobalToken struct {

	// [0] = [] gridBotState
	//
	// [1] = [] globalBalanceMint
	//
	// [2] = [WRITE] globalBalanceUser
	//
	// [3] = [WRITE] globalBalance
	//
	// [4] = [] tokenProgram
	//
	// [5] = [] associatedTokenProgram
	//
	// [6] = [WRITE, SIGNER] user
	//
	// [7] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewRegisterGlobalTokenInstructionBuilder creates a new `RegisterGlobalToken` instruction builder.
func NewRegisterGlobalTokenInstructionBuilder() *RegisterGlobalToken {
	nd := &RegisterGlobalToken{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 8),
	}
	return nd
}

// SetGridBotStateAccount sets the "gridBotState" account.
func (inst *RegisterGlobalToken) SetGridBotStateAccount(gridBotState ag_solanago.PublicKey) *RegisterGlobalToken {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(gridBotState)
	return inst
}

// GetGridBotStateAccount gets the "gridBotState" account.
func (inst *RegisterGlobalToken) GetGridBotStateAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetGlobalBalanceMintAccount sets the "globalBalanceMint" account.
func (inst *RegisterGlobalToken) SetGlobalBalanceMintAccount(globalBalanceMint ag_solanago.PublicKey) *RegisterGlobalToken {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(globalBalanceMint)
	return inst
}

// GetGlobalBalanceMintAccount gets the "globalBalanceMint" account.
func (inst *RegisterGlobalToken) GetGlobalBalanceMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetGlobalBalanceUserAccount sets the "globalBalanceUser" account.
func (inst *RegisterGlobalToken) SetGlobalBalanceUserAccount(globalBalanceUser ag_solanago.PublicKey) *RegisterGlobalToken {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(globalBalanceUser).WRITE()
	return inst
}

// GetGlobalBalanceUserAccount gets the "globalBalanceUser" account.
func (inst *RegisterGlobalToken) GetGlobalBalanceUserAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetGlobalBalanceAccount sets the "globalBalance" account.
func (inst *RegisterGlobalToken) SetGlobalBalanceAccount(globalBalance ag_solanago.PublicKey) *RegisterGlobalToken {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(globalBalance).WRITE()
	return inst
}

// GetGlobalBalanceAccount gets the "globalBalance" account.
func (inst *RegisterGlobalToken) GetGlobalBalanceAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *RegisterGlobalToken) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *RegisterGlobalToken {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *RegisterGlobalToken) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetAssociatedTokenProgramAccount sets the "associatedTokenProgram" account.
func (inst *RegisterGlobalToken) SetAssociatedTokenProgramAccount(associatedTokenProgram ag_solanago.PublicKey) *RegisterGlobalToken {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(associatedTokenProgram)
	return inst
}

// GetAssociatedTokenProgramAccount gets the "associatedTokenProgram" account.
func (inst *RegisterGlobalToken) GetAssociatedTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetUserAccount sets the "user" account.
func (inst *RegisterGlobalToken) SetUserAccount(user ag_solanago.PublicKey) *RegisterGlobalToken {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(user).WRITE().SIGNER()
	return inst
}

// GetUserAccount gets the "user" account.
func (inst *RegisterGlobalToken) GetUserAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *RegisterGlobalToken) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *RegisterGlobalToken {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *RegisterGlobalToken) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

func (inst RegisterGlobalToken) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_RegisterGlobalToken,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst RegisterGlobalToken) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *RegisterGlobalToken) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.GridBotState is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.GlobalBalanceMint is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.GlobalBalanceUser is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.GlobalBalance is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.AssociatedTokenProgram is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.User is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *RegisterGlobalToken) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("RegisterGlobalToken")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=8]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("          gridBotState", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("     globalBalanceMint", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("     globalBalanceUser", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("         globalBalance", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("          tokenProgram", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("associatedTokenProgram", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("                  user", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("         systemProgram", inst.AccountMetaSlice.Get(7)))
					})
				})
		})
}

func (obj RegisterGlobalToken) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *RegisterGlobalToken) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewRegisterGlobalTokenInstruction declares a new RegisterGlobalToken instruction with the provided parameters and accounts.
func NewRegisterGlobalTokenInstruction(
	// Accounts:
	gridBotState ag_solanago.PublicKey,
	globalBalanceMint ag_solanago.PublicKey,
	globalBalanceUser ag_solanago.PublicKey,
	globalBalance ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	associatedTokenProgram ag_solanago.PublicKey,
	user ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *RegisterGlobalToken {
	return NewRegisterGlobalTokenInstructionBuilder().
		SetGridBotStateAccount(gridBotState).
		SetGlobalBalanceMintAccount(globalBalanceMint).
		SetGlobalBalanceUserAccount(globalBalanceUser).
		SetGlobalBalanceAccount(globalBalance).
		SetTokenProgramAccount(tokenProgram).
		SetAssociatedTokenProgramAccount(associatedTokenProgram).
		SetUserAccount(user).
		SetSystemProgramAccount(systemProgram)
}
