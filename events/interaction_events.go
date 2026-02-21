package events

import (
	"sync"

	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
)

// InteractionResponderFunc is a function that can be used to respond to a discord.Interaction.
type InteractionResponderFunc func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error

// InteractionResponseState is a thread-safe state tracker for interaction responses.
type InteractionResponseState struct {
	Mu                sync.RWMutex
	ResponseTypeValue *discord.InteractionResponseType
}

func NewInteractionResponseState() *InteractionResponseState {
	return &InteractionResponseState{}
}

func (s *InteractionResponseState) ResponseType() *discord.InteractionResponseType {
	if s == nil {
		return nil
	}
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if s.ResponseTypeValue == nil {
		return nil
	}
	responseType := *s.ResponseTypeValue
	return &responseType
}

func (s *InteractionResponseState) IsDefferedReply() bool {
	responseType := s.ResponseType()
	return responseType != nil && *responseType == discord.InteractionResponseTypeDeferredCreateMessage
}

func (s *InteractionResponseState) IsDefferedUpdate() bool {
	responseType := s.ResponseType()
	return responseType != nil && *responseType == discord.InteractionResponseTypeDeferredUpdateMessage
}

// InteractionCreate indicates that a new interaction has been created.
type InteractionCreate struct {
	*GenericEvent
	discord.Interaction
	Respond       InteractionResponderFunc
	ResponseState *InteractionResponseState
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *InteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches.Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *InteractionCreate) ResponseType() *discord.InteractionResponseType {
	return e.ResponseState.ResponseType()
}

func (e *InteractionCreate) IsDefferedReply() bool {
	return e.ResponseState.IsDefferedReply()
}

func (e *InteractionCreate) IsDefferedUpdate() bool {
	return e.ResponseState.IsDefferedUpdate()
}

func (e *InteractionCreate) GetInteractionResponse(opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.GetInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *InteractionCreate) UpdateInteractionResponse(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageUpdate, opts...)
}

func (e *InteractionCreate) DeleteInteractionResponse(opts ...rest.RequestOpt) error {
	return e.Client().Rest.DeleteInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *InteractionCreate) GetFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.GetFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

func (e *InteractionCreate) CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), messageCreate, opts...)
}

func (e *InteractionCreate) UpdateFollowupMessage(messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.UpdateFollowupMessage(e.ApplicationID(), e.Token(), messageID, messageUpdate, opts...)
}

func (e *InteractionCreate) DeleteFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return e.Client().Rest.DeleteFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

// ApplicationCommandInteractionCreate is the base struct for all application command interaction create events.
type ApplicationCommandInteractionCreate struct {
	*GenericEvent
	discord.ApplicationCommandInteraction
	Respond       InteractionResponderFunc
	ResponseState *InteractionResponseState
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ApplicationCommandInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches.Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *ApplicationCommandInteractionCreate) ResponseType() *discord.InteractionResponseType {
	return e.ResponseState.ResponseType()
}

func (e *ApplicationCommandInteractionCreate) IsDefferedReply() bool {
	return e.ResponseState.IsDefferedReply()
}

func (e *ApplicationCommandInteractionCreate) IsDefferedUpdate() bool {
	return e.ResponseState.IsDefferedUpdate()
}

func (e *ApplicationCommandInteractionCreate) GetInteractionResponse(opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.GetInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *ApplicationCommandInteractionCreate) UpdateInteractionResponse(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageUpdate, opts...)
}

func (e *ApplicationCommandInteractionCreate) DeleteInteractionResponse(opts ...rest.RequestOpt) error {
	return e.Client().Rest.DeleteInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *ApplicationCommandInteractionCreate) GetFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.GetFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

func (e *ApplicationCommandInteractionCreate) CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), messageCreate, opts...)
}

func (e *ApplicationCommandInteractionCreate) UpdateFollowupMessage(messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.UpdateFollowupMessage(e.ApplicationID(), e.Token(), messageID, messageUpdate, opts...)
}

func (e *ApplicationCommandInteractionCreate) DeleteFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return e.Client().Rest.DeleteFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

// Acknowledge acknowledges the interaction.
//
// This is used strictly for acknowledging the HTTP interaction request from discord. This responds with 202 Accepted.
//
// When using this, your first http request must be [rest.Interactions.CreateInteractionResponse] or [rest.Interactions.CreateInteractionResponseWithCallback]
//
// This does not produce a visible loading state to the user.
// You are expected to send a new http request within 3 seconds to respond to the interaction.
// This allows you to gracefully handle errors with your sent response & access the resulting message.
//
// If you want to create a visible loading state, use DeferCreateMessage.
//
// Source docs: [Discord Source docs]
//
// [Discord Source docs]: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
func (e *ApplicationCommandInteractionCreate) Acknowledge(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAcknowledge, nil, opts...)
}

