package types

const (
	MCPServerToolSeparator = "::"
)

type MCPServerToolType string

func (a MCPServerToolType) String() string {
	return string(a)
}

const (
	FunctionCallerMcpTool   MCPServerToolType = "functionCaller"
	FunctionAnalyzerMcpTool MCPServerToolType = "functionAnalyzer"
)

type MCPServerCategory string

func (a MCPServerCategory) String() string {
	return string(a)
}

const (
	MCPCommunicationServerCategory MCPServerCategory = "communication"
	MCPFileServerCategory          MCPServerCategory = "file"
	MCPBlobServerCategory          MCPServerCategory = "blob"
	MCPProgrammingServerCategory   MCPServerCategory = "programming"
	MCPNabhikServerCategory        MCPServerCategory = "nabhik"
	MCPSearchServerCategory        MCPServerCategory = "search"
	MCPDataServerCategory          MCPServerCategory = "dataserver"
	MCPNabhikTaskServerCategory    MCPServerCategory = "taskserver"
)
