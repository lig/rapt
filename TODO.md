# TODO - MVP Implementation

This document outlines the essential commands and functionality needed to achieve MVP (Minimum Viable Product) status for the Rapt project.

## 🎯 Current State Analysis

**✅ Already Implemented:**
- `rapt init` - Install CRD in cluster
- `rapt add` - Add tool definitions  
- `rapt purge` - Remove CRD and resources
- `rapt version` - Show version information

**❌ Missing for MVP:**

## 🚀 Essential MVP Commands

### 1. **`rapt list`** - List Available Tools
```bash
rapt list [--namespace <ns>]
```
**Purpose**: Show all registered tools in the cluster
**Implementation**: Query Kubernetes for Tool CRDs
**Priority**: High

### 2. **`rapt run <tool-name>`** - Execute Tools ⭐ **CRITICAL**
```bash
rapt run <tool-name> [--arg key=value] [--namespace <ns>]
```
**Purpose**: Actually execute the tools (the core functionality!)
**Implementation**: Create Kubernetes Jobs from Tool definitions
**Priority**: Critical

### 3. **`rapt delete <tool-name>`** - Remove Tools
```bash
rapt delete <tool-name> [--namespace <ns>]
```
**Purpose**: Remove tool definitions from cluster
**Implementation**: Delete Tool CRDs
**Priority**: High

### 4. **`rapt describe <tool-name>`** - Show Tool Details
```bash
rapt describe <tool-name> [--namespace <ns>]
```
**Purpose**: Show detailed information about a specific tool
**Implementation**: Get and format Tool CRD details
**Priority**: Medium

### 5. **`rapt status`** - Cluster Status
```bash
rapt status [--namespace <ns>]
```
**Purpose**: Show Rapt installation status and cluster health
**Implementation**: Check CRD existence and cluster connectivity
**Priority**: Medium

### 6. **`rapt logs <tool-name>`** - View Job Logs
```bash
rapt logs <tool-name> [--follow] [--namespace <ns>]
```
**Purpose**: Show logs from tool executions
**Implementation**: Get logs from Kubernetes Jobs
**Priority**: Medium

## 🔧 Core Implementation Requirements

### **Job Management System** (Most Critical)
- [ ] Create Kubernetes Jobs from Tool definitions
- [ ] Handle job lifecycle (creation, monitoring, cleanup)
- [ ] Support job arguments and environment variables
- [ ] Implement job status tracking
- [ ] Add job resource limits and requests
- [ ] Handle job failures and retries

### **Tool Validation**
- [ ] Validate tool definitions before creation
- [ ] Check image availability and accessibility
- [ ] Validate argument formats and requirements
- [ ] Add comprehensive input validation

### **Error Handling**
- [ ] Comprehensive error messages
- [ ] Graceful failure handling
- [ ] User-friendly error reporting
- [ ] Proper exit codes

## 📊 MVP Priority Order

1. **`rapt run`** - Core functionality (execute tools) ⭐ **CRITICAL**
2. **`rapt list`** - Basic tool management
3. **`rapt delete`** - Tool lifecycle management
4. **`rapt describe`** - Tool inspection
5. **`rapt status`** - System health check
6. **`rapt logs`** - Debugging and monitoring

## 🎯 MVP Success Criteria

The project will be considered MVP-ready when users can:

1. ✅ **Initialize** Rapt in a cluster (`rapt init`)
2. ✅ **Add** tool definitions (`rapt add`)
3. ⏳ **List** available tools (`rapt list`)
4. ⏳ **Execute** tools (`rapt run <tool-name>`) ⭐ **CRITICAL**
5. ⏳ **View** tool details (`rapt describe <tool-name>`)
6. ⏳ **Remove** tools (`rapt delete <tool-name>`)
7. ⏳ **Check** system status (`rapt status`)
8. ⏳ **View** execution logs (`rapt logs <tool-name>`)

## 🏗️ Implementation Plan

### Phase 1: Core Execution (Critical Path)
- [ ] Implement `rapt run` command
- [ ] Create job management system
- [ ] Add basic error handling
- [ ] Test with simple tools

### Phase 2: Tool Management
- [ ] Implement `rapt list` command
- [ ] Implement `rapt delete` command
- [ ] Add tool validation

### Phase 3: Monitoring & Debugging
- [ ] Implement `rapt describe` command
- [ ] Implement `rapt status` command
- [ ] Implement `rapt logs` command

### Phase 4: Polish & Testing
- [ ] Comprehensive error handling
- [ ] Input validation
- [ ] Integration tests
- [ ] Documentation updates

## 🔍 Technical Considerations

### Job Management
- Use Kubernetes Jobs API for tool execution
- Implement proper job cleanup (TTL, garbage collection)
- Handle job status and completion
- Support job arguments and environment variables

### Kubernetes Integration
- Proper RBAC permissions for job creation
- Namespace handling and isolation
- Resource limits and requests
- Security contexts

### User Experience
- Clear command syntax and help text
- Consistent error messages
- Progress indicators for long-running operations
- Tab completion support

## 📝 Notes

- The **most critical missing piece** is the **`rapt run`** command, which is the core value proposition
- Focus on getting basic tool execution working before adding advanced features
- Ensure proper error handling and user feedback throughout
- Test thoroughly with various tool types and scenarios

---

**Last Updated**: January 2025  
**Status**: Planning Phase  
**Next Milestone**: Implement `rapt run` command