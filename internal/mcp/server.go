package mcp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
	"github.com/surrealwolf/unifi-protect-mcp/internal/unifi"
)

// Server represents the MCP server
type Server struct {
	protectClient *unifi.ProtectClient
	server        *server.MCPServer
	logger        *logrus.Entry
}

// NewServer creates a new MCP server
func NewServer(protectClient *unifi.ProtectClient) *Server {
	s := &Server{
		protectClient: protectClient,
		server:        server.NewMCPServer("unifi-protect-mcp", "0.1.0"),
		logger:        logrus.WithField("component", "MCPServer"),
	}

	s.registerTools()
	return s
}

func (s *Server) registerTools() {
	tools := []server.ServerTool{}

	// Helper to create tool definitions
	addTool := func(name, desc string, handler server.ToolHandlerFunc, properties map[string]any) {
		tools = append(tools, server.ServerTool{
			Tool: mcp.Tool{
				Name:        name,
				Description: desc,
				InputSchema: mcp.ToolInputSchema{
					Type:       "object",
					Properties: properties,
				},
			},
			Handler: handler,
		})
	}

	// Device and System queries
	addTool("get_protect_cameras", "Get all cameras from Unifi Protect", s.getProtectCameras, map[string]any{})
	addTool("get_protect_sensors", "Get all sensors from Unifi Protect", s.getProtectSensors, map[string]any{})
	addTool("get_protect_lights", "Get all lights from Unifi Protect", s.getProtectLights, map[string]any{})
	addTool("get_protect_chimes", "Get all chimes from Unifi Protect", s.getProtectChimes, map[string]any{})
	addTool("get_protect_liveviews", "Get all live views from Unifi Protect", s.getProtectLiveviews, map[string]any{})

	// Detailed resource information
	addTool("get_camera_detailed", "Get detailed information about a specific camera", s.getCameraDetailed, map[string]any{
		"camera_id": map[string]any{"type": "string", "description": "Camera ID (required)"},
	})
	addTool("get_sensor_detailed", "Get detailed information about a specific sensor", s.getSensorDetailed, map[string]any{
		"sensor_id": map[string]any{"type": "string", "description": "Sensor ID (required)"},
	})
	addTool("get_light_detailed", "Get detailed information about a specific light", s.getLightDetailed, map[string]any{
		"light_id": map[string]any{"type": "string", "description": "Light ID (required)"},
	})
	addTool("get_chime_detailed", "Get detailed information about a specific chime", s.getChimeDetailed, map[string]any{
		"chime_id": map[string]any{"type": "string", "description": "Chime ID (required)"},
	})
	addTool("get_liveview_detailed", "Get detailed information about a specific live view", s.getLiveviewDetailed, map[string]any{
		"liveview_id": map[string]any{"type": "string", "description": "Live view ID (required)"},
	})

	// System and Configuration
	addTool("get_protect_info", "Get system information from Unifi Protect", s.getProtectInfo, map[string]any{})
	addTool("get_protect_nvr", "Get NVR information from Unifi Protect", s.getProtectNVR, map[string]any{})
	addTool("get_protect_viewers", "Get all viewers from Unifi Protect", s.getProtectViewers, map[string]any{})
	addTool("get_protect_viewer_detailed", "Get detailed information about a specific viewer", s.getProtectViewerDetailed, map[string]any{
		"id": map[string]any{"type": "string", "description": "Viewer ID"},
	})

	// Modify Resources
	addTool("patch_protect_viewer", "Update viewer settings", s.patchProtectViewer, map[string]any{
		"id":       map[string]any{"type": "string", "description": "Viewer ID"},
		"settings": map[string]any{"type": "object", "description": "Viewer settings to update"},
	})

	// Camera Controls
	addTool("camera_start_ptz_patrol", "Start a PTZ patrol on a camera", s.cameraStartPTZPatrol, map[string]any{
		"camera_id": map[string]any{"type": "string", "description": "Camera ID"},
		"slot":      map[string]any{"type": "integer", "description": "Patrol slot number"},
	})
	addTool("camera_stop_ptz_patrol", "Stop a PTZ patrol on a camera", s.cameraStopPTZPatrol, map[string]any{
		"camera_id": map[string]any{"type": "string", "description": "Camera ID"},
	})
	addTool("camera_goto_ptz_preset", "Move camera to a PTZ preset position", s.cameraGotoPTZPreset, map[string]any{
		"camera_id": map[string]any{"type": "string", "description": "Camera ID"},
		"slot":      map[string]any{"type": "integer", "description": "Preset slot number"},
	})
	addTool("camera_create_rtsps_stream", "Create an RTSPS stream for a camera", s.cameraCreateRTSPSStream, map[string]any{
		"camera_id": map[string]any{"type": "string", "description": "Camera ID"},
		"config":    map[string]any{"type": "object", "description": "RTSPS stream configuration"},
	})
	addTool("camera_create_talkback_session", "Create a talkback session with a camera", s.cameraCreateTalkbackSession, map[string]any{
		"camera_id": map[string]any{"type": "string", "description": "Camera ID"},
		"config":    map[string]any{"type": "object", "description": "Talkback session configuration"},
	})
	addTool("camera_disable_mic_permanently", "Disable microphone permanently on a camera", s.cameraDisableMicPermanently, map[string]any{
		"camera_id": map[string]any{"type": "string", "description": "Camera ID"},
	})
	addTool("trigger_webhook_alarm", "Trigger a configured alarm webhook", s.triggerWebhookAlarm, map[string]any{
		"webhook_id": map[string]any{"type": "string", "description": "Webhook ID"},
		"payload":    map[string]any{"type": "object", "description": "Alarm trigger payload (optional)"},
	})

	// Events
	addTool("get_protect_events", "Get events from Unifi Protect", s.getProtectEvents, map[string]any{
		"limit":  map[string]any{"type": "integer", "description": "Number of events to retrieve (optional, default 50)"},
		"offset": map[string]any{"type": "integer", "description": "Offset for pagination (optional, default 0)"},
	})

	s.server.AddTools(tools...)
}

