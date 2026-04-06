package v1

import (
	"net/http"
	"strconv"

	"github.com/anthropics/firefly-iii-go/internal/adapter/transformer"
	"github.com/anthropics/firefly-iii-go/internal/entity"
	ruleuc "github.com/anthropics/firefly-iii-go/internal/usecase/rule"
	"github.com/anthropics/firefly-iii-go/pkg/pagination"
	"github.com/anthropics/firefly-iii-go/pkg/response"
	"github.com/gin-gonic/gin"
)

type RuleHandler struct {
	uc *ruleuc.UseCase
}

func NewRuleHandler(uc *ruleuc.UseCase) *RuleHandler {
	return &RuleHandler{uc: uc}
}

// --- RuleGroup endpoints ---

func (h *RuleHandler) RuleGroupIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 1000 {
		limit = 50
	}
	offset := (page - 1) * limit

	userGroupID := uint(0) // TODO: get from user context in later SP

	groups, total, err := h.uc.ListGroups(c.Request.Context(), userGroupID, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(groups))
	for i, rg := range groups {
		items[i] = transformer.TransformRuleGroup(&rg)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *RuleHandler) RuleGroupShow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid rule group ID")
		return
	}

	rg, err := h.uc.GetGroupByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Rule group not found")
		return
	}

	resource := transformer.TransformRuleGroup(rg)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type storeRuleGroupRequest struct {
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description"`
	Order          uint   `json:"order"`
	Active         *bool  `json:"active"`
	StopProcessing bool   `json:"stop_processing"`
}

func (h *RuleHandler) RuleGroupStore(c *gin.Context) {
	var req storeRuleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	rg := &entity.RuleGroup{
		Title:          req.Title,
		Description:    req.Description,
		Order:          req.Order,
		Active:         active,
		StopProcessing: req.StopProcessing,
	}

	if err := h.uc.CreateGroup(c.Request.Context(), rg); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resource := transformer.TransformRuleGroup(rg)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateRuleGroupRequest struct {
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description"`
	Order          uint   `json:"order"`
	Active         *bool  `json:"active"`
	StopProcessing bool   `json:"stop_processing"`
}

func (h *RuleHandler) RuleGroupUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid rule group ID")
		return
	}

	rg, err := h.uc.GetGroupByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Rule group not found")
		return
	}

	var req updateRuleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rg.Title = req.Title
	rg.Description = req.Description
	rg.Order = req.Order
	rg.StopProcessing = req.StopProcessing
	if req.Active != nil {
		rg.Active = *req.Active
	}

	if err := h.uc.UpdateGroup(c.Request.Context(), rg); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resource := transformer.TransformRuleGroup(rg)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *RuleHandler) RuleGroupDestroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid rule group ID")
		return
	}
	if err := h.uc.DeleteGroup(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Rule group not found")
		return
	}
	response.NoContent(c)
}

func (h *RuleHandler) RuleGroupListRules(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid rule group ID")
		return
	}

	rules, err := h.uc.ListGroupRules(c.Request.Context(), uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(rules))
	for i, r := range rules {
		triggers, _ := h.uc.GetTriggers(c.Request.Context(), r.ID)
		actions, _ := h.uc.GetActions(c.Request.Context(), r.ID)
		items[i] = transformer.TransformRule(&r, triggers, actions)
	}

	pg := pagination.NewMeta(len(rules), len(rules), 1)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

// --- Rule endpoints ---

func (h *RuleHandler) RuleIndex(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 1000 {
		limit = 50
	}
	offset := (page - 1) * limit

	userGroupID := uint(0) // TODO: get from user context in later SP

	rules, total, err := h.uc.ListRules(c.Request.Context(), userGroupID, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	items := make([]response.Resource, len(rules))
	for i, r := range rules {
		triggers, _ := h.uc.GetTriggers(c.Request.Context(), r.ID)
		actions, _ := h.uc.GetActions(c.Request.Context(), r.ID)
		items[i] = transformer.TransformRule(&r, triggers, actions)
	}

	pg := pagination.NewMeta(int(total), limit, page)
	response.JSONApi(c, http.StatusOK, response.Collection(items, pg))
}

func (h *RuleHandler) RuleShow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid rule ID")
		return
	}

	r, err := h.uc.GetRuleByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Rule not found")
		return
	}

	triggers, _ := h.uc.GetTriggers(c.Request.Context(), r.ID)
	actions, _ := h.uc.GetActions(c.Request.Context(), r.ID)
	resource := transformer.TransformRule(r, triggers, actions)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type ruleTriggerInput struct {
	Type           string `json:"type" binding:"required"`
	Value          string `json:"value"`
	Order          uint   `json:"order"`
	Active         *bool  `json:"active"`
	StopProcessing bool   `json:"stop_processing"`
}

