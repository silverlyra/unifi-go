package protect

// Bootstrap is a large omnibus JSON struct returned from the /api/bootstrap
// UniFi Protect route.
type Bootstrap struct {
	AccessKey  string `json:"accessKey"`
	AuthUserID string `json:"authUserId"`

	Recorder *Recorder `json:"nvr"`
	Cameras  []*Camera `json:"cameras"`

	LastUpdateID string `json:"lastUpdateId"`
}

// Recorder holds the configuration and state of a NVR (network video recorder).
type Recorder struct {
	// MAC is the hardware address of the recorder.
	MAC string `json:"mac"`

	// Name is the name of the UnfFi Protect installation.
	Name string `json:"name"`

	// Host is the primary host address of the NVR.
	Host string `json:"host"`

	// Hosts lists all listening interfaces (i.e., on different VLANs).
	Hosts []string `json:"hosts"`

	// Ports identifies listening ports of different protocols.
	Ports Ports `json:"ports"`

	LastSeen Time     `json:"lastSeen"`
	UpSince  Time     `json:"upSince"`
	Uptime   Duration `json:"uptime"`

	Retention Duration         `json:"recordingRetentionDurationMs"`
	Location  RecorderLocation `json:"locationSettings"`

	Version         string `json:"version"`
	ReleaseChannel  string `json:"releaseChannel"`
	FirmwareVersion string `json:"firmwareVersion"`
	AvailableUpdate string `json:"availableUpdate"`
	CanAutoUpdate   bool   `json:"canAutoUpdate"`

	ID               string `json:"id"`
	Type             string `json:"type"`
	IsHardware       bool   `json:"isHardware"`
	HardwareID       string `json:"hardwareId"`
	HardwareRevision string `json:"hardwareRevision"`
	HardwarePlatform string `json:"hardwarePlatform"`
	HostType         int    `json:"hostType"`
	HostShortname    string `json:"hostShortname"`
}

type Ports struct {
	UMP          int `json:"ump"`
	HTTP         int `json:"http"`
	HTTPS        int `json:"https"`
	RTSP         int `json:"rtsp"`
	RSTPS        int `json:"rstps"`
	DevicesWSS   int `json:"devicesWSS"`
	CameraEvents int `json:"cameraEvents"`
	CameraHTTPS  int `json:"cameraHttps"`
	CameraTCP    int `json:"cameraTcp"`
	LiveWS       int `json:"liveWs"`
	LiveWSS      int `json:"liveWss"`
	TCPBridge    int `json:"tcpBridge"`
	TCPStreams   int `json:"tcpStreams"`
	Playback     int `json:"playback"`
	UCore        int `json:"ucore"`
}

type RecorderLocation struct {
	Latitude  float64
	Longitude float64
	Radius    int

	IsAway              bool
	IsGeofencingEnabled bool
}

type System struct {
	CPU CPU `json:"cpu"`
}

type CPU struct {
	AverageLoad float64 `json:"averageLoad"`
	Temperature float64 `json:"temperature,omitempty"`
}

// Memory gives memory statistics for an NVR. All figures in megabytes (MB).
type Memory struct {
	Available int64 `json:"available"`
	Free      int64 `json:"free"`
	Total     int64 `json:"total"`
}

// Storage gives statistics for storage devices. All figures in bytes.
type Storage struct {
	Type        string          `json:"type"`
	Size        int64           `json:"size"`
	Used        int64           `json:"used"`
	Available   int64           `json:"available"`
	IsRecycling bool            `json:"isRecycling"`
	Devices     []StorageDevice `json:"devices"`
}

type StorageDevice struct {
	Model   string `json:"model"`
	Size    int64  `json:"size"`
	Healthy bool   `json:"healthy"`
}