// CreateMessage responds to the interaction with a new message.
func (e *ApplicationCommandInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *ApplicationCommandInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// Modal responds to the interaction with a new modal.
func (e *ApplicationCommandInteractionCreate) Modal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeModal, modalCreate, opts...)
}

// LaunchActivity responds to the interaction by launching activity associated with the app.
func (e *ApplicationCommandInteractionCreate) LaunchActivity(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeLaunchActivity, nil, opts...)
}

// ComponentInteractionCreate indicates that a new component interaction has been created.
type ComponentInteractionCreate struct {
	*GenericEvent
	discord.ComponentInteraction
	Respond       InteractionResponderFunc
	ResponseState *InteractionResponseState
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ComponentInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches.Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *ComponentInteractionCreate) ResponseType() *discord.InteractionResponseType {
	return e.ResponseState.ResponseType()
}

func (e *ComponentInteractionCreate) IsDefferedReply() bool {
	return e.ResponseState.IsDefferedReply()
}

func (e *ComponentInteractionCreate) IsDefferedUpdate() bool {
	return e.ResponseState.IsDefferedUpdate()
}

func (e *ComponentInteractionCreate) GetInteractionResponse(opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.GetInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *ComponentInteractionCreate) UpdateInteractionResponse(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageUpdate, opts...)
}

func (e *ComponentInteractionCreate) DeleteInteractionResponse(opts ...rest.RequestOpt) error {
	return e.Client().Rest.DeleteInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *ComponentInteractionCreate) GetFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.GetFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

func (e *ComponentInteractionCreate) CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), messageCreate, opts...)
}

func (e *ComponentInteractionCreate) UpdateFollowupMessage(messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.UpdateFollowupMessage(e.ApplicationID(), e.Token(), messageID, messageUpdate, opts...)
}

func (e *ComponentInteractionCreate) DeleteFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return e.Client().Rest.DeleteFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

// Acknowledge acknowledges the interaction.
//
// This is used strictly for acknowledging the HTTP interaction request from discord. This responds with 202 Accepted.
//
// When using this, your first http request must be [rest.Interactions.CreateInteractionResponse] or [rest.Interactions.CreateInteractionResponseWithCallback]
//
// This does not produce a visible loading state to the user.
// You are expected to send a new http request within 3 seconds to respond to the interaction.
// This allows you to gracefully handle errors with your sent response & access the resulting message.
//
// If you want to create a visible loading state, use DeferCreateMessage.
//
// Source docs: [Discord Source docs]
//
// [Discord Source docs]: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
func (e *ComponentInteractionCreate) Acknowledge(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAcknowledge, nil, opts...)
}

// CreateMessage responds to the interaction with a new message.
func (e *ComponentInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *ComponentInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// UpdateMessage responds to the interaction with updating the message the component is from.
func (e *ComponentInteractionCreate) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// DeferUpdateMessage responds to the interaction with nothing.
func (e *ComponentInteractionCreate) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}

// Modal responds to the interaction with a new modal.
func (e *ComponentInteractionCreate) Modal(modalCreate discord.ModalCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeModal, modalCreate, opts...)
}

// LaunchActivity responds to the interaction by launching activity associated with the app.
func (e *ComponentInteractionCreate) LaunchActivity(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeLaunchActivity, nil, opts...)
}

// AutocompleteInteractionCreate indicates that a new autocomplete interaction has been created.
type AutocompleteInteractionCreate struct {
	*GenericEvent
	discord.AutocompleteInteraction
	Respond       InteractionResponderFunc
	ResponseState *InteractionResponseState
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *AutocompleteInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches.Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *AutocompleteInteractionCreate) ResponseType() *discord.InteractionResponseType {
	return e.ResponseState.ResponseType()
}

func (e *AutocompleteInteractionCreate) IsDefferedReply() bool {
	return e.ResponseState.IsDefferedReply()
}

func (e *AutocompleteInteractionCreate) IsDefferedUpdate() bool {
	return e.ResponseState.IsDefferedUpdate()
}

// Acknowledge acknowledges the interaction.
//
// This is used strictly for acknowledging the HTTP interaction request from discord. This responds with 202 Accepted.
//
// When using this, your first http request must be [rest.Interactions.CreateInteractionResponse] or [rest.Interactions.CreateInteractionResponseWithCallback]
//
// This does not produce a visible loading state to the user.
// You are expected to send a new http request within 3 seconds to respond to the interaction.
// This allows you to gracefully handle errors with your sent response & access the resulting message.
//
// If you want to create a visible loading state, use DeferCreateMessage.
//
// Source docs: [Discord Source docs]
//
// [Discord Source docs]: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
func (e *AutocompleteInteractionCreate) Acknowledge(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAcknowledge, nil, opts...)
}