type ruleActionInput struct {
	Type           string `json:"type" binding:"required"`
	Value          string `json:"value"`
	Order          uint   `json:"order"`
	Active         *bool  `json:"active"`
	StopProcessing bool   `json:"stop_processing"`
}

type storeRuleRequest struct {
	Title          string             `json:"title" binding:"required"`
	Description    string             `json:"description"`
	RuleGroupID    uint               `json:"rule_group_id" binding:"required"`
	Order          uint               `json:"order"`
	Active         *bool              `json:"active"`
	Strict         bool               `json:"strict"`
	StopProcessing bool               `json:"stop_processing"`
	Triggers       []ruleTriggerInput `json:"triggers"`
	Actions        []ruleActionInput  `json:"actions"`
}

func (h *RuleHandler) RuleStore(c *gin.Context) {
	var req storeRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	active := true
	if req.Active != nil {
		active = *req.Active
	}

	r := &entity.Rule{
		Title:          req.Title,
		Description:    req.Description,
		RuleGroupID:    req.RuleGroupID,
		Order:          req.Order,
		Active:         active,
		Strict:         req.Strict,
		StopProcessing: req.StopProcessing,
	}

	if err := h.uc.CreateRule(c.Request.Context(), r); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Set triggers
	if len(req.Triggers) > 0 {
		triggers := make([]entity.RuleTrigger, len(req.Triggers))
		for i, t := range req.Triggers {
			triggerActive := true
			if t.Active != nil {
				triggerActive = *t.Active
			}
			triggers[i] = entity.RuleTrigger{
				RuleID:         r.ID,
				TriggerType:    t.Type,
				TriggerValue:   t.Value,
				Order:          t.Order,
				Active:         triggerActive,
				StopProcessing: t.StopProcessing,
			}
		}
		_ = h.uc.SetTriggers(c.Request.Context(), r.ID, triggers)
	}

	// Set actions
	if len(req.Actions) > 0 {
		actions := make([]entity.RuleAction, len(req.Actions))
		for i, a := range req.Actions {
			actionActive := true
			if a.Active != nil {
				actionActive = *a.Active
			}
			actions[i] = entity.RuleAction{
				RuleID:         r.ID,
				ActionType:     a.Type,
				ActionValue:    a.Value,
				Order:          a.Order,
				Active:         actionActive,
				StopProcessing: a.StopProcessing,
			}
		}
		_ = h.uc.SetActions(c.Request.Context(), r.ID, actions)
	}

	triggers, _ := h.uc.GetTriggers(c.Request.Context(), r.ID)
	actions, _ := h.uc.GetActions(c.Request.Context(), r.ID)
	resource := transformer.TransformRule(r, triggers, actions)
	response.JSONApi(c, http.StatusCreated, response.Single(resource.Type, resource.ID, resource.Attributes))
}

type updateRuleRequest struct {
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description"`
	RuleGroupID    uint   `json:"rule_group_id"`
	Order          uint   `json:"order"`
	Active         *bool  `json:"active"`
	Strict         bool   `json:"strict"`
	StopProcessing bool   `json:"stop_processing"`
}

func (h *RuleHandler) RuleUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid rule ID")
		return
	}

	r, err := h.uc.GetRuleByID(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "Rule not found")
		return
	}

	var req updateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	r.Title = req.Title
	r.Description = req.Description
	r.Order = req.Order
	r.Strict = req.Strict
	r.StopProcessing = req.StopProcessing
	if req.RuleGroupID > 0 {
		r.RuleGroupID = req.RuleGroupID
	}
	if req.Active != nil {
		r.Active = *req.Active
	}

	if err := h.uc.UpdateRule(c.Request.Context(), r); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	triggers, _ := h.uc.GetTriggers(c.Request.Context(), r.ID)
	actions, _ := h.uc.GetActions(c.Request.Context(), r.ID)
	resource := transformer.TransformRule(r, triggers, actions)
	response.JSONApi(c, http.StatusOK, response.Single(resource.Type, resource.ID, resource.Attributes))
}

func (h *RuleHandler) RuleDestroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid rule ID")
		return
	}
	if err := h.uc.DeleteRule(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, "Rule not found")
		return
	}
	response.NoContent(c)
}
