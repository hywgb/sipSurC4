package router

import (
	"context"
	"errors"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrNoAvailableAgent = errors.New("no available agent")
	ErrInvalidSkill     = errors.New("invalid skill requirement")
)

// RouteRequest 路由请求
type RouteRequest struct {
	CallID       uuid.UUID
	CustomerID   uuid.UUID
	RequiredSkills []string
	Priority     int
	Language     string
	VIP          bool
	Metadata     map[string]interface{}
}

// RouteResult 路由结果
type RouteResult struct {
	AgentID      uuid.UUID
	SkillMatched []string
	Score        float64
	WaitTime     time.Duration
}

// Agent 座席信息（用于路由）
type Agent struct {
	ID          uuid.UUID
	Status      AgentStatus
	Skills      []string
	Level       int
	CurrentLoad int
	Performance float64
	Languages   []string
}

// AgentStatus 座席状态
type AgentStatus string

const (
	AgentStatusAvailable AgentStatus = "available"
	AgentStatusBusy      AgentStatus = "busy"
	AgentStatusBreak     AgentStatus = "break"
	AgentStatusOffline   AgentStatus = "offline"
)

// Router 路由器接口
type Router interface {
	Route(ctx context.Context, req *RouteRequest) (*RouteResult, error)
	RegisterAgent(agent *Agent) error
	UnregisterAgent(agentID uuid.UUID) error
	UpdateAgentStatus(agentID uuid.UUID, status AgentStatus) error
	GetRouterStats() *RouterStats
}

// RouterStats 路由器统计
type RouterStats struct {
	TotalAgents     int
	AvailableAgents int
	BusyAgents      int
	TotalRouted     int64
	AvgWaitTime     time.Duration
}

// SkillBasedRouter 基于技能的路由器
type SkillBasedRouter struct {
	mu              sync.RWMutex
	agents          map[uuid.UUID]*Agent
	skillIndex      map[string][]uuid.UUID // skill -> agent IDs
	availableAgents map[uuid.UUID]bool
	stats           *RouterStats
	logger          *logrus.Logger
}

// NewRouter 创建路由器
func NewRouter() Router {
	return &SkillBasedRouter{
		agents:          make(map[uuid.UUID]*Agent),
		skillIndex:      make(map[string][]uuid.UUID),
		availableAgents: make(map[uuid.UUID]bool),
		stats:           &RouterStats{},
		logger:          logrus.New(),
	}
}

// Route 执行路由
func (r *SkillBasedRouter) Route(ctx context.Context, req *RouteRequest) (*RouteResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 查找符合条件的座席
	candidates := r.findCandidates(req)
	if len(candidates) == 0 {
		return nil, ErrNoAvailableAgent
	}

	// 对候选座席进行评分和排序
	scoredAgents := r.scoreAgents(candidates, req)
	sort.Slice(scoredAgents, func(i, j int) bool {
		return scoredAgents[i].score > scoredAgents[j].score
	})

	// 选择最佳座席
	bestAgent := scoredAgents[0]
	
	// 更新统计
	r.stats.TotalRouted++

	return &RouteResult{
		AgentID:      bestAgent.agent.ID,
		SkillMatched: bestAgent.matchedSkills,
		Score:        bestAgent.score,
		WaitTime:     0, // TODO: 实现预估等待时间
	}, nil
}

// findCandidates 查找候选座席
func (r *SkillBasedRouter) findCandidates(req *RouteRequest) []*Agent {
	candidates := make(map[uuid.UUID]*Agent)

	// 如果有技能要求，从技能索引中查找
	if len(req.RequiredSkills) > 0 {
		for _, skill := range req.RequiredSkills {
			if agentIDs, exists := r.skillIndex[skill]; exists {
				for _, agentID := range agentIDs {
					if agent, ok := r.agents[agentID]; ok && r.availableAgents[agentID] {
						candidates[agentID] = agent
					}
				}
			}
		}
	} else {
		// 没有技能要求，返回所有可用座席
		for agentID, agent := range r.agents {
			if r.availableAgents[agentID] {
				candidates[agentID] = agent
			}
		}
	}

	// 转换为切片
	result := make([]*Agent, 0, len(candidates))
	for _, agent := range candidates {
		result = append(result, agent)
	}

	return result
}

