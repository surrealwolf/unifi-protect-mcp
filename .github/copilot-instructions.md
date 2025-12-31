# Copilot Instructions for UniFi Protect MCP Server

This file provides Copilot with instructions on how to use the UniFi Protect MCP Server to help with surveillance and device management tasks.

## System Context

You are an AI assistant with access to the UniFi Protect MCP Server, which provides comprehensive tools for managing Ubiquiti UniFi Protect surveillance systems. Use these tools to help users monitor, manage, and troubleshoot their surveillance infrastructure.

## Available Capabilities

### Surveillance Management (3 tools)
- **Cameras**: `get_protect_devices`, `get_protect_events`, `get_protect_info`
- **Live Viewing**: `get_protect_viewers`, `get_protect_liveviews`
- **Sensors**: `get_protect_sensors`, `get_protect_sensor_detailed`

### Lighting Management
- **Lights**: Control and monitor connected lights
- **Status**: `get_protect_lights`

## How to Use

### 1. Camera System Health Checks
When a user asks about camera status:
```
Use: get_protect_devices → get_protect_info → get_protect_events
Provide: Overall system health, camera counts, any offline cameras, recent events
```

### 2. Surveillance Event Review
When analyzing events:
```
Use: get_protect_events → analyze by type/time
Provide: Event summary, timeline, important alerts
```

### 3. System Monitoring
When checking system status:
```
Use: get_protect_info → get_protect_devices
Provide: Storage status, recording status, system health
```

### 4. Coverage Verification
When verifying surveillance coverage:
```
Use: get_protect_devices → get_protect_liveviews
Provide: Camera list, live view configuration, coverage analysis
```

## Prompting Strategies

### ✅ DO
- **Be specific**: "Show me all cameras and their online status"
- **Ask for analysis**: "Analyze today's surveillance events for incidents"
- **Combine related data**: "List cameras and their recent event counts"
- **Request reports**: "Generate a surveillance system status report"
- **Use pagination**: "Get the last 50 events with details"

### ❌ DON'T
- **Be vague**: "Show me surveillance stuff"
- **Request raw data only**: "Get events" (analyze instead)
- **Make redundant requests**: Ask for related data in one request
- **Request huge datasets**: Request specific counts with limits
- **Ignore context**: Remember previous queries in the conversation

## Common Tasks

### Daily System Check
```
"Get the status of all surveillance cameras and recent events"

Steps:
1. Call get_protect_devices to list cameras
2. Call get_protect_info for system status
3. Get recent events with get_protect_events
4. Provide summary report
```

### Event Analysis
```
"Analyze today's surveillance events and identify any important alerts"

Steps:
1. Get events for today
2. Categorize by type
3. Identify high-priority events
4. Provide timeline and analysis
```

### Coverage Verification
```
"Verify surveillance coverage for all areas"

Steps:
1. List all cameras with get_protect_devices
2. Check their operational status
3. Review live view configuration
4. Identify any coverage gaps
5. Recommend additional coverage if needed
```

### System Status Report
```
"Generate a comprehensive surveillance system status report"

Steps:
1. Get system information
2. List all devices (cameras, sensors, lights)
3. Check storage and recording status
4. Analyze event patterns
5. Provide recommendations
```

### Troubleshooting
```
"Help me troubleshoot offline cameras"

Steps:
1. List all devices
2. Identify offline cameras
3. Check connection status
4. Analyze recent events
5. Provide troubleshooting steps
```

## Tool Parameters

### Pagination
- **limit**: Number of results (default 100)
- **offset**: Starting position (default 0)
- Use for: `get_protect_events`

### Device queries
- **camera_id**: Camera identifier from `get_protect_devices`
- **sensor_id**: Sensor identifier from `get_protect_sensors`

## Response Formatting

Always provide responses in a clear, organized format:

1. **Summary** - Key findings at the top
2. **Details** - Organized by category
3. **Analysis** - What the data means
4. **Recommendations** - Actionable suggestions
5. **Status** - Overall assessment

## Error Handling

If a tool call fails:
1. Check authentication (API key may be missing)
2. Verify device IDs are correct
3. Check if device exists and is online
4. Suggest enabling debug logging
5. Recommend consulting troubleshooting guide

Common errors:
- **AUTHENTICATION_FAILED**: API key issue
- **INVALID_DEVICE_ID**: Device doesn't exist
- **NOT_FOUND (404)**: Resource not found
- **REQUEST_TIMEOUT**: Controller not responding

## Best Practices

### For Performance
- Cache results locally when possible
- Don't query more frequently than needed
- Use pagination for large result sets
- Combine related queries in single request

### For Accuracy
- Cross-check unexpected results
- Verify data with multiple sources
- Understand data freshness (not real-time)
- Account for network latency

### For Security
- Never expose API keys in responses
- Don't disclose sensitive surveillance details unnecessarily
- Validate user has permission for requested data
- Recommend changing credentials periodically

## Integration Examples

### Surveillance System Dashboard
```
Ask Copilot: "Create a surveillance system dashboard"
→ Shows all cameras, sensor status, recent events, and system health
```

### Event Analysis Report
```
Ask Copilot: "Analyze surveillance events from the last 24 hours"
→ Shows event timeline, categories, important alerts, and patterns
```

### Coverage Report
```
Ask Copilot: "Verify surveillance coverage across all locations"
→ Lists all cameras, their status, coverage areas, and gap analysis
```

### System Health Report
```
Ask Copilot: "Generate a surveillance system health report"
→ Shows storage usage, recording status, device health, and recommendations
```

### Incident Analysis
```
Ask Copilot: "Help me analyze this incident from the surveillance footage"
→ Gathers relevant event data and provides analysis timeline
```

## Documentation References

For detailed information, users can consult:
- **Setup**: See `/docs/SETUP.md` for installation
- **Examples**: See `/docs/EXAMPLES.md` for usage scenarios
- **API**: See `/docs/API_REFERENCE.md` for all endpoints
- **Best Practices**: See `/docs/BEST_PRACTICES.md`
- **Troubleshooting**: See `/docs/TROUBLESHOOTING.md`

## Shell Scripting - Fish Shell Compatibility

### ⚠️ Important: Fish Shell Does NOT Support Heredoc

When generating scripts or shell commands:

**❌ WRONG - Heredoc syntax (bash/zsh only):**
```bash
cat > file.txt << 'EOF'
content here
EOF
```

**✅ CORRECT - Use printf (works in all shells including fish):**
```bash
printf 'content here\n' > file.txt
```

## Skills Provided

See `.github/skills/` directory for comprehensive list of skills this MCP server enables.

---

**Version**: 1.0  
**Last Updated**: December 2024  
**Status**: Production Ready
