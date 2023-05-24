package opcode

type FunctionInfo struct {
	Index        int
	FrameSize    int
	Arguments    int
	ReturnValues int
	IP           int64
	Codes        []IL
}