// GET Handlers

func (s *Server) getProtectCameras(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_cameras")

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	cameras, err := s.protectClient.GetCameras(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get cameras", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"cameras": cameras,
		"count":   len(cameras),
	})
}

func (s *Server) getProtectSensors(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_sensors")

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	sensors, err := s.protectClient.GetSensors(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get sensors", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"sensors": sensors,
		"count":   len(sensors),
	})
}

func (s *Server) getProtectLights(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_lights")

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	lights, err := s.protectClient.GetLights(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get lights", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"lights": lights,
		"count":  len(lights),
	})
}

func (s *Server) getProtectChimes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_chimes")

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	chimes, err := s.protectClient.GetChimes(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get chimes", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"chimes": chimes,
		"count":  len(chimes),
	})
}

func (s *Server) getProtectLiveviews(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_liveviews")

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	liveviews, err := s.protectClient.GetLiveviews(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get liveviews", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"liveviews": liveviews,
		"count":     len(liveviews),
	})
}

func (s *Server) getCameraDetailed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_camera_detailed")

	cameraID := request.GetString("camera_id", "")
	if cameraID == "" {
		return mcp.NewToolResultError("camera_id is required"), nil
	}

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	camera, err := s.protectClient.GetCameraDetailed(ctx, cameraID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get camera details", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"camera":    camera,
		"camera_id": cameraID,
	})
}

func (s *Server) getSensorDetailed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_sensor_detailed")

	sensorID := request.GetString("sensor_id", "")
	if sensorID == "" {
		return mcp.NewToolResultError("sensor_id is required"), nil
	}

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	sensor, err := s.protectClient.GetSensorDetailed(ctx, sensorID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get sensor details", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"sensor":    sensor,
		"sensor_id": sensorID,
	})
}

