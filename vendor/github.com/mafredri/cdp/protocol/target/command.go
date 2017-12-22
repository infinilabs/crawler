// Code generated by cdpgen. DO NOT EDIT.

package target

// ActivateTargetArgs represents the arguments for ActivateTarget in the Target domain.
type ActivateTargetArgs struct {
	TargetID ID `json:"targetId"` // No description.
}

// NewActivateTargetArgs initializes ActivateTargetArgs with the required arguments.
func NewActivateTargetArgs(targetID ID) *ActivateTargetArgs {
	args := new(ActivateTargetArgs)
	args.TargetID = targetID
	return args
}

// AttachToTargetArgs represents the arguments for AttachToTarget in the Target domain.
type AttachToTargetArgs struct {
	TargetID ID `json:"targetId"` // No description.
}

// NewAttachToTargetArgs initializes AttachToTargetArgs with the required arguments.
func NewAttachToTargetArgs(targetID ID) *AttachToTargetArgs {
	args := new(AttachToTargetArgs)
	args.TargetID = targetID
	return args
}

// AttachToTargetReply represents the return values for AttachToTarget in the Target domain.
type AttachToTargetReply struct {
	SessionID SessionID `json:"sessionId"` // Id assigned to the session.
}

// CloseTargetArgs represents the arguments for CloseTarget in the Target domain.
type CloseTargetArgs struct {
	TargetID ID `json:"targetId"` // No description.
}

// NewCloseTargetArgs initializes CloseTargetArgs with the required arguments.
func NewCloseTargetArgs(targetID ID) *CloseTargetArgs {
	args := new(CloseTargetArgs)
	args.TargetID = targetID
	return args
}

// CloseTargetReply represents the return values for CloseTarget in the Target domain.
type CloseTargetReply struct {
	Success bool `json:"success"` // No description.
}

// CreateBrowserContextReply represents the return values for CreateBrowserContext in the Target domain.
type CreateBrowserContextReply struct {
	BrowserContextID BrowserContextID `json:"browserContextId"` // The id of the context created.
}

// CreateTargetArgs represents the arguments for CreateTarget in the Target domain.
type CreateTargetArgs struct {
	URL              string            `json:"url"`                        // The initial URL the page will be navigated to.
	Width            *int              `json:"width,omitempty"`            // Frame width in DIP (headless chrome only).
	Height           *int              `json:"height,omitempty"`           // Frame height in DIP (headless chrome only).
	BrowserContextID *BrowserContextID `json:"browserContextId,omitempty"` // The browser context to create the page in (headless chrome only).
	// EnableBeginFrameControl Whether BeginFrames for this target will be controlled via DevTools (headless chrome only, not supported on MacOS yet, false by default).
	//
	// Note: This property is experimental.
	EnableBeginFrameControl *bool `json:"enableBeginFrameControl,omitempty"`
}

// NewCreateTargetArgs initializes CreateTargetArgs with the required arguments.
func NewCreateTargetArgs(url string) *CreateTargetArgs {
	args := new(CreateTargetArgs)
	args.URL = url
	return args
}

// SetWidth sets the Width optional argument. Frame width in DIP (headless chrome only).
func (a *CreateTargetArgs) SetWidth(width int) *CreateTargetArgs {
	a.Width = &width
	return a
}

// SetHeight sets the Height optional argument. Frame height in DIP (headless chrome only).
func (a *CreateTargetArgs) SetHeight(height int) *CreateTargetArgs {
	a.Height = &height
	return a
}

// SetBrowserContextID sets the BrowserContextID optional argument. The browser context to create the page in (headless chrome only).
func (a *CreateTargetArgs) SetBrowserContextID(browserContextID BrowserContextID) *CreateTargetArgs {
	a.BrowserContextID = &browserContextID
	return a
}

// SetEnableBeginFrameControl sets the EnableBeginFrameControl optional argument. Whether BeginFrames for this target will be controlled via DevTools (headless chrome only, not supported on MacOS yet, false by default).
//
// Note: This property is experimental.
func (a *CreateTargetArgs) SetEnableBeginFrameControl(enableBeginFrameControl bool) *CreateTargetArgs {
	a.EnableBeginFrameControl = &enableBeginFrameControl
	return a
}

// CreateTargetReply represents the return values for CreateTarget in the Target domain.
type CreateTargetReply struct {
	TargetID ID `json:"targetId"` // The id of the page opened.
}

// DetachFromTargetArgs represents the arguments for DetachFromTarget in the Target domain.
type DetachFromTargetArgs struct {
	SessionID *SessionID `json:"sessionId,omitempty"` // Session to detach.
	// TargetID is deprecated.
	//
	// Deprecated: Deprecated.
	TargetID *ID `json:"targetId,omitempty"`
}

// NewDetachFromTargetArgs initializes DetachFromTargetArgs with the required arguments.
func NewDetachFromTargetArgs() *DetachFromTargetArgs {
	args := new(DetachFromTargetArgs)

	return args
}

// SetSessionID sets the SessionID optional argument. Session to detach.
func (a *DetachFromTargetArgs) SetSessionID(sessionID SessionID) *DetachFromTargetArgs {
	a.SessionID = &sessionID
	return a
}