// AutocompleteResult responds to the interaction with a slice of choices.
func (e *AutocompleteInteractionCreate) AutocompleteResult(choices []discord.AutocompleteChoice, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAutocompleteResult, discord.AutocompleteResult{Choices: choices}, opts...)
}

// ModalSubmitInteractionCreate indicates that a new modal submit interaction has been created.
type ModalSubmitInteractionCreate struct {
	*GenericEvent
	discord.ModalSubmitInteraction
	Respond       InteractionResponderFunc
	ResponseState *InteractionResponseState
}

// Guild returns the guild that the interaction happened in if it happened in a guild.
// If the interaction happened in a DM, it returns nil.
// This only returns cached guilds.
func (e *ModalSubmitInteractionCreate) Guild() (discord.Guild, bool) {
	if e.GuildID() != nil {
		return e.Client().Caches.Guild(*e.GuildID())
	}
	return discord.Guild{}, false
}

func (e *ModalSubmitInteractionCreate) ResponseType() *discord.InteractionResponseType {
	return e.ResponseState.ResponseType()
}

func (e *ModalSubmitInteractionCreate) IsDefferedReply() bool {
	return e.ResponseState.IsDefferedReply()
}

func (e *ModalSubmitInteractionCreate) IsDefferedUpdate() bool {
	return e.ResponseState.IsDefferedUpdate()
}

func (e *ModalSubmitInteractionCreate) GetInteractionResponse(opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.GetInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *ModalSubmitInteractionCreate) UpdateInteractionResponse(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), messageUpdate, opts...)
}

func (e *ModalSubmitInteractionCreate) DeleteInteractionResponse(opts ...rest.RequestOpt) error {
	return e.Client().Rest.DeleteInteractionResponse(e.ApplicationID(), e.Token(), opts...)
}

func (e *ModalSubmitInteractionCreate) GetFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.GetFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

func (e *ModalSubmitInteractionCreate) CreateFollowupMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), messageCreate, opts...)
}

func (e *ModalSubmitInteractionCreate) UpdateFollowupMessage(messageID snowflake.ID, messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return e.Client().Rest.UpdateFollowupMessage(e.ApplicationID(), e.Token(), messageID, messageUpdate, opts...)
}

func (e *ModalSubmitInteractionCreate) DeleteFollowupMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return e.Client().Rest.DeleteFollowupMessage(e.ApplicationID(), e.Token(), messageID, opts...)
}

// Acknowledge acknowledges the interaction.
//
// This is used strictly for acknowledging the HTTP interaction request from discord. This responds with 202 Accepted.
//
// When using this, your first http request must be [rest.Interactions.CreateInteractionResponse] or [rest.Interactions.CreateInteractionResponseWithCallback]
//
// This does not produce a visible loading state to the user.
// You are expected to send a new http request within 3 seconds to respond to the interaction.
// This allows you to gracefully handle errors with your sent response & access the resulting message.
//
// If you want to create a visible loading state, use DeferCreateMessage.
//
// Source docs: [Discord Source docs]
//
// [Discord Source docs]: https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
func (e *ModalSubmitInteractionCreate) Acknowledge(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeAcknowledge, nil, opts...)
}

// CreateMessage responds to the interaction with a new message.
func (e *ModalSubmitInteractionCreate) CreateMessage(messageCreate discord.MessageCreate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeCreateMessage, messageCreate, opts...)
}

// DeferCreateMessage responds to the interaction with a "bot is thinking..." message which should be edited later.
func (e *ModalSubmitInteractionCreate) DeferCreateMessage(ephemeral bool, opts ...rest.RequestOpt) error {
	var data discord.InteractionResponseData
	if ephemeral {
		data = discord.MessageCreate{Flags: discord.MessageFlagEphemeral}
	}
	return e.Respond(discord.InteractionResponseTypeDeferredCreateMessage, data, opts...)
}

// UpdateMessage responds to the interaction with updating the message the component is from.
func (e *ModalSubmitInteractionCreate) UpdateMessage(messageUpdate discord.MessageUpdate, opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeUpdateMessage, messageUpdate, opts...)
}

// DeferUpdateMessage responds to the interaction with nothing.
func (e *ModalSubmitInteractionCreate) DeferUpdateMessage(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeDeferredUpdateMessage, nil, opts...)
}

// LaunchActivity responds to the interaction by launching activity associated with the app.
func (e *ModalSubmitInteractionCreate) LaunchActivity(opts ...rest.RequestOpt) error {
	return e.Respond(discord.InteractionResponseTypeLaunchActivity, nil, opts...)
}