func (s *Server) getLightDetailed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_light_detailed")

	lightID := request.GetString("light_id", "")
	if lightID == "" {
		return mcp.NewToolResultError("light_id is required"), nil
	}

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	light, err := s.protectClient.GetLightDetailed(ctx, lightID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get light details", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"light":    light,
		"light_id": lightID,
	})
}

func (s *Server) getChimeDetailed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_chime_detailed")

	chimeID := request.GetString("chime_id", "")
	if chimeID == "" {
		return mcp.NewToolResultError("chime_id is required"), nil
	}

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	chime, err := s.protectClient.GetChimeDetailed(ctx, chimeID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get chime details", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"chime":    chime,
		"chime_id": chimeID,
	})
}

func (s *Server) getLiveviewDetailed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_liveview_detailed")

	liveviewID := request.GetString("liveview_id", "")
	if liveviewID == "" {
		return mcp.NewToolResultError("liveview_id is required"), nil
	}

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	liveview, err := s.protectClient.GetLiveviewDetailed(ctx, liveviewID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get liveview details", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"liveview":    liveview,
		"liveview_id": liveviewID,
	})
}

func (s *Server) getProtectEvents(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_events")

	limit := request.GetInt("limit", 50)
	offset := request.GetInt("offset", 0)

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	events, err := s.protectClient.GetEvents(ctx, limit, offset)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get events", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"events": events,
		"count":  len(events),
		"limit":  limit,
		"offset": offset,
	})
}

func (s *Server) getProtectInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_info")

	if err := s.protectClient.Authenticate(ctx); err != nil {
		s.logger.WithError(err).Error("Failed to authenticate with Protect")
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	info, err := s.protectClient.GetSystemInfo(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get system info")
		return mcp.NewToolResultErrorFromErr("Failed to get system info", err), nil
	}

	result := map[string]interface{}{
		"version":             info.Version,
		"application_version": info.ApplicationVersion,
		"unique_id":           info.UniqueID,
		"system_type":         info.SystemType,
	}

	return mcp.NewToolResultJSON(result)
}

func (s *Server) getProtectNVR(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_nvr")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	nvr, err := s.protectClient.GetNVR(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get NVR information", err), nil
	}
	return mcp.NewToolResultJSON(nvr)
}

func (s *Server) getProtectViewers(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_viewers")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	viewers, err := s.protectClient.GetViewers(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get viewers", err), nil
	}
	result := map[string]interface{}{
		"viewers": viewers,
		"count":   len(viewers),
	}
	return mcp.NewToolResultJSON(result)
}

func (s *Server) getProtectViewerDetailed(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_viewer_detailed")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	viewerID := request.GetString("id", "")
	if viewerID == "" {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: id", nil), nil
	}
	viewer, err := s.protectClient.GetViewerDetailed(ctx, viewerID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get viewer details", err), nil
	}
	return mcp.NewToolResultJSON(viewer)
}

func (s *Server) patchProtectViewer(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: patch_protect_viewer")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	viewerID := request.GetString("id", "")
	if viewerID == "" {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: id", nil), nil
	}
	args := request.GetArguments()
	settings, ok := args["settings"].(map[string]interface{})
	if !ok || len(settings) == 0 {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: settings", nil), nil
	}
	viewer, err := s.protectClient.PatchViewer(ctx, viewerID, settings)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to update viewer", err), nil
	}
	return mcp.NewToolResultJSON(viewer)
}

func (s *Server) cameraStartPTZPatrol(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: camera_start_ptz_patrol")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	cameraID := request.GetString("camera_id", "")
	if cameraID == "" {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: camera_id", nil), nil
	}
	slot := request.GetInt("slot", 0)
	if slot < 0 {
		return mcp.NewToolResultErrorFromErr("Invalid slot number", nil), nil
	}
	result, err := s.protectClient.CameraStartPTZPatrol(ctx, cameraID, slot)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to start PTZ patrol", err), nil
	}
	return mcp.NewToolResultJSON(result)
}