// scoreAgents 对座席进行评分
func (r *SkillBasedRouter) scoreAgents(agents []*Agent, req *RouteRequest) []scoredAgent {
	scored := make([]scoredAgent, len(agents))

	for i, agent := range agents {
		score := 0.0
		matchedSkills := []string{}

		// 技能匹配度（权重：40%）
		skillScore := 0.0
		for _, requiredSkill := range req.RequiredSkills {
			for _, agentSkill := range agent.Skills {
				if requiredSkill == agentSkill {
					skillScore += 1.0
					matchedSkills = append(matchedSkills, requiredSkill)
					break
				}
			}
		}
		if len(req.RequiredSkills) > 0 {
			skillScore = (skillScore / float64(len(req.RequiredSkills))) * 40
		}

		// 座席等级（权重：20%）
		levelScore := float64(agent.Level) * 4 // 假设等级1-5

		// 当前负载（权重：20%）
		loadScore := (1.0 - float64(agent.CurrentLoad)/10.0) * 20 // 假设最大负载10

		// 历史表现（权重：20%）
		performanceScore := agent.Performance * 20

		// VIP优先
		if req.VIP && agent.Level >= 4 {
			score += 10 // VIP加分
		}

		// 语言匹配
		if req.Language != "" {
			for _, lang := range agent.Languages {
				if lang == req.Language {
					score += 5
					break
				}
			}
		}

		score = skillScore + levelScore + loadScore + performanceScore
		scored[i] = scoredAgent{
			agent:         agent,
			score:         score,
			matchedSkills: matchedSkills,
		}
	}

	return scored
}

// RegisterAgent 注册座席
func (r *SkillBasedRouter) RegisterAgent(agent *Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.agents[agent.ID] = agent
	
	// 更新技能索引
	for _, skill := range agent.Skills {
		r.skillIndex[skill] = append(r.skillIndex[skill], agent.ID)
	}

	// 如果座席状态是可用，加入可用列表
	if agent.Status == AgentStatusAvailable {
		r.availableAgents[agent.ID] = true
		r.stats.AvailableAgents++
	}

	r.stats.TotalAgents++
	r.logger.Infof("Agent %s registered with skills %v", agent.ID, agent.Skills)

	return nil
}

// UnregisterAgent 注销座席
func (r *SkillBasedRouter) UnregisterAgent(agentID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[agentID]
	if !exists {
		return errors.New("agent not found")
	}

	// 从技能索引中移除
	for _, skill := range agent.Skills {
		r.removeAgentFromSkill(skill, agentID)
	}

	// 从可用列表中移除
	if r.availableAgents[agentID] {
		delete(r.availableAgents, agentID)
		r.stats.AvailableAgents--
	}

	delete(r.agents, agentID)
	r.stats.TotalAgents--

	r.logger.Infof("Agent %s unregistered", agentID)
	return nil
}

// UpdateAgentStatus 更新座席状态
func (r *SkillBasedRouter) UpdateAgentStatus(agentID uuid.UUID, status AgentStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[agentID]
	if !exists {
		return errors.New("agent not found")
	}

	oldStatus := agent.Status
	agent.Status = status

	// 更新可用性
	if oldStatus == AgentStatusAvailable && status != AgentStatusAvailable {
		delete(r.availableAgents, agentID)
		r.stats.AvailableAgents--
		r.stats.BusyAgents++
	} else if oldStatus != AgentStatusAvailable && status == AgentStatusAvailable {
		r.availableAgents[agentID] = true
		r.stats.AvailableAgents++
		r.stats.BusyAgents--
	}

	r.logger.Infof("Agent %s status updated from %s to %s", agentID, oldStatus, status)
	return nil
}

// GetRouterStats 获取路由器统计
func (r *SkillBasedRouter) GetRouterStats() *RouterStats {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	stats := *r.stats
	return &stats
}

// removeAgentFromSkill 从技能索引中移除座席
func (r *SkillBasedRouter) removeAgentFromSkill(skill string, agentID uuid.UUID) {
	agents := r.skillIndex[skill]
	for i, id := range agents {
		if id == agentID {
			r.skillIndex[skill] = append(agents[:i], agents[i+1:]...)
			break
		}
	}
}

// scoredAgent 评分后的座席
type scoredAgent struct {
	agent         *Agent
	score         float64
	matchedSkills []string
}

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	mu     sync.RWMutex
	agents []*Agent
	index  int
}

// RoundRobin 轮询分配
func (lb *LoadBalancer) RoundRobin() *Agent {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.agents) == 0 {
		return nil
	}

	agent := lb.agents[lb.index]
	lb.index = (lb.index + 1) % len(lb.agents)
	return agent
}

// LeastConnections 最少连接分配
func (lb *LoadBalancer) LeastConnections() *Agent {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.agents) == 0 {
		return nil
	}

	minLoad := lb.agents[0].CurrentLoad
	selectedAgent := lb.agents[0]

	for _, agent := range lb.agents[1:] {
		if agent.CurrentLoad < minLoad {
			minLoad = agent.CurrentLoad
			selectedAgent = agent
		}
	}

	return selectedAgent
}

// Random 随机分配
func (lb *LoadBalancer) Random() *Agent {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	if len(lb.agents) == 0 {
		return nil
	}

	return lb.agents[rand.Intn(len(lb.agents))]
}