// SetTargetID sets the TargetID optional argument.
//
// Deprecated: Deprecated.
func (a *DetachFromTargetArgs) SetTargetID(targetID ID) *DetachFromTargetArgs {
	a.TargetID = &targetID
	return a
}

// DisposeBrowserContextArgs represents the arguments for DisposeBrowserContext in the Target domain.
type DisposeBrowserContextArgs struct {
	BrowserContextID BrowserContextID `json:"browserContextId"` // No description.
}

// NewDisposeBrowserContextArgs initializes DisposeBrowserContextArgs with the required arguments.
func NewDisposeBrowserContextArgs(browserContextID BrowserContextID) *DisposeBrowserContextArgs {
	args := new(DisposeBrowserContextArgs)
	args.BrowserContextID = browserContextID
	return args
}

// DisposeBrowserContextReply represents the return values for DisposeBrowserContext in the Target domain.
type DisposeBrowserContextReply struct {
	Success bool `json:"success"` // No description.
}

// GetTargetInfoArgs represents the arguments for GetTargetInfo in the Target domain.
type GetTargetInfoArgs struct {
	TargetID ID `json:"targetId"` // No description.
}

// NewGetTargetInfoArgs initializes GetTargetInfoArgs with the required arguments.
func NewGetTargetInfoArgs(targetID ID) *GetTargetInfoArgs {
	args := new(GetTargetInfoArgs)
	args.TargetID = targetID
	return args
}

// GetTargetInfoReply represents the return values for GetTargetInfo in the Target domain.
type GetTargetInfoReply struct {
	TargetInfo Info `json:"targetInfo"` // No description.
}

// GetTargetsReply represents the return values for GetTargets in the Target domain.
type GetTargetsReply struct {
	TargetInfos []Info `json:"targetInfos"` // The list of targets.
}

// SendMessageToTargetArgs represents the arguments for SendMessageToTarget in the Target domain.
type SendMessageToTargetArgs struct {
	Message   string     `json:"message"`             // No description.
	SessionID *SessionID `json:"sessionId,omitempty"` // Identifier of the session.
	// TargetID is deprecated.
	//
	// Deprecated: Deprecated.
	TargetID *ID `json:"targetId,omitempty"`
}

// NewSendMessageToTargetArgs initializes SendMessageToTargetArgs with the required arguments.
func NewSendMessageToTargetArgs(message string) *SendMessageToTargetArgs {
	args := new(SendMessageToTargetArgs)
	args.Message = message
	return args
}

// SetSessionID sets the SessionID optional argument. Identifier of the session.
func (a *SendMessageToTargetArgs) SetSessionID(sessionID SessionID) *SendMessageToTargetArgs {
	a.SessionID = &sessionID
	return a
}

// SetTargetID sets the TargetID optional argument.
//
// Deprecated: Deprecated.
func (a *SendMessageToTargetArgs) SetTargetID(targetID ID) *SendMessageToTargetArgs {
	a.TargetID = &targetID
	return a
}

// SetAttachToFramesArgs represents the arguments for SetAttachToFrames in the Target domain.
type SetAttachToFramesArgs struct {
	Value bool `json:"value"` // Whether to attach to frames.
}

// NewSetAttachToFramesArgs initializes SetAttachToFramesArgs with the required arguments.
func NewSetAttachToFramesArgs(value bool) *SetAttachToFramesArgs {
	args := new(SetAttachToFramesArgs)
	args.Value = value
	return args
}

// SetAutoAttachArgs represents the arguments for SetAutoAttach in the Target domain.
type SetAutoAttachArgs struct {
	AutoAttach             bool `json:"autoAttach"`             // Whether to auto-attach to related targets.
	WaitForDebuggerOnStart bool `json:"waitForDebuggerOnStart"` // Whether to pause new targets when attaching to them. Use `Runtime.runIfWaitingForDebugger` to run paused targets.
}

// NewSetAutoAttachArgs initializes SetAutoAttachArgs with the required arguments.
func NewSetAutoAttachArgs(autoAttach bool, waitForDebuggerOnStart bool) *SetAutoAttachArgs {
	args := new(SetAutoAttachArgs)
	args.AutoAttach = autoAttach
	args.WaitForDebuggerOnStart = waitForDebuggerOnStart
	return args
}

// SetDiscoverTargetsArgs represents the arguments for SetDiscoverTargets in the Target domain.
type SetDiscoverTargetsArgs struct {
	Discover bool `json:"discover"` // Whether to discover available targets.
}

// NewSetDiscoverTargetsArgs initializes SetDiscoverTargetsArgs with the required arguments.
func NewSetDiscoverTargetsArgs(discover bool) *SetDiscoverTargetsArgs {
	args := new(SetDiscoverTargetsArgs)
	args.Discover = discover
	return args
}

// SetRemoteLocationsArgs represents the arguments for SetRemoteLocations in the Target domain.
type SetRemoteLocationsArgs struct {
	Locations []RemoteLocation `json:"locations"` // List of remote locations.
}

// NewSetRemoteLocationsArgs initializes SetRemoteLocationsArgs with the required arguments.
func NewSetRemoteLocationsArgs(locations []RemoteLocation) *SetRemoteLocationsArgs {
	args := new(SetRemoteLocationsArgs)
	args.Locations = locations
	return args
}