func (s *Server) cameraStopPTZPatrol(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: camera_stop_ptz_patrol")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	cameraID := request.GetString("camera_id", "")
	if cameraID == "" {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: camera_id", nil), nil
	}
	result, err := s.protectClient.CameraStopPTZPatrol(ctx, cameraID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to stop PTZ patrol", err), nil
	}
	return mcp.NewToolResultJSON(result)
}

func (s *Server) cameraGotoPTZPreset(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: camera_goto_ptz_preset")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	cameraID := request.GetString("camera_id", "")
	if cameraID == "" {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: camera_id", nil), nil
	}
	slot := request.GetInt("slot", 0)
	if slot < 0 {
		return mcp.NewToolResultErrorFromErr("Invalid slot number", nil), nil
	}
	result, err := s.protectClient.CameraGotoPTZPreset(ctx, cameraID, slot)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to move to PTZ preset", err), nil
	}
	return mcp.NewToolResultJSON(result)
}

func (s *Server) cameraCreateRTSPSStream(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: camera_create_rtsps_stream")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	cameraID := request.GetString("camera_id", "")
	if cameraID == "" {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: camera_id", nil), nil
	}
	args := request.GetArguments()
	config, ok := args["config"].(map[string]interface{})
	if !ok {
		config = map[string]interface{}{}
	}
	stream, err := s.protectClient.CameraCreateRTSPSStream(ctx, cameraID, config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to create RTSPS stream", err), nil
	}
	return mcp.NewToolResultJSON(stream)
}

func (s *Server) cameraCreateTalkbackSession(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: camera_create_talkback_session")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	cameraID := request.GetString("camera_id", "")
	if cameraID == "" {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: camera_id", nil), nil
	}
	args := request.GetArguments()
	config, ok := args["config"].(map[string]interface{})
	if !ok {
		config = map[string]interface{}{}
	}
	session, err := s.protectClient.CameraCreateTalkbackSession(ctx, cameraID, config)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to create talkback session", err), nil
	}
	return mcp.NewToolResultJSON(session)
}

func (s *Server) cameraDisableMicPermanently(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: camera_disable_mic_permanently")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	cameraID := request.GetString("camera_id", "")
	if cameraID == "" {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: camera_id", nil), nil
	}
	result, err := s.protectClient.CameraDisableMicPermanently(ctx, cameraID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to disable microphone", err), nil
	}
	return mcp.NewToolResultJSON(result)
}

func (s *Server) triggerWebhookAlarm(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: trigger_webhook_alarm")
	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}
	webhookID := request.GetString("webhook_id", "")
	if webhookID == "" {
		return mcp.NewToolResultErrorFromErr("Missing required parameter: webhook_id", nil), nil
	}
	args := request.GetArguments()
	payload, ok := args["payload"].(map[string]interface{})
	if !ok {
		payload = map[string]interface{}{}
	}
	result, err := s.protectClient.TriggerWebhookAlarm(ctx, webhookID, payload)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to trigger webhook alarm", err), nil
	}
	return mcp.NewToolResultJSON(result)
}

// ServeStdio starts the MCP server with stdio transport
func (s *Server) ServeStdio(ctx context.Context) error {
	s.logger.Info("Starting UniFi Protect MCP Server")
	return server.ServeStdio(s.server)
}

// ServeHTTP starts the MCP server with HTTP transport
func (s *Server) ServeHTTP(addr string, ctx context.Context) error {
	s.logger.Infof("Starting UniFi Protect MCP Server on HTTP at %s", addr)

	http.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Parse the MCP request
		var requestData map[string]interface{}
		if err := json.Unmarshal(body, &requestData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Log the request
		s.logger.Debugf("HTTP MCP request received: %v", requestData)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]string{
			"status": "MCP HTTP transport is available",
			"info":   "This is an HTTP endpoint. Use stdio transport for full MCP protocol support.",
		}
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	return http.ListenAndServe(addr, nil)
}
