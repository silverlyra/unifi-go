package protect

type Camera struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`

	MAC            string `json:"mac"`
	Host           string `json:"host"`
	ConnectionHost string `json:"connectionHost"`

	ConnectedSince Time `json:"connectedSince"`
	UpSince        Time `json:"upSince"`
	LastSeen       Time `json:"lastSeen"`

	LastMotion Time `json:"lastMotion"`

	Recording      bool `json:"isRecording"`
	MicEnabled     bool `json:"isMicEnabled"`
	MotionDetected bool `json:"isMotionDetected"`

	HasSpeaker bool `json:"hasSpeaker"`

	ModelKey         string `json:"modelKey"`
	Platform         string `json:"platform"`
	HardwareRevision string `json:"hardwareRevision"`
	FirmwareVersion  string `json:"firmwareVersion"`
	FirmwareBuild    string `json:"firmwareBuild"`

	CanAdopt            bool `json:"canAdopt"`
	CanManage           bool `json:"canManage"`
	Adopted             bool `json:"isAdopted"`
	Adopting            bool `json:"isAdopting"`
	Managed             bool `json:"isManaged"`
	Updating            bool `json:"isUpdating"`
	Provisioned         bool `json:"isProvisioned"`
	Rebooting           bool `json:"isRebooting"`
	SSHEnabled          bool `json:"isSshEnabled"`
	AttemptingToConnect bool `json:"isAttemptingToConnect"`

	WiredConnection CameraWiredConnection `json:"wiredConnectionState"`

	HasWifi        bool                 `json:"hasWifi"`
	WifiConnection CameraWifiConnection `json:"wifiConnectionState"`
}

type CameraState string

const (
	CameraConnected    CameraState = "CONNECTED"
	CameraDisconnected CameraState = "DISCONNECTED"
)

type CameraChannel struct {
	ID      int    `json:"id"`
	VideoID string `json:"videoId"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`

	RTSPAlias   string `json:"rtspAlias"`
	RTSPEnabled bool   `json:"isRtspEnabled"`

	Width      int   `json:"width"`
	Height     int   `json:"height"`
	FPS        int   `json:"fps"`
	Bitrate    int   `json:"bitrate"`
	MinBitrate int   `json:"minBitrate"`
	MaxBitrate int   `json:"maxBitrate"`
	FPSValues  []int `json:"fpsValues"`
}

type CameraLED struct {
	Enabled bool `json:"isEnabled"`
}

type CameraStats struct {
	Traffic CameraTraffic `json:",inline"`
	Range   CameraRange   `json:"video"`
	Storage CameraStorage `json:"storage"`
}

type CameraTraffic struct {
	Received    int `json:"rxBytes"`
	Transmitted int `json:"txBytes"`
}

type CameraRange struct {
	RecordingStart Time `json:"recordingStart"`
	RecordingEnd   Time `json:"recordingEnd"`
	TimelapseStart Time `json:"timelapseStart"`
	TimelapseEnd   Time `json:"timelapseEnd"`
}

type CameraStorage struct {
	Used int     `json:"used"`
	Rate float64 `json:"rate"`
}

type CameraWiredConnection struct {
	PhysicalRate int `json:"phyRate"`
}

type CameraWifiConnection struct {
	Channel        int `json:"channel"`
	Frequency      int `json:"frequency"`
	PhysicalRate   int `json:"phyRate"`
	SignalQuality  int `json:"signalQuality"`
	SignalStrength int `json:"signalStrength"`
}
