package mcp

import (
	"context"

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
	addTool("get_protect_devices", "Get all devices from Unifi Protect", s.getProtectDevices, map[string]any{})
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
	
	// Events
	addTool("get_protect_events", "Get events from Unifi Protect", s.getProtectEvents, map[string]any{
		"limit":  map[string]any{"type": "integer", "description": "Number of events to retrieve (optional, default 50)"},
		"offset": map[string]any{"type": "integer", "description": "Offset for pagination (optional, default 0)"},
	})

	s.server.AddTools(tools...)
}

// GET Handlers

func (s *Server) getProtectDevices(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Tool called: get_protect_devices")

	if err := s.protectClient.Authenticate(ctx); err != nil {
		return mcp.NewToolResultErrorFromErr("Authentication failed", err), nil
	}

	devices, err := s.protectClient.GetDevices(ctx)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to get devices", err), nil
	}

	return mcp.NewToolResultJSON(map[string]interface{}{
		"devices": devices,
		"count":   len(devices),
	})
}

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

// ServeStdio starts the MCP server with stdio transport
func (s *Server) ServeStdio(ctx context.Context) error {
	s.logger.Info("Starting UniFi Protect MCP Server")
	return server.ServeStdio(s.server)
}
