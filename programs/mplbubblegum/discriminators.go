package mplbubblegum

// Bubblegum is an Anchor program: each instruction discriminator is the
// first 8 bytes of sha256("global:<instruction_name>"), not the 1-byte tags
// used elsewhere in this fork's hand-written and generated packages.
var (
	Instruction_Transfer   = [8]byte{163, 52, 200, 231, 140, 3, 69, 186}
	Instruction_TransferV2 = [8]byte{119, 40, 6, 235, 234, 221, 248, 49}
)
