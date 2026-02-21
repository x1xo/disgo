package handlers

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/rest"
)

func gatewayHandlerInteractionCreate(client *bot.Client, sequenceNumber int, shardID int, event gateway.EventInteractionCreate) {
	handleInteraction(client, sequenceNumber, shardID, nil, event.Interaction)
}

func respond(client *bot.Client, respondFunc httpserver.RespondFunc, interaction discord.Interaction, responseState *events.InteractionResponseState) events.InteractionResponderFunc {
	return func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opts ...rest.RequestOpt) error {
		if responseState.ResponseType() != nil {
			return discord.ErrInteractionAlreadyReplied
		}

		response := discord.InteractionResponse{
			Type: responseType,
			Data: data,
		}

		var err error
		if respondFunc != nil {
			err = respondFunc(response)
		} else {
			err = client.Rest.CreateInteractionResponse(interaction.ID(), interaction.Token(), response, opts...)
		}
		if err != nil {
			return err
		}

		responseState.Mu.Lock()
		defer responseState.Mu.Unlock()
		if responseState.ResponseTypeValue != nil {
			return discord.ErrInteractionAlreadyReplied
		}
		responseTypeCopy := responseType
		responseState.ResponseTypeValue = &responseTypeCopy
		return err
	}
}

func handleInteraction(client *bot.Client, sequenceNumber int, shardID int, respondFunc httpserver.RespondFunc, interaction discord.Interaction) {
	genericEvent := events.NewGenericEvent(client, sequenceNumber, shardID)
	responseState := events.NewInteractionResponseState()
	responder := respond(client, respondFunc, interaction, responseState)

	client.EventManager.DispatchEvent(&events.InteractionCreate{
		GenericEvent:  genericEvent,
		Interaction:   interaction,
		Respond:       responder,
		ResponseState: responseState,
	})

	switch i := interaction.(type) {
	case discord.ApplicationCommandInteraction:
		client.EventManager.DispatchEvent(&events.ApplicationCommandInteractionCreate{
			GenericEvent:                  genericEvent,
			ApplicationCommandInteraction: i,
			Respond:                       responder,
			ResponseState:                 responseState,
		})

	case discord.ComponentInteraction:
		client.EventManager.DispatchEvent(&events.ComponentInteractionCreate{
			GenericEvent:         genericEvent,
			ComponentInteraction: i,
			Respond:              responder,
			ResponseState:        responseState,
		})

	case discord.AutocompleteInteraction:
		client.EventManager.DispatchEvent(&events.AutocompleteInteractionCreate{
			GenericEvent:            genericEvent,
			AutocompleteInteraction: i,
			Respond:                 responder,
			ResponseState:           responseState,
		})

	case discord.ModalSubmitInteraction:
		client.EventManager.DispatchEvent(&events.ModalSubmitInteractionCreate{
			GenericEvent:           genericEvent,
			ModalSubmitInteraction: i,
			Respond:                responder,
			ResponseState:          responseState,
		})

	default:
		client.Logger.Error("unknown interaction", slog.String("type", fmt.Sprintf("%T", interaction)))
	}
}
