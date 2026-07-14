package evolution

import "fmt"

// EvolutionClient is the entry point for interacting with the Evolution API.
type EvolutionClient interface {
	Config() EvolutionConfig

	InstanceService() InstanceService
	MessageService() MessageService
	ChatService() ChatService
	GroupService() GroupService
	LabelService() LabelService
	SettingsService() SettingsService
	ProxyService() ProxyService
	WebhookService() WebhookService
	TemplateService() TemplateService
	CallService() CallService
	BusinessService() BusinessService
}

type evolutionClient struct {
	config EvolutionConfig

	provider HttpProvider

	instanceService InstanceService
	messageService  MessageService
	chatService     ChatService
	groupService    GroupService
	labelService    LabelService
	settingsService SettingsService
	proxyService    ProxyService
	webhookService  WebhookService
	templateService TemplateService
	callService     CallService
	businessService BusinessService
}

func NewEvolutionClient(config EvolutionConfig) (EvolutionClient, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid Evolution configuration: %w", err)
	}

	httpProvider := NewHttpProvider(config.BaseURL)

	return evolutionClient{
		config:          config,
		provider:        httpProvider,
		instanceService: NewInstanceService(httpProvider, config.APIKey),
		messageService:  NewMessageService(httpProvider, config.APIKey),
		chatService:     NewChatService(httpProvider, config.APIKey),
		groupService:    NewGroupService(httpProvider, config.APIKey),
		labelService:    NewLabelService(httpProvider, config.APIKey),
		settingsService: NewSettingsService(httpProvider, config.APIKey),
		proxyService:    NewProxyService(httpProvider, config.APIKey),
		webhookService:  NewWebhookService(httpProvider, config.APIKey),
		templateService: NewTemplateService(httpProvider, config.APIKey),
		callService:     NewCallService(httpProvider, config.APIKey),
		businessService: NewBusinessService(httpProvider, config.APIKey),
	}, nil
}

func (c evolutionClient) Config() EvolutionConfig {
	return c.config
}

func (c evolutionClient) InstanceService() InstanceService {
	return c.instanceService
}

func (c evolutionClient) MessageService() MessageService {
	return c.messageService
}

func (c evolutionClient) ChatService() ChatService {
	return c.chatService
}

func (c evolutionClient) GroupService() GroupService {
	return c.groupService
}

func (c evolutionClient) LabelService() LabelService {
	return c.labelService
}

func (c evolutionClient) SettingsService() SettingsService {
	return c.settingsService
}

func (c evolutionClient) ProxyService() ProxyService {
	return c.proxyService
}

func (c evolutionClient) WebhookService() WebhookService {
	return c.webhookService
}

func (c evolutionClient) TemplateService() TemplateService {
	return c.templateService
}

func (c evolutionClient) CallService() CallService {
	return c.callService
}

func (c evolutionClient) BusinessService() BusinessService {
	return c.businessService
}
