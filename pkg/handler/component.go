package handler

import (
	"github.com/bwmarrin/discordgo"
)

type Component interface {
	GetCommand() *discordgo.ApplicationCommand
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type ComponentHandler interface {
	AddComponent(component Component) error
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type componentHandler struct {
	components         []*discordgo.ApplicationCommand
	componentsHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func NewComponentHandler() ComponentHandler {
	handler := &componentHandler{}
	handler.components = make([]*discordgo.ApplicationCommand, 0)
	handler.componentsHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	return handler
}

func (h *componentHandler) AddComponent(component Component) error {
	h.components = append(h.components, component.GetCommand())
	h.componentsHandlers[component.GetCommand().Name] = component.Execute
	return nil
}

func (h *componentHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command, exist := h.componentsHandlers[i.MessageComponentData().CustomID]
	if !exist {
		return
	}
	command(s, i)
}
